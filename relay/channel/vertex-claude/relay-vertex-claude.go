package vertex_claude

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"one-api/common"
	"one-api/dto"
	relaymodel "one-api/dto"
	"one-api/relay/channel/claude"
	"one-api/service"
	"strings"
	"sync"
	"time"
)

var accessTokenMap sync.Map

func getAccessToken(json string) (string, error) {
	data, ok := accessTokenMap.Load(json)
	if ok {
		token := data.(oauth2.Token)
		if time.Now().Before(token.Expiry) {
			return token.AccessToken, nil
		}
	}
	creds, err := google.CredentialsFromJSON(context.Background(), []byte(json), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return "", err
	}
	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", err
	}
	accessTokenMap.Store(json, *token)
	return token.AccessToken, nil
}

func getRedirectModel(requestModel string) (string, error) {
	if model, ok := modelIdMap[requestModel]; ok {
		return model, nil
	}
	return "", errors.Errorf("model %s not found", requestModel)
}

func requestOpenAI2VertexClaude(request dto.GeneralOpenAIRequest) (*VertexClaudeRequest, error) {
	vertexClaudeRequest := VertexClaudeRequest{
		AnthropicVersion: "vertex-2023-10-16",
		Stream:           request.Stream,
	}
	if vertexClaudeRequest.MaxTokens == 0 {
		vertexClaudeRequest.MaxTokens = 4096
	}
	formatMessages := make([]dto.Message, 0)
	var lastMessage *dto.Message
	for i, message := range request.Messages {
		if message.Role == "" {
			request.Messages[i].Role = "user"
		}
		fmtMessage := dto.Message{
			Role:    message.Role,
			Content: message.Content,
		}
		if lastMessage != nil && lastMessage.Role == message.Role {
			if lastMessage.IsStringContent() && message.IsStringContent() {
				content, _ := json.Marshal(strings.Trim(fmt.Sprintf("%s %s", lastMessage.StringContent(), message.StringContent()), "\""))
				fmtMessage.Content = content
				// delete last message
				formatMessages = formatMessages[:len(formatMessages)-1]
			}
		}
		if fmtMessage.Content == nil {
			content, _ := json.Marshal("...")
			fmtMessage.Content = content
		}
		formatMessages = append(formatMessages, fmtMessage)
		lastMessage = &request.Messages[i]
	}
	claudeMessages := make([]claude.ClaudeMessage, 0)
	for _, message := range formatMessages {
		if message.Role == "system" {
			if message.IsStringContent() {
				vertexClaudeRequest.System = message.StringContent()
			} else {
				contents := message.ParseContent()
				content := ""
				for _, ctx := range contents {
					if ctx.Type == "text" {
						content += ctx.Text
					}
				}
				vertexClaudeRequest.System = content
			}
		} else {
			claudeMessage := claude.ClaudeMessage{
				Role: message.Role,
			}
			if message.IsStringContent() {
				claudeMessage.Content = message.StringContent()
			} else {
				claudeMediaMessages := make([]claude.ClaudeMediaMessage, 0)
				for _, mediaMessage := range message.ParseContent() {
					claudeMediaMessage := claude.ClaudeMediaMessage{
						Type: mediaMessage.Type,
					}
					if mediaMessage.Type == "text" {
						claudeMediaMessage.Text = mediaMessage.Text
					} else {
						imageUrl := mediaMessage.ImageUrl.(dto.MessageImageUrl)
						claudeMediaMessage.Type = "image"
						claudeMediaMessage.Source = &claude.ClaudeMessageSource{
							Type: "base64",
						}
						// 判断是否是url
						if strings.HasPrefix(imageUrl.Url, "http") {
							// 是url，获取图片的类型和base64编码的数据
							mimeType, data, _ := common.GetImageFromUrl(imageUrl.Url)
							claudeMediaMessage.Source.MediaType = mimeType
							claudeMediaMessage.Source.Data = data
						} else {
							_, format, base64String, err := common.DecodeBase64ImageData(imageUrl.Url)
							if err != nil {
								return nil, err
							}
							claudeMediaMessage.Source.MediaType = "image/" + format
							claudeMediaMessage.Source.Data = base64String
						}
					}
					claudeMediaMessages = append(claudeMediaMessages, claudeMediaMessage)
				}
				claudeMessage.Content = claudeMediaMessages
			}
			claudeMessages = append(claudeMessages, claudeMessage)
		}
	}
	vertexClaudeRequest.Messages = claudeMessages
	return &vertexClaudeRequest, nil
}

func vertexClaudeHandler(c *gin.Context, resp *http.Response) (*dto.OpenAIErrorWithStatusCode, *dto.Usage) {
	var claudeResponse claude.ClaudeResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	err = json.Unmarshal(responseBody, &claudeResponse)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil
	}
	openaiResp := claude.ResponseClaude2OpenAI(claude.RequestModeMessage, &claudeResponse)
	usage := relaymodel.Usage{
		PromptTokens:     claudeResponse.Usage.InputTokens,
		CompletionTokens: claudeResponse.Usage.OutputTokens,
		TotalTokens:      claudeResponse.Usage.InputTokens + claudeResponse.Usage.OutputTokens,
	}
	openaiResp.Usage = usage
	c.JSON(http.StatusOK, openaiResp)
	return nil, &usage
}

func vertexClaudeStreamHandler(c *gin.Context, resp *http.Response) (*relaymodel.OpenAIErrorWithStatusCode, *relaymodel.Usage) {
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := strings.Index(string(data), "\n"); i >= 0 {
			return i + 1, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	dataChan := make(chan string)
	stopChan := make(chan bool)
	go func() {
		for scanner.Scan() {
			data := scanner.Text()
			if len(data) < 5 { // ignore blank line or wrong format
				continue
			}
			if data[:5] != "data:" {
				continue
			}
			data = data[5:]
			dataChan <- data
		}
		stopChan <- true
	}()
	var id string
	var model string
	createdTime := common.GetTimestamp()
	var usage relaymodel.Usage
	service.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-dataChan:
			claudeResp := new(claude.ClaudeResponse)
			err := json.Unmarshal([]byte(data), &claudeResp)
			if err != nil {
				common.SysError("error unmarshalling stream response: " + err.Error())
				return true
			}
			response, claudeUsage := claude.StreamResponseClaude2OpenAI(claude.RequestModeMessage, claudeResp)

			if claudeUsage != nil {
				usage.PromptTokens += claudeUsage.InputTokens
				usage.CompletionTokens += claudeUsage.OutputTokens
			}

			if response == nil {
				return true
			}

			if response.Id != "" {
				id = response.Id
			}
			if response.Model != "" {
				model = response.Model
			}
			response.Created = createdTime
			response.Id = id
			response.Model = model

			jsonStr, err := json.Marshal(response)
			if err != nil {
				common.SysError("error marshalling stream response: " + err.Error())
				return true
			}
			c.Render(-1, common.CustomEvent{Data: "data: " + string(jsonStr)})
			return true
		case <-stopChan:
			c.Render(-1, common.CustomEvent{Data: "data: [DONE]"})
			return false
		}
	})
	err := resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	return nil, &usage
}
