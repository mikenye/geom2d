package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleRelationshipCircleCircle_String() {
	rel := geom2d.RelationshipCircleCircleExternallyTangent
	fmt.Println(rel.String())
	// Output:
	// RelationshipCircleCircleExternallyTangent
}

func ExampleRelationshipLineSegmentCircle_String() {
	rel := geom2d.RelationshipLineSegmentCircleIntersecting
	fmt.Println(rel.String())
	// Output:
	// RelationshipLineSegmentCircleIntersecting
}

func ExampleRelationshipRectangleCircle_String() {
	rel := geom2d.RelationshipRectangleCircleIntersection
	fmt.Println(rel.String())
	// Output:
	// RelationshipRectangleCircleIntersection
}
