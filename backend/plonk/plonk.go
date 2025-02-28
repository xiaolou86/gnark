// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package plonk implements PLONK Zero Knowledge Proof system.
//
// # See also
//
// https://eprint.iacr.org/2019/953
package plonk

import (
	"io"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/constraint"

	"github.com/consensys/gnark/backend/witness"
	cs_bls12377 "github.com/consensys/gnark/constraint/bls12-377"
	cs_bls12381 "github.com/consensys/gnark/constraint/bls12-381"
	cs_bls24315 "github.com/consensys/gnark/constraint/bls24-315"
	cs_bls24317 "github.com/consensys/gnark/constraint/bls24-317"
	cs_bn254 "github.com/consensys/gnark/constraint/bn254"
	cs_bw6633 "github.com/consensys/gnark/constraint/bw6-633"
	cs_bw6761 "github.com/consensys/gnark/constraint/bw6-761"

	plonk_bls12377 "github.com/consensys/gnark/backend/plonk/bls12-377"
	plonk_bls12381 "github.com/consensys/gnark/backend/plonk/bls12-381"
	plonk_bls24315 "github.com/consensys/gnark/backend/plonk/bls24-315"
	plonk_bls24317 "github.com/consensys/gnark/backend/plonk/bls24-317"
	plonk_bn254 "github.com/consensys/gnark/backend/plonk/bn254"
	plonk_bw6633 "github.com/consensys/gnark/backend/plonk/bw6-633"
	plonk_bw6761 "github.com/consensys/gnark/backend/plonk/bw6-761"

	fr_bls12377 "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	fr_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	fr_bls24315 "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	fr_bls24317 "github.com/consensys/gnark-crypto/ecc/bls24-317/fr"
	fr_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	fr_bw6633 "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	fr_bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"

	kzg_bls12377 "github.com/consensys/gnark-crypto/ecc/bls12-377/fr/kzg"
	kzg_bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381/fr/kzg"
	kzg_bls24315 "github.com/consensys/gnark-crypto/ecc/bls24-315/fr/kzg"
	kzg_bls24317 "github.com/consensys/gnark-crypto/ecc/bls24-317/fr/kzg"
	kzg_bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/kzg"
	kzg_bw6633 "github.com/consensys/gnark-crypto/ecc/bw6-633/fr/kzg"
	kzg_bw6761 "github.com/consensys/gnark-crypto/ecc/bw6-761/fr/kzg"

	gnarkio "github.com/consensys/gnark/io"
)

// Proof represents a Plonk proof generated by plonk.Prove
//
// it's underlying implementation is curve specific (see gnark/internal/backend)
type Proof interface {
	io.WriterTo
	io.ReaderFrom
	gnarkio.WriterRawTo
}

// ProvingKey represents a plonk ProvingKey
//
// it's underlying implementation is strongly typed with the curve (see gnark/internal/backend)
type ProvingKey interface {
	io.WriterTo
	io.ReaderFrom
	gnarkio.WriterRawTo
	gnarkio.UnsafeReaderFrom
	VerifyingKey() interface{}
}

// VerifyingKey represents a plonk VerifyingKey
//
// it's underlying implementation is strongly typed with the curve (see gnark/internal/backend)
type VerifyingKey interface {
	io.WriterTo
	io.ReaderFrom
	gnarkio.WriterRawTo
	gnarkio.UnsafeReaderFrom
	NbPublicWitness() int // number of elements expected in the public witness
	ExportSolidity(w io.Writer) error
}

// Setup prepares the public data associated to a circuit + public inputs.
func Setup(ccs constraint.ConstraintSystem, kzgSrs kzg.SRS) (ProvingKey, VerifyingKey, error) {

	switch tccs := ccs.(type) {
	case *cs_bn254.SparseR1CS:
		return plonk_bn254.Setup(tccs, *kzgSrs.(*kzg_bn254.SRS))
	case *cs_bls12381.SparseR1CS:
		return plonk_bls12381.Setup(tccs, *kzgSrs.(*kzg_bls12381.SRS))
	case *cs_bls12377.SparseR1CS:
		return plonk_bls12377.Setup(tccs, *kzgSrs.(*kzg_bls12377.SRS))
	case *cs_bw6761.SparseR1CS:
		return plonk_bw6761.Setup(tccs, *kzgSrs.(*kzg_bw6761.SRS))
	case *cs_bls24317.SparseR1CS:
		return plonk_bls24317.Setup(tccs, *kzgSrs.(*kzg_bls24317.SRS))
	case *cs_bls24315.SparseR1CS:
		return plonk_bls24315.Setup(tccs, *kzgSrs.(*kzg_bls24315.SRS))
	case *cs_bw6633.SparseR1CS:
		return plonk_bw6633.Setup(tccs, *kzgSrs.(*kzg_bw6633.SRS))
	default:
		panic("unrecognized SparseR1CS curve type")
	}

}

// Prove generates PLONK proof from a circuit, associated preprocessed public data, and the witness
// if the force flag is set:
//
//		will execute all the prover computations, even if the witness is invalid
//	 will produce an invalid proof
//		internally, the solution vector to the SparseR1CS will be filled with random values which may impact benchmarking
func Prove(ccs constraint.ConstraintSystem, pk ProvingKey, fullWitness witness.Witness, opts ...backend.ProverOption) (Proof, error) {

	switch tccs := ccs.(type) {
	case *cs_bn254.SparseR1CS:
		return plonk_bn254.Prove(tccs, pk.(*plonk_bn254.ProvingKey), fullWitness, opts...)

	case *cs_bls12381.SparseR1CS:
		return plonk_bls12381.Prove(tccs, pk.(*plonk_bls12381.ProvingKey), fullWitness, opts...)

	case *cs_bls12377.SparseR1CS:
		return plonk_bls12377.Prove(tccs, pk.(*plonk_bls12377.ProvingKey), fullWitness, opts...)

	case *cs_bw6761.SparseR1CS:
		return plonk_bw6761.Prove(tccs, pk.(*plonk_bw6761.ProvingKey), fullWitness, opts...)

	case *cs_bw6633.SparseR1CS:
		return plonk_bw6633.Prove(tccs, pk.(*plonk_bw6633.ProvingKey), fullWitness, opts...)

	case *cs_bls24317.SparseR1CS:
		return plonk_bls24317.Prove(tccs, pk.(*plonk_bls24317.ProvingKey), fullWitness, opts...)

	case *cs_bls24315.SparseR1CS:
		return plonk_bls24315.Prove(tccs, pk.(*plonk_bls24315.ProvingKey), fullWitness, opts...)

	default:
		panic("unrecognized SparseR1CS curve type")
	}
}

// Verify verifies a PLONK proof, from the proof, preprocessed public data, and public witness.
func Verify(proof Proof, vk VerifyingKey, publicWitness witness.Witness) error {

	switch _proof := proof.(type) {

	case *plonk_bn254.Proof:
		w, ok := publicWitness.Vector().(fr_bn254.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bn254.Verify(_proof, vk.(*plonk_bn254.VerifyingKey), w)

	case *plonk_bls12381.Proof:
		w, ok := publicWitness.Vector().(fr_bls12381.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bls12381.Verify(_proof, vk.(*plonk_bls12381.VerifyingKey), w)

	case *plonk_bls12377.Proof:
		w, ok := publicWitness.Vector().(fr_bls12377.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bls12377.Verify(_proof, vk.(*plonk_bls12377.VerifyingKey), w)

	case *plonk_bw6761.Proof:
		w, ok := publicWitness.Vector().(fr_bw6761.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bw6761.Verify(_proof, vk.(*plonk_bw6761.VerifyingKey), w)

	case *plonk_bw6633.Proof:
		w, ok := publicWitness.Vector().(fr_bw6633.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bw6633.Verify(_proof, vk.(*plonk_bw6633.VerifyingKey), w)

	case *plonk_bls24317.Proof:
		w, ok := publicWitness.Vector().(fr_bls24317.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bls24317.Verify(_proof, vk.(*plonk_bls24317.VerifyingKey), w)

	case *plonk_bls24315.Proof:
		w, ok := publicWitness.Vector().(fr_bls24315.Vector)
		if !ok {
			return witness.ErrInvalidWitness
		}
		return plonk_bls24315.Verify(_proof, vk.(*plonk_bls24315.VerifyingKey), w)

	default:
		panic("unrecognized proof type")
	}
}

// NewCS instantiate a concrete curved-typed SparseR1CS and return a ConstraintSystem interface
// This method exists for (de)serialization purposes
func NewCS(curveID ecc.ID) constraint.ConstraintSystem {
	var r1cs constraint.ConstraintSystem
	switch curveID {
	case ecc.BN254:
		r1cs = &cs_bn254.SparseR1CS{}
	case ecc.BLS12_377:
		r1cs = &cs_bls12377.SparseR1CS{}
	case ecc.BLS12_381:
		r1cs = &cs_bls12381.SparseR1CS{}
	case ecc.BW6_761:
		r1cs = &cs_bw6761.SparseR1CS{}
	case ecc.BLS24_317:
		r1cs = &cs_bls24317.SparseR1CS{}
	case ecc.BLS24_315:
		r1cs = &cs_bls24315.SparseR1CS{}
	case ecc.BW6_633:
		r1cs = &cs_bw6633.SparseR1CS{}
	default:
		panic("not implemented")
	}
	return r1cs
}

// NewProvingKey instantiates a curve-typed ProvingKey and returns an interface
// This function exists for serialization purposes
func NewProvingKey(curveID ecc.ID) ProvingKey {
	var pk ProvingKey
	switch curveID {
	case ecc.BN254:
		pk = &plonk_bn254.ProvingKey{}
	case ecc.BLS12_377:
		pk = &plonk_bls12377.ProvingKey{}
	case ecc.BLS12_381:
		pk = &plonk_bls12381.ProvingKey{}
	case ecc.BW6_761:
		pk = &plonk_bw6761.ProvingKey{}
	case ecc.BLS24_317:
		pk = &plonk_bls24317.ProvingKey{}
	case ecc.BLS24_315:
		pk = &plonk_bls24315.ProvingKey{}
	case ecc.BW6_633:
		pk = &plonk_bw6633.ProvingKey{}
	default:
		panic("not implemented")
	}

	return pk
}

// NewProof instantiates a curve-typed ProvingKey and returns an interface
// This function exists for serialization purposes
func NewProof(curveID ecc.ID) Proof {
	var proof Proof
	switch curveID {
	case ecc.BN254:
		proof = &plonk_bn254.Proof{}
	case ecc.BLS12_377:
		proof = &plonk_bls12377.Proof{}
	case ecc.BLS12_381:
		proof = &plonk_bls12381.Proof{}
	case ecc.BW6_761:
		proof = &plonk_bw6761.Proof{}
	case ecc.BLS24_317:
		proof = &plonk_bls24317.Proof{}
	case ecc.BLS24_315:
		proof = &plonk_bls24315.Proof{}
	case ecc.BW6_633:
		proof = &plonk_bw6633.Proof{}
	default:
		panic("not implemented")
	}

	return proof
}

// NewVerifyingKey instantiates a curve-typed VerifyingKey and returns an interface
// This function exists for serialization purposes
func NewVerifyingKey(curveID ecc.ID) VerifyingKey {
	var vk VerifyingKey
	switch curveID {
	case ecc.BN254:
		vk = &plonk_bn254.VerifyingKey{}
	case ecc.BLS12_377:
		vk = &plonk_bls12377.VerifyingKey{}
	case ecc.BLS12_381:
		vk = &plonk_bls12381.VerifyingKey{}
	case ecc.BW6_761:
		vk = &plonk_bw6761.VerifyingKey{}
	case ecc.BLS24_317:
		vk = &plonk_bls24317.VerifyingKey{}
	case ecc.BLS24_315:
		vk = &plonk_bls24315.VerifyingKey{}
	case ecc.BW6_633:
		vk = &plonk_bw6633.VerifyingKey{}
	default:
		panic("not implemented")
	}

	return vk
}
