package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

// KinesisClient defines the interface used to fetch Kinesis stream details.
type KinesisClient interface {
	DescribeStream(ctx context.Context, params *kinesis.DescribeStreamInput, optFns ...func(*kinesis.Options)) (*kinesis.DescribeStreamOutput, error)
}

// FetchKinesisStream retrieves attributes of a Kinesis stream by name.
func FetchKinesisStream(ctx context.Context, client KinesisClient, streamName string) (map[string]string, error) {
	if streamName == "" {
		return nil, fmt.Errorf("kinesis stream name must not be empty")
	}

	out, err := client.DescribeStream(ctx, &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	})
	if err != nil {
		return nil, fmt.Errorf("describing kinesis stream %q: %w", streamName, err)
	}

	if out.StreamDescription == nil {
		return nil, fmt.Errorf("kinesis stream %q not found", streamName)
	}

	d := out.StreamDescription
	attrs := map[string]string{
		"stream_name":   aws.ToString(d.StreamName),
		"stream_arn":    aws.ToString(d.StreamARN),
		"stream_status": string(d.StreamStatus),
		"shard_count":   fmt.Sprintf("%d", len(d.Shards)),
	}

	if d.RetentionPeriodHours != nil {
		attrs["retention_period_hours"] = fmt.Sprintf("%d", aws.ToInt32(d.RetentionPeriodHours))
	}

	return attrs, nil
}
