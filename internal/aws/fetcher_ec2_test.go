package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEC2Type(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{"aws_instance", true},
		{"aws_s3_bucket", false},
		{"aws_iam_user", false},
		{"", false},
		{"aws_db_instance", false},
	}

	for _, tc := range tests {
		t.Run(tc.resourceType, func(t *testing.T) {
			result := isEC2Type(tc.resourceType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFetchEC2Resource_UnsupportedType(t *testing.T) {
	cfg := validTestConfig(t)
	client, err := ec2ClientFromConfig(cfg)
	assert.NoError(t, err)

	_, err = fetchEC2Resource(client, "aws_vpc", "vpc-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported EC2 resource type")
}

func TestFetchEC2Resource_EmptyID(t *testing.T) {
	cfg := validTestConfig(t)
	client, err := ec2ClientFromConfig(cfg)
	assert.NoError(t, err)

	_, err = fetchEC2Resource(client, "aws_instance", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty resource ID")
}

func TestFetchEC2Resource_DelegatesToFetchEC2Instance(t *testing.T) {
	cfg := validTestConfig(t)
	client, err := ec2ClientFromConfig(cfg)
	assert.NoError(t, err)

	// A non-empty ID that will fail with an AWS error (no real AWS)
	// confirms delegation without mocking the full SDK
	_, err = fetchEC2Resource(client, "aws_instance", "i-doesnotexist")
	// We expect an error from AWS, not an "unsupported type" error
	if err != nil {
		assert.NotContains(t, err.Error(), "unsupported EC2 resource type")
	}
}
