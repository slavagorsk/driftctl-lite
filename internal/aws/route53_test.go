package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoute53Client struct {
	out *route53.GetHostedZoneOutput
	err error
}

func (m *mockRoute53Client) GetHostedZone(_ context.Context, _ *route53.GetHostedZoneInput, _ ...func(*route53.Options)) (*route53.GetHostedZoneOutput, error) {
	return m.out, m.err
}

func TestFetchRoute53HostedZone_Success(t *testing.T) {
	client := &mockRoute53Client{
		out: &route53.GetHostedZoneOutput{
			HostedZone: &types.HostedZone{
				Id:                     aws.String("/hostedzone/Z1234567890"),
				Name:                   aws.String("example.com."),
				ResourceRecordSetCount: aws.Int64(5),
				Config: &types.HostedZoneConfig{
					Comment:     aws.String("managed by terraform"),
					PrivateZone: false,
				},
			},
		},
	}

	attrs, err := FetchRoute53HostedZone(context.Background(), client, "Z1234567890")
	require.NoError(t, err)
	assert.Equal(t, "/hostedzone/Z1234567890", attrs["id"])
	assert.Equal(t, "example.com.", attrs["name"])
	assert.Equal(t, "managed by terraform", attrs["comment"])
	assert.Equal(t, "false", attrs["private_zone"])
	assert.Equal(t, "5", attrs["record_count"])
}

func TestFetchRoute53HostedZone_EmptyID(t *testing.T) {
	_, err := FetchRoute53HostedZone(context.Background(), &mockRoute53Client{}, "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestFetchRoute53HostedZone_Error(t *testing.T) {
	client := &mockRoute53Client{err: errors.New("api error")}
	_, err := FetchRoute53HostedZone(context.Background(), client, "Z999")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "fetching route53 hosted zone")
}

func TestFetchRoute53HostedZone_NilZone(t *testing.T) {
	client := &mockRoute53Client{out: &route53.GetHostedZoneOutput{HostedZone: nil}}
	_, err := FetchRoute53HostedZone(context.Background(), client, "Z999")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
