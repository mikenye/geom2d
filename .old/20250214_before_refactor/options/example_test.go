package options_test

import (
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/options"
)

func ExampleWithEpsilon() {

	s1 := linesegment.New[float64](1, 1, 4, 5)
	s2 := linesegment.New(1.0000001, 1.0000001, 4.0000001, 5.0000001)
	epsilon := 1e-6

	fmt.Printf(
		"Is line segment s1 %s equal to line segment s2 %s without epsilon: %t\n",
		s1,
		s2,
		s1.Eq(s2),
	)

	fmt.Printf(
		"Is line segment s1 %s equal to line segment s2 %s with an epsilon of %.0e: %t\n",
		s1,
		s2,
		epsilon,
		s1.Eq(s2, options.WithEpsilon(epsilon)),
	)

	// Output:
	// Is line segment s1 (1,1)(4,5) equal to line segment s2 (1.0000001,1.0000001)(4.0000001,5.0000001) without epsilon: false
	// Is line segment s1 (1,1)(4,5) equal to line segment s2 (1.0000001,1.0000001)(4.0000001,5.0000001) with an epsilon of 1e-06: true

}
