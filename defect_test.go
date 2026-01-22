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

var _ = Describe("Defect", func() {

	var (
		apiKey      string
		apiURL      = "http://myRallyUrl"
		rallyClient *RallyClient
		hrclient    *Defect
		ctx         = context.Background()
	)
	Describe(".QueryDefect", func() {

		var (
			fakeFormattedID = "DE624340"
			fakeClient      = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"FormattedID": "DE624340","Errors": [], "Warnings": []}]}}`)},
				},
			}
		)

		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewDefect(rallyClient)
		})
		Context("when called with a valid formattedID", func() {
			It("should return the requested array of defect results", func() {
				query := map[string]string{
					"FormattedID": fakeFormattedID,
				}
				hr, err := hrclient.QueryDefect(ctx, query)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(hr)).ShouldNot(Equal(0))
				Ω(hr[0].FormattedID).Should(Equal(fakeFormattedID))
			})
		})

	})

	Describe(".GetDefect", func() {
		var (
			fakeObjectID = "50137325678"
			ctrlID, _    = strconv.Atoi(fakeObjectID)
			fakeClient   = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"Defect": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
				},
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewDefect(rallyClient)
		})
		Context("when called with a valid objectID", func() {
			It("should return the Defect", func() {
				hr, err := hrclient.GetDefect(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.ObjectID).Should(Equal(ctrlID))
			})
		})

	})

	Describe(".CreateDefect", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,

					Body: &fakes.FakeResponseBody{bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewStory", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName   = "NewStory"
			newHrModel = models.Defect{
				Name: ctrlName,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"
			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewDefect(rallyClient)
		})
		Context("when called with a valid create request object", func() {
			It("should return the Defect object created", func() {
				hr, err := hrclient.CreateDefect(ctx, newHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})
	})

	Describe(".UpdateDefect", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedStoryName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName      = "UpdatedStoryName"
			updateHrModel = models.Defect{
				Name:     ctrlName,
				ObjectID: 50137325678,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewDefect(rallyClient)
		})

		Context("when called with a valid update request object", func() {
			It("should return the Defect object updated", func() {
				hr, err := hrclient.UpdateDefect(ctx, updateHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})

	})

	Describe(".DeleteDefect", func() {
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
			hrclient = NewDefect(rallyClient)
		})
		Context("when called with a valid delete objectID", func() {
			It("should return the correct operationalresponse struct", func() {
				err := hrclient.DeleteDefect(ctx, fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
