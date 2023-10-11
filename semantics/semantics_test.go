package semantics_test

import (
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

func parse(t *testing.T, src string) (*ast.AST, *types.Universe, ast.NodeID, []error) {
	t.Helper()

	parser := parser.New(token.NewFile("test.gos", []byte(src)))
	a, errs := parser.Parse()
	if len(errs) > 0 {
		return nil, nil, ast.InvalidNode, errs
	}

	root := a.Root()
	if a.Kind(root) != ast.DeclList {
		t.Errorf("Expected DeclList, but got %s", a.Kind(root))
	}

	tc := semantics.NewTypeChecker(a)

	_, errs = tc.Check(root)
	if errs != nil {
		return nil, nil, ast.InvalidNode, errs
	}

	return a, tc.Universe(), root, nil
}

func parseStmt(t *testing.T, src string) (*ast.AST, *types.Universe, ast.NodeID, []error) {
	t.Helper()

	a, uni, root, errs := parse(t, "func foo() {"+src+"}")
	if len(errs) > 0 {
		return nil, nil, ast.InvalidNode, errs
	}

	decl := a.Child(root, 0)
	if a.Kind(decl) != ast.FuncDecl {
		t.Errorf("Expected FuncDecl, but got %s", a.Kind(decl))
	}

	body := a.Child(decl, ast.FuncDeclBody)
	if a.Kind(body) != ast.StmtList {
		t.Errorf("Expected StmtList, but got %s", a.Kind(root))
	}

	stmt := a.Child(body, a.NumChildren(body)-1)
	return a, uni, stmt, nil
}
