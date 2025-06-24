# Package response

## Overview
The response package provides standardized structures and helper functions to construct HTTP responses consistently across an application. This includes utilities for error responses, redirection, and standard success formats, making API development more systematic.

## Key Components
- Error Responses: Predefined formats for error handling.
- Redirect Utilities: Helpers for managing HTTP redirection.
- Success Responses: Consistent structures for successful operations.

## Usage Example
~~~go
package main

import (
	"net/http"
	"github.com/omniful/go_commons/response"
)

func handler(w http.ResponseWriter, r *http.Request) {
	resp := response.Success{Message: "Operation completed successfully"}
	response.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	http.HandleFunc("/api", handler)
	http.ListenAndServe(":8080", nil)
}
~~~

## Notes
- Aims to standardize response formats across services.
