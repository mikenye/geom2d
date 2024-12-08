package geom2d

// todo for this file:
//   - is underscore ok, eg: "RelationshipLineSegmentCircle_ContainedByCircle"?
//   - possibly shorter names without losing clarity?

import "fmt"

// RelationshipCircleCircle defines the possible spatial relationships
// between two circles in a 2D plane.
//
// This enumeration categorizes how two circles relate to each other based on
// their positions and radii. It covers scenarios such as being disjoint, touching,
// overlapping, or one circle being contained within the other.
type RelationshipCircleCircle uint8

// Valid values for RelationshipCircleCircle
const (
	// RelationshipCircleCircleMiss - circles are disjoint.
	RelationshipCircleCircleMiss RelationshipCircleCircle = iota

	// RelationshipCircleCircleExternallyTangent - circles are externally tangent.
	RelationshipCircleCircleExternallyTangent

	// RelationshipCircleCircleIntersection - circles overlap and intersect at two points.
	RelationshipCircleCircleIntersection

	// RelationshipCircleCircleInternallyTangent - circles are internally tangent.
	RelationshipCircleCircleInternallyTangent

	// RelationshipCircleCircleContained - one circle is fully contained within the other, with no touching/intersection.
	RelationshipCircleCircleContained

	// RelationshipCircleCircleEqual - circles are identical.
	RelationshipCircleCircleEqual
)

// String returns the string representation of a [RelationshipCircleCircle].
//
// The function converts the [RelationshipCircleCircle] enum value into its
// corresponding string representation for readability and debugging purposes.
//
// Panics:
//   - If the [RelationshipCircleCircle] has an unsupported or undefined value, the function will panic with a descriptive error.
//
// Returns:
//   - string: The name of the [RelationshipCircleCircle] enum value.
func (r RelationshipCircleCircle) String() string {
	switch r {
	case RelationshipCircleCircleMiss:
		return "RelationshipCircleCircleMiss"
	case RelationshipCircleCircleExternallyTangent:
		return "RelationshipCircleCircleExternallyTangent"
	case RelationshipCircleCircleIntersection:
		return "RelationshipCircleCircleIntersection"
	case RelationshipCircleCircleInternallyTangent:
		return "RelationshipCircleCircleInternallyTangent"
	case RelationshipCircleCircleContained:
		return "RelationshipCircleCircleContained"
	case RelationshipCircleCircleEqual:
		return "RelationshipCircleCircleEqual"
	default:
		panic(fmt.Errorf("unsupported RelationshipCircleCircle"))
	}
}

// RelationshipLineSegmentCircle defines the possible spatial relationships
// between a [Circle] and a [LineSegment] in a 2D plane.
//
// This enumeration categorizes how a [LineSegment] relates to a [Circle] based on
// its position, intersection, and tangency.
type RelationshipLineSegmentCircle uint8

// Valid values for RelationshipLineSegmentCircle
const (
	// RelationshipLineSegmentCircleMiss indicates that the line segment lies completely outside the circle,
	// with no intersection or tangency.
	RelationshipLineSegmentCircleMiss RelationshipLineSegmentCircle = iota

	// RelationshipLineSegmentCircleContainedByCircle indicates that the line segment lies completely within the circle,
	// with both endpoints inside the circle's boundary.
	RelationshipLineSegmentCircleContainedByCircle

	// RelationshipLineSegmentCircleIntersecting indicates that the line segment intersects the circle at exactly
	// two distinct points.
	RelationshipLineSegmentCircleIntersecting

	// RelationshipLineSegmentCircleTangentToCircle indicates that the line segment is tangent to the circle, touching it
	// at exactly one point with the tangent forming a 90-degree angle to the circle's radius.
	RelationshipLineSegmentCircleTangentToCircle

	// RelationshipLineSegmentCircleEndOnCircumferenceOutside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies outside the circle.
	RelationshipLineSegmentCircleEndOnCircumferenceOutside

	// RelationshipLineSegmentCircleEndOnCircumferenceInside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies inside the circle.
	RelationshipLineSegmentCircleEndOnCircumferenceInside

	// RelationshipLineSegmentCircleBothEndsOnCircumference indicates that both endpoints of the line segment lie
	// exactly on the circle's boundary. The line segment does not extend inside or outside
	// the circle.
	RelationshipLineSegmentCircleBothEndsOnCircumference
)

// String returns the string representation of a [RelationshipLineSegmentCircle].
//
// This function converts the [RelationshipLineSegmentCircle] enum value into a
// corresponding string for improved readability and debugging.
//
// Panics:
//   - If the [RelationshipLineSegmentCircle] has an unsupported or undefined value,
//     the function will panic.
//
// Returns:
//   - string: The name of the [RelationshipLineSegmentCircle] enum value.
func (r RelationshipLineSegmentCircle) String() string {
	switch r {
	case RelationshipLineSegmentCircleMiss:
		return "RelationshipLineSegmentCircleMiss"
	case RelationshipLineSegmentCircleContainedByCircle:
		return "RelationshipLineSegmentCircleContainedByCircle"
	case RelationshipLineSegmentCircleIntersecting:
		return "RelationshipLineSegmentCircleIntersecting"
	case RelationshipLineSegmentCircleTangentToCircle:
		return "RelationshipLineSegmentCircleTangentToCircle"
	case RelationshipLineSegmentCircleEndOnCircumferenceOutside:
		return "RelationshipLineSegmentCircleEndOnCircumferenceOutside"
	case RelationshipLineSegmentCircleEndOnCircumferenceInside:
		return "RelationshipLineSegmentCircleEndOnCircumferenceInside"
	case RelationshipLineSegmentCircleBothEndsOnCircumference:
		return "RelationshipLineSegmentCircleBothEndsOnCircumference"
	default:
		panic(fmt.Errorf("unsupported RelationshipLineSegmentCircle"))
	}
}

// RelationshipCirclePolygon defines the possible spatial relationships
// between a [Circle] and the contour of a single polygon within a given [PolyTree].
//
// This type categorizes how a [Circle] relates to an individual polygon in the [PolyTree].
// The relationships include cases where the circle is inside, outside, tangential to, intersecting,
// containing, or contained by the polygon's contour.
//
// Notes:
//   - This type describes the relationship of the circle to the current polygon only.
//   - This type is intended to be used inside the [RelationshipCirclePolyTree] map.
//
// See Also:
//   - [PolyTree.RelationshipToCircle]: For computing relationships between a circle and all polygons
//     within a [PolyTree].
type RelationshipCirclePolygon uint8

// Valid values for RelationshipCirclePolygon
const (
	// RelationshipCirclePolyTreeMiss indicates that the circle lies entirely outside the polygon.
	RelationshipCirclePolyTreeMiss RelationshipCirclePolygon = iota

	// RelationshipCirclePolyTreeExternallyTouching indicates that the circle touches the polygon non-tangentially
	// at one or more points but does not intersect or overlap with it.
	// The circle and polygon share boundary points without crossing.
	RelationshipCirclePolyTreeExternallyTouching

	// RelationshipCirclePolyTreeExternallyTangent indicates that the circle is tangential to the polygon from the outside.
	// This means the circle touches the polygon at exactly one point and does not intersect or enter the polygon's interior.
	RelationshipCirclePolyTreeExternallyTangent

	// RelationshipCirclePolyTreeIntersection indicates that the circle intersects one or more edges of the polygon.
	// This means the circle crosses the polygon's boundary, entering and/or exiting the polygon.
	RelationshipCirclePolyTreeIntersection

	// RelationshipCirclePolyTreeInternallyTangent indicates that the circle is tangential to the polygon from the inside.
	// This means the circle touches the polygon at exactly one point without extending outside the polygon's boundary.
	RelationshipCirclePolyTreeInternallyTangent

	// RelationshipCirclePolyTreeInternallyTouching indicates that the circle's boundary touches the polygon's boundary
	// non-tangentially from the inside at one or more points but does not fully intersect or exit the polygon's interior.
	RelationshipCirclePolyTreeInternallyTouching

	// RelationshipCirclePolyTreeContainedByCircle indicates that the circle fully contains the polygon.
	// There is no intersection or touching between the circle's boundary and the polygon's boundary.
	RelationshipCirclePolyTreeContainedByCircle

	// RelationshipCirclePolyTreeContainedByPoly indicates that the polygon fully contains the circle.
	// There is no intersection or touching between the polygon's boundary and the circle's boundary.
	RelationshipCirclePolyTreeContainedByPoly
)

// String returns the name of the RelationshipCirclePolygon constant as a string.
// It panics if the value is unsupported.
func (r RelationshipCirclePolygon) String() string {
	switch r {
	case RelationshipCirclePolyTreeMiss:
		return "RelationshipCirclePolyTreeMiss"
	case RelationshipCirclePolyTreeExternallyTouching:
		return "RelationshipCirclePolyTreeExternallyTouching"
	case RelationshipCirclePolyTreeExternallyTangent:
		return "RelationshipCirclePolyTreeExternallyTangent"
	case RelationshipCirclePolyTreeIntersection:
		return "RelationshipCirclePolyTreeIntersection"
	case RelationshipCirclePolyTreeInternallyTangent:
		return "RelationshipCirclePolyTreeInternallyTangent"
	case RelationshipCirclePolyTreeInternallyTouching:
		return "RelationshipCirclePolyTreeInternallyTouching"
	case RelationshipCirclePolyTreeContainedByCircle:
		return "RelationshipCirclePolyTreeContainedByCircle"
	case RelationshipCirclePolyTreeContainedByPoly:
		return "RelationshipCirclePolyTreeContainedByPoly"
	default:
		panic(fmt.Errorf("unsupported RelationshipCirclePolygon value"))
	}
}

// RelationshipCirclePolyTree represents the spatial relationships between a [Circle] and the polygons
// in a [PolyTree].
//
// This type is a mapping of pointers to [PolyTree] nodes to their respective [RelationshipCirclePolygon]
// values. It indicates how a given [Circle] relates to each polygon in the [PolyTree], such as whether
// the circle is outside, intersects, is contained by, or contains the polygon.
//
// Type Parameters:
//   - T: A numeric type that satisfies the [SignedNumber] constraint (e.g., int, float64).
//
// Notes:
//   - Each [PolyTree] node represents a single polygon, which could be a solid polygon or a hole.
//   - This type does not automatically account for hierarchical relationships (e.g., parent-child relationships).
//     It simply provides the direct relationship of the [Circle] to each polygon in the [PolyTree].
type RelationshipCirclePolyTree[T SignedNumber] map[*PolyTree[T]]RelationshipCirclePolygon

// RelationshipRectangleCircle defines the possible spatial relationships
// between a [Circle] and a [Rectangle] in a 2D plane.
//
// This enumeration categorizes how a [Circle] relates to a [Rectangle] based on
// their positions and intersections. It can distinguish whether
// the [Circle] is fully inside, outside or intersecting the rectangle.
type RelationshipRectangleCircle uint8

// Valid values for RelationshipRectangleCircle
const (
	// RelationshipRectangleCircleMiss - Circle and rectangle are disjoint.
	RelationshipRectangleCircleMiss RelationshipRectangleCircle = iota

	// RelationshipRectangleCircleContainedByRectangle - Circle is fully contained within the rectangle with no touching/intersection.
	// todo: verify cases of touching return intersection
	RelationshipRectangleCircleContainedByRectangle

	// RelationshipRectangleCircleContainedByCircle - Rectangle is fully contained within the circle with no intersection/touching.
	// todo: verify cases of touching return intersection
	RelationshipRectangleCircleContainedByCircle

	// RelationshipRectangleCircleIntersection - Circle and rectangle intersect but are not fully contained.
	RelationshipRectangleCircleIntersection
)

// String returns the string representation of a [RelationshipRectangleCircle].
//
// This function converts the [RelationshipRectangleCircle] enum value into a
// corresponding string for improved readability and debugging.
//
// Panics:
//   - If the [RelationshipRectangleCircle] has an unsupported or undefined value,
//     the function will panic.
//
// Returns:
//   - string: The name of the [RelationshipRectangleCircle] enum value.
func (r RelationshipRectangleCircle) String() string {
	switch r {
	case RelationshipRectangleCircleMiss:
		return "RelationshipRectangleCircleMiss"
	case RelationshipRectangleCircleContainedByRectangle:
		return "RelationshipRectangleCircleContainedByRectangle"
	case RelationshipRectangleCircleContainedByCircle:
		return "RelationshipRectangleCircleContainedByCircle"
	case RelationshipRectangleCircleIntersection:
		return "RelationshipRectangleCircleIntersection"
	default:
		panic(fmt.Errorf("unsupported RelationshipRectangleCircle"))
	}
}

// RelationshipLineSegmentLineSegment defines the possible spatial relationships
// between two line segments in a 2D plane.
//
// This enumeration categorizes how two line segments relate to each other based on
// their positions, intersections, collinearity, etc.
type RelationshipLineSegmentLineSegment int8

// Valid values for RelationshipLineSegmentLineSegment
const (
	// RelationshipLineSegmentLineSegmentCollinearDisjoint - segments are collinear and do not intersect, overlap, or touch at any point.
	RelationshipLineSegmentLineSegmentCollinearDisjoint RelationshipLineSegmentLineSegment = iota - 1

	// RelationshipLineSegmentLineSegmentMiss - the segments are not collinear, disjoint and do not intersect, overlap, or touch at any point.
	RelationshipLineSegmentLineSegmentMiss

	// RelationshipLineSegmentLineSegmentIntersects - the segments intersect at a unique point that is not an endpoint.
	RelationshipLineSegmentLineSegmentIntersects

	// RelationshipLineSegmentLineSegmentAeqC - point A of segment AB coincides with Point C of segment CD.
	RelationshipLineSegmentLineSegmentAeqC

	// RelationshipLineSegmentLineSegmentAeqD - point A of segment AB coincides with Point D of segment CD.
	RelationshipLineSegmentLineSegmentAeqD

	// RelationshipLineSegmentLineSegmentBeqC - point End of segment AB coincides with Point C of segment CD.
	RelationshipLineSegmentLineSegmentBeqC

	// RelationshipLineSegmentLineSegmentBeqD - point End of segment AB coincides with Point D of segment CD.
	RelationshipLineSegmentLineSegmentBeqD

	// RelationshipLineSegmentLineSegmentAonCD - point A lies on LineSegment CD.
	RelationshipLineSegmentLineSegmentAonCD

	// RelationshipLineSegmentLineSegmentBonCD - point End lies on LineSegment CD.
	RelationshipLineSegmentLineSegmentBonCD

	// RelationshipLineSegmentLineSegmentConAB - point C lies on LineSegment AB.
	RelationshipLineSegmentLineSegmentConAB

	// RelationshipLineSegmentLineSegmentDonAB - [oint D lies on LineSegment AB.
	RelationshipLineSegmentLineSegmentDonAB

	// RelationshipLineSegmentLineSegmentCollinearAonCD - point A lies on LineSegment CD (partial overlap), and line segments are collinear.
	RelationshipLineSegmentLineSegmentCollinearAonCD

	// RelationshipLineSegmentLineSegmentCollinearBonCD - point End lies on LineSegment CD (partial overlap), and line segments are collinear.
	RelationshipLineSegmentLineSegmentCollinearBonCD

	// RelationshipLineSegmentLineSegmentCollinearABinCD - segment AB is fully contained within segment CD.
	RelationshipLineSegmentLineSegmentCollinearABinCD

	// RelationshipLineSegmentLineSegmentCollinearCDinAB - segment CD is fully contained within segment AB.
	RelationshipLineSegmentLineSegmentCollinearCDinAB

	// RelationshipLineSegmentLineSegmentCollinearEqual - the segments AB and CD are exactly equal, sharing both endpoints in the same locations.
	RelationshipLineSegmentLineSegmentCollinearEqual
)

// String provides a string representation of a [RelationshipLineSegmentLineSegment] value.
//
// This method converts the [RelationshipLineSegmentLineSegment] enum value into a human-readable string,
// allowing for easier debugging and output interpretation. Each enum value maps to its corresponding
// string name.
//
// Returns:
//   - string: A string representation of the [RelationshipLineSegmentLineSegment].
//
// Panics:
//   - If the enum value is not recognized, the method panics with an error indicating
//     an unsupported [RelationshipLineSegmentLineSegment] value.
func (r RelationshipLineSegmentLineSegment) String() string {
	switch r {
	case RelationshipLineSegmentLineSegmentCollinearDisjoint:
		return "RelationshipLineSegmentLineSegmentCollinearDisjoint"
	case RelationshipLineSegmentLineSegmentMiss:
		return "RelationshipLineSegmentLineSegmentMiss"
	case RelationshipLineSegmentLineSegmentIntersects:
		return "RelationshipLineSegmentLineSegmentIntersects"
	case RelationshipLineSegmentLineSegmentAeqC:
		return "RelationshipLineSegmentLineSegmentAeqC"
	case RelationshipLineSegmentLineSegmentAeqD:
		return "RelationshipLineSegmentLineSegmentAeqD"
	case RelationshipLineSegmentLineSegmentBeqC:
		return "RelationshipLineSegmentLineSegmentBeqC"
	case RelationshipLineSegmentLineSegmentBeqD:
		return "RelationshipLineSegmentLineSegmentBeqD"
	case RelationshipLineSegmentLineSegmentAonCD:
		return "RelationshipLineSegmentLineSegmentAonCD"
	case RelationshipLineSegmentLineSegmentBonCD:
		return "RelationshipLineSegmentLineSegmentBonCD"
	case RelationshipLineSegmentLineSegmentConAB:
		return "RelationshipLineSegmentLineSegmentConAB"
	case RelationshipLineSegmentLineSegmentDonAB:
		return "RelationshipLineSegmentLineSegmentDonAB"
	case RelationshipLineSegmentLineSegmentCollinearAonCD:
		return "RelationshipLineSegmentLineSegmentCollinearAonCD"
	case RelationshipLineSegmentLineSegmentCollinearBonCD:
		return "RelationshipLineSegmentLineSegmentCollinearBonCD"
	case RelationshipLineSegmentLineSegmentCollinearABinCD:
		return "RelationshipLineSegmentLineSegmentCollinearABinCD"
	case RelationshipLineSegmentLineSegmentCollinearCDinAB:
		return "RelationshipLineSegmentLineSegmentCollinearCDinAB"
	case RelationshipLineSegmentLineSegmentCollinearEqual:
		return "RelationshipLineSegmentLineSegmentCollinearEqual"
	default:
		panic(fmt.Errorf("unsupported RelationshipLineSegmentLineSegment"))
	}
}

// RelationshipPointCircle defines the possible spatial relationships between a point
// and a circle in a 2D plane. This type categorizes whether the point is inside,
// outside, or exactly on the circle's circumference.
type RelationshipPointCircle uint8

// Valid values for RelationshipPointCircle
const (
	// RelationshipPointCircleMiss - Point is outside the circle.
	RelationshipPointCircleMiss RelationshipPointCircle = iota

	// RelationshipPointCircleOnCircumference - Point lies exactly on the circle's circumference.
	RelationshipPointCircleOnCircumference

	// RelationshipPointCircleContainedByCircle - Point is inside the circle.
	RelationshipPointCircleContainedByCircle
)

func (r RelationshipPointCircle) String() string {
	switch r {
	case RelationshipPointCircleMiss:
		return "RelationshipPointCircleMiss"
	case RelationshipPointCircleOnCircumference:
		return "RelationshipPointCircleOnCircumference"
	case RelationshipPointCircleContainedByCircle:
		return "RelationshipPointCircleContainedByCircle"
	default:
		panic(fmt.Errorf("unsupported RelationshipPointCircle"))
	}
}

// RelationshipPointLineSegment defines the possible spatial relationships
// between a point and a line segment.
//
// This enumeration is used to classify how a point relates to a line segment,
// whether it lies on the infinite line, on the segment itself, or elsewhere.
type RelationshipPointLineSegment int8

// Valid values for RelationshipPointLineSegment
const (
	// RelationshipPointLineSegmentCollinearDisjoint - Point lies on the infinite line but not the segment.
	RelationshipPointLineSegmentCollinearDisjoint RelationshipPointLineSegment = iota - 1

	// RelationshipPointLineSegmentMiss - Point misses the line segment entirely.
	RelationshipPointLineSegmentMiss

	// RelationshipPointLineSegmentPointEqStart - Point coincides with the start of the segment.
	RelationshipPointLineSegmentPointEqStart

	// RelationshipPointLineSegmentPointEqEnd - Point coincides with the end of the segment.
	RelationshipPointLineSegmentPointEqEnd

	// RelationshipPointLineSegmentPointOnLineSegment - Point lies on the segment (not at an endpoint).
	RelationshipPointLineSegmentPointOnLineSegment
)

func (r RelationshipPointLineSegment) String() string {
	switch r {
	case RelationshipPointLineSegmentCollinearDisjoint:
		return "RelationshipPointLineSegmentCollinearDisjoint"
	case RelationshipPointLineSegmentMiss:
		return "RelationshipPointLineSegmentMiss"
	case RelationshipPointLineSegmentPointEqStart:
		return "RelationshipPointLineSegmentPointEqStart"
	case RelationshipPointLineSegmentPointEqEnd:
		return "RelationshipPointLineSegmentPointEqEnd"
	case RelationshipPointLineSegmentPointOnLineSegment:
		return "RelationshipPointLineSegmentPointOnLineSegment"
	default:
		panic(fmt.Errorf("unsupported RelationshipPointLineSegment"))
	}
}

// RelationshipPointPolyTree defines the possible spatial relationships between a point
// and a polygon in a 2D plane. This type accounts for the presence of holes, solid regions,
// and nested islands within the polygon.
//
// This enumeration provides fine-grained distinctions for where a point lies relative
// to the polygon's structure, including its boundaries, holes, and nested regions.
// todo: refactor. just have this represent the relationship with the current polygon in the polytree. can return a slice for each polygon
type RelationshipPointPolyTree int8

// Valid values for RelationshipPointPolyTree
const (
	// PPTRPointInHole indicates the point is inside a hole within the polygon.
	// Holes are void regions within the polygon that are not part of its solid area.
	PPTRPointInHole RelationshipPointPolyTree = iota - 1

	// PPTRPointOutside indicates the point lies outside the root polygon.
	// This includes points outside the boundary and not within any nested holes or islands.
	PPTRPointOutside

	// PPTRPointOnVertex indicates the point coincides with a vertex of the polygon,
	// including vertices of its holes or nested islands.
	PPTRPointOnVertex

	// PPTRPointOnEdge indicates the point lies exactly on an edge of the polygon.
	// This includes edges of the root polygon, its holes, or its nested islands.
	PPTRPointOnEdge

	// PPTRPointInside indicates the point is strictly inside the solid area of the polygon,
	// excluding any holes within the polygon.
	PPTRPointInside

	// PPTRPointInsideIsland indicates the point lies within a nested island inside the polygon.
	// Islands are solid regions contained within holes of the polygon.
	PPTRPointInsideIsland
)

func (r RelationshipPointPolyTree) String() string {
	switch r {
	case PPTRPointInHole:
		return "PPTRPointInHole"
	case PPTRPointOutside:
		return "PPTRPointOutside"
	case PPTRPointOnVertex:
		return "PPTRPointOnVertex"
	case PPTRPointOnEdge:
		return "PPTRPointOnEdge"
	case PPTRPointInside:
		return "PPTRPointInside"
	case PPTRPointInsideIsland:
		return "PPTRPointInsideIsland"
	default:
		panic(fmt.Errorf("unsupported RelationshipPointPolyTree"))
	}
}

// RelationshipPointRectangle defines the possible spatial relationships between a point
// and a rectangle in a 2D plane. This type categorizes whether the point is inside,
// outside, on an edge, or on a vertex of the rectangle.
type RelationshipPointRectangle uint8

// Valid values for RelationshipPointRectangle
const (
	// RelationshipPointRectangleMiss - The point lies outside the rectangle.
	RelationshipPointRectangleMiss RelationshipPointRectangle = iota

	// RelationshipPointRectangleContainedByRectangle - The point lies strictly inside the rectangle.
	RelationshipPointRectangleContainedByRectangle

	// RelationshipPointRectanglePointOnVertex - The point lies on a vertex of the rectangle.
	RelationshipPointRectanglePointOnVertex

	// RelationshipPointRectanglePointOnEdge - The point lies on an edge of the rectangle.
	RelationshipPointRectanglePointOnEdge
)

// RelationshipLineSegmentPolyTree defines the possible spatial relationships
// between a line segment and a PolyTree, which is a hierarchical structure
// of polygons with holes and nested islands.
//
// This enumeration categorizes how a line segment relates to the PolyTree
// based on its position, intersections, and interactions with the boundaries
// of polygons and holes.
// todo: refactor so relationship is just for current polygon
type RelationshipLineSegmentPolyTree uint8

// RelationshipLineSegmentPolyTree describes the relationship between a line segment
// and a PolyTree, capturing whether the segment is inside, outside, or intersects
// boundaries of the PolyTree's polygons.
const (
	// PTLRMiss indicates that the line segment lies entirely outside the PolyTree.
	PTLRMiss RelationshipLineSegmentPolyTree = iota

	// PTLRInsideSolid indicates that the line segment lies entirely within a solid polygon
	// in the PolyTree, without crossing any boundaries.
	PTLRInsideSolid

	// PTLRInsideHole indicates that the line segment lies entirely within a hole
	// in the PolyTree, without crossing any boundaries.
	PTLRInsideHole

	// PTLRIntersectsBoundary indicates that the line segment crosses one or more boundaries
	// within the PolyTree, transitioning between solid and hole regions.
	PTLRIntersectsBoundary
)

func (r RelationshipLineSegmentPolyTree) String() string {
	switch r {
	case PTLRMiss:
		return "PTLRMiss"
	case PTLRInsideSolid:
		return "PTLRInsideSolid"
	case PTLRInsideHole:
		return "PTLRInsideHole"
	case PTLRIntersectsBoundary:
		return "PTLRIntersectsBoundary"
	default:
		panic(fmt.Errorf("unsupported RelationshipLineSegmentPolyTree"))
	}
}

// RelationshipLineSegmentRectangle defines the possible spatial relationships between
// a line segment and a rectangle in a 2D plane. This type categorizes whether the segment
// is inside, outside, touches edges or vertices, or intersects the rectangle's boundary.
type RelationshipLineSegmentRectangle uint8

// Valid values for RelationshipLineSegmentRectangle
const (
	// RelationshipLineSegmentRectangleMiss - The segment lies entirely outside the rectangle.
	RelationshipLineSegmentRectangleMiss RelationshipLineSegmentRectangle = iota

	// RelationshipLineSegmentRectangleEndTouchesEdgeExternally - The segment lies outside the rectangle and one end touches an edge.
	RelationshipLineSegmentRectangleEndTouchesEdgeExternally

	// RelationshipLineSegmentRectangleEndTouchesVertexExternally - The segment lies outside the rectangle and one end touches a vertex.
	RelationshipLineSegmentRectangleEndTouchesVertexExternally

	// RelationshipLineSegmentRectangleContainedByRectangle - The segment lies entirely within the rectangle.
	RelationshipLineSegmentRectangleContainedByRectangle

	// RelationshipLineSegmentRectangleEndTouchesEdgeInternally - The segment lies within the rectangle and one end touches an edge without crossing the boundary.
	RelationshipLineSegmentRectangleEndTouchesEdgeInternally

	// RelationshipLineSegmentRectangleEndTouchesVertexInternally - The segment lies within the rectangle and one end touches a vertex without crossing the boundary.
	RelationshipLineSegmentRectangleEndTouchesVertexInternally

	// RelationshipLineSegmentRectangleEdgeCollinear - The segment lies entirely on one edge of the rectangle, not touching a vertex.
	RelationshipLineSegmentRectangleEdgeCollinear

	// RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex - The segment lies entirely on one edge of the rectangle and one or both ends touch a vertex.
	RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex

	// RelationshipLineSegmentRectangleIntersects - The segment crosses one or more edges of the rectangle.
	RelationshipLineSegmentRectangleIntersects

	// RelationshipLineSegmentRectangleEntersAndExits - The segment enters through one edge and exits through another.
	RelationshipLineSegmentRectangleEntersAndExits
)

// String returns the string representation of the RelationshipLineSegmentRectangle constant.
//
// The method maps each RelationshipLineSegmentRectangle value to a descriptive string. It
// panics if an unsupported value is encountered, ensuring that the relationship is always
// valid within its defined range.
//
// Returns:
//   - string: The string representation of the RelationshipLineSegmentRectangle.
func (r RelationshipLineSegmentRectangle) String() string {
	switch r {
	case RelationshipLineSegmentRectangleMiss:
		return "RelationshipLineSegmentRectangleMiss"
	case RelationshipLineSegmentRectangleEndTouchesEdgeExternally:
		return "RelationshipLineSegmentRectangleEndTouchesEdgeExternally"
	case RelationshipLineSegmentRectangleEndTouchesVertexExternally:
		return "RelationshipLineSegmentRectangleEndTouchesVertexExternally"
	case RelationshipLineSegmentRectangleContainedByRectangle:
		return "RelationshipLineSegmentRectangleContainedByRectangle"
	case RelationshipLineSegmentRectangleEndTouchesEdgeInternally:
		return "RelationshipLineSegmentRectangleEndTouchesEdgeInternally"
	case RelationshipLineSegmentRectangleEndTouchesVertexInternally:
		return "RelationshipLineSegmentRectangleEndTouchesVertexInternally"
	case RelationshipLineSegmentRectangleEdgeCollinear:
		return "RelationshipLineSegmentRectangleEdgeCollinear"
	case RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex:
		return "RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex"
	case RelationshipLineSegmentRectangleIntersects:
		return "RelationshipLineSegmentRectangleIntersects"
	case RelationshipLineSegmentRectangleEntersAndExits:
		return "RelationshipLineSegmentRectangleEntersAndExits"
	default:
		panic(fmt.Errorf("unsupported RelationshipLineSegmentRectangle: %d", r))
	}
}

// RelationshipRectangleRectangle defines the possible spatial relationships
// between two rectangles in a 2D plane.
//
// This enumeration categorizes how two rectangles relate to each other based on
// their positions, overlaps, and containment.
type RelationshipRectangleRectangle uint8

// Valid values for RelationshipRectangleRectangle
const (
	// RelationshipRectangleRectangleMiss - Rectangles are disjoint, with no overlap or touching.
	RelationshipRectangleRectangleMiss RelationshipRectangleRectangle = iota

	// RelationshipRectangleRectangleSharedEdge - Rectangles share a complete edge but do not overlap.
	RelationshipRectangleRectangleSharedEdge

	// RelationshipRectangleRectangleSharedVertex - Rectangles share a single vertex but do not overlap.
	RelationshipRectangleRectangleSharedVertex

	// RelationshipRectangleRectangleIntersection - Rectangles overlap but neither is fully contained in the other.
	// todo: standardise constant wording - should be "intersecting" across all
	RelationshipRectangleRectangleIntersection

	// RelationshipRectangleRectangleContained - One rectangle is fully contained within the other without touching edges.
	RelationshipRectangleRectangleContained

	// RelationshipRectangleRectangleContainedTouching - One rectangle is fully contained within the other and touches edges.
	RelationshipRectangleRectangleContainedTouching

	// RelationshipRectangleRectangleEqual - Rectangles are identical in position and size.
	RelationshipRectangleRectangleEqual
)
