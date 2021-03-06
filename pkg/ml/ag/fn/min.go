// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import mat "github.com/nlpodyssey/spago/pkg/mat32"

var _ Function = &Min{}

// Min is an operator to perform element-wise min.
// y = min(x1, x2)
type Min struct {
	x1 Operand
	x2 Operand
}

// NewMin returns a new Min Function.
func NewMin(x1, x2 Operand) *Min {
	return &Min{x1: x1, x2: x2}
}

// Forward computes the output of the function.
func (r *Min) Forward() mat.Matrix {
	x1v := r.x1.Value()
	x2v := r.x2.Value()
	if !(mat.SameDims(x1v, x2v) || mat.VectorsOfSameSize(x1v, x2v)) {
		panic("fn: matrices with not compatible size")
	}
	return x1v.Minimum(x2v)
}

// Backward computes the backward pass.
func (r *Min) Backward(gy mat.Matrix) {
	x1v := r.x1.Value()
	x2v := r.x2.Value()
	if !(mat.SameDims(x1v, gy) || mat.VectorsOfSameSize(x1v, gy)) &&
		!(mat.SameDims(x2v, gy) || mat.VectorsOfSameSize(x2v, gy)) {
		panic("fn: matrices with not compatible size")
	}

	n := gy.Size()
	gyData := gy.Data()
	x1vData := x1v.Data()
	x2vData := x2v.Data()

	if r.x1.RequiresGrad() {
		gx := x1v.ZerosLike()
		defer mat.ReleaseMatrix(gx)
		gxData := gx.Data()
		for i := 0; i < n; i++ {
			if x1vData[i] < x2vData[i] {
				gxData[i] = gyData[i]
			}
		}
		r.x1.PropagateGrad(gx)
	}
	if r.x2.RequiresGrad() {
		gx := x2v.ZerosLike()
		defer mat.ReleaseMatrix(gx)
		n := gy.Size()
		gxData := gx.Data()
		for i := 0; i < n; i++ {
			if x2vData[i] < x1vData[i] {
				gxData[i] = gyData[i]
			}
		}
		r.x2.PropagateGrad(gx)
	}
}
