// models - пакет c моделями для отправки запросов к ИИ
package models

// AIRequest - структура для
// запроса на генерацию ответа от Deepseek
type AIRequest struct {
	Model           string          `json:"model"`
	Messages        []Message       `json:"messages"`
	Thinking        *Thinking       `json:"thinking,omitempty"`
	ReasoningEffort string          `json:"reasoning_effort,omitempty"`
	Stream          bool            `json:"stream"`
	ResponseFormat  *ResponseFormat `json:"response_format,omitempty"`
}

// Message представляет сообщение диалога,
// передаваемое модели
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Thinking содержит настройки режима рассуждения модели
type Thinking struct {
	Type string `json:"type"`
}

// ResponseFormat определяет формат ответа модели
type ResponseFormat struct {
	Type string `json:"type"`
}

// AIResponse представляет ответ DeepSeek API
type AIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice представляет один вариант ответа модели
type Choice struct {
	Index        int     `json:"index"`
	FinishReason string  `json:"finish_reason"`
	Message      Message `json:"message"`
}

// Usage содержит статистику использования токенов
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
