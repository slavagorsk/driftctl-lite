package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

type mockWAFv2Client struct {
	out *wafv2.GetWebACLOutput
	err error
}

func (m *mockWAFv2Client) GetWebACL(_ context.Context, _ *wafv2.GetWebACLInput, _ ...func(*wafv2.Options)) (*wafv2.GetWebACLOutput, error) {
	return m.out, m.err
}

func TestFetchWAFv2WebACL_Success(t *testing.T) {
	client := &mockWAFv2Client{
		out: &wafv2.GetWebACLOutput{
			WebACL: &types.WebACL{
				Id:          aws.String("abc123"),
				Name:        aws.String("my-acl"),
				Description: aws.String("test acl"),
				ARN:         aws.String("arn:aws:wafv2:us-east-1:123456789012:regional/webacl/my-acl/abc123"),
			},
		},
	}
	attrs, err := FetchWAFv2WebACL(context.Background(), client, "abc123/my-acl/REGIONAL")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if attrs["id"] != "abc123" {
		t.Errorf("expected id 'abc123', got '%s'", attrs["id"])
	}
	if attrs["name"] != "my-acl" {
		t.Errorf("expected name 'my-acl', got '%s'", attrs["name"])
	}
	if attrs["description"] != "test acl" {
		t.Errorf("expected description 'test acl', got '%s'", attrs["description"])
	}
}

func TestFetchWAFv2WebACL_EmptyID(t *testing.T) {
	client := &mockWAFv2Client{}
	_, err := FetchWAFv2WebACL(context.Background(), client, "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
}

func TestFetchWAFv2WebACL_InvalidFormat(t *testing.T) {
	client := &mockWAFv2Client{}
	_, err := FetchWAFv2WebACL(context.Background(), client, "only-one-part")
	if err == nil {
		t.Fatal("expected error for invalid id format")
	}
}

func TestFetchWAFv2WebACL_UnsupportedScope(t *testing.T) {
	client := &mockWAFv2Client{}
	_, err := FetchWAFv2WebACL(context.Background(), client, "abc123/my-acl/GLOBAL")
	if err == nil {
		t.Fatal("expected error for unsupported scope")
	}
}

func TestFetchWAFv2WebACL_Error(t *testing.T) {
	client := &mockWAFv2Client{err: errors.New("api error")}
	_, err := FetchWAFv2WebACL(context.Background(), client, "abc123/my-acl/REGIONAL")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestFetchWAFv2WebACL_NilWebACL(t *testing.T) {
	client := &mockWAFv2Client{out: &wafv2.GetWebACLOutput{WebACL: nil}}
	_, err := FetchWAFv2WebACL(context.Background(), client, "abc123/my-acl/CLOUDFRONT")
	if err == nil {
		t.Fatal("expected error when WebACL is nil")
	}
}
