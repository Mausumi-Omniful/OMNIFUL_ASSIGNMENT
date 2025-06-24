# Go Commons Validator Package

A comprehensive validation package for internal services using [go-playground/validator v10](https://github.com/go-playground/validator). This package provides custom validations, error formatting, and integration with the internal error package for consistent validation across all internal services.

## Features

- **Custom Validations**: Phone numbers, passwords, slugs, ObjectIDs, and more
- **Error Formatting**: Consistent error format using the internal error package
- **Singleton Pattern**: Efficient validator instance management
- **JSON Tag Support**: Uses JSON field names in error messages
- **Custom Error Integration**: All validators return CustomError from the error package
- **Gin Framework Integration**: Automatic validation when using Gin binding methods
- **Internationalization**: Support for localized validation messages
- **Factory Pattern**: Modular and extensible validator implementation
- **Registry Pattern**: Automatic registration of validators 
- **Smart Empty Field Handling**: Skip validation for empty fields unless marked as required
- **Comprehensive Examples**: Real-world validation scenarios

## Important Notes

1. **Error Format**: All validation functions return `CustomError` from the internal error package. No custom response formats are created.
2. **Context Required**: All validation functions require a context parameter for i18n support.
3. **Gin Integration**: When using Gin's binding methods like `ShouldBindJSON`, validation happens automatically.
4. **Empty Field Handling**: Empty fields are going to be validated unless they are also marked as `omniempty`. This means a field with `validate:"phone"` will trigger validation errors, but a field with `validate:"omitempty,phone"` will not trigger if the field is empty or not provided.
5. **Struct and Custom Type Validation**: Validation tags on custom types (structs, custom types) will be applied regardless of whether the parent field is empty. Validation works consistently for both pointer and non-pointer structs.

## Installation

This package is part of the go_commons module. Import it in your service:

```go
import "github.com/omniful/go_commons/validator"
```

## Basic Usage

### 1. Validate a Struct Directly

```go
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,strong_password"`
    Name     string `json:"name" validate:"required,alphaspace,min=2,max=50"`
    Phone    string `json:"phone" validate:"omitempty,phone"`
}

func CreateUser(ctx context.Context, req CreateUserRequest) error {
    // All validation functions require context and return CustomError
    if cusErr := validator.ValidateStruct(ctx, req); cusErr.Exists() {
        // Handle validation error
        oresponse.NewErrorResponse(ctx, cusErr, analyticsError.CustomCodeToHttpCodeMapping)
        return
    }
    // Process user creation...
    return nil
}
```

### 2. Automatic Validation with Gin

```go
func CreateUserHandler(c *gin.Context) {
    var req CreateUserRequest
    
    var req UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Binding already includes validation
        oresponse.NewErrorResponseV2(ctx, validators.FormatValidationErrors(ctx, err))
        return
    }
    
    // Request is already validated, proceed with business logic
    user, err := userService.CreateUser(c.Request.Context(), req)
    if err != nil {
        errorHandler(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, user)
}
```

## Custom Validations

| Tag | Description | Example |
|-----|-------------|---------|
| `strong_password` | Min 8 chars, 3 of: upper/lower/digit/special | `P@ssw0rd123` |

## Built-in Validations (from validator v10)

Common built-in validations you can use:

| Tag | Description |
|-----|-------------|
| `required` | Field must not be empty |
| `email` | Valid email address |
| `min=n` | Minimum length/value |
| `max=n` | Maximum length/value |
| `len=n` | Exact length |
| `gte=n` | Greater than or equal |
| `lte=n` | Less than or equal |
| `oneof=foo bar` | Must be one of the values |
| `numeric` | Numeric string |
| `alpha` | Alphabetic characters only |
| `alphanum` | Alphanumeric characters |
| `uuid` | Valid UUID |
| `url` | Valid URL |
| `json` | Valid JSON string |
| `datetime=2006-01-02` | Date/time in specific format |

## Error Handling

All validation functions return `CustomError` from the internal error package with code `REQUEST_INVALID`.

### Working with Validation Errors

```go
// Using any validation function
ctx := context.Background()
customErr := validator.ValidateStruct(ctx, data)
if customErr.Exists() {
    // customErr.ErrorCode() == "REQUEST_INVALID"
    // customErr.ErrorMessage() == "Invalid Request" or localized message
    
    // Get detailed errors by field
    errorsMap := customErr.ErrorMap()
    for field, message := range errorsMap {
        fmt.Printf("Field %s: %s\n", field, message)
    }
    
    // Add extra context to the error
    customErr = customErr.WithParam("additional_info", "More context")
}
```

### Low-Level Validation Access

If you need access to the raw validation errors:

```go
// Using RawValidateStruct for direct access to validator errors
ctx := context.Background()
err := validator.RawValidateStruct(ctx, data)
if err != nil {
    // Convert to CustomError with formatted validation errors
    customErr := validator.FormatValidationErrors(ctx, err)
    
    // Access the error map
    errorsMap := customErr.ErrorMap()
}
```

## Gin Integration

The validator is integrated with the Gin framework for automatic validation:

```go
// In your main.go or init code
import (
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/omniful/go_commons/validator"
)

func init() {
    // This is already done in the validator package, but shown here for clarity
    binding.Validator = validator.NewGinValidator()
}

// In your handler
func CreateItem(c *gin.Context) {
    var req ItemRequest
    
    // This will automatically validate using our custom validator
    if err := c.ShouldBindJSON(&req); err != nil {
        // err will be a of type error from the gin
        oresponse.NewErrorResponseV2(ctx, validators.FormatValidationErrors(ctx, err))
        return
    }
    
    // Request is valid, continue with business logic
}
```

## Advanced Usage


### Nested Struct Validation

The validator now properly handles nested struct validation with improved behavior:

1. Validation consistently works for both pointer and non-pointer struct fields
2. Nested struct validation occurs even if the parent field is empty (unless the field is not required)
3. All validation tags inside nested structs are respected regardless of pointer status
4. The `dive` tag tells the validator to apply validation to each element of the slice or array

Example with nested struct validation:

```go
// Example 1: Nested struct validation
type Report struct {
    ID   string `json:"id" validate:"required"`
    Pipe Pipe   `json:"pipe" validate:"required"` // Struct level validation
}

type Pipe struct {
    Name string `json:"name" validate:"required"`
    Type string `json:"type" validate:"oneof=input output"`
}

// The validator will:
// 1. Check if Report.ID is provided
// 2. Check if Report.Pipe is provided (not empty)
// 3. If Pipe is provided, validate Pipe.Name and Pipe.Type
// Validation errors will show the full path: "pipe.name is required"

// Example 2: Slice validation with dive
type Order struct {
    Items []OrderItem `json:"items" validate:"required,min=1,dive"`
}

type OrderItem struct {
    ProductID string `json:"product_id" validate:"required,bsonid"`
    Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

// The validator will:
// 1. Check if Items slice exists and has at least one element
// 2. For each item in the slice, validate all its fields
```

Example with embedded structs:

```go
type UserProfile struct {
    User    User     `json:"user" validate:"required"` // User fields will be validated
    Address Address  `json:"address" validate:"required"` // Address fields will be validated
    Tags    []string `json:"tags" validate:"dive,max=10"` // Each tag will be validated
}
```

### Array/Slice Validation

```go
type Request struct {
    Tags       []string `json:"tags" validate:"required,min=1,max=10,dive,min=2,max=20"`
    Categories []string `json:"categories" validate:"dive,oneof=tech finance health"`
}
```

### Third-Party Type Validation

The validator supports proper validation of third-party types, including `null` types from libraries like `guregu/null`:

```go
import "gopkg.in/guregu/null.v4"

type Product struct {
    ID          string    `json:"id" validate:"required"`
    Name        string    `json:"name" validate:"required"`
    Stock       null.Int  `json:"stock" validate:"required,gte=0"`
    UpdatedAt   null.Time `json:"updated_at" validate:"required"`
}
```

Validation behavior for null types:
1. If the field is marked as `required`:
   - The field must be provided
   - Any Additional rules over the filed like gte, lte would panic

2. If the field is not marked as `required`:
   - The field can be omitted
   - If provided, it can be null (Valid=false)
   - Any Additional rules over the filed like gte, lte would panic


This validation behavior is consistent across all supported third-party types and works the same way whether the fields are pointers or values.

## Best Practices

1. **Use Context**: Always pass a context to validation functions for i18n support
2. **Use Gin Binding**: For Gin applications, use binding methods for automatic validation
3. **Handle Nested Errors**: Check for nested field validation errors in complex structs
4. **Validate Early**: Validate input as soon as it enters your system
5. **Review Struct Tags**: Ensure validation tags on nested structs and custom types align with your expectations, as validation now works consistently for all field types (pointers, non-pointers, custom types)
6. **Use `required` Appropriately**: Only mark fields as `required` if they must be present; other validation tags will be skipped for empty fields unless also marked as `required`

## Package Structure

The validator package uses a modular structure for custom validators:

```
validator/
  ├── validator.go             # Main validator package with core functions
  ├── errors.go                # Error formatting and handling
  ├── custom_validators/       # Custom validators directory
  │   ├── validator.go         # Base validator interface and registry
  │   ├── string_validators.go # String-related validators (phone, alphaspace, slug, password)
  │   ├── date_validators.go   # Date-related validators
  │   ├── numeric_validators.go # Numeric validators
  │   └── mongo_validators.go  # MongoDB-related validators
  └── tests/                   # Test cases
```

## Adding New Custom Validators

To add a new custom validator:

1. Determine which category your validator belongs to and add it to the appropriate file in `custom_validators/`:

```go
// For example, adding a new string validator to string_validators.go
package custom_validators

// myValidator validates something specific
type myValidator struct{}

// newMyValidator creates a new validator
func newMyValidator() customValidator {
    return &myValidator{}
}

// name returns the validation tag name
func (v *myValidator) name() string {
    return "my_custom"
}

// validate performs the validation logic
func (v *myValidator) validate(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    // Your validation logic
    return isValid
}
```

2. Register it in the `init` function in `custom_validators/validator.go`:

```go
func init() {
    // ... existing validators
    
    // Add your new validator
    registerValidator(newMyValidator())
}
```

3. Use it in your structs:

```go
type MyStruct struct {
    Field string `json:"field" validate:"my_custom"`
}
```

The registry pattern automatically registers all validators, so you only need to add your validator to the `init` function, and it will be automatically registered with the validator engine.

## Testing

When testing services that use this validator:

```go
func TestValidation(t *testing.T) {
    // Test validation
    req := CreateUserRequest{
        Email: "invalid-email",
    }
    
    ctx := context.Background()
    customErr := validator.ValidateStruct(ctx, req)
    assert.True(t, customErr.Exists())
    assert.Equal(t, customError.RequestInvalid, customErr.ErrorCode())
    
    // Check validation errors
    errorsMap := customErr.ErrorMap()
    assert.Contains(t, errorsMap, "email")
    assert.Contains(t, errorsMap["email"], "valid email")
}
```

## Contributing

When adding new features:

1. Add custom validators to the appropriate file in the `custom_validators` directory
2. For entirely new validation types, create a new file following the naming pattern
3. Register your validator in the `init` function in `validator.go`
4. Ensure error messages are clear and helpful
5. Always return CustomError from the error package
6. Update this README
7. Add tests for new validations

## Support

For questions or issues with the validator package, please contact the platform team or create an issue in the go_commons repository.

## Internationalization (i18n)

The validator supports internationalization through integration with the internal i18n package. Validation error messages can be displayed in multiple languages based on the context.

### How to Use i18n with Validator

1. Initialize the i18n package in your service with your localization directory:

```go
// Initialize service-specific localization
err := i18n.Initialize(i18n.WithRootPath("./localization"))
if err != nil {
    log.Errorf("unable to initialise localization, err: %v", err)
    panic("Error: unable to initialise localization")
}
```

2. Create message files following this structure:

```
localization/
  ├── messages.en.json  # English translations
  └── messages.ar.json  # Arabic translations (optional)
```

3. Add validation message translations to your JSON files:

```json
{
  "validation.failed": {
    "message": "Service-Specific Validation Failed"
  },
  "validation.required": {
    "message": "The {field} field is absolutely required!"
  }
}
```

4. Set language in your context (for API handlers):

```go
func ApiHandler(c *gin.Context) {
    // Get language from header, query param, etc.
    lang := c.GetHeader("Accept-Language")
    
    // Set language in context (already done in middleware)
    ctx := i18n.ContextWithLanguage(c.Request.Context(), lang)
    
    // Use this context for validation
    customErr := validator.ValidateStruct(ctx, data)
}
```

### Fallback Behavior

The system will look for translations in this order:
1. Service-specific translations in the requested language
2. Service-specific translations in English
3. Validator's built-in translations in the requested language
4. Validator's built-in translations in English
5. Default fallback message 