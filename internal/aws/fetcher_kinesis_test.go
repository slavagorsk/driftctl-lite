package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type mockKinesisClient struct {
	out *kinesis.DescribeStreamOutput
	err error
}

func (m *mockKinesisClient) DescribeStream(_ context.Context, _ *kinesis.DescribeStreamInput, _ ...func(*kinesis.Options)) (*kinesis.DescribeStreamOutput, error) {
	return m.out, m.err
}

func TestIsKinesisType(t *testing.T) {
	if !isKinesisType("aws_kinesis_stream") {
		t.Error("expected aws_kinesis_stream to be a kinesis type")
	}
	if isKinesisType("aws_s3_bucket") {
		t.Error("expected aws_s3_bucket not to be a kinesis type")
	}
}

func TestFetchKinesisResource_UnsupportedType(t *testing.T) {
	_, err := fetchKinesisResource(context.Background(), "us-east-1", "aws_sqs_queue", "my-stream")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestFetchKinesisResource_EmptyID(t *testing.T) {
	_, err := fetchKinesisResource(context.Background(), "us-east-1", "aws_kinesis_stream", "")
	if err == nil {
		t.Fatal("expected error for empty resource ID")
	}
}

func TestFetchKinesisResource_DelegatesToFetchKinesisStream(t *testing.T) {
	shards := []types.Shard{{ShardId: aws.String("shardId-000000000000")}}
	mock := &mockKinesisClient{
		out: &kinesis.DescribeStreamOutput{
			StreamDescription: &types.StreamDescription{
				StreamName:           aws.String("my-stream"),
				StreamARN:            aws.String("arn:aws:kinesis:us-east-1:123456789012:stream/my-stream"),
				StreamStatus:         types.StreamStatusActive,
				Shards:               shards,
				RetentionPeriodHours: aws.Int32(24),
			},
		},
	}

	attrs, err := FetchKinesisStream(context.Background(), mock, "my-stream")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["stream_name"] != "my-stream" {
		t.Errorf("expected stream_name=my-stream, got %s", attrs["stream_name"])
	}
	if attrs["shard_count"] != "1" {
		t.Errorf("expected shard_count=1, got %s", attrs["shard_count"])
	}
	if attrs["retention_period_hours"] != "24" {
		t.Errorf("expected retention_period_hours=24, got %s", attrs["retention_period_hours"])
	}
}

func TestFetchKinesisStream_Error(t *testing.T) {
	mock := &mockKinesisClient{err: errors.New("api error")}
	_, err := FetchKinesisStream(context.Background(), mock, "my-stream")
	if err == nil {
		t.Fatal("expected error from api")
	}
}
