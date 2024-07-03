package senzingchatservice

import (
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-chat/senzingchatapi"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicChatAPIService is...
type BasicChatAPIService struct {
	senzingchatapi.UnimplementedHandler
	// abstractFactory          senzing.SzAbstractFactory
	// abstractFactorySyncOnce  sync.Once
	GrpcDialOptions []grpc.DialOption
	GrpcTarget      string
	// logger                   logging.Logging
	LogLevelName             string
	ObserverOrigin           string
	Observers                []observer.Observer
	OpenAPISpecificationSpec []byte
	Port                     int
	Settings                 string
	SenzingInstanceName      string
	SenzingVerboseLogging    int64
	// szEngineSingleton        senzing.SzEngine
	// szEngineSyncOnce         sync.Once
	// szProductSingleton       senzing.SzProduct
	// szProductSyncOnce        sync.Once
	URLRoutePrefix string
}

// ----------------------------------------------------------------------------
// internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
// func (chatAPIService *BasicChatAPIService) getLogger() logging.Logging {
// 	var err error
// 	if chatAPIService.logger == nil {
// 		loggerOptions := []interface{}{
// 			&logging.OptionCallerSkip{Value: 3},
// 		}
// 		chatAPIService.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, loggerOptions...)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	return chatAPIService.logger
// }

// Log message.
// func (chatAPIService *BasicChatAPIService) log(messageNumber int, details ...interface{}) {
// 	chatAPIService.getLogger().Log(messageNumber, details...)
// }

// --- Errors -----------------------------------------------------------------

// // Create error.
// func (chatAPIService *BasicChatAPIService) error(messageNumber int, details ...interface{}) error {
// 	return chatAPIService.getLogger().NewError(messageNumber, details...)
// }

// --- Services ---------------------------------------------------------------

// func (chatAPIService *BasicChatAPIService) getAbstractFactory(ctx context.Context) senzing.SzAbstractFactory {
// 	_ = ctx
// 	var err error
// 	chatAPIService.abstractFactorySyncOnce.Do(func() {
// 		if len(chatAPIService.GrpcTarget) == 0 {
// 			chatAPIService.abstractFactory, err = szfactorycreator.CreateCoreAbstractFactory(chatAPIService.SenzingInstanceName, chatAPIService.Settings, chatAPIService.SenzingVerboseLogging, senzing.SzInitializeWithDefaultConfiguration)
// 			if err != nil {
// 				panic(err)
// 			}
// 		} else {
// 			grpcConnection, err := grpc.NewClient(chatAPIService.GrpcTarget, chatAPIService.GrpcDialOptions...)
// 			if err != nil {
// 				panic(err)
// 			}
// 			chatAPIService.abstractFactory, err = szfactorycreator.CreateGrpcAbstractFactory(grpcConnection)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	})
// 	return chatAPIService.abstractFactory
// }

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
// func (chatAPIService *BasicChatAPIService) getG2engine(ctx context.Context) senzing.SzEngine {
// 	var err error
// 	chatAPIService.szEngineSyncOnce.Do(func() {
// 		chatAPIService.szEngineSingleton, err = chatAPIService.getAbstractFactory(ctx).CreateSzEngine(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 	})
// 	return chatAPIService.szEngineSingleton
// }

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
// func (chatAPIService *BasicChatAPIService) getG2product(ctx context.Context) senzing.SzProduct {
// 	var err error
// 	chatAPIService.szProductSyncOnce.Do(func() {
// 		chatAPIService.szProductSingleton, err = chatAPIService.getAbstractFactory(ctx).CreateSzProduct(ctx)
// 		if err != nil {
// 			panic(err)
// 		}
// 	})
// 	return chatAPIService.szProductSingleton
// }

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
// 	var err error
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
