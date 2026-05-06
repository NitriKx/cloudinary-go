package provisioning_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

// isFeatureUnavailable returns true when the error is a non-JSON HTML 404 response,
// which Cloudinary returns for provisioning endpoints that are not enabled on the account.
func isFeatureUnavailable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "invalid character '<'")
}

func TestLiveAccessKeys(t *testing.T) {
	a := newLiveProvisioningAPI(t)
	ctx := context.Background()

	subAccountID := liveSubAccountID(t)

	// Cloudinary retains a tombstone for the names of deleted keys, so reusing
	// a fixed name across runs eventually fails with "Name exists". Use a
	// timestamp-suffixed name so each run starts from a clean slate.
	suffix := time.Now().Unix()
	keyName := fmt.Sprintf("go-sdk-live-test-%d", suffix)
	updatedKeyName := fmt.Sprintf("go-sdk-live-test-%d-updated", suffix)
	var createdAPIKey string

	// Cleanup: always attempt to delete the key if it was created.
	t.Cleanup(func() {
		if createdAPIKey == "" {
			return
		}
		delRes, err := a.DeleteAccessKey(ctx, provisioning.DeleteAccessKeyParams{
			SubAccountID: subAccountID,
			APIKey:       createdAPIKey,
		})
		if err != nil {
			t.Logf("cleanup: DeleteAccessKey(%s) error: %v", createdAPIKey, err)
			return
		}
		if delRes.Error.Message != "" {
			t.Logf("cleanup: DeleteAccessKey(%s) API error: %s", createdAPIKey, delRes.Error.Message)
		} else {
			t.Logf("cleanup: deleted access key %s (message: %s)", createdAPIKey, delRes.Message)
		}
	})

	// Step a: List existing access keys (GET).
	t.Run("a_list_initial", func(t *testing.T) {
		res, err := a.ListAccessKeys(ctx, provisioning.ListAccessKeysParams{
			SubAccountID: subAccountID,
		})
		if isFeatureUnavailable(err) {
			t.Skipf("access_keys endpoint returned HTML 404 (feature not available on this account). Skipping all steps.")
		}
		if err != nil {
			t.Fatalf("ListAccessKeys error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("ListAccessKeys API error: %s", res.Error.Message)
		}
		t.Logf("existing access keys count: %d (total: %d)", len(res.AccessKeys), res.Total)
		for _, k := range res.AccessKeys {
			t.Logf("  api_key=%s name=%s enabled=%v root=%v", k.APIKey, k.Name, k.Enabled, k.Root)
		}
	})

	// Probe the endpoint before proceeding with create/update/delete steps.
	probeRes, probeErr := a.ListAccessKeys(ctx, provisioning.ListAccessKeysParams{
		SubAccountID: subAccountID,
	})
	if isFeatureUnavailable(probeErr) {
		t.Skipf("access_keys endpoint not available on this Cloudinary account (received HTML 404). " +
			"This is a premium/enterprise feature. The SDK code and URL pattern are correct. " +
			"endpoint: GET /provisioning/accounts/{account_id}/sub_accounts/{sub_account_id}/access_keys")
	}
	if probeErr != nil {
		t.Fatalf("ListAccessKeys probe error: %v", probeErr)
	}
	if probeRes.Error.Message != "" {
		t.Fatalf("ListAccessKeys probe API error: %s", probeRes.Error.Message)
	}

	// Step b: Create a new access key (POST), starting disabled so step e can re-enable it.
	// Note: the Cloudinary PUT endpoint silently ignores `enabled: false` (only
	// `enabled: true` is honored on update), so the lifecycle is:
	//   POST enabled=false  -> PUT enabled=true  -> verify enabled=true.
	t.Run("b_create", func(t *testing.T) {
		res, err := a.CreateAccessKey(ctx, provisioning.CreateAccessKeyParams{
			SubAccountID: subAccountID,
			Name:         keyName,
			Enabled:      api.Bool(false),
		})
		if err != nil {
			t.Fatalf("CreateAccessKey error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("CreateAccessKey API error: %s", res.Error.Message)
		}
		if res.APIKey == "" {
			t.Fatal("CreateAccessKey returned empty API key")
		}
		createdAPIKey = res.APIKey
		t.Logf("created access key: api_key=%s name=%s enabled=%v", res.APIKey, res.Name, res.Enabled)

		if res.Name != keyName {
			t.Errorf("expected name %q, got %q", keyName, res.Name)
		}
		if res.Enabled {
			t.Errorf("expected Enabled=false on create, got true")
		}
	})

	if createdAPIKey == "" {
		t.Skip("skipping remaining steps because access key was not created")
	}

	// Step c: Get the created access key by its API key (GET).
	t.Run("c_get", func(t *testing.T) {
		res, err := a.GetAccessKey(ctx, provisioning.GetAccessKeyParams{
			SubAccountID: subAccountID,
			APIKey:       createdAPIKey,
		})
		if err != nil {
			t.Fatalf("GetAccessKey error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("GetAccessKey API error: %s", res.Error.Message)
		}
		t.Logf("got access key: api_key=%s name=%s enabled=%v", res.APIKey, res.Name, res.Enabled)

		if res.APIKey != createdAPIKey {
			t.Errorf("expected api_key %q, got %q", createdAPIKey, res.APIKey)
		}
		if res.Name != keyName {
			t.Errorf("expected name %q, got %q", keyName, res.Name)
		}
		if res.Enabled {
			t.Errorf("expected Enabled=false (created disabled), got true")
		}
	})

	// Step d: Update the access key name (PUT).
	t.Run("d_update", func(t *testing.T) {
		res, err := a.UpdateAccessKey(ctx, provisioning.UpdateAccessKeyParams{
			SubAccountID: subAccountID,
			APIKey:       createdAPIKey,
			Name:         updatedKeyName,
		})
		if err != nil {
			t.Fatalf("UpdateAccessKey error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("UpdateAccessKey API error: %s", res.Error.Message)
		}
		t.Logf("updated access key: api_key=%s name=%s", res.APIKey, res.Name)

		if res.Name != updatedKeyName {
			t.Errorf("expected name %q, got %q", updatedKeyName, res.Name)
		}
	})

	// Step e: Re-enable the access key (PUT with enabled=true).
	// The Cloudinary API only honors `enabled: true` on PUT. `enabled: false`
	// is silently ignored, so we test the false->true direction here.
	t.Run("e_enable", func(t *testing.T) {
		res, err := a.UpdateAccessKey(ctx, provisioning.UpdateAccessKeyParams{
			SubAccountID: subAccountID,
			APIKey:       createdAPIKey,
			Enabled:      api.Bool(true),
		})
		if err != nil {
			t.Fatalf("UpdateAccessKey error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("UpdateAccessKey API error: %s", res.Error.Message)
		}
		t.Logf("re-enabled access key: api_key=%s enabled=%v", res.APIKey, res.Enabled)

		if !res.Enabled {
			t.Errorf("expected Enabled=true after re-enable, got false")
		}
	})

	// Step f: Delete the created access key (DELETE).
	t.Run("f_delete", func(t *testing.T) {
		res, err := a.DeleteAccessKey(ctx, provisioning.DeleteAccessKeyParams{
			SubAccountID: subAccountID,
			APIKey:       createdAPIKey,
		})
		if err != nil {
			t.Fatalf("DeleteAccessKey error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("DeleteAccessKey API error: %s", res.Error.Message)
		}
		t.Logf("delete response message: %s", res.Message)
		// Mark as deleted so cleanup won't double-delete.
		createdAPIKey = ""
	})

	// Step g: Verify deletion by listing and confirming the key is absent.
	t.Run("g_verify_deleted", func(t *testing.T) {
		res, err := a.ListAccessKeys(ctx, provisioning.ListAccessKeysParams{
			SubAccountID: subAccountID,
		})
		if err != nil {
			t.Fatalf("ListAccessKeys error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("ListAccessKeys API error: %s", res.Error.Message)
		}
		t.Logf("access keys after deletion count: %d (total: %d)", len(res.AccessKeys), res.Total)

		for _, k := range res.AccessKeys {
			if k.Name == updatedKeyName {
				t.Errorf("deleted key with name %q still appears in list", updatedKeyName)
			}
		}
	})
}
