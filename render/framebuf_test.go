// Copyright 2021 Changkun Ou <changkun.de>. All rights reserved.
// Use of this source code is governed by a GPLv3 license that
// can be found in the LICENSE file.

package render_test

import (
	"image"
	"testing"

	"poly.red/render"
)

func BenchmarkBuffer_Clear(b *testing.B) {
	buf := render.NewBuffer(image.Rect(0, 0, 800, 800))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Clear()
	}
}
