package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func isCloudWatchType(resourceType string) bool {
	return resourceType == "aws_cloudwatch_metric_alarm"
}

func cloudWatchClientFromConfig(region string) (CloudWatchClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config for cloudwatch: %w", err)
	}
	return cloudwatch.NewFromConfig(cfg), nil
}

func fetchCloudWatchResource(ctx context.Context, region, resourceType, resourceID string) (map[string]string, error) {
	if !isCloudWatchType(resourceType) {
		return nil, fmt.Errorf("unsupported cloudwatch resource type: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}
	client, err := cloudWatchClientFromConfig(region)
	if err != nil {
		return nil, err
	}
	return FetchCloudWatchAlarm(ctx, client, resourceID)
}
