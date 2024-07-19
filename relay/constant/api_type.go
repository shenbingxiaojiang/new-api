package constant

import (
	"one-api/common"
)

const (
	APITypeOpenAI = iota
	APITypeAnthropic
	APITypePaLM
	APITypeBaidu
	APITypeZhipu
	APITypeAli
	APITypeXunfei
	APITypeAIProxyLibrary
	APITypeTencent
	APITypeGemini
	APITypeZhipuV4
	APITypeOllama
	APITypePerplexity
	APITypeAws
	APITypeCohere
	APITypeDify
	APITypeJina
	APITypeCloudflare

	APITypeScholarAI
	APITypeVertexClaude
	APITypeDummy // this one is only for count, do not add any channel after this
)

func ChannelType2APIType(channelType int) (int, bool) {
	apiType := -1
	switch channelType {
	case common.OpenAIChannel.Type:
		apiType = APITypeOpenAI
	case common.AnthropicChannel.Type:
		apiType = APITypeAnthropic
	case common.BaiduChannel.Type:
		apiType = APITypeBaidu
	case common.PaLMChannel.Type:
		apiType = APITypePaLM
	case common.ZhipuChannel.Type:
		apiType = APITypeZhipu
	case common.AliChannel.Type:
		apiType = APITypeAli
	case common.XunfeiChannel.Type:
		apiType = APITypeXunfei
	case common.AIProxyLibraryChannel.Type:
		apiType = APITypeAIProxyLibrary
	case common.TencentChannel.Type:
		apiType = APITypeTencent
	case common.GeminiChannel.Type:
		apiType = APITypeGemini
	case common.ZhipuV4Channel.Type:
		apiType = APITypeZhipuV4
	case common.OllamaChannel.Type:
		apiType = APITypeOllama
	case common.PerplexityChannel.Type:
		apiType = APITypePerplexity
	case common.AwsChannel.Type:
		apiType = APITypeAws
	case common.CohereChannel.Type:
		apiType = APITypeCohere
	case common.DifyChannel.Type:
		apiType = APITypeDify
	case common.JinaChannel.Type:
		apiType = APITypeJina
	case common.ChannelCloudflare:
		apiType = APITypeCloudflare
	case common.ScholarAIChannel.Type:
		apiType = APITypeScholarAI
	case common.VertexClaudeChannel.Type:
		apiType = APITypeVertexClaude
	}
	if apiType == -1 {
		return APITypeOpenAI, false
	}
	return apiType, true
}
