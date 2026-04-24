package drift

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintReport writes a formatted drift report to w.
func PrintReport(w io.Writer, deltas []Delta) {
	if len(deltas) == 0 {
		fmt.Fprintln(w, "✓ No drift detected.")
		return
	}

	fmt.Fprintf(w, "⚠ Drift detected: %d difference(s)\n\n", len(deltas))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "RESOURCE\tATTRIBUTE\tSTATE VALUE\tLIVE VALUE")
	fmt.Fprintln(tw, "--------\t---------\t-----------\t----------")
	for _, d := range deltas {
		fmt.Fprintf(tw, "%s.%s\t%s\t%s\t%s\n",
			d.ResourceType, d.ResourceName,
			d.Attribute,
			d.StateValue,
			d.LiveValue,
		)
	}
	_ = tw.Flush()
}
