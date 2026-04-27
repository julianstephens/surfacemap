package terraform

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type ExprValue interface {
	Type() string
	ToAny() any
}

type LiteralValue struct {
	Val cty.Value
	Pos hcl.Range
}

func (v LiteralValue) Type() string { return "literal" }
func (v LiteralValue) ToAny() any {
	if v.Val.IsNull() {
		return nil
	}

	switch v.Val.Type() {
	case cty.Number:
		var f float64
		_ = gocty.FromCtyValue(v.Val, &f)
		return f
	case cty.String:
		return v.Val.AsString()
	case cty.Bool:
		var b bool
		_ = gocty.FromCtyValue(v.Val, &b)
		return b
	default:
		return v.Val.AsValueMap()
	}
}

type TraversalValue struct {
	Segments []string
	Pos      hcl.Range
}

func (v TraversalValue) Type() string { return "traversal" }
func (v TraversalValue) ToAny() any {
	return strings.Join(v.Segments, ".")
}

type ListValue struct {
	Exprs []ExprValue
	Pos   hcl.Range
}

func (l ListValue) Type() string { return "list" }
func (l ListValue) ToAny() any {
	var res []any
	for _, expr := range l.Exprs {
		res = append(res, expr.ToAny())
	}
	return res
}

type ObjectValue struct {
	Fields map[string]ExprValue
	Pos    hcl.Range
}

func (o ObjectValue) Type() string { return "object" }
func (o ObjectValue) ToAny() any {
	res := make(map[string]any, len(o.Fields))

	for key, val := range o.Fields {
		res[key] = val.ToAny()
	}

	return res
}

type FunctionCallValue struct {
	Name string
	Args []ExprValue
	Pos  hcl.Range
}

func (f FunctionCallValue) Type() string { return "function" }
func (f FunctionCallValue) ToAny() any {
	argStrs := make([]string, len(f.Args))
	for i, arg := range f.Args {
		argStrs[i] = fmt.Sprintf("%v", arg.ToAny())
	}
	return fmt.Sprintf("%s(%s)", f.Name, strings.Join(argStrs, ", "))
}

type IndexValue struct {
	Collection ExprValue
	Key        ExprValue
	Pos        hcl.Range
}

func (i IndexValue) Type() string { return "index" }
func (i IndexValue) ToAny() any {
	collectionStr := fmt.Sprintf("%v", i.Collection.ToAny())
	keyStr := fmt.Sprintf("%v", i.Key.ToAny())
	return fmt.Sprintf("%s[%s]", collectionStr, keyStr)
}

type TemplateValue struct {
	Parts []ExprValue
	Pos   hcl.Range
}

func (t TemplateValue) Type() string { return "template" }
func (t TemplateValue) ToAny() any {
	var result strings.Builder
	for _, part := range t.Parts {
		fmt.Fprintf(&result, "%v", part.ToAny())
	}
	return result.String()
}

type UnknownValue struct {
	ExprType string
	Reason   string
	Pos      hcl.Range
}

func (u UnknownValue) Type() string { return "unknown" }
func (u UnknownValue) ToAny() any {
	return fmt.Sprintf("<unknown: %s at %s>", u.Reason, u.Pos.String())
}

// DecodeExpr takes an HCL expression and decodes it into a structured ExprValue that can be easily analyzed and converted to a string representation. It handles various expression types, including literals, traversals, templates, lists, objects, function calls, and index expressions. If the expression type is unsupported or if there are errors during decoding, it returns an appropriate error or an UnknownValue.
func DecodeExpr(expr hcl.Expression) (ExprValue, error) {
	switch ex := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		val, diag := ex.Value(&hcl.EvalContext{})
		if err := throwDecodeError(diag); err != nil {
			return nil, err
		}
		return LiteralValue{Val: val, Pos: ex.Range()}, nil
	case *hclsyntax.ScopeTraversalExpr:
		return TraversalValue{Segments: decodeTraversal(ex.Traversal), Pos: ex.Range()}, nil
	case *hclsyntax.TemplateExpr:
		parts := make([]ExprValue, len(ex.Parts))
		for i, part := range ex.Parts {
			decoded, err := DecodeExpr(part)
			if err != nil {
				return nil, err
			}
			parts[i] = decoded
		}
		return TemplateValue{Parts: parts, Pos: ex.Range()}, nil
	case *hclsyntax.TupleConsExpr:
		parts := make([]ExprValue, len(ex.Exprs))
		for i, part := range ex.Exprs {
			decoded, err := DecodeExpr(part)
			if err != nil {
				return nil, err
			}
			parts[i] = decoded
		}
		return ListValue{Exprs: parts, Pos: ex.Range()}, nil
	case *hclsyntax.ObjectConsExpr:
		obj := make(map[string]ExprValue)
		for _, item := range ex.Items {
			keyExpr, err := DecodeExpr(item.KeyExpr)
			if err != nil {
				return nil, err
			}

			keyStr, ok := keyExpr.ToAny().(string)
			if !ok {
				return nil, fmt.Errorf("object keys must be strings, got %T at %s", keyExpr.ToAny(), item.KeyExpr.Range().String())
			}

			valExpr, err := DecodeExpr(item.ValueExpr)
			if err != nil {
				return nil, err
			}

			obj[keyStr] = valExpr
		}
		return ObjectValue{Fields: obj, Pos: ex.Range()}, nil
	case *hclsyntax.FunctionCallExpr:
		parsedArgs := make([]ExprValue, 0, len(ex.Args))
		for _, arg := range ex.Args {
			valExpr, err := DecodeExpr(arg)
			if err != nil {
				return nil, err
			}
			parsedArgs = append(parsedArgs, valExpr)
		}
		return FunctionCallValue{Name: ex.Name, Args: parsedArgs, Pos: ex.Range()}, nil
	case *hclsyntax.IndexExpr:
		collectionExpr, err := DecodeExpr(ex.Collection)
		if err != nil {
			return nil, err
		}
		keyExpr, err := DecodeExpr(ex.Key)
		if err != nil {
			return nil, err
		}
		return IndexValue{Collection: collectionExpr, Key: keyExpr, Pos: ex.Range()}, nil
	case *hclsyntax.ObjectConsKeyExpr:
		return DecodeExpr(ex.Wrapped)
	default:
		return UnknownValue{
			ExprType: fmt.Sprintf("%T", expr),
			Reason:   "unsupported expression type",
			Pos:      expr.Range(),
		}, nil
	}
}

func decodeTraversal(traversal hcl.Traversal) (segments []string) {
	segments = make([]string, 0, len(traversal))
	for _, step := range traversal {
		switch s := step.(type) {
		case hcl.TraverseRoot:
			segments = append(segments, s.Name)
		case hcl.TraverseAttr:
			segments = append(segments, s.Name)
		case hcl.TraverseIndex:
			if s.Key.Type() == cty.String {
				segments = append(segments, fmt.Sprintf("[%s]", s.Key.AsString()))
			} else if s.Key.Type() == cty.Number {
				var idx int
				_ = gocty.FromCtyValue(s.Key, &idx)
				segments = append(segments, fmt.Sprintf("[%d]", idx))
			} else {
				segments = append(segments, fmt.Sprintf("[%v]", s.Key))
			}
		}
	}
	return
}

func throwDecodeError(diag hcl.Diagnostics) (err error) {
	if diag.HasErrors() {
		err = &TerraformParserError{
			Err:   ErrDecodeExpression,
			Cause: errors.New(diag.Error()),
		}
	}
	return
}

// DecodeBlock takes an HCL block body and recursively decodes its attributes and nested blocks into structured maps of ExprValue. It returns a map of attribute names to their decoded values, a map of block types to lists of their decoded contents, and any error encountered during decoding. This function allows for easy analysis and manipulation of the HCL configuration by converting it into a more accessible format.
func DecodeBlock(body *hclsyntax.Body) (map[string]ExprValue, map[string][]map[string]any, error) {
	attrs := make(map[string]ExprValue)
	blocks := make(map[string][]map[string]any)

	for name, attr := range body.Attributes {
		val, err := DecodeExpr(attr.Expr)
		if err != nil {
			return nil, nil, err
		}
		attrs[name] = val
	}

	for _, block := range body.Blocks {
		blockAttrs, nestedBlocks, err := DecodeBlock(block.Body)
		if err != nil {
			return nil, nil, err
		}

		blockData := make(map[string]any)
		for k, v := range blockAttrs {
			blockData[k] = v
		}

		for k, v := range nestedBlocks {
			blockData[k] = v
		}

		blocks[block.Type] = append(blocks[block.Type], blockData)
	}

	return attrs, blocks, nil
}
