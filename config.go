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
	"net/http"
	"os"
	"strconv"
	"time"
)

// Default configuration values
const (
	DefaultBaseURL    = "https://rally1.rallydev.com/slm/webservice/v2.0"
	DefaultTimeout    = 30
	DefaultMaxRetries = 3
	DefaultRetryDelay = 1000
)

// Config holds all configuration for the Rally client
type Config struct {
	// APIKey is the Rally API key for authentication (required)
	APIKey string
	// BaseURL is the Rally API base URL (optional, defaults to DefaultBaseURL)
	BaseURL string
	// Timeout is the HTTP timeout in seconds (optional, defaults to 30)
	Timeout int
	// MaxRetries is the maximum number of retry attempts for transient failures (optional, defaults to 3)
	MaxRetries int
	// RetryDelay is the initial retry delay in milliseconds (optional, defaults to 1000)
	RetryDelay int
}

// ErrAPIKeyRequired is returned when RALLY_API_KEY environment variable is not set
var ErrAPIKeyRequired = errors.New("RALLY_API_KEY environment variable is required")

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	apiKey := os.Getenv("RALLY_API_KEY")
	if apiKey == "" {
		return nil, ErrAPIKeyRequired
	}

	config := &Config{
		APIKey:     apiKey,
		BaseURL:    DefaultBaseURL,
		Timeout:    DefaultTimeout,
		MaxRetries: DefaultMaxRetries,
		RetryDelay: DefaultRetryDelay,
	}

	if baseURL := os.Getenv("RALLY_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if timeout := os.Getenv("RALLY_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil && t > 0 {
			config.Timeout = t
		}
	}

	if maxRetries := os.Getenv("RALLY_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil && r >= 0 {
			config.MaxRetries = r
		}
	}

	if retryDelay := os.Getenv("RALLY_RETRY_DELAY"); retryDelay != "" {
		if d, err := strconv.Atoi(retryDelay); err == nil && d >= 0 {
			config.RetryDelay = d
		}
	}

	return config, nil
}

// NewClientFromEnv creates a new RallyClient using configuration from environment variables
func NewClientFromEnv() (*RallyClient, error) {
	config, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	client := New(config.APIKey, config.BaseURL, httpClient)
	client.config = config

	return client, nil
}
