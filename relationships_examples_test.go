package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleCircleCircleRelationship_String() {
	rel := geom2d.RelationshipCircleCircleExternallyTangent
	fmt.Println(rel.String())
	// Output:
	// RelationshipCircleCircleExternallyTangent
}

func ExampleCircleLineSegmentRelationship_String() {
	rel := geom2d.RelationshipLineSegmentCircleIntersecting
	fmt.Println(rel.String())
	// Output:
	// RelationshipLineSegmentCircleIntersecting
}

func ExampleCircleRectangleRelationship_String() {
	rel := geom2d.RelationshipRectangleCircleIntersection
	fmt.Println(rel.String())
	// Output:
	// RelationshipRectangleCircleIntersection
}
