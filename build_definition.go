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
	"context"
	"strconv"

	"github.com/aleksofficial/go-rally-rest-toolkit/models"
)

// BuildDefinition - struct to hold client
type BuildDefinition struct {
	client *RallyClient
}

// QueryBuildDefinitionResponse - struct to contain query response
type QueryBuildDefinitionResponse struct {
	QueryResult struct {
		Results          []models.BuildDefinition
		TotalResultCount int
	}
}

// GetBuildDefinitionResponse - Struct to contain response
type GetBuildDefinitionResponse struct {
	BuildDefinition models.BuildDefinition
}

// CreateBuildDefinitionRequest - Struct to contain request
type CreateBuildDefinitionRequest struct {
	BuildDefinition models.BuildDefinition
}

type CreateBuildDefinitionResponse struct {
	CreateResult buildDefinitionResult
}

type buildDefinitionResult struct {
	Object models.BuildDefinition
}

// OperationResponse - struct to contain response
type buildDefinitionOperationResponse struct {
	OperationalResult buildDefinitionResult
}

// NewBuildDefinition - creates new BuildDefinition
func NewBuildDefinition(client *RallyClient) (de *BuildDefinition) {
	return &BuildDefinition{
		client: client,
	}
}

// QueryBuildDefinition - abstraction for QueryRequest
func (s *BuildDefinition) QueryBuildDefinition(ctx context.Context, query map[string]string) (des []models.BuildDefinition, err error) {
	qdes := new(QueryBuildDefinitionResponse)
	err = s.client.QueryRequest(ctx, query, "buildDefinition", &qdes)
	return qdes.QueryResult.Results, err
}

// GetBuildDefinition - abstraction for GetRequest
func (s *BuildDefinition) GetBuildDefinition(ctx context.Context, objectID string) (de models.BuildDefinition, err error) {
	gde := new(GetBuildDefinitionResponse)
	err = s.client.GetRequest(ctx, objectID, "buildDefinition", &gde)
	return gde.BuildDefinition, err
}

// CreateBuildDefinition - abstraction for CreateRequest
func (s *BuildDefinition) CreateBuildDefinition(ctx context.Context, buildDefinition models.BuildDefinition) (der models.BuildDefinition, err error) {
	createRequest := CreateBuildDefinitionRequest{
		BuildDefinition: buildDefinition,
	}
	ude := new(CreateBuildDefinitionResponse)
	err = s.client.CreateRequest(ctx, "buildDefinition", createRequest, &ude)
	der = ude.CreateResult.Object
	return der, err
}

// UpdateBuildDefinition - abstraction for UpdateRequest
func (s *BuildDefinition) UpdateBuildDefinition(ctx context.Context, buildDefinition models.BuildDefinition) (buildDefinitionr models.BuildDefinition, err error) {
	ude := new(buildDefinitionOperationResponse)
	err = s.client.UpdateRequest(ctx, strconv.Itoa(buildDefinition.ObjectID), "buildDefinition", buildDefinition, &ude)
	buildDefinitionr = ude.OperationalResult.Object
	return buildDefinitionr, err
}

// DeleteBuildDefinition - abstraction for DeleteRequest
func (s *BuildDefinition) DeleteBuildDefinition(ctx context.Context, objectID string) (err error) {
	ude := new(deOperationResponse)
	err = s.client.DeleteRequest(ctx, objectID, "buildDefinition", &ude)
	return err
}
