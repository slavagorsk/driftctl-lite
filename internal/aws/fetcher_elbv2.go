package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

const tfTypeALB = "aws_alb"
const tfTypeLB = "aws_lb"

func isELBv2Type(resourceType string) bool {
	return resourceType == tfTypeALB || resourceType == tfTypeLB
}

func elbv2ClientFromConfig(cfg ResourceFetcherConfig) (ELBv2Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config for elbv2: %w", err)
	}
	return elasticloadbalancingv2.NewFromConfig(awsCfg), nil
}

func fetchELBv2Resource(ctx context.Context, cfg ResourceFetcherConfig, resourceType, id string) (map[string]string, error) {
	if !isELBv2Type(resourceType) {
		return nil, fmt.Errorf("unsupported resource type for elbv2 fetcher: %s", resourceType)
	}
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("resource id must not be empty")
	}
	client, err := elbv2ClientFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return FetchELBv2LoadBalancer(ctx, client, id)
}
