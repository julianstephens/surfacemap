package terraform

import (
	"errors"
	"fmt"
)

const IDENTIFIER string = "terraformparser"

var (
	ErrDirWalk          = errors.New(IDENTIFIER + ": unable to traverse the provided root directory")
	ErrNoHCLFiles       = errors.New(IDENTIFIER + ": no valid HCL files found in directory")
	ErrInvalidHCLFiles  = errors.New(IDENTIFIER + ": one or more HCL files could not be parsed")
	ErrParse            = errors.New(IDENTIFIER + ": error parsing HCL files")
	ErrDecodeExpression = errors.New(IDENTIFIER + ": unable to decode expression")
)

// TerraformParserError is deprecated in favor of the rich error types in pkg/errors.
// Use FileError, ParseError, ExpressionError, etc. instead.
// This type is kept for backward compatibility during migration.
//
// Deprecated: Use pkg/errors types instead.
type TerraformParserError struct {
	Err       error
	Cause     error
	RootDir   string
	ModuleDir string
}

func (e *TerraformParserError) Error() string {
	msg := e.Err.Error()
	
	if e.RootDir != "" && e.ModuleDir != "" {
		msg = fmt.Sprintf("%s (root: %s, module: %s)", msg, e.RootDir, e.ModuleDir)
	} else if e.RootDir != "" {
		msg = fmt.Sprintf("%s (root: %s)", msg, e.RootDir)
	} else if e.ModuleDir != "" {
		msg = fmt.Sprintf("%s (module: %s)", msg, e.ModuleDir)
	}
	
	if e.Cause != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Cause)
	}
	
	return msg
}

func (e *TerraformParserError) Unwrap() error { return e.Err }
