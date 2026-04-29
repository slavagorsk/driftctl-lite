package drift

import "fmt"

// Report holds the full result of a drift detection run.
type Report struct {
	Drifts []Drift `json:"drifts"`
}

// Drift describes detected differences for a single cloud resource.
type Drift struct {
	ResourceType string       `json:"resource_type"`
	ResourceID   string       `json:"resource_id"`
	Differences  []Difference `json:"differences"`
}

// Difference captures a single attribute mismatch between state and reality.
type Difference struct {
	Attribute string `json:"attribute"`
	Expected  string `json:"expected"`
	Actual    string `json:"actual"`
}

// HasDrift returns true when the report contains at least one drift.
func (r Report) HasDrift() bool {
	return len(r.Drifts) > 0
}

// Summary returns a human-readable one-line summary of the report.
func (r Report) Summary() string {
	if !r.HasDrift() {
		return "No drift detected."
	}
	total := 0
	for _, d := range r.Drifts {
		total += len(d.Differences)
	}
	return formatSummary(len(r.Drifts), total)
}

// DriftedResources returns only the Drift entries that contain at least one
// Difference, which is useful when a Drift may have been recorded with an
// empty difference list.
func (r Report) DriftedResources() []Drift {
	result := make([]Drift, 0, len(r.Drifts))
	for _, d := range r.Drifts {
		if len(d.Differences) > 0 {
			result = append(result, d)
		}
	}
	return result
}

func formatSummary(resources, diffs int) string {
	return fmt.Sprintf("%d resource(s) drifted with %d attribute difference(s).", resources, diffs)
}
