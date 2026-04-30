package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func isSQSType(resourceType string) bool {
	return resourceType == "aws_sqs_queue"
}

func sqsClientFromConfig(region string) (SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for SQS: %w", err)
	}
	return sqs.NewFromConfig(cfg), nil
}

func fetchSQSResource(ctx context.Context, region, resourceType, resourceID string) (map[string]string, error) {
	if !isSQSType(resourceType) {
		return nil, fmt.Errorf("unsupported SQS resource type: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}
	client, err := sqsClientFromConfig(region)
	if err != nil {
		return nil, err
	}
	return FetchSQSQueue(ctx, client, resourceID)
}
