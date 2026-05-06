package provisioning_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/provisioning"
)

const (
	livePolicyName      = "go-sdk-live-test-policy"
	livePolicyStatement = "permit(principal, action, resource);"
	livePolicyScopeType = "account"
)

func TestLiveCustomPolicies(t *testing.T) {
	ctx := context.Background()
	a := newLiveProvisioningAPI(t)

	var createdPolicyID string

	// Always clean up even if a step fails mid-test.
	t.Cleanup(func() {
		if createdPolicyID == "" {
			return
		}
		delRes, err := a.DeleteCustomPolicy(ctx, provisioning.DeleteCustomPolicyParams{
			PolicyID: createdPolicyID,
		})
		if err != nil {
			t.Logf("[cleanup] DeleteCustomPolicy error: %v", err)
			return
		}
		if delRes.Error.Message != "" {
			t.Logf("[cleanup] DeleteCustomPolicy API error: %s", delRes.Error.Message)
			return
		}
		t.Logf("[cleanup] deleted policy %s", createdPolicyID)
	})

	// -------------------------------------------------------------------------
	// Step a: List existing custom policies
	// -------------------------------------------------------------------------
	t.Run("a_ListBefore", func(t *testing.T) {
		res, err := a.ListCustomPolicies(ctx, provisioning.ListCustomPoliciesParams{})
		if err != nil {
			t.Fatalf("ListCustomPolicies error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("ListCustomPolicies API error: %s", res.Error.Message)
		}
		t.Logf("existing policies count: %d", len(res.Data))
		for _, p := range res.Data {
			t.Logf("  - id=%s name=%q enabled=%v", p.ID, p.Name, p.Enabled)
		}
	})

	// -------------------------------------------------------------------------
	// Step b: Create a new custom policy
	// -------------------------------------------------------------------------
	t.Run("b_Create", func(t *testing.T) {
		enabled := true
		res, err := a.CreateCustomPolicy(ctx, provisioning.CreateCustomPolicyParams{
			Name:            livePolicyName,
			PolicyStatement: livePolicyStatement,
			ScopeType:       livePolicyScopeType,
			Enabled:         api.Bool(enabled),
		})
		if err != nil {
			t.Fatalf("CreateCustomPolicy error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("CreateCustomPolicy API error: %s", res.Error.Message)
		}
		if res.Data.ID == "" {
			t.Fatal("CreateCustomPolicy returned empty ID")
		}
		createdPolicyID = res.Data.ID
		t.Logf("created policy id=%s name=%q enabled=%v", res.Data.ID, res.Data.Name, res.Data.Enabled)
	})

	if createdPolicyID == "" {
		t.Skip("skipping remaining steps: policy was not created")
	}

	// -------------------------------------------------------------------------
	// Step c: Get the created custom policy by ID
	// -------------------------------------------------------------------------
	t.Run("c_Get", func(t *testing.T) {
		res, err := a.GetCustomPolicy(ctx, provisioning.GetCustomPolicyParams{
			PolicyID: createdPolicyID,
		})
		if err != nil {
			t.Fatalf("GetCustomPolicy error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("GetCustomPolicy API error: %s", res.Error.Message)
		}
		if res.Data.ID != createdPolicyID {
			t.Errorf("expected ID %q, got %q", createdPolicyID, res.Data.ID)
		}
		if res.Data.Name != livePolicyName {
			t.Errorf("expected Name %q, got %q", livePolicyName, res.Data.Name)
		}
		t.Logf("fetched policy id=%s name=%q policy_statement=%s", res.Data.ID, res.Data.Name, res.Data.PolicyStatement)
	})

	// -------------------------------------------------------------------------
	// Step d: Update the custom policy name (PUT requires all non-omitempty fields)
	// -------------------------------------------------------------------------
	updatedName := fmt.Sprintf("%s-updated", livePolicyName)
	t.Run("d_Update", func(t *testing.T) {
		res, err := a.UpdateCustomPolicy(ctx, provisioning.UpdateCustomPolicyParams{
			PolicyID:        createdPolicyID,
			Name:            updatedName,
			PolicyStatement: livePolicyStatement,
			ScopeType:       livePolicyScopeType,
			Enabled:         api.Bool(true),
		})
		if err != nil {
			t.Fatalf("UpdateCustomPolicy error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("UpdateCustomPolicy API error: %s", res.Error.Message)
		}
		if res.Data.Name != updatedName {
			t.Errorf("expected updated Name %q, got %q", updatedName, res.Data.Name)
		}
		t.Logf("updated policy id=%s name=%q", res.Data.ID, res.Data.Name)
	})

	// -------------------------------------------------------------------------
	// Step e: Delete the created custom policy
	// -------------------------------------------------------------------------
	t.Run("e_Delete", func(t *testing.T) {
		res, err := a.DeleteCustomPolicy(ctx, provisioning.DeleteCustomPolicyParams{
			PolicyID: createdPolicyID,
		})
		if err != nil {
			t.Fatalf("DeleteCustomPolicy error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("DeleteCustomPolicy API error: %s", res.Error.Message)
		}
		t.Logf("deleted policy %s", createdPolicyID)
		// Mark as already deleted so t.Cleanup skips the redundant call.
		createdPolicyID = ""
	})

	// -------------------------------------------------------------------------
	// Step f: Verify it's gone by listing again
	// -------------------------------------------------------------------------
	t.Run("f_ListAfter", func(t *testing.T) {
		res, err := a.ListCustomPolicies(ctx, provisioning.ListCustomPoliciesParams{})
		if err != nil {
			t.Fatalf("ListCustomPolicies error: %v", err)
		}
		if res.Error.Message != "" {
			t.Fatalf("ListCustomPolicies API error: %s", res.Error.Message)
		}
		t.Logf("policies after deletion count: %d", len(res.Data))
		for _, p := range res.Data {
			if p.Name == livePolicyName || p.Name == updatedName {
				t.Errorf("policy %q (id=%s) still present after deletion", p.Name, p.ID)
			}
		}
		t.Logf("verification passed: test policy is no longer listed")
	})
}
