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

func isRDSType(resourceType string) bool {
	return resourceType == "aws_db_instance"
}

func rdsClientFromConfig(cfg aws.Config) (*rds.Client, error) {
	return rds.NewFromConfig(cfg), nil
}

func newRDSClient(region string) (*rds.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for RDS: %w", err)
	}
	return rds.NewFromConfig(cfg), nil
}

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
