// Package provisioning is used for accessing Cloudinary Provisioning API functionality.
package provisioning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/logger"
)

// API is the Provisioning API struct.
type API struct {
	Config config.Configuration
	Logger *logger.Logger
	Client http.Client
}

// New creates a new Provisioning API instance from environment variables.
//
// It checks CLOUDINARY_ACCOUNT_URL first. The expected format is:
//
//	account://<ACCOUNT_API_KEY>:<ACCOUNT_API_SECRET>@<ACCOUNT_ID>
//
// If CLOUDINARY_ACCOUNT_URL is not set, it falls back to CLOUDINARY_URL.
func New() (*API, error) {
	if accountURL := os.Getenv("CLOUDINARY_ACCOUNT_URL"); accountURL != "" {
		c, err := config.NewFromAccountURL(accountURL)
		if err != nil {
			return nil, err
		}
		return NewWithConfiguration(c)
	}

	c, err := config.New()
	if err != nil {
		return nil, err
	}
	return NewWithConfiguration(c)
}

// NewWithConfiguration creates a new Provisioning API instance with the given Configuration.
func NewWithConfiguration(c *config.Configuration) (*API, error) {
	return &API{
		Config: *c,
		Client: http.Client{},
		Logger: logger.New(),
	}, nil
}

func (a *API) get(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodGet, path, requestParams, result)
}

func (a *API) post(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodPost, path, requestParams, result)
}

func (a *API) put(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodPut, path, requestParams, result)
}

func (a *API) delete(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodDelete, path, requestParams, result)
}

func (a *API) permGet(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callPermissionsAPI(ctx, http.MethodGet, path, requestParams, result)
}

func (a *API) permPost(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callPermissionsAPI(ctx, http.MethodPost, path, requestParams, result)
}

func (a *API) permPut(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callPermissionsAPI(ctx, http.MethodPut, path, requestParams, result)
}

func (a *API) permDelete(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callPermissionsAPI(ctx, http.MethodDelete, path, requestParams, result)
}

func (a *API) callPermissionsAPI(ctx context.Context, method string, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	var body io.Reader = nil
	var queryParams = ""

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		jsonReq, err := json.Marshal(requestParams)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	}

	if body == nil {
		params, err := api.StructToParams(requestParams)
		if err != nil {
			return nil, err
		}
		queryParams = params.Encode()
	}

	return a.executePermissionsRequest(ctx, method, path, body, queryParams, map[string]string{}, result)
}

func (a *API) executePermissionsRequest(ctx context.Context, method string, path interface{}, body io.Reader, queryParams string,
	headers map[string]string, result interface{}) (*http.Response, error) {

	if a.Config.Cloud.AccountID == "" {
		return nil, fmt.Errorf("AccountID is required for Provisioning API calls")
	}

	req, err := http.NewRequest(method,
		fmt.Sprintf("%v/v2/accounts/%v/%v",
			a.Config.API.UploadPrefix,
			a.Config.Cloud.AccountID,
			api.BuildPath(path)),
		body,
	)
	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	req.URL.RawQuery = queryParams

	req.Header.Set("User-Agent", api.GetUserAgent())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	setAuth(a, req)

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.Config.API.Timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a.Logger.Debug(string(bodyBytes))

	if len(bodyBytes) > 0 {
		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			return resp, err
		}

		err = api.HandleRawResponse(bodyBytes, result)
	}

	return resp, err
}

func (a *API) callAPI(ctx context.Context, method string, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	var body io.Reader = nil
	var queryParams = ""

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		jsonReq, err := json.Marshal(requestParams)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	}

	if body == nil {
		params, err := api.StructToParams(requestParams)
		if err != nil {
			return nil, err
		}
		queryParams = params.Encode()
	}

	return a.executeRequest(ctx, method, path, body, queryParams, map[string]string{}, result)
}

func (a *API) executeRequest(ctx context.Context, method string, path interface{}, body io.Reader, queryParams string,
	headers map[string]string, result interface{}) (*http.Response, error) {

	if a.Config.Cloud.AccountID == "" {
		return nil, fmt.Errorf("AccountID is required for Provisioning API calls")
	}

	apiVersion := ""
	if apiVersionRaw := ctx.Value("api_version"); apiVersionRaw != nil {
		apiVersion = apiVersionRaw.(string)
	}

	req, err := http.NewRequest(method,
		fmt.Sprintf("%v/provisioning/accounts/%v/%v",
			api.BaseURL(a.Config.API.UploadPrefix, apiVersion),
			a.Config.Cloud.AccountID,
			api.BuildPath(path)),
		body,
	)
	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	req.URL.RawQuery = queryParams

	req.Header.Set("User-Agent", api.GetUserAgent())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	setAuth(a, req)

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.Config.API.Timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a.Logger.Debug(string(bodyBytes))

	if len(bodyBytes) > 0 {
		err = json.Unmarshal(bodyBytes, result)
		if err != nil {
			return resp, err
		}

		err = api.HandleRawResponse(bodyBytes, result)
	}

	return resp, err
}

func setAuth(a *API, req *http.Request) {
	if a.Config.Cloud.OAuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Config.Cloud.OAuthToken))
	} else {
		req.SetBasicAuth(a.Config.Cloud.APIKey, a.Config.Cloud.APISecret)
	}
}
