package senzingchatservice

import (
	"context"
	_ "embed"
	"sync"

	"github.com/senzing/g2-sdk-go/g2api"
	"github.com/senzing/go-logging/logging"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/go-sdk-abstract-factory/factory"
	"github.com/senzing/serve-chat/senzingchatapi"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ChatApiServiceImpl is...
type ChatApiServiceImpl struct {
	senzingchatapi.UnimplementedHandler
	abstractFactory                factory.SdkAbstractFactory
	abstractFactorySyncOnce        sync.Once
	g2configmgrSingleton           g2api.G2configmgr
	g2configmgrSyncOnce            sync.Once
	g2configSingleton              g2api.G2config
	g2configSyncOnce               sync.Once
	g2productSingleton             g2api.G2product
	g2productSyncOnce              sync.Once
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	isTrace                        bool
	logger                         logging.LoggingInterface
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecificationSpec       []byte
	Port                           int
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
	UrlRoutePrefix                 string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

var debugOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

var traceOptions []interface{} = []interface{}{
	&logging.OptionCallerSkip{Value: 5},
}

// ----------------------------------------------------------------------------
// internal methods
// ----------------------------------------------------------------------------

// --- Logging ----------------------------------------------------------------

// Get the Logger singleton.
func (restApiService *ChatApiServiceImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if restApiService.logger == nil {
		loggerOptions := []interface{}{
			&logging.OptionCallerSkip{Value: 3},
		}
		restApiService.logger, err = logging.NewSenzingToolsLogger(ComponentId, IdMessages, loggerOptions...)
		if err != nil {
			panic(err)
		}
	}
	return restApiService.logger
}

// Log message.
func (restApiService *ChatApiServiceImpl) log(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Debug.
func (restApiService *ChatApiServiceImpl) debug(messageNumber int, details ...interface{}) {
	details = append(details, debugOptions...)
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method entry.
func (restApiService *ChatApiServiceImpl) traceEntry(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// Trace method exit.
func (restApiService *ChatApiServiceImpl) traceExit(messageNumber int, details ...interface{}) {
	restApiService.getLogger().Log(messageNumber, details...)
}

// --- Errors -----------------------------------------------------------------

// Create error.
func (restApiService *ChatApiServiceImpl) error(messageNumber int, details ...interface{}) error {
	return restApiService.getLogger().NewError(messageNumber, details...)
}

// --- Services ---------------------------------------------------------------

func (restApiService *ChatApiServiceImpl) getAbstractFactory() factory.SdkAbstractFactory {
	restApiService.abstractFactorySyncOnce.Do(func() {
		if len(restApiService.GrpcTarget) == 0 {
			restApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{}
		} else {
			restApiService.abstractFactory = &factory.SdkAbstractFactoryImpl{
				GrpcDialOptions: restApiService.GrpcDialOptions,
				GrpcTarget:      restApiService.GrpcTarget,
				ObserverOrigin:  restApiService.ObserverOrigin,
				Observers:       restApiService.Observers,
			}
		}
	})
	return restApiService.abstractFactory
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *ChatApiServiceImpl) getG2config(ctx context.Context) g2api.G2config {
	var err error = nil
	restApiService.g2configSyncOnce.Do(func() {
		restApiService.g2configSingleton, err = restApiService.getAbstractFactory().GetG2config(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2configSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2configSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2configSingleton
}

// Singleton pattern for g2config.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *ChatApiServiceImpl) getG2configmgr(ctx context.Context) g2api.G2configmgr {
	var err error = nil
	restApiService.g2configmgrSyncOnce.Do(func() {
		restApiService.g2configmgrSingleton, err = restApiService.getAbstractFactory().GetG2configmgr(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2configmgrSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2configmgrSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2configmgrSingleton
}

// Singleton pattern for g2product.
// See https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f
func (restApiService *ChatApiServiceImpl) getG2product(ctx context.Context) g2api.G2product {
	var err error = nil
	restApiService.g2productSyncOnce.Do(func() {
		restApiService.g2productSingleton, err = restApiService.getAbstractFactory().GetG2product(ctx)
		if err != nil {
			panic(err)
		}
		if restApiService.g2productSingleton.GetSdkId(ctx) == factory.ImplementedByBase {
			err = restApiService.g2productSingleton.Init(ctx, restApiService.SenzingModuleName, restApiService.SenzingEngineConfigurationJson, restApiService.SenzingVerboseLogging)
			if err != nil {
				panic(err)
			}
		}
	})
	return restApiService.g2productSingleton
}

// ----------------------------------------------------------------------------
// Interface methods
// See https://github.com/docktermj/go-rest-api-client/blob/main/senzingrestpapi/oas_unimplemented_gen.go
// ----------------------------------------------------------------------------
