package types

import (
	"strconv"
)

type Const interface {
	Type
	IsConst() bool
}

type int64Const int64
type boolConst bool

func (c int64Const) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (b boolConst) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (c int64Const) IsConst() bool { return true }
func (b boolConst) IsConst() bool  { return true }

func (c int64Const) Underlying() Type { return UntypedInt }
func (b boolConst) Underlying() Type  { return Bool }

func IntConst(v int64) Const {
	return int64Const(v)
}

func BoolConst(v bool) Const {
	return boolConst(v)
}

func Int64Value(c Const) (int64, bool) {
	if t, ok := c.(int64Const); ok {
		return int64(t), true
	}
	return 0, false
}

func BoolValue(c Const) (bool, bool) {
	if t, ok := c.(boolConst); ok {
		return bool(t), true
	}
	return false, false
}
