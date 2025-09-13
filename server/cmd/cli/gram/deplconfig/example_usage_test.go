package deplconfig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/speakeasy-api/gram/server/cmd/cli/gram/api"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/deplconfig"
)

func TestSourceReader_ImplementsAssetSource(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "openapi.yaml")
	testContent := `openapi: 3.0.0
info:
  title: Example API
  version: 1.0.0
paths:
  /users:
    get:
      summary: Get users
      responses:
        '200':
          description: OK`

	require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0600))

	source := deplconfig.Source{
		Type:     deplconfig.SourceTypeOpenAPIV3,
		Location: testFile,
	}

	reader := deplconfig.NewSourceReader(source)

	// Verify it satisfies SourceReader interface by using it as one.
	var assetSource api.SourceReader = reader

	// Test all AssetSource methods work
	require.Equal(t, "openapiv3", assetSource.GetType())
	require.Equal(t, "application/yaml", assetSource.GetContentType())

	rc, size, err := assetSource.Read()
	require.NoError(t, err)
	require.Positive(t, size)
	defer func() {
		require.NoError(t, rc.Close())
	}()

	t.Logf("Successfully created AssetSource for file: %s (size: %d bytes)", testFile, size)
}
