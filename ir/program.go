package ir

import "github.com/rj45/gosling/types"

type Program struct {
	Graph

	Packages []Package

	Types *types.Universe
}

type Package struct {
	Graph

	*Program

	Funcs []Function
}

type Function struct {
	Graph

	*Package

	Name string
	Sig  types.Type
}
