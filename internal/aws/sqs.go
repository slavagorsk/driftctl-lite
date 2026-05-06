package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSClient defines the interface for SQS operations used in fetching.
type SQSClient interface {
	GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)
	ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
}

// FetchSQSQueue retrieves attributes for an SQS queue by its URL.
func FetchSQSQueue(ctx context.Context, client SQSClient, queueURL string) (map[string]string, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("queue URL must not be empty")
	}

	out, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []sqstypes.QueueAttributeName{"All"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get SQS queue attributes for %q: %w", queueURL, err)
	}

	attrs := make(map[string]string, len(out.Attributes))
	for k, v := range out.Attributes {
		attrs[string(k)] = v
	}
	return attrs, nil
}

// ListSQSQueues returns all SQS queue URLs, optionally filtered by a prefix.
// An empty prefix returns all queues in the account and region.
func ListSQSQueues(ctx context.Context, client SQSClient, prefix string) ([]string, error) {
	input := &sqs.ListQueuesInput{}
	if prefix != "" {
		input.QueueNamePrefix = aws.String(prefix)
	}

	out, err := client.ListQueues(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list SQS queues: %w", err)
	}

	return out.QueueUrls, nil
}
