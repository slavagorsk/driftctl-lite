package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// rdsClient is the interface used for fetching RDS resources.
type rdsClient interface {
	DescribeDBInstances(ctx context.Context, params *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
}

// isRDSType reports whether the given Terraform resource type is managed
// by the RDS fetcher.
func isRDSType(resourceType string) bool {
	return resourceType == "aws_db_instance"
}

// rdsClientFromConfig creates a new RDS client from an existing AWS config.
func rdsClientFromConfig(cfg aws.Config) (*rds.Client, error) {
	return rds.NewFromConfig(cfg), nil
}

// newRDSClient creates a new RDS client configured for the given AWS region.
func newRDSClient(region string) (*rds.Client, error) {
	if region == "" {
		return nil, fmt.Errorf("region must not be empty")
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for RDS: %w", err)
	}
	return rds.NewFromConfig(cfg), nil
}

// fetchRDSResource retrieves the attributes of an RDS resource identified by
// resourceType and id using the provided client.
func fetchRDSResource(client rdsClient, resourceType, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("empty resource ID for type %s", resourceType)
	}
	switch resourceType {
	case "aws_db_instance":
		return FetchRDSInstance(client, id)
	default:
		return nil, fmt.Errorf("unsupported RDS resource type: %s", resourceType)
	}
}
