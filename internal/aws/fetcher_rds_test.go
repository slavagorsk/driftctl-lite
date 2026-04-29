package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRDSType(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     bool
	}{
		{"aws_db_instance", true},
		{"aws_instance", false},
		{"aws_s3_bucket", false},
		{"aws_iam_role", false},
		{"", false},
	}

	for _, tc := range tests {
		t.Run(tc.resourceType, func(t *testing.T) {
			result := isRDSType(tc.resourceType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFetchRDSResource_UnsupportedType(t *testing.T) {
	cfg := validTestConfig(t)
	awsCfg, err := awsConfigFromRegion(cfg.Region)
	assert.NoError(t, err)

	client, err := rdsClientFromConfig(awsCfg)
	assert.NoError(t, err)

	_, err = fetchRDSResource(client, "aws_rds_cluster", "cluster-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported RDS resource type")
}

func TestFetchRDSResource_EmptyID(t *testing.T) {
	cfg := validTestConfig(t)
	awsCfg, err := awsConfigFromRegion(cfg.Region)
	assert.NoError(t, err)

	client, err := rdsClientFromConfig(awsCfg)
	assert.NoError(t, err)

	_, err = fetchRDSResource(client, "aws_db_instance", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty resource ID")
}

func TestFetchRDSResource_DelegatesToFetchRDSInstance(t *testing.T) {
	cfg := validTestConfig(t)
	awsCfg, err := awsConfigFromRegion(cfg.Region)
	assert.NoError(t, err)

	client, err := rdsClientFromConfig(awsCfg)
	assert.NoError(t, err)

	_, err = fetchRDSResource(client, "aws_db_instance", "db-doesnotexist")
	if err != nil {
		assert.NotContains(t, err.Error(), "unsupported RDS resource type")
	}
}
