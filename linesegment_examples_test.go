package geom2d_test

import (
	"fmt"
	"geom2d"
	"math"
)

func ExampleNewLineSegment() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(3, 4))
	fmt.Println(segment.String())
	// Output:
	// LineSegment[(0, 0) -> (3, 4)]
}

func ExampleLineSegment_AddLineSegment() {
	segment1 := geom2d.NewLineSegment(geom2d.NewPoint(1, 1), geom2d.NewPoint(4, 5))
	segment2 := geom2d.NewLineSegment(geom2d.NewPoint(2, 3), geom2d.NewPoint(1, 2))
	result := segment1.AddLineSegment(segment2)
	fmt.Println(result)
	// Output:
	// LineSegment[(3, 4) -> (5, 7)]
}

func ExampleLineSegment_Area() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(1, 2), geom2d.NewPoint(3, 4))
	fmt.Println(segment.Area())
	// Output:
	// 0
}

func ExampleLineSegment_AsFloat() {
	// Create a LineSegment with integer coordinates
	intSegment := geom2d.NewLineSegment(
		geom2d.NewPoint(1, 2), // Start point
		geom2d.NewPoint(3, 4), // End point
	)

	// Convert the LineSegment to float64
	floatSegment := intSegment.AsFloat()

	// Print the converted LineSegment
	fmt.Println("Integer LineSegment:", intSegment)
	fmt.Println("Float64 LineSegment:", floatSegment)
	// Output:
	// Integer LineSegment: LineSegment[(1, 2) -> (3, 4)]
	// Float64 LineSegment: LineSegment[(1, 2) -> (3, 4)]
}

func ExampleLineSegment_AsInt() {
	// Create a LineSegment with floating-point coordinates
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(1.5, 2.7),
		geom2d.NewPoint(3.9, 4.2),
	)

	// Convert the LineSegment to integer coordinates
	intLine := line.AsInt()

	// Output the integer LineSegment
	fmt.Println(intLine)
	// Output:
	// LineSegment[(1, 2) -> (3, 4)]
}

func ExampleLineSegment_AsIntRounded() {
	// Create a LineSegment with floating-point coordinates
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(1.5, 2.7),
		geom2d.NewPoint(3.9, 4.2),
	)

	// Convert the LineSegment to integer coordinates with rounding
	roundedIntLine := line.AsIntRounded()

	// Output the rounded integer LineSegment
	fmt.Println(roundedIntLine)
	// Output:
	// LineSegment[(2, 3) -> (4, 4)]
}

func ExampleLineSegment_BoundingBox() {
	// Create a LineSegment
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(3, 1),
		geom2d.NewPoint(6, 4),
	)

	// Compute the bounding box of the LineSegment
	boundingBox := line.BoundingBox()

	// Output the bounding box
	fmt.Println(boundingBox)
	// Output:
	// Rectangle[(3, 1), (6, 1), (6, 4), (3, 4)]
}

func ExampleLineSegment_Center() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))
	center := lineSegment.Center()
	fmt.Printf("Center: %v\n", center)
	// Output:
	// Center: Point[(5, 5)]
}

func ExampleLineSegment_ContainsPoint() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))

	pointOnSegment := geom2d.NewPoint(5, 5)
	pointOffSegment := geom2d.NewPoint(7, 8)

	fmt.Printf("Point on segment: %v\n", lineSegment.ContainsPoint(pointOnSegment))
	fmt.Printf("Point off segment: %v\n", lineSegment.ContainsPoint(pointOffSegment))
	// Output:
	// Point on segment: true
	// Point off segment: false
}

func ExampleLineSegment_DistanceToLineSegment() {
	segmentAB := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(2, 2))
	segmentCD := geom2d.NewLineSegment(geom2d.NewPoint(3, 3), geom2d.NewPoint(5, 5))

	// Default behavior (no epsilon adjustment)
	distance := segmentAB.DistanceToLineSegment(segmentCD)
	distanceRounded := math.Round(distance*1000) / 1000
	fmt.Println(distanceRounded)
	// Output:
	// 1.414
}

func ExampleLineSegment_DistanceToPoint() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 0))

	point := geom2d.NewPoint(5, 5)
	distance := lineSegment.DistanceToPoint(point)

	fmt.Printf("Distance from point to line segment: %.2f\n", distance)
	// Output:
	// Distance from point to line segment: 5.00
}

func ExampleLineSegment_End() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(1, 2), geom2d.NewPoint(3, 4))
	fmt.Println(lineSegment.End())
	// Output:
	// Point[(3, 4)]
}

func ExampleLineSegment_Eq() {
	segment1 := geom2d.NewLineSegment(geom2d.NewPoint(1.0, 1.0), geom2d.NewPoint(4.0, 5.0))
	segment2 := geom2d.NewLineSegment(geom2d.NewPoint(1.0, 1.0), geom2d.NewPoint(4.0, 5.0))
	fmt.Println(segment1.Eq(segment2))
	// Output:
	// true
}

func ExampleLineSegment_Eq_epsilon() {
	// Approximate equality with epsilon
	segment1 := geom2d.NewLineSegment(geom2d.NewPoint(1.0, 1.0), geom2d.NewPoint(4.0, 5.0))
	segment3 := geom2d.NewLineSegment(geom2d.NewPoint(1.00001, 1.00001), geom2d.NewPoint(4.00001, 5.00001))
	fmt.Println(segment1.Eq(segment3, geom2d.WithEpsilon(1e-4)))
	// Output:
	// true
}

func ExampleLineSegment_IntersectionPoint() {
	// Define two line segments
	AB := geom2d.NewLineSegment(geom2d.NewPoint[float64](0, 0), geom2d.NewPoint[float64](4, 4))
	CD := geom2d.NewLineSegment(geom2d.NewPoint[float64](0, 4), geom2d.NewPoint[float64](4, 0))

	// Find the intersection point
	intersection, exists := AB.IntersectionPoint(CD)

	// Print the result
	if exists {
		fmt.Printf("Intersection point: (%.2f, %.2f)\n", intersection.X(), intersection.Y())
	} else {
		fmt.Println("No intersection point exists.")
	}

	// Output:
	// Intersection point: (2.00, 2.00)
}

func ExampleLineSegment_Perimeter() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(3, 4))
	fmt.Println(segment.Perimeter())
	// Output:
	// 5
}

func ExampleLineSegment_Points() {
	// Create a line segment with two endpoints
	line := geom2d.NewLineSegment(geom2d.NewPoint(1, 2), geom2d.NewPoint(3, 4))

	// Get the points as a slice
	points := line.Points()

	// Output the points
	fmt.Printf("Start Point: %v\n", points[0])
	fmt.Printf("End Point: %v\n", points[1])
	// Output:
	// Start Point: Point[(1, 2)]
	// End Point: Point[(3, 4)]
}

func ExampleLineSegment_Reflect() {
	// Create a line segment
	line := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](1, 2),
		geom2d.NewPoint[float64](3, 4),
	)

	// Reflect across the X-axis
	reflectedX := line.Reflect(geom2d.ReflectAcrossXAxis)

	// Reflect across the Y-axis
	reflectedY := line.Reflect(geom2d.ReflectAcrossYAxis)

	// Reflect across a custom line (e.g., y = x, represented as LineSegment[(0, 0), (1, 1)])
	customLine := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](1, 1),
	)
	reflectedCustom := line.Reflect(geom2d.ReflectAcrossCustomLine, customLine)

	// Output the results
	fmt.Printf("Original Line: %v\n", line)
	fmt.Printf("Reflected across X-axis: %v\n", reflectedX)
	fmt.Printf("Reflected across Y-axis: %v\n", reflectedY)
	fmt.Printf("Reflected across custom line (y = x): %v\n", reflectedCustom)
	// Output:
	// Original Line: LineSegment[(1, 2) -> (3, 4)]
	// Reflected across X-axis: LineSegment[(1, -2) -> (3, -4)]
	// Reflected across Y-axis: LineSegment[(-1, 2) -> (-3, 4)]
	// Reflected across custom line (y = x): LineSegment[(2, 1) -> (4, 3)]
}

func ExampleLineSegment_RelationshipToCircle() {
	// Define a circle with center (5, 5) and radius 3
	circle := geom2d.NewCircle(
		geom2d.NewPoint[float64](5, 5),
		3.0,
	)

	// Define various line segments to test against the circle
	lineOutside := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](1, 1),
	)
	lineTangent := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](2, 8),
		geom2d.NewPoint[float64](2, 2),
	)
	lineIntersecting := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](2, 2),
		geom2d.NewPoint[float64](8, 8),
	)
	lineInside := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](4, 4),
		geom2d.NewPoint[float64](6, 6),
	)
	lineOneEndOnCircumferenceOutside := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](5, 8),
		geom2d.NewPoint[float64](10, 10),
	)
	lineOneEndOnCircumferenceInside := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](5, 8),
		geom2d.NewPoint[float64](5, 6),
	)
	lineBothEndsOnCircumference := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](5, 8),
		geom2d.NewPoint[float64](5, 2),
	)

	// Analyze relationships
	fmt.Println(
		"Line Outside:",
		lineOutside.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"Line Tangent:",
		lineTangent.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"Line Intersecting:",
		lineIntersecting.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"Line Inside:",
		lineInside.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"One End On Circumference Outside:",
		lineOneEndOnCircumferenceOutside.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"One End On Circumference Inside:",
		lineOneEndOnCircumferenceInside.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)
	fmt.Println(
		"Both Ends On Circumference:",
		lineBothEndsOnCircumference.RelationshipToCircle(
			circle, geom2d.WithEpsilon(1e-10),
		).String(),
	)

	// Output:
	// Line Outside: RelationshipLineSegmentCircleMiss
	// Line Tangent: RelationshipLineSegmentCircleTangentToCircle
	// Line Intersecting: RelationshipLineSegmentCircleIntersecting
	// Line Inside: RelationshipLineSegmentCircleContainedByCircle
	// One End On Circumference Outside: RelationshipLineSegmentCircleEndOnCircumferenceOutside
	// One End On Circumference Inside: RelationshipLineSegmentCircleEndOnCircumferenceInside
	// Both Ends On Circumference: RelationshipLineSegmentCircleBothEndsOnCircumference
}

func ExampleLineSegment_RelationshipToLineSegment() {
	// Define a set of line segments to test
	AB := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 0),
	)
	CD := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](11, 0),
		geom2d.NewPoint[float64](15, 0),
	)
	EF := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 10),
	)
	GH := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 10),
		geom2d.NewPoint[float64](10, 0),
	)
	IJ := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](0, 10),
	)

	// Analyze relationships
	fmt.Println("Collinear Disjoint:", AB.RelationshipToLineSegment(CD, geom2d.WithEpsilon(1e-10)).String())
	fmt.Println("Intersects:", EF.RelationshipToLineSegment(GH, geom2d.WithEpsilon(1e-10)).String())
	fmt.Println("AeqC:", AB.RelationshipToLineSegment(IJ, geom2d.WithEpsilon(1e-10)).String())
	fmt.Println("Equal:", AB.RelationshipToLineSegment(AB, geom2d.WithEpsilon(1e-10)).String())
	// Output:
	// Collinear Disjoint: RelationshipLineSegmentLineSegmentCollinearDisjoint
	// Intersects: RelationshipLineSegmentLineSegmentIntersects
	// AeqC: RelationshipLineSegmentLineSegmentAeqC
	// Equal: RelationshipLineSegmentLineSegmentCollinearEqual
}

func ExampleLineSegment_RelationshipToPoint() {
	// Define a line segment AB
	AB := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 0),
	)

	// Define various test points
	startPoint := geom2d.NewPoint[float64](0, 0)
	endPoint := geom2d.NewPoint[float64](10, 0)
	onSegment := geom2d.NewPoint[float64](5, 0)
	onLineNotSegment := geom2d.NewPoint[float64](15, 0)
	missingPoint := geom2d.NewPoint[float64](5, 5)

	// Analyze relationships
	fmt.Println("Point at Start:", AB.RelationshipToPoint(startPoint).String())
	fmt.Println("Point at End:", AB.RelationshipToPoint(endPoint).String())
	fmt.Println("Point on Segment:", AB.RelationshipToPoint(onSegment).String())
	fmt.Println("Point on Infinite Line, Not Segment:", AB.RelationshipToPoint(onLineNotSegment).String())
	fmt.Println("Point Misses Segment:", AB.RelationshipToPoint(missingPoint).String())
	// Output:
	// Point at Start: RelationshipPointLineSegmentPointEqStart
	// Point at End: RelationshipPointLineSegmentPointEqEnd
	// Point on Segment: RelationshipPointLineSegmentPointOnLineSegment
	// Point on Infinite Line, Not Segment: RelationshipPointLineSegmentCollinearDisjoint
	// Point Misses Segment: RelationshipPointLineSegmentMiss
}

func ExampleLineSegment_RelationshipToRectangle() {
	// Define a rectangle with corners (0, 0), (10, 0), (10, 10), (0, 10)
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	})

	// Define various line segments to test
	lineOutside := geom2d.NewLineSegment(
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(20, 20),
	)
	lineInside := geom2d.NewLineSegment(
		geom2d.NewPoint(2, 2),
		geom2d.NewPoint(8, 8),
	)
	lineIntersects := geom2d.NewLineSegment(
		geom2d.NewPoint(5, -5),
		geom2d.NewPoint(5, 15),
	)
	lineTouchesVertex := geom2d.NewLineSegment(
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(-5, 15),
	)

	// Analyze relationships
	fmt.Println("Line Outside:", lineOutside.RelationshipToRectangle(rect).String())
	fmt.Println("Line Inside:", lineInside.RelationshipToRectangle(rect).String())
	fmt.Println("Line Intersects:", lineIntersects.RelationshipToRectangle(rect).String())
	fmt.Println("Line Touches Vertex:", lineTouchesVertex.RelationshipToRectangle(rect).String())
	// Output:
	// Line Outside: RelationshipLineSegmentRectangleMiss
	// Line Inside: RelationshipLineSegmentRectangleContainedByRectangle
	// Line Intersects: RelationshipLineSegmentRectangleEntersAndExits
	// Line Touches Vertex: RelationshipLineSegmentRectangleEndTouchesVertexExternally
}

func ExampleLineSegment_Rotate() {
	// Define a line segment from (2, 3) to (4, 6)
	line := geom2d.NewLineSegment(
		geom2d.NewPoint[int](2, 3),
		geom2d.NewPoint[int](4, 6),
	)

	// Define a pivot point at the origin (0, 0)
	pivot := geom2d.NewPoint[int](0, 0)

	// Rotate the line segment by 90 degrees (π/2 radians) counterclockwise
	rotatedLine := line.Rotate(pivot, math.Pi/2)

	// Print the rotated line segment's start and end points
	fmt.Printf("Rotated Line Start: %v\n", rotatedLine.Points()[0])
	fmt.Printf("Rotated Line End: %v\n", rotatedLine.Points()[1])
	// Output:
	// Rotated Line Start: Point[(-3, 2)]
	// Rotated Line End: Point[(-6, 4)]
}

func ExampleLineSegment_Scale() {
	// Define a line segment from (2, 3) to (4, 6)
	line := geom2d.NewLineSegment(
		geom2d.NewPoint[int](2, 3),
		geom2d.NewPoint[int](4, 6),
	)

	// Define a reference point for scaling
	ref := geom2d.NewPoint[int](0, 0)

	// Scale the line segment by a factor of 2 relative to the origin (0, 0)
	scaledLine := line.Scale(ref, 2)

	// Print the scaled line segment's start and end points
	fmt.Printf("Scaled Line Start: %v\n", scaledLine.Points()[0])
	fmt.Printf("Scaled Line End: %v\n", scaledLine.Points()[1])

	// Scale the line segment by a shrinking factor of 0.5, converting to floating-point type
	lineFloat := line.AsFloat()
	shrunkLine := lineFloat.Scale(ref.AsFloat(), 0.5)

	// Print the shrunk line segment's start and end points
	fmt.Printf("Shrunk Line Start: %v\n", shrunkLine.Points()[0])
	fmt.Printf("Shrunk Line End: %v\n", shrunkLine.Points()[1])

	// Scale the line segment relative to a custom point (3, 3)
	customRef := geom2d.NewPoint[int](3, 3)
	customScaledLine := line.Scale(customRef, 2)

	// Print the line segment scaled relative to the custom reference point
	fmt.Printf("Custom Scaled Line Start: %v\n", customScaledLine.Points()[0])
	fmt.Printf("Custom Scaled Line End: %v\n", customScaledLine.Points()[1])
	// Output:
	// Scaled Line Start: Point[(4, 6)]
	// Scaled Line End: Point[(8, 12)]
	// Shrunk Line Start: Point[(1, 1.5)]
	// Shrunk Line End: Point[(2, 3)]
	// Custom Scaled Line Start: Point[(1, 3)]
	// Custom Scaled Line End: Point[(5, 9)]
}

func ExampleLineSegment_Start() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(1, 2), geom2d.NewPoint(3, 4))
	fmt.Println(lineSegment.Start())
	// Output:
	// Point[(1, 2)]
}

func ExampleLineSegment_String() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(1, 1), geom2d.NewPoint(4, 5))
	fmt.Println(segment.String())
	// Output:
	// LineSegment[(1, 1) -> (4, 5)]
}

func ExampleLineSegment_SubLineSegment() {
	// Define two line segments
	AB := geom2d.NewLineSegment(
		geom2d.NewPoint[int](10, 10),
		geom2d.NewPoint[int](20, 20),
	)
	CD := geom2d.NewLineSegment(
		geom2d.NewPoint[int](5, 5),
		geom2d.NewPoint[int](15, 15),
	)

	// Subtract CD from AB
	result := AB.SubLineSegment(CD)

	// Print the resulting line segment's start and end points
	fmt.Printf("Resulting Line Start: %v\n", result.Points()[0])
	fmt.Printf("Resulting Line End: %v\n", result.Points()[1])
	// Output:
	// Resulting Line Start: Point[(5, 5)]
	// Resulting Line End: Point[(5, 5)]
}

func ExampleLineSegment_Translate() {
	// Define a line segment AB
	AB := geom2d.NewLineSegment(
		geom2d.NewPoint(1, 1),
		geom2d.NewPoint(4, 4),
	)

	// Define the translation vector
	delta := geom2d.NewPoint(2, 3)

	// Translate the line segment by the vector
	translatedAB := AB.Translate(delta)

	// Output the translated line segment
	fmt.Println("Translated Line Segment:")
	fmt.Println("Start Point:", translatedAB.Start())
	fmt.Println("End Point:", translatedAB.End())
	// Output:
	// Translated Line Segment:
	// Start Point: Point[(3, 4)]
	// End Point: Point[(6, 7)]
}

func ExampleRelationshipLineSegmentLineSegment_String() {
	rel := geom2d.RelationshipLineSegmentLineSegmentIntersects
	fmt.Println(rel.String())
	// Output:
	// RelationshipLineSegmentLineSegmentIntersects
}
