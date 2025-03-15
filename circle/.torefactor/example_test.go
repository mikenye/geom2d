package _torefactor_test

import (
	"fmt"
	"github.com/mikenye/geom2d/circle"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"math"
)

func ExampleNew() {
	c := circle.New(3, 4, 5)
	fmt.Printf("Circle: %s\n", c)

	// Output:
	// Circle: (3,4; r=5)
}

func ExampleNewFromPoint() {
	center := point.New(3, 4)
	c := circle.NewFromPoint(center, 5)
	fmt.Printf("Circle: %s\n", c)

	// Output:
	// Circle: (3,4; r=5)
}

func ExampleCircle_Area() {
	c := circle.New(3, 4, 5) // Creates a Circle with center (3, 4) and radius 5
	area := c.Area()         // area = π*5*5 = 25π
	fmt.Printf("Circle %s has area %.2f units²\n", c, area)

	// Output:
	// Circle (3,4; r=5) has area 78.54 units²
}

func ExampleCircle_AsFloat32() {
	intCircle := circle.New(3, 4, 5)
	fltCircle := intCircle.AsFloat32()
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	// Output:
	// intCircle is (3,4; r=5) of type: circle.Circle[int]
	// fltCircle is (3,4; r=5) of type: circle.Circle[float32]
}

func ExampleCircle_AsFloat64() {
	intCircle := circle.New(3, 4, 5)
	fltCircle := intCircle.AsFloat64()
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	// Output:
	// intCircle is (3,4; r=5) of type: circle.Circle[int]
	// fltCircle is (3,4; r=5) of type: circle.Circle[float64]
}

func ExampleCircle_AsInt() {
	fltCircle := circle.New(3.7, 4.9, 5.6)
	intCircle := fltCircle.AsInt()
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	// Output:
	// fltCircle is (3.7,4.9; r=5.6) of type: circle.Circle[float64]
	// intCircle is (3,4; r=5) of type: circle.Circle[int]
}

func ExampleCircle_AsIntRounded() {
	fltCircle := circle.New(3.7, 4.2, 5.6)
	intCircle := fltCircle.AsIntRounded()
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	// Output:
	// fltCircle is (3.7,4.2; r=5.6) of type: circle.Circle[float64]
	// intCircle is (4,4; r=6) of type: circle.Circle[int]
}

func ExampleCircle_Bresenham() {
	// Define a circle with center (5, 5) and radius 3.
	c := circle.New(5, 5, 3)

	// Generate the points on the circle's perimeter.
	fmt.Printf("Integer points on circle %s perimeter:\n", c)
	for p := range c.Bresenham {
		fmt.Printf("  point: %s\n", p)
	}

	// Output:
	// Integer points on circle (5,5; r=3) perimeter:
	//   point: (5,8)
	//   point: (5,8)
	//   point: (5,2)
	//   point: (5,2)
	//   point: (8,5)
	//   point: (2,5)
	//   point: (8,5)
	//   point: (2,5)
	//   point: (6,8)
	//   point: (4,8)
	//   point: (6,2)
	//   point: (4,2)
	//   point: (8,6)
	//   point: (2,6)
	//   point: (8,4)
	//   point: (2,4)
	//   point: (7,7)
	//   point: (3,7)
	//   point: (7,3)
	//   point: (3,3)
	//   point: (7,7)
	//   point: (3,7)
	//   point: (7,3)
	//   point: (3,3)
}

func ExampleCircle_Center() {
	c := circle.New(3, 4, 5)
	center := c.Center()
	fmt.Printf("The center point of circle %s is %s", c, center)

	// Output:
	// The center point of circle (3,4; r=5) is (3,4)
}

func ExampleCircle_Circumference() {
	c := circle.New(3, 4, 5)
	circumference := c.Circumference()
	fmt.Printf("The circumference of circle %s is %.2f", c, circumference)

	// Output:
	// The circumference of circle (3,4; r=5) is 31.42
}

func ExampleCircle_Eq() {
	c1 := circle.New(3, 4, 5)
	c2 := circle.New(3, 4, 5)
	c3 := circle.New(2, 3, 4)

	fmt.Printf("Is circle c1 %s equal to circle c2 %s: %t\n", c1, c2, c1.Eq(c2))
	fmt.Printf("Is circle c1 %s equal to circle c3 %s: %t\n", c1, c3, c1.Eq(c3))

	// Output:
	// Is circle c1 (3,4; r=5) equal to circle c2 (3,4; r=5): true
	// Is circle c1 (3,4; r=5) equal to circle c3 (2,3; r=4): false
}

func ExampleCircle_Eq_epsilon() {
	c1 := circle.New[float64](3, 4, 5)
	c2 := circle.New(3.0000001, 4.0000001, 5.0000001)
	epsilon := 1e-6

	fmt.Printf(
		"Is circle c1 %s equal to circle c2 %s without epsilon: %t\n",
		c1,
		c2,
		c1.Eq(c2),
	)

	fmt.Printf(
		"Is circle c1 %s equal to circle c2 %s with an epsilon of %.0e: %t\n",
		c1,
		c2,
		epsilon,
		c1.Eq(c2, options.WithEpsilon(epsilon)),
	)

	// Output:
	// Is circle c1 (3,4; r=5) equal to circle c2 (3.0000001,4.0000001; r=5.0000001) without epsilon: false
	// Is circle c1 (3,4; r=5) equal to circle c2 (3.0000001,4.0000001; r=5.0000001) with an epsilon of 1e-06: true
}

func ExampleCircle_Radius() {
	c := circle.New(3, 4, 5)
	radius := c.Radius()
	fmt.Printf("The radius of circle %s is %v", c, radius)

	// Output:
	// The radius of circle (3,4; r=5) is 5
}

func ExampleCircle_Rotate() {
	c := circle.New(3, 3, 5)

	pivot := point.New(1, 1)
	angle := math.Pi / 2 // Rotate 90 degrees

	rotated := c.Rotate(pivot, angle, options.WithEpsilon(1e-10)).AsInt()

	fmt.Printf("Circle %s rotated 90 degrees counter-clockwise around pivot point %s is %s", c, pivot, rotated)

	// Output:
	// Circle (3,3; r=5) rotated 90 degrees counter-clockwise around pivot point (1,1) is (-1,3; r=5)
}

func ExampleCircle_Scale() {
	c := circle.New(0, 0, 5)

	scaleFactor := 2
	scaled := c.Scale(scaleFactor).AsInt()

	fmt.Printf("Circle %s scaled by a factor of %v is %s", c, scaleFactor, scaled)

	// Output:
	// Circle (0,0; r=5) scaled by a factor of 2 is (0,0; r=10)
}

func ExampleCircle_String() {
	c := circle.New(3, 4, 5)

	// When fmt.Println is used to print a variable,
	// and that variable implements the Stringer interface (String() string),
	// Go automatically calls the String() method to generate the output,
	// rather than using the default representation of the type.
	// Thus:
	fmt.Println(c)
	// is the same as:
	fmt.Println(c.String())

	// Output:
	// (3,4; r=5)
	// (3,4; r=5)
}

func ExampleCircle_Translate() {
	c := circle.New(3, 4, 5)

	translationVector := point.New(2, 3) // point used as vector

	translatedCircle := c.Translate(translationVector)

	fmt.Printf("Circle %s translated by vector %s is %s", c, translationVector, translatedCircle)

	// Output:
	// Circle (3,4; r=5) translated by vector (2,3) is (5,7; r=5)
}
