package aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDynamoDBClient struct {
	output *dynamodb.DescribeTableOutput
	err    error
}

func (m *mockDynamoDBClient) DescribeTable(_ context.Context, _ *dynamodb.DescribeTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	return m.output, m.err
}

func TestFetchDynamoDBTable_Success(t *testing.T) {
	name := "my-table"
	status := types.TableStatusActive
	client := &mockDynamoDBClient{
		output: &dynamodb.DescribeTableOutput{
			Table: &types.TableDescription{
				TableName:   &name,
				TableStatus: status,
			},
		},
	}
	attrs, err := FetchDynamoDBTable(client, name)
	require.NoError(t, err)
	assert.Equal(t, "my-table", attrs["table_name"])
	assert.Equal(t, "ACTIVE", attrs["table_status"])
}

func TestFetchDynamoDBTable_EmptyID(t *testing.T) {
	client := &mockDynamoDBClient{}
	_, err := FetchDynamoDBTable(client, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "table name must not be empty")
}

func TestFetchDynamoDBTable_Error(t *testing.T) {
	client := &mockDynamoDBClient{err: fmt.Errorf("connection refused")}
	_, err := FetchDynamoDBTable(client, "my-table")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to describe DynamoDB table")
}

func TestFetchDynamoDBTable_NilTable(t *testing.T) {
	client := &mockDynamoDBClient{
		output: &dynamodb.DescribeTableOutput{Table: nil},
	}
	_, err := FetchDynamoDBTable(client, "missing-table")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
