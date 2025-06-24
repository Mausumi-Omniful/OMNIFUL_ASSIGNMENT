# Package set

## Overview
The set package provides a data structure to handle collections of unique elements and supports standard set operations like union, intersection, and difference.

## Key Components
- Set Data Structure: Encapsulates unique elements.
- Set Operations: Union, Intersection, Difference, etc.
- Utility Functions: Manage and display set contents.

## Usage Example
~~~go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/set"
)

func main() {
	s := set.New()
	s.Add("a")
	s.Add("b")
	fmt.Println("Set contents:", s.Items())
}
~~~

## Notes
- Ideal for ensuring uniqueness in data collections.
