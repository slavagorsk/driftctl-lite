package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// Delta represents a single attribute difference between state and live.
type Delta struct {
	ResourceType string
	ResourceName string
	Attribute    string
	StateValue   string
	LiveValue    string
}

// String returns a human-readable representation of the delta.
func (d Delta) String() string {
	return fmt.Sprintf("%s.%s → %s: state=%q live=%q",
		d.ResourceType, d.ResourceName, d.Attribute, d.StateValue, d.LiveValue)
}

// LiveFetcher is implemented by any type that can return live attributes
// for a given resource type and name.
type LiveFetcher interface {
	Fetch(resourceType, resourceName string) (map[string]string, error)
}

// Detect compares tfstate resources against live attributes returned by
// fetcher and returns a slice of Deltas for every mismatch found.
func Detect(state *tfstate.State, fetcher LiveFetcher) ([]Delta, error) {
	var deltas []Delta
	for _, res := range state.Resources {
		for _, inst := range res.Instances {
			live, err := fetcher.Fetch(res.Type, res.Name)
			if err != nil {
				return nil, fmt.Errorf("fetching %s.%s: %w", res.Type, res.Name, err)
			}
			for k, lv := range live {
				sv, _ := inst.Attributes[k].(string)
				if sv != lv {
					deltas = append(deltas, Delta{
						ResourceType: res.Type,
						ResourceName: res.Name,
						Attribute:    k,
						StateValue:   sv,
						LiveValue:    lv,
					})
				}
			}
		}
	}
	return deltas, nil
}
