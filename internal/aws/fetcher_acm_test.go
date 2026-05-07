package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsACMType(t *testing.T) {
	assert.True(t, isACMType("aws_acm_certificate"))
	assert.False(t, isACMType("aws_s3_bucket"))
	assert.False(t, isACMType(""))
}

func TestFetchACMResource_UnsupportedType(t *testing.T) {
	_, err := fetchACMResource(context.Background(), &mockACMClient{}, "aws_lambda_function", "arn:aws:acm:us-east-1:123:certificate/x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported ACM resource type")
}

func TestFetchACMResource_EmptyID(t *testing.T) {
	_, err := fetchACMResource(context.Background(), &mockACMClient{}, "aws_acm_certificate", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchACMResource_DelegatesToFetchACMCertificate(t *testing.T) {
	client := &mockACMClient{
		out: &acm.DescribeCertificateOutput{
			Certificate: &types.CertificateDetail{
				CertificateArn: aws.String("arn:aws:acm:us-east-1:123:certificate/abc"),
				DomainName:     aws.String("example.com"),
				Status:         types.CertificateStatusIssued,
			},
		},
	}
	attrs, err := fetchACMResource(context.Background(), client, "aws_acm_certificate", "arn:aws:acm:us-east-1:123:certificate/abc")
	require.NoError(t, err)
	assert.Equal(t, "example.com", attrs["domain"])
	assert.Equal(t, "ISSUED", attrs["status"])
}
