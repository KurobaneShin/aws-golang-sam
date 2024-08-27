package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Item struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func handler(ctx context.Context) error {
	// Load the SDK's configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("sa-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	item := Item{
		ID:   "example_id",
		Data: "example_data",
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatalf("got error marshalling map: %s", err)
	}

	// Put item in DynamoDB table
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("MyTable"),
	}

	_, err = svc.PutItem(ctx, input)
	if err != nil {
		log.Fatalf("got error calling PutItem: %s", err)
	}

	log.Println("Successfully added item to DynamoDB table")

	return nil
}

func main() {
	lambda.Start(handler)
}
