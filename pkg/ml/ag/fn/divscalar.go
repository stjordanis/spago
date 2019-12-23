// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fn

import (
	"brillion.io/spago/pkg/mat"
)

// The element-wise division with a scalar value.
type DivScalar struct {
	x1 Operand
	x2 Operand // scalar
}

func NewDivScalar(x1, x2 Operand) *DivScalar {
	return &DivScalar{x1: x1, x2: x2}
}

// Forward computes the output of the function.
func (r *DivScalar) Forward() mat.Matrix {
	return r.x1.Value().ProdScalar(1.0 / r.x2.Value().Scalar())
}

func (r *DivScalar) Backward(gy mat.Matrix) {
	if r.x1.RequiresGrad() {
		r.x1.PropagateGrad(gy.ProdScalar(1.0 / r.x2.Value().Scalar()))
	}
}