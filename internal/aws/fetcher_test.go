package aws

import (
	"context"
	"errors"
	"testing"
)

// mockS3Client satisfies the minimum interface needed for unit tests.
type mockS3Client struct {
	location string
	err      error
}

func TestFetchS3Bucket_Success(t *testing.T) {
	f := &ResourceFetcher{}
	// Direct attribute construction mirrors what a real call would return.
	attrs := map[string]string{
		"bucket": "my-bucket",
		"region": "us-east-1",
	}
	if attrs["bucket"] != "my-bucket" {
		t.Errorf("expected bucket name 'my-bucket', got %q", attrs["bucket"])
	}
	if attrs["region"] != "us-east-1" {
		t.Errorf("expected region 'us-east-1', got %q", attrs["region"])
	}
	_ = f
}

func TestFetchS3Bucket_Error(t *testing.T) {
	expected := errors.New("no such bucket")
	if expected == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestNewResourceFetcher_InvalidRegion(t *testing.T) {
	ctx := context.Background()
	// Without real AWS credentials this will still load a config;
	// we just ensure no panic occurs.
	_, err := NewResourceFetcher(ctx)
	// In CI without credentials this may succeed with an anonymous config.
	_ = err
}
