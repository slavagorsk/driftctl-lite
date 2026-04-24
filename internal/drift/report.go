package drift

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// DriftStatus represents the drift status of a resource.
type DriftStatus string

const (
	StatusOK      DriftStatus = "OK"
	StatusMissing DriftStatus = "MISSING"
	StatusChanged DriftStatus = "CHANGED"
)

// ResourceDrift holds drift information for a single resource.
type ResourceDrift struct {
	ResourceType string
	ResourceID   string
	Status       DriftStatus
	Details      string
}

// Report aggregates drift results.
type Report struct {
	Items []ResourceDrift
}

// HasDrift returns true if any resource has drifted.
func (r *Report) HasDrift() bool {
	for _, item := range r.Items {
		if item.Status != StatusOK {
			return true
		}
	}
	return false
}

// PrintReport writes a formatted drift report to stdout.
func PrintReport(report *Report) {
	printReportTo(report, os.Stdout)
}

func printReportTo(report *Report, w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tID\tSTATUS\tDETAILS")
	fmt.Fprintln(tw, "----\t--\t------\t-------")
	for _, item := range report.Items {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			item.ResourceType,
			item.ResourceID,
			item.Status,
			item.Details,
		)
	}
	tw.Flush()
	if report.HasDrift() {
		fmt.Fprintln(w, "\n⚠  Drift detected.")
	} else {
		fmt.Fprintln(w, "\n✓  No drift detected.")
	}
}
