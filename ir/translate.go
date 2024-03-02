package ir

import (
	"strconv"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/errors"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

type scope struct {
	parent  *scope
	names   map[string]NodeID
	control NodeID
}

type SoNTranslator struct {
	ast     *ast.AST
	program *Program
	pkg     *Package
	fn      *Function
	types   *types.Universe

	errs []error

	stack []NodeID
	scope *scope
}

func NewSoNTranslator(ast *ast.AST) *SoNTranslator {
	prog := &Program{Types: types.NewUniverse()}
	prog.Graph.Init(GlobalScope, nil, nil, &prog.Graph)
	return &SoNTranslator{
		ast:     ast,
		program: prog,
		types:   prog.Types,
		scope:   &scope{names: make(map[string]NodeID)},
	}
}

func (t *SoNTranslator) errorf(node ast.NodeID, msg string, args ...interface{}) {
	t.errs = append(t.errs, errors.Newf(t.ast.Src, t.ast.Token(node), msg, args...))
}

func (t *SoNTranslator) Translate() *Program {
	t.program.Packages = append(t.program.Packages, Package{Program: t.program})
	t.pkg = &t.program.Packages[len(t.program.Packages)-1]
	t.pkg.Graph.Init(PackageScope, nil, &t.pkg.Graph, &t.program.Graph)
	t.translateRoot(t.ast.Root())
	return t.program
}

func (t *SoNTranslator) translateRoot(node ast.NodeID) {
	switch t.ast.Kind(node) {
	case ast.DeclList:
		t.translateDeclList(node)
		return
	}

	panic("fixme: bad root node")
}

func (t *SoNTranslator) translateDeclList(node ast.NodeID) NodeID {
	for _, child := range t.ast.Children(node) {
		t.translateDecl(child)
	}
	return InvalidNode
}

func (t *SoNTranslator) translateDecl(node ast.NodeID) {
	switch t.ast.Kind(node) {
	case ast.FuncDecl:
		t.translateFuncDecl(node)
		return
	}
	panic("todo: implement more decls")
}

func (t *SoNTranslator) translateFuncDecl(node ast.NodeID) {
	t.pkg.Funcs = append(t.pkg.Funcs, Function{Package: t.pkg})
	t.fn = &t.pkg.Funcs[len(t.pkg.Funcs)-1]

	t.fn.Init(LocalScope, &t.fn.Graph, &t.pkg.Graph, &t.program.Graph)

	name := t.ast.Child(node, ast.FuncDeclName)
	params := t.ast.Child(node, ast.FuncDeclParams)
	body := t.ast.Child(node, ast.FuncDeclBody)

	t.fn.Name = t.ast.NodeString(name)

	t.fn.Sig = t.translateType(node)

	for _, child := range t.ast.Children(params) {
		// todo: make nodes for params
		_ = child
	}

	t.fn.start = t.fn.NewNodeWithIDs(OpStart, types.Void, t.ast.Token(node))
	t.startScope(t.fn.start)

	for _, child := range t.ast.Children(body) {
		t.translateStmt(child)
	}

	t.fn.end = t.endScope().control
}

func (t *SoNTranslator) translateStmtList(node ast.NodeID) {
	for _, child := range t.ast.Children(node) {
		t.translateStmt(child)
	}
}

func (t *SoNTranslator) translateStmt(node ast.NodeID) {
	switch t.ast.Kind(node) {
	case ast.ExprStmt:
		expr := t.ast.Child(node, ast.ExprStmtExpr)
		t.translateExpr(expr)
		return
	case ast.IfExpr:
		t.translateExpr(node)
		return
	case ast.AssignStmt:
		t.translateAssignStmt(node)
		return
	case ast.ReturnStmt:
		t.translateReturnStmt(node)
		return
	case ast.StmtList:
		t.translateStmtList(node)
		return
	case ast.EmptyStmt:
		return
	}
	panic("todo: implement stmt " + t.ast.Kind(node).String())
}

func (t *SoNTranslator) translateType(node ast.NodeID) types.Type {
	switch t.ast.Kind(node) {
	case ast.Name:
		return t.types.DeferredFor(t.ast.NodeString(node))
	case ast.FuncDecl:
		paramNode := t.ast.Child(node, ast.FuncDeclParams)
		retNode := t.ast.Child(node, ast.FuncDeclRet)

		params := make([]types.Type, 0, t.ast.NumChildren(paramNode))
		for _, child := range t.ast.Children(paramNode) {
			params = append(params, t.translateType(child))
		}
		ret := t.translateType(retNode)

		return t.types.FuncFor(params, ret)

	default:
		panic("todo: implement type: " + t.ast.StringOf(node))
	}
}

func (t *SoNTranslator) translateReturnStmt(node ast.NodeID) {
	rets := make([]NodeID, 0, t.ast.NumChildren(node)+1)
	rets = append(rets, t.scope.control)
	for _, child := range t.ast.Children(node) {
		t.translateExpr(child)
		rets = append(rets, t.pop())
	}
	t.scope.control = t.fn.NewNodeWithIDs(OpReturn, types.Void, t.ast.Token(node), rets...)
}

func (t *SoNTranslator) translateAssignStmt(node ast.NodeID) {
	lhs := t.ast.Child(node, ast.AssignStmtLHS)
	rhs := t.ast.Child(node, ast.AssignStmtRHS)

	if t.ast.Kind(lhs) != ast.Name {
		t.errorf(node, "cannot define non-name %s", t.ast.NodeString(lhs))
	}

	t.translateExpr(rhs)

	name := t.ast.NodeString(lhs)

	val := t.pop()

	id := t.fn.NewNodeWithIDs(OpCopy, types.Unknown, t.ast.Token(lhs), val)

	t.fn.reassignName(id, name)

	t.scope.names[name] = id
}

func (t *SoNTranslator) pop() NodeID {
	id := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	return id
}

func (t *SoNTranslator) translateExpr(node ast.NodeID) {
	switch t.ast.Kind(node) {
	case ast.BinaryExpr:
		t.translateBinaryExpr(node)
		return
	case ast.UnaryExpr:
		t.translateUnaryExpr(node)
		return
	case ast.Literal:
		t.translateLiteral(node)
		return
	case ast.Name:
		t.translateName(node)
		return
	case ast.IfExpr:
		t.translateIfExpr(node)
		return
	}
	panic("todo: implement more exprs")
}

func (t *SoNTranslator) translateBinaryExpr(node ast.NodeID) {
	opTok := t.ast.Token(node)
	lhs := t.ast.Child(node, ast.BinaryExprLHS)
	rhs := t.ast.Child(node, ast.BinaryExprRHS)

	t.translateExpr(lhs)
	lhsNode := t.pop()
	t.translateExpr(rhs)
	rhsNode := t.pop()

	var op Op
	switch opTok.Kind() {
	case token.Add:
		op = OpAdd
	case token.Sub:
		op = OpSub
	case token.Star:
		op = OpMul
	case token.Div:
		op = OpDiv
	default:
		panic("unimplemented binary operator " + opTok.Kind().String())
	}

	newNode := t.fn.NewNodeWithIDs(op, types.Unknown, opTok, lhsNode, rhsNode)

	t.stack = append(t.stack, newNode)
}

func (t *SoNTranslator) translateUnaryExpr(node ast.NodeID) {
	opTok := t.ast.Token(node)
	expr := t.ast.Child(node, ast.UnaryExprExpr)

	t.translateExpr(expr)
	exprNode := t.pop()

	var op Op
	switch opTok.Kind() {
	case token.Sub:
		op = OpNeg
	default:
		panic("unimplemented unary operator " + opTok.Kind().String())
	}

	newNode := t.fn.NewNodeWithIDs(op, types.Unknown, opTok, exprNode)

	t.stack = append(t.stack, newNode)
}

func (t *SoNTranslator) translateLiteral(node ast.NodeID) {
	var newNode NodeID
	switch t.ast.Token(node).Kind() {
	case token.Int:
		value, err := strconv.ParseInt(t.ast.NodeString(node), 10, 64)
		if err != nil {
			t.errorf(node, "invalid int literal: %v", err)
		}
		c := types.IntConst(value)
		typ := t.types.ConstFor(c)
		newNode = t.fn.NewNodeWithIDs(OpConst, typ, t.ast.Token(node))
	default:
		panic("unimplemented literal type " + t.ast.Token(node).Kind().String())
	}

	t.stack = append(t.stack, newNode)
}

func (t *SoNTranslator) translateName(node ast.NodeID) {
	name := t.ast.NodeString(node)
	if name == "true" || name == "false" {
		t.translateBooleanLiteral(node)
		return
	}
	scope := t.scope
	for scope != nil {
		id, ok := t.scope.names[name]
		if ok {
			t.stack = append(t.stack, id)
			return
		}
	}

	t.errorf(node, "undefined name %s", name)
	id := t.fn.NewNodeWithIDs(OpError, types.Unknown, t.ast.Token(node))
	t.stack = append(t.stack, id)
}

func (t *SoNTranslator) translateBooleanLiteral(node ast.NodeID) {
	var value bool
	switch t.ast.NodeString(node) {
	case "true":
		value = true
	case "false":
		value = false
	default:
		t.errorf(node, "invalid boolean literal %s", t.ast.NodeString(node))
	}

	c := types.BoolConst(value)
	typ := t.types.ConstFor(c)
	newNode := t.fn.NewNodeWithIDs(OpConst, typ, t.ast.Token(node))

	t.stack = append(t.stack, newNode)
}

func (t *SoNTranslator) startScope(control NodeID) {
	t.scope = &scope{parent: t.scope, names: make(map[string]NodeID), control: control}
}

func (t *SoNTranslator) endScope() *scope {
	oldScope := t.scope
	t.scope = t.scope.parent
	// todo: handle stack remains
	return oldScope
}

func (t *SoNTranslator) translateIfExpr(node ast.NodeID) {
	cond := t.ast.Child(node, ast.IfExprCond)
	then := t.ast.Child(node, ast.IfExprThen)
	els := t.ast.Child(node, ast.IfExprElse)

	t.translateExpr(cond)
	condNode := t.pop()

	ifNode := t.fn.NewNodeWithIDs(OpIf, types.Unknown, t.ast.Token(node), t.scope.control, condNode)
	thenStart := t.fn.NewNodeWithIDs(OpThen, types.Unknown, t.ast.Token(node), ifNode)

	t.startScope(thenStart)
	t.translateStmt(then)
	thenEnd := t.scope.control
	thenScope := t.endScope()

	var elseEnd NodeID
	var elseScope *scope
	if els != ast.InvalidNode {
		elseStart := t.fn.NewNodeWithIDs(OpElse, types.Unknown, t.ast.Token(node), ifNode)
		t.startScope(elseStart)
		t.translateStmt(els)
		elseEnd = t.scope.control
		elseScope = t.endScope()
	}

	var region NodeID
	if els == ast.InvalidNode {
		region = t.fn.NewNodeWithIDs(OpRegion, types.Unknown, t.ast.Token(node), t.scope.control, thenEnd)
	} else {
		region = t.fn.NewNodeWithIDs(OpRegion, types.Unknown, t.ast.Token(node), thenEnd, elseEnd)
	}

	var phiArgs map[string][]NodeID
	for _, scope := range []*scope{thenScope, elseScope} {
		if scope == nil {
			continue
		}
		for name, id := range scope.names {
			if t.scope.names[name] != InvalidNode && t.scope.names[name] != scope.names[name] {
				if phiArgs == nil {
					phiArgs = make(map[string][]NodeID)
				}
				if _, ok := phiArgs[name]; !ok {
					phiArgs[name] = []NodeID{region, t.scope.names[name]}
				}
				phiArgs[name] = append(phiArgs[name], id)
			}
		}
	}

	for name, args := range phiArgs {
		phi := t.fn.NewNodeWithIDs(OpPhi, types.Unknown, t.ast.Token(node), args...)
		t.fn.reassignName(phi, name)
		t.scope.names[name] = phi
	}

	t.scope.control = region
}
