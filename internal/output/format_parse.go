package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format, returning an error for unknown values.
// An empty string defaults to FormatText.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown output format %q: supported formats are %s", s, supportedFormatsString())
	}
}

// SupportedFormats returns all valid format strings.
func SupportedFormats() []string {
	return []string{string(FormatText), string(FormatJSON)}
}

// IsValidFormat reports whether s is a recognised output format.
func IsValidFormat(s string) bool {
	_, err := ParseFormat(s)
	return err == nil
}

// supportedFormatsString returns a human-readable list of supported formats
// suitable for use in error messages.
func supportedFormatsString() string {
	fmts := SupportedFormats()
	quoted := make([]string, len(fmts))
	for i, f := range fmts {
		quoted[i] = fmt.Sprintf("%q", f)
	}
	return strings.Join(quoted, " and ")
}
