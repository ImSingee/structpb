package structpb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
)

var _ json.Unmarshaler = (*List)(nil)
var _ json.Unmarshaler = (*Dict)(nil)
var _ json.Unmarshaler = (*Value)(nil)

func (x *List) UnmarshalJSON(p []byte) error {
	return json.Unmarshal(p, &x.Values)
}
func (x *List) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, p []byte) error {
	return x.UnmarshalJSON(p)
}

func (x *Dict) UnmarshalJSON(p []byte) error {
	if x.Fields == nil {
		x.Fields = make(map[string]*Value)
	}
	return json.Unmarshal(p, &x.Fields)
}
func (x *Dict) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, p []byte) error {
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
		var s Dict
		err := s.UnmarshalJSON(p)
		if err != nil {
			return err
		}
		x.Kind = &Value_DictValue{DictValue: &s}
		return nil
	case '[':
		var l List
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
