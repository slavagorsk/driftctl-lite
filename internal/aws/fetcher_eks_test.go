package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEKSType(t *testing.T) {
	assert.True(t, isEKSType("aws_eks_cluster"))
	assert.False(t, isEKSType("aws_s3_bucket"))
	assert.False(t, isEKSType(""))
}

func TestFetchEKSResource_UnsupportedType(t *testing.T) {
	_, err := fetchEKSResource(context.Background(), "us-east-1", "aws_lambda_function", "my-cluster")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported EKS resource type")
}

func TestFetchEKSResource_EmptyID(t *testing.T) {
	_, err := fetchEKSResource(context.Background(), "us-east-1", "aws_eks_cluster", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchEKSResource_DelegatesToFetchEKSCluster(t *testing.T) {
	origNewConfig := newEKSConfigFromRegion
	origClientFromConfig := eksClientFromConfig
	defer func() {
		newEKSConfigFromRegion = origNewConfig
		eksClientFromConfig = origClientFromConfig
	}()

	newEKSConfigFromRegion = func(_ context.Context, _ string) (aws.Config, error) {
		return aws.Config{}, nil
	}
	eksClientFromConfig = func(_ aws.Config) EKSClient {
		return &mockEKSClient{
			out: &eks.DescribeClusterOutput{
				Cluster: &types.Cluster{
					Name:    aws.String("my-cluster"),
					Status:  types.ClusterStatusActive,
					Version: aws.String("1.29"),
				},
			},
			err: nil,
		}
	}

	attrs, err := fetchEKSResource(context.Background(), "us-east-1", "aws_eks_cluster", "my-cluster")
	require.NoError(t, err)
	assert.Equal(t, "my-cluster", attrs["name"])
	assert.Equal(t, "ACTIVE", attrs["status"])
}

func TestFetchEKSResource_PropagatesClientError(t *testing.T) {
	origNewConfig := newEKSConfigFromRegion
	origClientFromConfig := eksClientFromConfig
	defer func() {
		newEKSConfigFromRegion = origNewConfig
		eksClientFromConfig = origClientFromConfig
	}()

	newEKSConfigFromRegion = func(_ context.Context, _ string) (aws.Config, error) {
		return aws.Config{}, nil
	}
	eksClientFromConfig = func(_ aws.Config) EKSClient {
		return &mockEKSClient{err: errors.New("throttled")}
	}

	_, err := fetchEKSResource(context.Background(), "us-east-1", "aws_eks_cluster", "my-cluster")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "throttled")
}
