# Package newrelic

## Overview
The newrelic package integrates New Relic performance monitoring into the application. It provides hooks to record metrics, errors, and transaction data for real-time monitoring.

## Key Components
- Integration Hooks: Attach monitoring to key application events.
- Metric Reporting: Send custom metrics to New Relic.
- Error Tracking: Log errors and performance issues.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/newrelic"
)

func main() {
	agent := newrelic.NewAgent(newrelic.Config{
		// Configuration options
	})
	agent.RecordMetric("Custom/Metric", 42)
	fmt.Println("Metric reported to New Relic")
}
~~~

## Error Reporting
For better error visibility in New Relic, it's recommended to use `NoticeCustomError` or `NoticeExpectedCustomError` methods.

~~~go
import (
	"context"
	customError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/newrelic"
)

// Using CustomError type for better error visibility
func reportError(ctx context.Context) {
	// Create a custom error
	err := customError.New("USER_NOT_FOUND", "User not found in database")
	
	// Report to New Relic with proper error class and message
	newrelic.NoticeCustomError(ctx, err)
	
	// For expected errors that shouldn't trigger alerts
	// newrelic.NoticeExpectedCustomError(ctx, err)
}
~~~

When using custom error types from `github.com/omniful/go_commons/error`, New Relic will display:
- The error code as the Error Class (e.g., "USER_NOT_FOUND")
- The error message with greater detail (e.g., "User not found in database")

Without custom error wrapping, errors will be reported with less structured information.

## Notes
- Requires proper agent configuration for effective monitoring.
