package depconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type SourceType string

const (
	SourceTypeOpenAPIV3 SourceType = "openapiv3"
)

type Source struct {
	Type SourceType `json:"type"`

	// Location is the filepath or remote URL of the asset source.
	Location string `json:"location"`

	// ProducerToken returns an API key with a `producer` scope.
	ProducerToken func() string
}

type DeploymentConfig struct {
	Project string   `json:"project"`
	Sources []Source `json:"sources"`
}

func ReadDeploymentConfig(filePath string) (*DeploymentConfig, error) {
	var cfg DeploymentConfig

	data, err := readFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &cfg, nil
}

// readFile validates that a file exists at `filePath` and that its mode is
// regular.
func readFile(filePath string) ([]byte, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}
	if !fi.Mode().IsRegular() {
		return nil, fmt.Errorf("path must be a regular file")
	}

	return os.ReadFile(filePath)
}
