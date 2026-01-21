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
	"fmt"
	"net/http"

	. "github.com/aleksofficial/go-rally-rest-toolkit"
	"github.com/aleksofficial/go-rally-rest-toolkit/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RallyClient", func() {
	var (
		ctrlResultCnt = 1
		err           error
		apiKey        string
		apiURL        = "http://myRallyUrl"
		rallyClient   *RallyClient
	)
	Describe(".QueryRequest(string, string, interface) error", func() {
		var (
			fakeOutput *fakes.FakeOutput
			fakeResult *fakes.FakeResult

			fakeFormattedID = "US624340"
			fakeQueryType   = "hierarchicalrequirement"
			fakeClient      = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
				},
			}
		)
		BeforeEach(func() {
			fakeOutput = new(fakes.FakeOutput)
		})
		Context("when called with a vaild query request and a valid api key", func() {
			It("should return the results to the caller", func() {
				apiKey = "abcdef"
				rallyClient = New(apiKey, apiURL, fakeClient)
				query := map[string]string{
					"FormattedID": fakeFormattedID,
				}
				err = rallyClient.QueryRequest(query, fakeQueryType, &fakeOutput)
				fmt.Printf("FakeOutput: %v\n", fakeOutput.QueryResult.TotalResultCount)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fakeOutput.QueryResult.TotalResultCount).Should(Equal(ctrlResultCnt))
			})
		})
		Context("when called with a invaild query request", func() {
			It("should return an error to the caller", func() {
				apiKey = "abcdef"
				fakeResult = new(fakes.FakeResult)
				rallyClient = New(apiKey, apiURL, fakeClient)
				query := map[string]string{
					"FormattedID": fakeFormattedID,
				}
				err = rallyClient.QueryRequest(query, fakeQueryType, &fakeResult)
				fmt.Printf("Error: %v\n", err.Error())
				Ω(err).Should(HaveOccurred())
			})
		})
	})
	Describe(".GetRequest(string, string, interface) error", func() {
		var (
			fakeOutput *fakes.FakeOutput
			fakeResult *fakes.FakeResult

			fakeObjectID  = "50137325678"
			fakeQueryType = "hierarchicalrequirement"
			fakeClient    = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"QueryResult": { "TotalResultCount": 1, "Results": [{"FakeValue": "fakeresponse"}]}}`)},
				},
			}
		)
		BeforeEach(func() {
			fakeOutput = new(fakes.FakeOutput)
		})
		Context("when called with a vaild get request and a valid api key", func() {
			It("should return the results to the caller", func() {
				apiKey = "abcdef"
				rallyClient = New(apiKey, apiURL, fakeClient)
				err = rallyClient.GetRequest(fakeObjectID, fakeQueryType, &fakeOutput)
				fmt.Printf("FakeOutput: %v\n", fakeOutput.QueryResult.TotalResultCount)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fakeOutput.QueryResult.TotalResultCount).Should(Equal(ctrlResultCnt))
			})
		})
		XContext("when called with an invalid return interface get request", func() {
			It("should return an error to the caller", func() {
				apiKey = "abcdef"
				fakeResult = new(fakes.FakeResult)
				rallyClient = New(apiKey, apiURL, fakeClient)
				err = rallyClient.GetRequest(fakeObjectID, fakeQueryType, &fakeResult)
				fmt.Printf("Error: %v\n", err.Error())
				Ω(err).Should(HaveOccurred())
			})
		})
	})
	Describe(".CreateRequest", func() {
		var (
			fakeOutput        *fakes.FakeCreateResponse
			fakeCreateType    = "hierarchicalrequirement"
			fakeCreateRequest = &fakes.FakeCreateRequest{
				FakeItem: fakes.FakeItem{
					Field1: "demostring",
				},
			}
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"CreateResult": { "FakeObject": {"Field1": "demostring"} }}`)},
				},
			}
		)

		BeforeEach(func() {
			fakeOutput = new(fakes.FakeCreateResponse)
		})
		Context("when called with a vaild create request and a valid api key", func() {
			It("should return the results to the caller", func() {
				apiKey = "abcdef"
				rallyClient = New(apiKey, apiURL, fakeClient)
				err = rallyClient.CreateRequest(fakeCreateType, fakeCreateRequest, &fakeOutput)
				fmt.Printf("FakeOutput: %v\n", fakeOutput.CreateResult)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fakeOutput.CreateResult.FakeObject["Field1"]).Should(Equal("demostring"))
			})
		})
		Context("when called with an invalid create request type", func() {
		})
	})
	Describe(".UpdateRequest", func() {
		var (
			fakeOutput        *fakes.FakeUpdateResponse
			fakeUpdateType    = "hierarchicalrequirement"
			fakeUpdateRequest = &fakes.FakeCreateRequest{
				FakeItem: fakes.FakeItem{
					Field1: "demostring",
				},
			}
			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationResult": { "FakeObject": {"Field1": "demostring"} }}`)},
				},
			}
		)

		BeforeEach(func() {
			fakeOutput = new(fakes.FakeUpdateResponse)
		})
		Context("when called with a vaild update request and a valid api key", func() {
			It("should return an error to the caller", func() {
				apiKey = "abcdef"
				rallyClient = New(apiKey, apiURL, fakeClient)
				err = rallyClient.UpdateRequest("12345", fakeUpdateType, fakeUpdateRequest, &fakeOutput)
				fmt.Printf("FakeOutput: %v\n", fakeOutput.OperationResult)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fakeOutput.OperationResult.FakeObject["Field1"]).Should(Equal("demostring"))
			})
		})
		Context("when called with an invalid update request type", func() {
			It("should return an error to the caller", func() {
			})
		})
	})
	Describe(".DeleteRequest", func() {
		var (
			fakeOutput     *fakes.FakeUpdateResponse
			fakeDeleteType = "hierarchicalrequirement"

			fakeClient = &fakes.FakeHTTPClient{
				FakeResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       &fakes.FakeResponseBody{bytes.NewBufferString(`{"OperationResult": { "Errors": [] }}`)},
				},
			}
		)

		BeforeEach(func() {
			fakeOutput = new(fakes.FakeUpdateResponse)
		})

		Context("when called with a vaild delete request and a valid api key", func() {
			It("should return the correct response to the caller", func() {
				apiKey = "abcdef"
				rallyClient = New(apiKey, apiURL, fakeClient)
				err = rallyClient.DeleteRequest("12345", fakeDeleteType, &fakeOutput)
				fmt.Printf("FakeOutput: %v\n", fakeOutput.OperationResult)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when called with an invalid delete request", func() {
			It("should return an error to the caller", func() {
			})
		})
	})
})
