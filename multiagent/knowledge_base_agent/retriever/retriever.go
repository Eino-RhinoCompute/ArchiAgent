package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	rr "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/redis/go-redis/v9"
)

const (
	embedderModel = "nomic-embed-text:v1.5"
)

func main() {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:          "localhost:6379",
		Protocol:      2,
		UnstableResp3: true,
	})
	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		Model:   embedderModel,
		BaseURL: "http://127.0.0.1:11434",
	})
	if err != nil {
		panic(err)
	}
	r, err := rr.NewRetriever(ctx, &rr.RetrieverConfig{
		Client:    client,
		Index:     "doc_index",
		Embedding: embedder,
	})
	if err != nil {
		panic(err)
	}
	docs, err := r.Retrieve(ctx, "dog")
	if err != nil {
		panic(err)
	}
	for _, v := range docs {
		fmt.Printf("ID:%s, CONTENT:%v \n", v.ID, v.Content)
	}
}
