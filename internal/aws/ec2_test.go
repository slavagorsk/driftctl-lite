package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// mockEC2Client implements a minimal EC2 API surface for testing.
type mockEC2Client struct {
	output *ec2.DescribeInstancesOutput
	err    error
}

func (m *mockEC2Client) DescribeInstances(
	_ context.Context,
	_ *ec2.DescribeInstancesInput,
	_ ...func(*ec2.Options),
) (*ec2.DescribeInstancesOutput, error) {
	return m.output, m.err
}

func TestFetchEC2Instance_Success(t *testing.T) {
	client := &mockEC2Client{
		output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{
				{
					Instances: []types.Instance{
						{
							InstanceId:   aws.String("i-0abc123def456"),
							InstanceType: types.InstanceTypeT3Micro,
							State: &types.InstanceState{
								Name: types.InstanceStateNameRunning,
							},
							Tags: []types.Tag{
								{Key: aws.String("Name"), Value: aws.String("web-server")},
								{Key: aws.String("Env"), Value: aws.String("prod")},
							},
						},
					},
				},
			},
		},
	}

	got, err := FetchEC2Instance(context.Background(), client, "i-0abc123def456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["id"] != "i-0abc123def456" {
		t.Errorf("expected id=i-0abc123def456, got %q", got["id"])
	}
	if got["instance_type"] != "t3.micro" {
		t.Errorf("expected instance_type=t3.micro, got %q", got["instance_type"])
	}
	if got["state"] != "running" {
		t.Errorf("expected state=running, got %q", got["state"])
	}
	if got["tag_Name"] != "web-server" {
		t.Errorf("expected tag_Name=web-server, got %q", got["tag_Name"])
	}
	if got["tag_Env"] != "prod" {
		t.Errorf("expected tag_Env=prod, got %q", got["tag_Env"])
	}
}

func TestFetchEC2Instance_Error(t *testing.T) {
	client := &mockEC2Client{
		err: errors.New("api failure"),
	}

	_, err := FetchEC2Instance(context.Background(), client, "i-0abc123def456")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchEC2Instance_NotFound(t *testing.T) {
	client := &mockEC2Client{
		output: &ec2.DescribeInstancesOutput{
			Reservations: []types.Reservation{},
		},
	}

	_, err := FetchEC2Instance(context.Background(), client, "i-doesnotexist")
	if err == nil {
		t.Fatal("expected not-found error, got nil")
	}
}

func TestFetchEC2Instance_EmptyID(t *testing.T) {
	client := &mockEC2Client{
		output: &ec2.DescribeInstancesOutput{},
	}

	_, err := FetchEC2Instance(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty instance ID, got nil")
	}
}

func TestFlattenEC2Tags(t *testing.T) {
	tags := []types.Tag{
		{Key: aws.String("Project"), Value: aws.String("driftctl")},
		{Key: aws.String("Owner"), Value: aws.String("team-a")},
	}

	attrs := make(map[string]string)
	flattenEC2Tags(attrs, tags)

	if attrs["tag_Project"] != "driftctl" {
		t.Errorf("expected tag_Project=driftctl, got %q", attrs["tag_Project"])
	}
	if attrs["tag_Owner"] != "team-a" {
		t.Errorf("expected tag_Owner=team-a, got %q", attrs["tag_Owner"])
	}
}
