package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

// ELBv2Client defines the interface used for fetching load balancers.
type ELBv2Client interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
}

// FetchELBv2LoadBalancer retrieves an Application/Network Load Balancer by ARN.
func FetchELBv2LoadBalancer(ctx context.Context, client ELBv2Client, arn string) (map[string]string, error) {
	if arn == "" {
		return nil, fmt.Errorf("load balancer ARN must not be empty")
	}

	out, err := client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{
		LoadBalancerArns: []string{arn},
	})
	if err != nil {
		return nil, fmt.Errorf("describe load balancer %s: %w", arn, err)
	}

	if len(out.LoadBalancers) == 0 {
		return nil, fmt.Errorf("load balancer %s not found", arn)
	}

	lb := out.LoadBalancers[0]
	attrs := map[string]string{
		"arn":    aws.ToString(lb.LoadBalancerArn),
		"name":   aws.ToString(lb.LoadBalancerName),
		"dns":    aws.ToString(lb.DNSName),
		"scheme": string(lb.Scheme),
		"state":  string(lb.State.Code),
		"type":   string(lb.Type),
		"vpc_id": aws.ToString(lb.VpcId),
	}
	return attrs, nil
}
