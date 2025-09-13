package deplconfig

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// SourceReader reads source content from local files or remote URLs.
type SourceReader struct {
	source Source
}

// NewSourceReader creates a new SourceReader for the given source.
func NewSourceReader(source Source) *SourceReader {
	return &SourceReader{
		source: source,
	}
}

// GetType returns the source type (e.g., "openapiv3").
func (sr *SourceReader) GetType() string {
	return string(sr.source.Type)
}

// GetContentType returns the MIME type of the content based on file extension.
func (sr *SourceReader) GetContentType() string {
	if isRemoteURL(sr.source.Location) {
		// For remote URLs, we'll need to determine content type differently.
		// For now, default to common OpenAPI types based on extension.
		return getContentTypeFromPath(sr.source.Location)
	}

	return getContentTypeFromPath(sr.source.Location)
}

// Read returns a reader for the asset content and its size.
func (sr *SourceReader) Read() (io.ReadCloser, int64, error) {
	if isRemoteURL(sr.source.Location) {
		return sr.readRemote()
	}
	return sr.readLocal()
}

// readLocal reads from a local file path.
func (sr *SourceReader) readLocal() (io.ReadCloser, int64, error) {
	data, err := readFile(sr.source.Location)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read local file: %w", err)
	}

	fi, err := os.Stat(sr.source.Location)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get file info: %w", err)
	}

	reader := strings.NewReader(string(data))
	return io.NopCloser(reader), fi.Size(), nil
}

// readRemote reads from a remote URL (placeholder for future implementation).
func (sr *SourceReader) readRemote() (io.ReadCloser, int64, error) {
	// TODO: Implement remote file reading (HTTP/HTTPS URLs)
	// This could use http.Get() and handle redirects, authentication, etc.
	return nil, 0, fmt.Errorf("remote URL reading not yet implemented: %s", sr.source.Location)
}

// isRemoteURL checks if the location is a remote URL.
func isRemoteURL(location string) bool {
	return strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://")
}

// getContentTypeFromPath determines content type from file path/extension.
func getContentTypeFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/yaml"
	default:
		if mimeType := mime.TypeByExtension(ext); mimeType != "" {
			return mimeType
		}
		defaultForOpenAPISpecs := "application/yaml"
		return defaultForOpenAPISpecs
	}
}
