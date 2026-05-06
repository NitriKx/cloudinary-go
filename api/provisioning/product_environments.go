package provisioning

import (
	"context"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
)

const subAccounts api.EndPoint = "sub_accounts"

// APIAccessKey represents an API access key embedded in a product environment response.
type APIAccessKey struct {
	Key     string `json:"key"`
	Secret  string `json:"secret"`
	Enabled bool   `json:"enabled"`
}

// ProductEnvironmentResult represents a single product environment (sub-account).
type ProductEnvironmentResult struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	CloudName        string            `json:"cloud_name"`
	APIAccessKeys    []APIAccessKey    `json:"api_access_keys,omitempty"`
	CustomAttributes map[string]string `json:"custom_attributes,omitempty"`
	Enabled          bool              `json:"enabled"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	// EdgeKey is only populated on the create response.
	EdgeKey string `json:"edge_key,omitempty"`
}

// CreateProductEnvironmentParams are the parameters for CreateProductEnvironment.
type CreateProductEnvironmentParams struct {
	Name             string            `json:"name"`
	CloudName        string            `json:"cloud_name,omitempty"`
	CustomAttributes map[string]string `json:"custom_attributes,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	// BaseSubAccountID, when set, tells Cloudinary to copy the settings (and
	// access bindings of users with global access) from the given sub-account
	// to the newly-created one.
	BaseSubAccountID string `json:"base_sub_account_id,omitempty"`
}

// CreateProductEnvironmentResult is the result of CreateProductEnvironment.
type CreateProductEnvironmentResult struct {
	ProductEnvironmentResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// CreateProductEnvironment creates a new product environment (sub-account).
//
// https://cloudinary.com/documentation/provisioning_api#create_sub_account
func (a *API) CreateProductEnvironment(ctx context.Context, params CreateProductEnvironmentParams) (*CreateProductEnvironmentResult, error) {
	res := &CreateProductEnvironmentResult{}
	_, err := a.post(ctx, subAccounts, params, res)
	return res, err
}

// ListProductEnvironmentsParams are the parameters for ListProductEnvironments.
type ListProductEnvironmentsParams struct {
	Enabled *bool    `json:"enabled,omitempty"`
	IDs     []string `json:"ids,omitempty"`
	// Prefix returns accounts whose name begins with the given case-insensitive
	// string. Ignored when IDs is non-empty.
	Prefix string `json:"prefix,omitempty"`
}

// ListProductEnvironmentsResult is the result of ListProductEnvironments.
type ListProductEnvironmentsResult struct {
	SubAccounts []ProductEnvironmentResult `json:"sub_accounts"`
	Error       api.ErrorResp              `json:"error,omitempty"`
	Response    interface{}
}

// ListProductEnvironments lists all product environments (sub-accounts).
//
// https://cloudinary.com/documentation/provisioning_api#get_sub_accounts
func (a *API) ListProductEnvironments(ctx context.Context, params ListProductEnvironmentsParams) (*ListProductEnvironmentsResult, error) {
	res := &ListProductEnvironmentsResult{}
	_, err := a.get(ctx, subAccounts, params, res)
	return res, err
}

// GetProductEnvironmentParams are the parameters for GetProductEnvironment.
type GetProductEnvironmentParams struct {
	SubAccountID string `json:"-"`
}

// GetProductEnvironmentResult is the result of GetProductEnvironment.
type GetProductEnvironmentResult struct {
	ProductEnvironmentResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// GetProductEnvironment gets a single product environment by ID.
//
// https://cloudinary.com/documentation/provisioning_api#get_sub_account
func (a *API) GetProductEnvironment(ctx context.Context, params GetProductEnvironmentParams) (*GetProductEnvironmentResult, error) {
	res := &GetProductEnvironmentResult{}
	_, err := a.get(ctx, api.BuildPath(subAccounts, params.SubAccountID), params, res)
	return res, err
}

// UpdateProductEnvironmentParams are the parameters for UpdateProductEnvironment.
type UpdateProductEnvironmentParams struct {
	SubAccountID     string            `json:"-"`
	Name             string            `json:"name,omitempty"`
	CloudName        string            `json:"cloud_name,omitempty"`
	CustomAttributes map[string]string `json:"custom_attributes,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
}

// UpdateProductEnvironmentResult is the result of UpdateProductEnvironment.
type UpdateProductEnvironmentResult struct {
	ProductEnvironmentResult
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// UpdateProductEnvironment updates an existing product environment.
//
// https://cloudinary.com/documentation/provisioning_api#update_sub_account
func (a *API) UpdateProductEnvironment(ctx context.Context, params UpdateProductEnvironmentParams) (*UpdateProductEnvironmentResult, error) {
	res := &UpdateProductEnvironmentResult{}
	_, err := a.put(ctx, api.BuildPath(subAccounts, params.SubAccountID), params, res)
	return res, err
}

// DeleteProductEnvironmentParams are the parameters for DeleteProductEnvironment.
type DeleteProductEnvironmentParams struct {
	SubAccountID string `json:"-"`
}

// DeleteProductEnvironmentResult is the result of DeleteProductEnvironment.
type DeleteProductEnvironmentResult struct {
	Message  string        `json:"message"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// DeleteProductEnvironment deletes a product environment.
//
// https://cloudinary.com/documentation/provisioning_api#delete_sub_account
func (a *API) DeleteProductEnvironment(ctx context.Context, params DeleteProductEnvironmentParams) (*DeleteProductEnvironmentResult, error) {
	res := &DeleteProductEnvironmentResult{}
	_, err := a.delete(ctx, api.BuildPath(subAccounts, params.SubAccountID), params, res)
	return res, err
}
