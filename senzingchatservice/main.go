package senzingchatservice

import (
	_ "embed"

	"github.com/senzing/serve-chat/senzingchatapi"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The ChatApiService interface is...
type ChatApiService interface {
	senzingchatapi.Handler
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the  package found messages having the format "senzing-6620xxxx".
// See https://github.com/Senzing/knowledge-base/blob/main/lists/senzing-component-ids.md
const ComponentId = 6620

// Log message prefix.
const Prefix = "serve-chat.chatapiservice."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	0001: "Example Trace log.",
	1000: "Example Debug log.",
	2000: "Example Info log.",
	3000: "Example Warn log.",
	4000: "Example Error log.",
	5000: "Example Fatal log.",
	6000: "Example Panic log.",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}

//go:embed openapi.json
var OpenApiSpecificationJson []byte

//go:embed openapi.yaml
var OpenApiSpecificationYaml []byte
