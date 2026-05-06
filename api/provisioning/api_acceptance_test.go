package provisioning_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
)

const testAccountID = "test_account_123"

// ProvisioningAPIRequestTest is a function that will be executed during the test.
type ProvisioningAPIRequestTest func(api *provisioning.API, ctx context.Context) (interface{}, error)

// APIType indicates which API base URL pattern to use.
type APIType int

const (
	// ProvisioningAPI uses /v1_1/provisioning/accounts/{account_id}/{path}
	ProvisioningAPI APIType = iota
	// PermissionsAPI uses /v2/accounts/{account_id}/{path}
	PermissionsAPI
)

// ExpectedProvisioningRequestParams are the expected request parameters for provisioning API tests.
type ExpectedProvisioningRequestParams struct {
	Method  string
	URI     string  // The path after the base prefix, e.g. "/access_keys" or "/permissions/policies/custom"
	Type    APIType // Which API pattern to expect (default: ProvisioningAPI)
	Params  *url.Values
	Body    *string
	Headers *map[string]string
}

// ProvisioningAPIAcceptanceTestCase is an acceptance test case for the provisioning API.
type ProvisioningAPIAcceptanceTestCase struct {
	Name              string                            // Name of the test case
	RequestTest       ProvisioningAPIRequestTest        // Function which will be called as an API request. Put SDK calls here.
	ResponseTest      cldtest.APIResponseTest           // Function which will be called to test an API response.
	ExpectedRequest   ExpectedProvisioningRequestParams // Expected HTTP request to be sent to the server
	JsonResponse      string                            // Mock of the JSON response from server. This is used to check JSON parsing.
	ExpectedStatus    string                            // Expected HTTP status of the request. This status will be returned from the HTTP mock.
	ExpectedCallCount int                               // Expected call count to the server.
	Config            *config.Configuration             // Configuration
}

// testProvisioningAPIByTestCases runs acceptance tests by given test cases.
func testProvisioningAPIByTestCases(cases []ProvisioningAPIAcceptanceTestCase, t *testing.T) {
	for num, test := range cases {
		if test.Name == "" {
			t.Skipf("Test name should be set for test #%d. Skipping it.", num)
		}

		t.Run(test.Name, func(t *testing.T) {
			callCounter := 0
			srv := getProvisioningServerMock(test.JsonResponse, t, &callCounter, test.ExpectedRequest)

			res, _ := test.RequestTest(getTestableProvisioningAPI(srv.URL, test.Config, t), ctx)
			test.ResponseTest(res, t)

			if callCounter != test.ExpectedCallCount {
				t.Errorf("Expected %d call(s), %d given", test.ExpectedCallCount, callCounter)
			}

			srv.Close()
		})
	}
}

// getTestableProvisioningAPI creates a configured provisioning API for testing.
func getTestableProvisioningAPI(mockServerURL string, c *config.Configuration, t *testing.T) *provisioning.API {
	if c == nil {
		var err error
		c, err = config.NewFromQueryParams(cldtest.CloudName, cldtest.APIKey, cldtest.APISecret,
			map[string][]string{"account_id": {testAccountID}})
		if err != nil {
			t.Error(err)
		}
	}

	c.API.UploadPrefix = mockServerURL

	api, err := provisioning.NewWithConfiguration(c)
	if err != nil {
		t.Error(err)
	}

	return api
}

// getProvisioningServerMock creates a mock HTTP server for provisioning API tests.
func getProvisioningServerMock(response string, t *testing.T, callCounter *int, ep ExpectedProvisioningRequestParams) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != ep.Method {
			t.Errorf("HTTP method should be %s, got %s", ep.Method, r.Method)
		}

		var expectedURI string
		if ep.Type == PermissionsAPI {
			expectedURI = "/v2/accounts/" + testAccountID + ep.URI
		} else {
			expectedURI = "/" + cldtest.APIVersion + "/provisioning/accounts/" + testAccountID + ep.URI
		}
		if expectedURI != r.URL.Path {
			t.Errorf("Expected request URI: %s, got: %s", expectedURI, r.URL.Path)
		}

		if ep.Params != nil && ep.Params.Encode() != r.URL.Query().Encode() {
			t.Errorf("Expected query string: %s, got: %s", ep.Params.Encode(), r.URL.Query().Encode())
		}

		if ep.Headers != nil {
			for expectedName, expectedValue := range *ep.Headers {
				value, present := r.Header[expectedName]
				if !present {
					t.Errorf("Expected request header: '%s' not found", expectedName)
				}
				stringValue := strings.Join(value, ", ")
				if expectedValue != stringValue {
					t.Errorf("Expected request header %s value: %s, got: %s", expectedName, expectedValue, stringValue)
				}
			}
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			if r.Body != nil && ep.Body != nil {
				bodyString, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				if string(bodyString) != *ep.Body {
					t.Errorf("Wrong request body. Expected: %s, given: %s", *ep.Body, string(bodyString))
				}
			}
		}

		*callCounter++
		_, err := io.WriteString(w, response)
		if err != nil {
			t.Error(err)
		}
	})

	return httptest.NewServer(handler)
}

// stringPtr returns a pointer to the given string.
func stringPtr(s string) *string {
	return &s
}
