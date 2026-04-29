package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIAMType(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{"aws_iam_user", true},
		{"aws_iam_role", true},
		{"aws_s3_bucket", false},
		{"aws_instance", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			result := isIAMType(tt.resourceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFetchIAMResource_UnsupportedType(t *testing.T) {
	cfg := validTestConfig(t)
	attrs, err := fetchIAMResource(cfg, "aws_iam_policy", "my-policy")
	assert.Nil(t, attrs)
	assert.ErrorContains(t, err, "unsupported IAM resource type")
}

func TestFetchIAMResource_EmptyID(t *testing.T) {
	cfg := validTestConfig(t)

	attrs, err := fetchIAMResource(cfg, "aws_iam_user", "")
	assert.Nil(t, attrs)
	assert.ErrorContains(t, err, "resource ID is empty")

	attrs, err = fetchIAMResource(cfg, "aws_iam_role", "")
	assert.Nil(t, attrs)
	assert.ErrorContains(t, err, "resource ID is empty")
}

func TestFetchIAMResource_DelegatesToFetchIAMUser(t *testing.T) {
	cfg := validTestConfig(t)
	// With a non-empty ID, the function should attempt to call FetchIAMUser.
	// Since we're using real AWS config in tests, we expect an AWS API error,
	// not an argument validation error.
	_, err := fetchIAMResource(cfg, "aws_iam_user", "test-user")
	assert.Error(t, err)
	// Should not be an "unsupported" or "empty ID" error
	assert.NotContains(t, err.Error(), "unsupported IAM resource type")
	assert.NotContains(t, err.Error(), "resource ID is empty")
}

func TestFetchIAMResource_DelegatesToFetchIAMRole(t *testing.T) {
	cfg := validTestConfig(t)
	_, err := fetchIAMResource(cfg, "aws_iam_role", "test-role")
	assert.Error(t, err)
	assert.NotContains(t, err.Error(), "unsupported IAM resource type")
	assert.NotContains(t, err.Error(), "resource ID is empty")
}
