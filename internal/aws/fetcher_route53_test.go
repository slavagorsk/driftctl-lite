package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsRoute53Type(t *testing.T) {
	assert.True(t, isRoute53Type("aws_route53_zone"))
	assert.False(t, isRoute53Type("aws_s3_bucket"))
	assert.False(t, isRoute53Type("aws_route53_record"))
	assert.False(t, isRoute53Type(""))
}

func TestFetchRoute53Resource_UnsupportedType(t *testing.T) {
	_, err := fetchRoute53Resource(context.Background(), &mockRoute53Client{}, "aws_route53_record", "Z123")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported route53 resource type")
}

func TestFetchRoute53Resource_EmptyID(t *testing.T) {
	_, err := fetchRoute53Resource(context.Background(), &mockRoute53Client{}, "aws_route53_zone", "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchRoute53Resource_DelegatesToFetchRoute53HostedZone(t *testing.T) {
	client := &mockRoute53Client{
		out: &route53.GetHostedZoneOutput{
			HostedZone: &types.HostedZone{
				Id:                     aws.String("/hostedzone/ZABC"),
				Name:                   aws.String("test.example.com."),
				ResourceRecordSetCount: aws.Int64(2),
				Config:                 &types.HostedZoneConfig{PrivateZone: true},
			},
		},
	}

	attrs, err := fetchRoute53Resource(context.Background(), client, "aws_route53_zone", "ZABC")
	require.NoError(t, err)
	assert.Equal(t, "test.example.com.", attrs["name"])
	assert.Equal(t, "true", attrs["private_zone"])
}
