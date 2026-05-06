package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsELBv2Type(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"aws_alb", true},
		{"aws_lb", true},
		{"aws_s3_bucket", false},
		{"aws_instance", false},
		{"", false},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, isELBv2Type(tc.input))
		})
	}
}

func TestFetchELBv2Resource_UnsupportedType(t *testing.T) {
	cfg := ResourceFetcherConfig{Region: "us-east-1"}
	_, err := fetchELBv2Resource(nil, cfg, "aws_rds_instance", "some-arn")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported resource type")
}

func TestFetchELBv2Resource_EmptyID(t *testing.T) {
	cfg := ResourceFetcherConfig{Region: "us-east-1"}
	_, err := fetchELBv2Resource(nil, cfg, "aws_lb", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchELBv2Resource_DelegatesToFetchELBv2LoadBalancer(t *testing.T) {
	// This test validates that fetchELBv2Resource attempts to build a real client
	// and will fail at the AWS config level in a unit test environment,
	// confirming the delegation path is reached past validation.
	cfg := ResourceFetcherConfig{Region: "us-east-1"}
	_, err := fetchELBv2Resource(nil, cfg, "aws_alb", "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc")
	// Error is expected (no real AWS creds), but it should not be a validation error
	assert.Error(t, err)
	assert.NotContains(t, err.Error(), "unsupported resource type")
	assert.NotContains(t, err.Error(), "must not be empty")
}
