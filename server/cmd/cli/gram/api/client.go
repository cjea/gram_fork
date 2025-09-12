package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/speakeasy-api/gram/server/cmd/cli/env"
	"github.com/speakeasy-api/gram/server/gen/deployments"
	depl_client "github.com/speakeasy-api/gram/server/gen/http/deployments/client"
	goahttp "goa.design/goa/v3/http"
)

var goaDoer = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

type Client struct {
	deploymentClient *deployments.Client
}

func NewClient() *Client {
	return &Client{
		deploymentClient: newDeploymentClient(),
	}
}

func (c *Client) ListDeployments(apiKey, projectSlug string) *deployments.ListDeploymentResult {
	ctx := context.Background()
	payload := &deployments.ListDeploymentsPayload{
		ApikeyToken:      &apiKey,
		SessionToken:     nil,
		ProjectSlugInput: &projectSlug,
		Cursor:           nil,
	}

	result, err := c.deploymentClient.ListDeployments(ctx, payload)
	if err != nil {
		log.Fatalf("Error calling ListDeployments: %v", err)
	}

	return result
}

func newDeploymentClient() *deployments.Client {
	h := httpClientForGoa()
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

func httpClientForGoa() *depl_client.Client {
	scheme := env.Fallback("GRAM_SCHEME", "https")
	host := env.Must("GRAM_HOST")
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	return depl_client.NewClient(scheme, host, goaDoer, enc, dec, restoreBody)
}

func ApiKeyFromEnv() string {
	return validateApiKey(env.Must("GRAM_API_KEY"))
}

func validateApiKey(key string) string {
	ok := strings.HasPrefix(key, "gram")

	if ok {
		return key
	} else {
		panic(fmt.Errorf("key is malformed"))
	}
}
