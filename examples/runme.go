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

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/models"
)

type QueryHR struct {
	QueryResult struct {
		Results          []models.HierarchicalRequirement
		TotalResultCount int
	}
}

type GetHR struct {
	HierarchicalRequirement models.HierarchicalRequirement
}

type CreateDefect struct {
	Defect models.Defect
}

type CreateResponse struct {
	CreateResult map[string]interface{}
}

type OperationResponse struct {
	OperationalResult map[string]interface{}
}

func main() {

	rallyURL := "https://rally1.rallydev.com/slm/webservice/v2.0"
	apiKey := os.Getenv("API_KEY")
	clientDoer := &http.Client{}

	rallyClient := rallyresttoolkit.New(apiKey, rallyURL, clientDoer)

	query := map[string]string{
		"FormattedID": "US624340",
	}

	output := new(QueryHR)
	err := rallyClient.QueryRequest(query, "hierarchicalrequirement", &output)

	if err == nil {
		fmt.Printf("output: %v \n", output)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
	}

	moarOutput := new(GetHR)
	err = rallyClient.GetRequest("29227987232", "hierarchicalrequirement", &moarOutput)

	if err == nil {
		fmt.Printf("MOAR output: %v \n", moarOutput)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
	}

	ref := string(moarOutput.HierarchicalRequirement.Ref)
	fmt.Printf("Ref output: %v \n", ref)
	storyRef := &models.Reference{
		Ref:   ref,
		Count: 1,
	}

	fmt.Printf("StoryRef output: %v \n", storyRef)

	newDefect := models.Defect{
		Name:        "TestDefect",
		Requirement: storyRef,
		Priority:    "Urgent",
		Severity:    "Urgent",
	}

	defectRequest := CreateDefect{
		Defect: newDefect,
	}

	createResponse := new(CreateResponse)

	err = rallyClient.CreateRequest("defect", defectRequest, &createResponse)

	object := createResponse.CreateResult["Object"].(map[string]interface{})
	resultID := object["ObjectID"].(float64)

	objectID := fmt.Sprintf("%v", int64(resultID))
	fmt.Printf("Create Response: %v", createResponse)
	updateDefect := CreateDefect{
		Defect: models.Defect{
			Resolution: "Not a Defect",
			State:      "Fixed",
		},
	}

	updateResponse := new(OperationResponse)
	_ = rallyClient.UpdateRequest(objectID, "defect", updateDefect, &updateResponse)

	_ = rallyClient.DeleteRequest(objectID, "defect", &updateResponse)
}
