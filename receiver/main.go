package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	log.Printf("Received event: %s", event.DetailType)
	log.Printf("Event detail: %s", event.Detail)

	return fmt.Sprintf("Successfully processed event: %s", event.DetailType), nil
}

func main() {
	lambda.Start(handleRequest)
}
