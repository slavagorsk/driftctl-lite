package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func TestIsCloudFrontType(t *testing.T) {
	if !isCloudFrontType("aws_cloudfront_distribution") {
		t.Error("expected aws_cloudfront_distribution to be a CloudFront type")
	}
	if isCloudFrontType("aws_s3_bucket") {
		t.Error("expected aws_s3_bucket not to be a CloudFront type")
	}
}

func TestFetchCloudFrontResource_UnsupportedType(t *testing.T) {
	_, err := fetchCloudFrontResource(context.Background(), "us-east-1", "aws_s3_bucket", "some-id")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestFetchCloudFrontResource_EmptyID(t *testing.T) {
	_, err := fetchCloudFrontResource(context.Background(), "us-east-1", "aws_cloudfront_distribution", "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

type mockCloudFrontClient struct {
	out *cloudfront.GetDistributionOutput
	err error
}

func (m *mockCloudFrontClient) GetDistribution(_ context.Context, _ *cloudfront.GetDistributionInput, _ ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error) {
	return m.out, m.err
}

func TestFetchCloudFrontResource_DelegatesToFetchCloudFrontDistribution(t *testing.T) {
	mock := &mockCloudFrontClient{
		out: &cloudfront.GetDistributionOutput{
			Distribution: &types.Distribution{
				DomainName: aws.String("d1234.cloudfront.net"),
				DistributionConfig: &types.DistributionConfig{
					Enabled: aws.Bool(true),
					Comment: aws.String("test dist"),
				},
			},
		},
	}

	attrs, err := FetchCloudFrontDistribution(context.Background(), mock, "EDFDVBD6EXAMPLE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["domain_name"] != "d1234.cloudfront.net" {
		t.Errorf("expected domain_name d1234.cloudfront.net, got %s", attrs["domain_name"])
	}
	if attrs["enabled"] != "true" {
		t.Errorf("expected enabled=true, got %s", attrs["enabled"])
	}
}

func TestFetchCloudFrontDistribution_Error(t *testing.T) {
	mock := &mockCloudFrontClient{err: errors.New("api error")}
	_, err := FetchCloudFrontDistribution(context.Background(), mock, "EDFDVBD6EXAMPLE")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
