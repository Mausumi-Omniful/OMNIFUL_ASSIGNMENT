# Migration Guide - Adopting the go_commons Validator

This guide helps teams migrate their existing validation logic to use the centralized go_commons validator package.

## Quick Start

### 1. Update Your Imports

Replace individual validator imports with the go_commons validator:

```go
// Before:
import "github.com/go-playground/validator/v10"

// After:
import "github.com/omniful/go_commons/validator"
```

### 2. Update Struct Tags

No changes needed if you're already using standard validator v10 tags. Our custom validations are additional:

```go
type User struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,strong_password"` // Custom validation
}
```

### 3. Gin Integration - Automatic Validation

The validator is now integrated with Gin. When using `ShouldBindJSON` (or other Gin binding methods), validation happens automatically:

```go
// Before:
func CreateUser(c *gin.Context) {
    var req UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        cusErr := commonError.NewCustomError(analyticsError.BadRequest, err.Error())
		oresponse.NewErrorResponse(ctx, cusErr, analyticsError.CustomCodeToHttpCodeMapping)
        return
    }
    
    // Separate validation step
    if err := validate.Struct(req); err != nil {
        // Handle validation error
        oresponse.NewErrorResponse(ctx, cusErr, analyticsError.CustomCodeToHttpCodeMapping)
        return
    }
    
    // Process valid request...
}

// After:
func CreateUser(c *gin.Context) {
    var req UserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Binding already includes validation
        oresponse.NewErrorResponseV2(ctx, validators.FormatValidationErrors(ctx, err))
        return
    }
    
    // Process valid request...
}
```

## Common Migration Scenarios

### 1. Custom Validation Functions

If you have service-specific validations, you can add them using the direct approach:

```go
// Create a service-specific validator
func init() {
    v := validator.GetValidator()
    v.RegisterValidation("my_custom", myCustomValidationFunc)
}

func myCustomValidationFunc(fl validator.FieldLevel) bool {
    // Your validation logic
    return true
}
```

For common validations that should be added to the go_commons package, you can contribute by:

1. Identifying the appropriate file in the `custom_validators` directory based on the type
2. Adding your validator to that file or creating a new file for a new category
3. Adding your validator to the registry in the `init` function:

```go
// In custom_validators/validator.go
func init() {
    // ... existing validators
    
    // Add your new validator
    registerValidator(newMyValidator())
}
```

The registry pattern automatically registers all validators when the package is initialized.

### 2. Testing

Update your tests to use the validator package:

```go
func TestUserValidation(t *testing.T) {
    user := User{
        Email: "invalid",
        Password: "123",
    }
    
    ctx := context.Background()
    customErr := validator.ValidateStruct(ctx, user)
    assert.True(t, customErr.Exists())
    
    // Access validation errors from the error map
    errorsMap := customErr.ErrorMap()
    assert.Len(t, errorsMap, 2)
    
    // Check for specific field errors
    assert.Contains(t, errorsMap, "email")
    assert.Contains(t, errorsMap, "password")
}
```

## Gradual Migration Strategy

1. **Phase 1**: Update imports and use go_commons validator alongside existing validation
2. **Phase 2**: Replace validation logic in new endpoints
3. **Phase 3**: Gradually migrate existing endpoints to use Gin binding
4. **Phase 4**: Remove old validation code

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

## Checklist

- [ ] Updated imports to use go_commons validator
- [ ] Reviewed and updated struct tags for custom validations
- [ ] Updated error handling to use the returned CustomError directly
- [ ] Removed redundant validation code after Gin binding
- [ ] Updated tests to work with CustomError return values
- [ ] Updated API documentation for error response format

## Need Help?

- Check the [README.md](./README.md) for detailed documentation
- Review the tests in the [tests](./tests) directory for usage examples
- Contact the platform team for assistance
- Submit PRs for new custom validations that would benefit other services

## Internationalization Support

The validator package now supports internationalization (i18n) through the internal i18n package. This means validation error messages can be displayed in multiple languages.

### How to Use i18n in Your Services

If you want to customize the validation messages for your service:

1. Set up localization files in your service:

```
localization/
  ├── messages.en.json  # English translations
  └── messages.ar.json  # Arabic translations (optional)
```

2. Initialize the i18n package in your service:

```go
// Initialize service-specific localization
err := i18n.Initialize(i18n.WithRootPath("./localization"))
if err != nil {
    log.Errorf("unable to initialise localization, err: %v", err)
    panic("Error: unable to initialise localization")
}
```

3. Add validation message translations to your JSON files:

```json
{
  "validation.required": {
    "message": "The {field} field is absolutely required!"
  },
  "validation.email": {
    "message": "Please provide a valid email address for {field}"
  }
}
```

You can use the placeholders `{field}`, `{value}`, `{tag}`, and `{param}` in your messages, which will be replaced with the actual values when the error message is generated.

The validator automatically uses the language set in the context when available.

## Changes in Validator Behavior

### Empty Field Handling (Latest Change)

In the latest update, validation behavior for empty fields has been improved:

- **Previous behavior**: Custom validators would run on all fields, even if they were empty or not provided, potentially causing validation errors for optional fields.
- **New behavior**: Custom validators now automatically skip empty fields unless they're also marked as `required`.

For example:
```go
type User struct {
    // This will NOT error if phone is empty or not provided
    Phone string `json:"phone" validate:"phone"`

    // This will error if email is empty or not provided
    Email string `json:"email" validate:"required,email"`
}
```

This change makes validation behavior more intuitive and consistent with the standard validation behavior. You may need to update your struct tags if you were relying on validators running on empty fields.

### Validation of Structs and Custom Types

Several improvements have been made to how structs and custom types are validated:

1. **Non-Pointer Struct Validation**:
   - **Previous behavior**: Validation tags inside a non-pointer struct didn't work properly.
   - **New behavior**: Validation tags are consistently applied to fields inside both pointer and non-pointer structs.

2. **Validation Tags on Nested Types**:
   - **Previous behavior**: In pointer case Validation tags on custom types (structs, custom types, etc.) were ignored when the field was empty or not passed.
   - **New behavior**: Validation tags on custom types will be applied and then it would be applied on keys inside it as well.

Example of the new validation behavior:

```go
type Report struct {
    ID   string `json:"id" validate:"required"`
    Pipe Pipe   `json:"pipe" validate:"required"` // Struct level validation
}

type Pipe struct {
    Name string `json:"name" validate:"required"`
    Type string `json:"type" validate:"oneof=input output"`
}

// Case 1: Empty/Missing Pipe struct
report := Report{
    ID: "123",
    // Pipe is not provided
}
// This will now fail validation with error: "pipe is required"
// Previously, this would have passed validation

// Case 2: Empty Pipe struct
report := Report{
    ID:   "123",
    Pipe: Pipe{}, // Empty struct
}
// This will fail validation with error: "pipe is required"
// Previously, this would have passed

// Case 3: Pipe provided but with invalid Name
report := Report{
    ID: "123",
    Pipe: Pipe{
        Type: "input",
        // Name is missing
    },
}
// This will fail validation with error: "pipe.name is required"
// Previously, the nested validation would have been skipped

// Case 4: Valid data
report := Report{
    ID: "123",
    Pipe: Pipe{
        Name: "main-pipe",
        Type: "input",
    },
}
// This will pass validation
```

Key differences in behavior:
1. If a struct field has a `validate:"required"` tag, the struct itself must be provided and non-empty
2. Once a struct is provided, all validation rules inside that struct are checked
3. Validation is now consistent whether the struct is a pointer or not
4. Validation errors properly show the full path to the invalid field (e.g., "pipe.name is required")

### Third-Party Type Validation (e.g., guregu/null)

The validator now properly handles validation of third-party types like `null.Int`, `null.String`, etc.:

```go
import "gopkg.in/guregu/null.v4"

type Product struct {
    ID          string    `json:"id" validate:"required"`
    Name        string    `json:"name" validate:"required"`
    Stock       null.Int  `json:"stock" validate:"required"` // Validation on null.Int
    UpdatedAt   null.Time `json:"updated_at" validate:"required"`
}

// Case 1: Stock not provided
product := Product{
    ID:   "123",
    Name: "Test Product",
    // Stock is not provided
}
// This will now fail validation with error: "stock is required"
// Previously, this validation would have been skipped

// Case 2: Stock provided but null
product := Product{
    ID:    "123",
    Name:  "Test Product",
    Stock: null.Int{}, // or null.NewInt(0, false)
}
// This will pass now
// Previously also, this would have passed

// Case 3: We have a tag of gte as well in stock 
product := Product{
    ID:    "123",
    Name:  "Test Product",
    Stock: null.IntFrom(anyvalue), // Negative value
}
// This will panic now
// Previously also, this would have paniced 

// Case 4: Valid data
product := Product{
    ID:    "123",
    Name:  "Test Product",
    Stock: null.IntFrom(10),
}
// This will pass validation
```

Note that the validation behavior is now consistent across all field types, including:
- Built-in types (string, int, etc.)
- Custom structs (like the Pipe example above)
- Third-party types (null.Int, null.String, etc.)
- Slices and maps