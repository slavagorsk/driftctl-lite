//go:build integration
// +build integration

package aws

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFetchEKSCluster_Integration verifies real AWS connectivity.
// Set env vars: AWS_REGION, TEST_EKS_CLUSTER_NAME.
func TestFetchEKSCluster_Integration(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		t.Skip("AWS_REGION not set; skipping integration test")
	}

	clusterName := os.Getenv("TEST_EKS_CLUSTER_NAME")
	if clusterName == "" {
		t.Skip("TEST_EKS_CLUSTER_NAME not set; skipping integration test")
	}

	attrs, err := fetchEKSResource(context.Background(), region, "aws_eks_cluster", clusterName)
	require.NoError(t, err)
	assert.Equal(t, clusterName, attrs["name"])
	assert.NotEmpty(t, attrs["arn"])
	assert.NotEmpty(t, attrs["status"])
	assert.NotEmpty(t, attrs["version"])
}
