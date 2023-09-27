package types

import "strings"

type Func struct {
	params []Type
	ret    Type
}

func (f *Func) String() string {
	params := make([]string, len(f.params))
	for i, p := range f.params {
		params[i] = p.String()
	}
	pstr := strings.Join(params, ", ")

	if f.ret != nil {
		return "func(" + pstr + ") " + f.ret.String()
	}
	return "func(" + pstr + ")"
}

func (f *Func) Underlying() Type {
	return f
}

// ReturnType returns the return type of the function.
func (f *Func) ReturnType() Type {
	return f.ret
}

// ParamTypes returns the parameter types of the function.
func (f *Func) ParamTypes() []Type {
	return f.params
}
