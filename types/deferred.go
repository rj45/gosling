package types

// Deferred represents a type that is not yet resolved.
type Deferred struct {
	uni  *Universe
	name string
}

func (d Deferred) String() string {
	return d.name
}
