package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

func isCloudFrontType(resourceType string) bool {
	return resourceType == "aws_cloudfront_distribution"
}

func cloudFrontClientFromConfig(region string) (CloudFrontClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config for cloudfront: %w", err)
	}
	return cloudfront.NewFromConfig(cfg), nil
}

func fetchCloudFrontResource(ctx context.Context, region, resourceType, id string) (map[string]string, error) {
	if !isCloudFrontType(resourceType) {
		return nil, fmt.Errorf("unsupported cloudfront resource type: %s", resourceType)
	}
	if id == "" {
		return nil, fmt.Errorf("resource id must not be empty")
	}

	client, err := cloudFrontClientFromConfig(region)
	if err != nil {
		return nil, err
	}

	return FetchCloudFrontDistribution(ctx, client, id)
}
