package main

import "suno-api/entity/po"

type SubmitGenSongReq struct {
	Prompt               string `json:"prompt"`
	Mv                   string `json:"mv"`
	Title                string `json:"title"`
	Tags                 string `json:"tags"`
	GptDescriptionPrompt string `json:"gpt_description_prompt,omitempty"`

	TaskID           string   `json:"task_id"`
	ContinueAt       *float64 `json:"continue_at,omitempty"`
	ContinueClipId   *string  `json:"continue_clip_id,omitempty"`
	MakeInstrumental bool     `json:"make_instrumental"`
}

type SubmitGenLyricsReq struct {
	Prompt string `json:"prompt"`
}

type FetchReq struct {
	IDs    []string `json:"ids"`
	Action string   `json:"action"`
}

type GenSongResponse struct {
	ID                string        `json:"id"`
	Clips             []po.SunoSong `json:"clips"`
	Metadata          any           `json:"metadata"`
	MajorModelVersion string        `json:"major_model_version"`
	Status            string        `json:"status"`
	CreatedAt         string        `json:"created_at"`
	BatchSize         int64         `json:"batch_size"`
}

type GenLyricsResponse struct {
	ID string `json:"id"`
}

type FetchLyricsResponse struct {
	Text   string `json:"text"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

type GeneralOpenAIRequest struct {
	Model            string      `json:"model,omitempty"`
	Messages         []Message   `json:"messages,omitempty"`
	Stream           bool        `json:"stream,omitempty"`
	MaxTokens        uint        `json:"max_tokens,omitempty"`
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"top_p,omitempty"`
	TopK             int         `json:"top_k,omitempty"`
	FunctionCall     interface{} `json:"function_call,omitempty"`
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`
	ToolChoice       string      `json:"tool_choice,omitempty"`
	Tools            []Tool      `json:"tools,omitempty"`
}

type Function struct {
	Url         string    `json:"url,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Parameters  Parameter `json:"parameters"`
	Arguments   string    `json:"arguments,omitempty"`
}
type Parameter struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type Tool struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created interface{}                           `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}

type ChatCompletionsStreamResponseChoice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
		Role    string `json:"role,omitempty"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Name    *string `json:"name,omitempty"`
}
