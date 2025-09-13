package app

import (
	"fmt"

	"github.com/speakeasy-api/gram/server/cmd/cli/gram/api"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/deplconfig"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/env"
	"github.com/speakeasy-api/gram/server/gen/assets"
	"github.com/speakeasy-api/gram/server/gen/deployments"
	"github.com/urfave/cli/v2"
)

type CLI interface {
	Run(args []string) error
}

type cliApp struct {
	app *cli.App
}

func NewCLI() CLI {
	app := &cli.App{
		Name:  "gram_cli",
		Usage: "Gram CLI tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Path to the deployment configuration file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage: fmt.Sprintf(
					"Project slug (falls back to %s environment variable)",
					env.VarNameProjectSlug),
			},
		},
		Action: mainAction,
	}

	return &cliApp{app: app}
}

func (c *cliApp) Run(args []string) error {
	if err := c.app.Run(args); err != nil {
		return fmt.Errorf("failed to run CLI app: %w", err)
	}
	return nil
}

func mainAction(c *cli.Context) error {
	fmt.Printf("Starting CLI.\n")

	filePath := c.String("file")
	deplconfig, err := deplconfig.ReadDeploymentConfig(filePath)
	if err != nil {
		return fmt.Errorf("error reading project config: %w", err)
	}

	apiKey := env.MustApiKey()

	projectSlug := c.String("project")
	if projectSlug == "" {
		projectSlug = env.MustProjectSlug()
	}

	fmt.Printf("Project: %s\n", projectSlug)
	fmt.Printf("Sources: %+v\n", deplconfig.Sources)

	deplclient := api.NewDeploymentsClient()
	result := deplclient.ListDeployments(apiKey, projectSlug)
	printDeployments(result)

	assetsClient := api.NewAssetsClient()
	assets := assetsClient.ListAssets(apiKey, projectSlug)
	printAssets(assets)

	return nil
}

func printDeployments(ds *deployments.ListDeploymentResult) {
	for i, deployment := range ds.Items {
		fmt.Printf("  [%d] %+v\n", i+1, deployment)
	}
}

func printAssets(as *assets.ListAssetsResult) {
	for i, asset := range as.Assets {
		fmt.Printf("  [%d] %+v\n", i+1, asset)
	}
}
