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

func TestQueryBuildDefinition_ValidName(t *testing.T) {
	fakeName := "concourse-1"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Name": "concourse-1","Errors": [], "Warnings": []}]}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	buildDefClient := NewBuildDefinition(rallyClient)
	ctx := context.Background()

	query := map[string]string{
		"Name": fakeName,
	}
	results, err := buildDefClient.QueryBuildDefinition(ctx, query)
	if err != nil {
		t.Fatalf("QueryBuildDefinition failed unexpectedly: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results, got empty slice")
	}
	if results[0].Name != fakeName {
		t.Errorf("expected Name=%s, got %s", fakeName, results[0].Name)
	}
}

func TestGetBuildDefinition_ValidObjectID(t *testing.T) {
	fakeObjectID := "50137325678"
	ctrlID := 50137325678
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"BuildDefinition": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	buildDefClient := NewBuildDefinition(rallyClient)
	ctx := context.Background()

	result, err := buildDefClient.GetBuildDefinition(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("GetBuildDefinition failed unexpectedly: %v", err)
	}
	if result.ObjectID != ctrlID {
		t.Errorf("expected ObjectID=%d, got %d", ctrlID, result.ObjectID)
	}
}

func TestCreateBuildDefinition_ValidRequest(t *testing.T) {
	ctrlName := "NewBuildDefinition"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewBuildDefinition", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	buildDefClient := NewBuildDefinition(rallyClient)
	ctx := context.Background()

	newBuildDef := models.BuildDefinition{
		Name: ctrlName,
	}
	result, err := buildDefClient.CreateBuildDefinition(ctx, newBuildDef)
	if err != nil {
		t.Fatalf("CreateBuildDefinition failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestUpdateBuildDefinition_ValidRequest(t *testing.T) {
	ctrlName := "UpdatedBuildDefinitionName"
	fakeClient := &fakes.FakeHTTPClient{
		FakeResponse: &http.Response{
			StatusCode: http.StatusOK,
			Body:       &fakes.FakeResponseBody{Reader: bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedBuildDefinitionName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
		},
	}

	apiKey := "abcdef"
	apiURL := "http://myRallyUrl"
	rallyClient := New(apiKey, apiURL, fakeClient)
	buildDefClient := NewBuildDefinition(rallyClient)
	ctx := context.Background()

	updateBuildDef := models.BuildDefinition{
		Name:     ctrlName,
		ObjectID: 50137325678,
	}
	result, err := buildDefClient.UpdateBuildDefinition(ctx, updateBuildDef)
	if err != nil {
		t.Fatalf("UpdateBuildDefinition failed unexpectedly: %v", err)
	}
	if result.Name != ctrlName {
		t.Errorf("expected Name=%s, got %s", ctrlName, result.Name)
	}
}

func TestDeleteBuildDefinition_ValidObjectID(t *testing.T) {
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
	buildDefClient := NewBuildDefinition(rallyClient)
	ctx := context.Background()

	err := buildDefClient.DeleteBuildDefinition(ctx, fakeObjectID)
	if err != nil {
		t.Fatalf("DeleteBuildDefinition failed unexpectedly: %v", err)
	}
}
