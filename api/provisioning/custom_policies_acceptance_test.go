package provisioning_test

// Acceptance tests for custom policies. See `TEST.md` for additional information.

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

// Acceptance test cases for custom policies methods
func getCustomPoliciesTestCases() []ProvisioningAPIAcceptanceTestCase {
	return []ProvisioningAPIAcceptanceTestCase{
		{
			Name: "CreateCustomPolicy",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.CreateCustomPolicy(ctx, provisioning.CreateCustomPolicyParams{
					Name:            "test-policy",
					PolicyStatement: "permit(principal, action, resource);",
					ScopeType:       "account",
					Enabled:         api.Bool(true),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.CreateCustomPolicyResult)
				if !ok {
					t.Errorf("Response should be type of CreateCustomPolicyResult, %s given", reflect.TypeOf(response))
				}
				if result.Data.ID != "policy123" {
					t.Errorf("ID should be 'policy123', got '%s'", result.Data.ID)
				}
				if result.Data.Name != "test-policy" {
					t.Errorf("Name should be 'test-policy', got '%s'", result.Data.Name)
				}
				if result.Data.ScopeType != "account" {
					t.Errorf("ScopeType should be 'account', got '%s'", result.Data.ScopeType)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "POST",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom",
				Body:   stringPtr(`{"policy_statement":"permit(principal, action, resource);","scope_type":"account","name":"test-policy","enabled":true}`),
			},
			JsonResponse:      `{"data":{"id":"policy123","policy_statement":"permit(principal, action, resource);","description":"","scope_type":"account","scope_id":"","name":"test-policy","enabled":true,"created_at":1777970973,"updated_at":1777970973}}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "ListCustomPolicies",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.ListCustomPolicies(ctx, provisioning.ListCustomPoliciesParams{})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.ListCustomPoliciesResult)
				if !ok {
					t.Errorf("Response should be type of ListCustomPoliciesResult, %s given", reflect.TypeOf(response))
				}
				if len(result.Data) != 1 {
					t.Errorf("Expected 1 custom policy, got %d", len(result.Data))
				}
				if result.Data[0].Name != "policy-1" {
					t.Errorf("Name should be 'policy-1', got '%s'", result.Data[0].Name)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom",
			},
			JsonResponse:      `{"data":[{"id":"policy1","policy_statement":"permit(principal, action, resource);","scope_type":"account","scope_id":"","name":"policy-1","enabled":true,"created_at":1777970973,"updated_at":1777970973}]}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "GetCustomPolicy",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetCustomPolicy(ctx, provisioning.GetCustomPolicyParams{PolicyID: "policy123"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetCustomPolicyResult)
				if !ok {
					t.Errorf("Response should be type of GetCustomPolicyResult, %s given", reflect.TypeOf(response))
				}
				if result.Data.ID != "policy123" {
					t.Errorf("ID should be 'policy123', got '%s'", result.Data.ID)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom/policy123",
			},
			JsonResponse:      `{"data":{"id":"policy123","policy_statement":"permit(principal, action, resource);","scope_type":"account","scope_id":"","name":"test-policy","enabled":true,"created_at":1777970973,"updated_at":1777970973}}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "UpdateCustomPolicy",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.UpdateCustomPolicy(ctx, provisioning.UpdateCustomPolicyParams{
					PolicyID: "policy123",
					Name:     "updated-policy",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.UpdateCustomPolicyResult)
				if !ok {
					t.Errorf("Response should be type of UpdateCustomPolicyResult, %s given", reflect.TypeOf(response))
				}
				if result.Data.Name != "updated-policy" {
					t.Errorf("Name should be 'updated-policy', got '%s'", result.Data.Name)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "PUT",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom/policy123",
				Body:   stringPtr(`{"name":"updated-policy"}`),
			},
			JsonResponse:      `{"data":{"id":"policy123","policy_statement":"permit(principal, action, resource);","scope_type":"account","scope_id":"","name":"updated-policy","enabled":true,"created_at":1777970973,"updated_at":1777970973}}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "DeleteCustomPolicy",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.DeleteCustomPolicy(ctx, provisioning.DeleteCustomPolicyParams{PolicyID: "policy123"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*provisioning.DeleteCustomPolicyResult)
				if !ok {
					t.Errorf("Response should be type of DeleteCustomPolicyResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "DELETE",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom/policy123",
			},
			JsonResponse:      ``,
			ExpectedCallCount: 1,
		},
		{
			Name: "CustomPolicy error check",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetCustomPolicy(ctx, provisioning.GetCustomPolicyParams{PolicyID: "nonexistent"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetCustomPolicyResult)
				if !ok {
					t.Errorf("Response should be type of GetCustomPolicyResult, %s given", reflect.TypeOf(response))
				}
				if result.Error.Message != "not found" {
					t.Errorf("Error message should be 'not found', got '%s'", result.Error.Message)
				}
				if result.Error.Code != "ACCOUNTS_00003" {
					t.Errorf("Error code should be 'ACCOUNTS_00003', got '%s'", result.Error.Code)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				Type:   PermissionsAPI,
				URI:    "/permissions/policies/custom/nonexistent",
			},
			JsonResponse:      `{"error":{"category":"user_error","code":"ACCOUNTS_00003","message":"not found","details":{"policy_id":"nonexistent"}}}`,
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestCustomPolicies_Acceptance(t *testing.T) {
	t.Parallel()
	testProvisioningAPIByTestCases(getCustomPoliciesTestCases(), t)
}
