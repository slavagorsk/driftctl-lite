package drift

import (
	"errors"
	"testing"

	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type stubFetcher struct {
	attrs map[string]string
	err   error
}

func (s *stubFetcher) Fetch(_, _ string) (map[string]string, error) {
	return s.attrs, s.err
}

func buildState(attrs map[string]interface{}) *tfstate.State {
	return &tfstate.State{
		Resources: []tfstate.Resource{
			{
				Type: "aws_s3_bucket",
				Name: "example",
				Instances: []tfstate.Instance{
					{Attributes: attrs},
				},
			},
		},
	}
}

func TestDetect_NoDrift(t *testing.T) {
	state := buildState(map[string]interface{}{"bucket": "my-bucket", "region": "us-east-1"})
	fetcher := &stubFetcher{attrs: map[string]string{"bucket": "my-bucket", "region": "us-east-1"}}
	deltas, err := Detect(state, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deltas) != 0 {
		t.Errorf("expected 0 deltas, got %d: %v", len(deltas), deltas)
	}
}

func TestDetect_WithDrift(t *testing.T) {
	state := buildState(map[string]interface{}{"bucket": "my-bucket", "region": "us-east-1"})
	fetcher := &stubFetcher{attrs: map[string]string{"bucket": "my-bucket", "region": "eu-west-1"}}
	deltas, err := Detect(state, fetcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deltas) != 1 {
		t.Fatalf("expected 1 delta, got %d", len(deltas))
	}
	if deltas[0].Attribute != "region" {
		t.Errorf("expected attribute 'region', got %q", deltas[0].Attribute)
	}
}

func TestDetect_FetchError(t *testing.T) {
	state := buildState(map[string]interface{}{})
	fetcher := &stubFetcher{err: errors.New("api error")}
	_, err := Detect(state, fetcher)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
