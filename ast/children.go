package ast

const (
	// BinaryExpr has LHS and RHS children
	BinaryExprLHS = 0
	BinaryExprRHS = 1

	// UnaryExpr has Expr child
	UnaryExprExpr = 0

	// StmtList has a list of statements

	// EmptyStmt has no children

	// ExprStmt has Expr child
	ExprStmtExpr = 0

	// AssignStmt has LHS and RHS children
	AssignStmtLHS = 0
	AssignStmtRHS = 1

	// ReturnStmt has list of Expr children

	// IfExpr has Cond, Then, and Else children
	IfExprCond = 0
	IfExprThen = 1
	IfExprElse = 2

	// ForStmt has Init, Cond, Post, and Body children
	ForStmtInit = 0
	ForStmtCond = 1
	ForStmtPost = 2
	ForStmtBody = 3

	// DerefExpr has Expr child
	DerefExprExpr = 0

	// AddrExpr has Expr child
	AddrExprExpr = 0
)
