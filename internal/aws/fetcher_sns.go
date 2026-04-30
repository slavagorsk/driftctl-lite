package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func isSNSType(resourceType string) bool {
	return resourceType == "aws_sns_topic"
}

func snsClientFromConfig(region string) (SNSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for SNS: %w", err)
	}
	return sns.NewFromConfig(cfg), nil
}

func fetchSNSResource(ctx context.Context, client SNSClient, resourceType, id string) (map[string]string, error) {
	if !isSNSType(resourceType) {
		return nil, fmt.Errorf("unsupported SNS resource type: %s", resourceType)
	}
	if id == "" {
		return nil, fmt.Errorf("resource ID must not be empty for type %s", resourceType)
	}
	return FetchSNSTopic(ctx, client, id)
}
