package config_test

import (
	"os"
	"testing"

	"github.com/example/driftctl-lite/internal/config"
)

func writeTempState(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "tfstate-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.Close()
	return f.Name()
}

func TestNew_Defaults(t *testing.T) {
	cfg := config.New("/tmp/state.tfstate", "", "", nil)
	if cfg.AWSRegion != config.DefaultRegion {
		t.Errorf("expected default region %q, got %q", config.DefaultRegion, cfg.AWSRegion)
	}
	if cfg.OutputFormat != config.DefaultOutputFormat {
		t.Errorf("expected default format %q, got %q", config.DefaultOutputFormat, cfg.OutputFormat)
	}
}

func TestNew_ResourceTypesTrimmed(t *testing.T) {
	cfg := config.New("/tmp/s.tfstate", "eu-west-1", "json", []string{" aws_s3_bucket ", "", "aws_instance"})
	if len(cfg.ResourceTypes) != 2 {
		t.Fatalf("expected 2 resource types, got %d", len(cfg.ResourceTypes))
	}
	if cfg.ResourceTypes[0] != "aws_s3_bucket" {
		t.Errorf("unexpected resource type: %q", cfg.ResourceTypes[0])
	}
}

func TestValidate_MissingStatePath(t *testing.T) {
	cfg := config.New("", "us-east-1", "text", nil)
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing state path")
	}
}

func TestValidate_StateFileNotExist(t *testing.T) {
	cfg := config.New("/nonexistent/path/state.tfstate", "us-east-1", "text", nil)
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for non-existent state file")
	}
}

func TestValidate_Valid(t *testing.T) {
	path := writeTempState(t)
	cfg := config.New(path, "us-west-2", "json", []string{"aws_s3_bucket"})
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_MissingRegion(t *testing.T) {
	path := writeTempState(t)
	cfg := &config.Config{
		StatePath:    path,
		AWSRegion:    "",
		OutputFormat: "text",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing region")
	}
}
