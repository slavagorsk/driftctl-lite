package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type mockSQSClient struct {
	attributes map[sqstypes.QueueAttributeName]string
	err        error
}

func (m *mockSQSClient) GetQueueAttributes(_ context.Context, _ *sqs.GetQueueAttributesInput, _ ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &sqs.GetQueueAttributesOutput{Attributes: m.attributes}, nil
}

func TestFetchSQSQueue_Success(t *testing.T) {
	client := &mockSQSClient{
		attributes: map[sqstypes.QueueAttributeName]string{
			"QueueArn":                 "arn:aws:sqs:us-east-1:123456789012:my-queue",
			"ApproximateNumberOfMessages": "0",
		},
	}

	attrs, err := FetchSQSQueue(context.Background(), client, "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attrs["QueueArn"] != "arn:aws:sqs:us-east-1:123456789012:my-queue" {
		t.Errorf("unexpected QueueArn: %s", attrs["QueueArn"])
	}
}

func TestFetchSQSQueue_EmptyID(t *testing.T) {
	client := &mockSQSClient{}
	_, err := FetchSQSQueue(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty queue URL, got nil")
	}
}

func TestFetchSQSQueue_Error(t *testing.T) {
	client := &mockSQSClient{err: errors.New("access denied")}
	_, err := FetchSQSQueue(context.Background(), client, "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errors.New("access denied")) && err.Error() == "" {
		t.Errorf("unexpected error message: %v", err)
	}
}
