package model

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/eino-contrib/ollama/api"
	"log"
)

const (
	modelName = "qwen3-vl:32b"
)

func NewChatModel() model.ToolCallingChatModel {
	ctx := context.Background()
	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL:  "http://localhost:11434",
		Model:    modelName,
		Thinking: &api.ThinkValue{Value: false},
	})
	if err != nil {
		log.Printf("NewChatModel failed, err=%v\n", err)
		return nil
	}
	return cm
}
