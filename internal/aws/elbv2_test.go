package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockELBv2Client struct {
	out *elasticloadbalancingv2.DescribeLoadBalancersOutput
	err error
}

func (m *mockELBv2Client) DescribeLoadBalancers(_ context.Context, _ *elasticloadbalancingv2.DescribeLoadBalancersInput, _ ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error) {
	return m.out, m.err
}

func TestFetchELBv2LoadBalancer_Success(t *testing.T) {
	client := &mockELBv2Client{
		out: &elasticloadbalancingv2.DescribeLoadBalancersOutput{
			LoadBalancers: []types.LoadBalancer{
				{
					LoadBalancerArn:  aws.String("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123"),
					LoadBalancerName: aws.String("my-lb"),
					DNSName:          aws.String("my-lb.us-east-1.elb.amazonaws.com"),
					Scheme:           types.LoadBalancerSchemeEnumInternetFacing,
					State:            &types.LoadBalancerState{Code: types.LoadBalancerStateEnumActive},
					Type:             types.LoadBalancerTypeEnumApplication,
					VpcId:            aws.String("vpc-abc123"),
				},
			},
		},
	}
	attrs, err := FetchELBv2LoadBalancer(context.Background(), client, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123")
	require.NoError(t, err)
	assert.Equal(t, "my-lb", attrs["name"])
	assert.Equal(t, "active", attrs["state"])
	assert.Equal(t, "vpc-abc123", attrs["vpc_id"])
}

func TestFetchELBv2LoadBalancer_EmptyID(t *testing.T) {
	_, err := FetchELBv2LoadBalancer(context.Background(), &mockELBv2Client{}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchELBv2LoadBalancer_Error(t *testing.T) {
	client := &mockELBv2Client{err: errors.New("api error")}
	_, err := FetchELBv2LoadBalancer(context.Background(), client, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "describe load balancer")
}

func TestFetchELBv2LoadBalancer_NotFound(t *testing.T) {
	client := &mockELBv2Client{out: &elasticloadbalancingv2.DescribeLoadBalancersOutput{}}
	_, err := FetchELBv2LoadBalancer(context.Background(), client, "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
