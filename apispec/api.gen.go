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

	// LockFilePath Identifies the file path of the lock file used to retrieve the dependency this event is about
	LockFilePath *string `json:"lock_file_path,omitempty"`

	// Verdict The verdict of the event
	Verdict models.Verdict `json:"verdict"`
}

// DetectionEvent defines model for DetectionEvent.
type DetectionEvent struct {
	Data struct {
		Body interface{} `json:"body"`
		Head interface{} `json:"head"`

		// UniqueId Unique SHA256 identifier
		UniqueId string `json:"unique_id"`
	} `json:"data"`
	GithubContext GitHubDetectionEventContext `json:"github_context"`
	Type          string                      `json:"type"`
}

// Error defines model for Error.
type Error struct {
	Identifier *string `json:"identifier,omitempty"`
	Message    string  `json:"message"`
}

// GitHubDependencyEventContext defines model for GitHubDependencyEventContext.
type GitHubDependencyEventContext = GitHubEventContext

// GitHubDetectionEventContext defines model for GitHubDetectionEventContext.
type GitHubDetectionEventContext = GitHubEventContext

// GitHubEventContext defines model for GitHubEventContext.
type GitHubEventContext struct {
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

// GitHubPipelineEventContext defines model for GitHubPipelineEventContext.
type GitHubPipelineEventContext = GitHubEventContext

// PipelineEvent defines model for PipelineEvent.
type PipelineEvent struct {
	Data          interface{}                `json:"data"`
	GithubContext GitHubPipelineEventContext `json:"github_context"`
	Type          string                     `json:"type"`
}

// PostApiV1DependenciesEventJSONRequestBody defines body for PostApiV1DependenciesEvent for application/json ContentType.
type PostApiV1DependenciesEventJSONRequestBody = DependencyEvent

// PostApiV1DetectionsEventJSONRequestBody defines body for PostApiV1DetectionsEvent for application/json ContentType.
type PostApiV1DetectionsEventJSONRequestBody = DetectionEvent

// PostApiV1PipelineEventJSONRequestBody defines body for PostApiV1PipelineEvent for application/json ContentType.
type PostApiV1PipelineEventJSONRequestBody = PipelineEvent

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

	// PostApiV1DetectionsEventWithBody request with any body
	PostApiV1DetectionsEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostApiV1DetectionsEvent(ctx context.Context, body PostApiV1DetectionsEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostApiV1PipelineEventWithBody request with any body
	PostApiV1PipelineEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostApiV1PipelineEvent(ctx context.Context, body PostApiV1PipelineEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
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

func (c *Client) PostApiV1DetectionsEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1DetectionsEventRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1DetectionsEvent(ctx context.Context, body PostApiV1DetectionsEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1DetectionsEventRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1PipelineEventWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1PipelineEventRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostApiV1PipelineEvent(ctx context.Context, body PostApiV1PipelineEventJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostApiV1PipelineEventRequest(c.Server, body)
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

// NewPostApiV1DetectionsEventRequest calls the generic PostApiV1DetectionsEvent builder with application/json body
func NewPostApiV1DetectionsEventRequest(server string, body PostApiV1DetectionsEventJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostApiV1DetectionsEventRequestWithBody(server, "application/json", bodyReader)
}

// NewPostApiV1DetectionsEventRequestWithBody generates requests for PostApiV1DetectionsEvent with any type of body
func NewPostApiV1DetectionsEventRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/detections/event")
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

// NewPostApiV1PipelineEventRequest calls the generic PostApiV1PipelineEvent builder with application/json body
func NewPostApiV1PipelineEventRequest(server string, body PostApiV1PipelineEventJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostApiV1PipelineEventRequestWithBody(server, "application/json", bodyReader)
}

// NewPostApiV1PipelineEventRequestWithBody generates requests for PostApiV1PipelineEvent with any type of body
func NewPostApiV1PipelineEventRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/pipeline/event")
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

	// PostApiV1DetectionsEventWithBodyWithResponse request with any body
	PostApiV1DetectionsEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1DetectionsEventResponse, error)

	PostApiV1DetectionsEventWithResponse(ctx context.Context, body PostApiV1DetectionsEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1DetectionsEventResponse, error)

	// PostApiV1PipelineEventWithBodyWithResponse request with any body
	PostApiV1PipelineEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1PipelineEventResponse, error)

	PostApiV1PipelineEventWithResponse(ctx context.Context, body PostApiV1PipelineEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1PipelineEventResponse, error)
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

type PostApiV1DetectionsEventResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostApiV1DetectionsEventResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostApiV1DetectionsEventResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostApiV1PipelineEventResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostApiV1PipelineEventResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostApiV1PipelineEventResponse) StatusCode() int {
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

// PostApiV1DetectionsEventWithBodyWithResponse request with arbitrary body returning *PostApiV1DetectionsEventResponse
func (c *ClientWithResponses) PostApiV1DetectionsEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1DetectionsEventResponse, error) {
	rsp, err := c.PostApiV1DetectionsEventWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1DetectionsEventResponse(rsp)
}

func (c *ClientWithResponses) PostApiV1DetectionsEventWithResponse(ctx context.Context, body PostApiV1DetectionsEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1DetectionsEventResponse, error) {
	rsp, err := c.PostApiV1DetectionsEvent(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1DetectionsEventResponse(rsp)
}

// PostApiV1PipelineEventWithBodyWithResponse request with arbitrary body returning *PostApiV1PipelineEventResponse
func (c *ClientWithResponses) PostApiV1PipelineEventWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostApiV1PipelineEventResponse, error) {
	rsp, err := c.PostApiV1PipelineEventWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1PipelineEventResponse(rsp)
}

func (c *ClientWithResponses) PostApiV1PipelineEventWithResponse(ctx context.Context, body PostApiV1PipelineEventJSONRequestBody, reqEditors ...RequestEditorFn) (*PostApiV1PipelineEventResponse, error) {
	rsp, err := c.PostApiV1PipelineEvent(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostApiV1PipelineEventResponse(rsp)
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

// ParsePostApiV1DetectionsEventResponse parses an HTTP response from a PostApiV1DetectionsEventWithResponse call
func ParsePostApiV1DetectionsEventResponse(rsp *http.Response) (*PostApiV1DetectionsEventResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostApiV1DetectionsEventResponse{
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

// ParsePostApiV1PipelineEventResponse parses an HTTP response from a PostApiV1PipelineEventWithResponse call
func ParsePostApiV1PipelineEventResponse(rsp *http.Response) (*PostApiV1PipelineEventResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostApiV1PipelineEventResponse{
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
	// Create a new detection event
	// (POST /api/v1/detections/event)
	PostApiV1DetectionsEvent(c *gin.Context)
	// Create a new pipeline event
	// (POST /api/v1/pipeline/event)
	PostApiV1PipelineEvent(c *gin.Context)
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

// PostApiV1DetectionsEvent operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1DetectionsEvent(c *gin.Context) {

	c.Set(JWTScopes, []string{"write:detectionevents", "project_id"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostApiV1DetectionsEvent(c)
}

// PostApiV1PipelineEvent operation middleware
func (siw *ServerInterfaceWrapper) PostApiV1PipelineEvent(c *gin.Context) {

	c.Set(JWTScopes, []string{"write:pipelineevents", "project_id"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.PostApiV1PipelineEvent(c)
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
	router.POST(options.BaseURL+"/api/v1/detections/event", wrapper.PostApiV1DetectionsEvent)
	router.POST(options.BaseURL+"/api/v1/pipeline/event", wrapper.PostApiV1PipelineEvent)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
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
