package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// SQSClient defines the interface for SQS operations used in drift detection.
type SQSClient interface {
	GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)
}

// FetchSQSQueue retrieves attributes for an SQS queue by its URL.
// Returns a map of attributes or an error if the queue is not found or the call fails.
func FetchSQSQueue(ctx context.Context, client SQSClient, queueURL string) (map[string]string, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("queue URL must not be empty")
	}

	out, err := client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []sqstypes.QueueAttributeName{"All"},
	})
	if err != nil {
		return nil, fmt.Errorf("fetching SQS queue %q: %w", queueURL, err)
	}

	attrs := make(map[string]string, len(out.Attributes))
	for k, v := range out.Attributes {
		attrs[string(k)] = v
	}
	return attrs, nil
}
