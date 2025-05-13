package httpserver_test

import (
	"testing"

	"github.com/senzing-garage/serve-chat/httpserver"
)

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestHTTPServerImpl_Serve(test *testing.T) {
	_ = test

	_ = httpserver.BasicHTTPServer{}
}
