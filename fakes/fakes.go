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
}

// Do - Fake HTTP client do method
func (s *FakeHTTPClient) Do(fakeRequest *http.Request) (*http.Response, error) {
	s.SpyRequest = fakeRequest
	return s.FakeResponse, s.FakeError
}
