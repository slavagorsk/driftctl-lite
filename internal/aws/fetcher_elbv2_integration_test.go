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

// TestFetchELBv2LoadBalancer_Integration requires real AWS credentials and a valid
// load balancer ARN set via environment variables.
//
// Run with: go test -tags integration ./internal/aws/...
func TestFetchELBv2LoadBalancer_Integration(t *testing.T) {
	arn := os.Getenv("TEST_LB_ARN")
	if arn == "" {
		t.Skip("TEST_LB_ARN not set, skipping integration test")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	cfg := ResourceFetcherConfig{Region: region}
	client, err := elbv2ClientFromConfig(cfg)
	require.NoError(t, err)

	attrs, err := FetchELBv2LoadBalancer(context.Background(), client, arn)
	require.NoError(t, err)
	assert.NotEmpty(t, attrs["name"])
	assert.NotEmpty(t, attrs["dns"])
	assert.Equal(t, arn, attrs["arn"])
}
