package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Instance represents the fetched state of an AWS EC2 instance.
type EC2Instance struct {
	InstanceID   string
	InstanceType string
	State        string
	AMI          string
	Tags         map[string]string
}

// FetchEC2Instance retrieves the current state of an EC2 instance by its ID
// from AWS and returns a map of attributes suitable for drift comparison.
func FetchEC2Instance(ctx context.Context, cfg aws.Config, instanceID string) (map[string]interface{}, error) {
	if instanceID == "" {
		return nil, fmt.Errorf("instance ID must not be empty")
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	resp, err := client.DescribeInstances(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("describing EC2 instance %q: %w", instanceID, err)
	}

	instance, err := extractInstance(resp.Reservations, instanceID)
	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"id":            aws.ToString(instance.InstanceId),
		"instance_type": string(instance.InstanceType),
		"ami":           aws.ToString(instance.ImageId),
		"state":         string(instance.State.Name),
		"tags":          flattenEC2Tags(instance.Tags),
	}

	return attrs, nil
}

// extractInstance finds and returns the first instance matching instanceID
// across all reservations in the DescribeInstances response.
func extractInstance(reservations []types.Reservation, instanceID string) (*types.Instance, error) {
	for _, r := range reservations {
		for i := range r.Instances {
			if aws.ToString(r.Instances[i].InstanceId) == instanceID {
				return &r.Instances[i], nil
			}
		}
	}
	return nil, fmt.Errorf("EC2 instance %q not found", instanceID)
}

// flattenEC2Tags converts the AWS SDK tag slice into a simple string map.
func flattenEC2Tags(tags []types.Tag) map[string]string {
	result := make(map[string]string, len(tags))
	for _, t := range tags {
		result[aws.ToString(t.Key)] = aws.ToString(t.Value)
	}
	return result
}
