package senzingchatservice

import (
	"context"
	"sync"

	"github.com/senzing-garage/g2-sdk-go/g2api"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/go-sdk-abstract-factory/factory"
	api "github.com/senzing-garage/serve-chat/senzingchatapi"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ChatApiServiceImpl is...
type ChatApiServiceImpl struct {
	api.UnimplementedHandler
	abstractFactory                factory.SdkAbstractFactory
	abstractFactorySyncOnce        sync.Once
	g2engineSingleton              g2api.G2engine
	g2engineSyncOnce               sync.Once
	g2productSingleton             g2api.G2product
	g2productSyncOnce              sync.Once
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

func (chatApiService *ChatApiServiceImpl) getAbstractFactory() factory.SdkAbstractFactory {
	chatApiService.abstractFactorySyncOnce.Do(func() {
		if len(chatApiService.GrpcTarget) == 0 {
			chatApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{}
		} else {
			chatApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{
				GrpcDialOptions: chatApiService.GrpcDialOptions,
				GrpcTarget:      chatApiService.GrpcTarget,
				ObserverOrigin:  chatApiService.ObserverOrigin,
				Observers:       chatApiService.Observers,
			}
		}
	})
	return chatApiService.abstractFactory
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (chatApiService *ChatApiServiceImpl) getG2engine(ctx context.Context) g2api.G2engine {
	var err error = nil
	chatApiService.g2engineSyncOnce.Do(func() {
		chatApiService.g2engineSingleton, err = chatApiService.getAbstractFactory().GetG2engine(ctx)
		if err != nil {
			panic(err)
		}
		if chatApiService.g2engineSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = chatApiService.g2engineSingleton.Init(ctx, chatApiService.SenzingModuleName, chatApiService.SenzingEngineConfigurationJson, chatApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return chatApiService.g2engineSingleton
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (chatApiService *ChatApiServiceImpl) getG2product(ctx context.Context) g2api.G2product {
	var err error = nil
	chatApiService.g2productSyncOnce.Do(func() {
		chatApiService.g2productSingleton, err = chatApiService.getAbstractFactory().GetG2product(ctx)
		if err != nil {
			panic(err)
		}
		if chatApiService.g2productSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = chatApiService.g2productSingleton.Init(ctx, chatApiService.SenzingModuleName, chatApiService.SenzingEngineConfigurationJson, chatApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return chatApiService.g2productSingleton
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
