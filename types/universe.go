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

func (u *Universe) Func(ret Type) *Func {
	for i, f := range u.funcs {
		if f.ret == ret {
			return &u.funcs[i]
		}
	}
	f := Func{ret: ret}
	u.funcs = append(u.funcs, f)
	return &u.funcs[len(u.funcs)-1]
}
