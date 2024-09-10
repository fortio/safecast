package safecast_test

import (
	"testing"

	"github.com/ldemailly/go-scratch/safecast"
)

func TestConvert(t *testing.T) {
	var inp uint32 = 42
	out, err := safecast.Convert[int8](inp)
	t.Logf("Out is %T: %v", out, out)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != 42 {
		t.Errorf("unexpected value: %v", out)
	}
	inp = 129
	_, err = safecast.Convert[int8](inp)
	t.Logf("Got err: %v", err)
	if err == nil {
		t.Errorf("expected error")
	}
	inp2 := int32(-42)
	_, err = safecast.Convert[uint8](inp2)
	t.Logf("Got err: %v", err)
	if err == nil {
		t.Errorf("expected error")
	}
	out, err = safecast.Convert[int8](inp2)
	t.Logf("Out is %T: %v", out, out)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out != -42 {
		t.Errorf("unexpected value: %v", out)
	}
	inp2 = -129
	_, err = safecast.Convert[uint8](inp2)
	t.Logf("Got err: %v", err)
	if err == nil {
		t.Errorf("expected error")
	}
	var a uint16 = 65535
	x, err := safecast.Convert[int16](a)
	if err == nil {
		t.Errorf("expected error, %d %d", a, x)
	}
	b := int8(-1)
	y, err := safecast.Convert[uint](b)
	if err == nil {
		t.Errorf("expected error, %d %d", b, y)
	}
	up := uintptr(42)
	b, err = safecast.Convert[int8](up)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
	}
	if b != 42 {
		t.Errorf("unexpected value: %v", b)
	}
	b = -1
	_, err = safecast.Convert[uintptr](b)
	if err == nil {
		t.Errorf("expected err")
	}
	ub := safecast.MustTruncate[uint8](255.6)
	if ub != 255 {
		t.Errorf("unexpected value: %v", ub)
	}
}

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic")
		} else {
			expected := "safecast: out of range for 255.5 (float32) to uint8"
			if r != expected {
				t.Errorf("unexpected panic: %q wanted %q", r, expected)
			}
		}
	}()
	safecast.MustRound[uint8](float32(255.5))
}
