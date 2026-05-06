package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func isKinesisType(resourceType string) bool {
	return resourceType == "aws_kinesis_stream"
}

func kinesisClientFromConfig(region string) (KinesisClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config for kinesis: %w", err)
	}
	return kinesis.NewFromConfig(cfg), nil
}

func fetchKinesisResource(ctx context.Context, region, resourceType, resourceID string) (map[string]string, error) {
	if !isKinesisType(resourceType) {
		return nil, fmt.Errorf("unsupported kinesis resource type: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}

	client, err := kinesisClientFromConfig(region)
	if err != nil {
		return nil, err
	}

	return FetchKinesisStream(ctx, client, resourceID)
}
