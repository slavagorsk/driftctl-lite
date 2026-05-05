package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSecretsManagerClient struct {
	out *secretsmanager.DescribeSecretOutput
	err error
}

func (m *mockSecretsManagerClient) DescribeSecret(_ context.Context, _ *secretsmanager.DescribeSecretInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error) {
	return m.out, m.err
}

func TestFetchSecretsManagerSecret_Success(t *testing.T) {
	client := &mockSecretsManagerClient{
		out: &secretsmanager.DescribeSecretOutput{
			Name:            aws.String("my-secret"),
			ARN:             aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret"),
			Description:     aws.String("test secret"),
			KmsKeyId:        aws.String("alias/aws/secretsmanager"),
			RotationEnabled: aws.Bool(true),
		},
	}

	attrs, err := FetchSecretsManagerSecret(context.Background(), client, "my-secret")
	require.NoError(t, err)
	assert.Equal(t, "my-secret", attrs["name"])
	assert.Equal(t, "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret", attrs["arn"])
	assert.Equal(t, "test secret", attrs["description"])
	assert.Equal(t, "alias/aws/secretsmanager", attrs["kms_key_id"])
	assert.Equal(t, "true", attrs["rotation_enabled"])
}

func TestFetchSecretsManagerSecret_EmptyID(t *testing.T) {
	client := &mockSecretsManagerClient{}
	_, err := FetchSecretsManagerSecret(context.Background(), client, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchSecretsManagerSecret_Error(t *testing.T) {
	client := &mockSecretsManagerClient{
		err: errors.New("ResourceNotFoundException"),
	}
	_, err := FetchSecretsManagerSecret(context.Background(), client, "missing-secret")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "describing secret")
}

func TestFetchSecretsManagerSecret_NoRotation(t *testing.T) {
	client := &mockSecretsManagerClient{
		out: &secretsmanager.DescribeSecretOutput{
			Name:            aws.String("plain-secret"),
			ARN:             aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:plain-secret"),
			RotationEnabled: aws.Bool(false),
		},
	}

	attrs, err := FetchSecretsManagerSecret(context.Background(), client, "plain-secret")
	require.NoError(t, err)
	assert.Equal(t, "false", attrs["rotation_enabled"])
	assert.NotContains(t, attrs, "description")
	assert.NotContains(t, attrs, "kms_key_id")
}
