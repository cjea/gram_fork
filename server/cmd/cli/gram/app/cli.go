package app

import (
	"fmt"

	"github.com/speakeasy-api/gram/server/cmd/cli/env"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/api"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/depconfig"
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
		},
		Action: mainAction,
	}

	return &cliApp{app: app}
}

func (c *cliApp) Run(args []string) error {
	return c.app.Run(args)
}

func mainAction(c *cli.Context) error {
	fmt.Printf("Starting CLI.\n")

	filePath := c.String("file")
	project, err := depconfig.ReadDeploymentConfig(filePath)
	if err != nil {
		return fmt.Errorf("error reading project config: %w", err)
	}

	fmt.Printf("Loaded project: %s\n", project.Project)
	fmt.Printf("Sources: %+v\n", project.Sources)

	apiKey := api.ApiKeyFromEnv()
	projectSlug := env.Must("GRAM_PROJECT_SLUG")

	client := api.NewClient()
	result := client.ListDeployments(apiKey, projectSlug)
	printDeployments(result)

	return nil
}

func printDeployments(ds *deployments.ListDeploymentResult) {
	for i, deployment := range ds.Items {
		fmt.Printf("  [%d] %+v\n", i+1, deployment)
	}
}
