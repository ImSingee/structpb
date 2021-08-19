package structpb

import (
	"go.starlark.net/starlark"
)

func (x *NullValue) ToStarlark() starlark.NoneType {
	return starlark.None
}

func (x *ListValue) ToStarlark() *starlark.List {
	if x == nil {
		return starlark.NewList(nil)
	}

	elems := make([]starlark.Value, len(x.Values))

	for i, v := range x.Values {
		elems[i] = v.ToStarlark()
	}

	return starlark.NewList(elems)
}

func (x *Struct) ToStarlark() *starlark.Dict {
	if x == nil {
		return starlark.NewDict(0)
	}

	dict := starlark.NewDict(len(x.Fields))

	for k, v := range x.Fields {
		_ = dict.SetKey(starlark.String(k), v.ToStarlark())
	}

	return dict
}

func (x *Value) ToStarlark() starlark.Value {
	if x == nil {
		return starlark.None
	}

	switch v := x.GetKind().(type) {
	case *Value_IntValue:
		return starlark.MakeInt64(v.IntValue)
	case *Value_FloatValue:
		return starlark.Float(v.FloatValue)
	case *Value_StringValue:
		return starlark.String(v.StringValue)
	case *Value_BoolValue:
		return starlark.Bool(v.BoolValue)
	case *Value_StructValue:
		return v.StructValue.ToStarlark()
	case *Value_ListValue:
		return v.ListValue.ToStarlark()
	default:
		return starlark.None
	}
}
