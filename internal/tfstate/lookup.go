package tfstate

import "fmt"

// ResourceKey uniquely identifies a resource within a state.
type ResourceKey struct {
	Type string
	Name string
}

// String returns a human-readable representation of the key.
func (k ResourceKey) String() string {
	return fmt.Sprintf("%s.%s", k.Type, k.Name)
}

// Index builds a map from ResourceKey to Resource for fast lookups.
func Index(s *State) map[ResourceKey]Resource {
	idx := make(map[ResourceKey]Resource, len(s.Resources))
	for _, r := range s.Resources {
		key := ResourceKey{Type: r.Type, Name: r.Name}
		idx[key] = r
	}
	return idx
}

// FindByType returns all resources of the given type from the state.
func FindByType(s *State, resourceType string) []Resource {
	var result []Resource
	for _, r := range s.Resources {
		if r.Type == resourceType {
			result = append(result, r)
		}
	}
	return result
}

// GetAttribute safely retrieves a string attribute from a resource.
func GetAttribute(r Resource, key string) (string, bool) {
	val, ok := r.Attributes[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}
