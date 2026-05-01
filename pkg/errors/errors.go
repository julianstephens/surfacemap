package errors

import (
	"fmt"
	"strings"

	"github.com/julianstephens/surfacemap/pkg/model"
)

// FileError represents errors related to file system operations
// such as reading, writing, traversing directories, or discovering files.
type FileError struct {
	// Sentinel is the base error type (e.g., ErrDirWalk, ErrNoHCLFiles)
	Sentinel error
	// Cause is the underlying error that triggered this error
	Cause error
	// Operation is the operation being performed (e.g., "read", "traverse", "discover")
	Operation string
	// Path is the file or directory path where the error occurred
	Path string
}

func (e *FileError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: failed to %s path '%s': %v", e.Sentinel.Error(), e.Operation, e.Path, e.Cause)
	}
	return fmt.Sprintf("%s: failed to %s path '%s'", e.Sentinel.Error(), e.Operation, e.Path)
}

func (e *FileError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// ParseError represents errors during HCL/Terraform file parsing
// with specific file location information (file, line, column).
type ParseError struct {
	// Sentinel is the base error type (e.g., ErrParse, ErrInvalidHCLFiles)
	Sentinel error
	// Cause is the underlying error (e.g., HCL diagnostics)
	Cause error
	// FilePath is the path to the file being parsed
	FilePath string
	// Line is the line number where the error occurred (1-indexed)
	Line int
	// Column is the column number where the error occurred (1-indexed)
	Column int
	// Message provides additional context about the parse error
	Message string
}

func (e *ParseError) Error() string {
	location := fmt.Sprintf("%s:%d:%d", e.FilePath, e.Line, e.Column)
	if e.Message != "" {
		return fmt.Sprintf("%s at %s: %s", e.Sentinel.Error(), location, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s at %s: %v", e.Sentinel.Error(), location, e.Cause)
	}
	return fmt.Sprintf("%s at %s", e.Sentinel.Error(), location)
}

func (e *ParseError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// ExpressionError represents errors during HCL expression decoding
// with information about the expression being evaluated.
type ExpressionError struct {
	// Sentinel is the base error type (e.g., ErrDecodeExpression)
	Sentinel error
	// Cause is the underlying error
	Cause error
	// FilePath is the path to the file containing the expression
	FilePath string
	// Line is the line number where the expression starts (1-indexed)
	Line int
	// Column is the column number where the expression starts (1-indexed)
	Column int
	// Expression is the raw HCL expression text (if available)
	Expression string
	// Message provides additional context about the expression error
	Message string
}

func (e *ExpressionError) Error() string {
	location := fmt.Sprintf("%s:%d:%d", e.FilePath, e.Line, e.Column)
	
	var parts []string
	parts = append(parts, e.Sentinel.Error())
	parts = append(parts, fmt.Sprintf("at %s", location))
	
	if e.Expression != "" {
		parts = append(parts, fmt.Sprintf("in expression '%s'", e.Expression))
	}
	
	if e.Message != "" {
		parts = append(parts, e.Message)
	} else if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	
	return strings.Join(parts, ": ")
}

func (e *ExpressionError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// AttributeError represents errors related to missing or invalid resource attributes
// during resource extraction or validation.
type AttributeError struct {
	// Sentinel is the base error type (e.g., ErrMissingAttribute, ErrInvalidAttributeType)
	Sentinel error
	// Cause is the underlying error (if any)
	Cause error
	// ResourceID is the Terraform resource address (e.g., "aws_lambda_function.my_function")
	ResourceID string
	// AttributeName is the name of the attribute that caused the error
	AttributeName string
	// Expected describes what was expected (e.g., "string", "non-empty")
	Expected string
	// Found is the actual value or type found (if applicable)
	Found interface{}
	// Location is the file location where the resource is defined
	Location model.FileLocation
}

func (e *AttributeError) Error() string {
	location := fmt.Sprintf("%s:%d:%d", e.Location.File, e.Location.Line, e.Location.Column)
	
	var parts []string
	parts = append(parts, e.Sentinel.Error())
	parts = append(parts, fmt.Sprintf("attribute '%s'", e.AttributeName))
	parts = append(parts, fmt.Sprintf("in resource '%s'", e.ResourceID))
	parts = append(parts, fmt.Sprintf("at %s", location))
	
	if e.Expected != "" && e.Found != nil {
		parts = append(parts, fmt.Sprintf("expected %s, found %T", e.Expected, e.Found))
	} else if e.Expected != "" {
		parts = append(parts, fmt.Sprintf("expected %s", e.Expected))
	}
	
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	
	return strings.Join(parts, ": ")
}

func (e *AttributeError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// ModuleError represents errors during Terraform module resolution,
// loading, or graph building operations.
type ModuleError struct {
	// Sentinel is the base error type (e.g., ErrModuleNotFound, ErrCyclicDependency)
	Sentinel error
	// Cause is the underlying error
	Cause error
	// ParentPath is the path to the parent module
	ParentPath string
	// ModuleSource is the Terraform module source string
	ModuleSource string
	// Operation is the operation being performed (e.g., "resolve", "load", "build graph")
	Operation string
	// AdditionalInfo provides extra context (e.g., cycle path)
	AdditionalInfo string
}

func (e *ModuleError) Error() string {
	var parts []string
	parts = append(parts, e.Sentinel.Error())
	
	if e.Operation != "" {
		parts = append(parts, fmt.Sprintf("during %s", e.Operation))
	}
	
	if e.ModuleSource != "" {
		parts = append(parts, fmt.Sprintf("module source '%s'", e.ModuleSource))
	}
	
	if e.ParentPath != "" {
		parts = append(parts, fmt.Sprintf("from parent '%s'", e.ParentPath))
	}
	
	if e.AdditionalInfo != "" {
		parts = append(parts, e.AdditionalInfo)
	}
	
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	
	return strings.Join(parts, ": ")
}

func (e *ModuleError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// ExtractionError represents errors during resource extraction from
// Terraform resources to typed AWS models.
type ExtractionError struct {
	// Sentinel is the base error type (e.g., ErrUnsupportedResourceType, ErrExtractionFailed)
	Sentinel error
	// Cause is the underlying error
	Cause error
	// ResourceID is the Terraform resource address
	ResourceID string
	// ResourceType is the Terraform resource type (e.g., "aws_lambda_function")
	ResourceType string
	// Location is the file location where the resource is defined
	Location model.FileLocation
	// Reason provides additional context about why extraction failed
	Reason string
}

func (e *ExtractionError) Error() string {
	location := fmt.Sprintf("%s:%d:%d", e.Location.File, e.Location.Line, e.Location.Column)
	
	var parts []string
	parts = append(parts, e.Sentinel.Error())
	parts = append(parts, fmt.Sprintf("resource '%s'", e.ResourceID))
	parts = append(parts, fmt.Sprintf("type '%s'", e.ResourceType))
	parts = append(parts, fmt.Sprintf("at %s", location))
	
	if e.Reason != "" {
		parts = append(parts, e.Reason)
	}
	
	if e.Cause != nil {
		parts = append(parts, e.Cause.Error())
	}
	
	return strings.Join(parts, ": ")
}

func (e *ExtractionError) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	}
	return e.Sentinel
}

// MultiError aggregates multiple errors from batch operations
// such as parsing multiple files or extracting multiple resources.
type MultiError struct {
	// Errors is the list of errors that occurred
	Errors []error
	// Operation describes the batch operation (e.g., "parsing files", "extracting resources")
	Operation string
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return "no errors"
	}
	
	if len(e.Errors) == 1 {
		if e.Operation != "" {
			return fmt.Sprintf("%s: %v", e.Operation, e.Errors[0])
		}
		return e.Errors[0].Error()
	}
	
	var sb strings.Builder
	if e.Operation != "" {
		sb.WriteString(fmt.Sprintf("%s: %d errors occurred:\n", e.Operation, len(e.Errors)))
	} else {
		sb.WriteString(fmt.Sprintf("%d errors occurred:\n", len(e.Errors)))
	}
	
	for i, err := range e.Errors {
		sb.WriteString(fmt.Sprintf("  %d. %v\n", i+1, err))
	}
	
	return sb.String()
}

func (e *MultiError) Unwrap() []error {
	return e.Errors
}

// IsMultiError returns true if the error is a MultiError
func (e *MultiError) Is(target error) bool {
	_, ok := target.(*MultiError)
	return ok
}

// Add appends an error to the MultiError. If err is nil, it's a no-op.
func (e *MultiError) Add(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

// HasErrors returns true if the MultiError contains any errors
func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

// ErrorOrNil returns the MultiError if it contains errors, otherwise nil
func (e *MultiError) ErrorOrNil() error {
	if e.HasErrors() {
		return e
	}
	return nil
}
