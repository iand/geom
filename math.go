package geom

import (
	"math"
)

// epsilon32 is the machine epsilon, or the upper bound on the relative error due to rounding in floating point arithmatic
var epsilon32 = float32(math.Nextafter32(1, 2) - 1)

const (
	pi    = 3.14159265358979323846264338327950288419716939937510582097494459 // http://oeis.org/A000796
	sqrt2 = 1.41421356237309504880168872420969807856967187537694807317667974 // http://oeis.org/A002193
	sqrt3 = 1.73205080756887729352744634150587236694280525381038062805580698 // http://oeis.org/A002194

	maxFloat32 = math.MaxFloat32
)

// Radians converts degrees into radians
func Radians(degrees float32) float32 {
	return degrees * pi / 180
}

// Degrees converts radians into degrees
func Degrees(radians float32) float32 {
	return radians * 180 / pi
}

// sqrt returns the positive square root of v
// TODO: replace with https://github.com/rkusa/gm
func sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}

// pow2 returns the next highest power of 2 or the number unchanged if it is already a power of 2.
// From https://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func pow2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	return v + 1
}

// cmp reports whether x and y are closely equal.
func cmp(a, b float32) bool {
	const maxDiff = 0.005
	const maxRelDiff = 1e-5

	diff := abs(a - b)
	if diff <= maxDiff {
		return true
	}
	a = abs(a)
	b = abs(b)

	largest := max(a, b)

	if diff <= largest*maxRelDiff {
		return true
	}
	return false
}

func clampZero(v float32) float32 {
	if cmp(v, 0) {
		return 0
	}
	return v
}

func clampZeroVec3(v Vec3) Vec3 {
	return Vec3{
		clampZero(v[0]),
		clampZero(v[1]),
		clampZero(v[2]),
	}
}

func ClampZeroVec2(v Vec2) Vec2 {
	return Vec2{
		clampZero(v[0]),
		clampZero(v[1]),
	}
}

// clamp constrains v to be >=l and <= u
func clamp(v, l, u float32) float32 {
	if v < l {
		return l
	}
	if v > u {
		return u
	}
	return v
}

// signbit32 returns true if x is negative or negative zero.
func signbit32(x float32) bool {
	return math.Float32bits(x)&(1<<31) != 0
}

// max returns the maximum of a or b
func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// min returns the minimum of a or b
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func abs(x float32) float32 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		// ensure abs(-0) = 0
		return 0
	}
	return x
}

// copysign returns a value with the magnitude
// of x and the sign of y.
func copysign(x, y float32) float32 {
	const sign = 1 << 31
	return math.Float32frombits(math.Float32bits(x)&^sign | math.Float32bits(y)&sign)
}

// nonzero returns v if it is non-zero, or a small number close to zero otherwise
func nonzero(v float32) float32 {
	if v != 0 {
		return v
	}

	return copysign(math.SmallestNonzeroFloat32, v)
}
