// Implementation: me (@ldemailly), idea: ccoVeille - https://github.com/ccoVeille/go-safecast
package safecast

import (
	"errors"
)

type Integer interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
}

var ErrOutOfRange = errors.New("out of range")

func Negative[T Integer](t T) bool {
	return t < 0
}

func SameSign[T1, T2 Integer](a T1, b T2) bool {
	return Negative(a) == Negative(b)
}

func Convert[IntOut Integer, IntIn Integer](orig IntIn) (converted IntOut, err error) {
	converted = IntOut(orig)
	if IntIn(converted) != orig {
		err = ErrOutOfRange
	}
	if !SameSign(orig, converted) {
		err = ErrOutOfRange
	}
	return
}

func MustConvert[IntOut Integer, IntIn Integer](orig IntIn) IntOut {
	converted, err := Convert[IntOut, IntIn](orig)
	if err != nil {
		panic(err)
	}
	return converted
}
