package semantics

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// checkExprChild checks types when a child is used in an expression.
// Useful for checking things that can be either statements or expressions.
func (tc *TypeChecker) checkExprChild(parent, child ast.NodeID) {
	switch tc.ast.Kind(child) {
	case ast.IfExpr:
		then := tc.ast.Child(child, ast.IfExprThen)
		els := tc.ast.Child(child, ast.IfExprElse)

		if els == ast.InvalidNode {
			// else branch will get the zero value of type
			return
		}

		thenType := tc.ast.Type(then)
		elsType := tc.ast.Type(els)
		if thenType != types.None && elsType != types.None {
			// todo: move this coercion logic to central place
			if thenType == types.UntypedInt && elsType == types.Int {
				tc.ast.SetType(then, types.Int)
				thenType = types.Int
			}

			if thenType == types.Int && elsType == types.UntypedInt {
				tc.ast.SetType(els, types.Int)
				elsType = types.Int
			}

			if thenType != elsType {
				tc.errorf(parent, "if branches have mismatched types: %s and %s", tc.uni.StringOf(thenType), tc.uni.StringOf(elsType))
			}
		}
	}
}

func (tc *TypeChecker) checkBinaryExpr(node ast.NodeID) {
	lhs := tc.ast.Type(tc.ast.Child(node, ast.BinaryExprLHS))
	rhs := tc.ast.Type(tc.ast.Child(node, ast.BinaryExprRHS))

	if lhs == types.None || rhs == types.None {
		return
	}

	lhs = tc.uni.Underlying(lhs)
	rhs = tc.uni.Underlying(rhs)

	switch tc.ast.Token(node).Kind() {
	case token.Eq, token.Ne:
		// todo: check if types are comparable and compatible
		tc.ast.SetType(node, types.Bool)
		return
	case token.Lt, token.Gt, token.Le, token.Ge:
		// todo: check if types are ordered and compatible
		tc.ast.SetType(node, types.Bool)
		return
	}

	// todo: this is far too naive
	if lhs == types.UntypedInt {
		tc.ast.SetType(node, rhs)
	}
	tc.ast.SetType(node, lhs)
}

func (tc *TypeChecker) checkUnaryExpr(node ast.NodeID) {
	child := tc.ast.Child(node, ast.UnaryExprExpr)
	tc.ast.SetType(node, tc.ast.Type(child))
}

func (tc *TypeChecker) checkDerefExpr(node ast.NodeID) {
	child := tc.ast.Child(node, ast.DerefExprExpr)
	typ := tc.ast.Type(child)
	if typ == types.None {
		return
	}

	if typ.Indirections() == 0 {
		tc.errorf(node, "cannot dereference non-pointer type %s", tc.uni.StringOf(typ))
		return
	}
	tc.ast.SetType(node, typ.Deref())
}

func (tc *TypeChecker) checkAddrExpr(node ast.NodeID) {
	child := tc.ast.Child(node, ast.AddrExprExpr)
	typ := tc.ast.Type(child)
	if typ == types.None {
		return
	}

	if tc.ast.Kind(child) != ast.Name {
		tc.errorf(node, "cannot take address of non-name")
		return
	}

	if typ.Indirections() >= 3 {
		tc.errorf(node, "cannot take address of triple pointer type %s", tc.uni.StringOf(typ))
		return
	}

	tc.ast.SetType(node, typ.Pointer())
}

func (tc *TypeChecker) checkCallExpr(node ast.NodeID) {
	name := tc.ast.Child(node, ast.CallExprName)

	sym := tc.symtab.Lookup(tc.ast.NodeString(name))
	if sym == nil {
		// remove the "undefined name" error
		tc.errs = tc.errs[:len(tc.errs)-1]
		tc.errorf(node, "cannot call undefined function %s", tc.ast.NodeString(name))
		return
	}

	typ := sym.Type

	if typ.Kind() != types.FuncType {
		tc.errorf(node, "cannot call non-function %s of type %s", tc.ast.NodeString(name), tc.uni.StringOf(typ))
		return
	}
	fnTyp := tc.uni.Func(typ)

	argsNode := tc.ast.Child(node, ast.CallExprArgs)
	args := tc.ast.Children(argsNode)
	if len(args) != len(fnTyp.ParamTypes()) {
		tc.errorf(node, "wrong number of arguments to %s: expected %d, got %d", tc.ast.NodeString(name), len(fnTyp.ParamTypes()), len(args))
		return
	}
	for i, arg := range args {
		typ := tc.ast.Type(arg)
		if typ == types.None {
			continue
		}
		if typ == types.UntypedInt && fnTyp.ParamTypes()[i] == types.Int {
			tc.ast.SetType(arg, types.Int)
			typ = types.Int
		}
		if typ != fnTyp.ParamTypes()[i] {
			tc.errorf(node, "wrong type for argument: expected %s, got %s", tc.uni.StringOf(fnTyp.ParamTypes()[i]), tc.uni.StringOf(typ))
		}
	}

	tc.ast.SetType(node, fnTyp.ReturnType())
}

func (tc *TypeChecker) checkLiteral(node ast.NodeID) {
	switch tc.ast.Token(node).Kind() {
	case token.Int:
		tc.ast.SetType(node, types.UntypedInt)
	default:
		panic("unknown literal type " + tc.ast.Token(node).Kind().String())
	}
}

func (tc *TypeChecker) checkName(node ast.NodeID) {
	sym := tc.symtab.Lookup(tc.ast.NodeString(node))
	if sym == nil {
		tc.errorf(node, "undefined name %s", tc.ast.NodeString(node))
		return
	}
	tc.ast.SetType(node, sym.Type)
}

func (tc *TypeChecker) checkBlock(node ast.NodeID) {
	children := tc.ast.Children(node)
	if len(children) == 0 {
		tc.ast.SetType(node, types.Void)
		return
	}

	for i, child := range children {
		if tc.ast.Kind(child) == ast.ReturnStmt {
			if i != len(children)-1 {
				tc.errorf(node, "return must be last statement in block")
			}
		}
	}

	tc.ast.SetType(node, tc.ast.Type(children[len(children)-1]))
}

func (tc *TypeChecker) checkIfExpr(node ast.NodeID) {
	cond := tc.ast.Child(node, ast.IfExprCond)
	condType := tc.ast.Type(cond)
	if condType == types.None {
		return
	}
	if tc.uni.Underlying(tc.ast.Type(cond)) != types.Bool {
		tc.errorf(node, "if condition must be bool but was %s", tc.uni.StringOf(condType))
		return
	}

	then := tc.ast.Child(node, ast.IfExprThen)
	if thenType := tc.ast.Type(then); thenType != types.None {
		tc.ast.SetType(node, thenType)
	}
}
