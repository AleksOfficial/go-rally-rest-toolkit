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
	"context"
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
func (s *RallyClient) QueryRequest(ctx context.Context, query map[string]string, queryType string, output interface{}) error {
	baseURL, err := url.Parse(strings.Join([]string{s.apiurl, queryType}, "/"))
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Add("fetch", "true")
	for idx, val := range query {
		params.Add("query", fmt.Sprintf("( %s = %s )", idx, val))
	}
	baseURL.RawQuery = params.Encode()

	urlStr := baseURL.String()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("ZSESSIONID", s.apikey)

	rallyResponse, err := s.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer rallyResponse.Body.Close()

	content, err := io.ReadAll(rallyResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if rallyResponse.StatusCode < 200 || rallyResponse.StatusCode >= 300 {
		return parseRallyError(rallyResponse.StatusCode, content)
	}

	if err := json.Unmarshal(content, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// GetRequest - Function to perform GET requests when objectID is known.
func (s *RallyClient) GetRequest(ctx context.Context, objectID string, queryType string, output interface{}) error {
	baseURL, err := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Add("fetch", "true")
	baseURL.RawQuery = params.Encode()

	urlStr := baseURL.String()

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("ZSESSIONID", s.apikey)

	rallyResponse, err := s.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer rallyResponse.Body.Close()

	content, err := io.ReadAll(rallyResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if rallyResponse.StatusCode < 200 || rallyResponse.StatusCode >= 300 {
		return parseRallyError(rallyResponse.StatusCode, content)
	}

	if err := json.Unmarshal(content, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func (s *RallyClient) CreateRequest(ctx context.Context, queryType string, input interface{}, output interface{}) error {
	baseURL, err := url.Parse(strings.Join([]string{s.apiurl, queryType, "create"}, "/"))
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	urlStr := baseURL.String()

	inputByteArray, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewReader(inputByteArray))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("ZSESSIONID", s.apikey)

	rallyResponse, err := s.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer rallyResponse.Body.Close()

	content, err := io.ReadAll(rallyResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if rallyResponse.StatusCode < 200 || rallyResponse.StatusCode >= 300 {
		return parseRallyError(rallyResponse.StatusCode, content)
	}

	if err := json.Unmarshal(content, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func (s *RallyClient) UpdateRequest(ctx context.Context, objectID string, queryType string, input interface{}, output interface{}) error {
	baseURL, err := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	urlStr := baseURL.String()

	inputByteArray, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, bytes.NewReader(inputByteArray))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("ZSESSIONID", s.apikey)

	rallyResponse, err := s.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer rallyResponse.Body.Close()

	content, err := io.ReadAll(rallyResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if rallyResponse.StatusCode < 200 || rallyResponse.StatusCode >= 300 {
		return parseRallyError(rallyResponse.StatusCode, content)
	}

	if err := json.Unmarshal(content, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func (s *RallyClient) DeleteRequest(ctx context.Context, objectID string, queryType string, output interface{}) error {
	baseURL, err := url.Parse(strings.Join([]string{s.apiurl, queryType, objectID}, "/"))
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Add("fetch", "true")
	baseURL.RawQuery = params.Encode()

	urlStr := baseURL.String()

	req, err := http.NewRequestWithContext(ctx, "DELETE", urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("ZSESSIONID", s.apikey)

	rallyResponse, err := s.HTTPClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer rallyResponse.Body.Close()

	content, err := io.ReadAll(rallyResponse.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if rallyResponse.StatusCode < 200 || rallyResponse.StatusCode >= 300 {
		return parseRallyError(rallyResponse.StatusCode, content)
	}

	if err := json.Unmarshal(content, output); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
