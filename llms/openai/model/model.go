package model

type Response struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Message struct {
	Role         Role          `json:"role"`
	Content      string        `json:"content"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type Role string

const (
	System    Role = "system"
	User      Role = "user"
	Assistant Role = "assistant"
)

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
	GPT_3_5_Turbo      Model = "gpt-3.5-turbo"
	GPT_3_5_Turbo_0613 Model = "gpt-3.5-turbo-0613"
	GPT_4              Model = "gpt-4"
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
