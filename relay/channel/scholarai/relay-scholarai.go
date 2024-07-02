package scholarai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"one-api/common"
	"one-api/dto"
	relaycommon "one-api/relay/common"
	"one-api/service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func requestOpenAI2ScholarAI(request dto.GeneralOpenAIRequest) *ScholarAIChatRequest {
	var combinedMessage strings.Builder
	for _, message := range request.Messages {
		combinedMessage.WriteString(fmt.Sprintf("%s: %s\n", message.Role, message.StringContent()))
	}
	scholarAIMessage := ScholarAIMessage{
		Role:    "user",
		Content: combinedMessage.String(),
	}
	return &ScholarAIChatRequest{
		Model:    request.Model,
		Messages: []ScholarAIMessage{scholarAIMessage},
		Stream:   request.Stream,
	}
}

func responseScholarAI2OpenAI(response *ScholarAITextResponse) *dto.OpenAITextResponse {
	fullTextResponse := dto.OpenAITextResponse{
		Id:      response.Id,
		Object:  response.Object,
		Created: response.Created,
	}
	for _, choice := range response.Choices {
		content, _ := json.Marshal(choice.ScholarAIMessage.Content)
		c := dto.OpenAITextResponseChoice{
			Index:        choice.Index,
			FinishReason: choice.FinishReason,
			Message: dto.Message{
				Content: content,
				Role:    choice.ScholarAIMessage.Role,
			},
		}
		fullTextResponse.Choices = append(fullTextResponse.Choices, c)
	}

	return &fullTextResponse
}

func scholarAIHandler(c *gin.Context, resp *http.Response, promptTokens int, model string) (*dto.OpenAIErrorWithStatusCode, *dto.Usage) {
	var ScholarAITextResponse ScholarAITextResponse
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil
	}
	err = resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil
	}
	err = json.Unmarshal(responseBody, &ScholarAITextResponse)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil
	}
	fullTextResponse := responseScholarAI2OpenAI(&ScholarAITextResponse)
	completionTokens, _ := service.CountTokenText(ScholarAITextResponse.Choices[0].ScholarAIMessage.Content, model)
	usage := dto.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
	}
	fullTextResponse.Usage = usage
	jsonResponse, err := json.Marshal(fullTextResponse)
	if err != nil {
		return service.OpenAIErrorWrapper(err, "marshal_response_body_failed", http.StatusInternalServerError), nil
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	c.Writer.Write(jsonResponse)
	return nil, &fullTextResponse.Usage
}

func scholarAIStreamHandler(c *gin.Context, resp *http.Response, info *relaycommon.RelayInfo) (*dto.OpenAIErrorWithStatusCode, string) {
	var responseText string
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
			dataChan <- data
		}
		stopChan <- true
	}()
	service.SetEventStreamHeaders(c)
	isFirst := true
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-dataChan:
			if isFirst {
				isFirst = false
				info.FirstResponseTime = time.Now()
			}
			var response dto.ChatCompletionsStreamResponse
			err := json.Unmarshal([]byte(data), &response)
			if err != nil {
				common.SysError("error unmarshalling stream response: " + err.Error())
				return true
			}
			if len(response.Choices) != 0 {
				responseText += response.Choices[0].Delta.GetContentString()
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				common.SysError("error marshalling stream response: " + err.Error())
				return true
			}
			c.Render(-1, common.CustomEvent{Data: "data: " + string(jsonResponse)})
			return true
		case <-stopChan:
			c.Render(-1, common.CustomEvent{Data: "data: [DONE]"})
			return false
		}
	})
	err := resp.Body.Close()
	if err != nil {
		return service.OpenAIErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), ""
	}
	return nil, responseText
}
