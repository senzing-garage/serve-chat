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
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-chat/httpserver"
	"github.com/senzing-garage/serve-chat/senzingchatservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	Long string = `
 serve-chat long description.
     `
	ReadHeaderTimeoutInSeconds        = 60
	Short                      string = "serve-chat short description"
	Use                        string = "serve-chat"
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.AvoidServe,
	option.Configuration,
	option.DatabaseURL,
	option.EnableAll,
	option.EnableSenzingChatAPI,
	option.EnableSwaggerUI,
	option.EngineInstanceName,
	option.EngineLogLevel,
	option.EngineSettings,
	option.GrpcURL,
	option.HTTPPort,
	option.LogLevel,
	option.ObserverOrigin,
	option.ObserverURL,
	option.ServerAddress,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

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

// Used in construction of cobra.Command.
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command.
func RunE(_ *cobra.Command, _ []string) error {
	var err error

	ctx := context.Background()

	senzingEngineConfigurationJSON, err := settings.BuildAndVerifySettings(ctx, viper.GetViper())
	if err != nil {
		return wraperror.Errorf(err, "BuildAndVerifySettings")
	}

	// Determine if gRPC is being used.

	grpcURL := viper.GetString(option.GrpcURL.Arg)
	grpcTarget := ""
	grpcDialOptions := []grpc.DialOption{}

	if len(grpcURL) > 0 {
		grpcTarget, grpcDialOptions, err = grpcurl.Parse(ctx, grpcURL)
		if err != nil {
			return wraperror.Errorf(err, "grpcurl.Parse: %s", grpcURL)
		}
	}

	// Build observers.
	//  viper.GetString(option.ObserverUrl),

	observers := []observer.Observer{}

	// Create object and Serve.

	httpServer := &httpserver.BasicHTTPServer{
		AvoidServing:          viper.GetBool(option.AvoidServe.Arg),
		ChatURLRoutePrefix:    "chat",
		EnableAll:             viper.GetBool(option.EnableAll.Arg),
		EnableSenzingChatAPI:  viper.GetBool(option.EnableSenzingChatAPI.Arg),
		EnableSwaggerUI:       viper.GetBool(option.EnableSwaggerUI.Arg),
		GrpcDialOptions:       grpcDialOptions,
		GrpcTarget:            grpcTarget,
		LogLevelName:          viper.GetString(option.LogLevel.Arg),
		ObserverOrigin:        viper.GetString(option.ObserverOrigin.Arg),
		Observers:             observers,
		OpenAPISpecification:  senzingchatservice.OpenAPISpecificationJSON,
		ReadHeaderTimeout:     ReadHeaderTimeoutInSeconds * time.Second,
		Setting:               senzingEngineConfigurationJSON,
		SenzingInstanceName:   viper.GetString(option.EngineInstanceName.Arg),
		SenzingVerboseLogging: viper.GetInt64(option.EngineLogLevel.Arg),
		ServerAddress:         viper.GetString(option.ServerAddress.Arg),
		ServerPort:            viper.GetInt(option.HTTPPort.Arg),
		SwaggerURLRoutePrefix: "swagger",
	}

	err = httpServer.Serve(ctx)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// Used in construction of cobra.Command.
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}
