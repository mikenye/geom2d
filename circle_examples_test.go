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

func ExampleCircle_AsFloat() {
	intCircle := geom2d.NewCircle(geom2d.NewPoint(3, 4), 5)
	fltCircle := intCircle.AsFloat()
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

func ExampleCircle_RelationshipToCircle_ccrmiss() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(-6, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(6, 0), 5)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCRMiss)
	// Output:
	// true
}

func ExampleCircle_RelationshipToCircle_ccrtouchingexternal() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(-5, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(5, 0), 5)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCRTouchingExternal)
	// Output:
	// true
}

func ExampleCircle_RelationshipToCircle_ccroverlapping() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(-4, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(4, 0), 5)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCROverlapping)
	// Output:
	// true
}

func ExampleCircle_RelationshipToCircle_ccrtouchinginternal() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(3, 0), 2)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCRTouchingInternal)
	// Output:
	// true
}

func ExampleCircle_RelationshipToCircle_ccrcontained() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 2)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCRContained)
	// Output:
	// true
}

func ExampleCircle_RelationshipToCircle_ccrequal() {
	c1 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	c2 := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	rel := c1.RelationshipToCircle(c2, geom2d.WithEpsilon(1e-10))
	fmt.Println(rel == geom2d.CCREqual)
	// Output:
	// true
}

func ExampleCircle_RelationshipToLineSegment_clroutside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(6, -6), geom2d.NewPoint(6, 6))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLROutside
}

func ExampleCircle_RelationshipToLineSegment_clrinside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(0, 4))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLRInside
}

func ExampleCircle_RelationshipToLineSegment_clrintersecting() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(0, 10))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLRIntersecting
}

func ExampleCircle_RelationshipToLineSegment_clrtangent() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(5, -6), geom2d.NewPoint(5, 6))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLRTangent
}

func ExampleCircle_RelationshipToLineSegment_clroneendoncircumferenceoutside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(-10, 0), geom2d.NewPoint(-5, 0))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLROneEndOnCircumferenceOutside
}

func ExampleCircle_RelationshipToLineSegment_clroneendoncircumferenceinside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(-5, 0))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLROneEndOnCircumferenceInside
}

func ExampleCircle_RelationshipToLineSegment_clrbothendsoncircumference() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	lineSeg := geom2d.NewLineSegment(geom2d.NewPoint(-4, 3), geom2d.NewPoint(4, -3))
	relationship := circle.RelationshipToLineSegment(lineSeg, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// CLRBothEndsOnCircumference
}

func ExampleCircle_RelationshipToPoint_pcroutside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	point := geom2d.NewPoint(4, 4)
	relationship := circle.RelationshipToPoint(point, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// PCROutside
}

func ExampleCircle_RelationshipToPoint_pcroncircumference() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	point := geom2d.NewPoint(0, -5)
	relationship := circle.RelationshipToPoint(point, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// PCROnCircumference
}

func ExampleCircle_RelationshipToPoint_pcrinside() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	point := geom2d.NewPoint(2, 3)
	relationship := circle.RelationshipToPoint(point, geom2d.WithEpsilon(1e-10))
	fmt.Println(relationship.String())
	// Output:
	// PCRInside
}

func ExampleCircle_RelationshipToRectangle_crrmiss() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	rectangle := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(-9, -6),
		geom2d.NewPoint(-4, -6),
		geom2d.NewPoint(-4, -4),
		geom2d.NewPoint(-9, -4),
	})
	relationship := circle.RelationshipToRectangle(rectangle)
	fmt.Println(relationship.String())
	// Output:
	// CRRMiss
}

func ExampleCircle_RelationshipToRectangle_crrcircleinrect() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	rectangle := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(-6, -6),
		geom2d.NewPoint(6, -6),
		geom2d.NewPoint(6, 6),
		geom2d.NewPoint(-6, 6),
	})
	relationship := circle.RelationshipToRectangle(rectangle)
	fmt.Println(relationship.String())
	// Output:
	// CRRCircleInRect
}

func ExampleCircle_RelationshipToRectangle_crrrectincircle() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	rectangle := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(-3, -3),
		geom2d.NewPoint(3, -3),
		geom2d.NewPoint(3, 3),
		geom2d.NewPoint(-3, 3),
	})
	relationship := circle.RelationshipToRectangle(rectangle)
	fmt.Println(relationship.String())
	// Output:
	// CRRRectInCircle
}

func ExampleCircle_RelationshipToRectangle_crrintersection() {
	circle := geom2d.NewCircle(geom2d.NewPoint(0, 0), 5)
	rectangle := geom2d.NewRectangle([]geom2d.Point[int]{
		geom2d.NewPoint(-4, -4),
		geom2d.NewPoint(4, -4),
		geom2d.NewPoint(4, 4),
		geom2d.NewPoint(-4, 4),
	})
	relationship := circle.RelationshipToRectangle(rectangle)
	fmt.Println(relationship.String())
	// Output:
	// CRRIntersection
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
