package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client defines the subset of the S3 API used for drift detection.
type S3Client interface {
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
	GetBucketVersioning(ctx context.Context, params *s3.GetBucketVersioningInput, optFns ...func(*s3.Options)) (*s3.GetBucketVersioningOutput, error)
}

// S3BucketAttributes holds the live attributes fetched from AWS for an S3 bucket.
type S3BucketAttributes struct {
	Bucket     string
	Region     string
	Versioning string
	Tags       map[string]string
}

// FetchS3Bucket retrieves the current state of an S3 bucket from AWS.
// It returns a map of attribute name to value suitable for drift comparison.
func FetchS3Bucket(ctx context.Context, client S3Client, bucketName string) (map[string]interface{}, error) {
	if bucketName == "" {
		return nil, fmt.Errorf("bucket name must not be empty")
	}

	attrs := make(map[string]interface{})
	attrs["bucket"] = bucketName

	// Fetch bucket region / location constraint.
	locOut, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("get bucket location for %q: %w", bucketName, err)
	}
	region := string(locOut.LocationConstraint)
	if region == "" {
		// Buckets in us-east-1 return an empty location constraint.
		region = "us-east-1"
	}
	attrs["region"] = region

	// Fetch versioning status.
	verOut, err := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("get bucket versioning for %q: %w", bucketName, err)
	}
	versioning := string(verOut.Status)
	if versioning == "" {
		versioning = "Disabled"
	}
	attrs["versioning"] = versioning

	// Fetch tags (best-effort: missing tags are treated as empty).
	tagOut, err := client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	})
	tags := map[string]string{}
	if err == nil && tagOut != nil {
		for _, t := range tagOut.TagSet {
			if t.Key != nil && t.Value != nil {
				tags[*t.Key] = *t.Value
			}
		}
	}
	attrs["tags"] = tags

	return attrs, nil
}
