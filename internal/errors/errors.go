package errors

import "fmt"

// ValidationError represents an error that occurs during input validation.
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

func NewValidationError(field, message string, err error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Err:     err,
	}
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// OperationError represents an error that occurs during an operation.
type OperationError struct {
	Operation string
	Context   string
	Err       error
}

func NewOperationError(operation, context string, err error) *OperationError {
	return &OperationError{
		Operation: operation,
		Context:   context,
		Err:       err,
	}
}

func (e *OperationError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s failed for %s: %v", e.Operation, e.Context, e.Err)
	}
	return fmt.Sprintf("%s failed: %v", e.Operation, e.Err)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

// ConfigError represents an error that occurs during configuration loading.
type ConfigError struct {
	Source  string
	Message string
	Err     error
}

func NewConfigError(source, message string, err error) *ConfigError {
	return &ConfigError{
		Source:  source,
		Message: message,
		Err:     err,
	}
}

func (e *ConfigError) Error() string {
	if e.Source != "" {
		return fmt.Sprintf("config error in %s: %s", e.Source, e.Message)
	}
	return fmt.Sprintf("config error: %s", e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

// RepositoryError represents an error that occurs during git operations.
type RepositoryError struct {
	Operation string
	Err       error
}

func NewRepositoryError(operation string, err error) *RepositoryError {
	return &RepositoryError{
		Operation: operation,
		Err:       err,
	}
}

func (e *RepositoryError) Error() string {
	return fmt.Sprintf("repository %s failed: %v", e.Operation, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// IsValidationError checks if the error is a ValidationError.
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsOperationError checks if the error is an OperationError.
func IsOperationError(err error) bool {
	_, ok := err.(*OperationError)
	return ok
}

// IsConfigError checks if the error is a ConfigError.
func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}

// IsRepositoryError checks if the error is a RepositoryError.
func IsRepositoryError(err error) bool {
	_, ok := err.(*RepositoryError)
	return ok
}
