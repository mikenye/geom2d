package geom2d_test

import (
	"fmt"
	"geom2d"
	"image"
	"math"
)

func ExampleConvexHull() {
	// Define a set of points
	points := []geom2d.Point[int]{
		geom2d.NewPoint(1, 1),
		geom2d.NewPoint(2, 5),
		geom2d.NewPoint(3, 3),
		geom2d.NewPoint(5, 3),
		geom2d.NewPoint(3, 1),
		geom2d.NewPoint(4, 4),
		geom2d.NewPoint(5, 5),
	}

	// Compute the convex hull
	hull := geom2d.ConvexHull(points)

	// Print the points that form the convex hull
	fmt.Println("Convex Hull:")
	for _, point := range hull {
		fmt.Println(point)
	}

	// Output:
	// Convex Hull:
	// Point[(1, 1)]
	// Point[(3, 1)]
	// Point[(5, 3)]
	// Point[(5, 5)]
	// Point[(2, 5)]
}

func ExampleEnsureClockwise() {

	// Define points in counter-clockwise order
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 0),
		geom2d.NewPoint(4, 4),
	}

	// Re-order points in clockwise order
	geom2d.EnsureClockwise(points)

	// Print re-ordered points
	fmt.Println(points)
	// Output:
	// [Point[(4, 4)] Point[(4, 0)] Point[(0, 0)]]
}

func ExampleEnsureCounterClockwise() {

	// Define points in clockwise order
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 4),
		geom2d.NewPoint(4, 0),
	}

	// Re-order points in counter-clockwise order
	geom2d.EnsureCounterClockwise(points)

	// Print re-ordered points
	fmt.Println(points)
	// Output:
	// [Point[(4, 0)] Point[(4, 4)] Point[(0, 0)]]
}

func ExampleNewPoint() {
	// Create a new point with integer coordinates
	pointInt := geom2d.NewPoint[int](10, 20)
	fmt.Printf("Integer Point: (%d, %d)\n", pointInt.X(), pointInt.Y())

	// Create a new point with floating-point coordinates
	pointFloat := geom2d.NewPoint[float64](10.5, 20.25)
	fmt.Printf("Floating-Point Point: (%.2f, %.2f)\n", pointFloat.X(), pointFloat.Y())

	// Create a new point with type inference.
	// As x and y are given as integer values, pointInferred will be of type Point[int].
	pointInferred := geom2d.NewPoint(10, 20)
	fmt.Printf("Inferred Point: (%v, %v)\n", pointInferred.X(), pointInferred.Y())
	// Output:
	// Integer Point: (10, 20)
	// Floating-Point Point: (10.50, 20.25)
	// Inferred Point: (10, 20)
}

func ExampleNewPointFromImagePoint() {
	// Define an image.Point
	imgPoint := image.Point{X: 10, Y: 20}

	// Convert the image.Point to a geom2d.Point[int]
	geomPoint := geom2d.NewPointFromImagePoint(imgPoint)

	// Print the result
	fmt.Printf("Image Point: %v\n", imgPoint)
	fmt.Printf("Geometry Point: %v\n", geomPoint)
	// Output:
	// Image Point: (10,20)
	// Geometry Point: Point[(10, 20)]
}

func ExamplePoint_AsFloat() {
	p := geom2d.NewPoint(3, 4)
	floatPoint := p.AsFloat()
	fmt.Printf("%s is type %T", floatPoint, floatPoint)
	// Output:
	// Point[(3, 4)] is type geom2d.Point[float64]
}

func ExamplePoint_AsInt() {
	p := geom2d.NewPoint(3.7, 4.9)
	intPoint := p.AsInt()
	fmt.Printf("%s is type %T", intPoint, intPoint)
	// Output:
	// Point[(3, 4)] is type geom2d.Point[int]
}

func ExamplePoint_AsIntRounded() {
	p := geom2d.NewPoint(3.7, 4.2)
	roundedPoint := p.AsIntRounded()
	fmt.Printf("%s is type %T", roundedPoint, roundedPoint)
	// Output:
	// Point[(4, 4)] is type geom2d.Point[int]
}

func ExamplePoint_CrossProduct() {
	// Define two points relative to the origin
	p := geom2d.NewPoint[int](1, 2) // Point p is at (1, 2)
	q := geom2d.NewPoint[int](2, 1) // Point q is at (2, 1)

	// Calculate the cross product of vectors from the origin to p and q
	crossProduct := p.CrossProduct(q)

	// Analyze the result
	fmt.Printf("Cross product: %d\n", crossProduct)
	if crossProduct > 0 {
		fmt.Println("The points indicate a counterclockwise turn (left turn).")
	} else if crossProduct < 0 {
		fmt.Println("The points indicate a clockwise turn (right turn).")
	} else {
		fmt.Println("The points are collinear.")
	}

	// Output:
	// Cross product: -3
	// The points indicate a clockwise turn (right turn).
}

func ExamplePoint_DistanceSquaredToPoint() {
	// Define two points
	p := geom2d.NewPoint[int](3, 4)
	q := geom2d.NewPoint[int](6, 8)

	// Calculate the squared Euclidean distance between the points
	distanceSquared := p.DistanceSquaredToPoint(q)

	// Display the result
	fmt.Printf("The squared distance between %v and %v is %d.\n", p, q, distanceSquared)
	// Output:
	// The squared distance between Point[(3, 4)] and Point[(6, 8)] is 25.
}

func ExamplePoint_DistanceToLineSegment() {
	// Define a point
	p := geom2d.NewPoint[float64](3, 4)

	// Define a line segment
	lineSegment := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](6, 0),
	)

	// Calculate the shortest distance from the point to the line segment
	distance := p.DistanceToLineSegment(lineSegment, geom2d.WithEpsilon(1e-10))

	// Display the result
	fmt.Printf("The shortest distance from %v to %v is %.2f.\n", p, lineSegment, distance)
	// Output:
	// The shortest distance from Point[(3, 4)] to LineSegment[(0, 0) -> (6, 0)] is 4.00.
}

func ExamplePoint_DistanceToPoint() {
	// Define two points
	p1 := geom2d.NewPoint[float64](3, 4)
	p2 := geom2d.NewPoint[float64](0, 0)

	// Calculate the Euclidean distance between the points
	distance := p1.DistanceToPoint(p2, geom2d.WithEpsilon(1e-10))

	// Display the result
	fmt.Printf("The Euclidean distance between %v and %v is %.2f.\n", p1, p2, distance)
	// Output:
	// The Euclidean distance between Point[(3, 4)] and Point[(0, 0)] is 5.00.
}

func ExamplePoint_DotProduct() {
	// Define two vectors as points
	p1 := geom2d.NewPoint[float64](3, 4)
	p2 := geom2d.NewPoint[float64](1, 2)

	// Calculate the dot product of the vectors
	dotProduct := p1.DotProduct(p2)

	// Display the result
	fmt.Printf("The dot product of vector %v and vector %v is %.2f.\n", p1, p2, dotProduct)
	// Output:
	// The dot product of vector Point[(3, 4)] and vector Point[(1, 2)] is 11.00.
}

func ExamplePoint_Eq() {
	p := geom2d.NewPoint(3.0, 4.0)
	q := geom2d.NewPoint(3.0, 4.0)

	// Exact equality
	isEqual := p.Eq(q)
	fmt.Printf("Are %s and %s equal? %t.\n", p, q, isEqual)

	// Approximate equality with epsilon
	r := geom2d.NewPoint(3.000001, 4.000001)
	isApproximatelyEqual := p.Eq(r, geom2d.WithEpsilon(1e-5))
	fmt.Printf("Are %s and %s equal within 5 decimal places? %t.\n", p, r, isApproximatelyEqual)
	// Output:
	// Are Point[(3, 4)] and Point[(3, 4)] equal? true.
	// Are Point[(3, 4)] and Point[(3.000001, 4.000001)] equal within 5 decimal places? true.
}

func ExamplePoint_Negate() {
	// Define a point
	p := geom2d.NewPoint(3, -4)

	// Negate the point
	negated := p.Negate()

	// Output the result
	fmt.Println("Original Point:", p)
	fmt.Println("Negated Point:", negated)
	// Output:
	// Original Point: Point[(3, -4)]
	// Negated Point: Point[(-3, 4)]
}

func ExamplePoint_ProjectOntoLineSegment() {
	// Define a point to be projected
	point := geom2d.NewPoint[float64](3, 7)

	// Define a line segment
	lineSegment := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 0),
	)

	// Project the point onto the line segment
	projectedPoint := point.ProjectOntoLineSegment(lineSegment)

	// Display the result
	fmt.Printf("%v projects onto the %v at %v.\n", point, lineSegment, projectedPoint)
	// Output:
	// Point[(3, 7)] projects onto the LineSegment[(0, 0) -> (10, 0)] at Point[(3, 0)].
}

func ExamplePoint_Reflect() {
	// Define the point to be reflected
	point := geom2d.NewPoint[float64](3, 4)

	// Reflect across the X-axis
	reflectedX := point.Reflect(geom2d.ReflectAcrossXAxis)

	// Reflect across the Y-axis
	reflectedY := point.Reflect(geom2d.ReflectAcrossYAxis)

	// Define a custom line for reflection
	line := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 10),
	)

	// Reflect across the custom line
	reflectedLine := point.Reflect(geom2d.ReflectAcrossCustomLine, line)

	// Print the results
	fmt.Printf("Original Point: %v\n", point)
	fmt.Printf("Reflected across X-axis: %v\n", reflectedX)
	fmt.Printf("Reflected across Y-axis: %v\n", reflectedY)
	fmt.Printf("Reflected across custom %v: %v\n", line, reflectedLine)
	// Output:
	// Original Point: Point[(3, 4)]
	// Reflected across X-axis: Point[(3, -4)]
	// Reflected across Y-axis: Point[(-3, 4)]
	// Reflected across custom LineSegment[(0, 0) -> (10, 10)]: Point[(4, 3)]
}

func ExamplePoint_RelationshipToLineSegment() {
	// Define a line segment
	segment := geom2d.NewLineSegment(
		geom2d.NewPoint[int](0, 0),
		geom2d.NewPoint[int](10, 0),
	)

	// Define points to analyze
	pointOnSegment := geom2d.NewPoint[int](5, 0)
	pointOnLine := geom2d.NewPoint[int](15, 0)
	pointAtStart := geom2d.NewPoint[int](0, 0)
	pointAtEnd := geom2d.NewPoint[int](10, 0)
	pointMiss := geom2d.NewPoint[int](5, 5)

	// Analyze relationships
	fmt.Println(pointOnSegment.RelationshipToLineSegment(segment))
	fmt.Println(pointOnLine.RelationshipToLineSegment(segment))
	fmt.Println(pointAtStart.RelationshipToLineSegment(segment))
	fmt.Println(pointAtEnd.RelationshipToLineSegment(segment))
	fmt.Println(pointMiss.RelationshipToLineSegment(segment))
	// Output:
	// RelationshipPointLineSegmentPointOnLineSegment
	// RelationshipPointLineSegmentCollinearDisjoint
	// RelationshipPointLineSegmentPointEqStart
	// RelationshipPointLineSegmentPointEqEnd
	// RelationshipPointLineSegmentMiss
}

func ExamplePoint_Rotate() {
	pivot := geom2d.NewPoint(1.0, 1.0)
	circle := geom2d.NewCircle(geom2d.NewPoint(3.0, 3.0), 5.0)

	// Rotates the circle 90 degrees around (1.0, 1.0)
	rotatedCircle := circle.Rotate(pivot, math.Pi/2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rotatedCircle)
	// Output:
	// Circle[center=(-1, 3), radius=5]
}

func ExamplePoint_Scale() {
	p := geom2d.NewPoint(3, 4)
	ref := geom2d.NewPoint(1, 1)
	scaled := p.Scale(ref, 2)
	fmt.Println(scaled)
	// Output:
	// Point[(5, 7)]
}

func ExamplePoint_String() {
	p := geom2d.NewPoint(3, 4)
	fmt.Println(p.String())
	// Output:
	// Point[(3, 4)]
}

func ExamplePoint_Translate() {
	// Create a point
	p := geom2d.NewPoint(5, 5)

	// Translate the point by a positive delta
	delta := geom2d.NewPoint(3, -2)
	pTranslated := p.Translate(delta)
	fmt.Println("Translated point:", pTranslated)

	// Subtract a vector by negating it and then using Translate
	subtractVector := geom2d.NewPoint(1, 1).Negate()
	pSubtracted := p.Translate(subtractVector)
	fmt.Println("Point after subtraction:", pSubtracted)
	// Output:
	// Translated point: Point[(8, 3)]
	// Point after subtraction: Point[(4, 4)]
}

func ExamplePoint_X() {
	p := geom2d.NewPoint(3, 4)
	fmt.Println(p.X())
	// Output:
	// 3
}

func ExamplePoint_Y() {
	p := geom2d.NewPoint(3, 4)
	fmt.Println(p.Y())
	// Output:
	// 4
}

func ExampleRelativeAngle() {
	A := geom2d.NewPoint[int](1, 0)
	B := geom2d.NewPoint[int](0, 1)
	radians := geom2d.RelativeAngle(A, B) // π/2 radians
	degrees := radians * 180 / math.Pi    // convert to degrees
	fmt.Println(degrees)
	// Output:
	// 90
}

func ExampleRelativeCosineOfAngle() {
	A := geom2d.NewPoint(1, 0)
	B := geom2d.NewPoint(0, 1)
	cosine := geom2d.RelativeCosineOfAngle(A, B) // cosine is 0 for a 90-degree angle
	fmt.Println(cosine)
	// Output:
	// 0
}

func ExampleSignedArea2X() {
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 0),
		geom2d.NewPoint(4, 3),
	}
	signedArea := geom2d.SignedArea2X(points)
	fmt.Println(signedArea)
	// Output:
	// 12
}
