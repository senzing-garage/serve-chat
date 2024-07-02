package senzingchatservice

import (
	"context"
	"sync"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	api "github.com/senzing-garage/serve-chat/senzingchatapi"
	"github.com/senzing-garage/sz-sdk-go/senzing"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ChatApiServiceImpl is...
type ChatApiServiceImpl struct {
	api.UnimplementedHandler
	abstractFactory                senzing.SzAbstractFactory
	abstractFactorySyncOnce        sync.Once
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	logger                         logging.LoggingInterface
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecificationSpec       []byte
	Port                           int
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int64
	szEngineSingleton              senzing.SzEngine
	szEngineSyncOnce               sync.Once
	szProductSingleton             senzing.SzProduct
	szProductSyncOnce              sync.Once
	UrlRoutePrefix                 string
}

// ----------------------------------------------------------------------------
// internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (chatApiService *ChatApiServiceImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if chatApiService.logger == nil {
		loggerOptions := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		chatApiService.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, loggerOptions...)
		if err != nil {
			panic(err)
		}
	}
	return chatApiService.logger
}

// Log message.
func (chatApiService *ChatApiServiceImpl) log(messageNumber int, details ...interface{}) {
	chatApiService.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (chatApiService *ChatApiServiceImpl) error(messageNumber int, details ...interface{}) error {
	return chatApiService.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

func (chatApiService *ChatApiServiceImpl) getAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
	var err error = nil
	chatApiService.abstractFactorySyncOnce.Do(func() {
		if len(chatApiService.GrpcTarget) == 0 {
			chatApiService.abstractFactory, err = szfactorycreator.CreateCoreAbstractFactory(chatApiService.SenzingModuleName, chatApiService.SenzingEngineConfigurationJson, chatApiService.SenzingVerboseLogging, senzing.SzInitializeWithDefaultConfiguration)
			if err != nil {
				panic(err)
			}
		} else {
			grpcConnection, err := grpc.DialContext(ctx, chatApiService.GrpcTarget, chatApiService.GrpcDialOptions...)
			if err != nil {
				panic(err)
			}
			chatApiService.abstractFactory, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
			if err != nil {
				panic(err)
			}
		}
	})
	return chatApiService.abstractFactory
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (chatApiService *ChatApiServiceImpl) getG2engine(ctx context.Context) senzing.SzEngine {
	var err error = nil
	chatApiService.szEngineSyncOnce.Do(func() {
		chatApiService.szEngineSingleton, err = chatApiService.getAbstractFactory(ctx).CreateSzEngine(ctx)
		if err != nil {
			panic(err)
		}
	})
	return chatApiService.szEngineSingleton
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (chatApiService *ChatApiServiceImpl) getG2product(ctx context.Context) senzing.SzProduct {
	var err error = nil
	chatApiService.szProductSyncOnce.Do(func() {
		chatApiService.szProductSingleton, err = chatApiService.getAbstractFactory(ctx).CreateSzProduct(ctx)
		if err != nil {
			panic(err)
		}
	})
	return chatApiService.szProductSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// See https://github.com/senzing-garage/serve-chat/blob/main/senzingchatpapi/oas_unimplemented_gen.go
// ----------------------------------------------------------------------------

// AddPet implements addPet operation.
//
// Add a new pet to the store.
//
// POST /pet
// func (chatApiService *ChatApiServiceImpl) AddPet(ctx context.Context, req *api.Pet) (r *api.Pet, _ error) {
// 	response, err := chatApiService.getG2product(ctx).Version(ctx)
// 	if err != nil {
// 		return r, err
// 	}
// 	parsedResponse, err := senzing.UnmarshalProductVersionResponse(ctx, response)
// 	if err != nil {
// 		return r, err
// 	}
// 	r = &api.Pet{
// 		ID:     api.NewOptInt64(1001),
// 		Name:   parsedResponse.BuildVersion,
// 		Status: api.NewOptPetStatus(api.PetStatusAvailable),
// 	}

// 	// Example logging.

// 	chatApiService.log(1, r, err)

// 	// Example error generation.

// 	newErr := chatApiService.error(2, "example error")
// 	if false {
// 		fmt.Printf(">>> An example error: %+v\n", newErr)
// 	}
// 	return r, err
// }

// GetPetById implements getPetById operation.
//
// Returns a single pet.
//
// GET /pet/{petId}
// func (chatApiService *ChatApiServiceImpl) GetPetById(ctx context.Context, params api.GetPetByIdParams) (r api.GetPetByIdRes, _ error) {
// 	var err error = nil
// 	response, err := chatApiService.getG2engine(ctx).GetEntityByEntityID(ctx, params.PetId)
// 	if err != nil {
// 		return r, err
// 	}

// 	r = &api.Pet{
// 		ID:     api.NewOptInt64(params.PetId),
// 		Name:   response,
// 		Status: api.NewOptPetStatus(api.PetStatusAvailable),
// 	}
// 	return r, err
// }
