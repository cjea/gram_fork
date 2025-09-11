package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/speakeasy-api/gram/server/gen/deployments"
	httpclient "github.com/speakeasy-api/gram/server/gen/http/deployments/client"
	goahttp "goa.design/goa/v3/http"
)

var API_KEY string = apiKeyFromEnv()
var PROJECT_SLUG string = mustEnv("GRAM_PROJECT_SLUG")

func main() {
	fmt.Printf("Starting CLI.\n")

	httpClient := httpClientForGoa()
	deploymentClient := newDeploymentClient(httpClient)

	result := listDeployments(deploymentClient)
	printDeployments(result)
}

func printDeployments(ds *deployments.ListDeploymentResult) {
	for i, deployment := range ds.Items {
		fmt.Printf("  [%d] %+v\n", i+1, deployment)
	}
}

func listDeployments(d *deployments.Client) *deployments.ListDeploymentResult {
	ctx := context.Background()
	payload := &deployments.ListDeploymentsPayload{
		ApikeyToken:      &API_KEY,
		ProjectSlugInput: &PROJECT_SLUG,
	}

	result, err := d.ListDeployments(ctx, payload)
	if err != nil {
		log.Fatalf("Error calling ListDeployments: %v", err)
	}

	return result
}

func newDeploymentClient(h *httpclient.Client) *deployments.Client {
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

func httpClientForGoa() *httpclient.Client {
	scheme := envOr("GRAM_SCHEME", "https")
	host := mustEnv("GRAM_HOST")
	doer := &http.Client{}
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	return httpclient.NewClient(scheme, host, doer, enc, dec, restoreBody)
}

func apiKeyFromEnv() string {
	return validateApiKey(mustEnv("GRAM_API_KEY"))
}

func validateApiKey(key string) string {
	ok := strings.HasPrefix(key, "gram")

	if ok {
		return key
	} else {
		panic(fmt.Errorf("key is malformed"))
	}
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		panic(fmt.Errorf("missing env: %s", key))
	}

	return val
}

func envOr(key string, fallback string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		return fallback
	} else {
		return val
	}
}
