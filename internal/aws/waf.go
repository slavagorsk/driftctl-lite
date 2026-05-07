package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

// WAFv2Client defines the interface for WAFv2 operations used in drift detection.
type WAFv2Client interface {
	GetWebACL(ctx context.Context, params *wafv2.GetWebACLInput, optFns ...func(*wafv2.Options)) (*wafv2.GetWebACLOutput, error)
}

// FetchWAFv2WebACL retrieves a WAFv2 Web ACL by its ID and name from AWS.
// The id parameter is expected in the format "id/name/scope" (e.g. "abc123/my-acl/REGIONAL").
func FetchWAFv2WebACL(ctx context.Context, client WAFv2Client, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("wafv2 web acl id must not be empty")
	}

	var aclID, name, scopeStr string
	_, err := fmt.Sscanf(id, "%s", &aclID)
	if err != nil {
		return nil, fmt.Errorf("invalid wafv2 web acl id format: %s", id)
	}

	// Parse composite id: id/name/scope
	parts := splitN(id, "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("wafv2 web acl id must be in format id/name/scope, got: %s", id)
	}
	aclID, name, scopeStr = parts[0], parts[1], parts[2]

	var scope types.Scope
	switch scopeStr {
	case "REGIONAL":
		scope = types.ScopeRegional
	case "CLOUDFRONT":
		scope = types.ScopeCloudfront
	default:
		return nil, fmt.Errorf("unsupported wafv2 scope: %s", scopeStr)
	}

	out, err := client.GetWebACL(ctx, &wafv2.GetWebACLInput{
		Id:    aws.String(aclID),
		Name:  aws.String(name),
		Scope: scope,
	})
	if err != nil {
		return nil, fmt.Errorf("fetching wafv2 web acl %s: %w", id, err)
	}
	if out.WebACL == nil {
		return nil, fmt.Errorf("wafv2 web acl %s not found", id)
	}

	attrs := map[string]string{
		"id":          aws.ToString(out.WebACL.Id),
		"name":        aws.ToString(out.WebACL.Name),
		"description": aws.ToString(out.WebACL.Description),
		"arn":         aws.ToString(out.WebACL.ARN),
	}
	return attrs, nil
}

// splitN splits s by sep up to n substrings.
func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
