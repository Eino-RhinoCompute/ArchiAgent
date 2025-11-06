package embedding

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
)

func main() {
	ctx := context.Background()
	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		Model:   "nomic-embed-text:v1.5",
		BaseURL: "http://127.0.0.1:11434",
	})
	if err != nil {
		panic(err)
	}
	embeddings, err := embedder.EmbedStrings(ctx, []string{"hello world"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v \n", embeddings)
}
