package safecast_test

import (
	"fmt"
	"testing"

	"fortio.org/safecast"
)

// TODO: steal the tests from https://github.com/ccoVeille/go-safecast

const all64bitsOne = ^uint64(0)

func FindNumIntBits[T safecast.Float](t *testing.T) int {
	var v T
	for i := 0; i < 64; i++ {
		bits := (all64bitsOne >> i)
		v = T(bits)
		if v != v-1 {
			return 64 - i
		}
	}
	panic("bug... didn't fine num bits")
}

func TestFloatBounds(t *testing.T) {
	float32bits := FindNumIntBits[float32](t)
	float64bits := FindNumIntBits[float64](t)
	t.Logf("float32: %d bits", float32bits)
	t.Logf("float64: %d bits", float64bits)
	f32, err := safecast.Convert[float32](all64bitsOne)
	if err == nil {
		t.Errorf("expected error, got %d -> %.0f", all64bitsOne, f32)
	}
	f64, err := safecast.Convert[float64](all64bitsOne)
	if err == nil {
		t.Errorf("expected error, got %d -> %.0f", all64bitsOne, f64)
	}
}

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
	ub = safecast.MustConvert[uint8](int64(255)) // shouldn't panic
	if ub != 255 {
		t.Errorf("unexpected value: %v", ub)
	}
}

func TestPanicMustRound(t *testing.T) {
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

func TestPanicMustTruncate(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic")
		} else {
			expected := "safecast: out of range for -1.5 (float32) to uint8"
			if r != expected {
				t.Errorf("unexpected panic: %q wanted %q", r, expected)
			}
		}
	}()
	safecast.MustTruncate[uint8](float32(-1.5))
}

func TestPanicMustConvert(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic")
		} else {
			expected := "safecast: out of range for 256 (int) to uint8"
			if r != expected {
				t.Errorf("unexpected panic: %q wanted %q", r, expected)
			}
		}
	}()
	safecast.MustConvert[uint8](256)
}

func Example() {
	var in int16 = 256
	// will error out
	out, err := safecast.Convert[uint8](in)
	fmt.Println(out, err)
	// will be fine
	out = safecast.MustRound[uint8](255.4)
	fmt.Println(out)
	// Also fine
	out = safecast.MustTruncate[uint8](255.6)
	fmt.Println(out)
	// Output: 0 out of range
	// 255
	// 255
}

func ExampleMustRound() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic:", r)
		}
	}()
	out := safecast.MustRound[int8](-128.6)
	fmt.Println("not reached", out) // not reached
	// Output: panic: safecast: out of range for -128.6 (float64) to int8
}
