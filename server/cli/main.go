package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/speakeasy-api/gram/server/gen/deployments"
	depl_client "github.com/speakeasy-api/gram/server/gen/http/deployments/client"
	"github.com/urfave/cli/v2"
	goahttp "goa.design/goa/v3/http"
)

type Source struct {
	Type string `json:"type"`
	Loc  string `json:"loc"`
}

type CliRequest struct {
	Project string   `json:"project"`
	Sources []Source `json:"sources"`
}

var goaDoer = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

var fileFlagName = "file"
var fileFlag = &cli.StringFlag{
	Name:     fileFlagName,
	Aliases:  []string{"f"},
	Usage:    "Path to the project configuration file",
	Value:    "server/cli/test.json",
	Required: false,
}

func main() {
	app := &cli.App{
		Name:   "gram_cli",
		Usage:  "Gram CLI tool",
		Flags:  []cli.Flag{fileFlag},
		Action: mainAction,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainAction(c *cli.Context) error {
	fmt.Printf("Starting CLI.\n")

	filePath := c.String(fileFlagName)
	project, err := readProjectConfig(filePath)
	if err != nil {
		return fmt.Errorf("error reading project config: %w", err)
	}

	fmt.Printf("Loaded project: %s\n", project.Project)
	fmt.Printf("Sources: %+v\n", project.Sources)

	apiKey := apiKeyFromEnv()
	projectSlug := mustEnv("GRAM_PROJECT_SLUG")

	deploymentClient := newDeploymentClient()

	result := listDeployments(deploymentClient, apiKey, projectSlug)
	printDeployments(result)

	return nil
}

func readProjectConfig(filePath string) (*CliRequest, error) {
	// #nosec G304 -- file path is controlled by CLI flag, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var project CliRequest
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &project, nil
}

func printDeployments(ds *deployments.ListDeploymentResult) {
	for i, deployment := range ds.Items {
		fmt.Printf("  [%d] %+v\n", i+1, deployment)
	}
}

func listDeployments(d *deployments.Client, apiKey, projectSlug string) *deployments.ListDeploymentResult {
	ctx := context.Background()
	payload := &deployments.ListDeploymentsPayload{
		ApikeyToken:      &apiKey,
		SessionToken:     nil,
		ProjectSlugInput: &projectSlug,
		Cursor:           nil,
	}

	result, err := d.ListDeployments(ctx, payload)
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
	scheme := envOr("GRAM_SCHEME", "https")
	host := mustEnv("GRAM_HOST")
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	return depl_client.NewClient(scheme, host, goaDoer, enc, dec, restoreBody)
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
