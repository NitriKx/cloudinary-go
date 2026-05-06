package provisioning_test

// Acceptance tests for product environments. See `TEST.md` for additional information.

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

// Acceptance test cases for product environments methods
func getProductEnvironmentsTestCases() []ProvisioningAPIAcceptanceTestCase {
	return []ProvisioningAPIAcceptanceTestCase{
		{
			Name: "CreateProductEnvironment",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.CreateProductEnvironment(ctx, provisioning.CreateProductEnvironmentParams{
					Name:      "test-env",
					CloudName: "test-cloud",
					Enabled:   api.Bool(true),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.CreateProductEnvironmentResult)
				if !ok {
					t.Errorf("Response should be type of CreateProductEnvironmentResult, %s given", reflect.TypeOf(response))
				}
				if result.ID != "sub123" {
					t.Errorf("ID should be 'sub123', got '%s'", result.ID)
				}
				if result.Name != "test-env" {
					t.Errorf("Name should be 'test-env', got '%s'", result.Name)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "POST",
				URI:    "/sub_accounts",
				Body:   stringPtr(`{"name":"test-env","cloud_name":"test-cloud","enabled":true}`),
			},
			JsonResponse:      `{"id":"sub123","name":"test-env","cloud_name":"test-cloud","enabled":true}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "ListProductEnvironments",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.ListProductEnvironments(ctx, provisioning.ListProductEnvironmentsParams{
					Prefix: "backmarket",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.ListProductEnvironmentsResult)
				if !ok {
					t.Errorf("Response should be type of ListProductEnvironmentsResult, %s given", reflect.TypeOf(response))
				}
				if len(result.SubAccounts) != 1 {
					t.Errorf("Expected 1 product environment, got %d", len(result.SubAccounts))
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts",
				Params: &url.Values{"prefix": []string{"backmarket"}},
			},
			JsonResponse:      `{"sub_accounts":[{"id":"sub1","name":"env-1","cloud_name":"cloud-1","enabled":true}]}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "GetProductEnvironment",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetProductEnvironment(ctx, provisioning.GetProductEnvironmentParams{SubAccountID: "sub123"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetProductEnvironmentResult)
				if !ok {
					t.Errorf("Response should be type of GetProductEnvironmentResult, %s given", reflect.TypeOf(response))
				}
				if result.ID != "sub123" {
					t.Errorf("ID should be 'sub123', got '%s'", result.ID)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts/sub123",
			},
			JsonResponse:      `{"id":"sub123","name":"test-env","cloud_name":"test-cloud","enabled":true}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "UpdateProductEnvironment",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.UpdateProductEnvironment(ctx, provisioning.UpdateProductEnvironmentParams{
					SubAccountID: "sub123",
					Name:         "updated-env",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.UpdateProductEnvironmentResult)
				if !ok {
					t.Errorf("Response should be type of UpdateProductEnvironmentResult, %s given", reflect.TypeOf(response))
				}
				if result.Name != "updated-env" {
					t.Errorf("Name should be 'updated-env', got '%s'", result.Name)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "PUT",
				URI:    "/sub_accounts/sub123",
				Body:   stringPtr(`{"name":"updated-env"}`),
			},
			JsonResponse:      `{"id":"sub123","name":"updated-env","cloud_name":"test-cloud","enabled":true}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "DeleteProductEnvironment",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.DeleteProductEnvironment(ctx, provisioning.DeleteProductEnvironmentParams{SubAccountID: "sub123"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.DeleteProductEnvironmentResult)
				if !ok {
					t.Errorf("Response should be type of DeleteProductEnvironmentResult, %s given", reflect.TypeOf(response))
				}
				if result.Message != "ok" {
					t.Errorf("Message should be 'ok', got '%s'", result.Message)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "DELETE",
				URI:    "/sub_accounts/sub123",
			},
			JsonResponse:      `{"message":"ok"}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "ProductEnvironment error check",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetProductEnvironment(ctx, provisioning.GetProductEnvironmentParams{SubAccountID: "nonexistent"})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetProductEnvironmentResult)
				if !ok {
					t.Errorf("Response should be type of GetProductEnvironmentResult, %s given", reflect.TypeOf(response))
				}
				if result.Error.Message != "not found" {
					t.Errorf("Error message should be 'not found', got '%s'", result.Error.Message)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts/nonexistent",
			},
			JsonResponse:      `{"error":{"message":"not found"}}`,
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestProductEnvironments_Acceptance(t *testing.T) {
	t.Parallel()
	testProvisioningAPIByTestCases(getProductEnvironmentsTestCases(), t)
}
