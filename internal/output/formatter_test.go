package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/driftctl-lite/internal/drift"
)

func buildReport(hasDrift bool) drift.Report {
	if !hasDrift {
		return drift.Report{Drifts: []drift.Drift{}}
	}
	return drift.Report{
		Drifts: []drift.Drift{
			{
				ResourceType: "aws_s3_bucket",
				ResourceID:   "my-bucket",
				Differences: []drift.Difference{
					{Attribute: "region", Expected: "us-east-1", Actual: "eu-west-1"},
				},
			},
		},
	}
}

func TestWriteText_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatText, &buf)
	if err := f.Write(buildReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestWriteText_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatText, &buf)
	if err := f.Write(buildReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket") {
		t.Errorf("expected resource type in output, got: %s", out)
	}
	if !strings.Contains(out, "eu-west-1") {
		t.Errorf("expected actual value in output, got: %s", out)
	}
}

func TestWriteJSON_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatJSON, &buf)
	if err := f.Write(buildReport(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result drift.Report
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result.Drifts) != 0 {
		t.Errorf("expected 0 drifts, got %d", len(result.Drifts))
	}
}

func TestWriteJSON_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(FormatJSON, &buf)
	if err := f.Write(buildReport(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result drift.Report
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result.Drifts) != 1 {
		t.Errorf("expected 1 drift, got %d", len(result.Drifts))
	}
}

func TestNewFormatter_DefaultsToStdout(t *testing.T) {
	f := NewFormatter(FormatText, nil)
	if f.out == nil {
		t.Error("expected non-nil writer")
	}
}
