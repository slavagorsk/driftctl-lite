package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client defines the interface for S3 operations used in fetching.
type S3Client interface {
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

// mockS3Client is used in tests.
type mockS3Client struct {
	headOutput *s3.HeadBucketOutput
	headErr    error
	tagsOutput *s3.GetBucketTaggingOutput
	tagsErr    error
}

func (m *mockS3Client) HeadBucket(_ context.Context, _ *s3.HeadBucketInput, _ ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return m.headOutput, m.headErr
}

func (m *mockS3Client) GetBucketTagging(_ context.Context, _ *s3.GetBucketTaggingInput, _ ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	return m.tagsOutput, m.tagsErr
}

// FetchS3Bucket retrieves live attributes for an S3 bucket.
func FetchS3Bucket(ctx context.Context, client S3Client, bucketName string) (map[string]string, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("s3: bucketName must not be empty")
	}

	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("s3: head bucket: %w", err)
	}

	attrs := map[string]string{
		"bucket": bucketName,
	}

	tagsOut, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})
	if err == nil && tagsOut != nil {
		for _, tag := range tagsOut.TagSet {
			attrs["tag:"+aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}
	}

	return attrs, nil
}
