// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package geometry_test

import (
	"testing"

	"poly.red/geometry"
	"poly.red/geometry/primitive"
	"poly.red/material"
)

func TestBufferedMesh(t *testing.T) {

	bm := geometry.NewBufferedMesh()
	bm.SetAttribute(geometry.AttributePos, &geometry.BufferAttribute{
		Stride: 3,
		Values: []float64{
			-0.363322, -0.387725, 0.85933, // 0
			-0.55029, -0.387725, -0.682297, // 1
			-0.038214, 0.990508, -0.126177, // 2
			0.951827, -0.215059, -0.050857, // 3
		},
	})
	bm.SetVertexIndex([]uint64{
		2, 3, 1,
		2, 0, 3,
		3, 0, 1,
		1, 0, 2,
	})

	bm.AABB()
	bm.Normalize()

	counter := 0
	bm.Faces(func(f primitive.Face, m material.Material) bool {
		counter++
		return true
	})

	if counter != 4 {
		t.Fatalf("expect 4 faces, but only got %v", counter)
	}
}
