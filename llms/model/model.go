package model

type Question struct {
	SystemPrompt string
	UserQuestion string
}

type Message struct {
	Role         string        `json:"role"`
	Content      string        `json:"content"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type Response struct {
	ID      string         `json:"id"`
	Model   string         `json:"model"`
	Choices []Choice       `json:"choices"`
	Usage   Usage          `json:"usage"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type EmbeddingResponse struct {
	Data []EmbeddingData `json:"data"`
}

type EmbeddingData struct {
	Vector []float32 `json:"embedding"`
}

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model Model  `json:"model"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Request struct {
	Model        Model         `json:"model"`
	Messages     []Message     `json:"messages"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
	Functions    []Function    `json:"functions,omitempty"`
	Temperature  *float64      `json:"temperature,omitempty"`
	MaxTokens    *uint         `json:"max_tokens,omitempty"`
}

type Model string

const (
	GPT_3_5_Turbo          Model = "gpt-3.5-turbo"
	GPT_3_5_Turbo_0613     Model = "gpt-3.5-turbo-0613"
	GPT_4                  Model = "gpt-4"
	Text_Embedding_Ada_002 Model = "text-embedding-ada-002"
)

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments,omitempty"`
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       ParameterType        `json:"type"`
	Properties map[string]Parameter `json:"properties"`
	Required   []string             `json:"required,omitempty"`
}

type Parameter struct {
	Type        ParameterType `json:"type"`
	Description string        `json:"description"`
	Enum        []any         `json:"enum,omitempty"`
	Items       *Item         `json:"items,omitempty"`
}

type ParameterType string

const (
	Object  ParameterType = "object"
	Integer ParameterType = "integer"
	String  ParameterType = "string"
	Boolean ParameterType = "boolean"
	Array   ParameterType = "array"
)

type Item struct {
	Type ParameterType `json:"type"`
}
