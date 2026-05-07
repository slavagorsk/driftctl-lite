package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
)

// ACMClient defines the interface for ACM operations used in fetching.
type ACMClient interface {
	DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error)
}

// FetchACMCertificate retrieves ACM certificate attributes by ARN.
func FetchACMCertificate(ctx context.Context, client ACMClient, arn string) (map[string]string, error) {
	if arn == "" {
		return nil, fmt.Errorf("certificate ARN must not be empty")
	}

	out, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	})
	if err != nil {
		return nil, fmt.Errorf("describing ACM certificate %q: %w", arn, err)
	}

	if out.Certificate == nil {
		return nil, fmt.Errorf("ACM certificate %q not found", arn)
	}

	attrs := map[string]string{
		"arn":    aws.ToString(out.Certificate.CertificateArn),
		"domain": aws.ToString(out.Certificate.DomainName),
		"status": string(out.Certificate.Status),
	}

	if out.Certificate.Type != "" {
		attrs["type"] = string(out.Certificate.Type)
	}

	return attrs, nil
}
