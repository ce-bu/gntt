package gntt_math

import (
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Ordered
}

func Min[T Number](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max[T Number](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}
