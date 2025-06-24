# Package rules

## Overview
The rules package provides a flexible and powerful business rules engine for implementing complex validation and filtering logic in Go applications. It supports dynamic rule creation, evaluation, and database querying with support for various data types and operators.

## Key Components

### Rule Structure
- **Rule**: Core entity that contains conditions and operators
  - Supports multiple conditions with AND/OR logic
  - Handles different entity types and tenant-based rules
  - Includes metadata like creation and update timestamps

### Condition Types
- **Data Types**: 
  - Integer
  - String
  - Float
  - Boolean

### Supported Operators
- **Rule Operators**:
  - AND
  - OR
- **Condition Operators**:
  - Equals
  - Not Equals
  - Greater Than
  - Less Than
  - In
  - Not In
  - All

## Usage Examples

### 1. Creating and Validating Rules

```go
package main

import (
	"fmt"
	"github.com/omniful/go_commons/rules"
	"github.com/omniful/go_commons/rules/models"
)

func main() {
	// Create a rule with conditions
	rule := &models.Rule{
		Name: models.OrderStatus,
		Operator: models.And,
		Conditions: []models.Condition{
			{
				Key: "status",
				Operator: models.Equals,
				DataType: models.String,
				Values: []string{"pending"},
			},
			{
				Key: "amount",
				Operator: models.GreaterThan,
				DataType: models.Float,
				Values: []string{"100.0"},
			},
		},
	}

	// Create a rule group
	ruleGroup := rules.NewRuleGroup([]*models.Rule{rule})

	// Validate input against rules
	input := map[string]string{
		"status": "pending",
		"amount": "150.0",
	}

	valid, err := ruleGroup.RuleValid(input, []models.Name{models.OrderStatus})
	if err != nil {
		fmt.Printf("Error validating rules: %v\n", err)
		return
	}

	if valid {
		fmt.Println("Input matches the rules")
	} else {
		fmt.Println("Input does not match the rules")
	}
}
```

### 2. Database Integration (PostgreSQL)

```go
package main

import (
	"github.com/omniful/go_commons/rules"
	"github.com/omniful/go_commons/rules/models"
	"gorm.io/gorm"
)

func filterOrders(db *gorm.DB, ruleGroup *rules.RuleGroup) error {
	// Define table mappings
	tables := map[string]string{
		"status": "orders.status",
		"amount": "orders.total_amount",
	}

	// Get database scopes
	scopes, err := ruleGroup.Scopes(rules.Postgres, []models.Name{models.OrderStatus}, tables)
	if err != nil {
		return err
	}

	// Apply rules to query
	if scope, ok := scopes[rules.Postgres].(func(*gorm.DB) *gorm.DB); ok {
		result := db.Scopes(scope).Find(&orders)
		return result.Error
	}

	return nil
}
```

## Features
1. **Type Safety**: Strong typing for operators and data types
2. **Database Integration**: Native support for PostgreSQL queries
3. **Flexible Validation**: Support for multiple comparison operators
4. **Multi-tenant Support**: Built-in tenant isolation
5. **Extensible**: Easy to add new rule types and operators

## Best Practices
1. Always validate rule configurations before applying them
2. Use appropriate data types for conditions
3. Consider performance implications when creating complex rule chains
4. Implement proper error handling for rule validation
5. Use table mappings when working with database queries

## Notes
- The package is designed for high performance and scalability
- Suitable for implementing complex business logic and validation rules
- Supports both in-memory validation and database query generation
- Thread-safe for concurrent rule evaluation
