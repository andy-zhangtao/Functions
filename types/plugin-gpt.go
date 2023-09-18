package types

type PluginGPTInput struct {
	Prompt struct {
		System string `json:"system"`
	} `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	Model       string  `json:"model"`
}
