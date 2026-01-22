# Go Rally REST Toolkit

A Go client library for the [Rally REST API](https://rally1.rallydev.com/slm/doc/webservice/).

## Features

- Full support for Rally REST API v2.0
- Context support for cancellation and timeouts
- Automatic retry with exponential backoff for transient failures
- Custom error types for Rally API errors
- Environment variable configuration

## Installation

```bash
go get github.com/aleksofficial/go-rally-rest-toolkit
```

Requires Go 1.21 or later.

## Quick Start

The simplest way to get started is using environment variables:

```go
package main

import (
    "context"
    "fmt"
    "log"

    rally "github.com/aleksofficial/go-rally-rest-toolkit"
)

func main() {
    // Create client from environment variables (requires RALLY_API_KEY)
    client, err := rally.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // Query for a user story
    ctx := context.Background()
    query := map[string]string{
        "FormattedID": "US12345",
    }

    var result struct {
        QueryResult struct {
            Results          []map[string]interface{}
            TotalResultCount int
        }
    }

    err = client.QueryRequest(ctx, query, "hierarchicalrequirement", &result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d results\n", result.QueryResult.TotalResultCount)
}
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `RALLY_API_KEY` | Yes | - | Rally API key for authentication |
| `RALLY_BASE_URL` | No | `https://rally1.rallydev.com/slm/webservice/v2.0` | Rally API base URL |
| `RALLY_TIMEOUT` | No | `30` | HTTP timeout in seconds |
| `RALLY_MAX_RETRIES` | No | `3` | Maximum retry attempts for transient failures |
| `RALLY_RETRY_DELAY` | No | `1000` | Initial retry delay in milliseconds |

## Manual Configuration

For more control, you can create a client with explicit parameters:

```go
package main

import (
    "context"
    "net/http"
    "time"

    rally "github.com/aleksofficial/go-rally-rest-toolkit"
)

func main() {
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }

    client := rally.New(
        "your-api-key",
        "https://rally1.rallydev.com/slm/webservice/v2.0",
        httpClient,
    )

    // Optionally set retry configuration
    client.SetConfig(&rally.Config{
        MaxRetries: 5,
        RetryDelay: 2000,
    })

    ctx := context.Background()
    // Use client...
}
```

## API Methods

All methods accept a `context.Context` as the first parameter for cancellation and timeout support.

### QueryRequest

Search for Rally artifacts using query parameters:

```go
query := map[string]string{
    "FormattedID": "DE12345",
}
var result DefectQueryResult
err := client.QueryRequest(ctx, query, "defect", &result)
```

### GetRequest

Retrieve a specific artifact by its ObjectID:

```go
var result DefectResult
err := client.GetRequest(ctx, "12345678", "defect", &result)
```

### CreateRequest

Create a new artifact:

```go
newDefect := CreateDefectRequest{
    Defect: Defect{
        Name:     "Bug in login",
        Priority: "High",
        Severity: "Major Problem",
    },
}
var result CreateResult
err := client.CreateRequest(ctx, "defect", newDefect, &result)
```

### UpdateRequest

Update an existing artifact:

```go
update := UpdateDefectRequest{
    Defect: Defect{
        State: "Fixed",
    },
}
var result OperationResult
err := client.UpdateRequest(ctx, "12345678", "defect", update, &result)
```

### DeleteRequest

Delete an artifact:

```go
var result OperationResult
err := client.DeleteRequest(ctx, "12345678", "defect", &result)
```

## Error Handling

The library provides structured error types for Rally API errors:

```go
err := client.QueryRequest(ctx, query, "defect", &result)
if err != nil {
    var rallyErr *rally.RallyAPIError
    if errors.As(err, &rallyErr) {
        fmt.Printf("Rally API error: %s (status %d)\n", rallyErr.Message, rallyErr.StatusCode)
        for _, e := range rallyErr.Errors {
            fmt.Printf("  - %s\n", e)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Retry Behavior

The client automatically retries requests that fail due to:

- Server errors (5xx status codes)
- Network timeouts
- Connection refused/reset errors

Retries use exponential backoff with jitter. Client errors (4xx) are not retried.

Configure retry behavior via environment variables or the `SetConfig` method.

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

Originally developed by Comcast Cable Communications Management, LLC.
