package point_test

import (
	"fmt"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"image"
	"math"
)

func ExampleNew() {
	// Create a new point with integer coordinates
	pointInt := point.New[int](10, 20)
	fmt.Printf("Integer Point: %s, type %T\n", pointInt, pointInt)

	// Create a new point with floating-point coordinates
	pointFloat := point.New[float64](10.5, 20.25)
	fmt.Printf("Floating-Point Point: %s, type %T\n", pointFloat, pointFloat)

	// Create a new point with type inference.
	// As x and y are given as integer values, pointInferred will be of type Point[int].
	pointInferred := point.New(10, 20)
	fmt.Printf("Inferred Point: %s, type %T\n", pointInferred, pointInferred)

	// Output:
	// Integer Point: (10,20), type point.Point[int]
	// Floating-Point Point: (10.5,20.25), type point.Point[float64]
	// Inferred Point: (10,20), type point.Point[int]
}

func ExampleNewFromImagePoint() {
	// Define an image.Point
	imgPoint := image.Point{X: 10, Y: 20}

	// Convert the image.Point to a geom2d.Point[int]
	geomPoint := point.NewFromImagePoint(imgPoint)

	// Print the result
	fmt.Printf("Image Point: %s, type %T\n", imgPoint, imgPoint)
	fmt.Printf("geom2d Point: %s, type %T\n", geomPoint, geomPoint)

	// Output:
	// Image Point: (10,20), type image.Point
	// geom2d Point: (10,20), type point.Point[int]

}

func ExamplePoint_AngleBetween() {
	origin := point.New(0, 0)
	pointA := point.New(10, 0)
	pointB := point.New(10, 10)

	radians := origin.AngleBetween(pointA, pointB, options.WithEpsilon(1e-8))
	degrees := radians * 180 / math.Pi

	fmt.Printf(
		"The angle between point %s and point %s relative to point %s is %0.0f degrees",
		pointA,
		pointB,
		origin,
		degrees,
	)

	// Output:
	// The angle between point (10,0) and point (10,10) relative to point (0,0) is 45 degrees
}

func ExamplePoint_AsFloat32() {
	intPoint := point.New(3, 4)
	fltPoint := intPoint.AsFloat32()
	fmt.Printf("intPoint is %v of type: %T\n", intPoint, intPoint)
	fmt.Printf("fltPoint is %v of type: %T\n", fltPoint, fltPoint)
	// Output:
	// intPoint is (3,4) of type: point.Point[int]
	// fltPoint is (3,4) of type: point.Point[float32]
}

func ExamplePoint_AsFloat64() {
	intPoint := point.New(3, 4)
	fltPoint := intPoint.AsFloat64()
	fmt.Printf("intPoint is %v of type: %T\n", intPoint, intPoint)
	fmt.Printf("fltPoint is %v of type: %T\n", fltPoint, fltPoint)
	// Output:
	// intPoint is (3,4) of type: point.Point[int]
	// fltPoint is (3,4) of type: point.Point[float64]
}

func ExamplePoint_AsInt() {
	fltPoint := point.New(3.7, 4.1)
	intPoint := fltPoint.AsInt()
	fmt.Printf("fltPoint is %v of type: %T\n", fltPoint, fltPoint)
	fmt.Printf("intPoint is %v of type: %T\n", intPoint, intPoint)
	// Output:
	// fltPoint is (3.7,4.1) of type: point.Point[float64]
	// intPoint is (3,4) of type: point.Point[int]
}

func ExamplePoint_AsIntRounded() {
	fltPoint := point.New(3.7, 4.1)
	intPoint := fltPoint.AsIntRounded()
	fmt.Printf("fltPoint is %v of type: %T\n", fltPoint, fltPoint)
	fmt.Printf("intPoint is %v of type: %T\n", intPoint, intPoint)
	// Output:
	// fltPoint is (3.7,4.1) of type: point.Point[float64]
	// intPoint is (4,4) of type: point.Point[int]
}

func ExamplePoint_Coordinates() {
	p := point.New(5, -3)

	x, y := p.Coordinates()
	fmt.Printf("Point coordinates: (%d, %d)\n", x, y)

	// Output:
	// Point coordinates: (5, -3)
}

func ExamplePoint_CosineOfAngleBetween() {
	origin := point.New(0, 0)
	pointA := point.New(10, 0)
	pointB := point.New(10, 10)

	cosineOfAngle := origin.CosineOfAngleBetween(pointA, pointB)

	fmt.Printf(
		"The cosine of the angle between point %s and point %s relative to point %s is %0.6f",
		pointA,
		pointB,
		origin,
		cosineOfAngle,
	)

	// Output:
	// The cosine of the angle between point (10,0) and point (10,10) relative to point (0,0) is 0.707107
}

func ExamplePoint_DistanceSquaredToPoint() {
	// Define two points
	p := point.New[int](3, 4)
	q := point.New[int](6, 8)

	// Calculate the squared Euclidean distance between the points
	distanceSquared := p.DistanceSquaredToPoint(q)

	// Display the result
	fmt.Printf("The squared distance between %v and %v is %d\n", p, q, distanceSquared)

	// Output:
	// The squared distance between (3,4) and (6,8) is 25
}

func ExamplePoint_DistanceToPoint() {
	// Define two points
	p1 := point.New[float64](3, 4)
	p2 := point.New[float64](0, 0)

	// Calculate the Euclidean distance between the points
	distance := p1.DistanceToPoint(p2, options.WithEpsilon(1e-8))

	// Display the result
	fmt.Printf("The Euclidean distance between %v and %v is %.2f\n", p1, p2, distance)

	// Output:
	// The Euclidean distance between (3,4) and (0,0) is 5.00
}

func ExamplePoint_DotProduct() {
	// Define two vectors as points
	p1 := point.New[float64](3, 4)
	p2 := point.New[float64](1, 2)

	// Calculate the dot product of the vectors
	dotProduct := p1.DotProduct(p2)

	// Display the result
	fmt.Printf("The dot product of vector %v and vector %v is %.2f\n", p1, p2, dotProduct)

	// Output:
	// The dot product of vector (3,4) and vector (1,2) is 11.00
}

func ExamplePoint_Eq() {
	p := point.New[float64](3, 4)
	q := point.New[float64](3, 4)

	isEqual := p.Eq(q)
	fmt.Printf("Are %s and %s equal: %t\n", p, q, isEqual)

	// Output:
	// Are (3,4) and (3,4) equal: true
}

func ExamplePoint_Eq_epsilon() {
	p := point.New[float64](3, 4)
	q := point.New[float64](3.00000000001, 4.00000000001)
	epsilon := 1e-8

	isEqual := p.Eq(q, options.WithEpsilon(epsilon))
	fmt.Printf("Are %s and %s equal: %t (with epsilon: %0.0e)\n", p, q, isEqual, epsilon)

	// Output:
	// Are (3,4) and (3.00000000001,4.00000000001) equal: true (with epsilon: 1e-08)
}

func ExamplePoint_Negate() {
	// Define a point
	p := point.New(3, -4)

	// Negate the point
	negated := p.Negate()

	// Output the result
	fmt.Println("Original Point:", p)
	fmt.Println("Negated Point:", negated)

	// Output:
	// Original Point: (3,-4)
	// Negated Point: (-3,4)
}

func ExamplePoint_Rotate() {
	pivot := point.New(0, 0)
	p := point.New(10, 0)
	radians := math.Pi / 2 // 90 degrees

	// Rotates the point 90 degrees counter-clockwise around (0, 0)
	rotated := p.Rotate(pivot, radians, options.WithEpsilon(1e-8))

	fmt.Printf(
		"Point %s rotated 90 degrees counter-clockwise around %s is: %s\n",
		p,
		pivot,
		rotated,
	)

	// Output:
	// Point (10,0) rotated 90 degrees counter-clockwise around (0,0) is: (0,10)
}

func ExamplePoint_Scale() {
	p := point.New(3, 4)
	ref := point.New(1, 1)
	factor := 2

	scaled := p.Scale(ref, factor)

	fmt.Printf(
		"Point %s scaled by a factor of %v relative to reference point %s is %s\n",
		p,
		factor,
		ref,
		scaled,
	)

	// Output:
	// Point (3,4) scaled by a factor of 2 relative to reference point (1,1) is (5,7)
}

func ExamplePoint_String() {
	p := point.New(1, 2)

	// When fmt.Println is used to print a variable,
	// and that variable implements the Stringer interface (String() string),
	// Go automatically calls the String() method to generate the output,
	// rather than using the default representation of the type.
	// Thus:
	fmt.Println(p)
	// is the same as:
	fmt.Println(p.String())

	// Output:
	// (1,2)
	// (1,2)
}

func ExamplePoint_Translate() {
	p := point.New(1, 2)
	delta := point.New(-2, -4)

	translated := p.Translate(delta)

	fmt.Printf("Point %s translated by %s is %s\n", p, delta, translated)

	// Output:
	// Point (1,2) translated by (-2,-4) is (-1,-2)
}

func ExamplePoint_X() {
	p := point.New(1, 2)

	fmt.Printf("The X coordinate of point %s is %d", p, p.X())

	// Output:
	// The X coordinate of point (1,2) is 1
}

func ExamplePoint_Y() {
	p := point.New(1, 2)

	fmt.Printf("The Y coordinate of point %s is %d", p, p.X())

	// Output:
	// The Y coordinate of point (1,2) is 1
}
