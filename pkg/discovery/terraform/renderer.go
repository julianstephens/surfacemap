package terraform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RenderToString converts an ExprValue to a human-readable string representation.
// This is used when CLI needs text output for display or debugging purposes.
// The rendering process calls ToAny() on the ExprValue and formats the result.
func RenderToString(val ExprValue) string {
	if val == nil {
		return "<nil>"
	}

	switch v := val.(type) {
	case LiteralValue:
		if v.Val.IsNull() {
			return "null"
		}
		return fmt.Sprintf("%v", v.ToAny())

	case TraversalValue:
		return v.ToAny().(string)

	case ListValue:
		items := v.ToAny().([]any)
		strItems := make([]string, len(items))
		for i, item := range items {
			strItems[i] = fmt.Sprintf("%v", item)
		}
		return "[" + strings.Join(strItems, ", ") + "]"

	case ObjectValue:
		obj := v.ToAny().(map[string]any)
		pairs := make([]string, 0, len(obj))
		for key, value := range obj {
			pairs = append(pairs, fmt.Sprintf("%s: %v", key, value))
		}
		return "{" + strings.Join(pairs, ", ") + "}"

	case FunctionCallValue:
		return v.ToAny().(string)

	case IndexValue:
		return v.ToAny().(string)

	case TemplateValue:
		return v.ToAny().(string)

	case UnknownValue:
		return v.ToAny().(string)

	default:
		return fmt.Sprintf("<unknown type: %T>", val)
	}
}

// RenderToJSON converts an ExprValue to JSON bytes.
// This is used when CLI needs structured JSON output for machine consumption.
// The rendering uses ToAny() to get Go native types which are then marshaled to JSON.
func RenderToJSON(val ExprValue) ([]byte, error) {
	if val == nil {
		return json.Marshal(nil)
	}

	// Convert to native Go type via ToAny(), then marshal to JSON
	return json.Marshal(val.ToAny())
}

// RenderAttributesToString converts a map of attributes (from DecodeBlock) to formatted string.
// Useful for displaying resource attributes in human-readable format.
func RenderAttributesToString(attrs map[string]ExprValue) string {
	if len(attrs) == 0 {
		return "{}"
	}

	pairs := make([]string, 0, len(attrs))
	for key, val := range attrs {
		pairs = append(pairs, fmt.Sprintf("  %s = %s", key, RenderToString(val)))
	}
	return "{\n" + strings.Join(pairs, "\n") + "\n}"
}

// RenderBlocksToString converts nested blocks to formatted string representation.
// Shows block types and their instances with attributes.
func RenderBlocksToString(blocks map[string][]map[string]any) string {
	if len(blocks) == 0 {
		return "{}"
	}

	blockStrs := make([]string, 0, len(blocks))
	for blockType, instances := range blocks {
		for i, instance := range instances {
			instancePairs := make([]string, 0, len(instance))
			for key, val := range instance {
				if exprVal, ok := val.(ExprValue); ok {
					instancePairs = append(instancePairs, fmt.Sprintf("    %s = %s", key, RenderToString(exprVal)))
				} else {
					instancePairs = append(instancePairs, fmt.Sprintf("    %s = %v", key, val))
				}
			}
			blockStrs = append(blockStrs, fmt.Sprintf("  %s[%d] {\n%s\n  }", blockType, i, strings.Join(instancePairs, "\n")))
		}
	}
	return "{\n" + strings.Join(blockStrs, "\n") + "\n}"
}
