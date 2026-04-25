package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type mockRDSClient struct {
	output *rds.DescribeDBInstancesOutput
	err    error
}

func (m *mockRDSClient) DescribeDBInstances(_ context.Context, _ *rds.DescribeDBInstancesInput, _ ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error) {
	return m.output, m.err
}

func TestFetchRDSInstance_Success(t *testing.T) {
	client := &mockRDSClient{
		output: &rds.DescribeDBInstancesOutput{
			DBInstances: []types.DBInstance{
				{
					DBInstanceIdentifier: aws.String("mydb"),
					DBInstanceClass:      aws.String("db.t3.micro"),
					Engine:               aws.String("mysql"),
					EngineVersion:        aws.String("8.0"),
					DBInstanceStatus:     aws.String("available"),
					MultiAZ:              false,
					PubliclyAccessible:   true,
					StorageType:          aws.String("gp2"),
					AllocatedStorage:     20,
				},
			},
		},
	}

	attrs, err := FetchRDSInstance(context.Background(), client, "mydb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["engine"] != "mysql" {
		t.Errorf("expected engine mysql, got %s", attrs["engine"])
	}
	if attrs["db_instance_class"] != "db.t3.micro" {
		t.Errorf("expected db.t3.micro, got %s", attrs["db_instance_class"])
	}
	if attrs["allocated_storage"] != "20" {
		t.Errorf("expected allocated_storage 20, got %s", attrs["allocated_storage"])
	}
}

func TestFetchRDSInstance_Error(t *testing.T) {
	client := &mockRDSClient{err: errors.New("api error")}
	_, err := FetchRDSInstance(context.Background(), client, "mydb")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchRDSInstance_NotFound(t *testing.T) {
	client := &mockRDSClient{
		output: &rds.DescribeDBInstancesOutput{DBInstances: []types.DBInstance{}},
	}
	_, err := FetchRDSInstance(context.Background(), client, "missing")
	if err == nil {
		t.Fatal("expected not-found error, got nil")
	}
}

func TestFetchRDSInstance_EmptyID(t *testing.T) {
	client := &mockRDSClient{}
	_, err := FetchRDSInstance(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
}
