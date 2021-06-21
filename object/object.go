// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package object

import "changkun.de/x/ddd/math"

type Type int

const (
	TypeGroup = iota
	TypeMesh
	TypeCamera
	TypeLight
)

type Object interface {
	Type() Type
	Rotate(r math.Vector, a float64)
	RotateX(a float64)
	RotateY(a float64)
	RotateZ(a float64)
	Translate(x, y, z float64)
	Scale(x, y, z float64)
	ModelMatrix() math.Matrix
}