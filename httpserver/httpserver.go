package httpserver

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/senzing/go-observing/observer"
	"github.com/senzing/serve-chat/senzingchatapi"
	"github.com/senzing/serve-chat/senzingchatservice"
	"google.golang.org/grpc"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// HttpServerImpl is the default implementation of the HttpServer interface.
type HttpServerImpl struct {
	ChatUrlRoutePrefix             string // FIXME: Only works with "chat"
	EnableAll                      bool
	EnableSenzingChatAPI           bool
	EnableSwaggerUI                bool
	GrpcDialOptions                []grpc.DialOption
	GrpcTarget                     string
	LogLevelName                   string
	ObserverOrigin                 string
	Observers                      []observer.Observer
	OpenApiSpecification           []byte
	ReadHeaderTimeout              time.Duration
	SenzingEngineConfigurationJson string
	SenzingModuleName              string
	SenzingVerboseLogging          int
	ServerAddress                  string
	ServerOptions                  []senzingchatapi.ServerOption
	ServerPort                     int
	SwaggerUrlRoutePrefix          string // FIXME: Only works with "swagger"
}

type TemplateVariables struct {
	HttpServerImpl
	ChatServerStatus string
	ChatServerUrl    string
	HtmlTitle        string
	RequestHost      string
	SwaggerStatus    string
	SwaggerUrl       string
}

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

//go:embed static/*
var static embed.FS

// ----------------------------------------------------------------------------
// Internal methods
// ----------------------------------------------------------------------------

func (httpServer *HttpServerImpl) getServerStatus(up bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}
	if up {
		result = "green"
	}
	return result
}

func (httpServer *HttpServerImpl) getServerUrl(up bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}
	if up {
		result = url
	}
	return result
}

func (httpServer *HttpServerImpl) openApiFunc(ctx context.Context, openApiSpecification []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bytesBuffer bytes.Buffer
		bufioWriter := bufio.NewWriter(&bytesBuffer)
		openApiSpecificationTemplate, err := template.New("OpenApiTemplate").Parse(string(httpServer.OpenApiSpecification))
		if err != nil {
			panic(err)
		}
		templateVariables := TemplateVariables{
			RequestHost: string(r.Host),
		}
		err = openApiSpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
func (httpServer *HttpServerImpl) populateStaticTemplate(responseWriter http.ResponseWriter, request *http.Request, filepath string, templateVariables TemplateVariables) {
	templateBytes, err := static.ReadFile(filepath)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	templateParsed, err := template.New("HtmlTemplate").Parse(string(templateBytes))
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
	err = templateParsed.Execute(responseWriter, templateVariables)
	if err != nil {
		http.Error(responseWriter, http.StatusText(500), 500)
		return
	}
}

// --- http.ServeMux ----------------------------------------------------------

func (httpServer *HttpServerImpl) getSenzingChatMux(ctx context.Context) *senzingchatapi.Server {
	service := &senzingchatservice.ChatApiServiceImpl{
		GrpcDialOptions:                httpServer.GrpcDialOptions,
		GrpcTarget:                     httpServer.GrpcTarget,
		LogLevelName:                   httpServer.LogLevelName,
		ObserverOrigin:                 httpServer.ObserverOrigin,
		Observers:                      httpServer.Observers,
		SenzingEngineConfigurationJson: httpServer.SenzingEngineConfigurationJson,
		SenzingModuleName:              httpServer.SenzingModuleName,
		SenzingVerboseLogging:          httpServer.SenzingVerboseLogging,
		UrlRoutePrefix:                 httpServer.ChatUrlRoutePrefix,
		OpenApiSpecificationSpec:       httpServer.OpenApiSpecification,
	}
	srv, err := senzingchatapi.NewServer(service, httpServer.ServerOptions...)
	if err != nil {
		log.Fatal(err)
	}
	return srv
}

func (httpServer *HttpServerImpl) getSwaggerUiMux(ctx context.Context) *http.ServeMux {
	swaggerMux := swaggerui.Handler([]byte{}) // OpenAPI specification handled by openApiFunc()
	swaggerFunc := swaggerMux.ServeHTTP
	submux := http.NewServeMux()
	submux.HandleFunc("/", swaggerFunc)
	submux.HandleFunc("/swagger_spec", httpServer.openApiFunc(ctx, httpServer.OpenApiSpecification))
	return submux
}

// --- Http Funcs -------------------------------------------------------------

func (httpServer *HttpServerImpl) siteFunc(w http.ResponseWriter, r *http.Request) {
	templateVariables := TemplateVariables{
		HttpServerImpl:   *httpServer,
		HtmlTitle:        "serve-chat",
		ChatServerUrl:    httpServer.getServerUrl(httpServer.EnableSenzingChatAPI, fmt.Sprintf("http://%s/chat", r.Host)),
		ChatServerStatus: httpServer.getServerStatus(httpServer.EnableSenzingChatAPI),
		SwaggerUrl:       httpServer.getServerUrl(httpServer.EnableSwaggerUI, fmt.Sprintf("http://%s/swagger", r.Host)),
		SwaggerStatus:    httpServer.getServerStatus(httpServer.EnableSwaggerUI),
	}
	w.Header().Set("Content-Type", "text/html")
	filePath := fmt.Sprintf("static/templates%s", r.RequestURI)
	httpServer.populateStaticTemplate(w, r, filePath, templateVariables)
}

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

func (httpServer *HttpServerImpl) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()
	var userMessage string = ""

	// Enable Senzing HTTP Chat API.

	if httpServer.EnableAll || httpServer.EnableSenzingChatAPI {
		senzingApiMux := httpServer.getSenzingChatMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.ChatUrlRoutePrefix), http.StripPrefix("/chat", senzingApiMux))
		userMessage = fmt.Sprintf("%sServing Senzing Chat API at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.ChatUrlRoutePrefix)
	}

	// Enable SwaggerUI.

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUiMux := httpServer.getSwaggerUiMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.SwaggerUrlRoutePrefix), http.StripPrefix("/swagger", swaggerUiMux))
		userMessage = fmt.Sprintf("%sServing SwaggerUI at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.SwaggerUrlRoutePrefix)
	}

	// Add route to template pages.

	rootMux.HandleFunc("/site/", httpServer.siteFunc)
	userMessage = fmt.Sprintf("%sServing Console at          http://localhost:%d\n", userMessage, httpServer.ServerPort)

	// Add route to static files.

	rootDir, err := fs.Sub(static, "static/root")
	if err != nil {
		panic(err)
	}
	rootMux.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(rootDir))))

	// Start service.

	listenOnAddress := fmt.Sprintf("%s:%v", httpServer.ServerAddress, httpServer.ServerPort)
	userMessage = fmt.Sprintf("%sStarting server on interface:port '%s'...\n", userMessage, listenOnAddress)
	fmt.Println(userMessage)
	server := http.Server{
		ReadHeaderTimeout: httpServer.ReadHeaderTimeout,
		Addr:              listenOnAddress,
		Handler:           rootMux,
	}
	return server.ListenAndServe()
}
