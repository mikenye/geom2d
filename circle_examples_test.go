package geom2d_test

import (
	"fmt"
	"geom2d"
	"math"
)

func ExampleNewCircle() {
	center := geom2d.NewPoint(3, 4) // Note type inference as int not explicitly specified
	radius := 5
	circle := geom2d.NewCircle(center, radius) // Creates a Circle with center (3, 4) and radius 5
	fmt.Println(circle)
	// Output:
	// Circle[center=(3, 4), radius=5]
}

func ExampleCircle_Area() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5) // Creates a Circle with center (3, 4) and radius 5
	area := circle.Area()                                // area = π*5*5 = 25π
	areaRounded := math.Round(area*100) / 100            // Round to two digits
	fmt.Println(areaRounded)
	// Output:
	// 78.54
}

func ExampleCircle_AsFloat64() {
	intCircle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	fltCircle := intCircle.AsFloat64()
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	// Output:
	// intCircle is Circle[center=(3, 4), radius=5] of type: geom2d.Circle[int]
	// fltCircle is Circle[center=(3, 4), radius=5] of type: geom2d.Circle[float64]
}

func ExampleCircle_AsInt() {
	fltCircle := geom2d.NewCircle(geom2d.NewPoint(3.7, 4.9), 5.6)
	intCircle := fltCircle.AsInt()
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	// Output:
	// fltCircle is Circle[center=(3.7, 4.9), radius=5.6] of type: geom2d.Circle[float64]
	// intCircle is Circle[center=(3, 4), radius=5] of type: geom2d.Circle[int]
}

func ExampleCircle_AsIntRounded() {
	fltCircle := geom2d.NewCircle(geom2d.NewPoint(3.7, 4.2), 5.6)
	intCircle := fltCircle.AsIntRounded()
	fmt.Printf("fltCircle is %v of type: %T\n", fltCircle, fltCircle)
	fmt.Printf("intCircle is %v of type: %T\n", intCircle, intCircle)
	// Output:
	// fltCircle is Circle[center=(3.7, 4.2), radius=5.6] of type: geom2d.Circle[float64]
	// intCircle is Circle[center=(4, 4), radius=6] of type: geom2d.Circle[int]
}

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

func ExampleCircle_Center() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	center := circle.Center()
	fmt.Println(center)
	// Output:
	// Point[(3, 4)]
}

func ExampleCircle_Circumference() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	circumference := circle.Circumference()
	circumferenceRounded := math.Round(circumference*100) / 100
	fmt.Println(circumferenceRounded)
	// Output:
	// 31.42
}

func ExampleCircle_Eq() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	isEqual := c1.Eq(c2)
	fmt.Println(isEqual)
	// Output:
	// true
}

func ExampleCircle_Eq_epsilon() {
	c1 := geom2d.NewCircle(geom2d.NewPoint[float64](3, 4), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(3.0001, 4.0001), 5.0001)
	isApproximatelyEqual := c1.Eq(c2, geom2d.WithEpsilon(1e-3))
	fmt.Println(isApproximatelyEqual)
	// Output:
	// true
}

func ExampleCircle_Radius() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	radius := circle.Radius()
	fmt.Println(radius)
	// Output:
	// 5
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

func ExampleCircle_Rotate() {
	pivot := geom2d.NewPoint(1, 1)
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 3), 5)
	angle := math.Pi / 2 // Rotate 90 degrees
	rotated := circle.Rotate(pivot, angle, geom2d.WithEpsilon(1e-10)).AsInt()
	fmt.Println(rotated)
	// Output:
	// Circle[center=(-1, 3), radius=5]
}

func ExampleCircle_Scale() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	scaled := circle.Scale(2).AsInt()
	fmt.Println(scaled)
	// Output:
	// Circle[center=(0, 0), radius=10]
}

func ExampleCircle_String() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	fmt.Println(circle.String())
	// Output:
	// Circle[center=(3, 4), radius=5]
}

func ExampleCircle_Translate() {
	circle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	translationVector := geom2d.NewPoint(2, 3)
	translatedCircle := circle.Translate(translationVector)
	fmt.Println(translatedCircle)
	// Output:
	// Circle[center=(5, 7), radius=5]
}
