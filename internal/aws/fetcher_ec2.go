package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// isEC2Type returns true if the given Terraform resource type is an EC2 resource.
func isEC2Type(resourceType string) bool {
	switch resourceType {
	case "aws_instance":
		return true
	}
	return false
}

// ec2ClientFromConfig constructs an EC2 client using the fetcher's AWS config.
func ec2ClientFromConfig(f *ResourceFetcher) *ec2.Client {
	return ec2.NewFromConfig(f.cfg)
}

// fetchEC2Resource dispatches EC2 resource fetching based on the resource type.
// It returns the fetched attributes as a map or an error if the fetch fails.
func fetchEC2Resource(f *ResourceFetcher, resourceType, id string) (map[string]interface{}, error) {
	if id == "" {
		return nil, fmt.Errorf("resource ID must not be empty for type %q", resourceType)
	}

	switch resourceType {
	case "aws_instance":
		client := ec2ClientFromConfig(f)
		return FetchEC2Instance(client, id)
	default:
		return nil, fmt.Errorf("unsupported EC2 resource type: %q", resourceType)
	}
}

// Ensure aws import is used (aws.String is referenced in ec2.go helpers).
var _ = aws.String
