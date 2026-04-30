package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

// CloudWatchClient defines the interface used for fetching CloudWatch alarms.
type CloudWatchClient interface {
	DescribeAlarms(ctx context.Context, params *cloudwatch.DescribeAlarmsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.DescribeAlarmsOutput, error)
}

// FetchCloudWatchAlarm retrieves attributes of a CloudWatch alarm by name.
func FetchCloudWatchAlarm(ctx context.Context, client CloudWatchClient, alarmName string) (map[string]string, error) {
	if alarmName == "" {
		return nil, fmt.Errorf("alarm name must not be empty")
	}

	out, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{alarmName},
	})
	if err != nil {
		return nil, fmt.Errorf("describe alarm %q: %w", alarmName, err)
	}

	if len(out.MetricAlarms) == 0 {
		return nil, fmt.Errorf("alarm %q not found", alarmName)
	}

	alarm := out.MetricAlarms[0]
	attrs := map[string]string{
		"alarm_name": aws.ToString(alarm.AlarmName),
		"alarm_arn":  aws.ToString(alarm.AlarmArn),
		"state":      string(alarm.StateValue),
	}
	if alarm.AlarmDescription != nil {
		attrs["alarm_description"] = aws.ToString(alarm.AlarmDescription)
	}
	return attrs, nil
}
