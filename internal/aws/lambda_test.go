package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type mockLambdaClient struct {
	out *lambda.GetFunctionOutput
	err error
}

func (m *mockLambdaClient) GetFunction(_ context.Context, _ *lambda.GetFunctionInput, _ ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error) {
	return m.out, m.err
}

func TestFetchLambdaFunction_Success(t *testing.T) {
	client := &mockLambdaClient{
		out: &lambda.GetFunctionOutput{
			Configuration: &types.FunctionConfiguration{
				FunctionName: aws.String("my-function"),
				FunctionArn:  aws.String("arn:aws:lambda:us-east-1:123456789012:function:my-function"),
				Runtime:      types.RuntimeNodejs18x,
				Handler:      aws.String("index.handler"),
				Role:         aws.String("arn:aws:iam::123456789012:role/my-role"),
				Description:  aws.String("test function"),
				State:        types.StateActive,
			},
		},
	}

	attrs, err := FetchLambdaFunction(context.Background(), client, "my-function")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["function_name"] != "my-function" {
		t.Errorf("expected function_name=my-function, got %q", attrs["function_name"])
	}
	if attrs["runtime"] != string(types.RuntimeNodejs18x) {
		t.Errorf("expected runtime=%s, got %q", types.RuntimeNodejs18x, attrs["runtime"])
	}
	if attrs["state"] != string(types.StateActive) {
		t.Errorf("expected state=Active, got %q", attrs["state"])
	}
}

func TestFetchLambdaFunction_EmptyID(t *testing.T) {
	client := &mockLambdaClient{}
	_, err := FetchLambdaFunction(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty function name")
	}
}

func TestFetchLambdaFunction_Error(t *testing.T) {
	client := &mockLambdaClient{err: errors.New("access denied")}
	_, err := FetchLambdaFunction(context.Background(), client, "my-function")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestFetchLambdaFunction_NilConfiguration(t *testing.T) {
	client := &mockLambdaClient{
		out: &lambda.GetFunctionOutput{Configuration: nil},
	}
	_, err := FetchLambdaFunction(context.Background(), client, "my-function")
	if err == nil {
		t.Fatal("expected error for nil configuration")
	}
}
