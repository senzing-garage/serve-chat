package senzingchatservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/stretchr/testify/assert"
)

var (
	chatApiServiceSingleton ChatApiService
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

func getTestObject(ctx context.Context, test *testing.T) ChatApiService {
	if chatApiServiceSingleton == nil {
		senzingEngineConfigurationJson, err := g2engineconfigurationjson.BuildSimpleSystemConfigurationJson("")
		if err != nil {
			test.Errorf("Error: %s", err)
		}
		chatApiServiceSingleton = &ChatApiServiceImpl{
			SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
			SenzingModuleName:              "go-rest-api-service-test",
			SenzingVerboseLogging:          0,
		}
	}
	return chatApiServiceSingleton
}

func testError(test *testing.T, ctx context.Context, err error) {
	if err != nil {
		test.Log("Error:", err.Error())
		assert.FailNow(test, err.Error())
	}
}

// ----------------------------------------------------------------------------
// Test harness
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	code := m.Run()
	err = teardown()
	if err != nil {
		fmt.Print(err)
	}
	os.Exit(code)
}

func setup() error {
	var err error = nil
	return err
}

func teardown() error {
	var err error = nil
	return err
}

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestChatApiServiceImpl_AddPet(test *testing.T) {
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleChatApiServiceImpl_AddPet() {

}
