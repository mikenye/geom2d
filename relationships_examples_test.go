package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleCircleCircleRelationship_String() {
	rel := geom2d.CCRTouchingExternal
	fmt.Println(rel.String())
	// Output:
	// CCRTouchingExternal
}

func ExampleCircleLineSegmentRelationship_String() {
	rel := geom2d.CLRIntersecting
	fmt.Println(rel.String())
	// Output:
	// CLRIntersecting
}

func ExampleCircleRectangleRelationship_String() {
	rel := geom2d.CRRIntersection
	fmt.Println(rel.String())
	// Output:
	// CRRIntersection
}
