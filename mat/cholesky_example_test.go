// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"fmt"

	"github.com/jingcheng-WU/gonum/mat"
)

func ExampleCholesky() {
	// Construct a symmetric positive definite matrix.
	tmp := mat.NewDense(4, 4, []float64{
		2, 6, 8, -4,
		1, 8, 7, -2,
		2, 2, 1, 7,
		8, -2, -2, 1,
	})
	var a mat.SymDense
	a.SymOuterK(1, tmp)

	fmt.Printf("a = %0.4v\n", mat.Formatted(&a, mat.Prefix("    ")))

	// Compute the cholesky factorization.
	var chol mat.Cholesky
	if ok := chol.Factorize(&a); !ok {
		fmt.Println("a matrix is not positive semi-definite.")
	}

	// Find the determinant.
	fmt.Printf("\nThe determinant of a is %0.4g\n\n", chol.Det())

	// Use the factorization to solve the system of equations a * x = b.
	b := mat.NewVecDense(4, []float64{1, 2, 3, 4})
	var x mat.VecDense
	if err := chol.SolveVecTo(&x, b); err != nil {
		fmt.Println("Matrix is near singular: ", err)
	}
	fmt.Println("Solve a * x = b")
	fmt.Printf("x = %0.4v\n", mat.Formatted(&x, mat.Prefix("    ")))

	// Extract the factorization and check that it equals the original matrix.
	var t mat.TriDense
	chol.LTo(&t)
	var test mat.Dense
	test.Mul(&t, t.T())
	fmt.Println()
	fmt.Printf("L * Lᵀ = %0.4v\n", mat.Formatted(&a, mat.Prefix("         ")))

	// Output:
	// a = ⎡120  114   -4  -16⎤
	//     ⎢114  118   11  -24⎥
	//     ⎢ -4   11   58   17⎥
	//     ⎣-16  -24   17   73⎦
	//
	// The determinant of a is 1.543e+06
	//
	// Solve a * x = b
	// x = ⎡  -0.239⎤
	//     ⎢  0.2732⎥
	//     ⎢-0.04681⎥
	//     ⎣  0.1031⎦
	//
	// L * Lᵀ = ⎡120  114   -4  -16⎤
	//          ⎢114  118   11  -24⎥
	//          ⎢ -4   11   58   17⎥
	//          ⎣-16  -24   17   73⎦
}

func ExampleCholesky_SymRankOne() {
	a := mat.NewSymDense(4, []float64{
		1, 1, 1, 1,
		0, 2, 3, 4,
		0, 0, 6, 10,
		0, 0, 0, 20,
	})
	fmt.Printf("A = %0.4v\n", mat.Formatted(a, mat.Prefix("    ")))

	// Compute the Cholesky factorization.
	var chol mat.Cholesky
	if ok := chol.Factorize(a); !ok {
		fmt.Println("matrix a is not positive definite.")
	}

	x := mat.NewVecDense(4, []float64{0, 0, 0, 1})
	fmt.Printf("\nx = %0.4v\n", mat.Formatted(x, mat.Prefix("    ")))

	// Rank-1 update the factorization.
	chol.SymRankOne(&chol, 1, x)
	// Rank-1 update the matrix a.
	a.SymRankOne(a, 1, x)

	var au mat.SymDense
	chol.ToSym(&au)

	// Print the matrix that was updated directly.
	fmt.Printf("\nA' =        %0.4v\n", mat.Formatted(a, mat.Prefix("            ")))
	// Print the matrix recovered from the factorization.
	fmt.Printf("\nU'ᵀ * U' =  %0.4v\n", mat.Formatted(&au, mat.Prefix("            ")))

	// Output:
	// A = ⎡ 1   1   1   1⎤
	//     ⎢ 1   2   3   4⎥
	//     ⎢ 1   3   6  10⎥
	//     ⎣ 1   4  10  20⎦
	//
	// x = ⎡0⎤
	//     ⎢0⎥
	//     ⎢0⎥
	//     ⎣1⎦
	//
	// A' =        ⎡ 1   1   1   1⎤
	//             ⎢ 1   2   3   4⎥
	//             ⎢ 1   3   6  10⎥
	//             ⎣ 1   4  10  21⎦
	//
	// U'ᵀ * U' =  ⎡ 1   1   1   1⎤
	//             ⎢ 1   2   3   4⎥
	//             ⎢ 1   3   6  10⎥
	//             ⎣ 1   4  10  21⎦
}
