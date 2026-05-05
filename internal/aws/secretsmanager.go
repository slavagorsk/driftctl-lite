package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerClient defines the interface used for fetching Secrets Manager secrets.
type SecretsManagerClient interface {
	DescribeSecret(ctx context.Context, params *secretsmanager.DescribeSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error)
}

// FetchSecretsManagerSecret retrieves metadata for a Secrets Manager secret by ARN or name.
func FetchSecretsManagerSecret(ctx context.Context, client SecretsManagerClient, id string) (map[string]string, error) {
	if id == "" {
		return nil, fmt.Errorf("secret id must not be empty")
	}

	out, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("describing secret %q: %w", id, err)
	}

	attrs := map[string]string{
		"name": aws.ToString(out.Name),
		"arn":  aws.ToString(out.ARN),
	}

	if out.Description != nil {
		attrs["description"] = aws.ToString(out.Description)
	}
	if out.KmsKeyId != nil {
		attrs["kms_key_id"] = aws.ToString(out.KmsKeyId)
	}
	if out.RotationEnabled != nil && *out.RotationEnabled {
		attrs["rotation_enabled"] = "true"
	} else {
		attrs["rotation_enabled"] = "false"
	}

	return attrs, nil
}
