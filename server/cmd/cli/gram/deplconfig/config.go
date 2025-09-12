package deplconfig

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/speakeasy-api/gram/server/cmd/cli/env"
)

type SourceType string

const (
	SourceTypeOpenAPIV3 SourceType = "openapiv3"
)

type Source struct {
	Type SourceType `json:"type"`

	// Location is the filepath or remote URL of the asset source.
	Location string `json:"location"`
}

type DeploymentConfig struct {
	// Project defines which project in which to create the deployment.
	Project string `json:"project"`

	// Sources defines the list of prospective assets to include in the
	// deployment.
	Sources []Source `json:"sources"`

	// GetProducerToken returns an API key with a `producer` scope.
	GetProducerToken func() string
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

	cfg.GetProducerToken = env.ReadApiKey

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
