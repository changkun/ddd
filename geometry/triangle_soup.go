// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package geometry

import (
	"changkun.de/x/ddd/geometry/primitive"
	"changkun.de/x/ddd/material"
	"changkun.de/x/ddd/math"
	"changkun.de/x/ddd/object"
)

var (
	_ Mesh = &TriangleSoup{}
)

// TriangleSoup implements a triangular mesh.
type TriangleSoup struct {
	// all faces of the triangle mesh.
	faces []*primitive.Triangle
	// the corresponding material of the triangle mesh.
	material material.Material
	// aabb must be transformed when applying the context.
	aabb *primitive.AABB

	math.TransformContext
}

func (f *TriangleSoup) Type() object.Type {
	return object.TypeMesh
}

func (f *TriangleSoup) NumTriangles() uint64 {
	return uint64(len(f.faces))
}

func (f *TriangleSoup) Faces(iter func(primitive.Face, material.Material) bool) {
	for i := range f.faces {
		if !iter(f.faces[i], f.material) {
			return
		}
	}
}

func (f *TriangleSoup) GetMaterial() material.Material {
	return f.material
}

func (t *TriangleSoup) SetMaterial(mat material.Material) {
	t.material = mat
}

// NewTriangleSoup returns a triangular soup.
func NewTriangleSoup(ts []*primitive.Triangle) *TriangleSoup {
	// Compute AABB at loading time.
	aabb := ts[0].AABB()
	for i := 1; i < len(ts); i++ {
		aabb.Add(ts[i].AABB())
	}

	ret := &TriangleSoup{
		faces: ts,
		aabb:  &aabb,
	}
	ret.ResetContext()
	return ret
}

func (m *TriangleSoup) AABB() primitive.AABB {
	if m.aabb == nil {
		// Compute AABB if not computed
		aabb := m.faces[0].AABB()
		lenth := len(m.faces)
		for i := 1; i < lenth; i++ {
			aabb.Add(m.faces[i].AABB())
		}
		m.aabb = &aabb
	}

	min := m.aabb.Min.Apply(m.ModelMatrix())
	max := m.aabb.Max.Apply(m.ModelMatrix())
	return primitive.AABB{Min: min, Max: max}
}

func (m *TriangleSoup) Center() math.Vector {
	aabb := m.AABB()
	return aabb.Min.Add(aabb.Max).Pos()
}

func (m *TriangleSoup) Radius() float64 {
	aabb := m.AABB()
	return aabb.Max.Sub(aabb.Min).Len() / 2
}

// Normalize rescales the mesh to the unit sphere centered at the origin.
func (m *TriangleSoup) Normalize() {
	aabb := m.AABB()
	center := aabb.Min.Add(aabb.Max).Pos()
	radius := aabb.Max.Sub(aabb.Min).Len() / 2
	fac := 1 / radius

	// scale all vertices
	for i := 0; i < len(m.faces); i++ {
		f := m.faces[i]
		f.V1.Pos = f.V1.Pos.Apply(m.ModelMatrix()).Translate(-center.X, -center.Y, -center.Z).Scale(fac, fac, fac, 1)
		f.V2.Pos = f.V2.Pos.Apply(m.ModelMatrix()).Translate(-center.X, -center.Y, -center.Z).Scale(fac, fac, fac, 1)
		f.V3.Pos = f.V3.Pos.Apply(m.ModelMatrix()).Translate(-center.X, -center.Y, -center.Z).Scale(fac, fac, fac, 1)
	}

	// update AABB after scaling
	min := aabb.Min.Translate(-center.X, -center.Y, -center.Z).Scale(fac, fac, fac, 1)
	max := aabb.Max.Translate(-center.X, -center.Y, -center.Z).Scale(fac, fac, fac, 1)
	m.aabb = &primitive.AABB{Min: min, Max: max}
	m.ResetContext()
}