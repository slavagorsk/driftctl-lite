package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type mockSNSClient struct {
	attrs map[string]string
	err   error
}

func (m *mockSNSClient) GetTopicAttributes(_ context.Context, _ *sns.GetTopicAttributesInput, _ ...func(*sns.Options)) (*sns.GetTopicAttributesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &sns.GetTopicAttributesOutput{Attributes: m.attrs}, nil
}

func TestFetchSNSTopic_Success(t *testing.T) {
	client := &mockSNSClient{
		attrs: map[string]string{
			"TopicArn":         "arn:aws:sns:us-east-1:123456789012:my-topic",
			"DisplayName":      "my-topic",
			"SubscriptionsConfirmed": "2",
		},
	}
	attrs, err := FetchSNSTopic(context.Background(), client, "arn:aws:sns:us-east-1:123456789012:my-topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["DisplayName"] != "my-topic" {
		t.Errorf("expected DisplayName=my-topic, got %q", attrs["DisplayName"])
	}
}

func TestFetchSNSTopic_EmptyID(t *testing.T) {
	client := &mockSNSClient{}
	_, err := FetchSNSTopic(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty ARN")
	}
}

func TestFetchSNSTopic_Error(t *testing.T) {
	client := &mockSNSClient{err: errors.New("topic not found")}
	_, err := FetchSNSTopic(context.Background(), client, "arn:aws:sns:us-east-1:123456789012:missing")
	if err == nil {
		t.Fatal("expected error from AWS client")
	}
}
