package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

// ElastiCacheClient defines the interface used for fetching ElastiCache clusters.
type ElastiCacheClient interface {
	DescribeCacheClusters(ctx context.Context, params *elasticache.DescribeCacheClustersInput, optFns ...func(*elasticache.Options)) (*elasticache.DescribeCacheClustersOutput, error)
}

// FetchElastiCacheCluster retrieves an ElastiCache cluster by its cluster ID.
func FetchElastiCacheCluster(ctx context.Context, client ElastiCacheClient, clusterID string) (map[string]string, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("elasticache cluster id must not be empty")
	}

	out, err := client.DescribeCacheClusters(ctx, &elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String(clusterID),
	})
	if err != nil {
		return nil, fmt.Errorf("describing elasticache cluster %q: %w", clusterID, err)
	}

	if len(out.CacheClusters) == 0 {
		return nil, fmt.Errorf("elasticache cluster %q not found", clusterID)
	}

	cluster := out.CacheClusters[0]
	attrs := map[string]string{
		"cluster_id":     aws.ToString(cluster.CacheClusterId),
		"engine":         aws.ToString(cluster.Engine),
		"engine_version": aws.ToString(cluster.EngineVersion),
		"node_type":      aws.ToString(cluster.CacheNodeType),
		"status":         aws.ToString(cluster.CacheClusterStatus),
		"num_nodes":      fmt.Sprintf("%d", cluster.NumCacheNodes),
	}

	if cluster.ConfigurationEndpoint != nil {
		attrs["endpoint"] = fmt.Sprintf("%s:%d",
			aws.ToString(cluster.ConfigurationEndpoint.Address),
			cluster.ConfigurationEndpoint.Port)
	}

	return attrs, nil
}
