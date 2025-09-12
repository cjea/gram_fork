package api

import (
	"context"
	"log"

	"github.com/speakeasy-api/gram/server/cmd/cli/env"
	"github.com/speakeasy-api/gram/server/gen/assets"
	assets_client "github.com/speakeasy-api/gram/server/gen/http/assets/client"
	goahttp "goa.design/goa/v3/http"
)

type AssetsClient struct {
	client *assets.Client
}

func NewAssetsClient() *AssetsClient {
	return &AssetsClient{
		client: newAssetsClient(),
	}
}

func (c *AssetsClient) ListAssets(apiKey, projectSlug string) *assets.ListAssetsResult {
	ctx := context.Background()
	payload := &assets.ListAssetsPayload{
		ApikeyToken:      &apiKey,
		ProjectSlugInput: &projectSlug,
		SessionToken:     nil,
	}

	result, err := c.client.ListAssets(ctx, payload)
	if err != nil {
		log.Fatalf("Error calling ListAssets: %v", err)
	}

	return result
}

func newAssetsClient() *assets.Client {
	h := assetsService()
	return assets.NewClient(
		h.ServeImage(),
		h.UploadImage(),
		h.UploadFunctions(),
		h.UploadOpenAPIv3(),
		h.ServeOpenAPIv3(),
		h.ListAssets(),
	)
}

func assetsService() *assets_client.Client {
	doer := goaSharedHTTPClient

	scheme := env.Fallback("GRAM_SCHEME", "https")
	host := env.Must("GRAM_HOST")
	enc := goahttp.RequestEncoder
	dec := goahttp.ResponseDecoder
	restoreBody := false

	return assets_client.NewClient(scheme, host, doer, enc, dec, restoreBody)
}
