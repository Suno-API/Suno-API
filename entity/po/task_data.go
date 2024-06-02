package po

import (
	"database/sql/driver"
	"encoding/json"
)

type TaskData interface {
	SunoSongs | SunoLyrics | json.RawMessage
}

type SunoSongs []SunoSong

type SunoSong struct {
	ID                string       `json:"id"`
	VideoURL          string       `json:"video_url"`
	AudioURL          string       `json:"audio_url"`
	ImageURL          string       `json:"image_url"`
	ImageLargeURL     string       `json:"image_large_url"`
	IsVideoPending    bool         `json:"is_video_pending"`
	MajorModelVersion string       `json:"major_model_version"`
	ModelName         string       `json:"model_name"`
	IsLiked           bool         `json:"is_liked"`
	UserID            string       `json:"user_id"`
	DisplayName       string       `json:"display_name"`
	Handle            string       `json:"handle"`
	IsHandleUpdated   bool         `json:"is_handle_updated"`
	IsTrashed         bool         `json:"is_trashed"`
	Reaction          interface{}  `json:"reaction"`
	CreatedAt         string       `json:"created_at"`
	Status            string       `json:"status"`
	Title             string       `json:"title"`
	PlayCount         int64        `json:"play_count"`
	UpvoteCount       int64        `json:"upvote_count"`
	IsPublic          bool         `json:"is_public"`
	Metadata          SunoMetadata `json:"metadata"`
}

type SunoMetadata struct {
	Tags                 string      `json:"tags"`
	Prompt               string      `json:"prompt"`
	GPTDescriptionPrompt interface{} `json:"gpt_description_prompt"`
	AudioPromptID        interface{} `json:"audio_prompt_id"`
	History              interface{} `json:"history"`
	ConcatHistory        interface{} `json:"concat_history"`
	Type                 string      `json:"type"`
	Duration             interface{} `json:"duration"`
	RefundCredits        interface{} `json:"refund_credits"`
	Stream               bool        `json:"stream"`
	ErrorType            interface{} `json:"error_type"`
	ErrorMessage         interface{} `json:"error_message"`
}

func (m *SunoSongs) Scan(val interface{}) error {
	bytesValue, _ := val.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m SunoSongs) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type SunoLyrics struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

func (m *SunoLyrics) Scan(val interface{}) error {
	bytesValue, _ := val.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m SunoLyrics) Value() (driver.Value, error) {
	return json.Marshal(m)
}
