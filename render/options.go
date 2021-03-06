// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package render

import (
	"image"
	"image/color"
	"sync"

	"poly.red/camera"
	"poly.red/light"
	"poly.red/math"
	"poly.red/object"
	"poly.red/scene"
)

type Option func(r *Renderer)

func WithSize(width, height int) Option {
	return func(r *Renderer) {
		r.width = width
		r.height = height
	}
}

func WithCamera(cam camera.Interface) Option {
	return func(r *Renderer) {
		r.renderCamera = cam
		if _, ok := cam.(*camera.Perspective); ok {
			r.renderPerspect = true
		}
	}
}

func WithScene(s *scene.Scene) Option {
	return func(r *Renderer) {
		r.scene = s
	}
}

func WithBackground(c color.RGBA) Option {
	return func(r *Renderer) {
		r.background = c
	}
}

func WithMSAA(n int) Option {
	return func(r *Renderer) {
		r.msaa = n
	}
}

func WithShadowMap(enable bool) Option {
	return func(r *Renderer) {
		r.useShadowMap = enable
	}
}

func WithGammaCorrection(enable bool) Option {
	return func(r *Renderer) {
		r.correctGamma = enable
	}
}

func WithBlendFunc(f BlendFunc) Option {
	return func(r *Renderer) {
		r.blendFunc = f
	}
}

func WithDebug(enable bool) Option {
	return func(r *Renderer) {
		r.debug = enable
	}
}

func WithConcurrency(n int32) Option {
	return func(r *Renderer) {
		r.concurrentSize = n
	}
}

func WithThreadLimit(n int) Option {
	return func(r *Renderer) {
		r.gomaxprocs = n
	}
}

func (r *Renderer) UpdateOptions(opts ...Option) {
	r.wait() // wait last frame to finish

	for _, opt := range opts {
		opt(r)
	}

	w := r.width * r.msaa
	h := r.height * r.msaa

	// calibrate rendering size
	r.lockBuf = make([]sync.Mutex, w*h)
	r.gBuf = make([]gInfo, w*h)
	r.frameBuf = image.NewRGBA(image.Rect(0, 0, w, h))

	r.lightSources = []light.Source{}
	r.lightEnv = []light.Environment{}
	if r.scene != nil {
		r.scene.IterObjects(func(o object.Object, modelMatrix math.Mat4) bool {
			if o.Type() != object.TypeLight {
				return true
			}

			switch l := o.(type) {
			case light.Source:
				r.lightSources = append(r.lightSources, l)
			case light.Environment:
				r.lightEnv = append(r.lightEnv, l)
			}
			return true
		})
	}

	// initialize shadow maps
	if r.scene != nil && r.useShadowMap {
		r.initShadowMaps()
	}

	r.resetGBuf()
	r.resetFrameBuf()
}
