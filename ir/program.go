package ir

import "github.com/rj45/gosling/types"

type Program struct {
	Packages []Package

	Types *types.Universe
}

type Package struct {
	*Program

	Funcs []Function
}

type Function struct {
	Graph

	*Package

	Name string
	Sig  types.Type

	Start NodeID
	End   NodeID
}
