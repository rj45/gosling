//go:build !nodebug

package assert

import "fmt"

func Equal[T comparable](a, b T) {
	if a != b {
		panic(fmt.Sprintf("assertion failed: %v == %v", a, b))
	}
}

func NotEqual[T comparable](a, b T) {
	if a == b {
		panic(fmt.Sprintf("assertion failed: %v != %v", a, b))
	}
}

func NonZero[T Zeroable](a T) {
	if a == 0 {
		panic(fmt.Sprintf("assertion failed: %v is non-zero", a))
	}
}

func True(a bool) {
	if !a {
		panic("assertion failed")
	}
}

func False(a bool) {
	if a {
		panic("assertion failed")
	}
}
