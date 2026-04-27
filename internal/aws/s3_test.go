package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type mockS3Client struct {
	headOutput *s3.HeadBucketOutput
	headErr    error
	tagsOutput *s3.GetBucketTaggingOutput
	tagsErr    error
}

func (m *mockS3Client) HeadBucket(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return m.headOutput, m.headErr
}

func (m *mockS3Client) GetBucketTagging(ctx context.Context, input *s3.GetBucketTaggingInput, opts ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return m.tagsOutput, m.tagsErr
}

func TestFetchS3Bucket_Success(t *testing.T) {
	client := &mockS3Client{
		headOutput: &s3.HeadBucketOutput{},
		tagsOutput: &s3.GetBucketTaggingOutput{
			TagSet: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("prod")},
			},
		},
	}
	res, err := FetchS3Bucket(context.Background(), client, "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "my-bucket" {
		t.Errorf("expected ID my-bucket, got %s", res.ID)
	}
	if res.Attributes["tags.env"] != "prod" {
		t.Errorf("expected tag env=prod, got %v", res.Attributes["tags.env"])
	}
}

func TestFetchS3Bucket_NotFound(t *testing.T) {
	client := &mockS3Client{
		headErr: &types.NoSuchBucket{},
	}
	_, err := FetchS3Bucket(context.Background(), client, "missing-bucket")
	if err == nil {
		t.Fatal("expected error for missing bucket")
	}
}

func TestFetchS3Bucket_EmptyID(t *testing.T) {
	client := &mockS3Client{}
	_, err := FetchS3Bucket(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty bucket ID")
	}
}

func TestFetchS3Bucket_TaggingError_Ignored(t *testing.T) {
	client := &mockS3Client{
		headOutput: &s3.HeadBucketOutput{},
		tagsErr:    errors.New("access denied"),
	}
	res, err := FetchS3Bucket(context.Background(), client, "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "my-bucket" {
		t.Errorf("expected ID my-bucket, got %s", res.ID)
	}
}
