package tfstate_test

import (
	"testing"

	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func buildTestState() *tfstate.State {
	return &tfstate.State{
		Version: 4,
		Resources: []tfstate.Resource{
			{
				Type: "aws_s3_bucket",
				Name: "bucket_a",
				Attributes: map[string]interface{}{"id": "bucket-a", "region": "us-east-1"},
			},
			{
				Type: "aws_s3_bucket",
				Name: "bucket_b",
				Attributes: map[string]interface{}{"id": "bucket-b"},
			},
			{
				Type: "aws_iam_role",
				Name: "role_a",
				Attributes: map[string]interface{}{"id": "role-a"},
			},
		},
	}
}

func TestIndex(t *testing.T) {
	idx := tfstate.Index(buildTestState())
	if len(idx) != 3 {
		t.Errorf("expected 3 entries in index, got %d", len(idx))
	}
	key := tfstate.ResourceKey{Type: "aws_iam_role", Name: "role_a"}
	if _, ok := idx[key]; !ok {
		t.Errorf("expected key %s to exist in index", key)
	}
}

func TestFindByType(t *testing.T) {
	buckets := tfstate.FindByType(buildTestState(), "aws_s3_bucket")
	if len(buckets) != 2 {
		t.Errorf("expected 2 s3 buckets, got %d", len(buckets))
	}
	none := tfstate.FindByType(buildTestState(), "aws_lambda_function")
	if len(none) != 0 {
		t.Errorf("expected 0 lambdas, got %d", len(none))
	}
}

func TestGetAttribute(t *testing.T) {
	res := tfstate.Resource{
		Type:       "aws_s3_bucket",
		Name:       "test",
		Attributes: map[string]interface{}{"id": "my-id"},
	}
	val, ok := tfstate.GetAttribute(res, "id")
	if !ok || val != "my-id" {
		t.Errorf("expected id=my-id, got %q ok=%v", val, ok)
	}
	_, ok = tfstate.GetAttribute(res, "missing")
	if ok {
		t.Error("expected missing key to return ok=false")
	}
}
