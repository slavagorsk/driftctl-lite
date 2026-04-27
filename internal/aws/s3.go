package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/snyk/driftctl-lite/internal/drift"
)

type s3Client interface {
	HeadBucket(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	GetBucketTagging(ctx context.Context, input *s3.GetBucketTaggingInput, opts ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

// FetchS3Bucket retrieves the live state of an S3 bucket by name.
func FetchS3Bucket(ctx context.Context, client s3Client, id string) (*drift.Resource, error) {
	if id == "" {
		return nil, errors.New("s3: bucket id must not be empty")
	}

	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &id,
	})
	if err != nil {
		var notFound *types.NoSuchBucket
		if errors.As(err, &notFound) {
			return nil, fmt.Errorf("s3: bucket %q not found", id)
		}
		return nil, fmt.Errorf("s3: head bucket %q: %w", id, err)
	}

	attrs := map[string]string{
		"bucket": id,
	}

	tagsOut, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: &id,
	})
	if err == nil {
		for _, tag := range tagsOut.TagSet {
			if tag.Key != nil && tag.Value != nil {
				attrs[fmt.Sprintf("tags.%s", *tag.Key)] = *tag.Value
			}
		}
	}
	// tagging errors (e.g. access denied) are non-fatal; we proceed without tags

	return &drift.Resource{
		Type:       "aws_s3_bucket",
		ID:         id,
		Attributes: attrs,
	}, nil
}
