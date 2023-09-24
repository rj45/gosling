package types

type SymbolKind uint8

const (
	VarSymbol SymbolKind = iota
	ConstSymbol
	TypeSymbol
)

type Symbol struct {
	Kind   SymbolKind
	Name   string
	Type   Type
	Const  Const
	Offset int
}

type SymbolID uint32

type SymTab struct {
	Sym     []Symbol
	NameSym map[string]SymbolID
}

func NewSymTab() *SymTab {
	symtab := &SymTab{
		NameSym: make(map[string]SymbolID),
	}

	symtab.InitSym("true", ConstSymbol, Bool, BoolConst(true))
	symtab.InitSym("false", ConstSymbol, Bool, BoolConst(false))

	symtab.InitSym("int", TypeSymbol, Int, nil)

	return symtab
}

func (t *SymTab) Lookup(name string) *Symbol {
	id, ok := t.NameSym[name]
	if !ok {
		id = SymbolID(len(t.Sym))
		t.Sym = append(t.Sym, Symbol{Kind: VarSymbol, Name: name, Offset: -1})
		t.NameSym[name] = id
		return &t.Sym[id]
	}
	return &t.Sym[id]
}

func (t *SymTab) InitSym(name string, kind SymbolKind, typ Type, c Const) {
	sym := t.Lookup(name)
	sym.Type = typ
	sym.Const = c
}
