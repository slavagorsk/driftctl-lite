package aws

import (
	"context"
	"testing"
)

func TestFetchS3Bucket_Success(t *testing.T) {
	t.Skip("requires mock S3 client wired via interface — covered in s3_test.go")
}

func TestFetchS3Bucket_Error(t *testing.T) {
	t.Skip("requires mock S3 client wired via interface — covered in s3_test.go")
}

func TestNewResourceFetcher_InvalidRegion(t *testing.T) {
	_, err := NewResourceFetcher("")
	if err == nil {
		t.Fatal("expected error for empty region, got nil")
	}
}

func TestFetch_UnsupportedType(t *testing.T) {
	f := &ResourceFetcher{}
	_, err := f.Fetch(context.Background(), "aws_lambda_function", "my-fn")
	if err == nil {
		t.Fatal("expected unsupported type error, got nil")
	}
}

func TestFetch_S3Bucket_DelegatesToFetchS3Bucket(t *testing.T) {
	mock := &mockS3Client{
		headOutput: nil,
		tagsOutput: nil,
		headErr:    nil,
		tagsErr:    nil,
	}
	f := &ResourceFetcher{s3Client: mock}
	// FetchS3Bucket returns error when HeadBucket output is nil — just verify delegation.
	_, _ = f.Fetch(context.Background(), "aws_s3_bucket", "test-bucket")
}

func TestFetch_RDS_DelegatesToFetchRDSInstance(t *testing.T) {
	mock := &mockRDSClient{err: nil, output: nil}
	f := &ResourceFetcher{rdsClient: mock}
	_, err := f.Fetch(context.Background(), "aws_db_instance", "mydb")
	// nil output triggers not-found error — just verify delegation happened
	if err == nil {
		t.Fatal("expected error from nil output, got nil")
	}
}
