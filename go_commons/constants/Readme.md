# Package constants

## Overview
The constants package defines shared constant values used throughout the application, including configuration keys, error codes, and other invariant values, promoting consistency and maintainability.

## Key Components
- Constant Definitions: Immutable values that guide application behavior.
- Grouped Constants: Logical grouping for easier management.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/constants"
)

func main() {
	fmt.Println("Application mode:", constants.AppMode)
}
~~~

## Notes
- Helps avoid the use of magic numbers and strings.
