package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
)

// Route53Client defines the interface used for fetching hosted zones.
type Route53Client interface {
	GetHostedZone(ctx context.Context, params *route53.GetHostedZoneInput, optFns ...func(*route53.Options)) (*route53.GetHostedZoneOutput, error)
}

// FetchRoute53HostedZone retrieves a Route53 hosted zone by its ID and returns
// a map of attributes suitable for drift comparison.
func FetchRoute53HostedZone(ctx context.Context, client Route53Client, zoneID string) (map[string]string, error) {
	if zoneID == "" {
		return nil, fmt.Errorf("hosted zone ID must not be empty")
	}

	out, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: aws.String(zoneID),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching route53 hosted zone %q: %w", zoneID, err)
	}

	if out.HostedZone == nil {
		return nil, fmt.Errorf("hosted zone %q not found", zoneID)
	}

	attrs := map[string]string{
		"id":            aws.ToString(out.HostedZone.Id),
		"name":          aws.ToString(out.HostedZone.Name),
		"comment":       "",
		"private_zone":  fmt.Sprintf("%v", out.HostedZone.Config != nil && out.HostedZone.Config.PrivateZone),
		"record_count":  fmt.Sprintf("%d", out.HostedZone.ResourceRecordSetCount),
	}

	if out.HostedZone.Config != nil {
		attrs["comment"] = aws.ToString(out.HostedZone.Config.Comment)
	}

	return attrs, nil
}
