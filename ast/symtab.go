package ast

import (
	"github.com/rj45/gosling/types"
)

type SymbolKind uint8

const (
	VarSymbol SymbolKind = iota
	FuncSymbol
	ConstSymbol
	TypeSymbol
)

type Symbol struct {
	ID     SymbolID
	Scope  ScopeID
	Kind   SymbolKind
	Name   string
	Type   types.Type
	Const  types.Const
	Offset int
}

type SymbolID uint32

type ScopeID uint32

const (
	InvalidScope = iota
	BuiltinScope
	GlobalScope
	LocalScope
)

type scope struct {
	parent     ScopeID
	node       NodeID
	level      int
	nameSym    map[string]SymbolID
	nextOffset int
}

type SymTab struct {
	sym    []Symbol
	scopes []scope
	scope  ScopeID

	nodeScope map[NodeID]ScopeID
}

func NewSymTab() *SymTab {
	symtab := &SymTab{
		// scope 0 is invalid
		scopes:    []scope{{}},
		nodeScope: make(map[NodeID]ScopeID),
	}

	// enter the builtin scope
	symtab.NewScope()

	symtab.NewSymbol("true", ConstSymbol, types.Bool).Const = types.BoolConst(true)
	symtab.NewSymbol("false", ConstSymbol, types.Bool).Const = types.BoolConst(false)

	symtab.NewSymbol("int", TypeSymbol, types.Int)
	symtab.NewSymbol("bool", TypeSymbol, types.Bool)

	return symtab
}

func (t *SymTab) NewScope() ScopeID {
	parent := &t.scopes[t.scope]

	id := ScopeID(len(t.scopes))
	t.scopes = append(t.scopes, scope{
		parent:  t.scope,
		level:   parent.level + 1,
		nameSym: make(map[string]SymbolID),
	})
	t.scope = id
	return id
}

func (t *SymTab) EnterScope(node NodeID) {
	if scope, ok := t.nodeScope[node]; ok {
		t.scope = scope
		return
	}
	scope := t.NewScope()
	t.nodeScope[node] = scope
	t.scopes[scope].node = node
}

func (t *SymTab) LeaveScope() {
	t.scope = t.scopes[t.scope].parent
}

func (t *SymTab) Scope() ScopeID {
	return t.scope
}

func (t *SymTab) Symbol(id SymbolID) *Symbol {
	return &t.sym[id]
}

func (t *SymTab) ScopeLevel(id ScopeID) int {
	return t.scopes[id].level
}

func (t *SymTab) ParentScope(id ScopeID) ScopeID {
	return t.scopes[id].parent
}

func (t *SymTab) ScopeNode(id ScopeID) NodeID {
	return t.scopes[id].node
}

func (t *SymTab) LocalScope() ScopeID {
	for scope := t.scope; t.scopes[scope].level > InvalidScope; scope = t.scopes[scope].parent {
		if t.scopes[scope].level == LocalScope {
			return scope
		}
	}
	return InvalidScope
}

func (t *SymTab) StackSize() int {
	localScopeID := t.LocalScope()
	if localScopeID == InvalidScope {
		return 0
	}

	localScope := &t.scopes[localScopeID]
	return localScope.nextOffset
}

func (t *SymTab) Lookup(name string) *Symbol {
	for scope := &t.scopes[t.scope]; scope.level > InvalidScope; scope = &t.scopes[scope.parent] {
		if id, ok := scope.nameSym[name]; ok {
			return &t.sym[id]
		}
	}

	return nil
}

func (t *SymTab) NewSymbol(name string, kind SymbolKind, typ types.Type) *Symbol {
	id := SymbolID(len(t.sym))

	offset := 0
	localScopeID := t.LocalScope()
	if localScopeID != InvalidScope {
		localScope := &t.scopes[localScopeID]
		offset = localScope.nextOffset
		localScope.nextOffset++
	}

	t.sym = append(t.sym, Symbol{ID: id, Scope: t.scope, Kind: kind, Name: name, Offset: offset})
	t.scopes[t.scope].nameSym[name] = id
	sym := &t.sym[id]
	sym.Type = typ
	return sym
}
