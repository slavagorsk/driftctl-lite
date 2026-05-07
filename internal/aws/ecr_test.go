package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type mockECRClient struct {
	out *ecr.DescribeRepositoriesOutput
	err error
}

func (m *mockECRClient) DescribeRepositories(_ context.Context, _ *ecr.DescribeRepositoriesInput, _ ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error) {
	return m.out, m.err
}

func TestFetchECRRepository_Success(t *testing.T) {
	client := &mockECRClient{
		out: &ecr.DescribeRepositoriesOutput{
			Repositories: []types.Repository{
				{
					RepositoryName: aws.String("my-repo"),
					RepositoryUri:  aws.String("123456789.dkr.ecr.us-east-1.amazonaws.com/my-repo"),
					RegistryId:     aws.String("123456789"),
					ImageTagMutability: types.ImageTagMutabilityMutable,
					ImageScanningConfiguration: &types.ImageScanningConfiguration{
						ScanOnPush: true,
					},
				},
			},
		},
	}

	attrs, err := FetchECRRepository(context.Background(), client, "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["repository_name"] != "my-repo" {
		t.Errorf("expected repository_name=my-repo, got %q", attrs["repository_name"])
	}
	if attrs["image_scanning_scan_on_push"] != "true" {
		t.Errorf("expected scan_on_push=true, got %q", attrs["image_scanning_scan_on_push"])
	}
	if attrs["image_tag_mutability"] != "MUTABLE" {
		t.Errorf("expected image_tag_mutability=MUTABLE, got %q", attrs["image_tag_mutability"])
	}
}

func TestFetchECRRepository_EmptyID(t *testing.T) {
	_, err := FetchECRRepository(context.Background(), &mockECRClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty repository name")
	}
}

func TestFetchECRRepository_Error(t *testing.T) {
	client := &mockECRClient{err: errors.New("api failure")}
	_, err := FetchECRRepository(context.Background(), client, "my-repo")
	if err == nil {
		t.Fatal("expected error from API failure")
	}
}

func TestFetchECRRepository_NotFound(t *testing.T) {
	client := &mockECRClient{
		out: &ecr.DescribeRepositoriesOutput{
			Repositories: []types.Repository{},
		},
	}
	_, err := FetchECRRepository(context.Background(), client, "missing-repo")
	if err == nil {
		t.Fatal("expected error when repository not found")
	}
}
