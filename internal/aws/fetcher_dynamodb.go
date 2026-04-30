package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func isDynamoDBType(resourceType string) bool {
	return resourceType == "aws_dynamodb_table"
}

func dynamoDBClientFromConfig(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

func newDynamoDBConfigFromRegion(region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
}

func fetchDynamoDBResource(region, resourceType, resourceID string) (map[string]string, error) {
	if !isDynamoDBType(resourceType) {
		return nil, fmt.Errorf("unsupported resource type for dynamodb fetcher: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}
	cfg, err := newDynamoDBConfigFromRegion(region)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := dynamoDBClientFromConfig(cfg)
	return FetchDynamoDBTable(client, resourceID)
}
