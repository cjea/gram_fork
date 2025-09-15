package app

import (
	"fmt"

	"github.com/speakeasy-api/gram/server/cmd/cli/gram/deploy"
	"github.com/speakeasy-api/gram/server/cmd/cli/gram/env"
	"github.com/urfave/cli/v2"
)

const Version = "0.1.0"

type CLI interface {
	Run(args []string) error
}

type cliApp struct {
	app *cli.App
}

func NewCLI() CLI {
	pushUsage := `Push a deployment to Gram.

Sample deployment file
======================
{
  "schema_version": "1.0.0",
  "type": "deployment",
  "sources": [
    {
      "type": "openapiv3",
      "location": "/path/to/spec.yaml",
      "name": "My API",
      "slug": "my-api"
    }
  ]
}

NOTE: Names and slugs must be unique across all sources.
`

	app := &cli.App{
		Name:    "gram",
		Usage:   "Remote MCP management",
		Version: Version,
		Commands: []*cli.Command{
			{
				Name:  "push",
				Usage: pushUsage,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "Path to the deployment file (relative locations resolve to the deployment file's directory)",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "project",
						Aliases: []string{"p"},
						EnvVars: []string{env.VarNameProjectSlug},
						Usage: fmt.Sprintf(
							"Project slug (falls back to %s environment variable)",
							env.VarNameProjectSlug),
					},
				},
				Action: pushAction,
			},
		},
	}

	return &cliApp{app: app}
}

func (c *cliApp) Run(args []string) error {
	if err := c.app.Run(args); err != nil {
		return fmt.Errorf("failed to run CLI app: %w", err)
	}
	return nil
}

func pushAction(c *cli.Context) error {
	filePath := c.String("file")
	projectSlug := c.String("project")

	fmt.Printf("Deploying to project: %s\n", projectSlug)

	req := deploy.CreateDeploymentFromFileRequest{
		FilePath: filePath,
		Project:  projectSlug,
	}
	result, err := deploy.CreateDeploymentFromFile(req)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Printf("Deployment created successfully: %+v\n", result.Deployment)
	return nil
}
