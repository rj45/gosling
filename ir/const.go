package ir

import (
	"strconv"
)

type Constant interface {
	String() string
	Value() any
	isConst()
}

type int64Const int64
type boolConst bool
type funcConst struct{ *Func }
type regConst RegMask

func (c int64Const) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (b boolConst) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (f funcConst) String() string {
	return f.Name
}

func (r regConst) String() string {
	return RegMask(r).String()
}

func (c int64Const) isConst() {}
func (b boolConst) isConst()  {}
func (f funcConst) isConst()  {}
func (r regConst) isConst()   {}

func (c int64Const) Value() any {
	return int64(c)
}

func (b boolConst) Value() any {
	return bool(b)
}

func (f funcConst) Value() any {
	return f.Func
}

func (r regConst) Value() any {
	return RegMask(r)
}

func IntConst(v int64) Constant {
	return int64Const(v)
}

func BoolConst(v bool) Constant {
	return boolConst(v)
}

func FuncConst(f *Func) Constant {
	return funcConst{f}
}

func RegConst(r RegMask) Constant {
	return regConst(r)
}

func Int64Value(c Constant) (int64, bool) {
	if t, ok := c.(int64Const); ok {
		return int64(t), true
	}
	return 0, false
}

func BoolValue(c Constant) (bool, bool) {
	if t, ok := c.(boolConst); ok {
		return bool(t), true
	}
	return false, false
}

func FuncValue(c Constant) (*Func, bool) {
	if t, ok := c.(funcConst); ok {
		return t.Func, true
	}
	return nil, false
}

func RegValue(c Constant) (RegMask, bool) {
	if t, ok := c.(regConst); ok {
		return RegMask(t), true
	}
	return 0, false
}
