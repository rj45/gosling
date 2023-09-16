package ast

const (
	// BinaryExpr has LHS and RHS children
	BinaryExprLHS = 0
	BinaryExprRHS = 1

	// UnaryExpr has Expr child
	UnaryExprExpr = 0

	// StmtList has a list of statements

	// ExprStmt has Expr child
	ExprStmtExpr = 0

	// AssignStmt has LHS and RHS children
	AssignStmtLHS = 0
	AssignStmtRHS = 1
)
