package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino/adk"

	"github.com/cloudwego/eino-examples/adk/common/prints"
	"github.com/cloudwego/eino-examples/adk/common/trace"
)

func main() {
	ctx := context.Background()

	traceCloseFn, startSpanFn := trace.AppendCozeLoopCallbackIfConfigured(ctx)
	defer traceCloseFn(ctx)

	sv, err := buildSupervisor(ctx)
	if err != nil {
		log.Fatalf("build layered supervisor failed: %v", err)
	}

	query := "find US and New York state GDP in 2024. what % of US GDP was New York state? " +
		"Then multiply that percentage by 1.589."

	ctx, endSpanFn := startSpanFn(ctx, "layered-supervisor", query)
	iter := adk.NewRunner(ctx, adk.RunnerConfig{
		EnableStreaming: true,
		Agent:           sv,
	}).Query(ctx, query)

	fmt.Println("\nuser query: ", query)

	var lastMessage adk.Message
	for {
		event, hasEvent := iter.Next()
		if !hasEvent {
			break
		}

		prints.Event(event)

		if event.Output != nil {
			lastMessage, _, err = adk.GetMessage(event)
		}
	}

	endSpanFn(ctx, lastMessage)

	// wait for all span to be ended
	time.Sleep(5 * time.Second)
}
