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
	"errors"
	"testing"
)

func TestRallyAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *RallyAPIError
		expected string
	}{
		{
			name: "with errors",
			err: &RallyAPIError{
				StatusCode: 400,
				Errors:     []string{"Invalid field", "Missing required field"},
			},
			expected: "Rally API error (status 400): Invalid field; Missing required field",
		},
		{
			name: "with message only",
			err: &RallyAPIError{
				StatusCode: 500,
				Message:    "Internal server error",
			},
			expected: "Rally API error (status 500): Internal server error",
		},
		{
			name: "status code only",
			err: &RallyAPIError{
				StatusCode: 404,
			},
			expected: "Rally API error (status 404)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("RallyAPIError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRallyAPIError_Is(t *testing.T) {
	err := &RallyAPIError{StatusCode: 400, Message: "Bad request"}

	t.Run("matches ErrRallyAPI sentinel", func(t *testing.T) {
		if !errors.Is(err, ErrRallyAPI) {
			t.Error("expected errors.Is(err, ErrRallyAPI) to be true")
		}
	})

	t.Run("matches same status code", func(t *testing.T) {
		target := &RallyAPIError{StatusCode: 400}
		if !errors.Is(err, target) {
			t.Error("expected errors.Is to match same status code")
		}
	})

	t.Run("does not match different status code", func(t *testing.T) {
		target := &RallyAPIError{StatusCode: 500}
		if errors.Is(err, target) {
			t.Error("expected errors.Is to not match different status code")
		}
	})

	t.Run("does not match non-RallyAPIError", func(t *testing.T) {
		target := errors.New("other error")
		if errors.Is(err, target) {
			t.Error("expected errors.Is to not match non-RallyAPIError")
		}
	})
}

func TestRallyAPIError_As(t *testing.T) {
	originalErr := &RallyAPIError{
		StatusCode: 400,
		Message:    "Bad request",
		Errors:     []string{"Invalid field"},
		Warnings:   []string{"Deprecated feature"},
	}

	t.Run("can extract as RallyAPIError", func(t *testing.T) {
		var apiErr *RallyAPIError
		if !errors.As(originalErr, &apiErr) {
			t.Fatal("expected errors.As to succeed")
		}
		if apiErr.StatusCode != 400 {
			t.Errorf("expected StatusCode=400, got %d", apiErr.StatusCode)
		}
		if len(apiErr.Errors) != 1 || apiErr.Errors[0] != "Invalid field" {
			t.Errorf("expected Errors=['Invalid field'], got %v", apiErr.Errors)
		}
	})
}

func TestParseRallyError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		body           string
		expectedErrors []string
	}{
		{
			name:           "OperationResult with errors",
			statusCode:     400,
			body:           `{"OperationResult": {"Errors": ["Invalid field", "Missing value"], "Warnings": []}}`,
			expectedErrors: []string{"Invalid field", "Missing value"},
		},
		{
			name:           "CreateResult with errors",
			statusCode:     400,
			body:           `{"CreateResult": {"Errors": ["Create failed"], "Warnings": ["Deprecated"]}}`,
			expectedErrors: []string{"Create failed"},
		},
		{
			name:           "QueryResult with errors",
			statusCode:     400,
			body:           `{"QueryResult": {"Errors": ["Query failed"], "Warnings": []}}`,
			expectedErrors: []string{"Query failed"},
		},
		{
			name:           "non-JSON body",
			statusCode:     500,
			body:           `Internal Server Error`,
			expectedErrors: nil,
		},
		{
			name:           "empty JSON",
			statusCode:     400,
			body:           `{}`,
			expectedErrors: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseRallyError(tt.statusCode, []byte(tt.body))
			if err.StatusCode != tt.statusCode {
				t.Errorf("expected StatusCode=%d, got %d", tt.statusCode, err.StatusCode)
			}
			if len(err.Errors) != len(tt.expectedErrors) {
				t.Errorf("expected %d errors, got %d: %v", len(tt.expectedErrors), len(err.Errors), err.Errors)
				return
			}
			for i, expected := range tt.expectedErrors {
				if err.Errors[i] != expected {
					t.Errorf("expected error[%d]=%q, got %q", i, expected, err.Errors[i])
				}
			}
		})
	}
}
