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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

//RallyClient - struct
type RallyClient struct {
	apikey string
	apiurl string
	client ClientDoer
	config *Config
}

//ClientDoer - interface
type ClientDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// New - creates a new RallyClient
func New(apikey string, apiurl string, client ClientDoer) *RallyClient {
	return &RallyClient{
		apikey: apikey,
		apiurl: apiurl,
		client: client,
	}
}

//HTTPClient - returns the internal client object
func (s *RallyClient) HTTPClient() ClientDoer {
	return s.client
}

// QueryRequest - function to search for an object.
func (s *RallyClient) QueryRequest(query map[string]string, queryType string, output interface{}) (err error) {

	baseURL, _ := url.Parse(strings.Join([]string{s.apiurl, queryType}, "/"))

	params := url.Values{}
	params.Add("fetch", "true")
	for idx, val := range query {
		params.Add("query", fmt.Sprintf("( %s = %s )", idx, val))
	}
	baseURL.RawQuery = params.Encode()

	urlStr := fmt.Sprintf("%v", baseURL)

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("ZSESSIONID", s.apikey)
	if rallyResponse, err := s.HTTPClient().Do(req); err == nil {
		content, _ := io.ReadAll(rallyResponse.Body)
		if err = json.Unmarshal(content, output); err != nil {
			return err
		}
	}
	return nil
}

// GetRequest - Function to perform GET requests when objectID is known.
func (s *RallyClient) GetRequest(objectID string, queryType string, output interface{}) (err error) {
	baseURL, _ := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))

	params := url.Values{}
	params.Add("fetch", "true")
	baseURL.RawQuery = params.Encode()

	urlStr := fmt.Sprintf("%v", baseURL)

	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("ZSESSIONID", s.apikey)

	if rallyResponse, err := s.HTTPClient().Do(req); err == nil {
		content, _ := io.ReadAll(rallyResponse.Body)
		if err = json.Unmarshal(content, output); err != nil {
			return err
		}
	}
	return nil
}

func (s *RallyClient) CreateRequest(queryType string, input interface{}, output interface{}) (err error) {
	baseURL, _ := url.Parse(strings.Join([]string{s.apiurl, queryType, "create"}, "/"))

	urlStr := fmt.Sprintf("%v", baseURL)

	inputByteArray, err := json.Marshal(input)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(inputByteArray))
	req.Header.Add("ZSESSIONID", s.apikey)

	if rallyResponse, err := s.HTTPClient().Do(req); err == nil {
		content, _ := io.ReadAll(rallyResponse.Body)
		if err = json.Unmarshal(content, output); err != nil {
			return err
		}
	}
	return nil
}

func (s *RallyClient) UpdateRequest(objectID string, queryType string, input interface{}, output interface{}) (err error) {
	baseURL, _ := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))

	urlStr := fmt.Sprintf("%v", baseURL)

	inputByteArray, err := json.Marshal(input)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(inputByteArray))
	req.Header.Add("ZSESSIONID", s.apikey)

	if rallyResponse, err := s.HTTPClient().Do(req); err == nil {
		content, _ := io.ReadAll(rallyResponse.Body)
		if err = json.Unmarshal(content, output); err != nil {
			return err
		}
	}
	return nil
}

func (s *RallyClient) DeleteRequest(objectID string, queryType string, output interface{}) (err error) {
	baseURL, _ := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))

	params := url.Values{}
	params.Add("fetch", "true")
	baseURL.RawQuery = params.Encode()

	urlStr := fmt.Sprintf("%v", baseURL)

	req, _ := http.NewRequest("DELETE", urlStr, nil)
	req.Header.Add("ZSESSIONID", s.apikey)

	if rallyResponse, err := s.HTTPClient().Do(req); err == nil {
		content, _ := io.ReadAll(rallyResponse.Body)
		if err = json.Unmarshal(content, output); err != nil {
			return err
		}
	}
	return nil
}
