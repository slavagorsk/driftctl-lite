package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"context"
)

func isLambdaType(resourceType string) bool {
	return resourceType == "aws_lambda_function"
}

func lambdaClientFromConfig(cfg aws.Config) *lambda.Client {
	return lambda.NewFromConfig(cfg)
}

func newLambdaConfigFromRegion(region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
}

func fetchLambdaResource(region, resourceType, resourceID string) (map[string]string, error) {
	if !isLambdaType(resourceType) {
		return nil, fmt.Errorf("unsupported resource type for lambda fetcher: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}
	cfg, err := newLambdaConfigFromRegion(region)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := lambdaClientFromConfig(cfg)
	return FetchLambdaFunction(client, resourceID)
}
