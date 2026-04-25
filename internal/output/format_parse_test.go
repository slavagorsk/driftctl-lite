package output

import (
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{" text ", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if got != tc.expected {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	invalid := []string{"yaml", "xml", "csv", "table"}
	for _, s := range invalid {
		t.Run(s, func(t *testing.T) {
			_, err := ParseFormat(s)
			if err == nil {
				t.Errorf("expected error for format %q, got nil", s)
			}
		})
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := SupportedFormats()
	if len(formats) < 2 {
		t.Errorf("expected at least 2 supported formats, got %d", len(formats))
	}
	found := map[string]bool{}
	for _, f := range formats {
		found[f] = true
	}
	if !found["text"] {
		t.Error("expected \"text\" in supported formats")
	}
	if !found["json"] {
		t.Error("expected \"json\" in supported formats")
	}
}
