package safecast_test

import (
	"fmt"
	"math"
	"testing"

	"fortio.org/safecast"
)

// TODO: steal the tests from https://github.com/ccoVeille/go-safecast

const all64bitsOne = ^uint64(0)

// Interesting part is the "true" for the first line, which is why we have to change the
// code in Convert to handle that 1 special case.
// safecast_test.go:22: bits 64: 1111111111111111111111111111111111111111111111111111111111111111
// : 18446744073709551615 -> float64 18446744073709551616 true.
func FindNumIntBits[T safecast.Float](t *testing.T) int {
	var v T
	for i := 0; i < 64; i++ {
		bits := (all64bitsOne >> i)
		v = T(bits)
		t.Logf("bits %02d: %b : %d -> %T %.0f %t", 64-i, bits, bits, v, v, uint64(v) == bits)
		if v != v-1 {
			return 64 - i
		}
	}
	panic("bug... didn't fine num bits")
}

func TestFloat32Bounds(t *testing.T) {
	float32bits := FindNumIntBits[float32](t)
	t.Logf("float32: %d bits", float32bits)
	float32int := uint64(1<<(float32bits) - 1) // 24 bits
	for i := 0; i <= 64-float32bits; i++ {
		t.Logf("float32int %b %d", float32int, float32int)
		f := safecast.MustConvert[float32](float32int)
		t.Logf("float32int -> %.0f", f)
		float32int <<= 1
	}
}

func TestFloat64Bounds(t *testing.T) {
	float64bits := FindNumIntBits[float64](t)
	t.Logf("float64: %d bits", float64bits)
	float64int := uint64(1<<(float64bits) - 1) // 53 bits
	for i := 0; i <= 64-float64bits; i++ {
		t.Logf("float64int %b %d", float64int, float64int)
		f := safecast.MustConvert[float64](float64int)
		t.Logf("float64int -> %.0f", f)
		float64int <<= 1
	}
}

func TestNonIntegerFloat(t *testing.T) {
	_, err := safecast.Convert[int](math.Pi)
	if err == nil {
		t.Errorf("expected error")
	}
	truncPi := math.Trunc(math.Pi) // math.Trunc returns a float64
	i, err := safecast.Convert[int](truncPi)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != 3 {
		t.Errorf("unexpected value: %v", i)
	}
	i, err = safecast.Truncate[int](math.Pi)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != 3 {
		t.Errorf("unexpected value: %v", i)
	}
	i, err = safecast.Round[int](math.Phi)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if i != 2 {
		t.Errorf("unexpected value: %v", i)
	}
}

// MaxUint64 special case and also MaxInt64+1.
func TestMaxInt64(t *testing.T) {
	f32, err := safecast.Convert[float32](all64bitsOne)
	if err == nil {
		t.Errorf("expected error, got %d -> %.0f", all64bitsOne, f32)
	}
	f64, err := safecast.Convert[float64](all64bitsOne)
	if err == nil {
		t.Errorf("expected error, got %d -> %.0f", all64bitsOne, f64)
	}
	minInt64p1 := int64(math.MinInt64) + 1 // not a power of 2
	t.Logf("minInt64p1 %b %d", minInt64p1, minInt64p1)
	_, err = safecast.Convert[float64](minInt64p1)
	f64 = float64(minInt64p1)
	int2 := int64(f64)
	t.Logf("minInt64p1 -> %.0f %d", f64, int2)
	if err == nil {
		t.Errorf("expected error, got %d -> %.0f", minInt64p1, f64)
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
	inp2 := int32(-1)
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
	if out != -1 {
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
