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

var _ = Describe("HierarchichalRequirement", func() {

	var (
		apiKey      string
		apiURL      = "http://myRallyUrl"
		rallyClient *RallyClient
		hrclient    *HierarchicalRequirement
		ctx         = context.Background()
	)
	Describe(".QueryHierarchicalRequirement", func() {

		var (
			fakeFormattedID = "US624340"
			fakeClient      = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"FormattedID": "US624340","Errors": [], "Warnings": []}]}}`)},
				},
			}
		)

		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewHierarchicalRequirement(rallyClient)
		})
		Context("when called with a valid formattedID", func() {
			It("should return the requested array of hierarchichal requirement results", func() {
				query := map[string]string{
					"FormattedID": fakeFormattedID,
				}
				hr, err := hrclient.QueryHierarchicalRequirement(ctx, query)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(hr)).ShouldNot(Equal(0))
				Ω(hr[0].FormattedID).Should(Equal(fakeFormattedID))
			})
		})

	})

	Describe(".GetHierarchicalRequirement", func() {
		var (
			fakeObjectID = "50137325678"
			ctrlID, _    = strconv.Atoi(fakeObjectID)
			fakeClient   = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"HierarchicalRequirement": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
				},
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewHierarchicalRequirement(rallyClient)
		})
		Context("when called with a valid objectID", func() {
			It("should return the hierarchichal requirement", func() {
				hr, err := hrclient.GetHierarchicalRequirement(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.ObjectID).Should(Equal(ctrlID))
			})
		})

	})

	Describe(".CreateHierarchicalRequirement", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,

					Body: &fakes.FakeResponseBody{bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewStory", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName   = "NewStory"
			newHrModel = models.HierarchicalRequirement{
				Name: ctrlName,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"
			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewHierarchicalRequirement(rallyClient)
		})
		Context("when called with a valid create request object", func() {
			It("should return the HierarchicalRequirement object created", func() {
				hr, err := hrclient.CreateHierarchicalRequirement(ctx, newHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})
	})

	Describe(".UpdateHierarchicalRequirement", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedStoryName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName      = "UpdatedStoryName"
			updateHrModel = models.HierarchicalRequirement{
				Name:     ctrlName,
				ObjectID: 50137325678,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewHierarchicalRequirement(rallyClient)
		})

		Context("when called with a valid update request object", func() {
			It("should return the HierarchicalRequirement object updated", func() {
				hr, err := hrclient.UpdateHierarchicalRequirement(ctx, updateHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})

	})

	Describe(".DeleteHierarchicalRequirement", func() {
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
			hrclient = NewHierarchicalRequirement(rallyClient)
		})
		Context("when called with a valid delete objectID", func() {
			It("should return the correct operationalresponse struct", func() {
				err := hrclient.DeleteHierarchicalRequirement(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
