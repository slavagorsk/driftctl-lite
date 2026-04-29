package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// iamClientFromConfig creates an IAM client from the provided AWS config.
func iamClientFromConfig(cfg aws.Config) *iam.Client {
	return iam.NewFromConfig(cfg)
}

// isIAMType returns true if the given Terraform resource type is an IAM resource.
func isIAMType(resourceType string) bool {
	switch resourceType {
	case "aws_iam_user", "aws_iam_role":
		return true
	}
	return false
}

// fetchIAMResource fetches an IAM resource by type and ID, returning its attributes.
func fetchIAMResource(cfg aws.Config, resourceType, id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, fmt.Errorf("resource ID is empty for type %s", resourceType)
	}

	client := iamClientFromConfig(cfg)

	switch resourceType {
	case "aws_iam_user":
		return FetchIAMUser(context.Background(), client, id)
	case "aws_iam_role":
		return FetchIAMRole(context.Background(), client, id)
	default:
		return nil, fmt.Errorf("unsupported IAM resource type: %s", resourceType)
	}
}

// newIAMConfigFromRegion creates an AWS config for the given region.
func newIAMConfigFromRegion(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config for IAM: %w", err)
	}
	return cfg, nil
}
