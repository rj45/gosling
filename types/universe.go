package types

// Universe represents the universe of Go types.
// It is used to avoid repeated allocations of basic types.
// It also allows to compare types by their ID.
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

// Unify returns the type that a and b can be unified to.
func (u *Universe) Unify(a, b Type) Type {
	if a == b {
		return a
	}

	if a.Kind() == BasicType && b.Kind() == BasicType {
		if a == UntypedInt && b == Int {
			return b
		}
		if a == Int && b == UntypedInt {
			return a
		}
	}

	return None
}

// IsAssignable returns whether src can be assigned to dst.
func (u *Universe) IsAssignable(dst, src Type) bool {
	switch dst.Kind() {
	case BasicType:
		if dst == Int && src == UntypedInt {
			return true
		}
	}
	return dst == src
}

// IsComparable returns whether t is comparable.
func (u *Universe) IsComparable(t Type) bool {
	return t.Kind() == BasicType && t != Void
}

// IsOrdered returns whether t is ordered.
func (u *Universe) IsOrdered(t Type) bool {
	return t.Kind() == BasicType && t != Void
}
