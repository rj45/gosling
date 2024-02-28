package ir

type Use struct {
	g  *Graph
	id inputID
}

func (u Use) User() Node {
	return Node{u.g, u.g.uses[u.id].user}
}

func (u Use) Next() Use {
	return Use{u.g, u.g.uses[u.id].nextUse}
}

func (u Use) Def() Node {
	return Node{u.g, u.g.inputs[u.id]}
}

func (u Use) Head() Use {
	return Use{u.g, u.g.useHeads[u.g.inputs[u.id]]}
}
