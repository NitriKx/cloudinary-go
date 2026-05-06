package provisioning

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api"
)

const permissionsPoliciesCustom api.EndPoint = "permissions/policies/custom"

// PermissionsErrorResp is the error response from the Permissions API (v2).
type PermissionsErrorResp struct {
	Category string                 `json:"category"`
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
}

// CustomPolicyResult represents a single custom policy.
type CustomPolicyResult struct {
	ID              string `json:"id"`
	PolicyStatement string `json:"policy_statement"`
	Description     string `json:"description"`
	ScopeType       string `json:"scope_type"`
	ScopeID         string `json:"scope_id"`
	Name            string `json:"name"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}

// CreateCustomPolicyParams are the parameters for CreateCustomPolicy.
type CreateCustomPolicyParams struct {
	PolicyStatement string `json:"policy_statement"`
	Description     string `json:"description,omitempty"`
	ScopeType       string `json:"scope_type"`
	ScopeID         string `json:"scope_id,omitempty"`
	Name            string `json:"name"`
	Enabled         *bool  `json:"enabled,omitempty"`
}

// CreateCustomPolicyResult is the result of CreateCustomPolicy.
type CreateCustomPolicyResult struct {
	Data     CustomPolicyResult   `json:"data"`
	Error    PermissionsErrorResp `json:"error,omitempty"`
	Response interface{}
}

// CreateCustomPolicy creates a new custom policy for the account.
func (a *API) CreateCustomPolicy(ctx context.Context, params CreateCustomPolicyParams) (*CreateCustomPolicyResult, error) {
	res := &CreateCustomPolicyResult{}
	_, err := a.permPost(ctx, permissionsPoliciesCustom, params, res)
	return res, err
}

// ListCustomPoliciesParams are the parameters for ListCustomPolicies.
type ListCustomPoliciesParams struct {
	ScopeType string `json:"scope_type,omitempty"`
	ScopeID   string `json:"scope_id,omitempty"`
}

// ListCustomPoliciesResult is the result of ListCustomPolicies.
type ListCustomPoliciesResult struct {
	Data     []CustomPolicyResult `json:"data"`
	Error    PermissionsErrorResp `json:"error,omitempty"`
	Response interface{}
}

// ListCustomPolicies lists all custom policies for the account.
func (a *API) ListCustomPolicies(ctx context.Context, params ListCustomPoliciesParams) (*ListCustomPoliciesResult, error) {
	res := &ListCustomPoliciesResult{}
	_, err := a.permGet(ctx, permissionsPoliciesCustom, params, res)
	return res, err
}

// GetCustomPolicyParams are the parameters for GetCustomPolicy.
type GetCustomPolicyParams struct {
	PolicyID string `json:"-"`
}

// GetCustomPolicyResult is the result of GetCustomPolicy.
type GetCustomPolicyResult struct {
	Data     CustomPolicyResult   `json:"data"`
	Error    PermissionsErrorResp `json:"error,omitempty"`
	Response interface{}
}

// GetCustomPolicy gets a single custom policy by ID.
func (a *API) GetCustomPolicy(ctx context.Context, params GetCustomPolicyParams) (*GetCustomPolicyResult, error) {
	res := &GetCustomPolicyResult{}
	_, err := a.permGet(ctx, api.BuildPath(permissionsPoliciesCustom, params.PolicyID), params, res)
	return res, err
}

// UpdateCustomPolicyParams are the parameters for UpdateCustomPolicy.
type UpdateCustomPolicyParams struct {
	PolicyID        string `json:"-"`
	PolicyStatement string `json:"policy_statement,omitempty"`
	Description     string `json:"description,omitempty"`
	ScopeType       string `json:"scope_type,omitempty"`
	ScopeID         string `json:"scope_id,omitempty"`
	Name            string `json:"name,omitempty"`
	Enabled         *bool  `json:"enabled,omitempty"`
}

// UpdateCustomPolicyResult is the result of UpdateCustomPolicy.
type UpdateCustomPolicyResult struct {
	Data     CustomPolicyResult   `json:"data"`
	Error    PermissionsErrorResp `json:"error,omitempty"`
	Response interface{}
}

// UpdateCustomPolicy updates an existing custom policy.
func (a *API) UpdateCustomPolicy(ctx context.Context, params UpdateCustomPolicyParams) (*UpdateCustomPolicyResult, error) {
	res := &UpdateCustomPolicyResult{}
	_, err := a.permPut(ctx, api.BuildPath(permissionsPoliciesCustom, params.PolicyID), params, res)
	return res, err
}

// DeleteCustomPolicyParams are the parameters for DeleteCustomPolicy.
type DeleteCustomPolicyParams struct {
	PolicyID string `json:"-"`
}

// DeleteCustomPolicyResult is the result of DeleteCustomPolicy.
type DeleteCustomPolicyResult struct {
	Error    PermissionsErrorResp `json:"error,omitempty"`
	Response interface{}
}

// DeleteCustomPolicy deletes a custom policy.
func (a *API) DeleteCustomPolicy(ctx context.Context, params DeleteCustomPolicyParams) (*DeleteCustomPolicyResult, error) {
	res := &DeleteCustomPolicyResult{}
	_, err := a.permDelete(ctx, api.BuildPath(permissionsPoliciesCustom, params.PolicyID), params, res)
	return res, err
}
