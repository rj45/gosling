// Package ir is a multi-level IR suitable for storing an
// intermediate representation for every level of the compiler
// after parsing. This includes a higher level IR, a low level
// IR closer to the machine, and the final assembly.
//
// The IR is designed to be as compact as possible in order
// to efficiently utilize the CPU cache. ValueIDs are used to
// identify values, and are simply an index into the current
// function's value list. This ValueID can be used in other
// arrays to store additional information about the value.
//
// If a value has a constant number of operands, it is preferable
// to identify operand indices with a constant so the code is more
// readable. See ast/operands.go for an example.
package ir
