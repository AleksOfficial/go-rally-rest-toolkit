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

// Task - struct to hold client
type Task struct {
	client *RallyClient
}

// QueryTaskResponse - struct to contain query response
type QueryTaskResponse struct {
	QueryResult struct {
		Results          []models.Task
		TotalResultCount int
	}
}

// GetTaskResponse - Struct to contain response
type GetTaskResponse struct {
	Task models.Task
}

// CreateTaskRequest - Struct to contain request
type CreateTaskRequest struct {
	Task models.Task
}

type CreateTaskResponse struct {
	CreateResult taskResult
}

type taskResult struct {
	Object models.Task
}

// OperationResponse - struct to contain response
type taskOperationResponse struct {
	OperationalResult taskResult
}

// NewTask - creates new Task
func NewTask(client *RallyClient) (de *Task) {
	return &Task{
		client: client,
	}
}

// QueryTask - abstraction for QueryRequest
func (s *Task) QueryTask(ctx context.Context, query map[string]string) (des []models.Task, err error) {
	qdes := new(QueryTaskResponse)
	err = s.client.QueryRequest(ctx, query, "task", &qdes)
	return qdes.QueryResult.Results, err
}

// GetTask - abstraction for GetRequest
func (s *Task) GetTask(ctx context.Context, objectID string) (de models.Task, err error) {
	gde := new(GetTaskResponse)
	err = s.client.GetRequest(ctx, objectID, "task", &gde)
	return gde.Task, err
}

// CreateTask - abstraction for CreateRequest
func (s *Task) CreateTask(ctx context.Context, task models.Task) (der models.Task, err error) {
	createRequest := CreateTaskRequest{
		Task: task,
	}
	ude := new(CreateTaskResponse)
	err = s.client.CreateRequest(ctx, "task", createRequest, &ude)
	der = ude.CreateResult.Object
	return der, err
}

// UpdateTask - abstraction for UpdateRequest
func (s *Task) UpdateTask(ctx context.Context, task models.Task) (taskr models.Task, err error) {
	ude := new(taskOperationResponse)
	err = s.client.UpdateRequest(ctx, strconv.Itoa(task.ObjectID), "task", task, &ude)
	taskr = ude.OperationalResult.Object
	return taskr, err
}

// DeleteTask - abstraction for DeleteRequest
func (s *Task) DeleteTask(ctx context.Context, objectID string) (err error) {
	ude := new(deOperationResponse)
	err = s.client.DeleteRequest(ctx, objectID, "task", &ude)
	return err
}
