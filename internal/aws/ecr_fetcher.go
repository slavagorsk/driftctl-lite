package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func isECRType(resourceType string) bool {
	return resourceType == "aws_ecr_repository"
}

func ecrClientFromConfig(cfg aws.Config) *ecr.Client {
	return ecr.NewFromConfig(cfg)
}

func newECRConfigFromRegion(ctx context.Context, region string) (aws.Config, error) {
	if region == "" {
		return aws.Config{}, fmt.Errorf("region must not be empty")
	}
	return config.LoadDefaultConfig(ctx, config.WithRegion(region))
}

func fetchECRResource(ctx context.Context, client *ecr.Client, resourceType, id string) (map[string]interface{}, error) {
	if !isECRType(resourceType) {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	if id == "" {
		return nil, fmt.Errorf("resource id must not be empty")
	}
	return FetchECRRepository(ctx, client, id)
}
