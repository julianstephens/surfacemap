package errors

import (
	"errors"
)

const identifier = "surfacemap"

// File operation errors
var (
	// ErrDirWalk indicates an error occurred while traversing a directory tree
	ErrDirWalk = errors.New(identifier + ": unable to traverse directory")
	
	// ErrNoHCLFiles indicates no valid HCL files were found in the specified directory
	ErrNoHCLFiles = errors.New(identifier + ": no valid HCL files found")
	
	// ErrFileRead indicates an error occurred while reading a file
	ErrFileRead = errors.New(identifier + ": unable to read file")
	
	// ErrFileWrite indicates an error occurred while writing a file
	ErrFileWrite = errors.New(identifier + ": unable to write file")
)

// Parse errors
var (
	// ErrParse indicates an error occurred during HCL/Terraform parsing
	ErrParse = errors.New(identifier + ": error parsing HCL")
	
	// ErrInvalidHCLFiles indicates one or more HCL files could not be parsed
	ErrInvalidHCLFiles = errors.New(identifier + ": invalid HCL files")
	
	// ErrInvalidSyntax indicates the HCL syntax is invalid
	ErrInvalidSyntax = errors.New(identifier + ": invalid HCL syntax")
)

// Expression errors
var (
	// ErrDecodeExpression indicates an error occurred while decoding an HCL expression
	ErrDecodeExpression = errors.New(identifier + ": unable to decode expression")
	
	// ErrDecodeBlock indicates an error occurred while decoding an HCL block
	ErrDecodeBlock = errors.New(identifier + ": unable to decode block")
	
	// ErrInvalidExpression indicates the HCL expression is invalid or malformed
	ErrInvalidExpression = errors.New(identifier + ": invalid expression")
	
	// ErrUnknownVariable indicates a variable reference could not be resolved
	ErrUnknownVariable = errors.New(identifier + ": unknown variable")
)

// Module errors
var (
	// ErrModuleNotFound indicates a Terraform module could not be located
	ErrModuleNotFound = errors.New(identifier + ": module not found")
	
	// ErrModuleLoad indicates an error occurred while loading a module
	ErrModuleLoad = errors.New(identifier + ": unable to load module")
	
	// ErrModuleResolve indicates an error occurred while resolving a module path
	ErrModuleResolve = errors.New(identifier + ": unable to resolve module")
	
	// ErrCyclicDependency indicates a circular dependency was detected in modules
	ErrCyclicDependency = errors.New(identifier + ": cyclic module dependency detected")
	
	// ErrMaxDepthExceeded indicates module nesting exceeded the maximum allowed depth
	ErrMaxDepthExceeded = errors.New(identifier + ": maximum module depth exceeded")
)

// Attribute errors
var (
	// ErrMissingAttribute indicates a required attribute is missing from a resource
	ErrMissingAttribute = errors.New(identifier + ": missing required attribute")
	
	// ErrInvalidAttributeType indicates an attribute has an unexpected type
	ErrInvalidAttributeType = errors.New(identifier + ": invalid attribute type")
	
	// ErrInvalidAttributeValue indicates an attribute has an invalid value
	ErrInvalidAttributeValue = errors.New(identifier + ": invalid attribute value")
)

// Extraction errors
var (
	// ErrUnsupportedResourceType indicates the resource type is not supported for extraction
	ErrUnsupportedResourceType = errors.New(identifier + ": unsupported resource type")
	
	// ErrExtractionFailed indicates resource extraction failed
	ErrExtractionFailed = errors.New(identifier + ": resource extraction failed")
)

// Validation errors
var (
	// ErrValidationFailed indicates resource validation failed
	ErrValidationFailed = errors.New(identifier + ": validation failed")
	
	// ErrInvalidConfiguration indicates the configuration is invalid
	ErrInvalidConfiguration = errors.New(identifier + ": invalid configuration")
)

// General errors
var (
	// ErrNotImplemented indicates functionality is not yet implemented
	ErrNotImplemented = errors.New(identifier + ": not yet implemented")
	
	// ErrInternal indicates an internal error occurred
	ErrInternal = errors.New(identifier + ": internal error")
)
