package gcp_claude

import "one-api/relay/channel/claude"

type GcpClaudeRequest struct {
	// vertex-2023-10-16
	AnthropicVersion string                 `json:"anthropic_version"`
	System           string                 `json:"system"`
	Messages         []claude.ClaudeMessage `json:"messages"`
	MaxTokens        uint                   `json:"max_tokens,omitempty"`
	Temperature      float64                `json:"temperature,omitempty"`
	TopP             float64                `json:"top_p,omitempty"`
	TopK             int                    `json:"top_k,omitempty"`
	StopSequences    []string               `json:"stop_sequences,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
}
