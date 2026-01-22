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

// Defect - struct to hold client
type Defect struct {
	client *RallyClient
}

// QueryDefectResponse - struct to contain query response
type QueryDefectResponse struct {
	QueryResult struct {
		Results          []models.Defect
		TotalResultCount int
	}
}

// GetDefectResponse - Struct to contain response
type GetDefectResponse struct {
	Defect models.Defect
}

// CreateDefectRequest - Struct to contain request
type CreateDefectRequest struct {
	Defect models.Defect
}

type CreateDefectResponse struct {
	CreateResult deResult
}

type deResult struct {
	Object models.Defect
}

// OperationResponse - struct to contain response
type deOperationResponse struct {
	OperationalResult deResult
}

// NewDefect - creates new Defect
func NewDefect(client *RallyClient) (de *Defect) {
	return &Defect{
		client: client,
	}
}

// QueryDefect - abstraction for QueryRequest
func (s *Defect) QueryDefect(ctx context.Context, query map[string]string) (des []models.Defect, err error) {
	qdes := new(QueryDefectResponse)
	err = s.client.QueryRequest(ctx, query, "defect", &qdes)
	return qdes.QueryResult.Results, err
}

// GetDefect - abstraction for GetRequest
func (s *Defect) GetDefect(ctx context.Context, objectID string) (de models.Defect, err error) {
	gde := new(GetDefectResponse)
	err = s.client.GetRequest(ctx, objectID, "defect", &gde)
	return gde.Defect, err
}

// CreateDefect - abstraction for CreateRequest
func (s *Defect) CreateDefect(ctx context.Context, de models.Defect) (der models.Defect, err error) {
	createRequest := CreateDefectRequest{
		Defect: de,
	}
	ude := new(CreateDefectResponse)
	err = s.client.CreateRequest(ctx, "defect", createRequest, &ude)
	der = ude.CreateResult.Object
	return der, err
}

// UpdateDefect - abstraction for UpdateRequest
func (s *Defect) UpdateDefect(ctx context.Context, de models.Defect) (der models.Defect, err error) {
	ude := new(deOperationResponse)
	err = s.client.UpdateRequest(ctx, strconv.Itoa(de.ObjectID), "Defect", de, &ude)
	der = ude.OperationalResult.Object
	return der, err
}

// DeleteDefect - abstraction for DeleteRequest
func (s *Defect) DeleteDefect(ctx context.Context, objectID string) (err error) {
	ude := new(deOperationResponse)
	err = s.client.DeleteRequest(ctx, objectID, "defect", &ude)
	return err
}
