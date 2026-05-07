package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

func isEKSType(resourceType string) bool {
	return resourceType == "aws_eks_cluster"
}

func eksClientFromConfig(cfg aws.Config) EKSClient {
	return eks.NewFromConfig(cfg)
}

func newEKSConfigFromRegion(ctx context.Context, region string) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion(region))
}

func fetchEKSResource(ctx context.Context, region, resourceType, resourceID string) (map[string]string, error) {
	if !isEKSType(resourceType) {
		return nil, fmt.Errorf("unsupported EKS resource type: %s", resourceType)
	}
	if resourceID == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}

	cfg, err := newEKSConfigFromRegion(ctx, region)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for EKS: %w", err)
	}

	client := eksClientFromConfig(cfg)
	return FetchEKSCluster(ctx, client, resourceID)
}
