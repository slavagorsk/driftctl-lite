package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

type mockElastiCacheClient struct {
	out *elasticache.DescribeCacheClustersOutput
	err error
}

func (m *mockElastiCacheClient) DescribeCacheClusters(_ context.Context, _ *elasticache.DescribeCacheClustersInput, _ ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error) {
	return m.out, m.err
}

func TestIsElastiCacheType(t *testing.T) {
	if !isElastiCacheType("aws_elasticache_cluster") {
		t.Error("expected aws_elasticache_cluster to be an elasticache type")
	}
	if isElastiCacheType("aws_s3_bucket") {
		t.Error("expected aws_s3_bucket not to be an elasticache type")
	}
}

func TestFetchElastiCacheResource_UnsupportedType(t *testing.T) {
	client := &mockElastiCacheClient{}
	_, err := fetchElastiCacheResource(context.Background(), client, "aws_lambda_function", "my-cluster")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestFetchElastiCacheResource_EmptyID(t *testing.T) {
	client := &mockElastiCacheClient{}
	_, err := fetchElastiCacheResource(context.Background(), client, "aws_elasticache_cluster", "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestFetchElastiCacheResource_DelegatesToFetchElastiCacheCluster(t *testing.T) {
	client := &mockElastiCacheClient{
		out: &elasticache.DescribeCacheClustersOutput{
			CacheClusters: []types.CacheCluster{
				{
					CacheClusterId:     aws.String("my-cluster"),
					Engine:             aws.String("redis"),
					EngineVersion:      aws.String("7.0.7"),
					CacheNodeType:      aws.String("cache.t3.micro"),
					CacheClusterStatus: aws.String("available"),
					NumCacheNodes:      1,
				},
			},
		},
	}
	attrs, err := fetchElastiCacheResource(context.Background(), client, "aws_elasticache_cluster", "my-cluster")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["engine"] != "redis" {
		t.Errorf("expected engine redis, got %s", attrs["engine"])
	}
	if attrs["status"] != "available" {
		t.Errorf("expected status available, got %s", attrs["status"])
	}
}

func TestFetchElastiCacheResource_PropagatesError(t *testing.T) {
	client := &mockElastiCacheClient{err: errors.New("aws error")}
	_, err := fetchElastiCacheResource(context.Background(), client, "aws_elasticache_cluster", "my-cluster")
	if err == nil {
		t.Fatal("expected error to be propagated")
	}
}
