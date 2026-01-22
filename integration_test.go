//go:build integration

/**
* Copyright 2014 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package rallyresttoolkit_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	rally "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/models"
)

// skipIfNoAPIKey skips the test if RALLY_API_KEY is not set
func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("RALLY_API_KEY") == "" {
		t.Skip("Skipping integration test: RALLY_API_KEY environment variable not set")
	}
}

// QueryResult represents the Rally API query response structure
type QueryResult struct {
	QueryResult struct {
		TotalResultCount int           `json:"TotalResultCount"`
		StartIndex       int           `json:"StartIndex"`
		PageSize         int           `json:"PageSize"`
		Results          []interface{} `json:"Results"`
	} `json:"QueryResult"`
}

// TestIntegration_NewClientFromEnv verifies that the client can be created from environment variables
func TestIntegration_NewClientFromEnv(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClientFromEnv() returned nil client")
	}

	// Verify the client has an HTTP client configured
	if client.HTTPClient() == nil {
		t.Error("client.HTTPClient() returned nil")
	}
}

// TestIntegration_QueryRequest_FetchProjects queries for projects to verify API connectivity
func TestIntegration_QueryRequest_FetchProjects(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result QueryResult
	// Query for projects - this should work for any valid API key
	// Using an empty query to fetch all accessible projects
	err = client.QueryRequest(ctx, map[string]string{}, "project", &result)
	if err != nil {
		t.Fatalf("QueryRequest for projects failed: %v", err)
	}

	t.Logf("Query returned %d projects", result.QueryResult.TotalResultCount)

	// We should have at least one project accessible
	if result.QueryResult.TotalResultCount < 1 {
		t.Error("Expected at least one project, got 0")
	}
}

// TestIntegration_QueryRequest_FetchWorkspaces queries for workspaces to verify API connectivity
func TestIntegration_QueryRequest_FetchWorkspaces(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result QueryResult
	// Query for workspaces - basic API endpoint that should work
	err = client.QueryRequest(ctx, map[string]string{}, "workspace", &result)
	if err != nil {
		t.Fatalf("QueryRequest for workspaces failed: %v", err)
	}

	t.Logf("Query returned %d workspaces", result.QueryResult.TotalResultCount)

	// We should have at least one workspace accessible
	if result.QueryResult.TotalResultCount < 1 {
		t.Error("Expected at least one workspace, got 0")
	}
}

// TestIntegration_QueryRequest_FetchDefects queries for defects with pagination
func TestIntegration_QueryRequest_FetchDefects(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var result QueryResult
	// Query for defects - common artifact type
	err = client.QueryRequest(ctx, map[string]string{}, "defect", &result)
	if err != nil {
		t.Fatalf("QueryRequest for defects failed: %v", err)
	}

	t.Logf("Query returned %d defects (TotalResultCount)", result.QueryResult.TotalResultCount)
	t.Logf("Query returned %d defects in Results array", len(result.QueryResult.Results))

	// TotalResultCount should be >= 0 (valid response)
	if result.QueryResult.TotalResultCount < 0 {
		t.Error("Expected TotalResultCount >= 0")
	}
}

// TestIntegration_QueryRequest_ContextCancellation verifies context cancellation works
func TestIntegration_QueryRequest_ContextCancellation(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result QueryResult
	err = client.QueryRequest(ctx, map[string]string{}, "project", &result)
	if err == nil {
		t.Error("Expected error due to cancelled context, got nil")
	}
}

// TestIntegration_CRUD_Defect tests the full CRUD lifecycle for a Defect
func TestIntegration_CRUD_Defect(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := rally.NewClientFromEnv()
	if err != nil {
		t.Fatalf("NewClientFromEnv() failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Track the created defect ID for cleanup
	var createdDefectID string

	// Register cleanup function to delete the defect even if test fails
	t.Cleanup(func() {
		if createdDefectID != "" {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cleanupCancel()

			defectClient := rally.NewDefect(client)
			deleteErr := defectClient.DeleteDefect(cleanupCtx, createdDefectID)
			if deleteErr != nil {
				t.Logf("Cleanup: failed to delete test defect %s: %v", createdDefectID, deleteErr)
			} else {
				t.Logf("Cleanup: successfully deleted test defect %s", createdDefectID)
			}
		}
	})

	defectClient := rally.NewDefect(client)

	// --- CREATE ---
	t.Log("Creating test defect...")
	testDefect := models.Defect{
		Name:        "Integration Test Defect - " + time.Now().Format(time.RFC3339),
		Description: "This is an automated integration test defect. It should be deleted after the test completes.",
		State:       "Submitted",
		Priority:    "Normal",
		Severity:    "Minor Problem",
	}

	createdDefect, err := defectClient.CreateDefect(ctx, testDefect)
	if err != nil {
		t.Fatalf("CreateDefect failed: %v", err)
	}

	if createdDefect.ObjectID == 0 {
		t.Fatal("Created defect has no ObjectID")
	}
	createdDefectID = fmt.Sprintf("%d", createdDefect.ObjectID)
	t.Logf("Created defect with ID: %s, FormattedID: %s", createdDefectID, createdDefect.FormattedID)

	if createdDefect.Name != testDefect.Name {
		t.Errorf("Created defect name mismatch: expected %q, got %q", testDefect.Name, createdDefect.Name)
	}

	// --- READ ---
	t.Log("Reading created defect...")
	readDefect, err := defectClient.GetDefect(ctx, createdDefectID)
	if err != nil {
		t.Fatalf("GetDefect failed: %v", err)
	}

	if readDefect.ObjectID != createdDefect.ObjectID {
		t.Errorf("Read defect ObjectID mismatch: expected %d, got %d", createdDefect.ObjectID, readDefect.ObjectID)
	}
	if readDefect.Name != createdDefect.Name {
		t.Errorf("Read defect Name mismatch: expected %q, got %q", createdDefect.Name, readDefect.Name)
	}
	t.Logf("Read defect successfully: %s (%s)", readDefect.Name, readDefect.FormattedID)

	// --- UPDATE ---
	t.Log("Updating defect...")
	updatedName := "Updated Integration Test Defect - " + time.Now().Format(time.RFC3339)
	readDefect.Name = updatedName
	readDefect.Priority = "High Attention"

	updatedDefect, err := defectClient.UpdateDefect(ctx, readDefect)
	if err != nil {
		t.Fatalf("UpdateDefect failed: %v", err)
	}

	// Verify update by reading again
	verifyDefect, err := defectClient.GetDefect(ctx, createdDefectID)
	if err != nil {
		t.Fatalf("GetDefect after update failed: %v", err)
	}

	if verifyDefect.Name != updatedName {
		t.Errorf("Updated defect name mismatch: expected %q, got %q", updatedName, verifyDefect.Name)
	}
	t.Logf("Updated defect successfully: %s (Priority: %s)", verifyDefect.Name, verifyDefect.Priority)
	_ = updatedDefect // avoid unused variable warning

	// --- DELETE ---
	t.Log("Deleting defect...")
	err = defectClient.DeleteDefect(ctx, createdDefectID)
	if err != nil {
		t.Fatalf("DeleteDefect failed: %v", err)
	}
	t.Logf("Deleted defect %s", createdDefectID)

	// Clear the ID so cleanup doesn't try to delete again
	createdDefectID = ""

	// Verify deletion by attempting to read (should fail)
	_, err = defectClient.GetDefect(ctx, fmt.Sprintf("%d", createdDefect.ObjectID))
	if err == nil {
		t.Error("Expected error when reading deleted defect, got nil")
	} else {
		t.Logf("Verified defect deletion: read attempt returned expected error: %v", err)
	}

	t.Log("CRUD test completed successfully")
}
