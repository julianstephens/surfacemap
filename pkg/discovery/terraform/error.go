package terraform

import (
	"errors"
)

const IDENTIFIER string = "terraformparser"

var (
	ErrDirWalk          = errors.New(IDENTIFIER + ": unable to traverse the provided root directory")
	ErrNoHCLFiles       = errors.New(IDENTIFIER + ": no valid HCL files found in directory")
	ErrInvalidHCLFiles  = errors.New(IDENTIFIER + ": one or more HCL files could not be parsed")
	ErrParse            = errors.New(IDENTIFIER + ": error parsing HCL files")
	ErrDecodeExpression = errors.New(IDENTIFIER + ": unable to decode expression")
)

type TerraformParserError struct {
	Err       error
	Cause     error
	RootDir   string
	ModuleDir string
}

func (e *TerraformParserError) Error() string { return e.Err.Error() }

func (e *TerraformParserError) Unwrap() error { return e.Err }
