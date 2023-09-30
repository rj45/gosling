package semantics

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

func (tc *TypeChecker) defineAssignStmt(node ast.NodeID) {
	if tc.ast.Token(node).Kind() != token.Define {
		return
	}

	lhs := tc.ast.Child(node, ast.AssignStmtLHS)

	if tc.ast.Kind(lhs) != ast.Name {
		tc.errorf(node, "cannot define non-name %s", tc.ast.NodeString(lhs))
		return
	}

	if tc.symtab.Lookup(tc.ast.NodeString(lhs)) != nil {
		tc.errorf(node, "cannot redefine %s", tc.ast.NodeString(lhs))
		return
	}

	rhs := tc.ast.Child(node, ast.AssignStmtRHS)
	tc.check(rhs)
	rhsType := tc.ast.Type(rhs)

	if rhsType == types.UntypedInt {
		rhsType = types.Int
	}

	tc.symtab.NewSymbol(tc.ast.NodeString(lhs), ast.VarSymbol, rhsType)
}

func (tc *TypeChecker) checkAssignStmt(node ast.NodeID) {
	lhs := tc.ast.Child(node, ast.AssignStmtLHS)
	rhs := tc.ast.Child(node, ast.AssignStmtRHS)

	lhsType := tc.ast.Type(lhs)
	rhsType := tc.ast.Type(rhs)

	if rhsType == types.None || lhsType == types.None {
		return
	}

	if rhsType == types.UntypedInt && lhsType != types.UntypedInt {
		rhsType = types.Int
	}

	if lhsType != rhsType && tc.ast.Token(node).Kind() != token.Define {
		tc.errorf(node, "cannot assign %s to %s", tc.uni.StringOf(rhsType), tc.uni.StringOf(lhsType))
		return
	}

	tc.ast.SetType(node, lhsType)
}

func (tc *TypeChecker) checkForStmt(node ast.NodeID) {
	cond := tc.ast.Child(node, ast.ForStmtCond)
	condType := tc.ast.Type(cond)
	if condType == types.None {
		return
	}
	if tc.uni.Underlying(tc.ast.Type(cond)) != types.Bool {
		tc.errorf(node, "for condition must be bool but was %s", tc.uni.StringOf(condType))
		return
	}
}

func (tc *TypeChecker) checkReturnStmt(node ast.NodeID) {
	children := tc.ast.Children(node)

	fnScope := tc.symtab.LocalScope()
	if fnScope == ast.InvalidScope {
		panic("return statement outside of function")
	}
	fnNode := tc.symtab.ScopeNode(fnScope)
	if fnNode == ast.InvalidNode {
		panic("function scope without function node")
	}

	fnName := tc.ast.Child(fnNode, ast.FuncDeclName)
	fnSym := tc.symtab.Lookup(tc.ast.NodeString(fnName))

	if fnSym == nil {
		panic("function scope without function symbol")
	}

	fnTypeT := fnSym.Type
	if fnTypeT == types.None {
		panic("function symbol without type")
	}
	fnType := tc.uni.Func(fnTypeT)
	retType := fnType.ReturnType()

	if retType == types.Void {
		if len(children) != 0 {
			tc.errorf(node, "cannot return value from void function")
		}
		tc.ast.SetType(node, types.Void)
		return
	}

	if len(children) != 1 {
		tc.errorf(node, "invalid return statement")
		return
	}

	typ := tc.ast.Type(children[0])
	if typ == types.None {
		return
	}

	// todo: centralize this coercion logic
	if retType == types.Int && typ == types.UntypedInt {
		typ = types.Int
	}

	if typ != retType {
		tc.errorf(node, "cannot return %s from function returning %s", tc.uni.StringOf(typ), tc.uni.StringOf(retType))
	}
	tc.ast.SetType(node, typ)
}
