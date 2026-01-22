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
	"strconv"

	. "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/fakes"
	"github.com/aleksofficial/go-rally-rest-toolkit/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {

	var (
		apiKey      string
		apiURL      = "http://myRallyUrl"
		rallyClient *RallyClient
		hrclient    *Task
		ctx         = context.Background()
	)
	Describe(".QueryTask", func() {

		var (
			fakeFormattedID = "TA624340"
			fakeClient      = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"FormattedID": "TA624340","Errors": [], "Warnings": []}]}}`)},
				},
			}
		)

		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewTask(rallyClient)
		})
		Context("when called with a valid formattedID", func() {
			It("should return the requested array of Task results", func() {
				query := map[string]string{
					"FormattedID": fakeFormattedID,
				}
				hr, err := hrclient.QueryTask(ctx, query)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(hr)).ShouldNot(Equal(0))
				Ω(hr[0].FormattedID).Should(Equal(fakeFormattedID))
			})
		})

	})

	Describe(".GetTask", func() {
		var (
			fakeObjectID = "50137325678"
			ctrlID, _    = strconv.Atoi(fakeObjectID)
			fakeClient   = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"Task": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
				},
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewTask(rallyClient)
		})
		Context("when called with a valid objectID", func() {
			It("should return the Task", func() {
				hr, err := hrclient.GetTask(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.ObjectID).Should(Equal(ctrlID))
			})
		})

	})

	Describe(".CreateTask", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,

					Body: &fakes.FakeResponseBody{bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewTask", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName   = "NewTask"
			newHrModel = models.Task{
				Name: ctrlName,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"
			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewTask(rallyClient)
		})
		Context("when called with a valid create request object", func() {
			It("should return the Task object created", func() {
				hr, err := hrclient.CreateTask(ctx, newHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})
	})

	Describe(".UpdateTask", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedTaskName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName      = "UpdatedTaskName"
			updateHrModel = models.Task{
				Name:     ctrlName,
				ObjectID: 50137325678,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewTask(rallyClient)
		})

		Context("when called with a valid update request object", func() {
			It("should return the Task object updated", func() {
				hr, err := hrclient.UpdateTask(ctx, updateHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})

	})

	Describe(".DeleteTask", func() {
		var (
			fakeObjectID = "50137325678"
			fakeClient   = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationalResult": {"Errors": [], "Warnings": []}}`)},
				},
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewTask(rallyClient)
		})
		Context("when called with a valid delete objectID", func() {
			It("should return the correct operationalresponse struct", func() {
				err := hrclient.DeleteTask(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
