// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package api

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
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for RecordOperation.
const (
	RecordOperationCREATE RecordOperation = "CREATE"
	RecordOperationDELETE RecordOperation = "DELETE"
	RecordOperationUPDATE RecordOperation = "UPDATE"
)

// Defines values for RecordWithIDOperation.
const (
	RecordWithIDOperationCREATE RecordWithIDOperation = "CREATE"
	RecordWithIDOperationDELETE RecordWithIDOperation = "DELETE"
	RecordWithIDOperationUPDATE RecordWithIDOperation = "UPDATE"
)

// IoK8sApiAuthenticationV1UserInfo UserInfo holds the information about the user needed to implement the user.Info interface.
type IoK8sApiAuthenticationV1UserInfo struct {
	// Extra Any additional information provided by the authenticator.
	Extra *map[string][]string `json:"extra,omitempty"`

	// Groups The names of groups this user is a part of.
	Groups *[]string `json:"groups,omitempty"`

	// Uid A unique value that identifies this user across time. If this user is deleted and another user by the same name is added, they will have different UIDs.
	Uid *string `json:"uid,omitempty"`

	// Username The name that uniquely identifies this user among all active users.
	Username *string `json:"username,omitempty"`
}

// IoK8sApimachineryPkgApisMetaV1FieldsV1 defines model for io.k8s.apimachinery.pkg.apis.meta.v1.FieldsV1.
type IoK8sApimachineryPkgApisMetaV1FieldsV1 = map[string]interface{}

// IoK8sApimachineryPkgApisMetaV1GroupVersionKind GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type IoK8sApimachineryPkgApisMetaV1GroupVersionKind struct {
	Group   string `json:"group"`
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

// IoK8sApimachineryPkgApisMetaV1GroupVersionResource defines model for io.k8s.apimachinery.pkg.apis.meta.v1.GroupVersionResource.
type IoK8sApimachineryPkgApisMetaV1GroupVersionResource struct {
	Group    string `json:"group"`
	Resource string `json:"resource"`
	Version  string `json:"version"`
}

// IoK8sApimachineryPkgApisMetaV1ManagedFieldsEntry defines model for io.k8s.apimachinery.pkg.apis.meta.v1.ManagedFieldsEntry.
type IoK8sApimachineryPkgApisMetaV1ManagedFieldsEntry struct {
	ApiVersion  *string                                 `json:"apiVersion,omitempty"`
	FieldsType  *string                                 `json:"fieldsType,omitempty"`
	FieldsV1    *IoK8sApimachineryPkgApisMetaV1FieldsV1 `json:"fieldsV1,omitempty"`
	Manager     *string                                 `json:"manager,omitempty"`
	Operation   *string                                 `json:"operation,omitempty"`
	Subresource *string                                 `json:"subresource,omitempty"`
	Time        *IoK8sApimachineryPkgApisMetaV1Time     `json:"time,omitempty"`
}

// IoK8sApimachineryPkgApisMetaV1ObjectMeta defines model for io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta.
type IoK8sApimachineryPkgApisMetaV1ObjectMeta struct {
	Annotations                *map[string]string                                  `json:"annotations,omitempty"`
	CreationTimestamp          *IoK8sApimachineryPkgApisMetaV1Time                 `json:"creationTimestamp,omitempty"`
	DeletionGracePeriodSeconds *int64                                              `json:"deletionGracePeriodSeconds,omitempty"`
	DeletionTimestamp          *IoK8sApimachineryPkgApisMetaV1Time                 `json:"deletionTimestamp,omitempty"`
	Finalizers                 *[]string                                           `json:"finalizers,omitempty"`
	GenerateName               *string                                             `json:"generateName,omitempty"`
	Generation                 *int64                                              `json:"generation,omitempty"`
	Labels                     *map[string]string                                  `json:"labels,omitempty"`
	ManagedFields              *[]IoK8sApimachineryPkgApisMetaV1ManagedFieldsEntry `json:"managedFields,omitempty"`
	Name                       *string                                             `json:"name,omitempty"`
	Namespace                  *string                                             `json:"namespace,omitempty"`
	OwnerReferences            *[]IoK8sApimachineryPkgApisMetaV1OwnerReference     `json:"ownerReferences,omitempty"`
	ResourceVersion            *string                                             `json:"resourceVersion,omitempty"`
	SelfLink                   *string                                             `json:"selfLink,omitempty"`
	Uid                        *string                                             `json:"uid,omitempty"`
}

// IoK8sApimachineryPkgApisMetaV1OwnerReference defines model for io.k8s.apimachinery.pkg.apis.meta.v1.OwnerReference.
type IoK8sApimachineryPkgApisMetaV1OwnerReference struct {
	// ApiVersion API version of the referent.
	ApiVersion         string `json:"apiVersion"`
	BlockOwnerDeletion *bool  `json:"blockOwnerDeletion,omitempty"`
	Controller         *bool  `json:"controller,omitempty"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	Uid                string `json:"uid"`
}

// IoK8sApimachineryPkgApisMetaV1Time defines model for io.k8s.apimachinery.pkg.apis.meta.v1.Time.
type IoK8sApimachineryPkgApisMetaV1Time = time.Time

// Record defines model for record.
type Record struct {
	ChangeTimestamp time.Time `json:"changeTimestamp"`
	Cluster         string    `json:"cluster"`
	DiffString      string    `json:"diffString"`
	Generation      int64     `json:"generation"`
	JsonPatch       string    `json:"jsonPatch"`
	JsonPatch6902   string    `json:"jsonPatch6902"`

	// Kind GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
	Kind       IoK8sApimachineryPkgApisMetaV1GroupVersionKind     `json:"kind"`
	Name       string                                             `json:"name"`
	Namespace  string                                             `json:"namespace"`
	ObjectMeta IoK8sApimachineryPkgApisMetaV1ObjectMeta           `json:"objectMeta"`
	Operation  RecordOperation                                    `json:"operation"`
	Resource   IoK8sApimachineryPkgApisMetaV1GroupVersionResource `json:"resource"`
	Uid        string                                             `json:"uid"`

	// UserInfo UserInfo holds the information about the user needed to implement the user.Info interface.
	UserInfo IoK8sApiAuthenticationV1UserInfo `json:"userInfo"`
}

// RecordOperation defines model for Record.Operation.
type RecordOperation string

// RecordWithID defines model for recordWithID.
type RecordWithID struct {
	ChangeTimestamp time.Time `json:"changeTimestamp"`
	Cluster         string    `json:"cluster"`
	DiffString      string    `json:"diffString"`
	Generation      int64     `json:"generation"`
	Id              string    `json:"id"`
	JsonPatch       string    `json:"jsonPatch"`
	JsonPatch6902   string    `json:"jsonPatch6902"`

	// Kind GroupVersionKind unambiguously identifies a kind.  It doesn't anonymously include GroupVersion to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
	Kind       IoK8sApimachineryPkgApisMetaV1GroupVersionKind     `json:"kind"`
	Name       string                                             `json:"name"`
	Namespace  string                                             `json:"namespace"`
	ObjectMeta IoK8sApimachineryPkgApisMetaV1ObjectMeta           `json:"objectMeta"`
	Operation  RecordWithIDOperation                              `json:"operation"`
	Resource   IoK8sApimachineryPkgApisMetaV1GroupVersionResource `json:"resource"`
	Uid        string                                             `json:"uid"`

	// UserInfo UserInfo holds the information about the user needed to implement the user.Info interface.
	UserInfo IoK8sApiAuthenticationV1UserInfo `json:"userInfo"`
}

// RecordWithIDOperation defines model for RecordWithID.Operation.
type RecordWithIDOperation string

// AddRecordJSONRequestBody defines body for AddRecord for application/json ContentType.
type AddRecordJSONRequestBody = Record

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
	// ListRecords request
	ListRecords(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// AddRecordWithBody request with any body
	AddRecordWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	AddRecord(ctx context.Context, body AddRecordJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListRecords(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListRecordsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddRecordWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddRecordRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) AddRecord(ctx context.Context, body AddRecordJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewAddRecordRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListRecordsRequest generates requests for ListRecords
func NewListRecordsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/record")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAddRecordRequest calls the generic AddRecord builder with application/json body
func NewAddRecordRequest(server string, body AddRecordJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddRecordRequestWithBody(server, "application/json", bodyReader)
}

// NewAddRecordRequestWithBody generates requests for AddRecord with any type of body
func NewAddRecordRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/record")
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
	// ListRecordsWithResponse request
	ListRecordsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListRecordsResponse, error)

	// AddRecordWithBodyWithResponse request with any body
	AddRecordWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddRecordResponse, error)

	AddRecordWithResponse(ctx context.Context, body AddRecordJSONRequestBody, reqEditors ...RequestEditorFn) (*AddRecordResponse, error)
}

type ListRecordsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]RecordWithID
}

// Status returns HTTPResponse.Status
func (r ListRecordsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListRecordsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type AddRecordResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *[]RecordWithID
}

// Status returns HTTPResponse.Status
func (r AddRecordResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r AddRecordResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListRecordsWithResponse request returning *ListRecordsResponse
func (c *ClientWithResponses) ListRecordsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListRecordsResponse, error) {
	rsp, err := c.ListRecords(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListRecordsResponse(rsp)
}

// AddRecordWithBodyWithResponse request with arbitrary body returning *AddRecordResponse
func (c *ClientWithResponses) AddRecordWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*AddRecordResponse, error) {
	rsp, err := c.AddRecordWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddRecordResponse(rsp)
}

func (c *ClientWithResponses) AddRecordWithResponse(ctx context.Context, body AddRecordJSONRequestBody, reqEditors ...RequestEditorFn) (*AddRecordResponse, error) {
	rsp, err := c.AddRecord(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseAddRecordResponse(rsp)
}

// ParseListRecordsResponse parses an HTTP response from a ListRecordsWithResponse call
func ParseListRecordsResponse(rsp *http.Response) (*ListRecordsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListRecordsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []RecordWithID
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseAddRecordResponse parses an HTTP response from a AddRecordWithResponse call
func ParseAddRecordResponse(rsp *http.Response) (*AddRecordResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &AddRecordResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest []RecordWithID
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /record)
	ListRecords(c *gin.Context)

	// (POST /record)
	AddRecord(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ListRecords operation middleware
func (siw *ServerInterfaceWrapper) ListRecords(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{"records:r"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListRecords(c)
}

// AddRecord operation middleware
func (siw *ServerInterfaceWrapper) AddRecord(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{"records:w"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AddRecord(c)
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

	router.GET(options.BaseURL+"/record", wrapper.ListRecords)
	router.POST(options.BaseURL+"/record", wrapper.AddRecord)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RYX2/bOBL/KgTvgHtRFKd3KG79lm7SwrvtNkjd9qHoA02OLNYUqZIjp97C330xlGxL",
	"tuQ6jftgQNaQP87Mb/6J37l0ReksWAx8/J0HmUMh4qN26eL/IRWlTkWFOVjUUqB2Nl1epe8D+InNHC1U",
	"EKTXJYn4mG8kLHdGBYY5MG0z54u4l4mZqzC+rQJ4ZgEUKIaO6aI0UIDdCdMIoy2Cz4SElCe89K4Ejxqi",
	"hvANvaAHoZQmdGHuOgs0QhEfcFUCH/OAXts5XyebF8J7saL/XRuu7YrtMDv6l94tNak8W0VFW65xnlRs",
	"kN3sC0gk6Ll3VRkOHTXNgVlRQGAuY/UihrkOtWN0YIKVwiNzGcGebkql1eFh16yy+msFbClMBQxzgUwr",
	"UjzT0D5XSO9CYKgLSNkk62qkwACCYsLSz2EOvpY1zgiiqG2K6isFKqH3K/agjWG5WAJTOsvAE83vJzeh",
	"5bCdQYRIIMMuq/WvLTKrAUMKZ+dMGMOERL2sQ6rvwHUPZ7vYL4TMtQW/SsvFnF6EtAAUlAMvNRgVPly1",
	"WHkkwiti/QP4oJ39U9se3vZXsMqKYqbnlatC13LBFtqqlLEJMuUg2P8gcWRXRbPUSlMpYG1EyjuxdFpR",
	"GDuKcMmkAy8pyztIVQAmBvbKKqArWCF8yIUx5NT9TI3h3Ru9i8bsA8GyPqZHtk64h6+V9qD4+FODvdvQ",
	"YH5OzkDKPQRXeRlD8VSLfGvPOa3a4v60ZW+EFXNQddzeWvSrQ7tEqT8M6pjwLO6dxteD4jol/u0h42P+",
	"r8tdh7ls2svl4/JrnfAiqu57DyX1BQ5pHKrZUUao0p1F3SkB/Xw1eRuXvwEUPaRY6zCaGI61u6HWsNNE",
	"eogwpGtAUZRntJyaqAFCf+WFhDvw2ql3IJ1VUbu6h/Ix1xaf/29Xh6nBE7UtgF+jXqatMPpv8I+ZDRL+",
	"7WJRzcBbQAgXpUCZXwT0AmG+4mNegJ9H9DlYCkP4q+lcB7DNgiZOT/CGETMwTyS8aKd8x+wnu7WnmvRM",
	"I3bIHXH0KcVAVroHC/4e4qgg4cyav+2An0Z55PliAcQ5TVgnx8Wm+ByrqgFM9lrbRa+wmefONbXsGf+D",
	"BrA3Rd5NWNOQaGilkc/XSNg7yc2Mk4t44E2T2S1LZs4ZEDaWJWfRO2M6Bb4lHxwSBqNr0GntHtsytTmj",
	"Qaz3/3SfnTYtZZvjSiBcxEaT9M0L0nl1SITMhZ1DpxSeBihNFXCgV9L4/a7+d5Ya9SU4e0ex3wu3lT7/",
	"bfTs6Pj35Jw+GKWfUHs6nfjp1WYHtz+sgK0KCsTf72+vp7c84e/vbuqHm9vXt9PbVgj2z5dnddt22B3M",
	"n/rbbPPdf9rhxy4P9hOySb5NALeJ2iboflq0XdoJ4Jay7Tjdj8pOTnTIb7m6KQjbbP2oMZ/cxHJpzNuM",
	"jz8dd0eT4+tkP8lPKVPx5M/r2Cdk5TWu3hFsDfAChAd/XWHMwFn893KTt398nPKkvtOJFTVKd3mcI5Z8",
	"TcC69yqHyn3mPKu113bOAlZZRgAaDSEsvPTKtz5TxvwqHaWjJtKtKDUf8//GVwkvBeZR6ctdzZsDHh78",
	"WgeMn+73cV1okzxRzYKdzEMonQ21Q56NRrF6Ootgse5opWmi75K4391znTxWdFiPY9W3Sb3vajTaH3rW",
	"B/dJ1CgN2eSyrUltPmP8tJn81MRZGPvIfcJLF3ocda1Uw86Bi66Vut9IKJog4AunVo/yzSkhfWjtNI4F",
	"JI33ejaAx5TdA1behmZoIOnkhrcDHX0F6wMyr349mT9ibxqvMMkK2HibPWjMu6acyOdDnctHF9OKfwIA",
	"AP//SQNe25kVAAA=",
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
