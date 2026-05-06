package provisioning_test

import (
	"context"
	"os"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
	"github.com/cloudinary/cloudinary-go/v2/config"
)

var ctx = context.Background()

// TestNew_AccountURL verifies that provisioning.New() picks up CLOUDINARY_ACCOUNT_URL.
func TestNew_AccountURL(t *testing.T) {
	t.Setenv("CLOUDINARY_ACCOUNT_URL", "account://mykey:mysecret@my-account-id")
	// Unset CLOUDINARY_URL so we don't accidentally fall through to it.
	t.Setenv("CLOUDINARY_URL", "")

	a, err := provisioning.New()
	if err != nil {
		t.Fatalf("provisioning.New() error: %v", err)
	}
	if a.Config.Cloud.APIKey != "mykey" {
		t.Errorf("expected APIKey 'mykey', got '%s'", a.Config.Cloud.APIKey)
	}
	if a.Config.Cloud.APISecret != "mysecret" {
		t.Errorf("expected APISecret 'mysecret', got '%s'", a.Config.Cloud.APISecret)
	}
	if a.Config.Cloud.AccountID != "my-account-id" {
		t.Errorf("expected AccountID 'my-account-id', got '%s'", a.Config.Cloud.AccountID)
	}
}

// Live test credentials are read from environment variables.
//
// Required:
//
//	CLOUDINARY_ACCOUNT_API_KEY    - Account management API key
//	CLOUDINARY_ACCOUNT_API_SECRET - Account management API secret
//	CLOUDINARY_ACCOUNT_ID         - Account ID (UUID)
//	CLOUDINARY_SUB_ACCOUNT_ID     - Product environment ID for access key tests
//
// Run with:
//
//	go test ./api/provisioning/... -v
func liveEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("skipping: environment variable %s is not set", key)
	}
	return v
}

func newLiveProvisioningAPI(t *testing.T) *provisioning.API {
	t.Helper()
	apiKey := liveEnv(t, "CLOUDINARY_ACCOUNT_API_KEY")
	apiSecret := liveEnv(t, "CLOUDINARY_ACCOUNT_API_SECRET")
	accountID := liveEnv(t, "CLOUDINARY_ACCOUNT_ID")

	c, err := config.NewFromQueryParams("", apiKey, apiSecret, map[string][]string{
		"account_id": {accountID},
	})
	if err != nil {
		t.Fatalf("failed to build config: %v", err)
	}
	a, err := provisioning.NewWithConfiguration(c)
	if err != nil {
		t.Fatalf("failed to create provisioning API: %v", err)
	}
	return a
}

func liveSubAccountID(t *testing.T) string {
	t.Helper()
	return liveEnv(t, "CLOUDINARY_SUB_ACCOUNT_ID")
}
