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
	"encoding/json"
	"fmt"
	"strings"
)

// RallyAPIError represents an error response from the Rally API.
// Rally API returns errors in a structured format with operation results
// containing warnings and errors arrays.
type RallyAPIError struct {
	// StatusCode is the HTTP status code returned by the Rally API
	StatusCode int
	// Message is a human-readable summary of the error
	Message string
	// Errors contains the list of error messages from the Rally API response
	Errors []string
	// Warnings contains the list of warning messages from the Rally API response
	Warnings []string
}

// Error implements the error interface for RallyAPIError.
func (e *RallyAPIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("Rally API error (status %d): %s", e.StatusCode, strings.Join(e.Errors, "; "))
	}
	if e.Message != "" {
		return fmt.Sprintf("Rally API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Rally API error (status %d)", e.StatusCode)
}

// Is implements errors.Is support for RallyAPIError.
// It returns true if the target is a *RallyAPIError with the same StatusCode,
// or if comparing against a sentinel error with StatusCode 0, it matches any RallyAPIError.
func (e *RallyAPIError) Is(target error) bool {
	t, ok := target.(*RallyAPIError)
	if !ok {
		return false
	}
	// If target has StatusCode 0, match any RallyAPIError (sentinel matching)
	if t.StatusCode == 0 {
		return true
	}
	return e.StatusCode == t.StatusCode
}

// ErrRallyAPI is a sentinel error that can be used with errors.Is to check
// if an error is any RallyAPIError.
var ErrRallyAPI = &RallyAPIError{}

// rallyErrorResponse represents the structure of a Rally API error response.
// Rally API wraps operation results in a key like "CreateResult", "QueryResult", etc.
type rallyErrorResponse struct {
	OperationResult *operationResult `json:"OperationResult,omitempty"`
	CreateResult    *operationResult `json:"CreateResult,omitempty"`
	QueryResult     *operationResult `json:"QueryResult,omitempty"`
}

// operationResult represents the common structure for Rally API operation results.
type operationResult struct {
	Errors   []string `json:"Errors"`
	Warnings []string `json:"Warnings"`
}

// parseRallyError attempts to parse a Rally API error response from the given body.
// If parsing fails or no errors are found, it returns a RallyAPIError with just
// the status code and raw body as the message.
func parseRallyError(statusCode int, body []byte) *RallyAPIError {
	apiErr := &RallyAPIError{
		StatusCode: statusCode,
		Message:    string(body),
	}

	// Try to parse as Rally API error response
	var resp rallyErrorResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return apiErr
	}

	// Check which result type is present
	var result *operationResult
	if resp.OperationResult != nil {
		result = resp.OperationResult
	} else if resp.CreateResult != nil {
		result = resp.CreateResult
	} else if resp.QueryResult != nil {
		result = resp.QueryResult
	}

	if result != nil {
		apiErr.Errors = result.Errors
		apiErr.Warnings = result.Warnings
		if len(result.Errors) > 0 {
			apiErr.Message = strings.Join(result.Errors, "; ")
		}
	}

	return apiErr
}
