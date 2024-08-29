package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func handleRequest(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := eventbridge.NewFromConfig(cfg)

	input := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: nil, // Uses default event bus
				Source:       aws.String("custom.myapp"),
				DetailType:   aws.String("MyCustomEvent"),
				Detail:       aws.String(`{"message": "Hello from sender!"}`),
			},
		},
	}

	result, err := client.PutEvents(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error putting events: %v", err)
	}

	return fmt.Sprintf("Successfully sent event. Entries: %v", result.Entries), nil
}

func main() {
	lambda.Start(handleRequest)
}
