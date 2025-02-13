package _old_test

import (
	"fmt"
	"github.com/mikenye/geom2d"
	"math"
)

func ExampleCircle_BoundingBox() {
	// Create a circle centered at (10, 10) with a radius of 5
	circle := geom2d.NewCircle(geom2d.NewPoint(10, 10), 5)

	// Calculate its bounding box
	boundingBox := circle.BoundingBox()

	// Output the bounding box
	fmt.Println("Bounding box:", boundingBox)
	// Output:
	// Bounding box: Rectangle[(5, 5), (15, 5), (15, 15), (5, 15)]
}

func ExampleCircle_RelationshipToCircle() {
	circle1 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 10)
	circle2 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 10)
	circle3 := geom2d.NewCircle(geom2d.NewPoint(3, 0), 5)
	circle4 := geom2d.NewCircle(geom2d.NewPoint(20, 0), 5)

	fmt.Println(circle1.RelationshipToCircle(circle2))
	fmt.Println(circle1.RelationshipToCircle(circle3))
	fmt.Println(circle3.RelationshipToCircle(circle1))
	fmt.Println(circle1.RelationshipToCircle(circle4))
	// Output:
	// RelationshipEqual
	// RelationshipContains
	// RelationshipContainedBy
	// RelationshipDisjoint
}

func ExampleCircle_RelationshipToLineSegment() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 10)

	line1 := geom2d.NewLineSegment(geom2d.NewPoint(15, 0), geom2d.NewPoint(20, 0))  // Outside
	line2 := geom2d.NewLineSegment(geom2d.NewPoint(0, -15), geom2d.NewPoint(0, 15)) // Intersects
	line3 := geom2d.NewLineSegment(geom2d.NewPoint(5, 5), geom2d.NewPoint(-5, -5))  // Fully contained

	fmt.Println(circle.RelationshipToLineSegment(line1))
	fmt.Println(circle.RelationshipToLineSegment(line2))
	fmt.Println(circle.RelationshipToLineSegment(line3))
	// Output:
	// RelationshipDisjoint
	// RelationshipIntersection
	// RelationshipContains
}

func ExampleCircle_RelationshipToPoint() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 10)

	point1 := geom2d.NewPoint(15, 0) // Outside the circle
	point2 := geom2d.NewPoint(10, 0) // On the circle boundary
	point3 := geom2d.NewPoint(5, 0)  // Inside the circle

	fmt.Println(circle.RelationshipToPoint(point1))
	fmt.Println(circle.RelationshipToPoint(point2))
	fmt.Println(circle.RelationshipToPoint(point3))
	// Output:
	// RelationshipDisjoint
	// RelationshipIntersection
	// RelationshipContains
}

func ExampleCircle_RelationshipToPolyTree() {
	// Create a polygon as a PolyTree
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Create a circle
	circle := geom2d.NewCircle(geom2d.NewPoint(50, 50), 10)

	// Check relationships
	relationships := circle.RelationshipToPolyTree(root)

	// Output the relationship
	fmt.Printf("Circle's relationship to root polygon is: %v\n", relationships[root])
	// Output:
	// Circle's relationship to root polygon is: RelationshipContainedBy
}

func ExampleCircle_RelationshipToRectangle() {
	// Define a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	})

	// Define a circle
	circle := geom2d.NewCircle(geom2d.NewPoint(50, 50), 30)

	// Determine the relationship
	rel := circle.RelationshipToRectangle(rect)

	// Print the result
	fmt.Println("Relationship:", rel)
	// Output:
	// Relationship: RelationshipContainedBy
}
