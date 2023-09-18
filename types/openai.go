package types

type OpenAIWithFunctionRequest struct {
	Model            string           `json:"model"`
	Messages         []OpenAIMessage  `json:"messages"`
	MaxTokens        int              `json:"max_tokens"`
	Temperature      float64          `json:"temperature,omitempty"`
	Functions        []OpenAIFunction `json:"functions,omitempty"`
	FunctionCallName interface{}      `json:"function_call,omitempty"`
}

type OpenAIMessage struct {
	Role         string              `json:"role"`
	Content      string              `json:"content"`
	Name         string              `json:"name,omitempty"`
	FunctionCall *OpenAIFunctionCall `json:"function_call,omitempty"`
}

type OpenAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type OpenAIFunction struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Parameters  OpenAIFunctionParameters `json:"parameters"`
}

type OpenAIFunctionParameters struct {
	Type       string                               `json:"type"`
	Properties map[string]OpenAiPropertyDescription `json:"properties"`
	Required   []string                             `json:"required"`
}

type OpenAiPropertyDescription struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type OpenAIResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int                  `json:"created"`
	Model   string               `json:"model"`
	Choices []OpenAIChoice       `json:"choices"`
	Usage   OpenAIUsage          `json:"usage"`
	Erorr   *OpenAIErrorResponse `json:"error,omitempty"`
}

type OpenAIErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type OpenAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type OpenAIFunctionCallName struct {
	Name string `json:"name"`
}

var OpenAIFunctionCALLAuto = PtrString("auto")

func PtrString(s string) *string { return &s }
