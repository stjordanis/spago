// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nn

import (
	"brillion.io/spago/pkg/mat"
	"brillion.io/spago/pkg/ml/ag"
	"brillion.io/spago/pkg/ml/optimizers/gd"
	"brillion.io/spago/pkg/utils"
	"io"
	"reflect"
	"strings"
)

// Model contains the serializable parameters.
type Model interface {
	utils.SerializerDeserializer
	// ForEachParam iterates all the parameters exploring the sub-parameters recursively.
	ForEachParam(callback func(param *Param))
	// NewProc returns a new processor to execute the forward step.
	NewProc(g *ag.Graph, opt ...interface{}) Processor
}

// Processor performs the operations on the computational graphs using the model's parameters.
type Processor interface {
	// Model returns the model the processor belongs to.
	Model() Model
	// Graph returns the computational graph on which the processor operates.
	Graph() *ag.Graph
	// Reset the processor to the initial configuration (e.g. clear all the states of recurrent networks), with the init options.
	Reset()
	// Whether the processor needs the complete sequence to start processing (as in the case of BiRNN and other bidirectional models), or not.
	RequiresFullSeq() bool
	// Forward performs the the forward step for each input and returns the result.
	// Recurrent networks treats the input nodes as a sequence.
	// Differently, feed-forward networks are stateless so every computation is independent.
	Forward(xs ...ag.Node) []ag.Node
}

// ForEachParam iterate all the parameters also exploring the sub-parameters recursively.
// TODO: don't loop the field every time, use a lazy initialized "params list" instead.
func ForEachParam(m Model, callback func(param *Param)) {
	utils.ForEachField(m, func(field interface{}, name string, tag reflect.StructTag) {
		switch item := field.(type) {
		case *Param:
			item.name = strings.ToLower(name)
			item.pType = ToType(tag.Get("type"))
			callback(item)
		case Model:
			item.ForEachParam(callback)
		}
	})
}

// TrackParams tells the optimizer to track all model parameters.
// TODO: move away from here
func TrackParams(m Model, o *gd.GradientDescent) {
	m.ForEachParam(func(param *Param) {
		o.Track(param)
	})
}

// ZeroGrad set the gradients of all model's parameters (including sub-params) to zeros.
func ZeroGrad(m Model) {
	m.ForEachParam(func(param *Param) {
		param.ZeroGrad()
	})
}

// ClearSupport clears the support structure of all model's parameters (including sub-params).
func ClearSupport(m Model) {
	m.ForEachParam(func(param *Param) {
		param.ClearSupport()
	})
}

// Serialize dumps the model to the writer.
func Serialize(model Model, w io.Writer) (n int, err error) {
	model.ForEachParam(func(param *Param) {
		cnt, err := mat.MarshalBinaryTo(param.Value(), w)
		n += cnt
		if err != nil {
			return
		}
	})
	return n, err
}

// Deserialize loads the model from the reader.
func Deserialize(model Model, r io.Reader) (n int, err error) {
	model.ForEachParam(func(param *Param) {
		cnt, err := mat.UnmarshalBinaryFrom(param.Value(), r)
		n += cnt
		if err != nil {
			return
		}
	})
	return n, err
}