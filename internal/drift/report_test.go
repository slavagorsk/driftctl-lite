package drift

import (
	"bytes"
	"strings"
	"testing"
)

func TestReport_HasDrift_NoDrift(t *testing.T) {
	report := &Report{
		Items: []ResourceDrift{
			{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Status: StatusOK, Details: ""},
		},
	}
	if report.HasDrift() {
		t.Error("expected no drift, but HasDrift returned true")
	}
}

func TestReport_HasDrift_WithDrift(t *testing.T) {
	report := &Report{
		Items: []ResourceDrift{
			{ResourceType: "aws_s3_bucket", ResourceID: "my-bucket", Status: StatusMissing, Details: "not found in AWS"},
		},
	}
	if !report.HasDrift() {
		t.Error("expected drift, but HasDrift returned false")
	}
}

func TestPrintReport_NoDrift(t *testing.T) {
	report := &Report{
		Items: []ResourceDrift{
			{ResourceType: "aws_s3_bucket", ResourceID: "bucket-1", Status: StatusOK, Details: ""},
		},
	}
	var buf bytes.Buffer
	printReportTo(report, &buf)
	out := buf.String()
	if !strings.Contains(out, "bucket-1") {
		t.Errorf("expected resource ID in output, got: %s", out)
	}
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message in output, got: %s", out)
	}
}

func TestPrintReport_WithDrift(t *testing.T) {
	report := &Report{
		Items: []ResourceDrift{
			{ResourceType: "aws_s3_bucket", ResourceID: "bucket-2", Status: StatusMissing, Details: "bucket not found"},
			{ResourceType: "aws_s3_bucket", ResourceID: "bucket-3", Status: StatusChanged, Details: "tag mismatch"},
		},
	}
	var buf bytes.Buffer
	printReportTo(report, &buf)
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING status in output, got: %s", out)
	}
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected CHANGED status in output, got: %s", out)
	}
	if !strings.Contains(out, "Drift detected") {
		t.Errorf("expected drift-detected message in output, got: %s", out)
	}
}

func TestPrintReport_Empty(t *testing.T) {
	report := &Report{Items: []ResourceDrift{}}
	var buf bytes.Buffer
	printReportTo(report, &buf)
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message for empty report, got: %s", out)
	}
}
