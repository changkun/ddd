// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

// +build darwin

package mtl_test

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"poly.red/internal/driver/mtl"
)

func TestRenderTriangle(t *testing.T) {
	device, err := mtl.CreateSystemDefaultDevice()
	if err != nil {
		t.Log(err)
		return
	}

	// Create a render pipeline state.
	const source = `#include <metal_stdlib>

using namespace metal;

struct Vertex {
	float4 position [[position]];
	float4 color;
};

vertex Vertex VertexShader(
	uint vertexID [[vertex_id]],
	device Vertex * vertices [[buffer(0)]]
) {
	return vertices[vertexID];
}

fragment float4 FragmentShader(Vertex in [[stage_in]]) {
	return in.color;
}
`
	lib, err := device.MakeLibrary(source, mtl.CompileOptions{})
	if err != nil {
		t.Fatal(err)
	}
	vs, err := lib.MakeFunction("VertexShader")
	if err != nil {
		t.Fatal(err)
	}
	fs, err := lib.MakeFunction("FragmentShader")
	if err != nil {
		t.Fatal(err)
	}
	var rpld mtl.RenderPipelineDescriptor
	rpld.VertexFunction = vs
	rpld.FragmentFunction = fs
	rpld.ColorAttachments[0].PixelFormat = mtl.PixelFormatRGBA8UNorm
	rps, err := device.MakeRenderPipelineState(rpld)
	if err != nil {
		t.Fatal(err)
	}

	// Create a vertex buffer.
	type Vertex struct {
		Position [4]float32
		Color    [4]float32
	}
	vertexData := [...]Vertex{
		{[4]float32{+0.00, +0.75, 0, 1}, [4]float32{1, 0, 0, 1}},
		{[4]float32{-0.75, -0.75, 0, 1}, [4]float32{0, 1, 0, 1}},
		{[4]float32{+0.75, -0.75, 0, 1}, [4]float32{0, 0, 1, 1}},
	}
	vertexBuffer := device.MakeBuffer(unsafe.Pointer(&vertexData[0]), unsafe.Sizeof(vertexData), mtl.ResourceStorageModeManaged)

	// Create an output texture to render into.
	td := mtl.TextureDescriptor{
		PixelFormat: mtl.PixelFormatRGBA8UNorm,
		Width:       512,
		Height:      512,
		StorageMode: mtl.StorageModeManaged,
	}
	texture := device.MakeTexture(td)

	cq := device.MakeCommandQueue()
	cb := cq.MakeCommandBuffer()

	// Encode all render commands.
	var rpd mtl.RenderPassDescriptor
	rpd.ColorAttachments[0].LoadAction = mtl.LoadActionClear
	rpd.ColorAttachments[0].StoreAction = mtl.StoreActionStore
	rpd.ColorAttachments[0].ClearColor = mtl.ClearColor{Red: 0.35, Green: 0.65, Blue: 0.85, Alpha: 1}
	rpd.ColorAttachments[0].Texture = texture
	rce := cb.MakeRenderCommandEncoder(rpd)
	rce.SetRenderPipelineState(rps)
	rce.SetVertexBuffer(vertexBuffer, 0, 0)
	rce.DrawPrimitives(mtl.PrimitiveTypeTriangle, 0, 3)
	rce.EndEncoding()

	// Encode all blit commands.
	bce := cb.MakeBlitCommandEncoder()
	bce.Synchronize(texture)
	bce.EndEncoding()

	cb.Commit()
	cb.WaitUntilCompleted()

	// Read pixels from output texture into an image.
	got := image.NewNRGBA(image.Rect(0, 0, texture.Width(), texture.Height()))
	bytesPerRow := 4 * texture.Width()
	region := mtl.RegionMake2D(0, 0, texture.Width(), texture.Height())
	texture.GetBytes(&got.Pix[0], uintptr(bytesPerRow), region, 0)

	want, err := readPNG(filepath.Join("testdata", "triangle.png"))
	if err != nil {
		t.Fatal(err)
	}

	if err := imageEq(got, want); err != nil {
		t.Errorf("got image != want: %v", err)
	}
}

func readPNG(name string) (image.Image, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// imageEq reports whether images m, n are considered equivalent. Two images are considered
// equivalent if they have same bounds, and all pixel colors are within a small margin.
func imageEq(m, n image.Image) error {
	if m.Bounds() != n.Bounds() {
		return fmt.Errorf("bounds don't match: %v != %v", m.Bounds(), n.Bounds())
	}
	for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
		for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
			c := color.NRGBAModel.Convert(m.At(x, y)).(color.NRGBA)
			d := color.NRGBAModel.Convert(n.At(x, y)).(color.NRGBA)
			if !colorEq(c, d) {
				return fmt.Errorf("pixel (%v, %v) doesn't match: %+v != %+v", x, y, c, d)
			}
		}
	}
	return nil
}

// colorEq reports whether colors c, d are considered equivalent, i.e., within a small margin.
func colorEq(c, d color.NRGBA) bool {
	return eqEpsilon(c.R, d.R) && eqEpsilon(c.G, d.G) && eqEpsilon(c.B, d.B) && eqEpsilon(c.A, d.A)
}

// eqEpsilon reports whether a and b are within epsilon of each other.
func eqEpsilon(a, b uint8) bool {
	const epsilon = 1
	return uint16(a)-uint16(b) <= epsilon || uint16(b)-uint16(a) <= epsilon
}