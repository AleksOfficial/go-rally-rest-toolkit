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

func TestQueryDefect_ValidFormattedID(t *testing.T) {
	fakeFormattedID := "DE624340"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"FormattedID": "DE624340","Errors": [], "Warnings": []}]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	defectClient := NewDefect(rallyClient)
	ctx := context.Background()

	query := map[string]string{
		"FormattedID": fakeFormattedID,
	}
	results, err := defectClient.QueryDefect(ctx, query)
	if err != nil {
		t.Fatalf("QueryDefect failed unexpectedly: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results, got empty slice")
	}
	if results[0].FormattedID != fakeFormattedID {
		t.Errorf("expected FormattedID=%s, got %s", fakeFormattedID, results[0].FormattedID)
	}
}

func TestGetDefect_ValidObjectID(t *testing.T) {
	fakeObjectID := "50137325678"
	ctrlID := 50137325678
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"Defect": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	defectClient := NewDefect(rallyClient)
	ctx := context.Background()

	result, err := defectClient.GetDefect(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("GetDefect failed unexpectedly: %v", err)
	}
	if result.ObjectID != ctrlID {
		t.Errorf("expected ObjectID=%d, got %d", ctrlID, result.ObjectID)
	}
}

func TestCreateDefect_ValidRequest(t *testing.T) {
	ctrlName := "NewStory"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewStory", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	defectClient := NewDefect(rallyClient)
	ctx := context.Background()

	newDefect := models.Defect{
		Name: ctrlName,
	}
	result, err := defectClient.CreateDefect(ctx, newDefect)
	if err != nil {
		t.Fatalf("CreateDefect failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestUpdateDefect_ValidRequest(t *testing.T) {
	ctrlName := "UpdatedStoryName"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedStoryName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	defectClient := NewDefect(rallyClient)
	ctx := context.Background()

	updateDefect := models.Defect{
		Name:     ctrlName,
		ObjectID: 50137325678,
	}
	result, err := defectClient.UpdateDefect(ctx, updateDefect)
	if err != nil {
		t.Fatalf("UpdateDefect failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestDeleteDefect_ValidObjectID(t *testing.T) {
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
	defectClient := NewDefect(rallyClient)
	ctx := context.Background()

	err := defectClient.DeleteDefect(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("DeleteDefect failed unexpectedly: %v", err)
	}
}
