package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSClient defines the interface for SNS operations used in drift detection.
type SNSClient interface {
	GetTopicAttributes(ctx context.Context, params *sns.GetTopicAttributesInput, optFns ...func(*sns.Options)) (*sns.GetTopicAttributesOutput, error)
}

// FetchSNSTopic retrieves attributes for an SNS topic by its ARN.
func FetchSNSTopic(ctx context.Context, client SNSClient, topicARN string) (map[string]string, error) {
	if topicARN == "" {
		return nil, fmt.Errorf("SNS topic ARN must not be empty")
	}

	out, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicARN),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching SNS topic %q: %w", topicARN, err)
	}

	attrs := make(map[string]string, len(out.Attributes))
	for k, v := range out.Attributes {
		attrs[k] = v
	}
	return attrs, nil
}
