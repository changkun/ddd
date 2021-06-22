// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package geometry

import (
	_ "image/jpeg" // for jpg encoding

	"changkun.de/x/ddd/geometry/primitive"
	"changkun.de/x/ddd/material"
	"changkun.de/x/ddd/object"
)

type Mesh interface {
	object.Object

	AABB() primitive.AABB
	Normalize()
	SetMaterial(m material.Material)
	GetMaterial() material.Material
	NumTriangles() uint64
	Faces(func(f primitive.Face, m material.Material) bool)
}
