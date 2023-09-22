package types

// Pointer represents a pointer type.
type Pointer struct {
	elem Type
}

func (p *Pointer) Underlying() Type { return p }
func (p *Pointer) String() string   { return "*" + p.elem.String() }

// Elem returns the element type of p.
func (p *Pointer) Elem() Type { return p.elem }
