package terraform

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func TestLiteralValue_Type(t *testing.T) {
	lv := LiteralValue{Val: cty.StringVal("test")}
	if got := lv.Type(); got != "literal" {
		t.Errorf("LiteralValue.Type() = %q, want %q", got, "literal")
	}
}

func TestLiteralValue_ToAny_String(t *testing.T) {
	lv := LiteralValue{Val: cty.StringVal("hello")}
	got := lv.ToAny()
	want := "hello"
	if got != want {
		t.Errorf("LiteralValue.ToAny() = %v, want %v", got, want)
	}
}

func TestLiteralValue_ToAny_Number(t *testing.T) {
	lv := LiteralValue{Val: cty.NumberIntVal(42)}
	got := lv.ToAny()
	want := 42.0
	if got != want {
		t.Errorf("LiteralValue.ToAny() = %v, want %v", got, want)
	}
}

func TestLiteralValue_ToAny_Float(t *testing.T) {
	lv := LiteralValue{Val: cty.NumberFloatVal(3.14)}
	got := lv.ToAny().(float64)
	want := 3.14
	if got != want {
		t.Errorf("LiteralValue.ToAny() = %v, want %v", got, want)
	}
}

func TestLiteralValue_ToAny_Bool(t *testing.T) {
	tests := []struct {
		name string
		val  cty.Value
		want bool
	}{
		{"true", cty.True, true},
		{"false", cty.False, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := LiteralValue{Val: tt.val}
			got := lv.ToAny()
			if got != tt.want {
				t.Errorf("LiteralValue.ToAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLiteralValue_ToAny_Null(t *testing.T) {
	lv := LiteralValue{Val: cty.NullVal(cty.String)}
	got := lv.ToAny()
	if got != nil {
		t.Errorf("LiteralValue.ToAny() = %v, want nil", got)
	}
}

func TestTraversalValue_Type(t *testing.T) {
	tv := TraversalValue{Segments: []string{"var", "foo"}}
	if got := tv.Type(); got != "traversal" {
		t.Errorf("TraversalValue.Type() = %q, want %q", got, "traversal")
	}
}

func TestTraversalValue_ToAny(t *testing.T) {
	tests := []struct {
		name     string
		segments []string
		want     string
	}{
		{
			name:     "simple variable",
			segments: []string{"var", "project"},
			want:     "var.project",
		},
		{
			name:     "nested attribute",
			segments: []string{"var", "lambda", "function_name"},
			want:     "var.lambda.function_name",
		},
		{
			name:     "local value",
			segments: []string{"local", "tags"},
			want:     "local.tags",
		},
		{
			name:     "data source",
			segments: []string{"data", "aws_iam_policy_document", "lambda", "json"},
			want:     "data.aws_iam_policy_document.lambda.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := TraversalValue{Segments: tt.segments}
			got := tv.ToAny()
			if got != tt.want {
				t.Errorf("TraversalValue.ToAny() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListValue_Type(t *testing.T) {
	lv := ListValue{Exprs: []ExprValue{}}
	if got := lv.Type(); got != "list" {
		t.Errorf("ListValue.Type() = %q, want %q", got, "list")
	}
}

func TestListValue_ToAny(t *testing.T) {
	tests := []struct {
		name  string
		exprs []ExprValue
		want  []any
	}{
		{
			name:  "empty list",
			exprs: []ExprValue{},
			want:  []any{},
		},
		{
			name: "list of strings",
			exprs: []ExprValue{
				LiteralValue{Val: cty.StringVal("a")},
				LiteralValue{Val: cty.StringVal("b")},
				LiteralValue{Val: cty.StringVal("c")},
			},
			want: []any{"a", "b", "c"},
		},
		{
			name: "list of numbers",
			exprs: []ExprValue{
				LiteralValue{Val: cty.NumberIntVal(1)},
				LiteralValue{Val: cty.NumberIntVal(2)},
				LiteralValue{Val: cty.NumberIntVal(3)},
			},
			want: []any{1.0, 2.0, 3.0},
		},
		{
			name: "mixed types",
			exprs: []ExprValue{
				LiteralValue{Val: cty.StringVal("foo")},
				LiteralValue{Val: cty.NumberIntVal(42)},
				TraversalValue{Segments: []string{"var", "bar"}},
			},
			want: []any{"foo", 42.0, "var.bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := ListValue{Exprs: tt.exprs}
			got := lv.ToAny().([]any)

			if len(got) != len(tt.want) {
				t.Fatalf("ListValue.ToAny() length = %d, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ListValue.ToAny()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestObjectValue_Type(t *testing.T) {
	ov := ObjectValue{Fields: map[string]ExprValue{}}
	if got := ov.Type(); got != "object" {
		t.Errorf("ObjectValue.Type() = %q, want %q", got, "object")
	}
}

func TestObjectValue_ToAny(t *testing.T) {
	tests := []struct {
		name    string
		exprMap map[string]ExprValue
		want    map[string]any
	}{
		{
			name:    "empty object",
			exprMap: map[string]ExprValue{},
			want:    map[string]any{},
		},
		{
			name: "simple object",
			exprMap: map[string]ExprValue{
				"name":    LiteralValue{Val: cty.StringVal("test")},
				"enabled": LiteralValue{Val: cty.True},
				"count":   LiteralValue{Val: cty.NumberIntVal(5)},
			},
			want: map[string]any{
				"name":    "test",
				"enabled": true,
				"count":   5.0,
			},
		},
		{
			name: "nested with traversal",
			exprMap: map[string]ExprValue{
				"table_name": TraversalValue{Segments: []string{"var", "table"}},
				"region":     LiteralValue{Val: cty.StringVal("us-east-1")},
			},
			want: map[string]any{
				"table_name": "var.table",
				"region":     "us-east-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ov := ObjectValue{Fields: tt.exprMap}
			got := ov.ToAny().(map[string]any)

			if len(got) != len(tt.want) {
				t.Fatalf("ObjectValue.ToAny() length = %d, want %d", len(got), len(tt.want))
			}

			for key, wantVal := range tt.want {
				gotVal, ok := got[key]
				if !ok {
					t.Errorf("ObjectValue.ToAny() missing key %q", key)
					continue
				}
				if gotVal != wantVal {
					t.Errorf("ObjectValue.ToAny()[%q] = %v, want %v", key, gotVal, wantVal)
				}
			}
		})
	}
}

func TestFunctionCallValue_Type(t *testing.T) {
	fcv := FunctionCallValue{Name: "test"}
	if got := fcv.Type(); got != "function" {
		t.Errorf("FunctionCallValue.Type() = %q, want %q", got, "function")
	}
}

func TestFunctionCallValue_ToAny(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		args     []ExprValue
		want     string
	}{
		{
			name:     "no args",
			funcName: "now",
			args:     []ExprValue{},
			want:     "now()",
		},
		{
			name:     "single string arg",
			funcName: "upper",
			args: []ExprValue{
				LiteralValue{Val: cty.StringVal("hello")},
			},
			want: "upper(hello)",
		},
		{
			name:     "multiple args",
			funcName: "format",
			args: []ExprValue{
				LiteralValue{Val: cty.StringVal("%s-bucket")},
				TraversalValue{Segments: []string{"var", "project"}},
			},
			want: "format(%s-bucket, var.project)",
		},
		{
			name:     "jsonencode with object",
			funcName: "jsonencode",
			args: []ExprValue{
				ObjectValue{
					Fields: map[string]ExprValue{
						"Version": LiteralValue{Val: cty.StringVal("2012-10-17")},
					},
				},
			},
			want: "jsonencode(map[Version:2012-10-17])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcv := FunctionCallValue{Name: tt.funcName, Args: tt.args}
			got := fcv.ToAny()
			if got != tt.want {
				t.Errorf("FunctionCallValue.ToAny() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIndexValue_Type(t *testing.T) {
	iv := IndexValue{
		Collection: TraversalValue{Segments: []string{"var", "list"}},
		Key:        LiteralValue{Val: cty.NumberIntVal(0)},
	}
	if got := iv.Type(); got != "index" {
		t.Errorf("IndexValue.Type() = %q, want %q", got, "index")
	}
}

func TestIndexValue_ToAny(t *testing.T) {
	tests := []struct {
		name       string
		collection ExprValue
		key        ExprValue
		want       string
	}{
		{
			name:       "list index with number",
			collection: TraversalValue{Segments: []string{"var", "subnets"}},
			key:        LiteralValue{Val: cty.NumberIntVal(0)},
			want:       "var.subnets[0]",
		},
		{
			name:       "map index with string",
			collection: TraversalValue{Segments: []string{"local", "config"}},
			key:        LiteralValue{Val: cty.StringVal("region")},
			want:       "local.config[region]",
		},
		{
			name:       "dynamic index",
			collection: TraversalValue{Segments: []string{"aws_subnet", "private"}},
			key:        TraversalValue{Segments: []string{"count", "index"}},
			want:       "aws_subnet.private[count.index]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := IndexValue{Collection: tt.collection, Key: tt.key}
			got := iv.ToAny()
			if got != tt.want {
				t.Errorf("IndexValue.ToAny() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTemplateValue_Type(t *testing.T) {
	tv := TemplateValue{Parts: []ExprValue{}}
	if got := tv.Type(); got != "template" {
		t.Errorf("TemplateValue.Type() = %q, want %q", got, "template")
	}
}

func TestTemplateValue_ToAny(t *testing.T) {
	tests := []struct {
		name  string
		parts []ExprValue
		want  string
	}{
		{
			name: "simple interpolation",
			parts: []ExprValue{
				TraversalValue{Segments: []string{"var", "project"}},
				LiteralValue{Val: cty.StringVal("-bucket")},
			},
			want: "var.project-bucket",
		},
		{
			name: "multiple interpolations",
			parts: []ExprValue{
				TraversalValue{Segments: []string{"var", "project"}},
				LiteralValue{Val: cty.StringVal("-")},
				TraversalValue{Segments: []string{"var", "environment"}},
				LiteralValue{Val: cty.StringVal("-bucket")},
			},
			want: "var.project-var.environment-bucket",
		},
		{
			name: "only literal",
			parts: []ExprValue{
				LiteralValue{Val: cty.StringVal("static-string")},
			},
			want: "static-string",
		},
		{
			name: "with numbers",
			parts: []ExprValue{
				LiteralValue{Val: cty.StringVal("port-")},
				LiteralValue{Val: cty.NumberIntVal(8080)},
			},
			want: "port-8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tv := TemplateValue{Parts: tt.parts}
			got := tv.ToAny()
			if got != tt.want {
				t.Errorf("TemplateValue.ToAny() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnknownValue_Type(t *testing.T) {
	uv := UnknownValue{ExprType: "test"}
	if got := uv.Type(); got != "unknown" {
		t.Errorf("UnknownValue.Type() = %q, want %q", got, "unknown")
	}
}

func TestUnknownValue_ToAny(t *testing.T) {
	tests := []struct {
		name         string
		exprType     string
		reason       string
		pos          hcl.Range
		wantContains []string
	}{
		{
			name:     "unsupported for expression",
			exprType: "*hclsyntax.ForExpr",
			reason:   "for expressions not supported",
			pos: hcl.Range{
				Filename: "test.tf",
				Start:    hcl.Pos{Line: 10, Column: 5},
				End:      hcl.Pos{Line: 10, Column: 45},
			},
			wantContains: []string{"for expressions not supported", "test.tf:10"},
		},
		{
			name:     "unsupported splat",
			exprType: "*hclsyntax.SplatExpr",
			reason:   "splat operator not supported",
			pos: hcl.Range{
				Filename: "lambda.tf",
				Start:    hcl.Pos{Line: 5, Column: 10},
			},
			wantContains: []string{"splat operator not supported", "lambda.tf:5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uv := UnknownValue{
				ExprType: tt.exprType,
				Reason:   tt.reason,
				Pos:      tt.pos,
			}
			got := uv.ToAny().(string)

			for _, want := range tt.wantContains {
				if !contains(got, want) {
					t.Errorf("UnknownValue.ToAny() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}

// Test nested structures
func TestNestedStructures(t *testing.T) {
	t.Run("list of objects", func(t *testing.T) {
		lv := ListValue{
			Exprs: []ExprValue{
				ObjectValue{
					Fields: map[string]ExprValue{
						"name": LiteralValue{Val: cty.StringVal("item1")},
					},
				},
				ObjectValue{
					Fields: map[string]ExprValue{
						"name": LiteralValue{Val: cty.StringVal("item2")},
					},
				},
			},
		}

		got := lv.ToAny().([]any)
		if len(got) != 2 {
			t.Fatalf("expected 2 items, got %d", len(got))
		}

		obj1 := got[0].(map[string]any)
		if obj1["name"] != "item1" {
			t.Errorf("first object name = %v, want %v", obj1["name"], "item1")
		}
	})

	t.Run("object with nested list", func(t *testing.T) {
		ov := ObjectValue{
			Fields: map[string]ExprValue{
				"actions": ListValue{
					Exprs: []ExprValue{
						LiteralValue{Val: cty.StringVal("dynamodb:GetItem")},
						LiteralValue{Val: cty.StringVal("dynamodb:PutItem")},
					},
				},
			},
		}

		got := ov.ToAny().(map[string]any)
		actions := got["actions"].([]any)

		if len(actions) != 2 {
			t.Fatalf("expected 2 actions, got %d", len(actions))
		}
		if actions[0] != "dynamodb:GetItem" {
			t.Errorf("first action = %v, want %v", actions[0], "dynamodb:GetItem")
		}
	})
}

// DecodeExpr tests
func TestDecodeExpr_LiteralValueExpr(t *testing.T) {
	tests := []struct {
		name    string
		expr    *hclsyntax.LiteralValueExpr
		want    any
		wantErr bool
	}{
		{
			name: "string literal",
			expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("hello")},
			want: "hello",
		},
		{
			name: "number literal",
			expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(42)},
			want: 42.0,
		},
		{
			name: "bool literal",
			expr: &hclsyntax.LiteralValueExpr{Val: cty.True},
			want: true,
		},
		{
			name: "null literal",
			expr: &hclsyntax.LiteralValueExpr{Val: cty.NullVal(cty.String)},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeExpr(tt.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeExpr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got.Type() != "literal" {
					t.Errorf("DecodeExpr() type = %v, want literal", got.Type())
				}
				if got.ToAny() != tt.want {
					t.Errorf("DecodeExpr().ToAny() = %v, want %v", got.ToAny(), tt.want)
				}
			}
		})
	}
}

func TestDecodeExpr_ScopeTraversalExpr(t *testing.T) {
	tests := []struct {
		name      string
		traversal hcl.Traversal
		want      string
	}{
		{
			name: "simple variable",
			traversal: hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: "project"},
			},
			want: "var.project",
		},
		{
			name: "nested attribute",
			traversal: hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: "lambda"},
				hcl.TraverseAttr{Name: "function_name"},
			},
			want: "var.lambda.function_name",
		},
		{
			name: "with index",
			traversal: hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: "subnets"},
				hcl.TraverseIndex{Key: cty.NumberIntVal(0)},
			},
			want: "var.subnets.[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &hclsyntax.ScopeTraversalExpr{Traversal: tt.traversal}
			got, err := DecodeExpr(expr)
			if err != nil {
				t.Errorf("DecodeExpr() error = %v", err)
				return
			}
			if got.Type() != "traversal" {
				t.Errorf("DecodeExpr() type = %v, want traversal", got.Type())
			}
			if got.ToAny() != tt.want {
				t.Errorf("DecodeExpr().ToAny() = %v, want %v", got.ToAny(), tt.want)
			}
		})
	}
}

func TestDecodeExpr_TemplateExpr(t *testing.T) {
	t.Run("template with interpolation", func(t *testing.T) {
		expr := &hclsyntax.TemplateExpr{
			Parts: []hclsyntax.Expression{
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("prefix-")},
				&hclsyntax.ScopeTraversalExpr{
					Traversal: hcl.Traversal{
						hcl.TraverseRoot{Name: "var"},
						hcl.TraverseAttr{Name: "name"},
					},
				},
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("-suffix")},
			},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}
		if got.Type() != "template" {
			t.Errorf("DecodeExpr() type = %v, want template", got.Type())
		}

		want := "prefix-var.name-suffix"
		if got.ToAny() != want {
			t.Errorf("DecodeExpr().ToAny() = %v, want %v", got.ToAny(), want)
		}
	})
}

func TestDecodeExpr_TupleConsExpr(t *testing.T) {
	t.Run("list of literals", func(t *testing.T) {
		expr := &hclsyntax.TupleConsExpr{
			Exprs: []hclsyntax.Expression{
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("a")},
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("b")},
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("c")},
			},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}
		if got.Type() != "list" {
			t.Errorf("DecodeExpr() type = %v, want list", got.Type())
		}

		result := got.ToAny().([]any)
		want := []any{"a", "b", "c"}

		if len(result) != len(want) {
			t.Fatalf("DecodeExpr().ToAny() length = %d, want %d", len(result), len(want))
		}
		for i := range result {
			if result[i] != want[i] {
				t.Errorf("DecodeExpr().ToAny()[%d] = %v, want %v", i, result[i], want[i])
			}
		}
	})
}

func TestDecodeExpr_ObjectConsExpr(t *testing.T) {
	t.Run("simple object", func(t *testing.T) {
		expr := &hclsyntax.ObjectConsExpr{
			Items: []hclsyntax.ObjectConsItem{
				{
					KeyExpr:   &hclsyntax.LiteralValueExpr{Val: cty.StringVal("name")},
					ValueExpr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("test")},
				},
				{
					KeyExpr:   &hclsyntax.LiteralValueExpr{Val: cty.StringVal("count")},
					ValueExpr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(5)},
				},
			},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}
		if got.Type() != "object" {
			t.Errorf("DecodeExpr() type = %v, want object", got.Type())
		}

		result := got.ToAny().(map[string]any)
		if result["name"] != "test" {
			t.Errorf("DecodeExpr().ToAny()[name] = %v, want test", result["name"])
		}
		if result["count"] != 5.0 {
			t.Errorf("DecodeExpr().ToAny()[count] = %v, want 5.0", result["count"])
		}
	})
}

func TestDecodeExpr_FunctionCallExpr(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		args     []hclsyntax.Expression
		want     string
	}{
		{
			name:     "no args",
			funcName: "now",
			args:     []hclsyntax.Expression{},
			want:     "now()",
		},
		{
			name:     "single arg",
			funcName: "upper",
			args: []hclsyntax.Expression{
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("hello")},
			},
			want: "upper(hello)",
		},
		{
			name:     "multiple args",
			funcName: "format",
			args: []hclsyntax.Expression{
				&hclsyntax.LiteralValueExpr{Val: cty.StringVal("%s-bucket")},
				&hclsyntax.ScopeTraversalExpr{
					Traversal: hcl.Traversal{
						hcl.TraverseRoot{Name: "var"},
						hcl.TraverseAttr{Name: "project"},
					},
				},
			},
			want: "format(%s-bucket, var.project)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &hclsyntax.FunctionCallExpr{
				Name: tt.funcName,
				Args: tt.args,
			}

			got, err := DecodeExpr(expr)
			if err != nil {
				t.Errorf("DecodeExpr() error = %v", err)
				return
			}
			if got.Type() != "function" {
				t.Errorf("DecodeExpr() type = %v, want function", got.Type())
			}
			if got.ToAny() != tt.want {
				t.Errorf("DecodeExpr().ToAny() = %v, want %v", got.ToAny(), tt.want)
			}
		})
	}
}

func TestDecodeExpr_IndexExpr(t *testing.T) {
	t.Run("list index", func(t *testing.T) {
		expr := &hclsyntax.IndexExpr{
			Collection: &hclsyntax.ScopeTraversalExpr{
				Traversal: hcl.Traversal{
					hcl.TraverseRoot{Name: "var"},
					hcl.TraverseAttr{Name: "subnets"},
				},
			},
			Key: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(0)},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}
		if got.Type() != "index" {
			t.Errorf("DecodeExpr() type = %v, want index", got.Type())
		}

		want := "var.subnets[0]"
		if got.ToAny() != want {
			t.Errorf("DecodeExpr().ToAny() = %v, want %v", got.ToAny(), want)
		}
	})
}

func TestDecodeExpr_ObjectConsKeyExpr(t *testing.T) {
	t.Run("unwraps key expression", func(t *testing.T) {
		expr := &hclsyntax.ObjectConsKeyExpr{
			Wrapped: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("key")},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}
		if got.Type() != "literal" {
			t.Errorf("DecodeExpr() type = %v, want literal", got.Type())
		}
		if got.ToAny() != "key" {
			t.Errorf("DecodeExpr().ToAny() = %v, want key", got.ToAny())
		}
	})
}

func TestDecodeExpr_UnknownType(t *testing.T) {
	t.Run("unsupported expression type", func(t *testing.T) {
		// Use a conditional expression which is not supported
		expr := &hclsyntax.ConditionalExpr{
			Condition:   &hclsyntax.LiteralValueExpr{Val: cty.True},
			TrueResult:  &hclsyntax.LiteralValueExpr{Val: cty.StringVal("yes")},
			FalseResult: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("no")},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() should not return error for unknown type, got %v", err)
			return
		}
		if got.Type() != "unknown" {
			t.Errorf("DecodeExpr() type = %v, want unknown", got.Type())
		}

		result := got.ToAny().(string)
		if !contains(result, "unsupported expression type") {
			t.Errorf("DecodeExpr().ToAny() = %v, want to contain 'unsupported expression type'", result)
		}
	})
}

func TestDecodeExpr_ErrorPropagation(t *testing.T) {
	t.Run("error in nested expression", func(t *testing.T) {
		// This would require creating an expression that causes an error
		// For now, we test that error propagation works with a nil pointer scenario
		// In practice, errors happen during recursive DecodeExpr calls

		// Create a template with a part that would error (simulated)
		// Most errors in DecodeExpr come from recursive calls or HCL evaluation
		// This is more of an integration test scenario
		t.Skip("Error propagation tested via integration tests")
	})
}

func TestDecodeExpr_ComplexNested(t *testing.T) {
	t.Run("list of objects with templates", func(t *testing.T) {
		expr := &hclsyntax.TupleConsExpr{
			Exprs: []hclsyntax.Expression{
				&hclsyntax.ObjectConsExpr{
					Items: []hclsyntax.ObjectConsItem{
						{
							KeyExpr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("name")},
							ValueExpr: &hclsyntax.TemplateExpr{
								Parts: []hclsyntax.Expression{
									&hclsyntax.LiteralValueExpr{Val: cty.StringVal("item-")},
									&hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(1)},
								},
							},
						},
					},
				},
			},
		}

		got, err := DecodeExpr(expr)
		if err != nil {
			t.Errorf("DecodeExpr() error = %v", err)
			return
		}

		result := got.ToAny().([]any)
		if len(result) != 1 {
			t.Fatalf("expected 1 item, got %d", len(result))
		}

		obj := result[0].(map[string]any)
		if obj["name"] != "item-1" {
			t.Errorf("nested object name = %v, want item-1", obj["name"])
		}
	})
}

// DecodeBlock tests
func TestDecodeBlock_AttributesOnly(t *testing.T) {
	t.Run("block with only attributes", func(t *testing.T) {
		// Simulate: resource { name = "test", count = 5 }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"name": &hclsyntax.Attribute{
					Name: "name",
					Expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("test")},
				},
				"count": &hclsyntax.Attribute{
					Name: "count",
					Expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(5)},
				},
			},
			Blocks: hclsyntax.Blocks{},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		// Check attributes
		if len(attrs) != 2 {
			t.Fatalf("expected 2 attributes, got %d", len(attrs))
		}
		if attrs["name"].ToAny() != "test" {
			t.Errorf("attrs[name] = %v, want test", attrs["name"].ToAny())
		}
		if attrs["count"].ToAny() != 5.0 {
			t.Errorf("attrs[count] = %v, want 5.0", attrs["count"].ToAny())
		}

		// Check blocks (should be empty)
		if len(blocks) != 0 {
			t.Errorf("expected 0 blocks, got %d", len(blocks))
		}
	})
}

func TestDecodeBlock_NestedBlocks(t *testing.T) {
	t.Run("block with nested blocks", func(t *testing.T) {
		// Simulate:
		// resource {
		//   name = "lambda"
		//   environment {
		//     variables = { KEY = "value" }
		//   }
		// }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"name": &hclsyntax.Attribute{
					Name: "name",
					Expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("lambda")},
				},
			},
			Blocks: hclsyntax.Blocks{
				&hclsyntax.Block{
					Type: "environment",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"variables": &hclsyntax.Attribute{
								Name: "variables",
								Expr: &hclsyntax.ObjectConsExpr{
									Items: []hclsyntax.ObjectConsItem{
										{
											KeyExpr:   &hclsyntax.LiteralValueExpr{Val: cty.StringVal("KEY")},
											ValueExpr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("value")},
										},
									},
								},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
			},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		// Check top-level attributes
		if len(attrs) != 1 {
			t.Fatalf("expected 1 attribute, got %d", len(attrs))
		}
		if attrs["name"].ToAny() != "lambda" {
			t.Errorf("attrs[name] = %v, want lambda", attrs["name"].ToAny())
		}

		// Check blocks
		if len(blocks) != 1 {
			t.Fatalf("expected 1 block type, got %d", len(blocks))
		}
		envBlocks, ok := blocks["environment"]
		if !ok {
			t.Fatal("expected environment block")
		}
		if len(envBlocks) != 1 {
			t.Fatalf("expected 1 environment block, got %d", len(envBlocks))
		}

		// Check nested block content
		envBlock := envBlocks[0]
		variables, ok := envBlock["variables"]
		if !ok {
			t.Fatal("expected variables attribute in environment block")
		}

		// Variables should be an ExprValue (ObjectValue)
		exprVal, ok := variables.(ExprValue)
		if !ok {
			t.Fatalf("variables should be ExprValue, got %T", variables)
		}

		varsMap := exprVal.ToAny().(map[string]any)
		if varsMap["KEY"] != "value" {
			t.Errorf("variables[KEY] = %v, want value", varsMap["KEY"])
		}
	})
}

func TestDecodeBlock_MultipleNestingLevels(t *testing.T) {
	t.Run("deeply nested blocks", func(t *testing.T) {
		// Simulate:
		// resource {
		//   outer {
		//     middle {
		//       inner = "deep"
		//     }
		//   }
		// }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{},
			Blocks: hclsyntax.Blocks{
				&hclsyntax.Block{
					Type: "outer",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{},
						Blocks: hclsyntax.Blocks{
							&hclsyntax.Block{
								Type: "middle",
								Body: &hclsyntax.Body{
									Attributes: hclsyntax.Attributes{
										"inner": &hclsyntax.Attribute{
											Name: "inner",
											Expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("deep")},
										},
									},
									Blocks: hclsyntax.Blocks{},
								},
							},
						},
					},
				},
			},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		// No top-level attributes
		if len(attrs) != 0 {
			t.Errorf("expected 0 top-level attributes, got %d", len(attrs))
		}

		// Check outer block
		outerBlocks := blocks["outer"]
		if len(outerBlocks) != 1 {
			t.Fatalf("expected 1 outer block, got %d", len(outerBlocks))
		}

		// Check middle block (nested in outer)
		middleBlocks, ok := outerBlocks[0]["middle"].([]map[string]any)
		if !ok {
			t.Fatalf("expected middle blocks to be []map[string]any, got %T", outerBlocks[0]["middle"])
		}
		if len(middleBlocks) != 1 {
			t.Fatalf("expected 1 middle block, got %d", len(middleBlocks))
		}

		// Check inner attribute (in middle block)
		innerVal, ok := middleBlocks[0]["inner"]
		if !ok {
			t.Fatal("expected inner attribute in middle block")
		}

		exprVal, ok := innerVal.(ExprValue)
		if !ok {
			t.Fatalf("inner should be ExprValue, got %T", innerVal)
		}

		if exprVal.ToAny() != "deep" {
			t.Errorf("inner = %v, want deep", exprVal.ToAny())
		}
	})
}

func TestDecodeBlock_MultipleSameTypeBlocks(t *testing.T) {
	t.Run("multiple blocks of same type", func(t *testing.T) {
		// Simulate:
		// resource {
		//   ingress { port = 80 }
		//   ingress { port = 443 }
		//   ingress { port = 8080 }
		// }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{},
			Blocks: hclsyntax.Blocks{
				&hclsyntax.Block{
					Type: "ingress",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"port": &hclsyntax.Attribute{
								Name: "port",
								Expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(80)},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
				&hclsyntax.Block{
					Type: "ingress",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"port": &hclsyntax.Attribute{
								Name: "port",
								Expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(443)},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
				&hclsyntax.Block{
					Type: "ingress",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"port": &hclsyntax.Attribute{
								Name: "port",
								Expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(8080)},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
			},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		if len(attrs) != 0 {
			t.Errorf("expected 0 attributes, got %d", len(attrs))
		}

		// Check that all 3 ingress blocks are present
		ingressBlocks := blocks["ingress"]
		if len(ingressBlocks) != 3 {
			t.Fatalf("expected 3 ingress blocks, got %d", len(ingressBlocks))
		}

		// Check each port
		expectedPorts := []float64{80.0, 443.0, 8080.0}
		for i, expected := range expectedPorts {
			portVal := ingressBlocks[i]["port"].(ExprValue)
			if portVal.ToAny() != expected {
				t.Errorf("ingress[%d].port = %v, want %v", i, portVal.ToAny(), expected)
			}
		}
	})
}

func TestDecodeBlock_EmptyBlock(t *testing.T) {
	t.Run("empty block", func(t *testing.T) {
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{},
			Blocks:     hclsyntax.Blocks{},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		if len(attrs) != 0 {
			t.Errorf("expected 0 attributes, got %d", len(attrs))
		}
		if len(blocks) != 0 {
			t.Errorf("expected 0 blocks, got %d", len(blocks))
		}
	})
}

func TestDecodeBlock_MixedContent(t *testing.T) {
	t.Run("attributes and blocks together", func(t *testing.T) {
		// Simulate:
		// resource {
		//   function_name = "my-lambda"
		//   timeout = 30
		//   environment { variables = { KEY = "val" } }
		//   logging_config { log_level = "INFO" }
		// }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"function_name": &hclsyntax.Attribute{
					Name: "function_name",
					Expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("my-lambda")},
				},
				"timeout": &hclsyntax.Attribute{
					Name: "timeout",
					Expr: &hclsyntax.LiteralValueExpr{Val: cty.NumberIntVal(30)},
				},
			},
			Blocks: hclsyntax.Blocks{
				&hclsyntax.Block{
					Type: "environment",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"variables": &hclsyntax.Attribute{
								Name: "variables",
								Expr: &hclsyntax.ObjectConsExpr{
									Items: []hclsyntax.ObjectConsItem{
										{
											KeyExpr:   &hclsyntax.LiteralValueExpr{Val: cty.StringVal("KEY")},
											ValueExpr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("val")},
										},
									},
								},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
				&hclsyntax.Block{
					Type: "logging_config",
					Body: &hclsyntax.Body{
						Attributes: hclsyntax.Attributes{
							"log_level": &hclsyntax.Attribute{
								Name: "log_level",
								Expr: &hclsyntax.LiteralValueExpr{Val: cty.StringVal("INFO")},
							},
						},
						Blocks: hclsyntax.Blocks{},
					},
				},
			},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		// Check attributes
		if len(attrs) != 2 {
			t.Fatalf("expected 2 attributes, got %d", len(attrs))
		}
		if attrs["function_name"].ToAny() != "my-lambda" {
			t.Errorf("function_name = %v, want my-lambda", attrs["function_name"].ToAny())
		}
		if attrs["timeout"].ToAny() != 30.0 {
			t.Errorf("timeout = %v, want 30.0", attrs["timeout"].ToAny())
		}

		// Check blocks
		if len(blocks) != 2 {
			t.Fatalf("expected 2 block types, got %d", len(blocks))
		}

		// Verify environment block
		if len(blocks["environment"]) != 1 {
			t.Errorf("expected 1 environment block, got %d", len(blocks["environment"]))
		}

		// Verify logging_config block
		if len(blocks["logging_config"]) != 1 {
			t.Errorf("expected 1 logging_config block, got %d", len(blocks["logging_config"]))
		}

		logLevel := blocks["logging_config"][0]["log_level"].(ExprValue)
		if logLevel.ToAny() != "INFO" {
			t.Errorf("log_level = %v, want INFO", logLevel.ToAny())
		}
	})
}

func TestDecodeBlock_WithTraversals(t *testing.T) {
	t.Run("attributes with variable references", func(t *testing.T) {
		// Simulate:
		// resource {
		//   name = var.project_name
		//   timeout = var.lambda.timeout
		// }
		body := &hclsyntax.Body{
			Attributes: hclsyntax.Attributes{
				"name": &hclsyntax.Attribute{
					Name: "name",
					Expr: &hclsyntax.ScopeTraversalExpr{
						Traversal: hcl.Traversal{
							hcl.TraverseRoot{Name: "var"},
							hcl.TraverseAttr{Name: "project_name"},
						},
					},
				},
				"timeout": &hclsyntax.Attribute{
					Name: "timeout",
					Expr: &hclsyntax.ScopeTraversalExpr{
						Traversal: hcl.Traversal{
							hcl.TraverseRoot{Name: "var"},
							hcl.TraverseAttr{Name: "lambda"},
							hcl.TraverseAttr{Name: "timeout"},
						},
					},
				},
			},
			Blocks: hclsyntax.Blocks{},
		}

		attrs, blocks, err := DecodeBlock(body)
		if err != nil {
			t.Errorf("DecodeBlock() error = %v", err)
			return
		}

		if len(attrs) != 2 {
			t.Fatalf("expected 2 attributes, got %d", len(attrs))
		}

		// Check traversal values
		if attrs["name"].Type() != "traversal" {
			t.Errorf("name type = %v, want traversal", attrs["name"].Type())
		}
		if attrs["name"].ToAny() != "var.project_name" {
			t.Errorf("name = %v, want var.project_name", attrs["name"].ToAny())
		}

		if attrs["timeout"].ToAny() != "var.lambda.timeout" {
			t.Errorf("timeout = %v, want var.lambda.timeout", attrs["timeout"].ToAny())
		}

		if len(blocks) != 0 {
			t.Errorf("expected 0 blocks, got %d", len(blocks))
		}
	})
}
