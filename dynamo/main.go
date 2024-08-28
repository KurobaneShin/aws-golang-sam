package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client

func init() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://172.17.0.1:8000",
				SigningRegion: "sa-east-1",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dynamoClient = dynamodb.NewFromConfig(cfg)
}

type Item struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

func insertItem(ctx context.Context, tableName string, item Item) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	}

	_, err = dynamoClient.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	return nil
}

func handleRequest(ctx context.Context) (string, error) {
	// Use dynamoClient to interact with your local DynamoDB
	// Example: List tables
	result, err := dynamoClient.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		return "", err
	}

	// Process the result
	for _, tableName := range result.TableNames {
		log.Printf("Table: %s", tableName)
	}

	// Insert an item
	item := Item{
		ID:   "1",
		Name: "Test Item",
	}
	err = insertItem(ctx, "MyTable", item)
	if err != nil {
		return "", err
	}

	fmt.Println("Successfully inserted item in table")

	return "Successfully used dynamoDB", nil
}

func main() {
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		log.Println("Running in local SAM environment")
	}
	lambda.Start(handleRequest)
}
