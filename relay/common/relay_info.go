package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"one-api/common"
	"one-api/relay/constant"
	"strings"
	"time"
)

type RelayInfo struct {
	ChannelType          int
	ChannelId            int
	TokenId              int
	UserId               int
	Group                string
	TokenUnlimited       bool
	StartTime            time.Time
	FirstResponseTime    time.Time
	setFirstResponse     bool
	ApiType              int
	IsStream             bool
	RelayMode            int
	RequestModelName     string
	UpstreamModelName    string
	RequestURLPath       string
	ApiVersion           string
	PromptTokens         int
	ApiKey               string
	Organization         string
	BaseUrl              string
	SupportStreamOptions bool
	ShouldIncludeUsage   bool
	Proxy                string
}

func GenRelayInfo(c *gin.Context) (*RelayInfo, error) {
	channelType := c.GetInt("channel")
	channelId := c.GetInt("channel_id")

	tokenId := c.GetInt("token_id")
	userId := c.GetInt("id")
	group := c.GetString("group")
	tokenUnlimited := c.GetBool("token_unlimited_quota")
	startTime := time.Now()
	// firstResponseTime = time.Now() - 1 second

	apiType, _ := constant.ChannelType2APIType(channelType)

	info := &RelayInfo{
		RelayMode:         constant.Path2RelayMode(c.Request.URL.Path),
		BaseUrl:           c.GetString("base_url"),
		RequestURLPath:    c.Request.URL.String(),
		ChannelType:       channelType,
		ChannelId:         channelId,
		TokenId:           tokenId,
		UserId:            userId,
		Group:             group,
		TokenUnlimited:    tokenUnlimited,
		StartTime:         startTime,
		FirstResponseTime: startTime.Add(-time.Second),
		ApiType:           apiType,
		ApiVersion:        c.GetString("api_version"),
		ApiKey:            strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "),
		Organization:      c.GetString("channel_organization"),
		Proxy:             c.GetString("proxy"),
	}
	if info.BaseUrl == "" {
		ch, exists := common.ChannelMap[channelType]
		if !exists {
			return nil, errors.New("channel not exists")
		}
		info.BaseUrl = ch.BaseUrl
	}
	if info.ChannelType == common.AzureChannel.Type {
		info.ApiVersion = GetAPIVersion(c)
	}
	if info.ChannelType == common.OpenAIChannel.Type || info.ChannelType == common.AnthropicChannel.Type ||
		info.ChannelType == common.AwsChannel.Type || info.ChannelType == common.GeminiChannel.Type ||
		info.ChannelType == common.CloudflareChannel.Type ||
		info.ChannelType == common.VertexClaudeChannel.Type || info.ChannelType == common.ScholarAIChannel.Type {
		info.SupportStreamOptions = true
	}
	return info, nil
}

func (info *RelayInfo) SetPromptTokens(promptTokens int) {
	info.PromptTokens = promptTokens
}

func (info *RelayInfo) SetIsStream(isStream bool) {
	info.IsStream = isStream
}

func (info *RelayInfo) SetFirstResponseTime() {
	if !info.setFirstResponse {
		info.FirstResponseTime = time.Now()
		info.setFirstResponse = true
	}
}

type TaskRelayInfo struct {
	ChannelType       int
	ChannelId         int
	TokenId           int
	UserId            int
	Group             string
	StartTime         time.Time
	ApiType           int
	RelayMode         int
	UpstreamModelName string
	RequestURLPath    string
	ApiKey            string
	BaseUrl           string

	Action       string
	OriginTaskID string

	ConsumeQuota bool
	Proxy        string
}

func GenTaskRelayInfo(c *gin.Context) (*TaskRelayInfo, error) {
	channelType := c.GetInt("channel")
	channelId := c.GetInt("channel_id")

	tokenId := c.GetInt("token_id")
	userId := c.GetInt("id")
	group := c.GetString("group")
	startTime := time.Now()

	apiType, _ := constant.ChannelType2APIType(channelType)

	info := &TaskRelayInfo{
		RelayMode:      constant.Path2RelayMode(c.Request.URL.Path),
		BaseUrl:        c.GetString("base_url"),
		RequestURLPath: c.Request.URL.String(),
		ChannelType:    channelType,
		ChannelId:      channelId,
		TokenId:        tokenId,
		UserId:         userId,
		Group:          group,
		StartTime:      startTime,
		ApiType:        apiType,
		ApiKey:         strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "),
		Proxy:          c.GetString("proxy"),
	}
	if info.BaseUrl == "" {
		ch, exists := common.ChannelMap[channelType]
		if !exists {
			return nil, errors.New("channel not exists")
		}
		info.BaseUrl = ch.BaseUrl
	}
	return info, nil
}
