package types

type Symbol struct {
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

	symtab.InitSym("true", Bool, BoolConst(true))
	symtab.InitSym("false", Bool, BoolConst(false))

	return symtab
}

func (t *SymTab) Lookup(name string) *Symbol {
	id, ok := t.NameSym[name]
	if !ok {
		id = SymbolID(len(t.Sym))
		t.Sym = append(t.Sym, Symbol{Name: name, Offset: -1})
		t.NameSym[name] = id
		return &t.Sym[id]
	}
	return &t.Sym[id]
}

func (t *SymTab) InitSym(name string, typ Type, c Const) {
	sym := t.Lookup(name)
	sym.Type = typ
	sym.Const = c
}
