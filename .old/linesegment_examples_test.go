package _old_test

import (
	"fmt"
	"github.com/mikenye/geom2d"
	"math"
)

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

func ExampleLineSegment_IntersectsLineSegment() {
	// Define two line segments
	line1 := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))   // Line segment from (0, 0) to (10, 10)
	line2 := geom2d.NewLineSegment(geom2d.NewPoint(5, 5), geom2d.NewPoint(15, 15))   // Line segment from (5, 5) to (15, 15)
	line3 := geom2d.NewLineSegment(geom2d.NewPoint(20, 20), geom2d.NewPoint(30, 30)) // Line segment far from line1 and line2

	// Check for intersection
	fmt.Printf("Line1 intersects Line2: %v\n", line1.IntersectsLineSegment(line2)) // Expected: true
	fmt.Printf("Line1 intersects Line3: %v\n", line1.IntersectsLineSegment(line3)) // Expected: false

	// Overlapping case
	line4 := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(5, 5))   // Line segment overlaps partially with Line1
	fmt.Printf("Line1 intersects Line4: %v\n", line1.IntersectsLineSegment(line4)) // Expected: true

	// Endpoint touching case
	line5 := geom2d.NewLineSegment(geom2d.NewPoint(10, 10), geom2d.NewPoint(20, 20)) // Shares an endpoint with Line1
	fmt.Printf("Line1 intersects Line5: %v\n", line1.IntersectsLineSegment(line5))   // Expected: true

	// Output:
	// Line1 intersects Line2: true
	// Line1 intersects Line3: false
	// Line1 intersects Line4: true
	// Line1 intersects Line5: true
}

// todo: update example below
//func ExampleLineSegment_IntersectionPoint() {
//	// Define two line segments
//	AB := geom2d.NewLineSegment(geom2d.NewPoint[float64](0, 0), geom2d.NewPoint[float64](4, 4))
//	CD := geom2d.NewLineSegment(geom2d.NewPoint[float64](0, 4), geom2d.NewPoint[float64](4, 0))
//
//	// Find the intersection point
//	intersection, exists := AB.IntersectionPoint(CD)
//
//	// Print the result
//	if exists {
//		fmt.Printf("Intersection point: (%.2f, %.2f)\n", intersection.X(), intersection.Y())
//	} else {
//		fmt.Println("No intersection point exists.")
//	}
//
//	// Output:
//	// Intersection point: (2.00, 2.00)
//}

func ExampleLineSegment_Normalize() {
	// Create a line segment where the start point is not leftmost-lowest
	original := geom2d.NewLineSegment(geom2d.NewPoint(3, 5), geom2d.NewPoint(1, 2))

	// Normalize the segment
	normalized := original.Normalize()

	// Print the result
	fmt.Printf("Original Line Segment: %s\n", original.String())
	fmt.Printf("Normalized Line Segment: %s\n", normalized.String())

	// Output:
	// Original Line Segment: LineSegment[(3, 5) -> (1, 2)]
	// Normalized Line Segment: LineSegment[(1, 2) -> (3, 5)]
}

func ExampleLineSegment_RelationshipToCircle() {
	// Define a circle with center (5, 5) and radius 5
	circle := geom2d.NewCircle(geom2d.NewPoint(5, 5), 5.0)

	// Define various line segments
	lineDisjoint := geom2d.NewLineSegment(
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(-2, -2),
	)
	lineIntersecting := geom2d.NewLineSegment(
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 10),
	)
	lineContained := geom2d.NewLineSegment(
		geom2d.NewPoint(5, 6),
		geom2d.NewPoint(5, 4),
	)

	// Evaluate relationships
	fmt.Println("Disjoint:", lineDisjoint.RelationshipToCircle(circle))
	fmt.Println("Intersecting:", lineIntersecting.RelationshipToCircle(circle))
	fmt.Println("Contained:", lineContained.RelationshipToCircle(circle))
	// Output:
	// Disjoint: RelationshipDisjoint
	// Intersecting: RelationshipIntersection
	// Contained: RelationshipContainedBy
}

func ExampleLineSegment_RelationshipToLineSegment() {
	// Define two line segments
	line1 := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))
	line2 := geom2d.NewLineSegment(geom2d.NewPoint(5, 5), geom2d.NewPoint(15, 15))
	line3 := geom2d.NewLineSegment(geom2d.NewPoint(20, 20), geom2d.NewPoint(30, 30))
	line4 := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))

	// Evaluate relationships
	fmt.Println("Line1 vs Line2:", line1.RelationshipToLineSegment(line2))
	fmt.Println("Line1 vs Line3:", line1.RelationshipToLineSegment(line3))
	fmt.Println("Line1 vs Line4:", line1.RelationshipToLineSegment(line4))
	// Output:
	// Line1 vs Line2: RelationshipIntersection
	// Line1 vs Line3: RelationshipDisjoint
	// Line1 vs Line4: RelationshipEqual
}

func ExampleLineSegment_RelationshipToPolyTree() {
	// Define a PolyTree with a root polygon and a hole
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	}, geom2d.PTSolid)

	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(4, 4),
		geom2d.NewPoint(6, 4),
		geom2d.NewPoint(6, 6),
		geom2d.NewPoint(4, 6),
	}, geom2d.PTHole)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define a line segment
	segment := geom2d.NewLineSegment(
		geom2d.NewPoint(2, 2),
		geom2d.NewPoint(8, 8),
	)

	// Evaluate relationships
	relationships := segment.RelationshipToPolyTree(root)

	// Print results
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContainedBy
	// Hole polygon relationship: RelationshipIntersection
}

func ExampleLineSegment_RelationshipToPoint() {
	// Define a line segment
	segment := geom2d.NewLineSegment(
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10, 10),
	)

	// Define some points
	point1 := geom2d.NewPoint[float64](5, 5)  // On the segment
	point2 := geom2d.NewPoint[float64](10, 0) // Disjoint
	point3 := geom2d.NewPoint[float64](0, 0)  // Coincides with an endpoint

	// Evaluate relationships
	fmt.Println("Point1 vs Line Segment:", segment.RelationshipToPoint(point1))
	fmt.Println("Point2 vs Line Segment:", segment.RelationshipToPoint(point2))
	fmt.Println("Point3 vs Line Segment:", segment.RelationshipToPoint(point3))
	// Output:
	// Point1 vs Line Segment: RelationshipIntersection
	// Point2 vs Line Segment: RelationshipDisjoint
	// Point3 vs Line Segment: RelationshipIntersection
}

func ExampleLineSegment_RelationshipToRectangle() {
	// Define a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	})

	// Define some line segments
	line1 := geom2d.NewLineSegment(geom2d.NewPoint(5, 5), geom2d.NewPoint(15, 15))   // Intersects
	line2 := geom2d.NewLineSegment(geom2d.NewPoint(1, 1), geom2d.NewPoint(9, 9))     // Contained
	line3 := geom2d.NewLineSegment(geom2d.NewPoint(20, 20), geom2d.NewPoint(30, 30)) // Disjoint

	// Evaluate relationships
	fmt.Println("Line1 vs Rectangle:", line1.RelationshipToRectangle(rect))
	fmt.Println("Line2 vs Rectangle:", line2.RelationshipToRectangle(rect))
	fmt.Println("Line3 vs Rectangle:", line3.RelationshipToRectangle(rect))
	// Output:
	// Line1 vs Rectangle: RelationshipIntersection
	// Line2 vs Rectangle: RelationshipContainedBy
	// Line3 vs Rectangle: RelationshipDisjoint
}

// ExampleLineSegment_RoundToEpsilon demonstrates the use of the RoundToEpsilon method
// to round the coordinates of a LineSegment to the nearest multiple of the given epsilon.
func ExampleLineSegment_RoundToEpsilon() {
	// Create a line segment
	ls := geom2d.NewLineSegment(geom2d.NewPoint[float64](1.2345, 4.5678), geom2d.NewPoint[float64](7.8912, 3.2109))

	// Round the coordinates to the nearest 0.1
	rounded := ls.RoundToEpsilon(0.1)

	// Print the rounded line segment
	fmt.Printf("LineSegment[(%.4f, %.4f) -> (%.4f, %.4f)]",
		rounded.Start().X(),
		rounded.Start().Y(),
		rounded.End().X(),
		rounded.End().Y(),
	)
	// Output: LineSegment[(1.2000, 4.6000) -> (7.9000, 3.2000)]
}

// ExampleLineSegment_Slope demonstrates the use of the Slope method
// to calculate the slope of a line segment.
func ExampleLineSegment_Slope() {
	// Create a line segment
	ls := geom2d.NewLineSegment(geom2d.NewPoint[float64](1, 1), geom2d.NewPoint[float64](3, 5))

	// Calculate the slope
	slope, ok := ls.Slope()

	// Print the slope and whether it is valid
	fmt.Printf("Slope: %.6f, Valid: %t\n", slope, ok)
	// Output: Slope: 2.000000, Valid: true
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

func ExampleLineSegment_XAtY() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(2, 2), geom2d.NewPoint(8, 5))

	x, ok := segment.XAtY(3) // Y = 3
	if ok {
		fmt.Printf("X at Y=3 is %.2f\n", x)
	} else {
		fmt.Println("Y=3 is out of bounds for the segment")
	}

	// Output:
	// X at Y=3 is 4.00

}

func ExampleLineSegment_YAtX() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(2, 2), geom2d.NewPoint(8, 5))

	y, ok := segment.YAtX(4) // X = 4
	if ok {
		fmt.Printf("Y at X=4 is %.2f\n", y)
	} else {
		fmt.Println("X=4 is out of bounds for the segment")
	}

	// Output:
	// Y at X=4 is 3.00
}
