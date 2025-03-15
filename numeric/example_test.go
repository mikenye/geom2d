package numeric_test

import (
	"fmt"
	"github.com/mikenye/geom2d/numeric"
)

func ExampleFloatEquals() {
	a := 0.3333333
	b := 1.0 / 3.0
	epsilon := 1e-7

	fmt.Printf("%.7f == 1/3 without epsilon: %t\n", a, a == b)
	fmt.Printf("%.7f == 1/3 with epsilon of %.7f: %t\n", a, epsilon, numeric.FloatEquals(a, b, epsilon))

	// Output:
	// 0.3333333 == 1/3 without epsilon: false
	// 0.3333333 == 1/3 with epsilon of 0.0000001: true
}

func ExampleFloatGreaterThan() {
	a := 0.33333333333334
	b := 1.0 / 3.0
	epsilon := 1e-7

	fmt.Printf("%.14f > 1/3 without epsilon: %t\n", a, a > b)
	fmt.Printf("%.14f > 1/3 with epsilon of %.7f: %t\n", a, epsilon, numeric.FloatGreaterThan(a, b, epsilon))

	// Output:
	// 0.33333333333334 > 1/3 without epsilon: true
	// 0.33333333333334 > 1/3 with epsilon of 0.0000001: false
}

func ExampleFloatGreaterThanOrEqualTo() {
	a := 0.33333333333332
	b := 1.0 / 3.0
	epsilon := 1e-7

	fmt.Printf("%.14f >= 1/3 without epsilon: %t\n", a, a >= b)
	fmt.Printf("%.14f >= 1/3 with epsilon of %.7f: %t\n", a, epsilon, numeric.FloatGreaterThanOrEqualTo(a, b, epsilon))

	// Output:
	// 0.33333333333332 >= 1/3 without epsilon: false
	// 0.33333333333332 >= 1/3 with epsilon of 0.0000001: true
}

func ExampleFloatLessThan() {
	a := 0.33333333333332
	b := 1.0 / 3.0
	epsilon := 1e-7

	fmt.Printf("%.14f < 1/3 without epsilon: %t\n", a, a < b)
	fmt.Printf("%.14f < 1/3 with epsilon of %.7f: %t\n", a, epsilon, numeric.FloatLessThan(a, b, epsilon))

	// Output:
	// 0.33333333333332 < 1/3 without epsilon: true
	// 0.33333333333332 < 1/3 with epsilon of 0.0000001: false
}

func ExampleFloatLessThanOrEqualTo() {
	a := 0.33333333333334
	b := 1.0 / 3.0
	epsilon := 1e-7

	fmt.Printf("%.14f <= 1/3 without epsilon: %t\n", a, a <= b)
	fmt.Printf("%.14f <= 1/3 with epsilon of %.7f: %t\n", a, epsilon, numeric.FloatLessThanOrEqualTo(a, b, epsilon))

	// Output:
	// 0.33333333333334 <= 1/3 without epsilon: false
	// 0.33333333333334 <= 1/3 with epsilon of 0.0000001: true
}

func ExampleSnapToEpsilon() {
	epsilon := 0.01

	// Values close to integers should snap
	fmt.Println(numeric.SnapToEpsilon(3.0001, epsilon))
	fmt.Println(numeric.SnapToEpsilon(4.9999, epsilon))

	// Values far from integers should remain unchanged
	fmt.Println(numeric.SnapToEpsilon(3.05, epsilon))
	fmt.Println(numeric.SnapToEpsilon(-2.02, epsilon))

	// Values exactly at integers should remain unchanged
	fmt.Println(numeric.SnapToEpsilon(7.0, epsilon))

	// Negative values close to integers should snap
	fmt.Println(numeric.SnapToEpsilon(-3.0005, epsilon))
	fmt.Println(numeric.SnapToEpsilon(-5.9999, epsilon))

	// Output:
	// 3
	// 5
	// 3.05
	// -2.02
	// 7
	// -3
	// -6
}
