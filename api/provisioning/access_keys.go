package provisioning

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
)

const accessKeysPath api.EndPoint = "access_keys"

// AccessKeyResult represents a single access key.
type AccessKeyResult struct {
	Name         string    `json:"name"`
	APIKey       string    `json:"api_key"`
	APISecret    string    `json:"api_secret"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Enabled      bool      `json:"enabled"`
	Root         bool      `json:"root"`
	DedicatedFor string    `json:"dedicated_for,omitempty"`
}

// CreateAccessKeyParams are the parameters for CreateAccessKey.
type CreateAccessKeyParams struct {
	SubAccountID string `json:"-"`
	Name         string `json:"name,omitempty"`
	Enabled      *bool  `json:"enabled,omitempty"`
}

// CreateAccessKeyResult is the result of CreateAccessKey.
type CreateAccessKeyResult struct {
	AccessKeyResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// CreateAccessKey creates a new access key for a product environment.
//
// https://cloudinary.com/documentation/provisioning_api#create_access_key
func (a *API) CreateAccessKey(ctx context.Context, params CreateAccessKeyParams) (*CreateAccessKeyResult, error) {
	res := &CreateAccessKeyResult{}
	_, err := a.post(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath), params, res)
	return res, err
}

// ListAccessKeysParams are the parameters for ListAccessKeys.
type ListAccessKeysParams struct {
	SubAccountID string `json:"-"`
	PageSize     int    `json:"page_size,omitempty"`
	Page         int    `json:"page,omitempty"`
	SortBy       string `json:"sort_by,omitempty"`
	SortOrder    string `json:"sort_order,omitempty"`
}

// ListAccessKeysResult is the result of ListAccessKeys.
type ListAccessKeysResult struct {
	AccessKeys []AccessKeyResult `json:"access_keys"`
	Total      int               `json:"total"`
	Active     int               `json:"active"`
	Error      api.ErrorResp     `json:"error,omitempty"`
	Response   interface{}
}

// ListAccessKeys lists all access keys for a product environment.
//
// https://cloudinary.com/documentation/provisioning_api#get_access_keys
func (a *API) ListAccessKeys(ctx context.Context, params ListAccessKeysParams) (*ListAccessKeysResult, error) {
	res := &ListAccessKeysResult{}
	_, err := a.get(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath), params, res)
	return res, err
}

// GetAccessKeyParams are the parameters for GetAccessKey.
type GetAccessKeyParams struct {
	SubAccountID string `json:"-"`
	APIKey       string `json:"-"`
}

// GetAccessKeyResult is the result of GetAccessKey.
type GetAccessKeyResult struct {
	AccessKeyResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// GetAccessKey gets a single access key by its API key.
//
// https://cloudinary.com/documentation/provisioning_api#get_access_key
func (a *API) GetAccessKey(ctx context.Context, params GetAccessKeyParams) (*GetAccessKeyResult, error) {
	res := &GetAccessKeyResult{}
	_, err := a.get(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath, params.APIKey), params, res)
	return res, err
}

// UpdateAccessKeyParams are the parameters for UpdateAccessKey.
type UpdateAccessKeyParams struct {
	SubAccountID string `json:"-"`
	APIKey       string `json:"-"`
	Name         string `json:"name,omitempty"`
	Enabled      *bool  `json:"enabled,omitempty"`
	// DedicatedFor associates the key with a specific Cloudinary feature.
	DedicatedFor string `json:"dedicated_for,omitempty"`
}

// UpdateAccessKeyResult is the result of UpdateAccessKey.
type UpdateAccessKeyResult struct {
	AccessKeyResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// UpdateAccessKey updates an existing access key.
//
// https://cloudinary.com/documentation/provisioning_api#update_access_key
func (a *API) UpdateAccessKey(ctx context.Context, params UpdateAccessKeyParams) (*UpdateAccessKeyResult, error) {
	res := &UpdateAccessKeyResult{}
	_, err := a.put(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath, params.APIKey), params, res)
	return res, err
}

// DeleteAccessKeyParams are the parameters for DeleteAccessKey.
type DeleteAccessKeyParams struct {
	SubAccountID string `json:"-"`
	APIKey       string `json:"-"`
}

// DeleteAccessKeyResult is the result of DeleteAccessKey.
type DeleteAccessKeyResult struct {
	Message  string        `json:"message"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// DeleteAccessKey deletes an access key.
//
// https://cloudinary.com/documentation/provisioning_api#delete_access_key
func (a *API) DeleteAccessKey(ctx context.Context, params DeleteAccessKeyParams) (*DeleteAccessKeyResult, error) {
	res := &DeleteAccessKeyResult{}
	_, err := a.delete(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath, params.APIKey), params, res)
	return res, err
}

// DeleteAccessKeyByNameParams are the parameters for DeleteAccessKeyByName.
type DeleteAccessKeyByNameParams struct {
	SubAccountID string `json:"-"`
	Name         string `json:"name"`
}

// DeleteAccessKeyByName deletes an access key by its display name. The request
// targets the access-keys list URL with the name in the body, mirroring the
// Node SDK's delete_access_key_by_name.
//
// https://cloudinary.com/documentation/provisioning_api#delete_access_key
func (a *API) DeleteAccessKeyByName(ctx context.Context, params DeleteAccessKeyByNameParams) (*DeleteAccessKeyResult, error) {
	res := &DeleteAccessKeyResult{}
	_, err := a.delete(ctx, api.BuildPath(subAccounts, params.SubAccountID, accessKeysPath), params, res)
	return res, err
}
