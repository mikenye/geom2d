package geom2d_test

import (
	"fmt"
	"geom2d"
	"image"
)

func ExampleNewRectangle() {
	// Define four points that form an axis-aligned rectangle
	points := []geom2d.Point[int]{
		geom2d.NewPoint(1, 1), // bottom-left
		geom2d.NewPoint(4, 1), // bottom-right
		geom2d.NewPoint(1, 3), // top-left
		geom2d.NewPoint(4, 3), // top-right
	}

	// Create the rectangle
	rect := geom2d.NewRectangle(points)

	// Print the corners of the rectangle
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Point[(1, 3)]
	// Point[(4, 3)]
	// Point[(4, 1)]
	// Point[(1, 1)]
}

func ExampleNewRectangleFromImageRect() {

	// Define an image.Rectangle
	imgRect := image.Rect(10, 20, 50, 80)

	// Create a Rectangle[int] from the image.Rectangle
	rect := geom2d.NewRectangleFromImageRect(imgRect)

	// Print the corners of the rectangle
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Point[(10, 80)]
	// Point[(50, 80)]
	// Point[(50, 20)]
	// Point[(10, 20)]
}

func ExampleRectangle_Area() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 5),
		geom2d.NewPoint(0, 5),
	})

	// Calculate the area
	area := rect.Area()

	// Print the area
	fmt.Printf("The area of the rectangle is: %d\n", area)

	// Output:
	// The area of the rectangle is: 50
}

func ExampleRectangle_AsFloat32() {
	// Create an integer-based rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 5),
		geom2d.NewPoint(0, 5),
	})

	// Convert the rectangle to a float32-based rectangle
	rectFloat32 := rect.AsFloat32()

	// Print the original and converted rectangle
	fmt.Println("Original rectangle (int):")
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}
	fmt.Println("Converted rectangle (float32):")
	for _, point := range rectFloat32.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Original rectangle (int):
	// Point[(0, 5)]
	// Point[(10, 5)]
	// Point[(10, 0)]
	// Point[(0, 0)]
	// Converted rectangle (float32):
	// Point[(0, 5)]
	// Point[(10, 5)]
	// Point[(10, 0)]
	// Point[(0, 0)]
}

func ExampleRectangle_AsFloat64() {
	// Create an integer-based rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 5),
		geom2d.NewPoint(0, 5),
	})

	// Convert the rectangle to a float64-based rectangle
	rectFloat64 := rect.AsFloat64()

	// Print the original and converted rectangle
	fmt.Println("Original rectangle (int):")
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}
	fmt.Println("Converted rectangle (float64):")
	for _, point := range rectFloat64.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Original rectangle (int):
	// Point[(0, 5)]
	// Point[(10, 5)]
	// Point[(10, 0)]
	// Point[(0, 0)]
	// Converted rectangle (float64):
	// Point[(0, 5)]
	// Point[(10, 5)]
	// Point[(10, 0)]
	// Point[(0, 0)]
}

func ExampleRectangle_AsInt() {
	// Create a float64-based rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[float64]{
		geom2d.NewPoint(0.5, 0.5),
		geom2d.NewPoint(10.5, 0.5),
		geom2d.NewPoint(10.5, 5.5),
		geom2d.NewPoint(0.5, 5.5),
	})

	// Convert the rectangle to an int-based rectangle
	rectInt := rect.AsInt()

	// Print the original and converted rectangle
	fmt.Println("Original rectangle (float64):")
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}
	fmt.Println("Converted rectangle (int):")
	for _, point := range rectInt.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Original rectangle (float64):
	// Point[(0.5, 5.5)]
	// Point[(10.5, 5.5)]
	// Point[(10.5, 0.5)]
	// Point[(0.5, 0.5)]
	// Converted rectangle (int):
	// Point[(0, 5)]
	// Point[(10, 5)]
	// Point[(10, 0)]
	// Point[(0, 0)]
}

func ExampleRectangle_AsIntRounded() {
	// Create a float64-based rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[float64]{
		geom2d.NewPoint(0.5, 0.5),
		geom2d.NewPoint(10.7, 0.5),
		geom2d.NewPoint(10.7, 5.3),
		geom2d.NewPoint(0.5, 5.3),
	})

	// Convert the rectangle to an int-based rectangle with rounding
	rectIntRounded := rect.AsIntRounded()

	// Print the original and converted rectangle
	fmt.Println("Original rectangle (float64):")
	for _, point := range rect.Contour() {
		fmt.Println(point)
	}
	fmt.Println("Converted rectangle (int):")
	for _, point := range rectIntRounded.Contour() {
		fmt.Println(point)
	}

	// Output:
	// Original rectangle (float64):
	// Point[(0.5, 5.3)]
	// Point[(10.7, 5.3)]
	// Point[(10.7, 0.5)]
	// Point[(0.5, 0.5)]
	// Converted rectangle (int):
	// Point[(1, 5)]
	// Point[(11, 5)]
	// Point[(11, 1)]
	// Point[(1, 1)]
}

func ExampleRectangle_ContainsPoint() {
	// Create an integer-based rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(0, 0),
	})

	// Check if the rectangle contains various points
	points := []geom2d.Point[int]{
		geom2d.NewPoint(5, 5),  // Inside
		geom2d.NewPoint(0, 10), // On the top-left corner
		geom2d.NewPoint(10, 0), // On the bottom-right corner
		geom2d.NewPoint(-1, 5), // Outside (left)
		geom2d.NewPoint(5, 11), // Outside (above)
	}

	for _, p := range points {
		fmt.Printf("Point %v contained: %v\n", p, rect.ContainsPoint(p))
	}

	// Output:
	// Point Point[(5, 5)] contained: true
	// Point Point[(0, 10)] contained: true
	// Point Point[(10, 0)] contained: true
	// Point Point[(-1, 5)] contained: false
	// Point Point[(5, 11)] contained: false
}

func ExampleRectangle_Edges() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 5),
		geom2d.NewPoint(0, 5),
	})

	// Get the edges of the rectangle
	edges := rect.Edges()

	// Print the edges
	for i, edge := range edges {
		fmt.Printf("Edge %d: Start %v, End %v\n", i+1, edge.Start(), edge.End())
	}

	// Output:
	// Edge 1: Start Point[(0, 0)], End Point[(10, 0)]
	// Edge 2: Start Point[(10, 0)], End Point[(10, 5)]
	// Edge 3: Start Point[(10, 5)], End Point[(0, 5)]
	// Edge 4: Start Point[(0, 5)], End Point[(0, 0)]
}

func ExampleRectangle_Eq() {
	// Create two identical rectangles
	rect1 := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(0, 0),
	})

	rect2 := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(0, 0),
	})

	// Create a different rectangle
	rect3 := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(1, 11),
		geom2d.NewPoint(11, 11),
		geom2d.NewPoint(11, 1),
		geom2d.NewPoint(1, 1),
	})

	// Compare the rectangles
	fmt.Println("rect1 equals rect2:", rect1.Eq(rect2))
	fmt.Println("rect1 equals rect3:", rect1.Eq(rect3))

	// Output:
	// rect1 equals rect2: true
	// rect1 equals rect3: false
}

func ExampleRectangle_Height() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(10, 0),  // bottom-right
		geom2d.NewPoint(0, 0),   // bottom-left
	})

	// Calculate and print the height of the rectangle
	fmt.Println("Height of the rectangle:", rect.Height())

	// Output:
	// Height of the rectangle: 10
}

func ExampleRectangle_Perimeter() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(0, 0),
	})

	// Calculate and print the perimeter of the rectangle
	fmt.Println("Perimeter of the rectangle:", rect.Perimeter())

	// Output:
	// Perimeter of the rectangle: 40
}

func ExampleRectangle_Contour() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(10, 0),  // bottom-right
		geom2d.NewPoint(0, 0),   // bottom-left
	})

	// Get the points of the rectangle
	points := rect.Contour()

	// Print the points
	fmt.Println("Contour of the rectangle:")
	for _, point := range points {
		fmt.Println(point)
	}

	// Output:
	// Contour of the rectangle:
	// Point[(0, 10)]
	// Point[(10, 10)]
	// Point[(10, 0)]
	// Point[(0, 0)]
}

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

func ExampleRectangle_Scale() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(0, 0),   // bottom-left
		geom2d.NewPoint(10, 0),  // bottom-right
	})

	// Scale the rectangle by a factor of 2 relative to the origin
	scaled := rect.Scale(geom2d.NewPoint(0, 0), 2)

	// Print the original and scaled rectangles
	fmt.Println("Original Rectangle Contour:", rect.Contour())
	fmt.Println("Scaled Rectangle Contour:", scaled.Contour())

	// Output:
	// Original Rectangle Contour: [Point[(0, 10)] Point[(10, 10)] Point[(10, 0)] Point[(0, 0)]]
	// Scaled Rectangle Contour: [Point[(0, 20)] Point[(20, 20)] Point[(20, 0)] Point[(0, 0)]]
}

func ExampleRectangle_ScaleHeight() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(0, 0),   // bottom-left
		geom2d.NewPoint(10, 0),  // bottom-right
	})

	// Scale the height of the rectangle by a factor of 1.5
	scaled := rect.ScaleHeight(1.5)

	// Print the original and scaled rectangles
	fmt.Println("Original Rectangle Contour:", rect.Contour())
	fmt.Println("Scaled Rectangle Contour:", scaled.Contour())

	// Output:
	// Original Rectangle Contour: [Point[(0, 10)] Point[(10, 10)] Point[(10, 0)] Point[(0, 0)]]
	// Scaled Rectangle Contour: [Point[(0, 25)] Point[(10, 25)] Point[(10, 10)] Point[(0, 10)]]
}

func ExampleRectangle_ScaleWidth() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(0, 0),   // bottom-left
		geom2d.NewPoint(10, 0),  // bottom-right
	})

	// Scale the width of the rectangle by a factor of 2.0
	scaled := rect.ScaleWidth(2.0)

	// Print the original and scaled rectangles
	fmt.Println("Original Rectangle Contour:", rect.Contour())
	fmt.Println("Scaled Rectangle Contour:", scaled.Contour())

	// Output:
	// Original Rectangle Contour: [Point[(0, 10)] Point[(10, 10)] Point[(10, 0)] Point[(0, 0)]]
	// Scaled Rectangle Contour: [Point[(0, 10)] Point[(20, 10)] Point[(20, 0)] Point[(0, 0)]]
}

func ExampleRectangle_String() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
	})

	// Get the string representation
	fmt.Println(rect.String())

	// Output:
	// Rectangle[(0, 0), (10, 0), (10, 10), (0, 10)]
}

func ExampleRectangle_Translate() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),  // top-left
		geom2d.NewPoint(10, 10), // top-right
		geom2d.NewPoint(0, 0),   // bottom-left
		geom2d.NewPoint(10, 0),  // bottom-right
	})

	// Translate the rectangle by (5, -5)
	translatedRect := rect.Translate(geom2d.NewPoint(5, -5))

	// Print the translated rectangle
	fmt.Println(translatedRect.String())

	// Output:
	// Rectangle[(5, -5), (15, -5), (15, 5), (5, 5)]
}

func ExampleRectangle_Width() {
	// Create a rectangle
	rect := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(20, 10),
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(20, 0),
	})

	// Calculate the width of the rectangle
	width := rect.Width()

	// Print the width
	fmt.Printf("The width of the rectangle is: %d\n", width)

	// Output:
	// The width of the rectangle is: 20
}
