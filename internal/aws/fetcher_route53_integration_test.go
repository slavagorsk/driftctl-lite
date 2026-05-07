//go:build integration
// +build integration

package aws

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetchRoute53HostedZone_Integration exercises the real AWS Route53 API.
// Set ROUTE53_ZONE_ID to a valid hosted zone ID before running.
//
//	ROUTE53_ZONE_ID=Z1234567890 go test -tags integration ./internal/aws/...
func TestFetchRoute53HostedZone_Integration(t *testing.T) {
	zoneID := os.Getenv("ROUTE53_ZONE_ID")
	if zoneID == "" {
		t.Skip("ROUTE53_ZONE_ID not set, skipping integration test")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	client, err := route53ClientFromConfig(region)
	require.NoError(t, err)

	attrs, err := FetchRoute53HostedZone(context.Background(), client, zoneID)
	require.NoError(t, err)
	assert.NotEmpty(t, attrs["id"])
	assert.NotEmpty(t, attrs["name"])
}
