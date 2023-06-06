package senzingchatservice

import "github.com/senzing/serve-chat/senzingchatapi"

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

// Identfier of the  package found messages having the format "senzing-6503xxxx".
const ComponentId = 9999

// Log message prefix.
const Prefix = "serve-chat.chatapiservice."

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for g2config implementations.
var IdMessages = map[int]string{
	10: "Enter " + Prefix + "InitializeSenzing().",
}

// Status strings for specific messages.
var IdStatuses = map[int]string{}
