package structpb

// NewList constructs a ListValue from a general-purpose Go slice.
// The slice elements are converted using NewValue.
func NewList(v []interface{}) (*List, error) {
	x := &List{Values: make([]*Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

func New2DList(v [][]interface{}) (*List, error) {
	x := &List{Values: make([]*Value, len(v))}
	for i, v := range v {
		v2, err := NewList(v)
		if err != nil {
			return nil, err
		}
		x.Values[i] = NewListValue(v2)
	}
	return x, nil
}
