package geom

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (

	// xzPlane3 is a plane perpendicular to the y axis, facing positive y
	xzPlane3 = Plane3{
		Normal:   Vec3{0, 1, 0},
		Distance: 0,
	}
	// xyPlane3 is a plane perpendicular to the z axis, facing positive z
	xyPlane3 = Plane3{
		Normal:   Vec3{0, 0, 1},
		Distance: 0,
	}
	// yzPlane3 is a plane perpendicular to the x axis, facing positive x
	yzPlane3 = Plane3{
		Normal:   Vec3{1, 0, 0},
		Distance: 0,
	}
	// xzInvPlane3 is a plane perpendicular to the y axis, facing negative y
	xzInvPlane3 = Plane3{
		Normal:   Vec3{0, -1, 0},
		Distance: 0,
	}
	// xyInvPlane3 is a plane perpendicular to the z axis, facing negative z
	xyInvPlane3 = Plane3{
		Normal:   Vec3{0, 0, -1},
		Distance: 0,
	}
	// yzInvPlane3 is a plane perpendicular to the x axis, facing negative x
	yzInvPlane3 = Plane3{
		Normal:   Vec3{-1, 0, 0},
		Distance: 0,
	}

	// xRay3 is a ray along the x-axis starting at a negative x, pointing towards the origin
	xRay3 = Ray3{
		Origin:    Vec3{-100, 0, 0},
		Direction: Vec3{1, 0, 0},
	}

	// xInvRay3 is a ray along the x-axis starting at a positive x, pointing towards the origin
	xInvRay3 = Ray3{
		Origin:    Vec3{100, 0, 0},
		Direction: Vec3{-1, 0, 0},
	}

	// yRay3 is a ray along the y-axis starting at a negative y, pointing towards the origin
	yRay3 = Ray3{
		Origin:    Vec3{0, -100, 0},
		Direction: Vec3{0, 1, 0},
	}

	// yInvRay3 is a ray along the y-axis starting at a positive y, pointing towards the origin
	yInvRay3 = Ray3{
		Origin:    Vec3{0, 100, 0},
		Direction: Vec3{0, -1, 0},
	}

	// zRay3 is a ray along the z-axis starting at a negative z, pointing towards the origin
	zRay3 = Ray3{
		Origin:    Vec3{0, 0, -100},
		Direction: Vec3{0, 0, 1},
	}

	// zInvRay3 is a ray along the z-axis starting at a positive z, pointing towards the origin
	zInvRay3 = Ray3{
		Origin:    Vec3{0, 0, 100},
		Direction: Vec3{0, 0, -1},
	}

	// aaOBB is an OBB that is axis aligned
	aaOBB = OBB{
		Position:    Point3{0, 0, 0},
		Size:        Vec3{2, 2, 2},
		Orientation: mgl32.QuatIdent(),
	}

	// planeOBB is an OBB that is aligned along the xy plane but has no z depth
	planeOBB = OBB{
		Position:    Point3{0, 0, 0},
		Size:        Vec3{2, 2, 0},
		Orientation: mgl32.QuatIdent(),
	}

	// tiltyOBB is an OBB that is tilted by 45degress along the y axis
	tiltyOBB = OBB{
		Position:    Point3{0, 0, 0},
		Size:        Vec3{2, 2, 2},
		Orientation: mgl32.QuatRotate(pi/4, Y3),
	}
)

func TestPlane3Raycast(t *testing.T) {
	testCases := []struct {
		p   Plane3
		r   Ray3
		hit bool
	}{
		{p: xyPlane3, r: zRay3, hit: false}, // ray origin is behind the plane
		{p: xyInvPlane3, r: zRay3, hit: true},
		{p: xyPlane3, r: zInvRay3, hit: true},
		{p: xyInvPlane3, r: zInvRay3, hit: false}, // ray origin is behind the plane

		{p: xzPlane3, r: yRay3, hit: false},
		{p: xzInvPlane3, r: yRay3, hit: true},
		{p: xzPlane3, r: yInvRay3, hit: true},
		{p: xzInvPlane3, r: yInvRay3, hit: false},

		{p: yzPlane3, r: xRay3, hit: false},
		{p: yzInvPlane3, r: xRay3, hit: true},
		{p: yzPlane3, r: xInvRay3, hit: true},
		{p: yzInvPlane3, r: xInvRay3, hit: false},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			rr, hit := tc.p.Raycast(tc.r)
			if hit != tc.hit {
				t.Errorf("got hit %v, wanted %v [fail=%v]", hit, tc.hit, rr.Fail)
			}
		})
	}
}

func TestPlane3ContainsPoint(t *testing.T) {
	testCases := []struct {
		p   Plane3
		pt  Point3
		hit bool
	}{

		{p: xyPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: xyInvPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: xyInvPlane3, pt: Point3{0, 0, 1}, hit: false},
		{p: xyInvPlane3, pt: Point3{0, 1, 0}, hit: true},
		{p: xyInvPlane3, pt: Point3{1, 0, 0}, hit: true},

		{p: xzPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: xzInvPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: xzInvPlane3, pt: Point3{0, 0, 1}, hit: true},
		{p: xzInvPlane3, pt: Point3{0, 1, 0}, hit: false},
		{p: xzInvPlane3, pt: Point3{1, 0, 0}, hit: true},

		{p: yzPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: yzInvPlane3, pt: Point3{0, 0, 0}, hit: true},
		{p: yzInvPlane3, pt: Point3{0, 0, 1}, hit: true},
		{p: yzInvPlane3, pt: Point3{0, 1, 0}, hit: true},
		{p: yzInvPlane3, pt: Point3{1, 0, 0}, hit: false},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			hit := tc.p.ContainsPoint3(tc.pt)
			if hit != tc.hit {
				t.Errorf("got hit %v, wanted %v", hit, tc.hit)
			}
		})
	}
}

func TestOBBContainsPoint3(t *testing.T) {
	testCases := []struct {
		o   OBB
		pt  Point3
		hit bool
	}{

		{o: aaOBB, pt: Point3{0, 0, 0}, hit: true},
		{o: aaOBB, pt: Point3{1, 0, 0}, hit: true},
		{o: aaOBB, pt: Point3{0, 1, 0}, hit: true},
		{o: aaOBB, pt: Point3{0, 0, 1}, hit: true},

		{o: aaOBB, pt: Point3{3, 0, 0}, hit: false},
		{o: aaOBB, pt: Point3{0, 3, 0}, hit: false},
		{o: aaOBB, pt: Point3{0, 0, 3}, hit: false},

		{o: aaOBB, pt: Point3{2, 0, 0}, hit: true},
		{o: aaOBB, pt: Point3{0, 2, 0}, hit: true},
		{o: aaOBB, pt: Point3{0, 0, 2}, hit: true},
		{o: aaOBB, pt: Point3{2, 2, 2}, hit: true},

		{o: aaOBB, pt: Point3{2.01, 0, 0}, hit: false},
		{o: aaOBB, pt: Point3{0, 2.01, 0}, hit: false},
		{o: aaOBB, pt: Point3{0, 0, 2.01}, hit: false},

		{o: planeOBB, pt: Point3{0, 0, 0}, hit: true},
		{o: planeOBB, pt: Point3{1, 0, 0}, hit: true},
		{o: planeOBB, pt: Point3{0, 1, 0}, hit: true},
		{o: planeOBB, pt: Point3{0, 0, 0.1}, hit: false},
		{o: planeOBB, pt: Point3{0, 0, -0.1}, hit: false},

		{o: tiltyOBB, pt: Point3{0, 0, 0}, hit: true},
		{o: tiltyOBB, pt: Point3{1, 0, 0}, hit: true},
		{o: tiltyOBB, pt: Point3{0, 1, 0}, hit: true},
		{o: tiltyOBB, pt: Point3{0, 0, 1}, hit: true},

		{o: tiltyOBB, pt: Point3{2, 0, 0}, hit: true},
		{o: tiltyOBB, pt: Point3{0, 2, 0}, hit: true},
		{o: tiltyOBB, pt: Point3{0, 0, 2}, hit: true},
		{o: tiltyOBB, pt: Point3{2, 2, 2}, hit: false},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			hit := tc.o.ContainsPoint3(tc.pt)
			if hit != tc.hit {
				t.Errorf("got hit %v, wanted %v (pt: %+v)", hit, tc.hit, tc.pt)
			}
		})
	}
}

// Package level variable for assignment, which avoids benchmarks being optimized away
var bres interface{}

func BenchmarkOBBContainsPoint3(b *testing.B) {
	testCases := []struct {
		name string
		o    OBB
		pt   Point3
	}{
		{name: "axis-aligned-inside", o: aaOBB, pt: Point3{1, 1, 1}},
		{name: "axis-aligned-outside", o: aaOBB, pt: Point3{3, 3, 3}},
		{name: "tilty-inside", o: tiltyOBB, pt: Point3{1, 1, 1}},
		{name: "tilty-outside", o: tiltyOBB, pt: Point3{3, 3, 3}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var contained bool
			for i := 0; i < b.N; i++ {
				contained = tc.o.ContainsPoint3(tc.pt)
			}
			b.StopTimer()
			bres = contained
		})
	}
}

func BenchmarkAABBContainsPoint3(b *testing.B) {
	aa := AABB{
		Position: Point3{0, 0, 0},
		Size:     Vec3{2, 2, 2},
	}

	testCases := []struct {
		name string
		a    AABB
		pt   Point3
	}{
		{name: "inside", a: aa, pt: Point3{1, 1, 1}},
		{name: "outside", a: aa, pt: Point3{3, 3, 3}},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var contained bool
			for i := 0; i < b.N; i++ {
				contained = tc.a.ContainsPoint3(tc.pt)
			}
			b.StopTimer()
			bres = contained
		})
	}
}

func BenchmarkOBBAxes(b *testing.B) {
	testCases := []struct {
		name string
		o    OBB
	}{
		{name: "axis-aligned", o: aaOBB},
		{name: "tilty", o: tiltyOBB},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var axes []Vec3
			for i := 0; i < b.N; i++ {
				axes = tc.o.Axes()
			}
			b.StopTimer()
			bres = axes
		})
	}
}

func BenchmarkABBAxes(b *testing.B) {
	aa := AABB{
		Position: Point3{0, 0, 0},
		Size:     Vec3{2, 2, 2},
	}

	b.ReportAllocs()

	var axes []Vec3
	for i := 0; i < b.N; i++ {
		axes = aa.Axes()
	}
	b.StopTimer()
	bres = axes
}

func BenchmarkOBBNormals(b *testing.B) {
	testCases := []struct {
		name string
		o    OBB
	}{
		{name: "axis-aligned", o: aaOBB},
		{name: "tilty", o: tiltyOBB},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var normals []Vec3
			for i := 0; i < b.N; i++ {
				normals = tc.o.Normals()
			}
			b.StopTimer()
			bres = normals
		})
	}
}

func BenchmarkABBNormals(b *testing.B) {
	aa := AABB{
		Position: Point3{0, 0, 0},
		Size:     Vec3{2, 2, 2},
	}

	b.ReportAllocs()

	var normals []Vec3
	for i := 0; i < b.N; i++ {
		normals = aa.Normals()
	}
	b.StopTimer()
	bres = normals
}

func BenchmarkIntersectsBox3(b *testing.B) {
	aa1 := AABB{
		Position: Point3{0, 0, 0},
		Size:     Vec3{2, 2, 2},
	}
	aa2 := AABB{
		Position: Point3{1, 1, 1},
		Size:     Vec3{2, 2, 2},
	}
	aa3 := AABB{
		Position: Point3{0, 0, 5},
		Size:     Vec3{2, 2, 2},
	}
	o1 := OBB{
		Position:    Point3{1, 1, 1},
		Size:        Vec3{2, 2, 2},
		Orientation: mgl32.QuatRotate(pi/4, Y3),
	}
	o2 := OBB{
		Position:    Point3{0, 0, 5},
		Size:        Vec3{2, 2, 2},
		Orientation: mgl32.QuatRotate(pi/4, Y3),
	}

	testCases := []struct {
		name string
		a    Box3
		b    Box3
	}{
		{name: "aabb-aabb-intersect", a: &aa1, b: &aa2},
		{name: "aabb-aabb-nonintersect", a: &aa1, b: &aa3},
		{name: "obb-aligned-aabb-intersect", a: &aaOBB, b: &aa2},
		{name: "obb-aligned-aabb-nonintersect", a: &aaOBB, b: &aa3},
		{name: "obb-oriented-aabb-intersect", a: &tiltyOBB, b: &aa2},
		{name: "obb-oriented-aabb-nonintersect", a: &tiltyOBB, b: &aa3},
		{name: "obb-aligned-obb-oriented-intersect", a: &aaOBB, b: &o1},
		{name: "obb-aligned-obb-oriented-nonintersect", a: &aaOBB, b: &o2},
		{name: "obb-oriented-obb-oriented-intersect", a: &tiltyOBB, b: &o1},
		{name: "obb-oriented-obb-oriented-nonintersect", a: &tiltyOBB, b: &o2},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			var intersects bool
			for i := 0; i < b.N; i++ {
				intersects = IntersectsBox3(tc.a, tc.b)
			}
			b.StopTimer()
			bres = intersects
		})
	}
}
