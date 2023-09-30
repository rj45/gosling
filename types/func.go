package types

import "strings"

type Func struct {
	uni    *Universe
	params []Type
	ret    Type
}

func (f *Func) String() string {
	params := make([]string, len(f.params))
	for i, p := range f.params {
		params[i] = f.uni.StringOf(p)
	}
	pstr := strings.Join(params, ", ")

	if f.ret != Void {
		return "func(" + pstr + ") " + f.uni.StringOf(f.ret)
	}
	return "func(" + pstr + ")"
}

// ReturnType returns the return type of the function.
func (f *Func) ReturnType() Type {
	return f.ret
}

// ParamTypes returns the parameter types of the function.
func (f *Func) ParamTypes() []Type {
	return f.params
}
