package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/acm"
)

const tfACMCertificate = "aws_acm_certificate"

func isACMType(resourceType string) bool {
	return resourceType == tfACMCertificate
}

func acmClientFromConfig(region string) (ACMClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config for ACM: %w", err)
	}
	return acm.NewFromConfig(cfg), nil
}

func fetchACMResource(ctx context.Context, client ACMClient, resourceType, id string) (map[string]string, error) {
	if !isACMType(resourceType) {
		return nil, fmt.Errorf("unsupported ACM resource type: %s", resourceType)
	}
	if id == "" {
		return nil, fmt.Errorf("resource ID must not be empty")
	}
	return FetchACMCertificate(ctx, client, id)
}
