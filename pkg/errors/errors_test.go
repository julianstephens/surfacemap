package errors

import (
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/julianstephens/surfacemap/pkg/model"
)

func TestFileError(t *testing.T) {
	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("permission denied")
		err := NewFileError(ErrDirWalk, "traverse", "/some/path", cause)
		
		// Test Error() includes all context
		errMsg := err.Error()
		if !strings.Contains(errMsg, "traverse") {
			t.Errorf("error message should contain operation: %s", errMsg)
		}
		if !strings.Contains(errMsg, "/some/path") {
			t.Errorf("error message should contain path: %s", errMsg)
		}
		if !strings.Contains(errMsg, "permission denied") {
			t.Errorf("error message should contain cause: %s", errMsg)
		}
		
		// Test Unwrap()
		if !errors.Is(err, cause) {
			t.Error("Unwrap should return cause")
		}
	})
	
	t.Run("without cause", func(t *testing.T) {
		err := NewFileError(ErrNoHCLFiles, "discover", "/some/path", nil)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "discover") {
			t.Errorf("error message should contain operation: %s", errMsg)
		}
		if !strings.Contains(errMsg, "/some/path") {
			t.Errorf("error message should contain path: %s", errMsg)
		}
		
		// Test Unwrap() returns sentinel
		if !errors.Is(err, ErrNoHCLFiles) {
			t.Error("Unwrap should return sentinel")
		}
	})
}

func TestParseError(t *testing.T) {
	t.Run("with all fields", func(t *testing.T) {
		err := NewParseError(ErrParse, "/path/to/file.tf", 10, 5, "unexpected token", nil)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "/path/to/file.tf:10:5") {
			t.Errorf("error message should contain location: %s", errMsg)
		}
		if !strings.Contains(errMsg, "unexpected token") {
			t.Errorf("error message should contain message: %s", errMsg)
		}
	})
	
	t.Run("from HCL diagnostics", func(t *testing.T) {
		diags := hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Test error",
				Subject: &hcl.Range{
					Filename: "test.tf",
					Start:    hcl.Pos{Line: 5, Column: 2},
				},
			},
		}
		
		err := NewParseErrorFromDiagnostics(ErrInvalidHCLFiles, "/path/to/test.tf", diags)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "5") {
			t.Errorf("error message should contain line number: %s", errMsg)
		}
	})
}

func TestExpressionError(t *testing.T) {
	t.Run("with expression text", func(t *testing.T) {
		err := NewExpressionError(
			ErrDecodeExpression,
			"/path/to/file.tf",
			10, 5,
			"var.foo",
			"unknown variable",
			nil,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "var.foo") {
			t.Errorf("error message should contain expression: %s", errMsg)
		}
		if !strings.Contains(errMsg, "unknown variable") {
			t.Errorf("error message should contain message: %s", errMsg)
		}
		if !strings.Contains(errMsg, "10:5") {
			t.Errorf("error message should contain location: %s", errMsg)
		}
	})
	
	t.Run("from HCL Range", func(t *testing.T) {
		hclRange := hcl.Range{
			Filename: "test.tf",
			Start:    hcl.Pos{Line: 5, Column: 2},
			End:      hcl.Pos{Line: 5, Column: 10},
		}
		
		err := NewExpressionErrorFromRange(
			ErrInvalidExpression,
			"/path/to/test.tf",
			hclRange,
			"foo[\"bar\"]",
			"invalid index",
			nil,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "5:2") {
			t.Errorf("error message should contain location from range: %s", errMsg)
		}
	})
}

func TestAttributeError(t *testing.T) {
	location := model.FileLocation{
		File:   "main.tf",
		Line:   10,
		Column: 2,
	}
	
	t.Run("with type information", func(t *testing.T) {
		err := NewAttributeErrorWithType(
			ErrMissingAttribute,
			"aws_lambda_function.my_func",
			"function_name",
			"string",
			12345, // found int instead
			location,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "function_name") {
			t.Errorf("error message should contain attribute name: %s", errMsg)
		}
		if !strings.Contains(errMsg, "aws_lambda_function.my_func") {
			t.Errorf("error message should contain resource ID: %s", errMsg)
		}
		if !strings.Contains(errMsg, "string") {
			t.Errorf("error message should contain expected type: %s", errMsg)
		}
		if !strings.Contains(errMsg, "int") {
			t.Errorf("error message should contain found type: %s", errMsg)
		}
		if !strings.Contains(errMsg, "main.tf:10:2") {
			t.Errorf("error message should contain location: %s", errMsg)
		}
	})
	
	t.Run("without type information", func(t *testing.T) {
		err := NewAttributeError(
			ErrMissingAttribute,
			"aws_lambda_function.my_func",
			"role",
			location,
			nil,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "role") {
			t.Errorf("error message should contain attribute name: %s", errMsg)
		}
	})
}

func TestModuleError(t *testing.T) {
	t.Run("basic module error", func(t *testing.T) {
		err := NewModuleError(
			ErrModuleNotFound,
			"resolve",
			"/path/to/parent",
			"../child",
			nil,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "resolve") {
			t.Errorf("error message should contain operation: %s", errMsg)
		}
		if !strings.Contains(errMsg, "../child") {
			t.Errorf("error message should contain module source: %s", errMsg)
		}
		if !strings.Contains(errMsg, "/path/to/parent") {
			t.Errorf("error message should contain parent path: %s", errMsg)
		}
	})
	
	t.Run("with additional info", func(t *testing.T) {
		err := NewModuleErrorWithInfo(
			ErrCyclicDependency,
			"build graph",
			"/path/to/parent",
			"./child",
			"cycle: A -> B -> A",
			nil,
		)
		
		errMsg := err.Error()
		if !strings.Contains(errMsg, "cycle: A -> B -> A") {
			t.Errorf("error message should contain additional info: %s", errMsg)
		}
	})
}

func TestExtractionError(t *testing.T) {
	location := model.FileLocation{
		File:   "resources.tf",
		Line:   25,
		Column: 1,
	}
	
	err := NewExtractionError(
		ErrUnsupportedResourceType,
		"aws_some_resource.example",
		"aws_some_resource",
		location,
		"resource type not yet supported",
		nil,
	)
	
	errMsg := err.Error()
	if !strings.Contains(errMsg, "aws_some_resource.example") {
		t.Errorf("error message should contain resource ID: %s", errMsg)
	}
	if !strings.Contains(errMsg, "aws_some_resource") {
		t.Errorf("error message should contain resource type: %s", errMsg)
	}
	if !strings.Contains(errMsg, "resources.tf:25:1") {
		t.Errorf("error message should contain location: %s", errMsg)
	}
	if !strings.Contains(errMsg, "not yet supported") {
		t.Errorf("error message should contain reason: %s", errMsg)
	}
}

func TestMultiError(t *testing.T) {
	t.Run("empty multi error", func(t *testing.T) {
		multi := NewMultiError("test operation")
		
		if multi.HasErrors() {
			t.Error("empty multi error should not have errors")
		}
		
		if err := multi.ErrorOrNil(); err != nil {
			t.Error("ErrorOrNil should return nil for empty multi error")
		}
	})
	
	t.Run("single error", func(t *testing.T) {
		multi := NewMultiError("test operation")
		multi.Add(errors.New("error 1"))
		
		if !multi.HasErrors() {
			t.Error("multi error with errors should return true for HasErrors")
		}
		
		errMsg := multi.Error()
		if !strings.Contains(errMsg, "error 1") {
			t.Errorf("error message should contain the error: %s", errMsg)
		}
	})
	
	t.Run("multiple errors", func(t *testing.T) {
		multi := NewMultiError("parsing files")
		multi.Add(errors.New("error 1"))
		multi.Add(errors.New("error 2"))
		multi.Add(errors.New("error 3"))
		
		if len(multi.Errors) != 3 {
			t.Errorf("expected 3 errors, got %d", len(multi.Errors))
		}
		
		errMsg := multi.Error()
		if !strings.Contains(errMsg, "3 errors occurred") {
			t.Errorf("error message should mention count: %s", errMsg)
		}
		if !strings.Contains(errMsg, "error 1") || !strings.Contains(errMsg, "error 2") {
			t.Errorf("error message should contain all errors: %s", errMsg)
		}
	})
	
	t.Run("ignores nil errors", func(t *testing.T) {
		multi := NewMultiError("test")
		multi.Add(nil)
		multi.Add(errors.New("real error"))
		multi.Add(nil)
		
		if len(multi.Errors) != 1 {
			t.Errorf("expected 1 error (nil ignored), got %d", len(multi.Errors))
		}
	})
	
	t.Run("Unwrap returns slice", func(t *testing.T) {
		multi := NewMultiError("test")
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		multi.Add(err1)
		multi.Add(err2)
		
		unwrapped := multi.Unwrap()
		if len(unwrapped) != 2 {
			t.Errorf("expected 2 unwrapped errors, got %d", len(unwrapped))
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("IsFileError", func(t *testing.T) {
		fileErr := NewFileError(ErrDirWalk, "traverse", "/path", nil)
		if !IsFileError(fileErr) {
			t.Error("IsFileError should return true for FileError")
		}
		
		otherErr := errors.New("some error")
		if IsFileError(otherErr) {
			t.Error("IsFileError should return false for non-FileError")
		}
	})
	
	t.Run("HasSentinel", func(t *testing.T) {
		fileErr := NewFileError(ErrDirWalk, "traverse", "/path", nil)
		if !HasSentinel(fileErr, ErrDirWalk) {
			t.Error("HasSentinel should find the sentinel error")
		}
		
		if HasSentinel(fileErr, ErrNoHCLFiles) {
			t.Error("HasSentinel should not find wrong sentinel")
		}
	})
	
	t.Run("GetFileError", func(t *testing.T) {
		originalErr := NewFileError(ErrDirWalk, "traverse", "/my/path", nil)
		
		fileErr, ok := GetFileError(originalErr)
		if !ok {
			t.Error("GetFileError should successfully extract FileError")
		}
		if fileErr.Path != "/my/path" {
			t.Errorf("extracted FileError should have correct path, got %s", fileErr.Path)
		}
	})
	
	t.Run("WrapWithContext", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrapped := WrapWithContext(originalErr, "context: %s", "some context")
		
		if !errors.Is(wrapped, originalErr) {
			t.Error("WrapWithContext should preserve error chain")
		}
		
		wrappedMsg := wrapped.Error()
		if !strings.Contains(wrappedMsg, "context: some context") {
			t.Errorf("wrapped error should contain context: %s", wrappedMsg)
		}
		if !strings.Contains(wrappedMsg, "original error") {
			t.Errorf("wrapped error should contain original: %s", wrappedMsg)
		}
	})
	
	t.Run("WrapWithContext handles nil", func(t *testing.T) {
		wrapped := WrapWithContext(nil, "context")
		if wrapped != nil {
			t.Error("WrapWithContext should return nil for nil error")
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	sentinels := []error{
		ErrDirWalk,
		ErrNoHCLFiles,
		ErrFileRead,
		ErrFileWrite,
		ErrParse,
		ErrInvalidHCLFiles,
		ErrInvalidSyntax,
		ErrDecodeExpression,
		ErrDecodeBlock,
		ErrInvalidExpression,
		ErrUnknownVariable,
		ErrModuleNotFound,
		ErrModuleLoad,
		ErrModuleResolve,
		ErrCyclicDependency,
		ErrMaxDepthExceeded,
		ErrMissingAttribute,
		ErrInvalidAttributeType,
		ErrInvalidAttributeValue,
		ErrUnsupportedResourceType,
		ErrExtractionFailed,
		ErrValidationFailed,
		ErrInvalidConfiguration,
		ErrNotImplemented,
		ErrInternal,
	}
	
	// Verify all sentinels are distinct
	seen := make(map[string]bool)
	for _, sentinel := range sentinels {
		msg := sentinel.Error()
		if seen[msg] {
			t.Errorf("duplicate sentinel error message: %s", msg)
		}
		seen[msg] = true
	}
	
	// Verify all contain identifier
	for _, sentinel := range sentinels {
		if !strings.Contains(sentinel.Error(), "surfacemap") {
			t.Errorf("sentinel %v should contain identifier", sentinel)
		}
	}
}
