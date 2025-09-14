package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/speakeasy-api/gram/server/cmd/cli/gram/env"
	"github.com/speakeasy-api/gram/server/gen/deployments"
	depl_client "github.com/speakeasy-api/gram/server/gen/http/deployments/client"
	goahttp "goa.design/goa/v3/http"
)

type DeploymentsClient struct {
	client *deployments.Client
}

func NewDeploymentsClient() *DeploymentsClient {
	return &DeploymentsClient{
		client: newDeploymentClient(),
	}
}

func (c *DeploymentsClient) ListDeployments(
	apiKey string,
	projectSlug string,
) *deployments.ListDeploymentResult {
	ctx := context.Background()
	payload := &deployments.ListDeploymentsPayload{
		ApikeyToken:      &apiKey,
		ProjectSlugInput: &projectSlug,
		SessionToken:     nil,
		Cursor:           nil,
	}

	result, err := c.client.ListDeployments(ctx, payload)
	if err != nil {
		log.Fatalf("Error calling ListDeployments: %v", err)
	}

	return result
}

// DeploymentCreator represents a request for creating a deployment
type DeploymentCreator interface {
	CredentialGetter

	// GetIdempotencyKey returns a unique identifier that will mitigate against
	// duplicate deployments.
	GetIdempotencyKey() string

	// GetOpenAPIv3Assets returns the OpenAPI v3 assets to include in the
	// deployment.
	GetOpenAPIv3Assets() []*deployments.AddOpenAPIv3DeploymentAssetForm
}

// CreateDeployment creates a remote deployment.
func (c *DeploymentsClient) CreateDeployment(
	dc DeploymentCreator,
) (*deployments.CreateDeploymentResult, error) {
	ctx := context.Background()

	apiKey := dc.GetApiKey()
	projectSlug := dc.GetProjectSlug()

	payload := &deployments.CreateDeploymentPayload{
		ApikeyToken:      &apiKey,
		ProjectSlugInput: &projectSlug,
		IdempotencyKey:   dc.GetIdempotencyKey(),
		Openapiv3Assets:  dc.GetOpenAPIv3Assets(),
		SessionToken:     nil,
		GithubRepo:       nil,
		GithubPr:         nil,
		GithubSha:        nil,
		ExternalID:       nil,
		ExternalURL:      nil,
		Packages:         nil,
	}

	result, err := c.client.CreateDeployment(ctx, payload)
	if err != nil {
		enhancedErr := enhanceDeploymentError(err)
		return nil, fmt.Errorf("failed to create deployment: %w", enhancedErr)
	}

	return result, nil
}

// enhanceDeploymentError provides more context for deployment errors,
// especially decode errors that may indicate server issues.
func enhanceDeploymentError(err error) error {
	errStr := err.Error()

	// Check if this is a decode error that suggests the server returned HTML instead of JSON
	if strings.Contains(errStr, "can't decode") && strings.Contains(errStr, "text/html") {
		return fmt.Errorf("%w\n\nThis error typically occurs when the server returns an HTML error page (e.g., 500 error) instead of the expected JSON response. Check server logs or try again later", err)
	}

	return err
}

func newDeploymentClient() *deployments.Client {
	h := deploymentService()
	return deployments.NewClient(
		h.GetDeployment(),
		h.GetLatestDeployment(),
		h.CreateDeployment(),
		h.Evolve(),
		h.Redeploy(),
		h.ListDeployments(),
		h.GetDeploymentLogs(),
	)
}

func deploymentService() *depl_client.Client {
	doer := &debugHTTPClient{client: goaSharedHTTPClient}

	scheme := env.Fallback("GRAM_SCHEME", "https")
	host := env.Must("GRAM_HOST")
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := true // Enable body restoration to allow reading raw response on decode errors

	return depl_client.NewClient(scheme, host, doer, enc, dec, restoreBody)
}

// debugHTTPClient wraps the HTTP client to log response details for debugging
type debugHTTPClient struct {
	client *http.Client
}

func (d *debugHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Log outgoing request details for deployment endpoints
	if strings.Contains(req.URL.Path, "deployments") {
		fmt.Printf("DEBUG REQUEST: %s %s\n", req.Method, req.URL.String())

		// Log request headers, filtering out sensitive information
		fmt.Printf("DEBUG REQUEST HEADERS:\n")
		for name, values := range req.Header {
			// Never log the Gram-Key header for security reasons
			if strings.ToLower(name) == "gram-key" {
				fmt.Printf("  %s: [REDACTED]\n", name)
			} else {
				fmt.Printf("  %s: %v\n", name, values)
			}
		}

		// Log request body if present
		if req.Body != nil && req.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(req.Body)
			if err == nil {
				fmt.Printf("DEBUG REQUEST BODY:\n%s\n", bodyBytes)
				// Restore the body so the request can still be sent
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error making HTTP request: %w", err)
	}

	// Log response details for deployment endpoints
	if strings.Contains(req.URL.Path, "deployments") {
		fmt.Printf("DEBUG RESPONSE: %s %s -> HTTP %d, Content-Type: %s\n",
			req.Method, req.URL.Path, resp.StatusCode, resp.Header.Get("Content-Type"))
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG RESPONSE BODY:\n%s\n", bodyBytes)

		// rewind so Goa can still try decoding
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return resp, nil
}
