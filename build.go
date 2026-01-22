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

// Build - struct to hold client
type Build struct {
	client *RallyClient
}

// QueryBuildResponse - struct to contain query response
type QueryBuildResponse struct {
	QueryResult struct {
		Results          []models.Build
		TotalResultCount int
	}
}

// GetBuildResponse - Struct to contain response
type GetBuildResponse struct {
	Build models.Build
}

// CreateBuildRequest - Struct to contain request
type CreateBuildRequest struct {
	Build models.Build
}

type CreateBuildResponse struct {
	CreateResult buildResult
}

type buildResult struct {
	Object models.Build
}

// OperationResponse - struct to contain response
type buildOperationResponse struct {
	OperationalResult buildResult
}

// NewBuild - creates new Build
func NewBuild(client *RallyClient) (de *Build) {
	return &Build{
		client: client,
	}
}

// QueryBuild - abstraction for QueryRequest
func (s *Build) QueryBuild(ctx context.Context, query map[string]string) (des []models.Build, err error) {
	qdes := new(QueryBuildResponse)
	err = s.client.QueryRequest(ctx, query, "build", &qdes)
	return qdes.QueryResult.Results, err
}

// GetBuild - abstraction for GetRequest
func (s *Build) GetBuild(ctx context.Context, objectID string) (de models.Build, err error) {
	gde := new(GetBuildResponse)
	err = s.client.GetRequest(ctx, objectID, "build", &gde)
	return gde.Build, err
}

// CreateBuild - abstraction for CreateRequest
func (s *Build) CreateBuild(ctx context.Context, build models.Build) (der models.Build, err error) {
	createRequest := CreateBuildRequest{
		Build: build,
	}
	ude := new(CreateBuildResponse)
	err = s.client.CreateRequest(ctx, "build", createRequest, &ude)
	der = ude.CreateResult.Object
	return der, err
}

// UpdateBuild - abstraction for UpdateRequest
func (s *Build) UpdateBuild(ctx context.Context, build models.Build) (buildr models.Build, err error) {
	ude := new(buildOperationResponse)
	err = s.client.UpdateRequest(ctx, strconv.Itoa(build.ObjectID), "build", build, &ude)
	buildr = ude.OperationalResult.Object
	return buildr, err
}

// DeleteBuild - abstraction for DeleteRequest
func (s *Build) DeleteBuild(ctx context.Context, objectID string) (err error) {
	ude := new(deOperationResponse)
	err = s.client.DeleteRequest(ctx, objectID, "build", &ude)
	return err
}
