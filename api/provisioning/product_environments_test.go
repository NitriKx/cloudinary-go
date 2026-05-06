package provisioning_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

// TestLiveProductEnvironments exercises the Product Environments (sub-accounts)
// endpoints against the live Cloudinary Provisioning API.
func TestLiveProductEnvironments(t *testing.T) {
	a := newLiveProvisioningAPI(t)
	subAccountID := liveSubAccountID(t)
	ctx := context.Background()

	// Use a timestamp-suffixed name. Cloudinary cloud_names are globally unique
	// and stay locked even after deletion, so reusing a fixed name across runs
	// fails with "Customers is invalid".
	suffix := time.Now().Unix()
	envName := fmt.Sprintf("backmarket-test-%d", suffix)
	renamedEnvName := fmt.Sprintf("%s-renamed", envName)

	var createdEnvID string

	// Cleanup: best-effort delete in case a step fails between create and delete.
	t.Cleanup(func() {
		if createdEnvID == "" {
			return
		}
		delRes, err := a.DeleteProductEnvironment(ctx, provisioning.DeleteProductEnvironmentParams{
			SubAccountID: createdEnvID,
		})
		if err != nil {
			t.Logf("[cleanup] DeleteProductEnvironment(%s) error: %v", createdEnvID, err)
			return
		}
		if delRes.Error.Message != "" {
			t.Logf("[cleanup] DeleteProductEnvironment(%s) API error: %s", createdEnvID, delRes.Error.Message)
			return
		}
		t.Logf("[cleanup] deleted product environment %s", createdEnvID)
	})

	// -------------------------------------------------------------------------
	// Step a: List all product environments
	// -------------------------------------------------------------------------
	t.Run("a_list", func(t *testing.T) {
		res, err := a.ListProductEnvironments(ctx, provisioning.ListProductEnvironmentsParams{})
		if err != nil {
			t.Fatalf("ListProductEnvironments error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("ListProductEnvironments API error: %s", res.Error.Message)
		}

		t.Logf("total product environments returned: %d", len(res.SubAccounts))
		for _, env := range res.SubAccounts {
			t.Logf("  id=%s name=%q cloud_name=%q enabled=%v api_access_keys=%d",
				env.ID, env.Name, env.CloudName, env.Enabled, len(env.APIAccessKeys))
		}

		if res.SubAccounts == nil {
			t.Error("expected SubAccounts field to be non-nil (slice), got nil")
		}
	})

	// -------------------------------------------------------------------------
	// Step b: Get the known sub-account by ID
	// -------------------------------------------------------------------------
	t.Run("b_get_known", func(t *testing.T) {
		res, err := a.GetProductEnvironment(ctx, provisioning.GetProductEnvironmentParams{
			SubAccountID: subAccountID,
		})
		if err != nil {
			t.Fatalf("GetProductEnvironment error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("GetProductEnvironment API error: %s", res.Error.Message)
		}

		t.Logf("id=%s name=%q cloud_name=%q enabled=%v api_access_keys=%d",
			res.ID, res.Name, res.CloudName, res.Enabled, len(res.APIAccessKeys))

		if res.ID != subAccountID {
			t.Errorf("expected ID %q, got %q", subAccountID, res.ID)
		}
		if res.Name == "" {
			t.Error("expected non-empty Name")
		}
		if res.CloudName == "" {
			t.Error("expected non-empty CloudName")
		}
		if !res.Enabled {
			t.Errorf("expected Enabled=true for sub-account %s, got false", subAccountID)
		}
		if len(res.APIAccessKeys) == 0 {
			t.Errorf("expected at least one APIAccessKey for sub-account %s, got none", subAccountID)
		}
	})

	// -------------------------------------------------------------------------
	// Step c: Create a new product environment
	// -------------------------------------------------------------------------
	t.Run("c_create", func(t *testing.T) {
		res, err := a.CreateProductEnvironment(ctx, provisioning.CreateProductEnvironmentParams{
			Name:      envName,
			CloudName: envName,
		})
		if err != nil {
			t.Fatalf("CreateProductEnvironment error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("CreateProductEnvironment API error: %s", res.Error.Message)
		}
		if res.ID == "" {
			t.Fatal("CreateProductEnvironment returned empty ID")
		}
		createdEnvID = res.ID

		t.Logf("created product environment id=%s name=%q cloud_name=%q enabled=%v edge_key=%q",
			res.ID, res.Name, res.CloudName, res.Enabled, res.EdgeKey)

		if res.Name != envName {
			t.Errorf("expected Name %q, got %q", envName, res.Name)
		}
		if res.CloudName != envName {
			t.Errorf("expected CloudName %q, got %q", envName, res.CloudName)
		}
		if !res.Enabled {
			t.Errorf("expected Enabled=true on create, got false")
		}
		if res.EdgeKey == "" {
			t.Errorf("expected non-empty EdgeKey on create response")
		}
		// A newly-created sub-account always has at least one auto-generated key.
		if len(res.APIAccessKeys) == 0 {
			t.Errorf("expected at least one APIAccessKey on create, got none")
		}
	})

	if createdEnvID == "" {
		t.Skip("skipping remaining steps because product environment was not created")
	}

	// -------------------------------------------------------------------------
	// Step d: Update the product environment name
	// -------------------------------------------------------------------------
	// Note: the credentials used for tests may not have the `update` action on
	// newly-created sub-accounts. The step skips itself on a 403, so granting
	// update permissions later activates it automatically.
	t.Run("d_update", func(t *testing.T) {
		res, err := a.UpdateProductEnvironment(ctx, provisioning.UpdateProductEnvironmentParams{
			SubAccountID: createdEnvID,
			Name:         renamedEnvName,
		})
		if err != nil {
			t.Fatalf("UpdateProductEnvironment error: %v", err)
		}
		if msg := res.Error.Message; msg != "" {
			if strings.Contains(msg, `actions=["update"]`) {
				t.Skipf("update permission not granted on this account: %s", msg)
			}
			t.Fatalf("UpdateProductEnvironment API error: %s", msg)
		}

		t.Logf("updated product environment id=%s name=%q", res.ID, res.Name)

		if res.Name != renamedEnvName {
			t.Errorf("expected updated Name %q, got %q", renamedEnvName, res.Name)
		}
	})

	// -------------------------------------------------------------------------
	// Step e: Delete the created product environment
	// -------------------------------------------------------------------------
	t.Run("e_delete", func(t *testing.T) {
		res, err := a.DeleteProductEnvironment(ctx, provisioning.DeleteProductEnvironmentParams{
			SubAccountID: createdEnvID,
		})
		if err != nil {
			t.Fatalf("DeleteProductEnvironment error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("DeleteProductEnvironment API error: %s", res.Error.Message)
		}
		t.Logf("delete response message: %s", res.Message)
		// Mark as deleted so cleanup won't double-delete.
		createdEnvID = ""
	})
}
