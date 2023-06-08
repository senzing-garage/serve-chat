/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/senzing/go-common/g2engineconfigurationjson"
	"github.com/senzing/go-grpcing/grpcurl"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/senzing-tools/constant"
	"github.com/senzing/senzing-tools/envar"
	"github.com/senzing/senzing-tools/helper"
	"github.com/senzing/senzing-tools/option"
	"github.com/senzing/serve-chat/httpserver"
	"github.com/senzing/serve-chat/senzingchatservice"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	defaultConfiguration           string = ""
	defaultDatabaseUrl             string = ""
	defaultEnableAll               bool   = false
	defaultEnableSenzingChatApi    bool   = false
	defaultEnableSwaggerUI         bool   = false
	defaultEngineConfigurationJson string = ""
	defaultEngineLogLevel          int    = 0
	defaultGrpcUrl                        = ""
	defaultHttpPort                int    = 8262
	defaultLogLevel                string = "INFO"
	defaultObserverOrigin          string = "serve-chat"
	defaultObserverUrl             string = ""
	defaultServerAddress           string = "0.0.0.0"
	envarEnableAll                 string = "SENZING_TOOLS_ENABLE_ALL"
	envarEnableSenzingChatApi      string = "SENZING_TOOLS_ENABLE_SENZING_CHAT_API"
	envarServerAddress             string = "SENZING_TOOLS_SERVER_ADDRESS"
	optionEnableAll                string = "enable-all"
	optionEnableSenzingChatApi     string = "enable-senzing-chat-api"
	optionServerAddress            string = "server-address"
	Short                          string = "serve-chat short description"
	Use                            string = "serve-chat"
	Long                           string = `
 serve-chat long description.
	 `
)

var (
	defaultEngineModuleName string = fmt.Sprintf("serve-chat-%d", time.Now().Unix())
)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	RootCmd.Flags().Bool(option.EnableSwaggerUi, defaultEnableSwaggerUI, fmt.Sprintf("Enable the Swagger UI service [%s]", envar.EnableSwaggerUi))
	RootCmd.Flags().Bool(optionEnableAll, defaultEnableSwaggerUI, fmt.Sprintf("Enable all services [%s]", envarEnableAll))
	RootCmd.Flags().Bool(optionEnableSenzingChatApi, defaultEnableSwaggerUI, fmt.Sprintf("Enable the Senzing Chat API service [%s]", envarEnableSenzingChatApi))
	RootCmd.Flags().Int(option.EngineLogLevel, defaultEngineLogLevel, fmt.Sprintf("Log level for Senzing Engine [%s]", envar.EngineLogLevel))
	RootCmd.Flags().Int(option.HttpPort, defaultHttpPort, fmt.Sprintf("Port to serve HTTP [%s]", envar.HttpPort))
	RootCmd.Flags().String(option.Configuration, defaultConfiguration, fmt.Sprintf("Path to configuration file [%s]", envar.Configuration))
	RootCmd.Flags().String(option.DatabaseUrl, defaultDatabaseUrl, fmt.Sprintf("URL of database to initialize [%s]", envar.DatabaseUrl))
	RootCmd.Flags().String(option.EngineConfigurationJson, defaultEngineConfigurationJson, fmt.Sprintf("JSON string sent to Senzing's init() function [%s]", envar.EngineConfigurationJson))
	RootCmd.Flags().String(option.EngineModuleName, defaultEngineModuleName, fmt.Sprintf("Identifier given to the Senzing engine [%s]", envar.EngineModuleName))
	RootCmd.Flags().String(option.GrpcUrl, defaultGrpcUrl, fmt.Sprintf("URL of Senzing gRPC service [%s]", envar.GrpcUrl))
	RootCmd.Flags().String(option.LogLevel, defaultLogLevel, fmt.Sprintf("Log level [%s]", envar.LogLevel))
	RootCmd.Flags().String(option.ObserverOrigin, defaultObserverOrigin, fmt.Sprintf("Identify this instance to the Observer [%s]", envar.ObserverOrigin))
	RootCmd.Flags().String(option.ObserverUrl, defaultObserverUrl, fmt.Sprintf("URL of Observer [%s]", envar.ObserverUrl))
	RootCmd.Flags().String(optionServerAddress, defaultServerAddress, fmt.Sprintf("IP interface server listens on [%s]", envarServerAddress))
}

// If a configuration file is present, load it.
func loadConfigurationFile(cobraCommand *cobra.Command) {
	configuration := ""
	configFlag := cobraCommand.Flags().Lookup(option.Configuration)
	if configFlag != nil {
		configuration = configFlag.Value.String()
	}
	if configuration != "" { // Use configuration file specified as a command line option.
		viper.SetConfigFile(configuration)
	} else { // Search for a configuration file.

		// Determine home directory.

		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Specify configuration file name.

		viper.SetConfigName("serve-chat")
		viper.SetConfigType("yaml")

		// Define search path order.

		viper.AddConfigPath(home + "/.serve-chat")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/serve-chat")
	}

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Applying configuration file:", viper.ConfigFileUsed())
	}
}

// Configure Viper with user-specified options.
func loadOptions(cobraCommand *cobra.Command) {
	var err error = nil
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix(constant.SetEnvPrefix)

	// Bools

	boolOptions := map[string]bool{
		option.EnableSwaggerUi:     defaultEnableSwaggerUI,
		optionEnableAll:            defaultEnableAll,
		optionEnableSenzingChatApi: defaultEnableSenzingChatApi,
	}
	for optionKey, optionValue := range boolOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// Ints

	intOptions := map[string]int{
		option.EngineLogLevel: defaultEngineLogLevel,
		option.HttpPort:       defaultHttpPort,
	}
	for optionKey, optionValue := range intOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// Strings

	stringOptions := map[string]string{
		option.Configuration:           defaultConfiguration,
		option.DatabaseUrl:             defaultDatabaseUrl,
		option.EngineConfigurationJson: defaultEngineConfigurationJson,
		option.EngineModuleName:        defaultEngineModuleName,
		option.GrpcUrl:                 defaultGrpcUrl,
		option.LogLevel:                defaultLogLevel,
		option.ObserverOrigin:          defaultObserverOrigin,
		option.ObserverUrl:             defaultObserverUrl,
		optionServerAddress:            defaultServerAddress,
	}
	for optionKey, optionValue := range stringOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

	// StringSlice

	stringSliceOptions := map[string][]string{}
	for optionKey, optionValue := range stringSliceOptions {
		viper.SetDefault(optionKey, optionValue)
		err = viper.BindPFlag(optionKey, cobraCommand.Flags().Lookup(optionKey))
		if err != nil {
			panic(err)
		}
	}

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
	loadConfigurationFile(cobraCommand)
	loadOptions(cobraCommand)
	cobraCommand.SetVersionTemplate(constant.VersionTemplate)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	var err error = nil
	ctx := context.TODO()

	// Build senzingEngineConfigurationJson.

	senzingEngineConfigurationJson := viper.GetString(option.EngineConfigurationJson)
	if len(senzingEngineConfigurationJson) == 0 {
		senzingEngineConfigurationJson, err = g2engineconfigurationjson.BuildSimpleSystemConfigurationJson(viper.GetString(option.DatabaseUrl))
		if err != nil {
			return err
		}
	}

	// Determine if gRPC is being used.

	grpcUrl := viper.GetString(option.GrpcUrl)
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
		EnableAll:                      viper.GetBool(optionEnableAll),
		EnableSenzingChatAPI:           viper.GetBool(optionEnableSenzingChatApi),
		EnableSwaggerUI:                viper.GetBool(option.EnableSwaggerUi),
		GrpcDialOptions:                grpcDialOptions,
		GrpcTarget:                     grpcTarget,
		LogLevelName:                   viper.GetString(option.LogLevel),
		ObserverOrigin:                 viper.GetString(option.ObserverOrigin),
		Observers:                      observers,
		OpenApiSpecification:           senzingchatservice.OpenApiSpecificationJson,
		ReadHeaderTimeout:              60 * time.Second,
		SenzingEngineConfigurationJson: senzingEngineConfigurationJson,
		SenzingModuleName:              viper.GetString(option.EngineModuleName),
		SenzingVerboseLogging:          viper.GetInt(option.EngineLogLevel),
		ServerAddress:                  viper.GetString(optionServerAddress),
		ServerPort:                     viper.GetInt(option.HttpPort),
		SwaggerUrlRoutePrefix:          "swagger",
	}
	err = httpServer.Serve(ctx)
	return err
}

// Used in construction of cobra.Command
func Version() string {
	return helper.MakeVersion(githubVersion, githubIteration)
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
