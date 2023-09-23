package semantics

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/errors"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// TypeChecker is a semantic analyzer that checks the AST for type errors.
type TypeChecker struct {
	uni  *types.Universe
	ast  *ast.AST
	errs []error
}

// NewTypeChecker creates a new TypeChecker.
func NewTypeChecker(ast *ast.AST) *TypeChecker {
	return &TypeChecker{
		uni: types.NewUniverse(),
		ast: ast,
	}
}

// Check checks the AST for type errors, labeling the AST with types.
func (tc *TypeChecker) Check(node ast.NodeID) []error {
	tc.check(node)
	errs := tc.errs
	tc.errs = nil
	return errs
}

func (tc *TypeChecker) errorf(node ast.NodeID, msg string, args ...interface{}) {
	tc.errs = append(tc.errs, errors.Newf(tc.ast.Src, tc.ast.Token(node), msg, args...))
}

func (tc *TypeChecker) check(node ast.NodeID) {
	if node == ast.InvalidNode || tc.ast.Type(node) != nil {
		return
	}

	for _, child := range tc.ast.Children(node) {
		tc.check(child)
	}

	// extra checks for children used in expressions
	if tc.ast.Kind(node) != ast.StmtList {
		for _, child := range tc.ast.Children(node) {
			tc.checkExprChild(node, child)
		}
	}

	switch tc.ast.Kind(node) {
	case ast.ExprList:
		// todo: maybe have a tuple type?
	case ast.BinaryExpr:
		tc.checkBinaryExpr(node)
	case ast.UnaryExpr:
		tc.checkUnaryExpr(node)
	case ast.DerefExpr:
		tc.checkDerefExpr(node)
	case ast.AddrExpr:
		tc.checkAddrExpr(node)
	case ast.CallExpr:
		tc.checkCallExpr(node)
	case ast.Literal:
		tc.checkLiteral(node)
	case ast.Name:
		tc.checkName(node)
	case ast.AssignStmt:
		tc.checkAssignStmt(node)
	case ast.IfExpr:
		tc.checkIfExpr(node)
	case ast.ForStmt:
		tc.checkForStmt(node)
	case ast.ReturnStmt:
		tc.checkReturnStmt(node)
	case ast.ExprStmt:
		tc.ast.SetType(node, tc.ast.Type(tc.ast.Child(node, ast.ExprStmtExpr)))
	case ast.StmtList:
		tc.ast.SetType(node, tc.ast.Type(tc.ast.Child(node, tc.ast.NumChildren(node)-1)))
	case ast.EmptyStmt:
		// nothing to do
	default:
		panic("todo: implement node kind " + tc.ast.Kind(node).String())
	}
}

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
		if thenType != nil && elsType != nil {
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
				tc.errorf(parent, "if branches have mismatched types: %s and %s", thenType, elsType)
			}
		}
	}
}

func (tc *TypeChecker) checkBinaryExpr(node ast.NodeID) {
	lhs := tc.ast.Type(tc.ast.Child(node, ast.BinaryExprLHS))
	rhs := tc.ast.Type(tc.ast.Child(node, ast.BinaryExprRHS))

	if lhs == nil || rhs == nil {
		return
	}

	lhs = lhs.Underlying()
	rhs = rhs.Underlying()

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
	if typ == nil {
		return
	}

	ptr, ok := typ.(*types.Pointer)
	if !ok {
		tc.errorf(node, "cannot dereference non-pointer type %s", typ)
		return
	}
	tc.ast.SetType(node, ptr.Elem())
}

func (tc *TypeChecker) checkAddrExpr(node ast.NodeID) {
	child := tc.ast.Child(node, ast.AddrExprExpr)
	typ := tc.ast.Type(child)
	if typ == nil {
		return
	}

	if tc.ast.Kind(child) != ast.Name {
		tc.errorf(node, "cannot take address of non-name")
		return
	}

	tc.ast.SetType(node, tc.uni.Pointer(typ))
}

func (tc *TypeChecker) checkCallExpr(node ast.NodeID) {
	name := tc.ast.Child(node, ast.CallExprName)

	sym := tc.ast.SymbolOf(name)
	if sym.Type == nil {
		tc.errorf(node, "cannot call undefined function %s", tc.ast.NodeString(name))
		return
	}

	typ := tc.ast.Type(name)
	if typ == nil {
		return
	}

	if _, ok := typ.(*types.Func); !ok {
		tc.errorf(node, "cannot call non-function %s of type %s", tc.ast.NodeString(name), typ)
		return
	}
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
	sym := tc.ast.SymbolOf(node)
	if sym == nil {
		tc.errorf(node, "undefined name %s", tc.ast.NodeString(node))
		return
	}
	tc.ast.SetType(node, sym.Type)
}

func (tc *TypeChecker) checkAssignStmt(node ast.NodeID) {
	lhs := tc.ast.Child(node, ast.AssignStmtLHS)
	rhs := tc.ast.Child(node, ast.AssignStmtRHS)

	isDefine := tc.ast.Token(node).Kind() == token.Define

	lhsType := tc.ast.Type(lhs)
	rhsType := tc.ast.Type(rhs)

	if rhsType == nil {
		return
	}

	if rhsType == types.UntypedInt && lhsType != types.UntypedInt {
		rhsType = types.Int
	}

	if tc.ast.Kind(lhs) == ast.Name {
		sym := tc.ast.SymbolOf(lhs)

		if isDefine && sym.Type != nil {
			tc.errorf(node, "cannot redefine %s", tc.ast.NodeString(lhs))
			return
		}

		if lhsType == nil {

			if sym.Type == nil {
				if isDefine {
					sym.Type = rhsType
					tc.ast.SetType(lhs, rhsType)
				} else {
					tc.errorf(node, "undefined name %s", tc.ast.NodeString(lhs))
					return
				}
			}

			lhsType = sym.Type
		}
	}

	if lhsType != rhsType {
		tc.errorf(node, "cannot assign %s to %s", rhsType, lhsType)
		return
	}

	tc.ast.SetType(node, lhsType)
}

func (tc *TypeChecker) checkIfExpr(node ast.NodeID) {
	cond := tc.ast.Child(node, ast.IfExprCond)
	condType := tc.ast.Type(cond)
	if condType == nil {
		return
	}
	if tc.ast.Type(cond).Underlying() != types.Bool {
		tc.errorf(node, "if condition must be bool but was %s", condType)
		return
	}

	then := tc.ast.Child(node, ast.IfExprThen)
	if thenType := tc.ast.Type(then); thenType != nil {
		tc.ast.SetType(node, thenType)
	}
}

func (tc *TypeChecker) checkForStmt(node ast.NodeID) {
	cond := tc.ast.Child(node, ast.ForStmtCond)
	condType := tc.ast.Type(cond)
	if condType == nil {
		return
	}
	if tc.ast.Type(cond).Underlying() != types.Bool {
		tc.errorf(node, "for condition must be bool but was %s", condType)
		return
	}
}

func (tc *TypeChecker) checkReturnStmt(node ast.NodeID) {
	children := tc.ast.Children(node)

	if len(children) != 1 {
		tc.errorf(node, "invalid return statement")
		return
	}

	typ := tc.ast.Type(children[0])
	if typ == nil {
		return
	}
	if typ != types.Int && typ != types.UntypedInt {
		tc.errorf(node, "cannot return %s", typ)
	}
	tc.ast.SetType(node, typ)
}
