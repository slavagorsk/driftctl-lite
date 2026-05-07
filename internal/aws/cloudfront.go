package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

// CloudFrontClient defines the interface for CloudFront operations.
type CloudFrontClient interface {
	GetDistribution(ctx context.Context, params *cloudfront.GetDistributionInput, optFns ...func(*cloudfront.Options)) (*cloudfront.GetDistributionOutput, error)
}

// FetchCloudFrontDistribution retrieves a CloudFront distribution by its ID.
func FetchCloudFrontDistribution(ctx context.Context, client CloudFrontClient, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("cloudfront distribution id must not be empty")
	}

	out, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{
		Id: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching cloudfront distribution %q: %w", id, err)
	}

	if out.Distribution == nil || out.Distribution.DistributionConfig == nil {
		return nil, fmt.Errorf("cloudfront distribution %q returned nil config", id)
	}

	cfg := out.Distribution.DistributionConfig
	attrs := map[string]string{
		"id":      id,
		"enabled": fmt.Sprintf("%v", aws.ToBool(cfg.Enabled)),
		"comment": aws.ToString(cfg.Comment),
	}

	if out.Distribution.DomainName != nil {
		attrs["domain_name"] = aws.ToString(out.Distribution.DomainName)
	}

	if cfg.DefaultCacheBehavior != nil {
		attrs["viewer_protocol_policy"] = string(cfg.DefaultCacheBehavior.ViewerProtocolPolicy)
	}

	return attrs, nil
}
