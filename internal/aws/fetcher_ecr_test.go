package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func TestIsECRType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"aws_ecr_repository", true},
		{"aws_s3_bucket", false},
		{"aws_iam_role", false},
		{"", false},
	}
	for _, tt := range tests {
		got := isECRType(tt.input)
		if got != tt.expected {
			t.Errorf("isECRType(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestFetchECRResource_UnsupportedType(t *testing.T) {
	client := ecr.NewFromConfig(aws.Config{})
	_, err := fetchECRResource(context.Background(), client, "aws_s3_bucket", "my-repo")
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}

func TestFetchECRResource_EmptyID(t *testing.T) {
	client := ecr.NewFromConfig(aws.Config{})
	_, err := fetchECRResource(context.Background(), client, "aws_ecr_repository", "")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
}

func TestFetchECRResource_DelegatesToFetchECRRepository(t *testing.T) {
	mockClient := &mockECRClient{
		describeOutput: &ecr.DescribeRepositoriesOutput{
			Repositories: []types.Repository{
				{
					RepositoryName: aws.String("my-repo"),
					RepositoryArn:  aws.String("arn:aws:ecr:us-east-1:123456789012:repository/my-repo"),
				},
			},
		},
	}
	attrs, err := FetchECRRepository(context.Background(), mockClient, "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["repository_name"] != "my-repo" {
		t.Errorf("expected repository_name=my-repo, got %v", attrs["repository_name"])
	}
}
