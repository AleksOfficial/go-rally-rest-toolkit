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
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	. "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/fakes"
)

func TestQueryRequest_ValidQueryWithValidAPIKey(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err != nil {
		t.Fatalf("QueryRequest failed unexpectedly: %v", err)
	}
	if fakeOutput.QueryResult.TotalResultCount != 1 {
		t.Errorf("expected TotalResultCount=1, got %d", fakeOutput.QueryResult.TotalResultCount)
	}
}

func TestQueryRequest_HTTPError(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationResult": {"Errors": ["Server error"]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err == nil {
		t.Error("expected error for 500 status code, got nil")
	}
}

func TestGetRequest_ValidGetWithValidAPIKey(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	err := rallyClient.GetRequest(ctx, "50137325678", "hierarchicalrequirement", &fakeOutput)
	if err != nil {
		t.Fatalf("GetRequest failed unexpectedly: %v", err)
	}
	if fakeOutput.QueryResult.TotalResultCount != 1 {
		t.Errorf("expected TotalResultCount=1, got %d", fakeOutput.QueryResult.TotalResultCount)
	}
}

func TestCreateRequest_ValidCreateWithValidAPIKey(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"CreateResult": { "FakeObject": {"Field1": "demostring"} }}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeCreateRequest := &fakes.FakeCreateRequest{
		FakeItem: fakes.FakeItem{
			Field1: "demostring",
		},
	}
	fakeOutput := new(fakes.FakeCreateResponse)

	err := rallyClient.CreateRequest(ctx, "hierarchicalrequirement", fakeCreateRequest, &fakeOutput)
	if err != nil {
		t.Fatalf("CreateRequest failed unexpectedly: %v", err)
	}
	if fakeOutput.CreateResult.FakeObject["Field1"] != "demostring" {
		t.Errorf("expected Field1='demostring', got %v", fakeOutput.CreateResult.FakeObject["Field1"])
	}
}

func TestUpdateRequest_ValidUpdateWithValidAPIKey(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationResult": { "FakeObject": {"Field1": "demostring"} }}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeUpdateRequest := &fakes.FakeCreateRequest{
		FakeItem: fakes.FakeItem{
			Field1: "demostring",
		},
	}
	fakeOutput := new(fakes.FakeUpdateResponse)

	err := rallyClient.UpdateRequest(ctx, "12345", "hierarchicalrequirement", fakeUpdateRequest, &fakeOutput)
	if err != nil {
		t.Fatalf("UpdateRequest failed unexpectedly: %v", err)
	}
	if fakeOutput.OperationResult.FakeObject["Field1"] != "demostring" {
		t.Errorf("expected Field1='demostring', got %v", fakeOutput.OperationResult.FakeObject["Field1"])
	}
}

func TestDeleteRequest_ValidDeleteWithValidAPIKey(t *testing.T) {
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationResult": { "Errors": [] }}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	ctx := context.Background()

	fakeOutput := new(fakes.FakeUpdateResponse)

	err := rallyClient.DeleteRequest(ctx, "12345", "hierarchicalrequirement", &fakeOutput)
	if err != nil {
		t.Fatalf("DeleteRequest failed unexpectedly: %v", err)
	}
}

func TestQueryRequest_RetryOn5xxSuccess(t *testing.T) {
	// First call returns 500, second call returns 200 (success)
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponses: []*http.Response{
			{
				StatusCode: http.StatusInternalServerError,
				Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationResult": {"Errors": ["Server error"]}}`)},
			},
			{
				StatusCode: http.StatusOK,
				Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
			},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	// Configure retry with minimal delay for faster tests
	rallyClient.SetConfig(&Config{
		MaxRetries: 3,
		RetryDelay: 1, // 1ms for fast tests
	})
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err != nil {
		t.Fatalf("QueryRequest should have succeeded after retry: %v", err)
	}
	if fakeClient.CallCount != 2 {
		t.Errorf("expected 2 calls (1 failure + 1 success), got %d", fakeClient.CallCount)
	}
	if fakeOutput.QueryResult.TotalResultCount != 1 {
		t.Errorf("expected TotalResultCount=1, got %d", fakeOutput.QueryResult.TotalResultCount)
	}
}

func TestQueryRequest_NoRetryOn4xx(t *testing.T) {
	// 400 error should not be retried
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponses: []*http.Response{
			{
				StatusCode: http.StatusBadRequest,
				Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationResult": {"Errors": ["Bad request"]}}`)},
			},
			{
				StatusCode: http.StatusOK,
				Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1}}`)},
			},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	rallyClient.SetConfig(&Config{
		MaxRetries: 3,
		RetryDelay: 1,
	})
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err == nil {
		t.Fatal("QueryRequest should have failed on 400 error")
	}
	// Should only be called once - no retry on 4xx
	if fakeClient.CallCount != 1 {
		t.Errorf("expected 1 call (no retry on 4xx), got %d", fakeClient.CallCount)
	}
}

func TestQueryRequest_RetryOnTransientError(t *testing.T) {
	// First call returns timeout error, second call succeeds
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponses: []*http.Response{
			nil,
			{
				StatusCode: http.StatusOK,
				Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
			},
		},
		FakeErrors: []error{
			errors.New("connection timeout"),
			nil,
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	rallyClient.SetConfig(&Config{
		MaxRetries: 3,
		RetryDelay: 1,
	})
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err != nil {
		t.Fatalf("QueryRequest should have succeeded after retry: %v", err)
	}
	if fakeClient.CallCount != 2 {
		t.Errorf("expected 2 calls (1 failure + 1 success), got %d", fakeClient.CallCount)
	}
}

func TestQueryRequest_MaxRetriesExceeded(t *testing.T) {
	// All calls return 500 - should fail after max retries
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponses: []*http.Response{
			{StatusCode: http.StatusInternalServerError, Body: &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{}`)}},
			{StatusCode: http.StatusInternalServerError, Body: &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{}`)}},
			{StatusCode: http.StatusInternalServerError, Body: &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{}`)}},
			{StatusCode: http.StatusInternalServerError, Body: &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{}`)}},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	rallyClient.SetConfig(&Config{
		MaxRetries: 3, // 1 initial + 3 retries = 4 total attempts
		RetryDelay: 1,
	})
	ctx := context.Background()

	fakeOutput := new(fakes.FakeOutput)
	query := map[string]string{
		"FormattedID": "US624340",
	}

	err := rallyClient.QueryRequest(ctx, query, "hierarchicalrequirement", &fakeOutput)
	if err == nil {
		t.Fatal("QueryRequest should have failed after max retries")
	}
	// Should be called 4 times: 1 initial + 3 retries
	if fakeClient.CallCount != 4 {
		t.Errorf("expected 4 calls (1 initial + 3 retries), got %d", fakeClient.CallCount)
	}
}
