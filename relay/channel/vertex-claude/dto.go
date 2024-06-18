package vertex_claude

import "one-api/relay/channel/claude"

type VertexClaudeRequest struct {
	// vertex-2023-10-16
	AnthropicVersion string                 `json:"anthropic_version"`
	System           string                 `json:"system,omitempty"`
	Messages         []claude.ClaudeMessage `json:"messages"`
	MaxTokens        int                    `json:"max_tokens,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
}
