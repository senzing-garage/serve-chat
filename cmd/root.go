/*
 */
package cmd

import (
	"context"
	"os"
	"time"

	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/settings"
	"github.com/senzing-garage/go-grpcing/grpcurl"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-chat/httpserver"
	"github.com/senzing-garage/serve-chat/senzingchatservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	Short string = "serve-chat short description"
	Use   string = "serve-chat"
	Long  string = `
 serve-chat long description.
	 `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.DatabaseUrl,
	option.EnableAll,
	option.EnableSenzingChatApi,
	option.EnableSwaggerUi,
	option.EngineConfigurationJson,
	option.EngineLogLevel,
	option.EngineModuleName,
	option.GrpcUrl,
	option.HttpPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverUrl,
	option.ServerAddress,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.Background()

	senzingEngineConfigurationJson, err := settings.BuildAndVerifySenzingEngineConfigurationJson(ctx, viper.GetViper())
	if err != nil {
		return err
	}

	// Determine if gRPC is being used.

	grpcUrl := viper.GetString(option.GrpcUrl.Arg)
	grpcTarget := ""
	grpcDialOptions := []grpc.DialOption{}
	if len(grpcUrl) > 0 {
		grpcTarget, grpcDialOptions, err = grpcurl.Parse(ctx, grpcUrl)
		if err != nil {
			return err
		}
	}

	// Build observers.
	//  viper.GetString(option.ObserverUrl),

	observers := []observer.Observer{}

	// Create object and Serve.

	httpServer := &httpserver.HttpServerImpl{
		ChatUrlRoutePrefix:             "chat",
		EnableAll:                      viper.GetBool(option.EnableAll.Arg),
		EnableSenzingChatAPI:           viper.GetBool(option.EnableSenzingChatApi.Arg),
		EnableSwaggerUI:                viper.GetBool(option.EnableSwaggerUi.Arg),
		GrpcDialOptions:                grpcDialOptions,
		GrpcTarget:                     grpcTarget,
		LogLevelName:                   viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin.Arg),
		Observers:                      observers,
		OpenApiSpecification:           senzingchatservice.OpenApiSpecificationJson,
		ReadHeaderTimeout:              60 * time.Second,
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName.Arg),
		SenzingVerboseLogging:          viper.GetInt64(option.EngineLogLevel.Arg),
		ServerAddress:                  viper.GetString(option.ServerAddress.Arg),
		ServerPort:                     viper.GetInt(option.HttpPort.Arg),
		SwaggerUrlRoutePrefix:          "swagger",
	}
	return httpServer.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
