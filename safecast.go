package safecast

import (
	"errors"
)

type Integer interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
}

var ErrOutOfRange = errors.New("out of range")

func Convert[Tout Integer, Tin Integer](orig Tin) (converted Tout, err error) {
	converted = Tout(orig)
	if Tin(converted) != orig {
		err = ErrOutOfRange
	}
	return
}
