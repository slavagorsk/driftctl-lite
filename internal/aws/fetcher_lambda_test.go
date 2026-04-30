package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLambdaType(t *testing.T) {
	assert.True(t, isLambdaType("aws_lambda_function"))
	assert.False(t, isLambdaType("aws_s3_bucket"))
	assert.False(t, isLambdaType("aws_iam_user"))
	assert.False(t, isLambdaType(""))
}

func TestFetchLambdaResource_UnsupportedType(t *testing.T) {
	_, err := fetchLambdaResource("us-east-1", "aws_s3_bucket", "my-bucket")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported resource type for lambda fetcher")
}

func TestFetchLambdaResource_EmptyID(t *testing.T) {
	_, err := fetchLambdaResource("us-east-1", "aws_lambda_function", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "resource ID must not be empty")
}

func TestFetchLambdaResource_DelegatesToFetchLambdaFunction(t *testing.T) {
	// This test verifies that fetchLambdaResource delegates correctly.
	// A real AWS call would be made; we verify the error path for invalid region.
	_, err := fetchLambdaResource("invalid-region-xyz", "aws_lambda_function", "my-function")
	// We expect either a config error or an AWS API error, not a logic error.
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "unsupported resource type")
	assert.NotContains(t, err.Error(), "resource ID must not be empty")
}
