package common

import "sync"

type Channel struct {
	Type    int
	BaseUrl string
}

var (
	UnknownChannel        = Channel{Type: 0, BaseUrl: ""}
	OpenAIChannel         = Channel{Type: 1, BaseUrl: "https://api.openai.com"}
	MidjourneyChannel     = Channel{Type: 2, BaseUrl: "https://oa.api2d.net"}
	AzureChannel          = Channel{Type: 3, BaseUrl: ""}
	OllamaChannel         = Channel{Type: 4, BaseUrl: "http://localhost:11434"}
	MidjourneyPlusChannel = Channel{Type: 5, BaseUrl: "https://api.openai-sb.com"}
	OpenAIMaxChannel      = Channel{Type: 6, BaseUrl: "https://api.openaimax.com"}
	OhMyGPTChannel        = Channel{Type: 7, BaseUrl: "https://api.ohmygpt.com"}
	CustomChannel         = Channel{Type: 8, BaseUrl: ""}
	AILSChannel           = Channel{Type: 9, BaseUrl: "https://api.caipacity.com"}
	AIProxyChannel        = Channel{Type: 10, BaseUrl: "https://api.aiproxy.io"}
	PaLMChannel           = Channel{Type: 11, BaseUrl: ""}
	API2GPTChannel        = Channel{Type: 12, BaseUrl: "https://api.api2gpt.com"}
	AIGC2DChannel         = Channel{Type: 13, BaseUrl: "https://api.aigc2d.com"}
	AnthropicChannel      = Channel{Type: 14, BaseUrl: "https://api.anthropic.com"}
	BaiduChannel          = Channel{Type: 15, BaseUrl: "https://aip.baidubce.com"}
	ZhipuChannel          = Channel{Type: 16, BaseUrl: "https://open.bigmodel.cn"}
	AliChannel            = Channel{Type: 17, BaseUrl: "https://dashscope.aliyuncs.com"}
	XunfeiChannel         = Channel{Type: 18, BaseUrl: ""}
	AI360Channel          = Channel{Type: 19, BaseUrl: "https://ai.360.cn"}
	OpenRouterChannel     = Channel{Type: 20, BaseUrl: "https://openrouter.ai/api"}
	AIProxyLibraryChannel = Channel{Type: 21, BaseUrl: "https://api.aiproxy.io"}
	FastGPTChannel        = Channel{Type: 22, BaseUrl: "https://fastgpt.run/api/openapi"}
	TencentChannel        = Channel{Type: 23, BaseUrl: "https://hunyuan.cloud.tencent.com"}
	GeminiChannel         = Channel{Type: 24, BaseUrl: "https://generativelanguage.googleapis.com"}
	MoonshotChannel       = Channel{Type: 25, BaseUrl: "https://api.moonshot.cn"}
	ZhipuV4Channel        = Channel{Type: 26, BaseUrl: "https://open.bigmodel.cn"}
	PerplexityChannel     = Channel{Type: 27, BaseUrl: "https://api.perplexity.ai"}
	LingYiWanWuChannel    = Channel{Type: 31, BaseUrl: "https://api.lingyiwanwu.com"}
	AwsChannel            = Channel{Type: 33, BaseUrl: ""}
	CohereChannel         = Channel{Type: 34, BaseUrl: "https://api.cohere.ai"}
	MiniMaxChannel        = Channel{Type: 35, BaseUrl: "https://api.minimax.chat"}
	SunoAPIChannel        = Channel{Type: 36, BaseUrl: ""}
	DifyChannel           = Channel{Type: 37, BaseUrl: ""}
	JinaChannel           = Channel{Type: 38, BaseUrl: "https://api.jina.ai"}
	CloudflareChannel     = Channel{Type: 39, BaseUrl: "https://api.cloudflare.com"}
	SiliconFlowChannel    = Channel{Type: 40, BaseUrl: "https://api.siliconflow.cn"}

	ScholarAIChannel = Channel{Type: 10001, BaseUrl: "https://api.scholarai.io"}
	DoubaoChannel    = Channel{Type: 10002, BaseUrl: "https://ark.cn-beijing.volces.com"}
	GcpClaudeChannel = Channel{Type: 10003, BaseUrl: ""}
)

var ChannelList = []Channel{
	UnknownChannel,
	OpenAIChannel,
	MidjourneyChannel,
	AzureChannel,
	OllamaChannel,
	MidjourneyPlusChannel,
	OpenAIMaxChannel,
	OhMyGPTChannel,
	CustomChannel,
	AILSChannel,
	AIProxyChannel,
	PaLMChannel,
	API2GPTChannel,
	AIGC2DChannel,
	AnthropicChannel,
	BaiduChannel,
	ZhipuChannel,
	AliChannel,
	XunfeiChannel,
	AI360Channel,
	OpenRouterChannel,
	AIProxyLibraryChannel,
	FastGPTChannel,
	TencentChannel,
	GeminiChannel,
	MoonshotChannel,
	ZhipuV4Channel,
	PerplexityChannel,
	LingYiWanWuChannel,
	AwsChannel,
	CohereChannel,
	MiniMaxChannel,
	SunoAPIChannel,
	DifyChannel,
	JinaChannel,
	CloudflareChannel,
	SiliconFlowChannel,

	ScholarAIChannel,
	DoubaoChannel,
	GcpClaudeChannel,
}

var ChannelMap map[int]Channel
var ChannelMapRWMutex sync.RWMutex

func InitChannelMap() {
	ChannelMapRWMutex.Lock()
	ChannelMap = make(map[int]Channel)
	for _, channel := range ChannelList {
		ChannelMap[channel.Type] = channel
	}
	ChannelMapRWMutex.Unlock()
}
