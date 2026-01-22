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

package fakes

import (
	"io"
	"net/http"
)

type FakeOutput struct {
	QueryResult struct {
		Results          []FakeResult
		TotalResultCount int
	}
}

type FakeCreateResponse struct {
	CreateResult FakeObject
}

type FakeUpdateResponse struct {
	OperationResult FakeObject
	Errors          []map[string]interface{}
}

type FakeObject struct {
	FakeObject map[string]interface{}
}
type FakeCreateRequest struct {
	FakeItem FakeItem
}

type FakeItem struct {
	Field1 string
}

type FakeResult struct {
	FakeValue string
}

//FakeResponseBody - a fake response body object
type FakeResponseBody struct {
	Reader io.Reader
}

// Read implements io.Reader
func (f *FakeResponseBody) Read(p []byte) (n int, err error) {
	return f.Reader.Read(p)
}

//Close - close fake body
func (FakeResponseBody) Close() error { return nil }

//FakeRequestBody - a fake response body object
type FakeRequestBody struct {
	io.Reader
}

//Close - close fake body
func (FakeRequestBody) Close() error { return nil }

//FakeHTTPClient - a fake http client
type FakeHTTPClient struct {
	http.Client
	SpyRequest   *http.Request
	FakeResponse *http.Response
	FakeError    error
	// CallCount tracks how many times Do was called (for retry testing)
	CallCount int
	// FakeResponses allows returning different responses on subsequent calls (for retry testing)
	// If set, FakeResponse is ignored and FakeResponses[CallCount] is returned instead
	FakeResponses []*http.Response
	// FakeErrors allows returning different errors on subsequent calls (for retry testing)
	FakeErrors []error
}

// Do - Fake HTTP client do method
func (s *FakeHTTPClient) Do(fakeRequest *http.Request) (*http.Response, error) {
	s.SpyRequest = fakeRequest
	idx := s.CallCount
	s.CallCount++

	// If FakeResponses or FakeErrors are set, use them based on call count
	if len(s.FakeResponses) > 0 || len(s.FakeErrors) > 0 {
		var resp *http.Response
		var err error

		if idx < len(s.FakeResponses) {
			resp = s.FakeResponses[idx]
		}
		if idx < len(s.FakeErrors) {
			err = s.FakeErrors[idx]
		}
		return resp, err
	}

	return s.FakeResponse, s.FakeError
}
