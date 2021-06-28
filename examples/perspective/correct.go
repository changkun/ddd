package main

import (
	"changkun.de/x/polyred/camera"
	"changkun.de/x/polyred/color"
	"changkun.de/x/polyred/geometry"
	"changkun.de/x/polyred/image"
	"changkun.de/x/polyred/io"
	"changkun.de/x/polyred/light"
	"changkun.de/x/polyred/material"
	"changkun.de/x/polyred/math"
	"changkun.de/x/polyred/render"
	"changkun.de/x/polyred/scene"
	"changkun.de/x/polyred/utils"
)

func main() {
	s := scene.NewScene()
	s.SetCamera(camera.NewPerspective(
		math.NewVector(2, 2, 2, 1),
		math.NewVector(0, 0, 0, 1),
		math.NewVector(0, 1, 0, 0),
		45,
		1,
		0.1, 10,
	))
	s.Add(light.NewPoint(
		light.WithPointLightIntensity(10),
		light.WithPointLightColor(color.RGBA{0, 128, 255, 255}),
		light.WithPointLightPosition(math.NewVector(2, 2, 2, 1)),
	))

	m := geometry.NewPlane(1, 1)
	m.SetMaterial(material.NewBlinnPhong(
		material.WithBlinnPhongTexture(image.NewTexture(
			image.WithSource(io.MustLoadImage("../../testdata/uvgrid2.png")),
			image.WithIsotropicMipMap(true),
		)),
		material.WithBlinnPhongFactors(0.6, 0.5),
		material.WithBlinnPhongShininess(150),
	))
	m.Scale(2, 2, 2)
	m.Rotate(math.NewVector(0, 1, 0, 0), math.Pi/4)
	s.Add(m)

	r := render.NewRenderer(
		render.WithSize(500, 500),
		render.WithScene(s),
		render.WithBackground(color.FromHex("#ffffff")),
	)
	utils.Save(r.Render(), "./persp.png")

}
