package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

// RDSClient defines the interface for RDS operations used in fetching.
type RDSClient interface {
	DescribeDBInstances(ctx context.Context, params *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
}

// FetchRDSInstance retrieves attributes of an RDS DB instance by its identifier.
func FetchRDSInstance(ctx context.Context, client RDSClient, dbInstanceID string) (map[string]string, error) {
	if dbInstanceID == "" {
		return nil, fmt.Errorf("rds: dbInstanceID must not be empty")
	}

	out, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbInstanceID),
	})
	if err != nil {
		return nil, fmt.Errorf("rds: describe db instances: %w", err)
	}

	if len(out.DBInstances) == 0 {
		return nil, fmt.Errorf("rds: instance %q not found", dbInstanceID)
	}

	db := out.DBInstances[0]
	attrs := map[string]string{
		"db_instance_identifier": aws.ToString(db.DBInstanceIdentifier),
		"db_instance_class":      aws.ToString(db.DBInstanceClass),
		"engine":                 aws.ToString(db.Engine),
		"engine_version":         aws.ToString(db.EngineVersion),
		"db_name":                aws.ToString(db.DBName),
		"status":                 aws.ToString(db.DBInstanceStatus),
		"multi_az":               fmt.Sprintf("%t", db.MultiAZ),
		"publicly_accessible":    fmt.Sprintf("%t", db.PubliclyAccessible),
		"storage_type":           aws.ToString(db.StorageType),
		"allocated_storage":      fmt.Sprintf("%d", db.AllocatedStorage),
	}

	return attrs, nil
}
