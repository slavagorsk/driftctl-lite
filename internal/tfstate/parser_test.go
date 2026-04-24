package tfstate_test

import (
	"testing"

	"github.com/your-org/driftctl-lite/internal/tfstate"
)

const sampleState = `{
  "version": 4,
  "resources": [
    {
      "type": "aws_s3_bucket",
      "name": "my_bucket",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "id": "my-bucket-name",
            "region": "us-east-1"
          }
        }
      ]
    }
  ]
}`

func TestParse_ValidState(t *testing.T) {
	state, err := tfstate.Parse([]byte(sampleState))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	res := state.Resources[0]
	if res.Type != "aws_s3_bucket" {
		t.Errorf("expected type aws_s3_bucket, got %s", res.Type)
	}
	if res.Name != "my_bucket" {
		t.Errorf("expected name my_bucket, got %s", res.Name)
	}
	if res.Attributes["id"] != "my-bucket-name" {
		t.Errorf("unexpected id attribute: %v", res.Attributes["id"])
	}
}

func TestParse_EmptyResources(t *testing.T) {
	state, err := tfstate.Parse([]byte(`{"version":4,"resources":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(state.Resources))
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := tfstate.Parse([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
