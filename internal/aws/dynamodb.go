package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBDescribeTableAPI defines the interface for describing a DynamoDB table.
type DynamoDBDescribeTableAPI interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

// FetchDynamoDBTable retrieves attributes of a DynamoDB table by name.
func FetchDynamoDBTable(client DynamoDBDescribeTableAPI, tableName string) (map[string]string, error) {
	if tableName == "" {
		return nil, fmt.Errorf("table name must not be empty")
	}

	out, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe DynamoDB table %q: %w", tableName, err)
	}
	if out.Table == nil {
		return nil, fmt.Errorf("DynamoDB table %q not found", tableName)
	}

	attrs := map[string]string{
		"table_name": *out.Table.TableName,
		"table_status": string(out.Table.TableStatus),
	}
	if out.Table.BillingModeSummary != nil {
		attrs["billing_mode"] = string(out.Table.BillingModeSummary.BillingMode)
	}
	return attrs, nil
}
