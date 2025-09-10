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

func main() {
	fmt.Printf("Starting CLI.")
	scheme := "https"
	host := "app.getgram.ai"

	doer := &http.Client{}
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	httpClient := httpclient.NewClient(scheme, host, doer, enc, dec, restoreBody)
	client := deployments.NewClient(
		httpClient.GetDeployment(),
		httpClient.GetLatestDeployment(),
		httpClient.CreateDeployment(),
		httpClient.Evolve(),
		httpClient.Redeploy(),
		httpClient.ListDeployments(),
		httpClient.GetDeploymentLogs(),
	)

	ctx := context.Background()
	payload := &deployments.ListDeploymentsPayload{ApikeyToken: &API_KEY}

	result, err := client.ListDeployments(ctx, payload)
	if err != nil {
		log.Fatalf("Error calling ListDeployments: %v", err)
	}

	fmt.Printf("Success! Got %d deployments\n", len(result.Items))
	for i, deployment := range result.Items {
		fmt.Printf("  [%d] %+v\n", i+1, deployment)
	}
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
