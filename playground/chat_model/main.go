package main

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
	"log"
)

func main() {
	ctx := context.Background()
	modelName := "qwen3-vl:32b"

	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   modelName,
	})
	if err != nil {
		log.Printf("NewChatModel failed, err=%v\n", err)
		return
	}

	resp, err := chatModel.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "你是谁？",
		},
	})
	if err != nil {
		log.Printf("Generate failed, err=%v\n", err)
		return
	}
	log.Printf("resp=%+v\n", resp.Content)
}
