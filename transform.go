package geom

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	position    Vec3
	scale       Vec3
	orientation Quat
	matrix      *Mat4
}

func NewTransform() Transform {
	return Transform{
		position:    Vec3{0, 0, 0},
		scale:       Vec3{1, 1, 1},
		orientation: mgl32.QuatIdent(),
	}
}

// SetPosition sets the position of the object.
func (t *Transform) SetPosition(v Vec3) {
	t.position = clampZeroVec3(v)
	t.matrix = nil
}

func (t *Transform) SetScale(v Vec3) {
	t.scale = clampZeroVec3(v)
	t.matrix = nil
}

func (t *Transform) SetOrientation(q Quat) {
	t.orientation = q.Normalize()
	t.matrix = nil
}

func (t *Transform) Matrix() Mat4 {
	if t.matrix == nil {
		trans := mgl32.Translate3D(t.position[0], t.position[1], t.position[2])
		scale := mgl32.Scale3D(t.scale[0], t.scale[1], t.scale[2])
		rot := t.orientation.Mat4()

		m := trans.Mul4(rot).Mul4(scale)
		t.matrix = &m
	}
	return *t.matrix
}

func (t *Transform) SetMatrix(m Mat4) {
	t.scale[0], t.scale[1], t.scale[2] = mgl32.Extract3DScale(m)
	t.position[0], t.position[1], t.position[2] = m[12], m[13], m[14]

	rot := Mat4{
		// col 0
		m[0] / t.scale[0],
		m[1] / t.scale[0],
		m[2] / t.scale[0],
		0,
		// col 1
		m[4] / t.scale[1],
		m[5] / t.scale[1],
		m[6] / t.scale[1],
		0,
		// col 2
		m[8] / t.scale[2],
		m[9] / t.scale[2],
		m[10] / t.scale[2],
		0,
		// col 3
		0,
		0,
		0,
		1,
	}

	t.orientation = mgl32.Mat4ToQuat(rot)
	t.matrix = nil
}

// Scale returns the scale of the object
func (t *Transform) Scale() Vec3 {
	return t.scale
}

// SetScaleUniform sets the scale of the object to v along all axes
func (t *Transform) SetScaleUniform(v float32) {
	t.SetScale(Vec3{v, v, v})
}

// ScaleBy changes the scale of the object by x, y and z along three axes.
func (t *Transform) ScaleBy(x, y, z float32) {
	t.SetScale(Vec3{t.scale[0] * x, t.scale[1] * y, t.scale[2] * z})
}

// ScaleUniformBy changes the scale of the object by v along all axes.
func (t *Transform) ScaleUniformBy(v float32) {
	t.SetScale(t.scale.Mul(v))
}

// Pos returns the position of the object.
func (t *Transform) Pos() Vec3 {
	return t.position
}

// Translate translates the object by v.
func (t *Transform) Translate(v Vec3) {
	t.SetPosition(t.position.Add(v))
}

// Translate translates the object by v along axis.
func (t *Transform) TranslateAlong(v float32, axis Vec3) {
	t.Translate(axis.Mul(v))
}

// Orientation returns the orientation of the object.
func (t *Transform) Orientation() Quat {
	return t.orientation
}

// SetAngleAbout sets the orientation of the object to the angle in radians about axis.
func (t *Transform) SetAngleAbout(axis Vec3, angle float32) {
	t.SetOrientation(mgl32.QuatRotate(angle, axis))
}

// RotateAbout rotates the object about axis by the angle given in radians.
func (t *Transform) RotateAbout(axis Vec3, angle float32) {
	t.Rotate(mgl32.QuatRotate(angle, axis))
}

// Rotate rotates the object using the rotation defined by q.
func (t *Transform) Rotate(q Quat) {
	t.SetOrientation(q.Normalize().Mul(t.orientation))
}

// RotateToward rotates the object so that it faces the target. Its Front vector
// will point toward target.
func (t *Transform) RotateToward(target Vec3) {
	desiredFront := target.Sub(t.position).Normalize()
	rotation := mgl32.QuatBetweenVectors(t.Front(), desiredFront).Normalize()
	t.Rotate(rotation)
}

// Front returns the direction the front of the object is facing. The vector will point along
// the object's local Z axis.
func (t *Transform) Front() Vec3 {
	return clampZeroVec3(t.orientation.Rotate(Vec3{0, 0, 1}).Normalize())
}

// Top returns the direction the top of the object is facing. The vector will point along
// the object's local Y axis.
func (t *Transform) Top() Vec3 {
	return clampZeroVec3(t.orientation.Rotate(Vec3{0, 1, 0}).Normalize())
}

// // Right returns the direction the right of the object is facing. The vector will point along
// // the object's local negative X axis.
// func (t *Transform) Right() Vec3 {
// 	return ClampZeroVec3(t.orientation.Rotate(Vec3{-1, 0, 0}).Normalize())
// }

// Left returns the direction the left of the object is facing. The vector will point along
// the object's local X axis.
func (t *Transform) Left() Vec3 {
	return clampZeroVec3(t.orientation.Rotate(Vec3{1, 0, 0}).Normalize())
}

// Advance moves the object along the direction it is facing without rotating.
func (t *Transform) Advance(s float32) {
	t.Translate(t.Front().Mul(s))
}

// Ascend moves the object along the direction of its top vector without rotating.
func (t *Transform) Ascend(s float32) {
	t.Translate(t.Top().Mul(s))
}

// Strafe moves the object along its left pointing vector without rotating.
func (t *Transform) Strafe(s float32) {
	t.Translate(t.Left().Mul(s))
}

// Pitch rotates the object about its left pointing vector.
func (t *Transform) Pitch(angle float32) {
	t.RotateAbout(t.Left(), angle)
}

// Yaw rotates the object about the direction of its top vector.
func (t *Transform) Yaw(angle float32) {
	t.RotateAbout(t.Top(), angle)
}

// Roll rotates the object about about the direction of its front vector.
// does not change, nor does the point it is looking at.
func (t *Transform) Roll(angle float32) {
	t.RotateAbout(t.Front(), angle)
}
