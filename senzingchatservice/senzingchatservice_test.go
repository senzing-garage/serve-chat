package senzingchatservice

import (
	"testing"
)

var (
// chatAPIServiceSingleton ChatAPIService
)

// ----------------------------------------------------------------------------
// Internal functions
// ----------------------------------------------------------------------------

// func getTestObject(ctx context.Context, test *testing.T) ChatAPIService {
// 	_ = ctx
// 	if chatAPIServiceSingleton == nil {
// 		senzingEngineConfigurationJSON, err := settings.BuildSimpleSettingsUsingEnvVars()
// 		require.NoError(test, err)
// 		chatAPIServiceSingleton = &BasicChatAPIService{
// 			Settings:              senzingEngineConfigurationJSON,
// 			SenzingInstanceName:   "go-rest-api-service-test",
// 			SenzingVerboseLogging: 0,
// 		}
// 	}
// 	return chatAPIServiceSingleton
// }

// ----------------------------------------------------------------------------
// Test interface functions
// ----------------------------------------------------------------------------

func TestBasicChatAPIService_AddPet(test *testing.T) {
	_ = test
}

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

// func ExampleBasicChatAPIService_AddPet() {

// }
