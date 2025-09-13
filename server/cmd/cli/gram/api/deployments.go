package api

import (
	"context"
	"fmt"
	"log"

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

func (c *DeploymentsClient) ListDeployments(apiKey, projectSlug string) *deployments.ListDeploymentResult {
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

	// GetIdempotencyKey returns a unique identifier that will mitigate against duplicate deployments.
	GetIdempotencyKey() string

	// GetOpenAPIv3Assets returns the OpenAPI v3 assets to include in the deployment.
	GetOpenAPIv3Assets() []*deployments.AddOpenAPIv3DeploymentAssetForm
}

func (c *DeploymentsClient) CreateDeployment(dc DeploymentCreator) (*deployments.CreateDeploymentResult, error) {
	ctx := context.Background()

	apiKey := dc.GetApiKey()
	projectSlug := dc.GetProjectSlug()

	payload := &deployments.CreateDeploymentPayload{
		ApikeyToken:      &apiKey,
		ProjectSlugInput: &projectSlug,
		SessionToken:     nil,
		IdempotencyKey:   dc.GetIdempotencyKey(),
		Openapiv3Assets:  dc.GetOpenAPIv3Assets(),
	}

	result, err := c.client.CreateDeployment(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return result, nil
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
	doer := goaSharedHTTPClient

	scheme := env.Fallback("GRAM_SCHEME", "https")
	host := env.Must("GRAM_HOST")
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	return depl_client.NewClient(scheme, host, doer, enc, dec, restoreBody)
}
