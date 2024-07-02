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
	ChatURLRoutePrefix    string // FIXME: Only works with "chat"
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
	SwaggerURLRoutePrefix string // FIXME: Only works with "swagger"
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
// Internal methods
// ----------------------------------------------------------------------------

func (httpServer *BasicHTTPServer) getServerStatus(up bool) string {
	result := "red"
	if httpServer.EnableAll {
		result = "green"
	}
	if up {
		result = "green"
	}
	return result
}

func (httpServer *BasicHTTPServer) getServerURL(up bool, url string) string {
	result := ""
	if httpServer.EnableAll {
		result = url
	}
	if up {
		result = url
	}
	return result
}

func (httpServer *BasicHTTPServer) openAPIFunc(ctx context.Context, openAPISpecification []byte) http.HandlerFunc {
	_ = ctx
	_ = openAPISpecification
	return func(w http.ResponseWriter, r *http.Request) {
		var bytesBuffer bytes.Buffer
		bufioWriter := bufio.NewWriter(&bytesBuffer)
		openAPISpecificationTemplate, err := template.New("OpenApiTemplate").Parse(string(httpServer.OpenAPISpecification))
		if err != nil {
			panic(err)
		}
		templateVariables := TemplateVariables{
			RequestHost: string(r.Host),
		}
		err = openAPISpecificationTemplate.Execute(bufioWriter, templateVariables)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytesBuffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}
func (httpServer *BasicHTTPServer) populateStaticTemplate(responseWriter http.ResponseWriter, request *http.Request, filepath string, templateVariables TemplateVariables) {
	_ = request
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
		log.Fatal(err)
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

func (httpServer *BasicHTTPServer) siteFunc(w http.ResponseWriter, r *http.Request) {
	templateVariables := TemplateVariables{
		BasicHTTPServer:  *httpServer,
		HTMLTitle:        "serve-chat",
		ChatServerURL:    httpServer.getServerURL(httpServer.EnableSenzingChatAPI, fmt.Sprintf("http://%s/chat", r.Host)),
		ChatServerStatus: httpServer.getServerStatus(httpServer.EnableSenzingChatAPI),
		SwaggerURL:       httpServer.getServerURL(httpServer.EnableSwaggerUI, fmt.Sprintf("http://%s/swagger", r.Host)),
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

func (httpServer *BasicHTTPServer) Serve(ctx context.Context) error {
	rootMux := http.NewServeMux()
	var userMessage = ""

	// Enable Senzing HTTP Chat API.

	if httpServer.EnableAll || httpServer.EnableSenzingChatAPI {
		senzingAPIMux := httpServer.getSenzingChatMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.ChatURLRoutePrefix), http.StripPrefix("/chat", senzingAPIMux))
		userMessage = fmt.Sprintf("%sServing Senzing Chat API at http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.ChatURLRoutePrefix)
	}

	// Enable SwaggerUI.

	if httpServer.EnableAll || httpServer.EnableSwaggerUI {
		swaggerUIMux := httpServer.getSwaggerUIMux(ctx)
		rootMux.Handle(fmt.Sprintf("/%s/", httpServer.SwaggerURLRoutePrefix), http.StripPrefix("/swagger", swaggerUIMux))
		userMessage = fmt.Sprintf("%sServing SwaggerUI at        http://localhost:%d/%s\n", userMessage, httpServer.ServerPort, httpServer.SwaggerURLRoutePrefix)
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

	if !httpServer.AvoidServing {
		return server.ListenAndServe()
	}
	return nil
}
