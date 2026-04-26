package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIAMClient struct {
	userOut *iam.GetUserOutput
	userErr error
	roleOut *iam.GetRoleOutput
	roleErr error
}

func (m *mockIAMClient) GetUser(_ context.Context, _ *iam.GetUserInput, _ ...func(*iam.Options)) (*iam.GetUserOutput, error) {
	return m.userOut, m.userErr
}

func (m *mockIAMClient) GetRole(_ context.Context, _ *iam.GetRoleInput, _ ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	return m.roleOut, m.roleErr
}

func TestFetchIAMUser_Success(t *testing.T) {
	client := &mockIAMClient{
		userOut: &iam.GetUserOutput{
			User: &types.User{
				UserName: aws.String("alice"),
				Arn:      aws.String("arn:aws:iam::123456789012:user/alice"),
				UserId:   aws.String("AIDAXXXXXXXXXXXXXXXXX"),
				Path:     aws.String("/"),
			},
		},
	}
	attrs, err := FetchIAMUser(context.Background(), client, "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["user_name"] != "alice" {
		t.Errorf("expected user_name=alice, got %q", attrs["user_name"])
	}
	if attrs["arn"] != "arn:aws:iam::123456789012:user/alice" {
		t.Errorf("unexpected arn: %q", attrs["arn"])
	}
}

func TestFetchIAMUser_EmptyID(t *testing.T) {
	_, err := FetchIAMUser(context.Background(), &mockIAMClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty user name")
	}
}

func TestFetchIAMUser_Error(t *testing.T) {
	client := &mockIAMClient{userErr: errors.New("no such user")}
	_, err := FetchIAMUser(context.Background(), client, "ghost")
	if err == nil {
		t.Fatal("expected error from client")
	}
}

func TestFetchIAMRole_Success(t *testing.T) {
	client := &mockIAMClient{
		roleOut: &iam.GetRoleOutput{
			Role: &types.Role{
				RoleName:                 aws.String("my-role"),
				Arn:                      aws.String("arn:aws:iam::123456789012:role/my-role"),
				RoleId:                   aws.String("AROAXXXXXXXXXXXXXXXXX"),
				Path:                     aws.String("/"),
				AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17"}`),
			},
		},
	}
	attrs, err := FetchIAMRole(context.Background(), client, "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["role_name"] != "my-role" {
		t.Errorf("expected role_name=my-role, got %q", attrs["role_name"])
	}
}

func TestFetchIAMRole_EmptyID(t *testing.T) {
	_, err := FetchIAMRole(context.Background(), &mockIAMClient{}, "")
	if err == nil {
		t.Fatal("expected error for empty role name")
	}
}

func TestFetchIAMRole_Error(t *testing.T) {
	client := &mockIAMClient{roleErr: errors.New("no such role")}
	_, err := FetchIAMRole(context.Background(), client, "missing-role")
	if err == nil {
		t.Fatal("expected error from client")
	}
}
