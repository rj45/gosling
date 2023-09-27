package types

// Universe represents the universe of Go types.
// It is used to avoid repeated allocations of basic types.
// It also allows to compare types by their address.
type Universe struct {
	pointers map[Type]*Pointer
	funcs    []Func
}

func NewUniverse() *Universe {
	return &Universe{
		pointers: make(map[Type]*Pointer),
	}
}

func (u *Universe) Pointer(elem Type) *Pointer {
	p, ok := u.pointers[elem]
	if !ok {
		p = &Pointer{elem: elem}
		u.pointers[elem] = p
	}
	return p
}

func (u *Universe) Func(params []Type, ret Type) *Func {
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
			return &u.funcs[i]
		}
	}
	f := Func{params: params, ret: ret}
	u.funcs = append(u.funcs, f)
	return &u.funcs[len(u.funcs)-1]
}
