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

type mockEKSClient struct {
	out *eks.DescribeClusterOutput
	err error
}

func (m *mockEKSClient) DescribeCluster(_ context.Context, _ *eks.DescribeClusterInput, _ ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	return m.out, m.err
}

func TestFetchEKSCluster_Success(t *testing.T) {
	client := &mockEKSClient{
		out: &eks.DescribeClusterOutput{
			Cluster: &types.Cluster{
				Name:     aws.String("my-cluster"),
				Arn:      aws.String("arn:aws:eks:us-east-1:123456789012:cluster/my-cluster"),
				Status:   types.ClusterStatusActive,
				Version:  aws.String("1.29"),
				Endpoint: aws.String("https://EXAMPLE.gr7.us-east-1.eks.amazonaws.com"),
				RoleArn:  aws.String("arn:aws:iam::123456789012:role/eks-role"),
			},
		},
	}

	attrs, err := FetchEKSCluster(context.Background(), client, "my-cluster")
	require.NoError(t, err)
	assert.Equal(t, "my-cluster", attrs["name"])
	assert.Equal(t, "ACTIVE", attrs["status"])
	assert.Equal(t, "1.29", attrs["version"])
}

func TestFetchEKSCluster_EmptyID(t *testing.T) {
	_, err := FetchEKSCluster(context.Background(), &mockEKSClient{}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchEKSCluster_Error(t *testing.T) {
	client := &mockEKSClient{err: errors.New("access denied")}
	_, err := FetchEKSCluster(context.Background(), client, "my-cluster")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestFetchEKSCluster_NilCluster(t *testing.T) {
	client := &mockEKSClient{out: &eks.DescribeClusterOutput{Cluster: nil}}
	_, err := FetchEKSCluster(context.Background(), client, "my-cluster")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
