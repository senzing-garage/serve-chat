package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-observing/observer"
	"github.com/senzing-garage/serve-chat/senzingchatapi"
	"github.com/senzing-garage/serve-chat/senzingchatservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// BasicHTTPServer is the default implementation of the HttpServer interface.
type BasicHTTPServer struct {
	AvoidServing          bool
	ChatURLRoutePrefix    string // IMPROVE: Only works with "chat"
	EnableAll             bool
	EnableSenzingChatAPI  bool
	EnableSwaggerUI       bool
	GrpcDialOptions       []grpc.DialOption
	GrpcTarget            string
	LogLevelName          string
	ObserverOrigin        string
	Observers             []observer.Observer
	OpenAPISpecification  []byte
	ReadHeaderTimeout     time.Duration
	Setting               string
	SenzingInstanceName   string
	SenzingVerboseLogging int64
	ServerAddress         string
	ServerOptions         []senzingchatapi.ServerOption
	ServerPort            int
	SwaggerURLRoutePrefix string // IMPROVE: Only works with "swagger"
}

type TemplateVariables struct {
	BasicHTTPServer
	ChatServerStatus string
	ChatServerURL    string
	HTMLTitle        string
	RequestHost      string
	SwaggerStatus    string
	SwaggerURL       string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

/*
The Serve method simply prints the 'Something' value in the type-struct.

Input
  - ctx: A context to control lifecycle.

Output
  - Nothing is returned, except for an error.  However, something is printed.
    See the example output.
*/

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()

	var userMessages []string

	// Add to root Mux.

	userMessages = append(userMessages, httpServer.addChatToMux(ctx, rootMux)...)
	userMessages = append(userMessages, httpServer.addSwagerToMux(ctx, rootMux)...)

	// Add route to template pages.

	rootMux.HandleFunc("/site/", httpServer.siteFunc)
	userMessages = append(
		userMessages,
		fmt.Sprintf("Serving Console at          http://localhost:%d\n", httpServer.ServerPort),
	)

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}

	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessages = append(userMessages,
		fmt.Sprintf("Starting server on interface:port '%s'...", listenOnAddress))

	for userMessage := range userMessages {
		outputln(userMessage)
	}

	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}

	if !httpServer.AvoidServing {
		err = server.ListenAndServe()

		return wraperror.Errorf(err, "ListenAndServe")
	}

	return nil
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) addChatToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	if httpServer.EnableAll || httpServer.EnableSenzingChatAPI {
		senzingAPIMux := httpServer.getSenzingChatMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.ChatURLRoutePrefix), http.StripPrefix("/chat", senzingAPIMux))
		result = append(result,
			fmt.Sprintf(
				"Serving Senzing Chat API at http://localhost:%d/%s",
				httpServer.ServerPort,
				httpServer.ChatURLRoutePrefix))
	}

	return result
}

func (httpServer *BasicHTTPServer) addSwagerToMux(
	ctx context.Context,
	rootMux *http.ServeMux,
) []string {
	var result []string

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUIMux := httpServer.getSwaggerUIMux(ctx)
		rootMux.Handle(
			fmt.Sprintf("/%s/", httpServer.SwaggerURLRoutePrefix),
			http.StripPrefix("/swagger", swaggerUIMux),
		)

		result = append(result,
			fmt.Sprintf(
				"Serving SwaggerUI at        http://localhost:%d/%s",
				httpServer.ServerPort,
				httpServer.SwaggerURLRoutePrefix))
	}

	return result
}

func (httpServer *BasicHTTPServer) getServerStatus(active bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}

	if active {
		result = "green"
	}

	return result
}

func (httpServer *BasicHTTPServer) getServerURL(active bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}

	if active {
		result = url
	}

	return result
}

func (httpServer *BasicHTTPServer) openAPIFunc(ctx context.Context, openAPISpecification []byte) http.HandlerFunc {
	_ = ctx
	_ = openAPISpecification

	return func(writer http.ResponseWriter, request *http.Request) {
		var bytesBuffer bytes.Buffer

		bufioWriter := bufio.NewWriter(&bytesBuffer)

		openAPISpecificationTemplate, err := template.New("OpenApiTemplate").
			Parse(string(httpServer.OpenAPISpecification))
		if err != nil {
			panic(err)
		}

		templateVariables := TemplateVariables{
			RequestHost: request.Host,
		}

		err = openAPISpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}

		_, err = writer.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

func (httpServer *BasicHTTPServer) populateStaticTemplate(
	responseWriter http.ResponseWriter,
	request *http.Request,
	filepath string,
	templateVariables TemplateVariables,
) {
	_ = request

	templateBytes, err := static.ReadFile(filepath)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	err = templateParsed.Execute(responseWriter, templateVariables)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *BasicHTTPServer) getSenzingChatMux(ctx context.Context) *senzingchatapi.Server {
	_ = ctx
	service := &senzingchatservice.BasicChatAPIService{
		GrpcDialOptions:          httpServer.GrpcDialOptions,
		GrpcTarget:               httpServer.GrpcTarget,
		LogLevelName:             httpServer.LogLevelName,
		ObserverOrigin:           httpServer.ObserverOrigin,
		Observers:                httpServer.Observers,
		Settings:                 httpServer.Setting,
		SenzingInstanceName:      httpServer.SenzingInstanceName,
		SenzingVerboseLogging:    httpServer.SenzingVerboseLogging,
		URLRoutePrefix:           httpServer.ChatURLRoutePrefix,
		OpenAPISpecificationSpec: httpServer.OpenAPISpecification,
	}

	srv, err := senzingchatapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		panic(err)
	}

	return srv
}

func (httpServer *BasicHTTPServer) getSwaggerUIMux(ctx context.Context) *http.ServeMux {
	swaggerMux := swaggerui.Handler([]byte{}) // OpenAPI specification handled by openApiFunc()
	swaggerFunc := swaggerMux.ServeHTTP
	submux := http.NewServeMux()
	submux.HandleFunc("/", swaggerFunc)
	submux.HandleFunc("/swagger_spec", httpServer.openAPIFunc(ctx, httpServer.OpenAPISpecification))

	return submux
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *BasicHTTPServer) siteFunc(writer http.ResponseWriter, request *http.Request) {
	templateVariables := TemplateVariables{
		BasicHTTPServer: *httpServer,
		HTMLTitle:       "serve-chat",
		ChatServerURL: httpServer.getServerURL(
			httpServer.EnableSenzingChatAPI,
			fmt.Sprintf("http://%s/chat", request.Host),
		),
		ChatServerStatus: httpServer.getServerStatus(httpServer.EnableSenzingChatAPI),
		SwaggerURL: httpServer.getServerURL(
			httpServer.EnableSwaggerUI,
			fmt.Sprintf("http://%s/swagger", request.Host),
		),
		SwaggerStatus: httpServer.getServerStatus(httpServer.EnableSwaggerUI),
	}

	writer.Header().Set("Content-Type", "text/html")

	filePath := "static/templates" + request.RequestURI
	httpServer.populateStaticTemplate(writer, request, filePath, templateVariables)
}

func outputln(message ...any) {
	fmt.Println(message...) //nolint
}
