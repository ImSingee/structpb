package structpb

import (
	"encoding/base64"
	"google.golang.org/protobuf/runtime/protoimpl"
	"math"
	"unicode/utf8"
)

// NewValue constructs a Value from a general-purpose Go interface.
//
//	╔════════════════════════╤════════════════════════════════════════════╗
//	║ Go type                │ Conversion                                 ║
//	╠════════════════════════╪════════════════════════════════════════════╣
//	║ nil                    │ stored as NullValue                        ║
//	║ bool                   │ stored as BoolValue                        ║
//	║ int, int*              │ stored as IntValue                         ║
//	║ uint*                  │ stored as IntValue                         ║
//	║ uint, uint64           │ stored as IntValue / FloatValue            ║
//	║ float32, float64       │ stored as FloatValue                       ║
//	║ string                 │ stored as StringValue; must be valid UTF-8 ║
//	║ []byte                 │ stored as StringValue; base64-encoded      ║
//	║ map[string]interface{} │ stored as StructValue                      ║
//	║ []interface{}          │ stored as ListValue                        ║
//	╚════════════════════════╧════════════════════════════════════════════╝
//
// When converting an int64 or uint64 to a NumberValue, numeric precision loss
// is possible since they are stored as a float64.
func NewValue(v interface{}) (*Value, error) {
	switch v := v.(type) {
	case nil:
		return NewNullValue(), nil
	case bool:
		return NewBoolValue(v), nil
	case int:
		return NewIntValue(int64(v)), nil
	case int8:
		return NewIntValue(int64(v)), nil
	case int16:
		return NewIntValue(int64(v)), nil
	case int32:
		return NewIntValue(int64(v)), nil
	case int64:
		return NewIntValue(int64(v)), nil
	case uint8:
		return NewIntValue(int64(v)), nil
	case uint16:
		return NewIntValue(int64(v)), nil
	case uint32:
		return NewIntValue(int64(v)), nil
	case uint:
		if v > math.MaxInt64 {
			return NewFloatValue(float64(v)), nil
		} else {
			return NewIntValue(int64(v)), nil
		}
	case uint64:
		if v > math.MaxInt64 {
			return NewFloatValue(float64(v)), nil
		} else {
			return NewIntValue(int64(v)), nil
		}
	case float32:
		return NewFloatValue(float64(v)), nil
	case float64:
		return NewFloatValue(float64(v)), nil
	case string:
		if !utf8.ValidString(v) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", v)
		}
		return NewStringValue(v), nil
	case []byte:
		s := base64.StdEncoding.EncodeToString(v)
		return NewStringValue(s), nil
	case map[string]interface{}:
		v2, err := NewDict(v)
		if err != nil {
			return nil, err
		}
		return NewStructValue(v2), nil
	case []interface{}:
		v2, err := NewList(v)
		if err != nil {
			return nil, err
		}
		return NewListValue(v2), nil
	default:
		return nil, protoimpl.X.NewError("invalid type: %T", v)
	}
}

// NewNullValue constructs a new null Value.
func NewNullValue() *Value {
	return &Value{Kind: &Value_NullValue{NullValue: NullValue_NULL_VALUE}}
}

// NewBoolValue constructs a new boolean Value.
func NewBoolValue(v bool) *Value {
	return &Value{Kind: &Value_BoolValue{BoolValue: v}}
}

// NewIntValue constructs a new integer number Value.
func NewIntValue(v int64) *Value {
	return &Value{Kind: &Value_IntValue{IntValue: v}}
}

// NewFloatValue constructs a new float number Value.
func NewFloatValue(v float64) *Value {
	return &Value{Kind: &Value_FloatValue{FloatValue: v}}
}

// NewStringValue constructs a new string Value.
func NewStringValue(v string) *Value {
	return &Value{Kind: &Value_StringValue{StringValue: v}}
}

// NewStructValue constructs a new struct Value.
func NewStructValue(v *Dict) *Value {
	return &Value{Kind: &Value_DictValue{DictValue: v}}
}

// NewListValue constructs a new list Value.
func NewListValue(v *List) *Value {
	return &Value{Kind: &Value_ListValue{ListValue: v}}
}

// Unwrap returns the underlying value
// it's type may be nil, int64, float64, string, bool, *Dict, *List
//
// Call from a nil is safe
func (x *Value) Unwrap() interface{} {
	switch v := x.GetKind().(type) {
	case *Value_IntValue:
		if v != nil {
			return v.IntValue
		}
	case *Value_FloatValue:
		if v != nil {
			return v.FloatValue
		}
	case *Value_StringValue:
		if v != nil {
			return v.StringValue
		}
	case *Value_BoolValue:
		if v != nil {
			return v.BoolValue
		}
	case *Value_DictValue:
		if v != nil {
			return v.DictValue
		}
	case *Value_ListValue:
		if v != nil {
			return v.ListValue
		}
	}

	return nil
}

// AsInterface converts x to a general-purpose Go interface.
//
// Unlike Unwrap, this function will always return Go's basic types.
//
// For Null, Int, String, Bool, the return value is same as Unwrap (will return nil, int64, string, bool)
// For Float, this may return a float64 or "NaN" "Infinity" "-Infinity" string
// For Dict, this will return a map[string]interface{}, which is returned from (*Dict).AsMap()
// For List, this will return a []interface{}, which is returned from (*List).AsSlice()
//
// Call from nil is safe
func (x *Value) AsInterface() interface{} {
	switch v := x.GetKind().(type) {
	case *Value_IntValue:
		if v != nil {
			return v.IntValue
		}
	case *Value_FloatValue:
		if v != nil {
			switch {
			case math.IsNaN(v.FloatValue):
				return "NaN"
			case math.IsInf(v.FloatValue, +1):
				return "Infinity"
			case math.IsInf(v.FloatValue, -1):
				return "-Infinity"
			default:
				return v.FloatValue
			}
		}
	case *Value_StringValue:
		if v != nil {
			return v.StringValue
		}
	case *Value_BoolValue:
		if v != nil {
			return v.BoolValue
		}
	case *Value_DictValue:
		if v != nil {
			return v.DictValue.AsMap()
		}
	case *Value_ListValue:
		if v != nil {
			return v.ListValue.AsSlice()
		}
	}
	return nil
}
