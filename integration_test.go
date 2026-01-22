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
	"os"
	"testing"
	"time"

	rally "github.com/aleksofficial/go-rally-rest-toolkit"
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
