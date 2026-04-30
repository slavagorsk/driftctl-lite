package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsSQSType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"aws_sqs_queue", true},
		{"aws_sns_topic", false},
		{"aws_s3_bucket", false},
		{"", false},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expected, isSQSType(tc.input))
		})
	}
}

func TestFetchSQSResource_UnsupportedType(t *testing.T) {
	_, err := fetchSQSResource(context.Background(), "us-east-1", "aws_lambda_function", "some-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported SQS resource type")
}

func TestFetchSQSResource_EmptyID(t *testing.T) {
	_, err := fetchSQSResource(context.Background(), "us-east-1", "aws_sqs_queue", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchSQSResource_DelegatesToFetchSQSQueue(t *testing.T) {
	mock := &mockSQSClient{
		attributes: map[string]string{
			"QueueArn":                 "arn:aws:sqs:us-east-1:123456789012:my-queue",
			"ApproximateNumberOfMessages": "0",
		},
	}
	attrs, err := FetchSQSQueue(context.Background(), mock, "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue")
	require.NoError(t, err)
	assert.Equal(t, "arn:aws:sqs:us-east-1:123456789012:my-queue", attrs["QueueArn"])
}
