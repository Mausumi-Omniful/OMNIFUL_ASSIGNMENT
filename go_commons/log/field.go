package log

import (
	"time"

	"go.uber.org/zap"
)

type Field = zap.Field

// Int creates a Field with the given key and integer value.
// It is used to add structured logging fields to log entries.
//
// Parameters:
//   - key: The name of the field.
//   - val: The integer value to associate with the field.
//
// Returns:
//
//	A Field that can be used with a logger to include the key-value pair in log output.
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// String creates a new Field with a string value.
// It takes a key and a string value as arguments and returns a Field
// that can be used for structured logging.
//
// Parameters:
//   - key: The name of the field.
//   - val: The string value to associate with the key.
//
// Returns:
//
//	A Field containing the key-value pair.
func String(key string, val string) Field {
	return zap.String(key, val)
}

// Duration creates a Field that represents a duration value.
// The key parameter specifies the name of the field, and the val parameter
// is the duration value to be logged. This function wraps the zap.Duration
// method to provide a convenient way to log duration fields.
//
// Parameters:
//   - key: The name of the field.
//   - val: The duration value to be logged.
//
// Returns:
//
//	A Field representing the duration value.
func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

// ErrorField creates a new logging field that represents an error.
// It takes an error value as input and returns a Field that can be
// used with structured logging to include error details in log entries.
//
// Parameters:
//   - val: The error value to be logged.
//
// Returns:
//
//	A Field containing the error information.
func ErrorField(val error) Field {
	return zap.Error(val)
}

// Any creates a new Field with the given key and value.
// The value can be of any type, and it will be serialized
// appropriately by the underlying zap.Any function.
// This is useful for logging arbitrary data.
//
// Parameters:
//   - key: The name of the field.
//   - val: The value to associate with the field.
//
// Returns:
//
//	A Field that can be used in structured logging.
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

// Bool creates a Field with a boolean value.
// The Field can be used to add a key-value pair to a structured log entry.
//
// Parameters:
//   - key: The name of the field.
//   - val: The boolean value to associate with the key.
//
// Returns:
//
//	A Field containing the key and boolean value.
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// Float64 creates a Field that carries a float64 value with the specified key.
// This is typically used for structured logging, allowing you to associate
// a floating-point number with a descriptive key in the log entry.
//
// Parameters:
//   - key: The name of the field to be logged.
//   - val: The float64 value to associate with the key.
//
// Returns:
//
//	A Field that can be used with a logger to include the key-value pair in a log entry.
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Int64 creates a Field with the given key and int64 value.
// It is used to add structured logging fields to log entries.
//
// Parameters:
//   - key: The name of the field.
//   - val: The int64 value to associate with the key.
//
// Returns:
//
//	A Field that can be used with a logger for structured logging.
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}
