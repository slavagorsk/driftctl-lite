package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type mockCloudWatchClient struct {
	out *cloudwatch.DescribeAlarmsOutput
	err error
}

func (m *mockCloudWatchClient) DescribeAlarms(_ context.Context, _ *cloudwatch.DescribeAlarmsInput, _ ...func(*cloudwatch.Options)) (*cloudwatch.DescribeAlarmsOutput, error) {
	return m.out, m.err
}

func TestIsCloudWatchType(t *testing.T) {
	if !isCloudWatchType("aws_cloudwatch_metric_alarm") {
		t.Error("expected aws_cloudwatch_metric_alarm to be a cloudwatch type")
	}
	if isCloudWatchType("aws_s3_bucket") {
		t.Error("expected aws_s3_bucket not to be a cloudwatch type")
	}
}

func TestFetchCloudWatchResource_UnsupportedType(t *testing.T) {
	_, err := fetchCloudWatchResource(context.Background(), "us-east-1", "aws_lambda_function", "my-fn")
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestFetchCloudWatchResource_EmptyID(t *testing.T) {
	_, err := fetchCloudWatchResource(context.Background(), "us-east-1", "aws_cloudwatch_metric_alarm", "")
	if err == nil {
		t.Fatal("expected error for empty resource ID")
	}
}

func TestFetchCloudWatchResource_DelegatesToFetchCloudWatchAlarm(t *testing.T) {
	client := &mockCloudWatchClient{
		out: &cloudwatch.DescribeAlarmsOutput{
			MetricAlarms: []types.MetricAlarm{
				{
					AlarmName: aws.String("my-alarm"),
					AlarmArn:  aws.String("arn:aws:cloudwatch:us-east-1:123456789012:alarm:my-alarm"),
					StateValue: types.StateValueOk,
				},
			},
		},
	}
	attrs, err := FetchCloudWatchAlarm(context.Background(), client, "my-alarm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["alarm_name"] != "my-alarm" {
		t.Errorf("expected alarm_name=my-alarm, got %s", attrs["alarm_name"])
	}
	if attrs["state"] != "OK" {
		t.Errorf("expected state=OK, got %s", attrs["state"])
	}
}

func TestFetchCloudWatchAlarm_Error(t *testing.T) {
	client := &mockCloudWatchClient{err: errors.New("api error")}
	_, err := FetchCloudWatchAlarm(context.Background(), client, "my-alarm")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFetchCloudWatchAlarm_NotFound(t *testing.T) {
	client := &mockCloudWatchClient{
		out: &cloudwatch.DescribeAlarmsOutput{MetricAlarms: []types.MetricAlarm{}},
	}
	_, err := FetchCloudWatchAlarm(context.Background(), client, "missing-alarm")
	if err == nil {
		t.Fatal("expected not-found error")
	}
}
