package semantics

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/errors"
	"github.com/rj45/gosling/types"
)

// TypeChecker is a semantic analyzer that checks the AST for type errors.
type TypeChecker struct {
	uni    *types.Universe
	symtab *ast.SymTab
	ast    *ast.AST
	errs   []error
}

// NewTypeChecker creates a new TypeChecker.
func NewTypeChecker(a *ast.AST) *TypeChecker {
	return &TypeChecker{
		uni:    types.NewUniverse(),
		ast:    a,
		symtab: ast.NewSymTab(),
	}
}

func (tc *TypeChecker) Universe() *types.Universe {
	return tc.uni
}

// Check checks the AST for type errors, labeling the AST with types.
func (tc *TypeChecker) Check(node ast.NodeID) (*ast.SymTab, []error) {
	tc.check(node)
	errs := tc.errs
	tc.errs = nil
	return tc.symtab, errs
}

func (tc *TypeChecker) errorf(node ast.NodeID, msg string, args ...interface{}) {
	tc.errs = append(tc.errs, errors.Newf(tc.ast.Src, tc.ast.Token(node), msg, args...))
}

func (tc *TypeChecker) check(node ast.NodeID) {
	if node == ast.InvalidNode || tc.ast.Type(node) != types.None {
		return
	}

	// pre-checks before checking children
	switch tc.ast.Kind(node) {
	case ast.DeclList:
		tc.symtab.EnterScope(node)
		defer tc.symtab.LeaveScope()
		for _, child := range tc.ast.Children(node) {
			if tc.ast.Kind(child) == ast.FuncDecl {
				// define functions first
				// this way they can be in any order
				tc.defineFunc(child)
			}
		}

	case ast.AssignStmt:
		// ensure defined variables are created in the symtab
		tc.defineAssignStmt(node)
	case ast.FuncDecl:
		tc.symtab.EnterScope(node)
		defer tc.symtab.LeaveScope()
		tc.defineFuncParams(node)
	case ast.StmtList:
		tc.symtab.EnterScope(node)
		defer tc.symtab.LeaveScope()
	}

	// check children
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
	case ast.DeclList:
		// nothing to do
	case ast.FuncDecl:
		tc.checkFuncDecl(node)
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
		tc.checkBlock(node)
	case ast.Field:
		tc.ast.SetType(node, tc.ast.Type(tc.ast.Child(node, ast.FieldTyp)))
	case ast.EmptyStmt, ast.FieldList:
		// nothing to do
	default:
		panic("todo: implement node kind " + tc.ast.Kind(node).String())
	}
}
