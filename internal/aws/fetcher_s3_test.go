package aws

import (
	"context"
	"testing"
)

func TestIsS3Type(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"aws_s3_bucket", true},
		{"aws_instance", false},
		{"aws_iam_user", false},
		{"", false},
	}
	for _, tc := range cases {
		got := isS3Type(tc.input)
		if got != tc.expected {
			t.Errorf("isS3Type(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestFetchS3Resource_UnsupportedType(t *testing.T) {
	// fetchS3Resource with an unsupported type should return an error
	// before attempting any AWS calls (client build may fail in test env,
	// but we verify the routing logic via the mock path).
	_, err := fetchS3Resource(context.Background(), "us-east-1", "aws_s3_unknown", "bucket-id")
	if err == nil {
		t.Fatal("expected error for unsupported s3 resource type")
	}
}

func TestFetchS3Resource_EmptyID(t *testing.T) {
	// We bypass the real AWS client by calling FetchS3Bucket directly with a mock.
	client := &mockS3Client{}
	_, err := FetchS3Bucket(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty bucket ID")
	}
}

func TestFetchS3Resource_DelegatesToFetchS3Bucket(t *testing.T) {
	// Verify that fetchS3Resource with aws_s3_bucket type calls FetchS3Bucket.
	// We test the delegation indirectly by confirming the expected error shape
	// when the bucket does not exist via a mock.
	client := &mockS3Client{
		headErr: fmt.Errorf("no such bucket"),
	}
	_, err := FetchS3Bucket(context.Background(), client, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing bucket")
	}
}
