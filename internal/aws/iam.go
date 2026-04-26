package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMClient defines the subset of IAM API methods used.
type IAMClient interface {
	GetUser(ctx context.Context, params *iam.GetUserInput, optFns ...func(*iam.Options)) (*iam.GetUserOutput, error)
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
}

// FetchIAMUser retrieves an IAM user by name and returns its attributes.
func FetchIAMUser(ctx context.Context, client IAMClient, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("IAM user name must not be empty")
	}

	out, err := client.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching IAM user %q: %w", id, err)
	}
	if out.User == nil {
		return nil, fmt.Errorf("IAM user %q not found", id)
	}

	attrs := map[string]string{
		"user_name": aws.ToString(out.User.UserName),
		"arn":       aws.ToString(out.User.Arn),
		"user_id":   aws.ToString(out.User.UserId),
		"path":      aws.ToString(out.User.Path),
	}
	return attrs, nil
}

// FetchIAMRole retrieves an IAM role by name and returns its attributes.
func FetchIAMRole(ctx context.Context, client IAMClient, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("IAM role name must not be empty")
	}

	out, err := client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching IAM role %q: %w", id, err)
	}
	if out.Role == nil {
		return nil, fmt.Errorf("IAM role %q not found", id)
	}

	attrs := map[string]string{
		"role_name":            aws.ToString(out.Role.RoleName),
		"arn":                  aws.ToString(out.Role.Arn),
		"role_id":              aws.ToString(out.Role.RoleId),
		"path":                 aws.ToString(out.Role.Path),
		"assume_role_policy":   aws.ToString(out.Role.AssumeRolePolicyDocument),
	}
	return attrs, nil
}
