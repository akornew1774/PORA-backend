// models - пакет c моделями для отправки запросов к ИИ
package models

// Content получает контент одного из ответов Deepseek на запрос
func (r *AIResponse) Content() string {
	if len(r.Choices) == 0 {
		return ""
	}

	return r.Choices[0].Message.Content
}
