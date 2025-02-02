package linesegment_test

import (
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"math"
)

func ExampleNew() {
	segment := linesegment.New(0, 0, 3, 4)
	fmt.Printf("LineSegment: %s\n", segment)

	// Output:
	// LineSegment: (0,0)(3,4)
}

func ExampleNewFromPoints() {
	start := point.New(0, 0)
	end := point.New(3, 4)

	segment := linesegment.NewFromPoints(start, end)

	fmt.Printf("LineSegment: %s\n", segment)

	// Output:
	// LineSegment: (0,0)(3,4)
}

func ExampleFindIntersectionsSlow() {
	// define octothorpe ("#") shape
	octothorpe := []linesegment.LineSegment[int]{
		// Horizontal lines
		linesegment.New[int](0, 7, 10, 7),
		linesegment.New[int](0, 3, 10, 3),
		// Vertical lines
		linesegment.New[int](3, 10, 3, 0),
		linesegment.New[int](7, 10, 7, 0),
	}

	// find intersections via brute force method
	intersections := linesegment.FindIntersectionsSlow(octothorpe, options.WithEpsilon(1e-8))

	fmt.Print("Octothorpe shape with the following line segments: ")
	first := true
	for _, seg := range octothorpe {
		if first {
			fmt.Print(seg)
			first = false
			continue
		}
		fmt.Printf(", %s", seg)
	}
	fmt.Printf("\n...containts the following intersections:\n")
	for _, intersection := range intersections {
		fmt.Println(intersection)
	}

	// Output:
	// Octothorpe shape with the following line segments: (0,7)(10,7), (0,3)(10,3), (3,10)(3,0), (7,10)(7,0)
	// ...containts the following intersections:
	// Intersection type: IntersectionPoint: (3,7) from segments: (0,7)(10,7), (3,10)(3,0)
	// Intersection type: IntersectionPoint: (7,7) from segments: (0,7)(10,7), (7,10)(7,0)
	// Intersection type: IntersectionPoint: (3,3) from segments: (0,3)(10,3), (3,10)(3,0)
	// Intersection type: IntersectionPoint: (7,3) from segments: (0,3)(10,3), (7,10)(7,0)
}

func ExampleFindIntersectionsFast() {
	// define octothorpe ("#") shape
	octothorpe := []linesegment.LineSegment[int]{
		// Horizontal lines
		linesegment.New[int](0, 7, 10, 7),
		linesegment.New[int](0, 3, 10, 3),
		// Vertical lines
		linesegment.New[int](3, 10, 3, 0),
		linesegment.New[int](7, 10, 7, 0),
	}

	// find intersections via sweep line method
	intersections := linesegment.FindIntersectionsFast(octothorpe, options.WithEpsilon(1e-8))

	fmt.Print("Octothorpe shape with the following line segments: ")
	first := true
	for _, seg := range octothorpe {
		if first {
			fmt.Print(seg)
			first = false
			continue
		}
		fmt.Printf(", %s", seg)
	}
	fmt.Printf("\n...containts the following intersections:\n")
	for _, intersection := range intersections {
		fmt.Println(intersection)
	}

	// Output:
	// Octothorpe shape with the following line segments: (0,7)(10,7), (0,3)(10,3), (3,10)(3,0), (7,10)(7,0)
	// ...containts the following intersections:
	// Intersection type: IntersectionPoint: (3,7) from segments: (0,7)(10,7), (3,10)(3,0)
	// Intersection type: IntersectionPoint: (7,7) from segments: (0,7)(10,7), (7,10)(7,0)
	// Intersection type: IntersectionPoint: (3,3) from segments: (0,3)(10,3), (3,10)(3,0)
	// Intersection type: IntersectionPoint: (7,3) from segments: (0,3)(10,3), (7,10)(7,0)
}

func ExampleLineSegment_AsFloat32() {
	intSegment := linesegment.New(3, 4, 5, 6)
	fltSegment := intSegment.AsFloat32()
	fmt.Printf("intSegment is %v of type: %T\n", intSegment, intSegment)
	fmt.Printf("fltSegment is %v of type: %T\n", fltSegment, fltSegment)
	// Output:
	// intSegment is (3,4)(5,6) of type: linesegment.LineSegment[int]
	// fltSegment is (3,4)(5,6) of type: linesegment.LineSegment[float32]
}

func ExampleLineSegment_AsFloat64() {
	intSegment := linesegment.New(3, 4, 5, 6)
	fltSegment := intSegment.AsFloat64()
	fmt.Printf("intSegment is %v of type: %T\n", intSegment, intSegment)
	fmt.Printf("fltSegment is %v of type: %T\n", fltSegment, fltSegment)
	// Output:
	// intSegment is (3,4)(5,6) of type: linesegment.LineSegment[int]
	// fltSegment is (3,4)(5,6) of type: linesegment.LineSegment[float64]
}

func ExampleLineSegment_AsInt() {
	fltSegment := linesegment.New(3.7, 4.1, 5.6, 6.2)
	intSegment := fltSegment.AsInt()
	fmt.Printf("fltSegment is %v of type: %T\n", fltSegment, fltSegment)
	fmt.Printf("intSegment is %v of type: %T\n", intSegment, intSegment)
	// Output:
	// fltSegment is (3.7,4.1)(5.6,6.2) of type: linesegment.LineSegment[float64]
	// intSegment is (3,4)(5,6) of type: linesegment.LineSegment[int]
}

func ExampleLineSegment_AsIntRounded() {
	fltSegment := linesegment.New(3.7, 4.1, 5.6, 6.2)
	intSegment := fltSegment.AsIntRounded()
	fmt.Printf("fltSegment is %v of type: %T\n", fltSegment, fltSegment)
	fmt.Printf("intSegment is %v of type: %T\n", intSegment, intSegment)
	// Output:
	// fltSegment is (3.7,4.1)(5.6,6.2) of type: linesegment.LineSegment[float64]
	// intSegment is (4,4)(6,6) of type: linesegment.LineSegment[int]
}

func ExampleLineSegment_Bresenham() {
	segment := linesegment.New(0, 0, 4, 5)

	// Generate the points on the circle's perimeter.
	fmt.Printf("Integer points on line segment %s:\n", segment)
	for p := range segment.Bresenham {
		fmt.Printf("  point: %s\n", p)
	}

	// Output:
	// Integer points on line segment (0,0)(4,5):
	//   point: (0,0)
	//   point: (1,1)
	//   point: (2,2)
	//   point: (2,3)
	//   point: (3,4)
	//   point: (4,5)
}

func ExampleLineSegment_Center() {
	segment := linesegment.New(0, 0, 10, 10)
	center := segment.Center()
	fmt.Printf("The center of line segment %s is %s\n", segment, center)
	// Output:
	// The center of line segment (0,0)(10,10) is (5,5)
}

func ExampleLineSegment_ContainsPoint() {
	segment := linesegment.New(0, 0, 10, 10)

	pointOnSegment := point.New(5, 5)
	pointOffSegment := point.New(7, 8)

	fmt.Printf("Is point %s on line segment %s: %t\n", pointOnSegment, segment, segment.ContainsPoint(pointOnSegment))
	fmt.Printf("Is point %s on line segment %s: %t\n", pointOffSegment, segment, segment.ContainsPoint(pointOffSegment))
	// Output:
	// Is point (5,5) on line segment (0,0)(10,10): true
	// Is point (7,8) on line segment (0,0)(10,10): false
}

func ExampleLineSegment_DistanceToLineSegment() {
	segmentA := linesegment.New(2, 0, 2, 10)
	segmentB := linesegment.New(8, 0, 8, 10)

	distance := segmentA.DistanceToLineSegment(segmentB)

	fmt.Printf("The Euclidean distance between line segments %s and %s is %0.2f", segmentA, segmentB, distance)
	// Output:
	// The Euclidean distance between line segments (2,0)(2,10) and (8,0)(8,10) is 6.00
}

func ExampleLineSegment_DistanceToPoint() {
	segment := linesegment.New(2, 0, 2, 10)
	p := point.New(5, 5)

	distance := segment.DistanceToPoint(p)

	fmt.Printf("The Euclidean distance between line segment %s and point %s is %0.2f", segment, p, distance)
	// Output:
	// The Euclidean distance between line segment (2,0)(2,10) and point (5,5) is 3.00
}

func ExampleLineSegment_End() {
	segment := linesegment.New(1, 2, 3, 4)
	fmt.Printf("The end point of line segment %s is %s", segment, segment.End())
	// Output:
	// The end point of line segment (1,2)(3,4) is (3,4)
}

func ExampleLineSegment_Eq() {
	s1 := linesegment.New(1.0, 1.0, 4.0, 5.0)
	s2 := linesegment.New(1.0, 1.0, 4.0, 5.0)
	s3 := linesegment.New(2.0, 3.0, 4.0, 5.0)

	fmt.Printf("Is line segment s1 %s equal to line segment s2 %s: %t\n", s1, s2, s1.Eq(s2))
	fmt.Printf("Is line segment s1 %s equal to line segment s3 %s: %t\n", s1, s3, s1.Eq(s3))

	// Output:
	// Is line segment s1 (1,1)(4,5) equal to line segment s2 (1,1)(4,5): true
	// Is line segment s1 (1,1)(4,5) equal to line segment s3 (2,3)(4,5): false
}

func ExampleLineSegment_Eq_epsilon() {
	s1 := linesegment.New[float64](1, 1, 4, 5)
	s2 := linesegment.New(1.0000001, 1.0000001, 4.0000001, 5.0000001)
	epsilon := 1e-6

	fmt.Printf(
		"Is line segment s1 %s equal to line segment s2 %s without epsilon: %t\n",
		s1,
		s2,
		s1.Eq(s2),
	)

	fmt.Printf(
		"Is line segment s1 %s equal to line segment s2 %s with an epsilon of %.0e: %t\n",
		s1,
		s2,
		epsilon,
		s1.Eq(s2, options.WithEpsilon(epsilon)),
	)

	// Output:
	// Is line segment s1 (1,1)(4,5) equal to line segment s2 (1.0000001,1.0000001)(4.0000001,5.0000001) without epsilon: false
	// Is line segment s1 (1,1)(4,5) equal to line segment s2 (1.0000001,1.0000001)(4.0000001,5.0000001) with an epsilon of 1e-06: true
}

func ExampleLineSegment_Flip() {
	segOrig := linesegment.New(0, 0, 10, 10)
	segFlip := segOrig.Flip()
	fmt.Printf("Original orientation: %s\n", segOrig)
	fmt.Printf("Flipped orientation:  %s\n", segFlip)

	// Output:
	// Original orientation: (0,0)(10,10)
	// Flipped orientation:  (10,10)(0,0)
}

func ExampleLineSegment_Intersection() {
	s1 := linesegment.New(0, 0, 10, 10)
	s2 := linesegment.New(0, 10, 10, 0)
	s3 := linesegment.New(5, 5, 15, 15)
	s4 := linesegment.New(2, 0, 8, 0)

	fmt.Printf("s1 intersection with s2:\n%s\n", s1.Intersection(s2))
	fmt.Printf("s1 overlap with s3:\n%s\n", s1.Intersection(s3))
	fmt.Printf("s1 no intersection or overlap with s4:\n%s\n", s1.Intersection(s4))

	// Output:
	// s1 intersection with s2:
	// Intersection type: IntersectionPoint: (5,5) from segments: (0,0)(10,10), (0,10)(10,0)
	// s1 overlap with s3:
	// Intersection type: IntersectionOverlappingSegment: (5,5)(10,10) from segments: (0,0)(10,10), (5,5)(15,15)
	// s1 no intersection or overlap with s4:
	// Intersection type: IntersectionNone, from segments: (0,0)(10,10), (2,0)(8,0)
}

func ExampleLineSegment_Length() {
	segment := linesegment.New(0, 0, 10, 10)
	fmt.Printf("The segment %s is %0.2f units in length\n", segment, segment.Length())

	// Output:
	// The segment (0,0)(10,10) is 14.14 units in length
}

func ExampleLineSegment_Points() {
	segment := linesegment.New(0, 0, 10, 10)

	start := segment.Start()
	end := segment.End()

	fmt.Printf("The segment %s has start point %s and end point %s\n", segment, start, end)

	// Output:
	// The segment (0,0)(10,10) has start point (0,0) and end point (10,10)
}

func ExampleLineSegment_ProjectPoint() {
	segment := linesegment.New(0, 0, 10, 10)
	p := point.New(10, 0)
	closest := segment.ProjectPoint(p)

	fmt.Printf("The closest point to %s on line segment %s is point %s\n", p, segment, closest)

	// Output:
	// The closest point to (10,0) on line segment (0,0)(10,10) is point (5,5)
}

func ExampleLineSegment_ReflectLineSegment() {
	segment := linesegment.New(0, 0, 10, 10)       // line segment to reflect
	xAxis := linesegment.New(0, 0, 1, 0)           // line segment representing x-asix
	reflected := xAxis.ReflectLineSegment(segment) // reflected line segment

	fmt.Printf("Line segment %s reflected across the x-axis is: %s\n", segment, reflected)

	// Output:
	// Line segment (0,0)(10,10) reflected across the x-axis is: (0,0)(10,-10)
}

func ExampleLineSegment_ReflectPoint() {
	p := point.New(10, 10)               // point to reflect
	xAxis := linesegment.New(0, 0, 1, 0) // line segment representing x-asix
	reflected := xAxis.ReflectPoint(p)   // reflected point

	fmt.Printf("Point %s reflected across the x-axis is: %s\n", p, reflected)

	// Output:
	// Point (10,10) reflected across the x-axis is: (10,-10)
}

func ExampleLineSegment_Rotate() {
	segment := linesegment.New(0, 0, 10, 10)

	pivot := point.New(0, 0)
	angle := math.Pi / 2 // Rotate 90 degrees

	rotated := segment.Rotate(pivot, angle)

	fmt.Printf("Line segment %s rotated 90 degrees counter-clockwise around pivot point %s is %s", segment, pivot, rotated)

	// Output:
	// Line segment (0,0)(10,10) rotated 90 degrees counter-clockwise around pivot point (0,0) is (0,0)(-10,10)
}

func ExampleLineSegment_Scale() {
	segment := linesegment.New(0, 0, 10, 10)

	ref := point.New(0, 0)
	scaleFactor := 2

	scaled := segment.Scale(ref, scaleFactor)

	fmt.Printf("Line segment %s scaled by a factor of %v from reference point %s is %s", segment, scaleFactor, ref, scaled)

	// Output:
	// Line segment (0,0)(10,10) scaled by a factor of 2 from reference point (0,0) is (0,0)(20,20)
}

func ExampleLineSegment_Slope() {
	vert := linesegment.New(0, 0, 0, 10)
	horz := linesegment.New(0, 0, 10, 0)
	pos := linesegment.New(0, 0, 10, 10)
	neg := linesegment.New(0, 0, 10, -10)

	fmt.Printf("Slope of vertical line segmemt %s is %f\n", vert, vert.Slope())
	fmt.Printf("Slope of horizontal line segmemt %s is %f\n", horz, horz.Slope())
	fmt.Printf("Slope of diagonal line segmemt %s is %f\n", pos, pos.Slope())
	fmt.Printf("Slope of diagonal line segmemt %s is %f\n", neg, neg.Slope())

	// Output:
	// Slope of vertical line segmemt (0,0)(0,10) is NaN
	// Slope of horizontal line segmemt (0,0)(10,0) is 0.000000
	// Slope of diagonal line segmemt (0,0)(10,10) is 1.000000
	// Slope of diagonal line segmemt (0,0)(10,-10) is -1.000000
}

func ExampleLineSegment_Start() {
	segment := linesegment.New(1, 2, 3, 4)
	fmt.Printf("The start point of line segment %s is %s", segment, segment.Start())
	// Output:
	// The start point of line segment (1,2)(3,4) is (1,2)
}

func ExampleLineSegment_String() {
	segment := linesegment.New(1, 2, 3, 4)

	// When fmt.Println is used to print a variable,
	// and that variable implements the Stringer interface (String() string),
	// Go automatically calls the String() method to generate the output,
	// rather than using the default representation of the type.
	// Thus:
	fmt.Println(segment)
	// is the same as:
	fmt.Println(segment.String())

	// Output:
	// (1,2)(3,4)
	// (1,2)(3,4)
}

func ExampleLineSegment_Translate() {
	segment := linesegment.New(1, 2, 3, 4)

	translationVector := point.New(2, 3) // point used as vector

	translatedSegment := segment.Translate(translationVector)

	fmt.Printf("Line segment %s translated by vector %s is %s", segment, translationVector, translatedSegment)

	// Output:
	// Line segment (1,2)(3,4) translated by vector (2,3) is (3,5)(5,7)
}

func ExampleLineSegment_XAtY() {
	segment := linesegment.New(0, 0, 10, 10)
	fmt.Printf("On line segment %s, when y=5.0, x=%0.1f\n", segment, segment.XAtY(5))
	// Output:
	// On line segment (0,0)(10,10), when y=5.0, x=5.0
}

func ExampleLineSegment_YAtX() {
	segment := linesegment.New(0, 0, 10, 10)
	fmt.Printf("On line segment %s, when x=5.0, y=%0.1f\n", segment, segment.YAtX(5))
	// Output:
	// On line segment (0,0)(10,10), when x=5.0, y=5.0
}
