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
	"net/http"
	"testing"

	. "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/fakes"
	"github.com/aleksofficial/go-rally-rest-toolkit/models"
)

func TestQueryChangeset_ValidMessage(t *testing.T) {
	fakeMessage := "concourse-1"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Message": "concourse-1","Errors": [], "Warnings": []}]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	changesetClient := NewChangeset(rallyClient)
	ctx := context.Background()

	query := map[string]string{
		"Message": fakeMessage,
	}
	results, err := changesetClient.QueryChangeset(ctx, query)
	if err != nil {
		t.Fatalf("QueryChangeset failed unexpectedly: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results, got empty slice")
	}
	if results[0].Message != fakeMessage {
		t.Errorf("expected Message=%s, got %s", fakeMessage, results[0].Message)
	}
}

func TestGetChangeset_ValidObjectID(t *testing.T) {
	fakeObjectID := "50137325678"
	ctrlID := 50137325678
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"Changeset": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	changesetClient := NewChangeset(rallyClient)
	ctx := context.Background()

	result, err := changesetClient.GetChangeset(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("GetChangeset failed unexpectedly: %v", err)
	}
	if result.ObjectID != ctrlID {
		t.Errorf("expected ObjectID=%d, got %d", ctrlID, result.ObjectID)
	}
}

func TestCreateChangeset_ValidRequest(t *testing.T) {
	ctrlName := "NewChangeset"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewChangeset", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	changesetClient := NewChangeset(rallyClient)
	ctx := context.Background()

	newChangeset := models.Changeset{
		Name: ctrlName,
	}
	result, err := changesetClient.CreateChangeset(ctx, newChangeset)
	if err != nil {
		t.Fatalf("CreateChangeset failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestUpdateChangeset_ValidRequest(t *testing.T) {
	ctrlName := "UpdatedChangesetName"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedChangesetName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	changesetClient := NewChangeset(rallyClient)
	ctx := context.Background()

	updateChangeset := models.Changeset{
		Name:     ctrlName,
		ObjectID: 50137325678,
	}
	result, err := changesetClient.UpdateChangeset(ctx, updateChangeset)
	if err != nil {
		t.Fatalf("UpdateChangeset failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestDeleteChangeset_ValidObjectID(t *testing.T) {
	fakeObjectID := "50137325678"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationalResult": {"Errors": [], "Warnings": []}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	changesetClient := NewChangeset(rallyClient)
	ctx := context.Background()

	err := changesetClient.DeleteChangeset(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("DeleteChangeset failed unexpectedly: %v", err)
	}
}
