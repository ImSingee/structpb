package structpb

func NewEmptyList() *List {
	return &List{}
}

// NewStringList constructs a ListValue from a gstring Go slice.
// The slice elements are converted using NewValue.
func NewStringList(v []string) (*List, error) {
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

// AsSlice converts x to a general-purpose Go slice.
// The slice elements are converted by calling Value.AsInterface.
func (x *List) AsSlice() []interface{} {
	vs := make([]interface{}, len(x.GetValues()))
	for i, v := range x.GetValues() {
		vs[i] = v.AsInterface()
	}
	return vs
}
