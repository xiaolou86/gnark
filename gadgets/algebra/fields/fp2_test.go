/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fields

import (
	"testing"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gurvy"
	"github.com/consensys/gurvy/bls377"
	"github.com/consensys/gurvy/bls377/fp"
)

//--------------------------------------------------------------------
// utils

func newOperandFp2(circuit *frontend.CS, s string) Fp2Elmt {
	res := NewFp2Elmt(circuit,
		circuit.SECRET_INPUT(s+"0"),
		circuit.SECRET_INPUT(s+"1"))
	return res
}

func assignOperandFp2(inputs map[string]interface{}, s string, w bls377.E2) {
	inputs[s+"0"] = w.A0.String()
	inputs[s+"1"] = w.A1.String()
}

//--------------------------------------------------------------------
// tests

func TestAddFp2(t *testing.T) {

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, b, c bls377.E2
	a.SetRandom()
	b.SetRandom()
	c.Add(&a, &b)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))
	fp2b := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("b0"), circuit.SECRET_INPUT("b1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.Add(&circuit, &fp2a, &fp2b)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()
	inputs["b0"] = b.A0.String()
	inputs["b1"] = b.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error AddFp2")
		}
	}
}

func TestSubFp2(t *testing.T) {

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, b, c bls377.E2
	a.SetRandom()
	b.SetRandom()
	c.Sub(&a, &b)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))
	fp2b := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("b0"), circuit.SECRET_INPUT("b1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.Sub(&circuit, &fp2a, &fp2b)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()
	inputs["b0"] = b.A0.String()
	inputs["b1"] = b.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error SubFp2")
		}
	}
}

func TestMulFp2(t *testing.T) {

	ext := Extension{uSquare: 5}

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, b, c bls377.E2
	a.SetRandom()
	b.SetRandom()
	c.Mul(&a, &b)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))
	fp2b := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("b0"), circuit.SECRET_INPUT("b1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.Mul(&circuit, &fp2a, &fp2b, ext)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()
	inputs["b0"] = b.A0.String()
	inputs["b1"] = b.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error MulFp2")
		}
	}
}

func TestMulByFpFp2(t *testing.T) {

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, c bls377.E2
	var b fp.Element
	a.SetRandom()
	b.SetRandom()
	c.MulByElement(&a, &b)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))
	fpb := circuit.SECRET_INPUT("b0")

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.MulByFp(&circuit, &fp2a, fpb)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()
	inputs["b0"] = b.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error MulByFpFp2")
		}
	}
}

func TestMulByImFp2(t *testing.T) {

	ext := Extension{uSquare: 5}

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, c bls377.E2
	a.SetRandom()
	t.Skip("missing e2.MulByNonSquare")
	// TODO c.MulByNonSquare(&a)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.MulByIm(&circuit, &fp2a, ext)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error MulByImFp2")
		}
	}
}

func TestConjugateFp2(t *testing.T) {

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, c bls377.E2
	a.SetRandom()
	c.Conjugate(&a)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.Conjugate(&circuit, &fp2a)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error ConjugateFp2")
		}
	}
}

func TestInverseFp2(t *testing.T) {

	ext := Extension{uSquare: 5}

	circuit := frontend.NewConstraintSystem()

	// witness values
	var a, c bls377.E2
	a.SetRandom()
	c.Inverse(&a)

	fp2a := NewFp2Elmt(&circuit, circuit.SECRET_INPUT("a0"), circuit.SECRET_INPUT("a1"))

	fp2c := NewFp2Elmt(&circuit, nil, nil)
	fp2c.Inverse(&circuit, &fp2a, ext)

	circuit.Tag(fp2c.X, "c0")
	circuit.Tag(fp2c.Y, "c1")

	inputs := make(map[string]interface{})
	inputs["a0"] = a.A0.String()
	inputs["a1"] = a.A1.String()

	expectedValues := make(map[string]*fp.Element)
	expectedValues["c0"] = &c.A0
	expectedValues["c1"] = &c.A1

	r1cs := circuit.ToR1CS().ToR1CS(gurvy.BW761)

	res, err := r1cs.Inspect(inputs, false)
	if err != nil {
		t.Fatal(err)
	}

	// TODO here we use string because we can't compare bls377.fp to bw761.fr elmts (add a raw cast?)
	for k, v := range res {
		var _v fp.Element
		_v.SetInterface(v)
		if !expectedValues[k].Equal(&_v) {
			t.Fatal("error InverseFp2")
		}
	}
}
