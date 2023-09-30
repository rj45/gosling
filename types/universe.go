package types

// Universe represents the universe of Go types.
// It is used to avoid repeated allocations of basic types.
// It also allows to compare types by their address.
type Universe struct {
	funcs []Func
}

func NewUniverse() *Universe {
	return &Universe{}
}

func (u *Universe) FuncFor(params []Type, ret Type) Type {
outer:
	for i, f := range u.funcs {
		if len(f.params) != len(params) {
			continue
		}
		for j, p := range f.params {
			if p != params[j] {
				continue outer
			}
		}
		if f.ret == ret {
			return newType(FuncType, i, 0)
		}
	}
	f := Func{uni: u, params: params, ret: ret}
	u.funcs = append(u.funcs, f)
	return newType(FuncType, len(u.funcs)-1, 0)
}

func (u *Universe) Basic(t Type) *Basic {
	if t.Kind() != BasicType {
		panic("not a basic type")
	}
	return &basicInfos[t.Index()]
}

func (u *Universe) Func(t Type) *Func {
	if t.Kind() != FuncType {
		panic("not a func type")
	}
	return &u.funcs[t.Index()]
}

func (u *Universe) StringOf(t Type) string {
	prefix := ""
	for i := 0; i < t.Indirections(); i++ {
		prefix += "*"
	}

	switch t.Kind() {
	case BasicType:
		return prefix + u.Basic(t).String()
	case FuncType:
		return prefix + u.Func(t).String()
	default:
		panic("unknown type kind")
	}
}

func (u *Universe) Underlying(t Type) Type {
	// todo: implement this
	return t
}
