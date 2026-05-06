package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

func isElastiCacheType(resourceType string) bool {
	return resourceType == "aws_elasticache_cluster"
}

func elastiCacheClientFromConfig(region string) (ElastiCacheClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("loading aws config for elasticache: %w", err)
	}
	return elasticache.NewFromConfig(cfg), nil
}

func fetchElastiCacheResource(ctx context.Context, client ElastiCacheClient, resourceType, id string) (map[string]string, error) {
	if !isElastiCacheType(resourceType) {
		return nil, fmt.Errorf("unsupported elasticache resource type: %s", resourceType)
	}
	if id == "" {
		return nil, fmt.Errorf("resource id must not be empty")
	}
	return FetchElastiCacheCluster(ctx, client, id)
}
