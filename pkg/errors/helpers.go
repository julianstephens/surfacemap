package errors

import (
	"errors"
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/julianstephens/surfacemap/pkg/model"
)

// NewFileError creates a new FileError with the given parameters
func NewFileError(sentinel error, operation, path string, cause error) *FileError {
	return &FileError{
		Sentinel:  sentinel,
		Cause:     cause,
		Operation: operation,
		Path:      path,
	}
}

// NewParseError creates a new ParseError with the given parameters
func NewParseError(sentinel error, filePath string, line, column int, message string, cause error) *ParseError {
	return &ParseError{
		Sentinel: sentinel,
		Cause:    cause,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Message:  message,
	}
}

// NewParseErrorFromDiagnostics creates a ParseError from HCL diagnostics
func NewParseErrorFromDiagnostics(sentinel error, filePath string, diags hcl.Diagnostics) *ParseError {
	if len(diags) == 0 {
		return &ParseError{
			Sentinel: sentinel,
			FilePath: filePath,
			Message:  "unknown parse error",
		}
	}

	// Use the first diagnostic for location and message
	diag := diags[0]
	line, column := 0, 0
	if diag.Subject != nil {
		line = diag.Subject.Start.Line
		column = diag.Subject.Start.Column
	}

	return &ParseError{
		Sentinel: sentinel,
		FilePath: filePath,
		Line:     line,
		Column:   column,
		Message:  diags.Error(),
		Cause:    errors.New(diags.Error()),
	}
}

// NewExpressionError creates a new ExpressionError with the given parameters
func NewExpressionError(sentinel error, filePath string, line, column int, expression, message string, cause error) *ExpressionError {
	return &ExpressionError{
		Sentinel:   sentinel,
		Cause:      cause,
		FilePath:   filePath,
		Line:       line,
		Column:     column,
		Expression: expression,
		Message:    message,
	}
}

// NewExpressionErrorFromRange creates an ExpressionError from an HCL Range
func NewExpressionErrorFromRange(sentinel error, filePath string, hclRange hcl.Range, expression, message string, cause error) *ExpressionError {
	return &ExpressionError{
		Sentinel:   sentinel,
		Cause:      cause,
		FilePath:   filePath,
		Line:       hclRange.Start.Line,
		Column:     hclRange.Start.Column,
		Expression: expression,
		Message:    message,
	}
}

// NewAttributeError creates a new AttributeError with the given parameters
func NewAttributeError(sentinel error, resourceID, attributeName string, location model.FileLocation, cause error) *AttributeError {
	return &AttributeError{
		Sentinel:      sentinel,
		Cause:         cause,
		ResourceID:    resourceID,
		AttributeName: attributeName,
		Location:      location,
	}
}

// NewAttributeErrorWithType creates an AttributeError with expected and found types
func NewAttributeErrorWithType(sentinel error, resourceID, attributeName, expected string, found interface{}, location model.FileLocation) *AttributeError {
	return &AttributeError{
		Sentinel:      sentinel,
		ResourceID:    resourceID,
		AttributeName: attributeName,
		Expected:      expected,
		Found:         found,
		Location:      location,
	}
}

// NewModuleError creates a new ModuleError with the given parameters
func NewModuleError(sentinel error, operation, parentPath, moduleSource string, cause error) *ModuleError {
	return &ModuleError{
		Sentinel:     sentinel,
		Cause:        cause,
		ParentPath:   parentPath,
		ModuleSource: moduleSource,
		Operation:    operation,
	}
}

// NewModuleErrorWithInfo creates a ModuleError with additional information
func NewModuleErrorWithInfo(sentinel error, operation, parentPath, moduleSource, additionalInfo string, cause error) *ModuleError {
	return &ModuleError{
		Sentinel:       sentinel,
		Cause:          cause,
		ParentPath:     parentPath,
		ModuleSource:   moduleSource,
		Operation:      operation,
		AdditionalInfo: additionalInfo,
	}
}

// NewExtractionError creates a new ExtractionError with the given parameters
func NewExtractionError(sentinel error, resourceID, resourceType string, location model.FileLocation, reason string, cause error) *ExtractionError {
	return &ExtractionError{
		Sentinel:     sentinel,
		Cause:        cause,
		ResourceID:   resourceID,
		ResourceType: resourceType,
		Location:     location,
		Reason:       reason,
	}
}

// NewMultiError creates a new MultiError with the given operation
func NewMultiError(operation string) *MultiError {
	return &MultiError{
		Operation: operation,
		Errors:    make([]error, 0),
	}
}

// NewMultiErrorWithErrors creates a MultiError with initial errors
func NewMultiErrorWithErrors(operation string, errs []error) *MultiError {
	return &MultiError{
		Operation: operation,
		Errors:    errs,
	}
}

// IsFileError returns true if the error is or wraps a FileError
func IsFileError(err error) bool {
	var fileErr *FileError
	return errors.As(err, &fileErr)
}

// IsParseError returns true if the error is or wraps a ParseError
func IsParseError(err error) bool {
	var parseErr *ParseError
	return errors.As(err, &parseErr)
}

// IsExpressionError returns true if the error is or wraps an ExpressionError
func IsExpressionError(err error) bool {
	var exprErr *ExpressionError
	return errors.As(err, &exprErr)
}

// IsAttributeError returns true if the error is or wraps an AttributeError
func IsAttributeError(err error) bool {
	var attrErr *AttributeError
	return errors.As(err, &attrErr)
}

// IsModuleError returns true if the error is or wraps a ModuleError
func IsModuleError(err error) bool {
	var modErr *ModuleError
	return errors.As(err, &modErr)
}

// IsExtractionError returns true if the error is or wraps an ExtractionError
func IsExtractionError(err error) bool {
	var extErr *ExtractionError
	return errors.As(err, &extErr)
}

// IsMultiError returns true if the error is or wraps a MultiError
func IsMultiError(err error) bool {
	var multiErr *MultiError
	return errors.As(err, &multiErr)
}

// HasSentinel checks if the error chain contains the given sentinel error
func HasSentinel(err, sentinel error) bool {
	return errors.Is(err, sentinel)
}

// GetFileError extracts a FileError from the error chain, if present
func GetFileError(err error) (*FileError, bool) {
	var fileErr *FileError
	ok := errors.As(err, &fileErr)
	return fileErr, ok
}

// GetParseError extracts a ParseError from the error chain, if present
func GetParseError(err error) (*ParseError, bool) {
	var parseErr *ParseError
	ok := errors.As(err, &parseErr)
	return parseErr, ok
}

// GetAttributeError extracts an AttributeError from the error chain, if present
func GetAttributeError(err error) (*AttributeError, bool) {
	var attrErr *AttributeError
	ok := errors.As(err, &attrErr)
	return attrErr, ok
}

// GetModuleError extracts a ModuleError from the error chain, if present
func GetModuleError(err error) (*ModuleError, bool) {
	var modErr *ModuleError
	ok := errors.As(err, &modErr)
	return modErr, ok
}

// GetExtractionError extracts an ExtractionError from the error chain, if present
func GetExtractionError(err error) (*ExtractionError, bool) {
	var extErr *ExtractionError
	ok := errors.As(err, &extErr)
	return extErr, ok
}

// GetMultiError extracts a MultiError from the error chain, if present
func GetMultiError(err error) (*MultiError, bool) {
	var multiErr *MultiError
	ok := errors.As(err, &multiErr)
	return multiErr, ok
}

// WrapWithContext wraps an error with additional context message
func WrapWithContext(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}
