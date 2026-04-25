package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds the runtime configuration for driftctl-lite.
type Config struct {
	StatePath  string
	AWSRegion  string
	OutputFormat string
	ResourceTypes []string
}

// Defaults
const (
	DefaultRegion       = "us-east-1"
	DefaultOutputFormat = "text"
)

// Validate checks that all required fields are set and valid.
func (c *Config) Validate() error {
	if c.StatePath == "" {
		return errors.New("state path must not be empty")
	}
	if _, err := os.Stat(c.StatePath); os.IsNotExist(err) {
		return errors.New("state file does not exist: " + c.StatePath)
	}
	if c.AWSRegion == "" {
		return errors.New("AWS region must not be empty")
	}
	if c.OutputFormat == "" {
		return errors.New("output format must not be empty")
	}
	return nil
}

// New returns a Config with defaults applied.
func New(statePath, region, outputFormat string, resourceTypes []string) *Config {
	if region == "" {
		region = DefaultRegion
	}
	if outputFormat == "" {
		outputFormat = DefaultOutputFormat
	}
	// Normalise resource type list
	normalised := make([]string, 0, len(resourceTypes))
	for _, rt := range resourceTypes {
		rt = strings.TrimSpace(rt)
		if rt != "" {
			normalised = append(normalised, rt)
		}
	}
	return &Config{
		StatePath:     statePath,
		AWSRegion:     region,
		OutputFormat:  outputFormat,
		ResourceTypes: normalised,
	}
}
