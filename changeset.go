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

package rallyresttoolkit

import (
	"strconv"

	"github.com/aleksofficial/go-rally-rest-toolkit/models"
)

// Changeset - struct to hold client
type Changeset struct {
	client *RallyClient
}

// QueryChangesetResponse - struct to contain query response
type QueryChangesetResponse struct {
	QueryResult struct {
		Results          []models.Changeset
		TotalResultCount int
	}
}

// GetChangesetResponse - Struct to contain response
type GetChangesetResponse struct {
	Changeset models.Changeset
}

// CreateChangesetRequest - Struct to contain request
type CreateChangesetRequest struct {
	Changeset models.Changeset
}

type CreateChangesetResponse struct {
	CreateResult changesetResult
}

type changesetResult struct {
	Object models.Changeset
}

// OperationResponse - struct to contain response
type changesetOperationResponse struct {
	OperationalResult changesetResult
}

// NewChangeset - creates new Changeset
func NewChangeset(client *RallyClient) (cs *Changeset) {
	return &Changeset{
		client: client,
	}
}

// QueryChangeset - abstraction for QueryRequest
func (s *Changeset) QueryChangeset(query map[string]string) (des []models.Changeset, err error) {
	qdes := new(QueryChangesetResponse)
	err = s.client.QueryRequest(query, "changeset", &qdes)
	return qdes.QueryResult.Results, err
}

// GetChangeset - abstraction for GetRequest
func (s *Changeset) GetChangeset(objectID string) (de models.Changeset, err error) {
	gde := new(GetChangesetResponse)
	err = s.client.GetRequest(objectID, "changeset", &gde)
	return gde.Changeset, err
}

// CreateChangeset - abstraction for CreateRequest
func (s *Changeset) CreateChangeset(changeset models.Changeset) (der models.Changeset, err error) {
	createRequest := CreateChangesetRequest{
		Changeset: changeset,
	}
	ude := new(CreateChangesetResponse)
	err = s.client.CreateRequest("changeset", createRequest, &ude)
	der = ude.CreateResult.Object
	return der, err
}

// UpdateChangeset - abstraction for UpdateRequest
func (s *Changeset) UpdateChangeset(changeset models.Changeset) (changesetr models.Changeset, err error) {
	ude := new(changesetOperationResponse)
	err = s.client.UpdateRequest(strconv.Itoa(changeset.ObjectID), "changeset", changeset, &ude)
	changesetr = ude.OperationalResult.Object
	return changesetr, err
}

// DeleteChangeset - abstraction for DeleteRequest
func (s *Changeset) DeleteChangeset(objectID string) (err error) {
	ude := new(deOperationResponse)
	err = s.client.DeleteRequest(objectID, "changeset", &ude)
	return err
}
