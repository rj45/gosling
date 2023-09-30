package semantics

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/types"
)

// defineFunc gathers function declarations in the symtab before checking them.
func (tc *TypeChecker) defineFunc(node ast.NodeID) {
	name := tc.ast.Child(node, ast.FuncDeclName)
	if tc.symtab.Lookup(tc.ast.NodeString(name)) != nil {
		tc.errorf(node, "cannot redefine function %s", tc.ast.NodeString(name))
	}

	paramsNode := tc.ast.Child(node, ast.FuncDeclParams)
	paramFields := tc.ast.Children(paramsNode)

	params := make([]types.Type, len(paramFields))
	for i, paramField := range paramFields {
		paramTyp := tc.ast.Child(paramField, ast.FieldTyp)
		tc.check(paramTyp)
		params[i] = tc.ast.Type(paramTyp)
	}

	ret := tc.ast.Child(node, ast.FuncDeclRet)

	var typ types.Type

	if ret == ast.InvalidNode {
		typ = tc.uni.FuncFor(params, types.Void)
	} else {
		tc.check(ret)
		typ = tc.uni.FuncFor(params, tc.ast.Type(ret))
	}

	tc.symtab.NewSymbol(tc.ast.NodeString(name), ast.FuncSymbol, typ)
}

func (tc *TypeChecker) defineFuncParams(node ast.NodeID) {
	paramsNode := tc.ast.Child(node, ast.FuncDeclParams)
	paramFields := tc.ast.Children(paramsNode)

	for _, paramField := range paramFields {
		paramName := tc.ast.Child(paramField, ast.FieldName)
		paramTyp := tc.ast.Child(paramField, ast.FieldTyp)

		typ := tc.ast.Type(paramTyp)

		tc.symtab.NewSymbol(tc.ast.NodeString(paramName), ast.VarSymbol, typ)
	}
}

func (tc *TypeChecker) checkFuncDecl(node ast.NodeID) {
	name := tc.ast.Child(node, ast.FuncDeclName)

	sym := tc.symtab.Lookup(tc.ast.NodeString(name))

	tc.ast.SetType(node, sym.Type)

	if sym.Type == types.None {
		return
	}

	retType := tc.uni.Func(sym.Type).ReturnType()
	if retType != types.Void && retType != types.None {
		body := tc.ast.Child(node, ast.FuncDeclBody)
		if !tc.returns(body) {
			tc.errorf(node, "missing return statement in function %s", tc.ast.NodeString(name))
		}
	}
}

func (tc *TypeChecker) returns(node ast.NodeID) bool {
	switch tc.ast.Kind(node) {
	case ast.ReturnStmt:
		return true
	case ast.StmtList:
		num := tc.ast.NumChildren(node)
		if num == 0 {
			return false
		}
		return tc.returns(tc.ast.Child(node, num-1))
	case ast.IfExpr:
		then := tc.ast.Child(node, ast.IfExprThen)
		els := tc.ast.Child(node, ast.IfExprElse)
		return tc.returns(then) && tc.returns(els)
	case ast.ForStmt:
		body := tc.ast.Child(node, ast.ForStmtBody)
		return tc.returns(body)
	}
	return false
}
