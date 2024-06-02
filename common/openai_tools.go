package common

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"time"
)

func ConstructChatCompletionStreamReponse(model, answerID string, answer string) openai.ChatCompletionStreamResponse {
	resp := openai.ChatCompletionStreamResponse{
		ID:      answerID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []openai.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Content: answer,
				},
			},
		},
	}
	return resp
}

func SendChatData(w http.ResponseWriter, model string, chatID, data string) {
	dataRune := []rune(data)
	for _, d := range dataRune {
		respData := ConstructChatCompletionStreamReponse(model, chatID, string(d))
		byteData, _ := json.Marshal(respData)
		_, _ = fmt.Fprintf(w, "data: %s\n\n", string(byteData))
		w.(http.Flusher).Flush()
		time.Sleep(1 * time.Millisecond)
	}
}

func SendChatDone(w http.ResponseWriter) {
	_, _ = fmt.Fprintf(w, "data: [DONE]")
	w.(http.Flusher).Flush()
}
