package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// iamClientFromConfig builds a real IAM client using the fetcher's AWS config.
func iamClientFromConfig(f *ResourceFetcher) IAMClient {
	return iam.NewFromConfig(f.cfg)
}

// fetchIAMResource dispatches to the appropriate IAM fetch function based on
// the resource type suffix ("aws_iam_user" or "aws_iam_role").
func fetchIAMResource(ctx context.Context, f *ResourceFetcher, resourceType, id string) (map[string]string, error) {
	client := iamClientFromConfig(f)

	switch {
	case strings.HasSuffix(resourceType, "iam_user"):
		return FetchIAMUser(ctx, client, id)
	case strings.HasSuffix(resourceType, "iam_role"):
		return FetchIAMRole(ctx, client, id)
	default:
		return nil, fmt.Errorf("unsupported IAM resource type: %q", resourceType)
	}
}

// isIAMType returns true when the resource type belongs to the IAM service.
func isIAMType(resourceType string) bool {
	return strings.Contains(resourceType, "iam_user") ||
		strings.Contains(resourceType, "iam_role")
}
