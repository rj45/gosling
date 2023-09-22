package types

// Universe represents the universe of Go types.
// It is used to avoid repeated allocations of basic types.
// It also allows to compare types by their address.
type Universe struct {
	pointers map[Type]*Pointer
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
