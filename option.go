package cli

import "reflect"

type option struct {
	short string
	long  string
	desc  string
}

func newOption(sf reflect.StructField, rv reflect.Value) (*option, error) {
	return &option{name: sf.Name}, nil
}
