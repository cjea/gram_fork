package deplconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/speakeasy-api/gram/server/cmd/cli/gram/env"
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
	// SchemaVersion defines the version of the configuration schema.
	SchemaVersion string `json:"schema_version"`

	// Sources defines the list of prospective assets to include in the
	// deployment.
	Sources []Source `json:"sources"`
}

// GetProducerToken returns an API key with a `producer` scope.
func (dc DeploymentConfig) GetProducerToken() string {
	return env.MustApiKey()
}

var ValidSchemaVersions = []string{"1.0.0"}

// SchemaValid returns true if the incoming schema version is valid.
func (dc DeploymentConfig) SchemaValid() bool {
	return slices.Contains(ValidSchemaVersions, dc.SchemaVersion)
}

func (dc DeploymentConfig) Validate() error {
	if !dc.SchemaValid() {
		msg := "unsupported schema version: '%s'. Expected one of %+v"

		return fmt.Errorf(msg, dc.SchemaVersion, ValidSchemaVersions)
	}

	if len(dc.Sources) < 1 {
		return fmt.Errorf("must specify at least one source")
	}

	return nil
}

// ReadDeploymentConfig reads a deployment config.
func ReadDeploymentConfig(filePath string) (*DeploymentConfig, error) {
	var cfg DeploymentConfig

	data, err := readFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
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

	data, err := os.ReadFile(filePath) // #nosec G304
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}
