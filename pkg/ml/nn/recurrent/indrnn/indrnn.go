// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package indrnn

import (
	"encoding/gob"
	mat "github.com/nlpodyssey/spago/pkg/mat32"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"log"
)

var (
	_ nn.Model = &Model{}
)

// Model contains the serializable parameters.
type Model struct {
	nn.BaseModel
	W          nn.Param  `spago:"type:weights"`
	WRec       nn.Param  `spago:"type:weights"`
	B          nn.Param  `spago:"type:biases"`
	Activation ag.OpName // output activation
	States     []*State  `spago:"scope:processor"`
}

// State represent a state of the IndRNN recurrent network.
type State struct {
	Y ag.Node
}

func init() {
	gob.Register(&Model{})
}

// New returns a new model with parameters initialized to zeros.
func New(in, out int, activation ag.OpName) *Model {
	return &Model{
		W:          nn.NewParam(mat.NewEmptyDense(out, in)),
		WRec:       nn.NewParam(mat.NewEmptyVecDense(out)),
		B:          nn.NewParam(mat.NewEmptyVecDense(out)),
		Activation: activation,
	}
}

// SetInitialState sets the initial state of the recurrent network.
// It panics if one or more states are already present.
func (m *Model) SetInitialState(state *State) {
	if len(m.States) > 0 {
		log.Fatal("indrnn: the initial state must be set before any input")
	}
	m.States = append(m.States, state)
}

// Forward performs the forward step for each input node and returns the result.
func (m *Model) Forward(xs ...ag.Node) []ag.Node {
	ys := make([]ag.Node, len(xs))
	for i, x := range xs {
		s := m.forward(x)
		m.States = append(m.States, s)
		ys[i] = s.Y
	}
	return ys
}

// LastState returns the last state of the recurrent network.
// It returns nil if there are no states.
func (m *Model) LastState() *State {
	n := len(m.States)
	if n == 0 {
		return nil
	}
	return m.States[n-1]
}

// y = f(w (dot) x + wRec * yPrev + b)
func (m *Model) forward(x ag.Node) (s *State) {
	g := m.Graph()
	s = new(State)
	yPrev := m.prev()
	h := nn.Affine(g, m.B, m.W, x)
	if yPrev != nil {
		h = g.Add(h, g.Prod(m.WRec, yPrev))
	}
	s.Y = g.Invoke(m.Activation, h)
	return
}

func (m *Model) prev() (yPrev ag.Node) {
	s := m.LastState()
	if s != nil {
		yPrev = s.Y
	}
	return
}
