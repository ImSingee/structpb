package structpb

import (
	"google.golang.org/protobuf/runtime/protoimpl"
	"unicode/utf8"
)

func NewEmptyStruct() *Struct {
	return &Struct{Fields: map[string]*Value{}}
}

// NewStruct constructs a Struct from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue.
func NewStruct(v map[string]interface{}) (*Struct, error) {
	x := &Struct{Fields: make(map[string]*Value, len(v))}
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", k)
		}
		var err error
		x.Fields[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// Get equals x.Fields[key] but can be called from nil safely
func (x *Struct) Get(key string) *Value {
	if x == nil || x.Fields == nil {
		return nil
	}
	return x.Fields[key]
}

// Set equals x.Fields[key] = value but is more safe
// return value is false only if struct == nil
func (x *Struct) Set(key string, value *Value) bool {
	if x == nil {
		return false
	}

	if x.Fields == nil {
		x.Fields = map[string]*Value{key: value}
	} else {
		x.Fields[key] = value
	}
	return true
}

// AsMap converts x to a general-purpose Go map.
// The map values are converted by calling Value.AsInterface.
func (x *Struct) AsMap() map[string]interface{} {
	vs := make(map[string]interface{})
	for k, v := range x.GetFields() {
		vs[k] = v.AsInterface()
	}
	return vs
}
