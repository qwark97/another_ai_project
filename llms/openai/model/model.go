package model

type Response struct {
	ID      string
	Model   string
	Choices []Choice
	Usage   Usage
}

type Choice struct {
	Message Message `json:"message,omitempty"`
}

type Message struct {
	Role    Role   `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
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
	Model    Model     `json:"model,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

type Model string

const (
	GPT_3_5_Turbo Model = "gpt-3.5-turbo"
	GPT_4         Model = "gpt-4"
)
