package structpb

// NewStringList constructs a ListValue from a gstring Go slice.
// The slice elements are converted using NewValue.
func NewStringList(v []string) (*ListValue, error) {
	x := &ListValue{Values: make([]*Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// AsSlice converts x to a general-purpose Go slice.
// The slice elements are converted by calling Value.AsInterface.
func (x *ListValue) AsSlice() []interface{} {
	vs := make([]interface{}, len(x.GetValues()))
	for i, v := range x.GetValues() {
		vs[i] = v.AsInterface()
	}
	return vs
}
