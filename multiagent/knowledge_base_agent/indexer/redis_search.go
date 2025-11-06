package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// --- 1. 连接到 Redis ---
	rdb := redis.NewClient(&redis.Options{
		Addr:          "localhost:6379", // Redis Stack 服务的地址
		UnstableResp3: true,
		Protocol:      2,
	})
	indexName := "doc_index"
	// 构建 KNN 查询
	// `*=>[KNN 2 @vector_content $blob]` 的含义:
	// - `*`: 匹配所有文档 (我们不过滤元数据)。
	// - `=>`: 表示这是一个混合查询，我们主要关心右边的向量部分。
	// - `[KNN 2 @vector_content $blob]`: 在 `vector_content` 字段上执行一个 K-最近邻查询，
	//   查找 2 个最近邻。`$blob` 是一个参数，我们将把查询向量的二进制数据传递给它。
	// DIALECT 2 是必须的，用于支持这种现代的查询语法。
	k := 2
	query := fmt.Sprintf("*=>[KNN %d @vector_content $blob AS score]", k)
	searchContent := "That is a happy person"
	embedder, err := ollama.NewEmbedder(ctx, &ollama.EmbeddingConfig{
		Model:   "modelscope.cn/nomic-ai/nomic-embed-text-v1.5-GGUF:latest",
		BaseURL: "http://127.0.0.1:11434",
	})
	if err != nil {
		panic(err)
	}
	embeddings, err := embedder.EmbedStrings(ctx, []string{searchContent})
	if err != nil {
		panic(err)
	}
	// 使用 redis.NewSearch 来构建带参数的查询
	searchResult, err := rdb.FTSearchWithArgs(ctx, indexName, query, &redis.FTSearchOptions{
		Params: map[string]interface{}{
			"blob": vector2Bytes(embeddings[0]),
		},
		DialectVersion: 2,
		Return: []redis.FTSearchReturn{
			{
				FieldName: "content",
			},
			{
				FieldName: "score",
			},
		},
	}).Result()
	if err != nil {
		panic(err)
	}
	for _, v := range searchResult.Docs {
		fmt.Printf("%v:%v \n", v.Fields["content"], v.Fields["score"])
	}
}

func vector2Bytes(vector []float64) []byte {
	float32Arr := make([]float32, len(vector))
	for i, v := range vector {
		float32Arr[i] = float32(v)
	}
	bytes := make([]byte, len(float32Arr)*4)
	for i, v := range float32Arr {
		binary.LittleEndian.PutUint32(bytes[i*4:], math.Float32bits(v))
	}
	return bytes
}
