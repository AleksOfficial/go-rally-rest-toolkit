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

// HierarchicalRequirement - struct to hold client
type HierarchicalRequirement struct {
	client *RallyClient
}

// QueryHierarchicalRequirementResponse - struct to contain query response
type QueryHierarchicalRequirementResponse struct {
	QueryResult struct {
		Results          []models.HierarchicalRequirement
		TotalResultCount int
	}
}

// GetHierarchicalRequirementResponse - Struct to contain response
type GetHierarchicalRequirementResponse struct {
	HierarchicalRequirement models.HierarchicalRequirement
}

// HierarchicalRequirementRequest - Struct to contain request for creates and updates.
type HierarchicalRequirementRequest struct {
	HierarchicalRequirement models.HierarchicalRequirement
}

// CreateHierarchicalRequirementResponse - reponse struct
type CreateHierarchicalRequirementResponse struct {
	CreateResult HrResult
}

// HrResult - struct to contain result from request.
type HrResult struct {
	Object models.HierarchicalRequirement
}

// OperationResponse - struct to contain response
type OperationResponse struct {
	OperationalResult HrResult
}

// NewHierarchicalRequirement - creates new HierarchicalRequirement
func NewHierarchicalRequirement(client *RallyClient) (hr *HierarchicalRequirement) {
	return &HierarchicalRequirement{
		client: client,
	}
}

// QueryHierarchicalRequirement - abstraction for QueryRequest
func (s *HierarchicalRequirement) QueryHierarchicalRequirement(query map[string]string) (hrs []models.HierarchicalRequirement, err error) {
	qhrs := new(QueryHierarchicalRequirementResponse)
	err = s.client.QueryRequest(query, "HierarchicalRequirement", &qhrs)
	return qhrs.QueryResult.Results, err
}

// GetHierarchicalRequirement - abstraction for GetRequest
func (s *HierarchicalRequirement) GetHierarchicalRequirement(objectID string) (hr models.HierarchicalRequirement, err error) {
	ghr := new(GetHierarchicalRequirementResponse)
	err = s.client.GetRequest(objectID, "HierarchicalRequirement", &ghr)
	return ghr.HierarchicalRequirement, err
}

// CreateHierarchicalRequirement - abstraction for CreateRequest
func (s *HierarchicalRequirement) CreateHierarchicalRequirement(hr models.HierarchicalRequirement) (hrr models.HierarchicalRequirement, err error) {
	createRequest := HierarchicalRequirementRequest{
		HierarchicalRequirement: hr,
	}
	uhr := new(CreateHierarchicalRequirementResponse)
	err = s.client.CreateRequest("HierarchicalRequirement", createRequest, &uhr)
	hrr = uhr.CreateResult.Object
	return hrr, err
}

// UpdateHierarchicalRequirement - abstraction for UpdateRequest
func (s *HierarchicalRequirement) UpdateHierarchicalRequirement(hr models.HierarchicalRequirement) (hrr models.HierarchicalRequirement, err error) {
	updateRequest := HierarchicalRequirementRequest{
		HierarchicalRequirement: hr,
	}
	uhr := new(OperationResponse)
	err = s.client.UpdateRequest(strconv.Itoa(hr.ObjectID), "HierarchicalRequirement", updateRequest, &uhr)
	hrr = uhr.OperationalResult.Object
	return hrr, err
}

// DeleteHierarchicalRequirement - abstraction for DeleteRequest
func (s *HierarchicalRequirement) DeleteHierarchicalRequirement(objectID string) (err error) {
	uhr := new(OperationResponse)
	err = s.client.DeleteRequest(objectID, "HierarchicalRequirement", &uhr)
	return err
}
