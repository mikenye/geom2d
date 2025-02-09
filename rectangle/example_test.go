package rectangle_test

import (
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/rectangle"
)

func ExampleRectangle_Area() {
	r := rectangle.New(0, 0, 10, 10)
	w := r.Width()
	h := r.Height()
	area := r.Area()

	fmt.Printf("The area of rectangle %s is %d*%d = %d units²", r, w, h, area)

	// Output:
	// The area of rectangle [(0,0),(10,10)] is 10*10 = 100 units²
}

func ExampleRectangle_AsFloat32() {
	intRect := rectangle.New(0, 0, 10, 10)
	fltRect := intRect.AsFloat32()
	fmt.Printf("intRect is %v of type: %T\n", intRect, intRect)
	fmt.Printf("fltRect is %v of type: %T\n", fltRect, fltRect)
	// Output:
	// intRect is [(0,0),(10,10)] of type: rectangle.Rectangle[int]
	// fltRect is [(0,0),(10,10)] of type: rectangle.Rectangle[float32]
}

func ExampleRectangle_AsFloat64() {
	intRect := linesegment.New(3, 4, 5, 6)
	fltRect := intRect.AsFloat64()
	fmt.Printf("intRect is %v of type: %T\n", intRect, intRect)
	fmt.Printf("fltRect is %v of type: %T\n", fltRect, fltRect)
	// Output:
	// intRect is (3,4)(5,6) of type: linesegment.LineSegment[int]
	// fltRect is (3,4)(5,6) of type: linesegment.LineSegment[float64]
}

func ExampleRectangle_AsInt() {
	fltRect := linesegment.New(3.7, 4.1, 5.6, 6.2)
	intRect := fltRect.AsInt()
	fmt.Printf("fltRect is %v of type: %T\n", fltRect, fltRect)
	fmt.Printf("intRect is %v of type: %T\n", intRect, intRect)
	// Output:
	// fltRect is (3.7,4.1)(5.6,6.2) of type: linesegment.LineSegment[float64]
	// intRect is (3,4)(5,6) of type: linesegment.LineSegment[int]
}

func ExampleRectangle_AsIntRounded() {
	fltRect := linesegment.New(3.7, 4.1, 5.6, 6.2)
	intRect := fltRect.AsIntRounded()
	fmt.Printf("fltRect is %v of type: %T\n", fltRect, fltRect)
	fmt.Printf("intRect is %v of type: %T\n", intRect, intRect)
	// Output:
	// fltRect is (3.7,4.1)(5.6,6.2) of type: linesegment.LineSegment[float64]
	// intRect is (4,4)(6,6) of type: linesegment.LineSegment[int]
}

func ExampleRectangle_ContainsPoint() {
	r := rectangle.New(5, 5, 10, 10)
	p1 := point.New(2, 1)
	p2 := point.New(7, 8)

	fmt.Printf("Is point %s inside rectangle %s: %t\n", p1, r, r.ContainsPoint(p1))
	fmt.Printf("Is point %s inside rectangle %s: %t\n", p2, r, r.ContainsPoint(p2))

	// Output:
	// Is point (2,1) inside rectangle [(5,5),(10,10)]: false
	// Is point (7,8) inside rectangle [(5,5),(10,10)]: true
}

func ExampleRectangle_Contour() {
	r := rectangle.New(5, 5, 10, 10)
	bl, br, tr, tl := r.Contour()
	fmt.Printf("The bottom-left corner of rectangle %s is: %s\n", r, bl)
	fmt.Printf("The bottom-right corner of rectangle %s is: %s\n", r, br)
	fmt.Printf("The top-left corner of rectangle %s is: %s\n", r, tl)
	fmt.Printf("The top-right corner of rectangle %s is: %s\n", r, tr)

	// Output:
	// The bottom-left corner of rectangle [(5,5),(10,10)] is: (5,5)
	// The bottom-right corner of rectangle [(5,5),(10,10)] is: (10,5)
	// The top-left corner of rectangle [(5,5),(10,10)] is: (5,10)
	// The top-right corner of rectangle [(5,5),(10,10)] is: (10,10)

}

func ExampleRectangle_Edges() {
	r := rectangle.New(5, 5, 10, 10)
	bottom, right, top, left := r.Edges()
	fmt.Printf("The bottom edge of rectangle %s is: %s\n", r, bottom)
	fmt.Printf("The right edge of rectangle %s is: %s\n", r, right)
	fmt.Printf("The top edge of rectangle %s is: %s\n", r, top)
	fmt.Printf("The left edge of rectangle %s is: %s\n", r, left)

	// Output:
	// The bottom edge of rectangle [(5,5),(10,10)] is: (5,5)(10,5)
	// The right edge of rectangle [(5,5),(10,10)] is: (10,5)(10,10)
	// The top edge of rectangle [(5,5),(10,10)] is: (10,10)(5,10)
	// The left edge of rectangle [(5,5),(10,10)] is: (5,10)(5,5)
}

func ExampleRectangle_EdgesIter() {
	r := rectangle.New(0, 0, 4, 3)

	fmt.Printf("Edges in rectangle: %s\n", r.String())
	for seg := range r.EdgesIter {
		fmt.Printf("  edge: %s\n", seg.String())
	}

	// Output:
	// Edges in rectangle: [(0,0),(4,3)]
	//   edge: (0,0)(4,0)
	//   edge: (4,0)(4,3)
	//   edge: (4,3)(0,3)
	//   edge: (0,3)(0,0)
}

func ExampleRectangle_Eq() {
	r1 := rectangle.New(0, 0, 10, 10)
	r2 := rectangle.New(0, 0, 10, 10)
	r3 := rectangle.New(1, 1, 11, 11)
	fmt.Printf("Are rectangles %s and %s equal: %t\n", r1, r2, r1.Eq(r2))
	fmt.Printf("Are rectangles %s and %s equal: %t\n", r1, r3, r1.Eq(r3))

	// Output:
	// Are rectangles [(0,0),(10,10)] and [(0,0),(10,10)] equal: true
	// Are rectangles [(0,0),(10,10)] and [(1,1),(11,11)] equal: false
}

func ExampleRectangle_Eq_epsilon() {
	r1 := rectangle.New[float64](0, 0, 10, 10)
	r2 := rectangle.New[float64](0, 0, 10.0000000001, 10.0000000001)
	epsilon := 1e-8
	fmt.Printf("Are rectangles %s and %s equal: %t\n", r1, r2, r1.Eq(r2))
	fmt.Printf(
		"Are rectangles %s and %s equal with epsilon of %0.0e: %t\n",
		r1,
		r2,
		epsilon,
		r1.Eq(r2, options.WithEpsilon(epsilon)),
	)

	// Output:
	// Are rectangles [(0,0),(10,10)] and [(0,0),(10.0000000001,10.0000000001)] equal: false
	// Are rectangles [(0,0),(10,10)] and [(0,0),(10.0000000001,10.0000000001)] equal with epsilon of 1e-08: true
}

func ExampleRectangle_Height() {
	r := rectangle.New(3, 6, 9, 12)
	h := r.Height()
	fmt.Printf("The height of rectangle %s is %d\n", r, h)

	// Output:
	// The height of rectangle [(3,6),(9,12)] is 6
}

func ExampleRectangle_Perimeter() {
	r := rectangle.New(3, 6, 9, 12)
	p := r.Perimeter()
	fmt.Printf("The perimeter of rectangle %s is %d\n", r, p)

	// Output:
	// The perimeter of rectangle [(3,6),(9,12)] is 24
}

func ExampleRectangle_Scale() {
	r := rectangle.New(0, 0, 10, 10)
	ref := point.New(0, 0)
	factor := 2
	scaled := r.Scale(ref, factor)
	fmt.Printf(
		"Rectangle %s scaled by a factor of %d from reference point %s is: %s",
		r,
		factor,
		ref,
		scaled,
	)

	// Output:
	// Rectangle [(0,0),(10,10)] scaled by a factor of 2 from reference point (0,0) is: [(0,0),(20,20)]
}

func ExampleRectangle_ScaleHeight() {
	r := rectangle.New(0, 0, 10, 10)
	factor := 2
	scaled := r.ScaleHeight(factor)
	fmt.Printf(
		"Rectangle %s with height scaled by a factor of %d is: %s",
		r,
		factor,
		scaled,
	)

	// Output:
	// Rectangle [(0,0),(10,10)] with height scaled by a factor of 2 is: [(0,0),(10,20)]
}

func ExampleRectangle_ScaleWidth() {
	r := rectangle.New(0, 0, 10, 10)
	factor := 2
	scaled := r.ScaleWidth(factor)
	fmt.Printf(
		"Rectangle %s with width scaled by a factor of %d is: %s",
		r,
		factor,
		scaled,
	)

	// Output:
	// Rectangle [(0,0),(10,10)] with width scaled by a factor of 2 is: [(0,0),(20,10)]
}

func ExampleRectangle_String() {
	r := rectangle.New(-3, 4, 8, 13)

	// When fmt.Println is used to print a variable,
	// and that variable implements the Stringer interface (String() string),
	// Go automatically calls the String() method to generate the output,
	// rather than using the default representation of the type.
	// Thus:
	fmt.Println(r)
	// is the same as:
	fmt.Println(r.String())

	// Output:
	// [(-3,4),(8,13)]
	// [(-3,4),(8,13)]
}

func ExampleRectangle_ToImageRect() {
	r := rectangle.New(-3, 4, 8, 13)
	ir := r.ToImageRect()

	fmt.Printf("r is %s, type %T\n", r, r)
	fmt.Printf("ir is %s, type %T\n", ir, ir)

	// Output:
	// r is [(-3,4),(8,13)], type rectangle.Rectangle[int]
	// ir is (-3,4)-(8,13), type image.Rectangle
}

func ExampleRectangle_Translate() {
	r := rectangle.New(0, 0, 10, 10)
	p := point.New(5, 5)
	translated := r.Translate(p)

	fmt.Printf("Rectangle %s translated by %s is: %s\n", r, p, translated)

	// Output:
	// Rectangle [(0,0),(10,10)] translated by (5,5) is: [(5,5),(15,15)]
}

func ExampleRectangle_Width() {
	r := rectangle.New(3, 6, 9, 12)
	w := r.Width()
	fmt.Printf("The width of rectangle %s is %d\n", r, w)

	// Output:
	// The width of rectangle [(3,6),(9,12)] is 6
}
