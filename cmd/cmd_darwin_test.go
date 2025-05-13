//go:build darwin

package cmd_test

import (
	"bytes"
	"testing"

	"github.com/senzing-garage/serve-chat/cmd"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test private functions
// ----------------------------------------------------------------------------

func Test_docsAction(test *testing.T) {
	test.Parallel()
	var buffer bytes.Buffer
	err := cmd.DocsAction(&buffer, "/tmp")
	require.NoError(test, err)
}
