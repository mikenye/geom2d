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

func ExampleOrientation() {
	// Define three sets of points to test different orientations

	// Counterclockwise orientation
	p0 := geom2d.NewPoint(0, 0)
	p1 := geom2d.NewPoint(1, 1)
	p2 := geom2d.NewPoint(2, 0)
	orientation1 := geom2d.Orientation(p0, p1, p2)
	fmt.Printf("Orientation of p0: %v, p1: %v, p2: %v: %v\n", p0, p1, p2, orientation1)

	// Clockwise orientation
	p3 := geom2d.NewPoint(0, 0)
	p4 := geom2d.NewPoint(2, 0)
	p5 := geom2d.NewPoint(1, 1)
	orientation2 := geom2d.Orientation(p3, p4, p5)
	fmt.Printf("Orientation of p3: %v, p4: %v, p5: %v: %v\n", p3, p4, p5, orientation2)

	// Collinear orientation
	p6 := geom2d.NewPoint(0, 0)
	p7 := geom2d.NewPoint(1, 1)
	p8 := geom2d.NewPoint(2, 2)
	orientation3 := geom2d.Orientation(p6, p7, p8)
	fmt.Printf("Orientation of points p6%v, p7%v, p8%v: %v\n", p6, p7, p8, orientation3)
	// Output:
	// Orientation of p0: Point[(0, 0)], p1: Point[(1, 1)], p2: Point[(2, 0)]: PointsClockwise
	// Orientation of p3: Point[(0, 0)], p4: Point[(2, 0)], p5: Point[(1, 1)]: PointsCounterClockwise
	// Orientation of points p6Point[(0, 0)], p7Point[(1, 1)], p8Point[(2, 2)]: PointsCollinear
}

func ExamplePoint_AsFloat64() {
	p := geom2d.NewPoint(3, 4)
	floatPoint := p.AsFloat64()
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

func ExamplePoint_RelationshipToCircle() {
	// Define a point
	point := geom2d.NewPoint(1, 1)

	// Define a circle with center (0, 0) and radius 5
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)

	// Calculate the relationship between the point and the circle
	relationship := point.RelationshipToCircle(circle, geom2d.WithEpsilon(1e-10))

	// Output the relationship
	fmt.Println(relationship)
	// Output:
	// RelationshipContainedBy
}

func ExamplePoint_RelationshipToLineSegment() {
	// Define a line segment from (0, 0) to (10, 0)
	segment := geom2d.NewLineSegment(
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
	)

	// Define points to test
	pointOnSegment := geom2d.NewPoint(5, 0) // Lies on the segment
	pointAbove := geom2d.NewPoint(5, 2)     // Lies above the segment

	// Analyze relationships
	fmt.Println(pointOnSegment.RelationshipToLineSegment(segment)) // Intersection
	fmt.Println(pointAbove.RelationshipToLineSegment(segment))     // Disjoint
	// Output:
	// RelationshipIntersection
	// RelationshipDisjoint
}

func ExamplePoint_RelationshipToPoint() {
	// Define two points
	pointA := geom2d.NewPoint(5, 5)
	pointB := geom2d.NewPoint(5, 5)
	pointC := geom2d.NewPoint(10, 10)

	// Analyze relationships
	fmt.Println(pointA.RelationshipToPoint(pointB)) // Equal
	fmt.Println(pointA.RelationshipToPoint(pointC)) // Disjoint
	// Output:
	// RelationshipEqual
	// RelationshipDisjoint
}

func ExamplePoint_RelationshipToPolyTree() {
	// Define a PolyTree with a root polygon and a hole
	root, _ := geom2d.NewPolyTree(
		[]geom2d.Point[int]{
			geom2d.NewPoint(0, 0),
			geom2d.NewPoint(0, 100),
			geom2d.NewPoint(100, 100),
			geom2d.NewPoint(100, 0),
		},
		geom2d.PTSolid,
	)
	hole, _ := geom2d.NewPolyTree(
		[]geom2d.Point[int]{
			geom2d.NewPoint(20, 20),
			geom2d.NewPoint(20, 80),
			geom2d.NewPoint(80, 80),
			geom2d.NewPoint(80, 20),
		},
		geom2d.PTHole,
	)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define points to test
	pointA := geom2d.NewPoint(50, 50)   // Inside both the hole and the root poly
	pointB := geom2d.NewPoint(5, 5)     // Inside the root polygon only
	pointC := geom2d.NewPoint(0, 10)    // On the root polygon's edge
	pointD := geom2d.NewPoint(-15, -15) // Outside the entire PolyTree

	// Analyze relationships
	relationships := pointA.RelationshipToPolyTree(root)
	fmt.Println(relationships[root]) // RelationshipContainedBy
	fmt.Println(relationships[hole]) // RelationshipContainedBy

	relationships = pointB.RelationshipToPolyTree(root)
	fmt.Println(relationships[root]) // RelationshipContainedBy
	fmt.Println(relationships[hole]) // RelationshipDisjoint

	relationships = pointC.RelationshipToPolyTree(root)
	fmt.Println(relationships[root]) // RelationshipIntersection
	fmt.Println(relationships[hole]) // RelationshipDisjoint

	relationships = pointD.RelationshipToPolyTree(root)
	fmt.Println(relationships[root]) // RelationshipDisjoint
	fmt.Println(relationships[hole]) // RelationshipDisjoint
	// Output:
	// RelationshipContainedBy
	// RelationshipContainedBy
	// RelationshipContainedBy
	// RelationshipDisjoint
	// RelationshipIntersection
	// RelationshipDisjoint
	// RelationshipDisjoint
	// RelationshipDisjoint
}

func ExamplePoint_RelationshipToRectangle() {
	// Define a rectangle
	rect := geom2d.NewRectangleByOppositeCorners(
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 10),
	)

	// Define points
	pointA := geom2d.NewPoint(5, 5)  // Inside
	pointB := geom2d.NewPoint(10, 5) // On edge
	pointC := geom2d.NewPoint(15, 5) // Outside

	// Analyze relationships
	fmt.Println(pointA.RelationshipToRectangle(rect)) // Contained
	fmt.Println(pointB.RelationshipToRectangle(rect)) // Intersection
	fmt.Println(pointC.RelationshipToRectangle(rect)) // Disjoint
	// Output:
	// RelationshipContainedBy
	// RelationshipIntersection
	// RelationshipDisjoint
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
