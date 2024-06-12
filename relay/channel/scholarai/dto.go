package scholarai

type ScholarAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ScholarAIChatRequest struct {
	Model    string             `json:"model,omitempty"`
	Messages []ScholarAIMessage `json:"messages,omitempty"`
	Stream   bool               `json:"stream,omitempty"`
}

type ScholarAITextResponseChoice struct {
	Index            int `json:"index"`
	ScholarAIMessage `json:"message"`
	FinishReason     string `json:"finish_reason"`
}

type ScholarAITextResponse struct {
	Id                string                        `json:"id"`
	Object            string                        `json:"object"`
	Created           int64                         `json:"created"`
	Model             string                        `json:"model"`
	Logprods          string                        `json:"logprods"`
	Choices           []ScholarAITextResponseChoice `json:"choices"`
	SystemFingerprint string                        `json:"system_fingerprint"`
}
