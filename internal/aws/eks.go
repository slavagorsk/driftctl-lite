package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// EKSClient defines the interface for EKS operations used in fetching.
type EKSClient interface {
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
}

// FetchEKSCluster retrieves an EKS cluster by name and returns its attributes.
func FetchEKSCluster(ctx context.Context, client EKSClient, clusterName string) (map[string]string, error) {
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name must not be empty")
	}

	out, err := client.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("describing EKS cluster %q: %w", clusterName, err)
	}

	if out.Cluster == nil {
		return nil, fmt.Errorf("EKS cluster %q not found", clusterName)
	}

	c := out.Cluster
	attrs := map[string]string{
		"name":        aws.ToString(c.Name),
		"arn":         aws.ToString(c.Arn),
		"status":      string(c.Status),
		"version":     aws.ToString(c.Version),
		"endpoint":    aws.ToString(c.Endpoint),
		"role_arn":    aws.ToString(c.RoleArn),
	}

	return attrs, nil
}
