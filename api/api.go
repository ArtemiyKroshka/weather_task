// Package api embeds the OpenAPI specification for the Weather Subscription API.
package api

import _ "embed"

//go:embed openapi.yaml
var OpenAPISpec []byte
