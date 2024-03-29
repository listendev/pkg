// Package apispec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package apispec

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/listendev/pkg/models"
	"github.com/listendev/pkg/type/int64string"
)

const (
	JWTScopes = "JWT.Scopes"
)

// DependencyEvent defines model for DependencyEvent.
type DependencyEvent struct {
	GithubContext GitHubDependencyEventContext `json:"github_context"`

	// Verdict The verdict of the event
	Verdict models.Verdict `json:"verdict"`
}

// Error defines model for Error.
type Error struct {
	Identifier *string `json:"identifier,omitempty"`
	Message    string  `json:"message"`
}

// GitHubDependencyEventContext defines model for GitHubDependencyEventContext.
type GitHubDependencyEventContext struct {
	Action            string                  `json:"action"`
	ActionPath        *string                 `json:"action_path,omitempty"`
	ActionRepository  *string                 `json:"action_repository,omitempty"`
	Actor             string                  `json:"actor"`
	ActorId           int64string.Int64String `json:"actor_id"`
	EventName         string                  `json:"event_name"`
	Job               string                  `json:"job"`
	Ref               string                  `json:"ref"`
	RefName           string                  `json:"ref_name"`
	RefProtected      bool                    `json:"ref_protected"`
	RefType           string                  `json:"ref_type"`
	Repository        string                  `json:"repository"`
	RepositoryId      int64string.Int64String `json:"repository_id"`
	RepositoryOwner   string                  `json:"repository_owner"`
	RepositoryOwnerId int64string.Int64String `json:"repository_owner_id"`
	RunAttempt        int64string.Int64String `json:"run_attempt"`
	RunId             int64string.Int64String `json:"run_id"`
	RunNumber         int64string.Int64String `json:"run_number"`
	RunnerArch        string                  `json:"runner_arch"`
	RunnerDebug       *bool                   `json:"runner_debug,omitempty"`
	RunnerOs          string                  `json:"runner_os"`
	ServerUrl         string                  `json:"server_url"`
	Sha               string                  `json:"sha"`
	TriggeringActor   string                  `json:"triggering_actor"`
	Workflow          string                  `json:"workflow"`
	WorkflowRef       string                  `json:"workflow_ref"`
	Workspace         string                  `json:"workspace"`
}

// PostApiV1DependenciesEventJSONRequestBody defines body for PostApiV1DependenciesEvent for application/json ContentType.
type PostApiV1DependenciesEventJSONRequestBody = DependencyEvent

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostApiV1DependenciesEventWithBody request with any body
	PostApiV1DependenciesEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostApiV1DependenciesEvent(ctx context.Context, body PostApiV1DependenciesEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostApiV1DependenciesEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1DependenciesEventRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1DependenciesEvent(ctx context.Context, body PostApiV1DependenciesEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1DependenciesEventRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostApiV1DependenciesEventRequest calls the generic PostApiV1DependenciesEvent builder with application/json body
func NewPostApiV1DependenciesEventRequest(server string, body PostApiV1DependenciesEventJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostApiV1DependenciesEventRequestWithBody(server, "application/json", bodyReader)
}

// NewPostApiV1DependenciesEventRequestWithBody generates requests for PostApiV1DependenciesEvent with any type of body
func NewPostApiV1DependenciesEventRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/dependencies/event")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostApiV1DependenciesEventWithBodyWithResponse request with any body
	PostApiV1DependenciesEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1DependenciesEventResponse, error)

	PostApiV1DependenciesEventWithResponse(ctx context.Context, body PostApiV1DependenciesEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1DependenciesEventResponse, error)
}

type PostApiV1DependenciesEventResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostApiV1DependenciesEventResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostApiV1DependenciesEventResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostApiV1DependenciesEventWithBodyWithResponse request with arbitrary body returning *PostApiV1DependenciesEventResponse
func (c *ClientWithResponses) PostApiV1DependenciesEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1DependenciesEventResponse, error) {
	rsp, err := c.PostApiV1DependenciesEventWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1DependenciesEventResponse(rsp)
}

func (c *ClientWithResponses) PostApiV1DependenciesEventWithResponse(ctx context.Context, body PostApiV1DependenciesEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1DependenciesEventResponse, error) {
	rsp, err := c.PostApiV1DependenciesEvent(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1DependenciesEventResponse(rsp)
}

// ParsePostApiV1DependenciesEventResponse parses an HTTP response from a PostApiV1DependenciesEventWithResponse call
func ParsePostApiV1DependenciesEventResponse(rsp *http.Response) (*PostApiV1DependenciesEventResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostApiV1DependenciesEventResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new dependency event
	// (POST /api/v1/dependencies/event)
	PostApiV1DependenciesEvent(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// PostApiV1DependenciesEvent operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1DependenciesEvent(c *gin.Context) {

	c.Set(JWTScopes, []string{"write:dependencyevents", "project_id"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostApiV1DependenciesEvent(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/api/v1/dependencies/event", wrapper.PostApiV1DependenciesEvent)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RXXW/bNhT9KwK3R1v0Z+zoqV0brBn2UKBdO6AIBIq6sphKJEdSdrzA/30gKTuiLKfd",
	"0r1Z4j2HV/fz+BFRUUvBgRuNkkekaQk1cT/fggSeA6f7my1wY19JJSQow8AZbJgpmyylght4cOc/KyhQ",
	"gn7CT6S4ZcS/MvOuyXqkb1rsYYS2oHJGHU0OmiomDRMcJehjCVF7GIkiMiVE4BwaIbOXgBIksnug9vlh",
	"vBHj9mUtcqh0/KmlPRxGSMFfDVOQo+TL6bpR/zPu+rSHEbpRSijrGTyQWlZgf9agNdlYuz84aUwpFPsb",
	"cmsdRonlwA0rGDiCllobxfjGGp9ozs56Dh8Nh/x7NrZneSPUR7bzOe07myv6VTSd2D556k1SSUwZYnEp",
	"asCq4RwU3gn1FactHVYgxVjs7AEnNYxFMfZHY3uCt9NnLrIWmhmh9uF1z3NeIOylDwlqBCXmonXK8hAw",
	"nc0Xy6vVGSAsOsbN1cKfxLf294cTravZ1HocEtuAFZXYpTnTkhhaDvl0L7IQlTWsyocsXQeG4So0LoHk",
	"Gtf7NFOED1+hoBjwbjqb4xrUBi5BpBIGqIEwXEY1cAJkQlRA+BHh3z4O0Q0nvM0VfgdVJcafhbr05Uf4",
	"peStr/97+jrsrva+u5z6wB9fWarhKTEGaml6zC/jPHP0arlez6aL6/nLeHlTZ/0AvozShpUo2ptLf66v",
	"BhPi7XPImk0AKEilh6vWI4QO+X9nvHkYukGD2oJKG1WFgNIYqROM/caJqagH0SUJYUVB6HK+giuaZcX1",
	"fJFN1qvFksxX6+v5bDWb5cVsvSLL+RCZUWyzAfuU/rspeBxLIeCDqMGUPj0XIenZDDp2cOk6eGc7GMc+",
	"CviIsuNpfHyI93X16nsnlwVpSSh8ay3Ve7cixnbMBQ/nrL3l267M4yrpLIlgsPtB7YdwZ6L2J2VnDgZT",
	"rz/DBqbO8DwJZ8Cpe4N2CxulW9RBwfr6GyicTkn0Ut3NwLk6saFkvBA2OVZjEa/xoCasQgkCvmEc3EWv",
	"NkRxMDFhaIT8KnLi7+bJJPoIpLY5DyVixbQBHuewjZz00dHr97f2I5jx6kZayi0o7e2n8SSeWBohgdvD",
	"BM3jSWw7yMob1+eYSIa3U5wfVRUDjeGkhIUelKpMR8BzKRg3EalsWUeFUE60UgXE2lkRSyIOOy9jY+Tc",
	"UO7sNkcJei+0eS3Zp+nbztU3reS1VQna/CLy/TGirUtEyopRR4PvtRd4Xnx/S5r3lf4hLH+7zN0LLQXX",
	"XkPOJrPzz3do/52QR7qhFLQumqra21gvJtMf5rEX5M7P0IW+Fl9OJv//pbfcgOKkinwfRdAa2saijWJm",
	"j5Ivj+i3zx/tLNkpZiA5ldXeVYHtQqmEbRnbuXeHuxHSTV0TK4bQGxfStmqekO3/oMPhtHO0u8gtntO2",
	"oUJB/NQh6HB3+CcAAP//2T/kNvYNAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
