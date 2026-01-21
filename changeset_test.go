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
	"net/http"
	"strconv"

	. "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/fakes"
	"github.com/aleksofficial/go-rally-rest-toolkit/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Changeset", func() {

	var (
		apiKey      string
		apiURL      = "http://myRallyUrl"
		rallyClient *RallyClient
		hrclient    *Changeset
	)
	Describe(".QueryChangeset", func() {

		var (
			fakeMessage = "concourse-1"
			fakeClient  = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Message": "concourse-1","Errors": [], "Warnings": []}]}}`)},
				},
			}
		)

		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewChangeset(rallyClient)
		})
		Context("when called with a valid formattedID", func() {
			It("should return the requested array of Changeset results", func() {
				query := map[string]string{
					"Message": fakeMessage,
				}
				hr, err := hrclient.QueryChangeset(query)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(len(hr)).ShouldNot(Equal(0))
				Ω(hr[0].Message).Should(Equal(fakeMessage))
			})
		})

	})

	Describe(".GetChangeset", func() {
		var (
			fakeObjectID = "50137325678"
			ctrlID, _    = strconv.Atoi(fakeObjectID)
			fakeClient   = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"Changeset": {"CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}`)},
				},
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewChangeset(rallyClient)
		})
		Context("when called with a valid objectID", func() {
			It("should return the Changeset", func() {
				hr, err := hrclient.GetChangeset(fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.ObjectID).Should(Equal(ctrlID))
			})
		})

	})

	Describe(".CreateChangeset", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,

					Body: &fakes.FakeResponseBody{bytes.NewBufferString(`{"CreateResult": {"Object": {"Name": "NewChangeset", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName   = "NewChangeset"
			newHrModel = models.Changeset{
				Name: ctrlName,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"
			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewChangeset(rallyClient)
		})
		Context("when called with a valid create request object", func() {
			It("should return the Changeset object created", func() {
				hr, err := hrclient.CreateChangeset(newHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})
	})

	Describe(".UpdateChangeset", func() {
		var (
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationalResult": {"Object": {"Name": "UpdatedChangesetName", "CreationDate": "2016-01-21T21:47:08.551Z", "ObjectID": 50137325678,"Errors": [], "Warnings": []}}}`)},
				},
			}
			ctrlName      = "UpdatedChangesetName"
			updateHrModel = models.Changeset{
				Name:     ctrlName,
				ObjectID: 50137325678,
			}
		)
		BeforeEach(func() {
			apiKey = "abcdef"

			rallyClient = New(apiKey, apiURL, fakeClient)
			hrclient = NewChangeset(rallyClient)
		})

		Context("when called with a valid update request object", func() {
			It("should return the Changeset object updated", func() {
				hr, err := hrclient.UpdateChangeset(updateHrModel)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(hr.Name).Should(Equal(ctrlName))
			})
		})

	})

	Describe(".DeleteChangeset", func() {
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
			hrclient = NewChangeset(rallyClient)
		})
		Context("when called with a valid delete objectID", func() {
			It("should return the correct operationalresponse struct", func() {
				err := hrclient.DeleteChangeset(fakeObjectID)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
