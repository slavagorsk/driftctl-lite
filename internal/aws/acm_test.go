package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockACMClient struct {
	out *acm.DescribeCertificateOutput
	err error
}

func (m *mockACMClient) DescribeCertificate(_ context.Context, _ *acm.DescribeCertificateInput, _ ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
	return m.out, m.err
}

func TestFetchACMCertificate_Success(t *testing.T) {
	client := &mockACMClient{
		out: &acm.DescribeCertificateOutput{
			Certificate: &types.CertificateDetail{
				CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/abc"),
				DomainName:     aws.String("example.com"),
				Status:         types.CertificateStatusIssued,
				Type:           types.CertificateTypeAmazonIssued,
			},
		},
	}
	attrs, err := FetchACMCertificate(context.Background(), client, "arn:aws:acm:us-east-1:123456789012:certificate/abc")
	require.NoError(t, err)
	assert.Equal(t, "example.com", attrs["domain"])
	assert.Equal(t, "ISSUED", attrs["status"])
	assert.Equal(t, "AMAZON_ISSUED", attrs["type"])
}

func TestFetchACMCertificate_EmptyID(t *testing.T) {
	_, err := FetchACMCertificate(context.Background(), &mockACMClient{}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchACMCertificate_Error(t *testing.T) {
	client := &mockACMClient{err: errors.New("api error")}
	_, err := FetchACMCertificate(context.Background(), client, "arn:aws:acm:us-east-1:123:certificate/x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "describing ACM certificate")
}

func TestFetchACMCertificate_NilCertificate(t *testing.T) {
	client := &mockACMClient{out: &acm.DescribeCertificateOutput{}}
	_, err := FetchACMCertificate(context.Background(), client, "arn:aws:acm:us-east-1:123:certificate/x")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
