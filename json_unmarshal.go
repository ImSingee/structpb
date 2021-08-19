package structpb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
)

var _ json.Unmarshaler = (*ListValue)(nil)
var _ json.Unmarshaler = (*Struct)(nil)
var _ json.Unmarshaler = (*Value)(nil)

func (x *ListValue) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, &x.Values)
}
func (x *ListValue) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, p []byte) error {
	return x.UnmarshalJSON(p)
}

func (x *Struct) UnmarshalJSON(p []byte) error {
	if x.Fields == nil {
		x.Fields = make(map[string]*Value)
	}
	return json.Unmarshal(p, &x.Fields)
}
func (x *Struct) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, p []byte) error {
	return x.UnmarshalJSON(p)
}

func (x *Value) UnmarshalJSON(p []byte) error {
	p = bytes.TrimSpace(p)
	if len(p) == 0 {
		return fmt.Errorf("unexpected end")
	}

	switch p[0] {
	case '"':
		var s string
		err := json.Unmarshal(p, &s)
		if err != nil {
			return err
		}
		x.Kind = &Value_StringValue{StringValue: s}
		return nil
	case '{':
		var s Struct
		err := s.UnmarshalJSON(p)
		if err != nil {
			return err
		}
		x.Kind = &Value_StructValue{StructValue: &s}
		return nil
	case '[':
		var l ListValue
		err := l.UnmarshalJSON(p)
		if err != nil {
			return err
		}
		x.Kind = &Value_ListValue{ListValue: &l}
		return nil
	}

	switch {
	case bytes.Equal(p, []byte("null")):
		x.Kind = &Value_NullValue{NullValue: NullValue_NULL_VALUE}
		return nil
	case bytes.Equal(p, []byte("true")):
		x.Kind = &Value_BoolValue{BoolValue: true}
		return nil
	case bytes.Equal(p, []byte("false")):
		x.Kind = &Value_BoolValue{BoolValue: false}
		return nil
	}

	num := json.Number(p)
	if intValue, err := num.Int64(); err == nil {
		x.Kind = &Value_IntValue{IntValue: intValue}
		return nil
	}
	if floatValue, err := num.Float64(); err == nil {
		x.Kind = &Value_FloatValue{FloatValue: floatValue}
		return nil
	}

	return fmt.Errorf("invalid json data %s", p)
}
func (x *Value) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, p []byte) error {
	return x.UnmarshalJSON(p)
}
