package extractors_test

import (
	"errors"
	"testing"

	"github.com/julianstephens/surfacemap/pkg/discovery/terraform/extractors"
	pkgerrors "github.com/julianstephens/surfacemap/pkg/errors"
	"github.com/julianstephens/surfacemap/pkg/model"
)

func TestLambdaExtractor_Extract(t *testing.T) {
	// Helper variables for nil slices (matching implementation behavior)
	nilArchitectures := []string(nil)
	nilSubnetIds := []string(nil)
	nilSecurityGroupIds := []string(nil)

	tests := []struct {
		name     string
		resource model.Resource
		want     *model.LambdaFunction
		wantErr  bool
		sentinel error
	}{
		{
			name: "successful extraction with required fields only",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function",
					"role":          "arn:aws:iam::123456789012:role/test-role",
				},
				Blocks: map[string][]map[string]any{},
				Location: model.FileLocation{
					File:   "main.tf",
					Line:   1,
					Column: 0,
				},
			},
			want: &model.LambdaFunction{
				Resource: model.Resource{
					ID:   "aws_lambda_function.test",
					Type: "aws_lambda_function",
					Attributes: map[string]any{
						"function_name": "test-function",
						"role":          "arn:aws:iam::123456789012:role/test-role",
					},
					Blocks: map[string][]map[string]any{},
					Location: model.FileLocation{
						File:   "main.tf",
						Line:   1,
						Column: 0,
					},
				},
				FunctionName:   "test-function",
				RoleArn:        "arn:aws:iam::123456789012:role/test-role",
				Handler:        stringPtr(""),
				Runtime:        stringPtr(""),
				SourceCodePath: stringPtr(""),
				Architectures:  &nilArchitectures,
				Environment: &struct {
					Variables map[string]string "json:\"variables\""
				}{
					Variables: nil,
				},
				VPCConfig: nil,
			},
			wantErr: false,
		},
		{
			name: "successful extraction with all optional fields",
			resource: model.Resource{
				ID:   "aws_lambda_function.test_full",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function-full",
					"role":          "arn:aws:iam::123456789012:role/test-role",
					"handler":       "index.handler",
					"runtime":       "nodejs18.x",
					"filename":      "./lambda.zip",
					"architectures": []string{"x86_64"},
				},
				Blocks: map[string][]map[string]any{
					"environment": {
						{
							"variables": map[string]string{
								"ENV": "production",
								"KEY": "value",
							},
						},
					},
					"vpc_config": {
						{
							"subnet_ids":         []string{"subnet-123", "subnet-456"},
							"security_group_ids": []string{"sg-123"},
						},
					},
				},
				Location: model.FileLocation{
					File:   "main.tf",
					Line:   10,
					Column: 0,
				},
			},
			want: &model.LambdaFunction{
				Resource: model.Resource{
					ID:   "aws_lambda_function.test_full",
					Type: "aws_lambda_function",
					Attributes: map[string]any{
						"function_name": "test-function-full",
						"role":          "arn:aws:iam::123456789012:role/test-role",
						"handler":       "index.handler",
						"runtime":       "nodejs18.x",
						"filename":      "./lambda.zip",
						"architectures": []string{"x86_64"},
					},
					Blocks: map[string][]map[string]any{
						"environment": {
							{
								"variables": map[string]string{
									"ENV": "production",
									"KEY": "value",
								},
							},
						},
						"vpc_config": {
							{
								"subnet_ids":         []string{"subnet-123", "subnet-456"},
								"security_group_ids": []string{"sg-123"},
							},
						},
					},
					Location: model.FileLocation{
						File:   "main.tf",
						Line:   10,
						Column: 0,
					},
				},
				FunctionName:   "test-function-full",
				RoleArn:        "arn:aws:iam::123456789012:role/test-role",
				Handler:        stringPtr("index.handler"),
				Runtime:        stringPtr("nodejs18.x"),
				SourceCodePath: stringPtr("./lambda.zip"),
				Architectures:  &[]string{"x86_64"},
				Environment: &struct {
					Variables map[string]string "json:\"variables\""
				}{
					Variables: map[string]string{
						"ENV": "production",
						"KEY": "value",
					},
				},
				VPCConfig: &model.VPCConfig{
					SubnetIds:        &[]string{"subnet-123", "subnet-456"},
					SecurityGroupIds: &[]string{"sg-123"},
				},
			},
			wantErr: false,
		},
		{
			name: "error when function_name is missing",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"role": "arn:aws:iam::123456789012:role/test-role",
				},
				Blocks:   map[string][]map[string]any{},
				Location: model.FileLocation{},
			},
			want:     nil,
			wantErr:  true,
			sentinel: pkgerrors.ErrMissingAttribute,
		},
		{
			name: "error when function_name is not a string",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": 12345,
					"role":          "arn:aws:iam::123456789012:role/test-role",
				},
				Blocks:   map[string][]map[string]any{},
				Location: model.FileLocation{},
			},
			want:     nil,
			wantErr:  true,
			sentinel: pkgerrors.ErrMissingAttribute,
		},
		{
			name: "error when role is missing",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function",
				},
				Blocks:   map[string][]map[string]any{},
				Location: model.FileLocation{},
			},
			want:     nil,
			wantErr:  true,
			sentinel: pkgerrors.ErrMissingAttribute,
		},
		{
			name: "error when role is not a string",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function",
					"role":          []string{"not", "a", "string"},
				},
				Blocks:   map[string][]map[string]any{},
				Location: model.FileLocation{},
			},
			want:     nil,
			wantErr:  true,
			sentinel: pkgerrors.ErrMissingAttribute,
		},
		{
			name: "extraction with environment but no variables",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function",
					"role":          "arn:aws:iam::123456789012:role/test-role",
				},
				Blocks: map[string][]map[string]any{
					"environment": {
						{
							"other_field": "value",
						},
					},
				},
				Location: model.FileLocation{},
			},
			want: &model.LambdaFunction{
				Resource: model.Resource{
					ID:   "aws_lambda_function.test",
					Type: "aws_lambda_function",
					Attributes: map[string]any{
						"function_name": "test-function",
						"role":          "arn:aws:iam::123456789012:role/test-role",
					},
					Blocks: map[string][]map[string]any{
						"environment": {
							{
								"other_field": "value",
							},
						},
					},
					Location: model.FileLocation{},
				},
				FunctionName:   "test-function",
				RoleArn:        "arn:aws:iam::123456789012:role/test-role",
				Handler:        stringPtr(""),
				Runtime:        stringPtr(""),
				SourceCodePath: stringPtr(""),
				Architectures:  &nilArchitectures,
				Environment: &struct {
					Variables map[string]string "json:\"variables\""
				}{
					Variables: nil,
				},
				VPCConfig: nil,
			},
			wantErr: false,
		},
		{
			name: "extraction with vpc_config but missing fields",
			resource: model.Resource{
				ID:   "aws_lambda_function.test",
				Type: "aws_lambda_function",
				Attributes: map[string]any{
					"function_name": "test-function",
					"role":          "arn:aws:iam::123456789012:role/test-role",
				},
				Blocks: map[string][]map[string]any{
					"vpc_config": {
						{
							"other_field": "value",
						},
					},
				},
				Location: model.FileLocation{},
			},
			want: &model.LambdaFunction{
				Resource: model.Resource{
					ID:   "aws_lambda_function.test",
					Type: "aws_lambda_function",
					Attributes: map[string]any{
						"function_name": "test-function",
						"role":          "arn:aws:iam::123456789012:role/test-role",
					},
					Blocks: map[string][]map[string]any{
						"vpc_config": {
							{
								"other_field": "value",
							},
						},
					},
					Location: model.FileLocation{},
				},
				FunctionName:   "test-function",
				RoleArn:        "arn:aws:iam::123456789012:role/test-role",
				Handler:        stringPtr(""),
				Runtime:        stringPtr(""),
				SourceCodePath: stringPtr(""),
				Architectures:  &nilArchitectures,
				Environment: &struct {
					Variables map[string]string "json:\"variables\""
				}{
					Variables: nil,
				},
				VPCConfig: &model.VPCConfig{
					SubnetIds:        &nilSubnetIds,
					SecurityGroupIds: &nilSecurityGroupIds,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := &extractors.LambdaExtractor{}
			got, err := extractor.Extract(tt.resource)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LambdaExtractor.Extract() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.sentinel != nil && !errors.Is(err, tt.sentinel) {
					t.Errorf("LambdaExtractor.Extract() error = %v, want sentinel error %v", err, tt.sentinel)
				}
				// Verify it's an AttributeError with proper context
				var attrErr *pkgerrors.AttributeError
				if tt.sentinel == pkgerrors.ErrMissingAttribute && !errors.As(err, &attrErr) {
					t.Errorf("LambdaExtractor.Extract() error is not an AttributeError: %v", err)
				}
				return
			}

			if err != nil {
				t.Errorf("LambdaExtractor.Extract() unexpected error = %v", err)
				return
			}

			if got == nil {
				t.Errorf("LambdaExtractor.Extract() got nil, want non-nil")
				return
			}

			// Compare basic fields
			if got.FunctionName != tt.want.FunctionName {
				t.Errorf("LambdaExtractor.Extract() FunctionName = %v, want %v", got.FunctionName, tt.want.FunctionName)
			}
			if got.RoleArn != tt.want.RoleArn {
				t.Errorf("LambdaExtractor.Extract() RoleArn = %v, want %v", got.RoleArn, tt.want.RoleArn)
			}

			// Compare pointer fields
			if !equalStringPtr(got.Handler, tt.want.Handler) {
				t.Errorf("LambdaExtractor.Extract() Handler = %v, want %v", ptrToString(got.Handler), ptrToString(tt.want.Handler))
			}
			if !equalStringPtr(got.Runtime, tt.want.Runtime) {
				t.Errorf("LambdaExtractor.Extract() Runtime = %v, want %v", ptrToString(got.Runtime), ptrToString(tt.want.Runtime))
			}
			if !equalStringPtr(got.SourceCodePath, tt.want.SourceCodePath) {
				t.Errorf("LambdaExtractor.Extract() SourceCodePath = %v, want %v", ptrToString(got.SourceCodePath), ptrToString(tt.want.SourceCodePath))
			}

			// Compare architectures
			if !equalStringSlicePtr(got.Architectures, tt.want.Architectures) {
				t.Errorf("LambdaExtractor.Extract() Architectures = %v, want %v", got.Architectures, tt.want.Architectures)
			}

			// Compare environment
			if (got.Environment == nil) != (tt.want.Environment == nil) {
				t.Errorf("LambdaExtractor.Extract() Environment nil mismatch: got %v, want %v", got.Environment == nil, tt.want.Environment == nil)
			} else if got.Environment != nil {
				if !equalStringMap(got.Environment.Variables, tt.want.Environment.Variables) {
					t.Errorf("LambdaExtractor.Extract() Environment.Variables = %v, want %v", got.Environment.Variables, tt.want.Environment.Variables)
				}
			}

			// Compare VPC Config
			if (got.VPCConfig == nil) != (tt.want.VPCConfig == nil) {
				t.Errorf("LambdaExtractor.Extract() VPCConfig nil mismatch: got %v, want %v", got.VPCConfig == nil, tt.want.VPCConfig == nil)
			} else if got.VPCConfig != nil {
				if !equalStringSlicePtr(got.VPCConfig.SubnetIds, tt.want.VPCConfig.SubnetIds) {
					t.Errorf("LambdaExtractor.Extract() VPCConfig.SubnetIds = %v, want %v", got.VPCConfig.SubnetIds, tt.want.VPCConfig.SubnetIds)
				}
				if !equalStringSlicePtr(got.VPCConfig.SecurityGroupIds, tt.want.VPCConfig.SecurityGroupIds) {
					t.Errorf("LambdaExtractor.Extract() VPCConfig.SecurityGroupIds = %v, want %v", got.VPCConfig.SecurityGroupIds, tt.want.VPCConfig.SecurityGroupIds)
				}
			}
		})
	}
}

// Helper functions for tests

func stringPtr(s string) *string {
	return &s
}

func ptrToString(ptr *string) string {
	if ptr == nil {
		return "<nil>"
	}
	return *ptr
}

func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalStringSlicePtr(a, b *[]string) bool {
	// Treat both nil pointers as equal
	if a == nil && b == nil {
		return true
	}
	// If one is nil and the other is not, they're not equal
	if a == nil || b == nil {
		return false
	}
	// Dereference and compare
	aVal := *a
	bVal := *b
	// Treat nil slice and empty slice as equal
	if len(aVal) == 0 && len(bVal) == 0 {
		return true
	}
	if len(aVal) != len(bVal) {
		return false
	}
	for i, v := range aVal {
		if v != bVal[i] {
			return false
		}
	}
	return true
}

func equalStringMap(a, b map[string]string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || v != bv {
			return false
		}
	}
	return true
}
