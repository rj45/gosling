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

	if !tc.uni.IsAssignable(lhsType, rhsType) && tc.ast.Token(node).Kind() != token.Define {
		tc.errorf(node, "cannot assign %s to %s", tc.uni.StringOf(rhsType), tc.uni.StringOf(lhsType))
		return
	}

	tc.ast.SetType(node, tc.uni.Unify(lhsType, rhsType))
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
