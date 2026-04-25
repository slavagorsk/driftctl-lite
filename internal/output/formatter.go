package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/driftctl-lite/internal/drift"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift reports in a specific format.
type Formatter struct {
	format Format
	out    io.Writer
}

// NewFormatter creates a new Formatter writing to the given writer.
func NewFormatter(format Format, out io.Writer) *Formatter {
	if out == nil {
		out = os.Stdout
	}
	return &Formatter{format: format, out: out}
}

// Write outputs the drift report in the configured format.
func (f *Formatter) Write(report drift.Report) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(report)
	default:
		return f.writeText(report)
	}
}

func (f *Formatter) writeJSON(report drift.Report) error {
	enc := json.NewEncoder(f.out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (f *Formatter) writeText(report drift.Report) error {
	if len(report.Drifts) == 0 {
		fmt.Fprintln(f.out, "✅  No drift detected.")
		return nil
	}

	fmt.Fprintf(f.out, "⚠️  Drift detected: %d resource(s) differ\n\n", len(report.Drifts))

	w := tabwriter.NewWriter(f.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RESOURCE TYPE\tRESOURCE ID\tATTRIBUTE\tEXPECTED\tACTUAL")
	fmt.Fprintln(w, "-------------\t-----------\t---------\t--------\t------")
	for _, d := range report.Drifts {
		for _, diff := range d.Differences {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				d.ResourceType, d.ResourceID, diff.Attribute, diff.Expected, diff.Actual)
		}
	}
	return w.Flush()
}
