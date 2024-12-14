package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleRectangle_RelationshipToCircle() {
	// Define a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint[int](0, 0),
		geom2d.NewPoint[int](100, 0),
		geom2d.NewPoint[int](100, 100),
		geom2d.NewPoint[int](0, 100),
	})

	// Define circles to test
	circleInside := geom2d.NewCircle(geom2d.NewPoint[int](50, 50), 10)
	circleIntersecting := geom2d.NewCircle(geom2d.NewPoint[int](50, 50), 60)
	circleOutside := geom2d.NewCircle(geom2d.NewPoint[int](200, 200), 20)

	// Check relationships
	fmt.Println("Circle inside:", rect.RelationshipToCircle(circleInside))
	fmt.Println("Circle intersecting:", rect.RelationshipToCircle(circleIntersecting))
	fmt.Println("Circle outside:", rect.RelationshipToCircle(circleOutside))
	// Output:
	// Circle inside: RelationshipContains
	// Circle intersecting: RelationshipIntersection
	// Circle outside: RelationshipDisjoint
}

func ExampleRectangle_RelationshipToLineSegment() {
	// Define a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint[int](0, 0),
		geom2d.NewPoint[int](100, 0),
		geom2d.NewPoint[int](100, 100),
		geom2d.NewPoint[int](0, 100),
	})

	// Define line segments to check against the rectangle
	segmentInside := geom2d.NewLineSegment(geom2d.NewPoint(10, 10), geom2d.NewPoint(90, 90))
	segmentIntersecting := geom2d.NewLineSegment(geom2d.NewPoint(-10, 50), geom2d.NewPoint(110, 50))
	segmentOutside := geom2d.NewLineSegment(geom2d.NewPoint(200, 200), geom2d.NewPoint(300, 300))

	// Check relationships
	fmt.Println("Segment inside:", rect.RelationshipToLineSegment(segmentInside))
	fmt.Println("Segment intersecting:", rect.RelationshipToLineSegment(segmentIntersecting))
	fmt.Println("Segment outside:", rect.RelationshipToLineSegment(segmentOutside))
	// Output:
	// Segment inside: RelationshipContains
	// Segment intersecting: RelationshipIntersection
	// Segment outside: RelationshipDisjoint
}

func ExampleRectangle_RelationshipToPoint() {
	// Define a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	})

	// Define points to check against the rectangle
	pointInside := geom2d.NewPoint(50, 50)
	pointOnEdge := geom2d.NewPoint(0, 50)
	pointOnVertex := geom2d.NewPoint(0, 0)
	pointOutside := geom2d.NewPoint(200, 200)

	// Check relationships
	fmt.Println("Point inside:", rect.RelationshipToPoint(pointInside))
	fmt.Println("Point on edge:", rect.RelationshipToPoint(pointOnEdge))
	fmt.Println("Point on vertex:", rect.RelationshipToPoint(pointOnVertex))
	fmt.Println("Point outside:", rect.RelationshipToPoint(pointOutside))

	// Output:
	// Point inside: RelationshipContains
	// Point on edge: RelationshipIntersection
	// Point on vertex: RelationshipIntersection
	// Point outside: RelationshipDisjoint
}

func ExampleRectangle_RelationshipToRectangle() {
	// Define two rectangles
	r1 := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	})
	r2 := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(5, 5),
		geom2d.NewPoint(15, 5),
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(5, 15),
	})

	// Determine the relationship
	relationship := r1.RelationshipToRectangle(r2, geom2d.WithEpsilon(1e-10))

	// Output the result
	fmt.Println("Relationship:", relationship)
	// Output:
	// Relationship: RelationshipIntersection
}
