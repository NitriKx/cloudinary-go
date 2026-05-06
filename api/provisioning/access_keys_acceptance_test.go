package provisioning_test

// Acceptance tests for access keys. See `TEST.md` for additional information.

import (
	"context"
	"reflect"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

const testSubAccountID = "sub_account_123"

// Acceptance test cases for access keys methods
func getAccessKeysTestCases() []ProvisioningAPIAcceptanceTestCase {
	return []ProvisioningAPIAcceptanceTestCase{
		{
			Name: "CreateAccessKey",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.CreateAccessKey(ctx, provisioning.CreateAccessKeyParams{
					SubAccountID: testSubAccountID,
					Name:         "test-key",
					Enabled:      api.Bool(true),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.CreateAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of CreateAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Name != "test-key" {
					t.Errorf("Name should be 'test-key', got '%s'", result.Name)
				}
				if result.APIKey != "key123" {
					t.Errorf("APIKey should be 'key123', got '%s'", result.APIKey)
				}
				if result.Root != false {
					t.Errorf("Root should be false")
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "POST",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys",
				Body:   stringPtr(`{"name":"test-key","enabled":true}`),
			},
			JsonResponse:      `{"name":"test-key","api_key":"key123","api_secret":"secret123","created_at":"2026-05-05T08:51:34Z","updated_at":"2026-05-05T08:51:34Z","enabled":true,"root":false}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "ListAccessKeys",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.ListAccessKeys(ctx, provisioning.ListAccessKeysParams{
					SubAccountID: testSubAccountID,
					PageSize:     10,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.ListAccessKeysResult)
				if !ok {
					t.Errorf("Response should be type of ListAccessKeysResult, %s given", reflect.TypeOf(response))
				}
				if len(result.AccessKeys) != 1 {
					t.Errorf("Expected 1 access key, got %d", len(result.AccessKeys))
				}
				if result.Total != 1 {
					t.Errorf("Total should be 1, got %d", result.Total)
				}
				if result.Active != 1 {
					t.Errorf("Active should be 1, got %d", result.Active)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys",
			},
			JsonResponse:      `{"access_keys":[{"name":"key-1","api_key":"key1","api_secret":"secret1","enabled":true,"root":false}],"total":1,"active":1}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "GetAccessKey",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetAccessKey(ctx, provisioning.GetAccessKeyParams{
					SubAccountID: testSubAccountID,
					APIKey:       "key123",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of GetAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.APIKey != "key123" {
					t.Errorf("APIKey should be 'key123', got '%s'", result.APIKey)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys/key123",
			},
			JsonResponse:      `{"name":"test","api_key":"key123","api_secret":"secret123","enabled":true,"root":false}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "UpdateAccessKey",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.UpdateAccessKey(ctx, provisioning.UpdateAccessKeyParams{
					SubAccountID: testSubAccountID,
					APIKey:       "key123",
					Name:         "updated-key",
					DedicatedFor: "webhooks",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.UpdateAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of UpdateAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Name != "updated-key" {
					t.Errorf("Name should be 'updated-key', got '%s'", result.Name)
				}
				if result.DedicatedFor != "webhooks" {
					t.Errorf("DedicatedFor should be 'webhooks', got '%s'", result.DedicatedFor)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "PUT",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys/key123",
				Body:   stringPtr(`{"name":"updated-key","dedicated_for":"webhooks"}`),
			},
			JsonResponse:      `{"name":"updated-key","api_key":"key123","api_secret":"*****secret","enabled":true,"root":false,"dedicated_for":"webhooks"}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "DisableAccessKey",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.UpdateAccessKey(ctx, provisioning.UpdateAccessKeyParams{
					SubAccountID: testSubAccountID,
					APIKey:       "key123",
					Enabled:      api.Bool(false),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.UpdateAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of UpdateAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Enabled {
					t.Errorf("Enabled should be false, got true")
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "PUT",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys/key123",
				Body:   stringPtr(`{"enabled":false}`),
			},
			JsonResponse:      `{"name":"test-key","api_key":"key123","api_secret":"*****secret","enabled":false,"root":false}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "DeleteAccessKeyByName",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.DeleteAccessKeyByName(ctx, provisioning.DeleteAccessKeyByNameParams{
					SubAccountID: testSubAccountID,
					Name:         "test-key",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.DeleteAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of DeleteAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Message != "ok" {
					t.Errorf("Message should be 'ok', got '%s'", result.Message)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "DELETE",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys",
				Body:   stringPtr(`{"name":"test-key"}`),
			},
			JsonResponse:      `{"message":"ok"}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "DeleteAccessKey",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.DeleteAccessKey(ctx, provisioning.DeleteAccessKeyParams{
					SubAccountID: testSubAccountID,
					APIKey:       "key123",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.DeleteAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of DeleteAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Message != "ok" {
					t.Errorf("Message should be 'ok', got '%s'", result.Message)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "DELETE",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys/key123",
			},
			JsonResponse:      `{"message":"ok"}`,
			ExpectedCallCount: 1,
		},
		{
			Name: "AccessKey error check",
			RequestTest: func(a *provisioning.API, ctx context.Context) (interface{}, error) {
				return a.GetAccessKey(ctx, provisioning.GetAccessKeyParams{
					SubAccountID: testSubAccountID,
					APIKey:       "nonexistent",
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				result, ok := response.(*provisioning.GetAccessKeyResult)
				if !ok {
					t.Errorf("Response should be type of GetAccessKeyResult, %s given", reflect.TypeOf(response))
				}
				if result.Error.Message != "not found" {
					t.Errorf("Error message should be 'not found', got '%s'", result.Error.Message)
				}
			},
			ExpectedRequest: ExpectedProvisioningRequestParams{
				Method: "GET",
				URI:    "/sub_accounts/" + testSubAccountID + "/access_keys/nonexistent",
			},
			JsonResponse:      `{"error":{"message":"not found"}}`,
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestAccessKeys_Acceptance(t *testing.T) {
	t.Parallel()
	testProvisioningAPIByTestCases(getAccessKeysTestCases(), t)
}
