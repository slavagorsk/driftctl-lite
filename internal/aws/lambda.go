package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// LambdaClient defines the interface for Lambda API calls used in drift detection.
type LambdaClient interface {
	GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
}

// FetchLambdaFunction retrieves a Lambda function's attributes by name/ARN.
// Returns a map of key attributes or an error if the function is not found.
func FetchLambdaFunction(ctx context.Context, client LambdaClient, functionName string) (map[string]string, error) {
	if functionName == "" {
		return nil, fmt.Errorf("lambda function name must not be empty")
	}

	out, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return nil, fmt.Errorf("fetching lambda function %q: %w", functionName, err)
	}

	if out.Configuration == nil {
		return nil, fmt.Errorf("lambda function %q returned empty configuration", functionName)
	}

	attrs := map[string]string{
		"function_name": aws.ToString(out.Configuration.FunctionName),
		"arn":           aws.ToString(out.Configuration.FunctionArn),
		"runtime":       string(out.Configuration.Runtime),
		"handler":       aws.ToString(out.Configuration.Handler),
		"role":          aws.ToString(out.Configuration.Role),
		"description":   aws.ToString(out.Configuration.Description),
		"state":         string(out.Configuration.State),
	}

	return attrs, nil
}
