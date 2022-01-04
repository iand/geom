package geom

import (
	"github.com/go-gl/mathgl/mgl32"
)

type (
	Vec2   = mgl32.Vec2
	Vec3   = mgl32.Vec3
	Vec4   = mgl32.Vec4
	Mat4   = mgl32.Mat4
	Point2 = Vec2
	Point3 = Vec3
	Quat   = mgl32.Quat
)

type (
	Vec2i   [2]int32
	Vec3i   [3]int32
	Point2i = Vec2i
	Point3i = Vec3i
)

var (
	X3 = Vec3{1, 0, 0} // X-Axis in 3 dimensions
	Y3 = Vec3{0, 1, 0} // Y-Axis in 3 dimensions
	Z3 = Vec3{0, 0, 1} // Z-Axis in 3 dimensions
	X2 = Vec2{1, 0}    // X-Axis in 3 dimensions
	Y2 = Vec2{0, 1}    // Y-Axis in 3 dimensions
)

type Interval struct {
	Min, Max float32
}

func (i *Interval) Overlaps(i2 Interval) bool {
	return ((i2.Min <= i.Max) && (i.Min <= i2.Max))
}

// GetOverlap returns the amount of overlap between two intervals
func (i *Interval) GetOverlap(i2 Interval) float32 {
	if !i.Overlaps(i2) {
		return 0
	}
	return min(i.Max, i2.Max) - max(i.Min, i2.Min)
}

// OverlapOnAxis projects a and b onto the axis and tests if they overlap.
// If they do not overlap then we can guarantee that a and b do not overlap
func OverlapOnAxis(a, b Projecter, axis Vec3) bool {
	i1 := a.ProjectOntoAxis(axis)
	i2 := b.ProjectOntoAxis(axis)
	return i1.Overlaps(i2)
}

type Projecter interface {
	ProjectOntoAxis(axis Vec3) Interval
}

type Raycastable interface {
	Raycast(ray Ray3) (RaycastResult, bool)
}

// Box3 is a 3 dimensional cuboid
type Box3 interface {
	Projecter
	Axes() []Vec3
	Corners() []Point3
	Normals() []Vec3
	ContainsPoint3(pt Point3) bool
	// TODO
	// Raycastable
	Raycast(ray Ray3) (RaycastResult, bool)
}

// IntersectsBox3 uses the Separating Axis Theorem (SAT) which tests the axes from a, from b and from
// the cross-products of the axes from the two objects. Two objects only overlap if all axes
// overlap. See http://www.dyn4j.org/2010/01/sat/
func IntersectsBox3(a, b Box3) bool {
	axesa := a.Axes()
	axesb := b.Axes()

	for j := 0; j < len(axesb); j++ {
		if !OverlapOnAxis(a, b, axesb[j]) {
			// A separating axis was found
			return false
		}
	}

	for i := 0; i < len(axesa); i++ {
		if !OverlapOnAxis(a, b, axesa[i]) {
			// A separating axis was found
			return false
		}

		// Check the cross product of this axis with each of b's axes
		for j := 0; i < len(axesb); i++ {
			if !OverlapOnAxis(a, b, axesb[j].Cross(axesa[i])) {
				// A separating axis was found
				return false
			}
		}
	}

	// No separating axis was fund
	return true
}

// Ray2 is 2 dimensional ray that starts from the origin and projects an infinite distance in the specified direction.
type Ray2 struct {
	Origin    Point2
	Direction Vec2 // The direction of the ray, always normalised
}

// Point returns the coordinates of the point at a distance d from the ray's origin.
func (r *Ray2) Point(d float32) Point2 {
	return r.Origin.Add(r.Direction.Mul(d))
}

// Ray3 is 3 dimensional ray that starts from the origin and projects an infinite distance in the specified direction.
type Ray3 struct {
	Origin    Point3
	Direction Vec3 // The direction of the ray, always normalised
}

// Point returns the coordinates of the point at a distance d from the ray's origin.
func (r *Ray3) Point(d float32) Point3 {
	return r.Origin.Add(r.Direction.Mul(d))
}

// ClosestPoint returns the point along the ray that is closest to p
func (r *Ray3) ClosestPoint(p Point3) Point3 {
	// Project point onto ray,
	t := p.Sub(r.Origin).Dot(r.Direction)
	t = max(t, 0) // clamp found point to the ray's origin

	return r.Origin.Add(r.Direction.Mul(t))
}

// Inverse returns a ray with the same origin but pointing in the opposite direction.
func (r *Ray3) Inverse() Ray3 {
	return Ray3{
		Origin:    r.Origin,
		Direction: r.Direction.Mul(-1),
	}
}

func (r *Ray3) ApproxEqual(r2 Ray3) bool {
	return r.Origin.ApproxEqual(r2.Origin) && r.Direction.ApproxEqual(r2.Direction)
}

func (r *Ray3) ApproxEqualThreshold(r2 Ray3, threshold float32) bool {
	return r.Origin.ApproxEqualThreshold(r2.Origin, threshold) && r.Direction.ApproxEqualThreshold(r2.Direction, threshold)
}

// Line3 is 3 dimensional straight line that starts at one point and ends at another.
type Line3 struct {
	Start Point3
	End   Point3
}

// RaycastResult is the result of a raycast test.
type RaycastResult struct {
	Point    Point3
	Normal   Vec3
	Distance float32
	Fail     RaycastFail
}

type RaycastFail int

const (
	RaycastFailUnknown RaycastFail = iota
	RaycastFailOutsideBounds
	RaycastFailTargetBehindRayOrigin
	RaycastFailPlaneFacesAwayFromRay
)

func (r RaycastFail) String() string {
	switch r {
	case RaycastFailOutsideBounds:
		return "outside bounds"
	case RaycastFailTargetBehindRayOrigin:
		return "behind ray origin"
	case RaycastFailPlaneFacesAwayFromRay:
		return "faces away from ray"
	default:
		return "unknown"
	}
}

// Rect is a 2 dimensional axis-aligned rectangle
type Rect struct {
	Position Point2 // Centre of the rectangle
	Size     Vec2   // HALF SIZE!
}

// Min returns the minimum point of the Rect
func (r Rect) Min() Point2 {
	p1 := r.Position.Add(r.Size)
	p2 := r.Position.Sub(r.Size)

	return Point2{
		min(p1[0], p2[0]),
		min(p1[1], p2[1]),
	}
}

// Max returns the maximum point of the Rect
func (r Rect) Max() Point2 {
	p1 := r.Position.Add(r.Size)
	p2 := r.Position.Sub(r.Size)

	return Point2{
		max(p1[0], p2[0]),
		max(p1[1], p2[1]),
	}
}

func (r Rect) TopLeft() Point2 {
	return Vec2{r.Position[0] - r.Size[0], r.Position[1] - r.Size[1]}
}

func (r Rect) TopRight() Point2 {
	return Vec2{r.Position[0] + r.Size[0], r.Position[1] - r.Size[1]}
}

func (r Rect) BottomLeft() Point2 {
	return Vec2{r.Position[0] - r.Size[0], r.Position[1] + r.Size[1]}
}

func (r Rect) BottomRight() Point2 {
	return Vec2{r.Position[0] + r.Size[0], r.Position[1] + r.Size[1]}
}

func (r Rect) Shrink(v float32) Rect {
	return Rect{
		Position: r.Position,
		Size:     Vec2{r.Size[0] - v, r.Size[1] - v},
	}
}

func (r Rect) Width() float32  { return r.Size[0] * 2 }
func (r Rect) Height() float32 { return r.Size[1] * 2 }

// Contains reports whether p is contained within the bounds of the Rect
func (r *Rect) ContainsPoint2(pt Point2) bool {
	min := r.Min()
	max := r.Max()

	return min[0] <= pt[0] && min[1] <= pt[1] &&
		pt[0] <= max[0] && pt[1] <= max[1]
}

func (r Rect) IntersectsRect(r2 Rect) bool {
	rMin := r.Min()
	rMax := r.Max()
	r2Min := r2.Min()
	r2Max := r2.Max()

	return (rMin[0] <= r2Max[0] && rMax[0] >= r2Min[0]) &&
		(rMin[1] <= r2Max[1] && rMax[1] >= r2Min[1])
}

// MTVRect returns the MTV (Minimum Translation Vector) for an overlapping Rect. The MTV is
// the vector that should be applied to r2 to ensure it does not overlap r
func (r Rect) MTVRect(r2 *Rect) (bool, Vec2) {
	rMin := r.Min()
	rMax := r.Max()
	r2Min := r2.Min()
	r2Max := r2.Max()

	// axis 0 interval
	rInterval0 := Interval{Min: rMin[0], Max: rMax[0]}
	r2Interval0 := Interval{Min: r2Min[0], Max: r2Max[0]}

	overlap0 := rInterval0.GetOverlap(r2Interval0)

	// axis 1 interval
	rInterval1 := Interval{Min: rMin[1], Max: rMax[1]}
	r2Interval1 := Interval{Min: r2Min[1], Max: r2Max[1]}

	overlap1 := rInterval1.GetOverlap(r2Interval1)

	// Both axes must overlap
	if overlap0 == 0 || overlap1 == 0 {
		return false, Vec2{}
	}

	if overlap0 < overlap1 {
		if rMin[0] < r2Min[0] {
			return true, Vec2{overlap0, 0}
		}
		return true, Vec2{-overlap0, 0}
	}

	if rMin[1] < r2Min[1] {
		return true, Vec2{0, overlap1}
	}
	return true, Vec2{0, -overlap1}
}

var _ Box3 = (*AABB)(nil)

var (
	aabbAxes    = [3]Vec3{X3, Y3, Z3}
	aabbNormals = [6]Vec3{
		{-1, 0, 0},
		{1, 0, 0},
		{0, -1, 0},
		{0, 1, 0},
		{0, 0, -1},
		{0, 0, 1},
	}
)

// AABB is a 3 dimensional axis-aligned bounding box
type AABB struct {
	Position Point3
	Size     Vec3      // HALF SIZE, i.e. the size in each direction
	corners  [8]Point3 // pre-allocated space to avoid allocations during calls to Corners
}

func AABBFromCorners(pmin, pmax Point3) AABB {
	a := AABB{
		Size: Vec3{
			(pmax[0] - pmin[0]) / 2,
			(pmax[1] - pmin[1]) / 2,
			(pmax[2] - pmin[2]) / 2,
		},
	}

	a.Position[0] = min(pmin[0], pmax[0]) + a.Size[0]
	a.Position[1] = min(pmin[1], pmax[1]) + a.Size[1]
	a.Position[2] = min(pmin[2], pmax[2]) + a.Size[2]
	return a
}

// Min returns the minimum point of the AABB
func (a *AABB) Min() Point3 {
	p1 := a.Position.Add(a.Size)
	p2 := a.Position.Sub(a.Size)

	return Point3{
		min(p1[0], p2[0]),
		min(p1[1], p2[1]),
		min(p1[2], p2[2]),
	}
}

// Max returns the maximum point of the AABB
func (a *AABB) Max() Point3 {
	p1 := a.Position.Add(a.Size)
	p2 := a.Position.Sub(a.Size)

	return Point3{
		max(p1[0], p2[0]),
		max(p1[1], p2[1]),
		max(p1[2], p2[2]),
	}
}

// Corners returns the points at the eight corners of the box.
func (a *AABB) Corners() []Point3 {
	min := a.Min()
	max := a.Max()

	a.corners[0] = Point3{min[0], max[1], max[2]}
	a.corners[0] = Point3{min[0], max[1], min[2]}
	a.corners[0] = Point3{min[0], min[1], max[2]}
	a.corners[0] = Point3{min[0], min[1], min[2]}
	a.corners[0] = Point3{max[0], max[1], max[2]}
	a.corners[0] = Point3{max[0], max[1], min[2]}
	a.corners[0] = Point3{max[0], min[1], max[2]}
	a.corners[0] = Point3{max[0], min[1], min[2]}
	return a.corners[:]
}

func (a *AABB) Axes() []Vec3 {
	return aabbAxes[:]
}

func (a *AABB) Normals() []Vec3 {
	return aabbNormals[:]
}

// Contains reports whether p is contained within the bounds of the AABB
func (a *AABB) ContainsPoint3(pt Point3) bool {
	min := a.Min()
	max := a.Max()

	if pt[0] < min[0] || pt[1] < min[1] || pt[2] < min[2] {
		return false
	}
	if pt[0] > max[0] || pt[1] > max[1] || pt[2] > max[2] {
		return false
	}

	return true
}

// ClosestPoint returns the point in the AABB that is closest to p
func (a *AABB) ClosestPoint(p Point3) Point3 {
	min := a.Min()
	max := a.Max()

	if p[0] < min[0] {
		p[0] = min[0]
	}
	if p[1] < min[1] {
		p[1] = min[1]
	}
	if p[2] < min[2] {
		p[2] = min[2]
	}

	if p[0] > max[0] {
		p[0] = max[0]
	}
	if p[1] > max[1] {
		p[1] = max[1]
	}
	if p[2] > max[2] {
		p[2] = max[2]
	}

	return p
}

func (a *AABB) IntersectsAABB(b *AABB) bool {
	aMin := a.Min()
	aMax := a.Max()
	bMin := b.Min()
	bMax := b.Max()

	return (aMin[0] <= bMax[0] && aMax[0] >= bMin[0]) &&
		(aMin[1] <= bMax[1] && aMax[1] >= bMin[1]) &&
		(aMin[2] <= bMax[2] && aMax[2] >= bMin[2])
}

// MTVAABB returns the MTV (Minimum Translation Vector) for an overlapping AABB
func (a *AABB) MTVAABB(b *AABB) (bool, Vec3) {
	aMin := a.Min()
	aMax := a.Max()
	bMin := b.Min()
	bMax := b.Max()

	if !((aMin[0] <= bMax[0] && aMax[0] >= bMin[0]) &&
		(aMin[1] <= bMax[1] && aMax[1] >= bMin[1]) &&
		(aMin[2] <= bMax[2] && aMax[2] >= bMin[2])) {
		return false, Vec3{}
	}

	var axis Vec3
	var minOverlap float32 = maxFloat32
	var sign float32 = 1

	for i := 0; i < 3; i++ {
		aint := Interval{Min: aMin[i], Max: aMax[i]}
		bint := Interval{Min: bMin[i], Max: bMax[i]}
		overlap := aint.GetOverlap(bint)
		if overlap < minOverlap {
			minOverlap = overlap
			switch i {
			case 0:
				axis = X3
				if a.Position[0] < b.Position[0] {
					sign = -1
				} else {
					sign = 1
				}

			case 1:
				axis = Y3
				if a.Position[1] < b.Position[1] {
					sign = -1
				} else {
					sign = 1
				}
			case 2:
				axis = Z3
				if a.Position[2] < b.Position[2] {
					sign = -1
				} else {
					sign = 1
				}
			}
		}
	}

	if minOverlap <= 0 {
		return false, Vec3{}
	}
	return true, axis.Mul(sign * minOverlap)
}

func (a *AABB) ProjectOntoAxis(axis Vec3) Interval {
	vertex := a.Corners()

	var in Interval
	in.Min = axis.Dot(vertex[0])
	in.Max = in.Min

	for i := 1; i < 8; i++ {
		projection := axis.Dot(vertex[i])
		if projection < in.Min {
			in.Min = projection
		}
		if projection > in.Max {
			in.Max = projection
		}
	}

	return in
}

// Raycast tests whether the ray intersects the AABB
func (a *AABB) Raycast(ray Ray3) (RaycastResult, bool) {
	var res RaycastResult
	amin := a.Min()
	amax := a.Max()

	// debug("aabb min=", amin, "max=", amax)
	// debug("ray origin=", ray.Origin, "direction=", ray.Direction)
	// Any component of direction could be 0!
	// Address this by using a small number, close to
	// 0 in case any of directions components are 0
	t := [6]float32{
		(amin[0] - ray.Origin[0]) / nonzero(ray.Direction[0]),
		(amax[0] - ray.Origin[0]) / nonzero(ray.Direction[0]),
		(amin[1] - ray.Origin[1]) / nonzero(ray.Direction[1]),
		(amax[1] - ray.Origin[1]) / nonzero(ray.Direction[1]),
		(amin[2] - ray.Origin[2]) / nonzero(ray.Direction[2]),
		(amax[2] - ray.Origin[2]) / nonzero(ray.Direction[2]),
	}

	tmin := max(max(min(t[0], t[1]), min(t[2], t[3])), min(t[4], t[5]))
	tmax := min(min(max(t[0], t[1]), max(t[2], t[3])), max(t[4], t[5]))

	// if tmax < 0, ray is intersecting AABB
	// but entire AABB is behind it's origin
	if tmax < 0 {
		// debug("tmax=", tmax, " < 0, entire aabb is behind ray's origin")
		res.Fail = RaycastFailTargetBehindRayOrigin
		return res, false
	}

	// if tmin > tmax, ray doesn't intersect AABB
	if tmin > tmax {
		// debug("tmin > tmax, ray doesn't intersect AABB")
		res.Fail = RaycastFailOutsideBounds
		return res, false
	}

	res.Distance = tmin

	// If tmin is < 0, tmax is closer
	if tmin < 0 {
		res.Distance = tmax
	}

	res.Point = ray.Point(res.Distance)

	// Find closest side to the ray
	normals := [6]Vec3{
		{-1, 0, 0},
		{1, 0, 0},
		{0, -1, 0},
		{0, 1, 0},
		{0, 0, -1},
		{0, 0, 1},
	}

	for i := 0; i < 6; i++ {
		if cmp(res.Distance, t[i]) {
			res.Normal = normals[i]
		}
	}

	return res, true
}

func (a *AABB) OBB(tx *Transform) OBB {
	o := OBB{
		Position:    tx.Pos(),
		Size:        a.Size,
		Orientation: tx.Orientation(),
	}

	scale := tx.Scale()
	o.Size[0] *= scale[0]
	o.Size[1] *= scale[1]
	o.Size[2] *= scale[2]

	return o
}

// Plane3 is a plane in 3 dimensions
type Plane3 struct {
	Normal   Vec3    // Must be normalized
	Distance float32 // distance from origin
}

// Raycast tests whether the ray intersects the Plane.
// See https://www.cs.princeton.edu/courses/archive/fall00/cs426/lectures/raycast/sld017.htm
func (p *Plane3) Raycast(ray Ray3) (RaycastResult, bool) {
	var res RaycastResult

	nd := ray.Direction.Dot(p.Normal)
	pn := ray.Origin.Dot(p.Normal)

	// if nd is positive, the ray and plane normals
	// point in the same direction. No intersection.
	if nd >= 0 {
		res.Fail = RaycastFailPlaneFacesAwayFromRay
		return res, false
	}

	t := -(p.Distance + pn) / nd

	// t must be positive
	if t >= 0.0 {
		res.Distance = t
		res.Point = ray.Origin.Add(ray.Direction.Mul(t))
		res.Normal = p.Normal.Normalize() // TODO: isn't this the ray direction?
		return res, true
	}

	res.Fail = RaycastFailTargetBehindRayOrigin
	return res, false
}

// ClosestPoint returns the point in the plane that is closest to point
func (p *Plane3) ClosestPoint(point Point3) Point3 {
	// This works assuming plane.Normal is normalized, which it should be
	distance := p.Normal.Dot(point) - p.Distance
	return point.Sub(p.Normal.Mul(distance))
}

// ContainsPoint3 reports whether the point lies on the plane.
func (p *Plane3) ContainsPoint3(point Point3) bool {
	return cmp(point.Dot(p.Normal)-p.Distance, 0)
}

// Add performs element-wise addition between two vectors.
func (v1 Vec2i) Add(v2 Vec2i) Vec2i {
	return Vec2i{v1[0] + v2[0], v1[1] + v2[1]}
}

// Sub performs element-wise subtraction between two vectors.
func (v1 Vec2i) Sub(v2 Vec2i) Vec2i {
	return Vec2i{v1[0] - v2[0], v1[1] - v2[1]}
}

// Mul performs a scalar multiplication between the vector and some constant value
func (v1 Vec2i) Mul(c float32) Vec2i {
	return Vec2i{int32(float32(v1[0]) * c), int32(float32(v1[1]) * c)}
}

// Mul2 performs an element-wise scalar multiplication between the vector and another vector
func (v1 Vec2i) Mul2(v Vec2) Vec2i {
	return Vec2i{int32(float32(v1[0]) * v[0]), int32(float32(v1[1]) * v[1])}
}

// Add performs element-wise addition between two vectors.
func (v1 Vec3i) Add(v2 Vec3i) Vec3i {
	return Vec3i{v1[0] + v2[0], v1[1] + v2[1], v1[2] + v2[2]}
}

// Sub performs element-wise subtraction between two vectors.
func (v1 Vec3i) Sub(v2 Vec3i) Vec3i {
	return Vec3i{v1[0] - v2[0], v1[1] - v2[1], v1[2] - v2[2]}
}

type Sphere struct {
	Position Point3
	Radius   float32
}

// ClosestPoint returns the point on the sphere that is closest to point
func (s *Sphere) ClosestPoint(point Point3) Point3 {
	sphereToPoint := point.Sub(s.Position).Normalize()
	sphereToPoint.Mul(s.Radius)
	return sphereToPoint.Mul(s.Radius).Add(s.Position)
}

// ContainsPoint3 reports whether the point lies in the sphere.
func (s *Sphere) ContainsPoint3(point Point3) bool {
	e := point.Sub(s.Position)
	eMagnitudeSquared := e.Dot(e)
	rSquared := s.Radius * s.Radius

	return eMagnitudeSquared < rSquared
}

// Raycast tests whether the ray intersects the Sphere.
func (s *Sphere) Raycast(ray Ray3) (RaycastResult, bool) {
	var res RaycastResult

	e := s.Position.Sub(ray.Origin)
	rSquared := s.Radius * s.Radius

	eMagnitudeSquared := e.Dot(e)
	a := e.Dot(ray.Direction)

	bSquared := eMagnitudeSquared - (a * a)
	f := sqrt(abs(rSquared - bSquared))

	// Assume normal intersection
	t := a - f

	if rSquared-bSquared < 0 {
		// No collision has happened

		res.Fail = RaycastFailOutsideBounds
		return res, false
	} else if eMagnitudeSquared < rSquared {
		// Ray starts inside the sphere
		// Reverse direction
		t = a + f
	}

	res.Distance = t
	res.Point = ray.Origin.Add(ray.Direction.Mul(t))
	res.Normal = res.Point.Sub(s.Position).Normalize()
	return res, true
}

// Rect is a 2 dimensional axis-aligned rectangle
type Recti struct {
	Position Point2i // Centre of the rectangle
	Size     Vec2i   // half the width and height
}

func (r Recti) TopLeft() Point2i {
	return Point2i{r.Position[0] - r.Size[0], r.Position[1] - r.Size[1]}
}

func (r Recti) TopRight() Point2i {
	return Point2i{r.Position[0] + r.Size[0], r.Position[1] - r.Size[1]}
}

func (r Recti) BottomLeft() Point2i {
	return Point2i{r.Position[0] - r.Size[0], r.Position[1] + r.Size[1]}
}

func (r Recti) BottomRight() Point2i {
	return Point2i{r.Position[0] + r.Size[0], r.Position[1] + r.Size[1]}
}

func (r Recti) Shrink(v int32) Recti {
	return Recti{
		Position: r.Position,
		Size:     Vec2i{r.Size[0] - v, r.Size[1] - v},
	}
}

func (r Recti) Width() int32  { return r.Size[0] * 2 }
func (r Recti) Height() int32 { return r.Size[1] * 2 }

// Contains reports whether p is contained within the bounds of the Rect
func (r Recti) ContainsPoint2i(pt Point2i) bool {
	min := r.Min()
	max := r.Max()

	return min[0] <= pt[0] && min[1] <= pt[1] &&
		pt[0] <= max[0] && pt[1] <= max[1]
}

// Min returns the minimum point of the Rect
func (r Recti) Min() Point2i {
	p1 := r.Position.Add(r.Size)
	p2 := r.Position.Sub(r.Size)

	return Point2i{
		mini(p1[0], p2[0]),
		mini(p1[1], p2[1]),
	}
}

// Max returns the maximum point of the Rect
func (r Recti) Max() Point2i {
	p1 := r.Position.Add(r.Size)
	p2 := r.Position.Sub(r.Size)

	return Point2i{
		maxi(p1[0], p2[0]),
		maxi(p1[1], p2[1]),
	}
}

func mini(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func maxi(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func RectiFromCorners(tl, br Point2i) Recti {
	size := Point2i{(br[0] - tl[0]) / 2, (br[1] - tl[1]) / 2}
	return Recti{
		Position: Point2i{tl[0] + size[0], tl[1] + size[1]},
		Size:     size,
	}
}

// Tri3 is a triangle whose corners are 3 points in 3 dimensions. A,B and C
// are assume to  be in counter clockwise order.
type Tri3 struct {
	A, B, C Point3
}

// The Centroid of a triangle is the intersection of the three medians of the triangle
func (t Tri3) Centroid() Vec3 {
	var result Vec3
	result[0] = (t.A[0] + t.B[0] + t.C[0]) / 3
	result[1] = (t.A[1] + t.B[1] + t.C[1]) / 3
	result[2] = (t.A[2] + t.B[2] + t.C[2]) / 3
	return result
}

func (t Tri3) ContainsPoint3(pt Point3) bool {
	// Move the triangle so that the point is
	// now at the origin of the triangle
	a := t.A.Sub(pt)
	b := t.B.Sub(pt)
	c := t.C.Sub(pt)

	// The point should be moved too, so they are both
	// relative, but because we don't use p in the
	// equation anymore, we don't need it!
	// p -= p; // This would just equal the zero vector!

	normPBC := b.Cross(c) // Normal of PBC (u)
	normPCA := c.Cross(a) // Normal of PCA (v)
	normPAB := a.Cross(b) // Normal of PAB (w)

	// Test to see if the normals are facing
	// the same direction, return false if not
	if normPBC.Dot(normPCA) < 0.0 {
		return false
	} else if normPBC.Dot(normPAB) < 0.0 {
		return false
	}

	// All normals facing the same way, return true
	return true
}

// BarycentricPoint3 returns the barycentric coordinates of pt which must be within the triangle.
func (t Tri3) BarycentricPoint3(pt Point3) Vec3 {
	v0 := t.B.Sub(t.A)
	v1 := t.C.Sub(t.A)
	v2 := pt.Sub(t.A)

	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	denom := d00*d11 - d01*d01

	if cmp(denom, 0.0) {
		return Vec3{}
	}

	var result Vec3
	result[1] = (d11*d20 - d01*d21) / denom
	result[2] = (d00*d21 - d01*d20) / denom
	result[0] = 1.0 - result[1] - result[2]
	return result
}

func (t *Tri3) SortCCW(normal Vec3) {
	// See https://stackoverflow.com/a/14371081/325180
	c := t.Centroid()
	angle := normal.Dot(t.A.Sub(c).Cross(t.B.Sub(c)))
	// if angle is positive then B is counterclockwise from A
	if angle < 0 {
		// Swap them
		t.A, t.B = t.B, t.A
	}

	angle = normal.Dot(t.B.Sub(c).Cross(t.C.Sub(c)))
	if angle < 0 {
		// Swap them
		t.B, t.C = t.C, t.B
	}
}

// TODO: ensure edges are sorted CCW?
func (t *Tri3) Edges() []Line3 {
	return []Line3{
		{Start: t.A, End: t.B},
		{Start: t.B, End: t.C},
		{Start: t.C, End: t.A},
	}
}

// Plane3FromTri3 returns the plane that lies on the triangle
func Plane3FromTri3(t Tri3) Plane3 {
	var result Plane3
	result.Normal = t.B.Sub(t.A).Cross(t.C.Sub(t.A)).Normalize()
	result.Distance = result.Normal.Dot(t.A)
	return result
}

// Tri2 is a triangle whose corners are 3 points in 2 dimensions. A, B and C
// are assume to  be in counter clockwise order.
type Tri2 struct {
	A, B, C Point2
}

// The Centroid of a triangle is the intersection of the three medians of the triangle
func (t Tri2) Centroid() Vec2 {
	var result Vec2
	result[0] = (t.A[0] + t.B[0] + t.C[0]) / 3
	result[1] = (t.A[1] + t.B[1] + t.C[1]) / 3
	return result
}

func (t Tri2) ContainsPoint2(pt Point2) bool {
	b := t.BarycentricPoint2(pt)

	// Point is inside triangle if all barycentric coordinates are in range [0,1]
	// Can be inaccurate due to float precision.
	// See http://totologic.blogspot.co.uk/2014/01/accurate-point-in-triangle-test.html
	return 0 <= b[0] && b[0] <= 1 &&
		0 <= b[1] && b[1] <= 1 &&
		0 <= b[2] && b[2] <= 1
}

// BarycentricPoint2 returns the barycentric coordinates of pt which must be within the triangle.
func (t Tri2) BarycentricPoint2(pt Point2) Vec3 {
	v0 := t.C.Sub(t.A)
	v1 := t.B.Sub(t.A)
	v2 := pt.Sub(t.A)

	dot00 := v0.Dot(v0)
	dot01 := v0.Dot(v1)
	dot02 := v0.Dot(v2)
	dot11 := v1.Dot(v1)
	dot12 := v1.Dot(v2)

	denom := (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) / denom
	v := (dot00*dot12 - dot01*dot02) / denom

	return Vec3{
		1 - u - v,
		v,
		u,
	}
}

// CircumCircle returns the circle that circumscribes the triangle
func (t Tri2) CircumCircle() Circle {
	var c Circle

	x1, y1 := t.A[0], t.A[1]
	x2, y2 := t.B[0], t.B[1]
	x3, y3 := t.C[0], t.C[1]

	var m1, m2, mx1, mx2, my1, my2 float32

	fabsy1y2 := abs(y1 - y2)
	fabsy2y3 := abs(y2 - y3)

	// Check for coincident points
	if fabsy1y2 < epsilon32 && fabsy2y3 < epsilon32 {
		return c
	}

	if fabsy1y2 < epsilon32 {
		m2 = -(x3 - x2) / (y3 - y2)
		mx2 = (x2 + x3) / 2.0
		my2 = (y2 + y3) / 2.0
		c.Centre[0] = (x2 + x1) / 2.0
		c.Centre[1] = m2*(c.Centre[0]-mx2) + my2
	} else if fabsy2y3 < epsilon32 {
		m1 = -(x2 - x1) / (y2 - y1)
		mx1 = (x1 + x2) / 2.0
		my1 = (y1 + y2) / 2.0
		c.Centre[0] = (x3 + x2) / 2.0
		c.Centre[1] = m1*(c.Centre[0]-mx1) + my1
	} else {
		m1 = -(x2 - x1) / (y2 - y1)
		m2 = -(x3 - x2) / (y3 - y2)
		mx1 = (x1 + x2) / 2.0
		mx2 = (x2 + x3) / 2.0
		my1 = (y1 + y2) / 2.0
		my2 = (y2 + y3) / 2.0
		c.Centre[0] = (m1*mx1 - m2*mx2 + my2 - my1) / (m1 - m2)
		if fabsy1y2 > fabsy2y3 {
			c.Centre[1] = m1*(c.Centre[0]-mx1) + my1
		} else {
			c.Centre[1] = m2*(c.Centre[0]-mx2) + my2
		}
	}

	dx := x2 - c.Centre[0]
	dy := y2 - c.Centre[1]
	c.Radius = dx*dx + dy*dy

	return c
}

type Circle struct {
	Centre Point2
	Radius float32
}

func (c Circle) ContainsPoint2(pt Point2) bool {
	dx := pt[0] - c.Centre[0]
	dy := pt[1] - c.Centre[1]
	distance := dx*dx + dy*dy

	return (distance - c.Radius) <= epsilon32
}

func DistanceSquared3(a, b Vec3) float32 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]

	return dx*dx + dy*dy + dz*dz
}

var _ Box3 = (*OBB)(nil)

// An OBB is an oriented bounding box
type OBB struct {
	Position    Point3
	Size        Vec3 // HALF SIZE, i.e. the size in each direction
	Orientation mgl32.Quat
	axes        [3]Vec3   // pre-allocated space to avoid allocations during calls to Axes
	corners     [8]Point3 // pre-allocated space to avoid allocations during calls to Corners
}

// ContainsPoint3 reports whether the point lies within the OBB.
func (o *OBB) ContainsPoint3(pt Point3) bool {
	if o.Orientation == mgl32.QuatIdent() {
		return (&AABB{Position: o.Position, Size: o.Size}).ContainsPoint3(pt)
	}

	dir := pt.Sub(o.Position)

	axes := o.Axes()
	for i := 0; i < 3; i++ {
		distance := dir.Dot(axes[i])
		if distance > o.Size[i] {
			return false
		}
		if distance < -o.Size[i] {
			return false
		}
	}

	return true
}

// Corners returns the points at the eight corners of the box.
func (o *OBB) Corners() []Point3 {
	if o.Orientation == mgl32.QuatIdent() {
		return (&AABB{Position: o.Position, Size: o.Size}).Corners()
	}
	o.corners[0] = o.Orientation.Rotate(Vec3{o.Position[0] + o.Size[0], o.Position[1] + o.Size[1], o.Position[2] + o.Size[2]})
	o.corners[1] = o.Orientation.Rotate(Vec3{o.Position[0] + o.Size[0], o.Position[1] + o.Size[1], o.Position[2] - o.Size[2]})
	o.corners[2] = o.Orientation.Rotate(Vec3{o.Position[0] + o.Size[0], o.Position[1] - o.Size[1], o.Position[2] + o.Size[2]})
	o.corners[3] = o.Orientation.Rotate(Vec3{o.Position[0] + o.Size[0], o.Position[1] - o.Size[1], o.Position[2] - o.Size[2]})
	o.corners[4] = o.Orientation.Rotate(Vec3{o.Position[0] - o.Size[0], o.Position[1] + o.Size[1], o.Position[2] + o.Size[2]})
	o.corners[5] = o.Orientation.Rotate(Vec3{o.Position[0] - o.Size[0], o.Position[1] + o.Size[1], o.Position[2] - o.Size[2]})
	o.corners[6] = o.Orientation.Rotate(Vec3{o.Position[0] - o.Size[0], o.Position[1] - o.Size[1], o.Position[2] + o.Size[2]})
	o.corners[7] = o.Orientation.Rotate(Vec3{o.Position[0] - o.Size[0], o.Position[1] - o.Size[1], o.Position[2] - o.Size[2]})
	return o.corners[:]
}

func (o *OBB) Axes() []Vec3 {
	if o.Orientation == mgl32.QuatIdent() {
		return (&AABB{Position: o.Position, Size: o.Size}).Axes()
	}
	o.axes[0] = o.Orientation.Rotate(X3).Normalize()
	o.axes[1] = o.Orientation.Rotate(Y3).Normalize()
	o.axes[2] = o.Orientation.Rotate(Z3).Normalize()

	return o.axes[:]
}

func (o *OBB) Normals() []Vec3 {
	if o.Orientation == mgl32.QuatIdent() {
		return (&AABB{Position: o.Position, Size: o.Size}).Normals()
	}
	return []Vec3{
		o.Orientation.Rotate(Vec3{-1, 0, 0}).Normalize(),
		o.Orientation.Rotate(Vec3{1, 0, 0}).Normalize(),
		o.Orientation.Rotate(Vec3{0, -1, 0}).Normalize(),
		o.Orientation.Rotate(Vec3{0, 1, 0}).Normalize(),
		o.Orientation.Rotate(Vec3{0, 0, -1}).Normalize(),
		o.Orientation.Rotate(Vec3{0, 0, 1}).Normalize(),
	}
}

func (o *OBB) ProjectOntoAxis(axis Vec3) Interval {
	if o.Orientation == mgl32.QuatIdent() {
		return (&AABB{Position: o.Position, Size: o.Size}).ProjectOntoAxis(axis)
	}
	vertex := o.Corners()

	var in Interval
	in.Min = axis.Dot(vertex[0])
	in.Max = in.Min

	for i := 1; i < 8; i++ {
		projection := axis.Dot(vertex[i])
		if projection < in.Min {
			in.Min = projection
		}
		if projection > in.Max {
			in.Max = projection
		}
	}

	return in
}

func (o *OBB) Raycast(ray Ray3) (RaycastResult, bool) {
	var res RaycastResult

	axes := o.Axes()
	f := [3]float32{
		axes[0].Dot(ray.Direction),
		axes[1].Dot(ray.Direction),
		axes[2].Dot(ray.Direction),
	}

	p := o.Position.Sub(ray.Origin)
	e := [3]float32{
		axes[0].Dot(p),
		axes[1].Dot(p),
		axes[2].Dot(p),
	}

	if cmp(f[0], 0) {
		if -e[0]-o.Size[0] > 0 || -e[0]+o.Size[0] < 0 {
			res.Fail = RaycastFailOutsideBounds
			return res, false
		}
		f[0] = nonzero(f[0]) // Avoid div by 0!
	} else if cmp(f[1], 0) {
		if -e[1]-o.Size[1] > 0 || -e[1]+o.Size[1] < 0 {
			res.Fail = RaycastFailOutsideBounds
			return res, false
		}
		f[1] = nonzero(f[1]) // Avoid div by 0!
	} else if cmp(f[2], 0) {
		if -e[2]-o.Size[2] > 0 || -e[2]+o.Size[2] < 0 {
			res.Fail = RaycastFailOutsideBounds
			return res, false
		}
		f[2] = nonzero(f[2]) // Avoid div by 0!
	}

	t := [6]float32{
		(e[0] + o.Size[0]) / f[0],
		(e[0] - o.Size[0]) / f[0],
		(e[1] + o.Size[1]) / f[1],
		(e[1] - o.Size[1]) / f[1],
		(e[2] + o.Size[2]) / f[2],
		(e[2] - o.Size[2]) / f[2],
	}

	tmin := max(max(min(t[0], t[1]), min(t[2], t[3])), min(t[4], t[5]))
	tmax := min(min(max(t[0], t[1]), max(t[2], t[3])), max(t[4], t[5]))

	// if tmax < 0, ray is intersecting AABB
	// but entire AABB is behing it's origin
	if tmax < 0 {
		res.Fail = RaycastFailTargetBehindRayOrigin
		return res, false
	}

	// if tmin > tmax, ray doesn't intersect AABB
	if tmin > tmax {
		res.Fail = RaycastFailOutsideBounds
		return res, false
	}

	// If tmin is < 0, tmax is closer
	res.Distance = tmin

	// If tmin is < 0, tmax is closer
	if tmin < 0 {
		res.Distance = tmax
	}

	res.Point = ray.Point(res.Distance)

	// Find closest side to the ray
	normals := [6]Vec3{
		axes[0],         // +x
		axes[0].Mul(-1), // -x
		axes[1],         // +y
		axes[1].Mul(-1), // -y
		axes[2],         // +z
		axes[2].Mul(-1), // -z
	}

	for i := 0; i < 6; i++ {
		if cmp(res.Distance, t[i]) {
			res.Normal = normals[i].Normalize()
		}
	}
	return res, true
}
