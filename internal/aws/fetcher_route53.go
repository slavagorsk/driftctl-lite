package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func isRoute53Type(resourceType string) bool {
	return resourceType == "aws_route53_zone"
}

func route53ClientFromConfig(region string) (Route53Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config for route53: %w", err)
	}
	return route53.NewFromConfig(cfg), nil
}

func fetchRoute53Resource(ctx context.Context, client Route53Client, resourceType, resourceID string) (map[string]string, error) {
	if !isRoute53Type(resourceType) {
		return nil, fmt.Errorf("unsupported route53 resource type: %s", resourceType)
	}
	if strings.TrimSpace(resourceID) == "" {
		return nil, fmt.Errorf("resource ID must not be empty for type %s", resourceType)
	}
	return FetchRoute53HostedZone(ctx, client, resourceID)
}
