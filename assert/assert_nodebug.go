//go:build nodebug

package assert

func Equal[T comparable](a, b T) {}

func NotEqual[T comparable](a, b T) {}

func NonZero[T Zeroable](a T) {}

func True(a bool) {}

func False(a bool) {}
