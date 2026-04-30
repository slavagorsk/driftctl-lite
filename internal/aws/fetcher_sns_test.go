package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func TestIsSNSType(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"aws_sns_topic", true},
		{"aws_s3_bucket", false},
		{"aws_iam_user", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := isSNSType(tc.input); got != tc.expected {
			t.Errorf("isSNSType(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestFetchSNSResource_UnsupportedType(t *testing.T) {
	client := &mockSNSClient{}
	_, err := fetchSNSResource(context.Background(), client, "aws_lambda_function", "some-arn")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestFetchSNSResource_EmptyID(t *testing.T) {
	client := &mockSNSClient{}
	_, err := fetchSNSResource(context.Background(), client, "aws_sns_topic", "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestFetchSNSResource_DelegatesToFetchSNSTopic(t *testing.T) {
	client := &mockSNSClient{
		attrs: map[string]string{
			"TopicArn":    "arn:aws:sns:us-east-1:123456789012:test-topic",
			"DisplayName": "test-topic",
		},
	}
	attrs, err := fetchSNSResource(context.Background(), client, "aws_sns_topic", "arn:aws:sns:us-east-1:123456789012:test-topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["DisplayName"] != "test-topic" {
		t.Errorf("expected DisplayName=test-topic, got %q", attrs["DisplayName"])
	}
}

var _ SNSClient = (*sns.Client)(nil)
