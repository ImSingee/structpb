package structpb

import (
	"bytes"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
)

func marshalJSON(m *jsonpb.Marshaler, data interface{}) ([]byte, error) {
	result, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if m != nil && m.Indent != "" {
		buf := new(bytes.Buffer)
		err = json.Indent(buf, result, "", m.Indent)
		if err != nil {
			return result, nil
		} else {
			return buf.Bytes(), nil
		}
	} else {
		return result, nil
	}
}

func (x *Dict) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	return marshalJSON(m, x.AsMap())
}
func (x *Dict) MarshalJSON() ([]byte, error) { return x.MarshalJSONPB(nil) }

func (x *List) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	return marshalJSON(m, x.AsSlice())
}
func (x *List) MarshalJSON() ([]byte, error) { return x.MarshalJSONPB(nil) }

func (x *Value) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	return marshalJSON(m, x.AsInterface())
}
func (x *Value) MarshalJSON() ([]byte, error) { return x.MarshalJSONPB(nil) }
