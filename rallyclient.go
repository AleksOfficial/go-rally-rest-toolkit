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
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// SetConfig sets the configuration for the RallyClient
func (s *RallyClient) SetConfig(config *Config) {
	s.config = config
}

// isRetryableStatusCode returns true if the HTTP status code indicates a transient error
// that should be retried (5xx server errors)
func isRetryableStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// isRetryableError returns true if the error is a transient error that should be retried
// (timeouts, connection errors, etc.)
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	// Check for context deadline exceeded (timeout)
	if err == context.DeadlineExceeded {
		return true
	}
	// Check for network errors by looking for common patterns
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporary failure")
}

// doWithRetry executes an HTTP request with retry logic and exponential backoff
// It retries on 5xx errors and transient network errors, but not on 4xx errors
func (s *RallyClient) doWithRetry(req *http.Request, body []byte) (*http.Response, error) {
	maxRetries := DefaultMaxRetries
	retryDelay := DefaultRetryDelay
	if s.config != nil {
		maxRetries = s.config.MaxRetries
		retryDelay = s.config.RetryDelay
	}

	var lastErr error
	var lastResp *http.Response

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// If this is a retry and we have a body, we need to reset the request body
		if attempt > 0 && body != nil {
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		resp, err := s.client.Do(req)

		if err != nil {
			lastErr = err
			// Check if the error is retryable
			if !isRetryableError(err) || attempt == maxRetries {
				return nil, err
			}
		} else {
			// Check if we should retry based on status code
			if !isRetryableStatusCode(resp.StatusCode) || attempt == maxRetries {
				return resp, nil
			}
			// Close the response body before retrying to avoid resource leak
			resp.Body.Close()
			lastResp = resp
			lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
		}

		// Calculate delay with exponential backoff: delay * 2^attempt
		delay := time.Duration(retryDelay) * time.Millisecond * (1 << attempt)

		// Add jitter: random value between 0 and 50% of the delay to prevent thundering herd
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter

		// Wait before retrying, respecting context cancellation
		select {
		case <-req.Context().Done():
			if lastResp != nil {
				return nil, fmt.Errorf("context cancelled after %d retries: %w", attempt, req.Context().Err())
			}
			return nil, fmt.Errorf("context cancelled after %d retries: %w", attempt, req.Context().Err())
		case <-time.After(delay):
			// Continue to next retry attempt
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
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

	rallyResponse, err := s.doWithRetry(req, nil)
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

	rallyResponse, err := s.doWithRetry(req, nil)
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

	rallyResponse, err := s.doWithRetry(req, inputByteArray)
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

	rallyResponse, err := s.doWithRetry(req, inputByteArray)
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

	rallyResponse, err := s.doWithRetry(req, nil)
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
