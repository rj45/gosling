package types

type Symbol struct {
	Name   string
	Offset int
}

type SymbolID uint32

type SymTab struct {
	Sym     []Symbol
	NameSym map[string]SymbolID
}

func NewSymTab() *SymTab {
	return &SymTab{
		NameSym: make(map[string]SymbolID),
	}
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
