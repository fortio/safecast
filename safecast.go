// Implementation: me (@ldemailly), idea: ccoVeille - https://github.com/ccoVeille/go-safecast
package safecast

import (
	"errors"
	"fmt"
	"math"
)

// Same as golang.org/x/contraints.Integer but without importing the whole thing for 1 line.
type Integer interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64 | ~uintptr
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Integer | Float
}

var ErrOutOfRange = errors.New("out of range")

func Negative[T Number](t T) bool {
	return t < 0
}

func SameSign[T1, T2 Number](a T1, b T2) bool {
	return Negative(a) == Negative(b)
}

func Convert[NumOut Number, NumIn Number](orig NumIn) (converted NumOut, err error) {
	converted = NumOut(orig)
	if !SameSign(orig, converted) {
		err = ErrOutOfRange
		return
	}
	if NumIn(converted) != orig {
		err = ErrOutOfRange
	}
	return
}

func MustConvert[NumOut Number, NumIn Number](orig NumIn) NumOut {
	converted, err := Convert[NumOut, NumIn](orig)
	if err != nil {
		doPanic(err, orig, converted)
	}
	return converted
}

func Truncate[NumOut Number, NumIn Float](orig NumIn) (converted NumOut, err error) {
	return Convert[NumOut](math.Trunc(float64(orig)))
}

func Round[NumOut Number, NumIn Float](orig NumIn) (converted NumOut, err error) {
	return Convert[NumOut](math.Round(float64(orig)))
}

func MustTruncate[NumOut Number, NumIn Float](orig NumIn) NumOut {
	converted, err := Truncate[NumOut, NumIn](orig)
	if err != nil {
		doPanic(err, orig, converted)
	}
	return converted
}

func MustRound[NumOut Number, NumIn Float](orig NumIn) NumOut {
	converted, err := Round[NumOut, NumIn](orig)
	if err != nil {
		doPanic(err, orig, converted)
	}
	return converted
}

func doPanic[NumOut Number, NumIn Number](err error, orig NumIn, converted NumOut) {
	panic(fmt.Sprintf("safecast: %v for %v (%T) to %T", err, orig, orig, converted))
}
