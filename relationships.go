package geom2d

import "fmt"

// CircleCircleRelationship defines the possible spatial relationships
// between two circles in a 2D plane.
//
// This enumeration categorizes how two circles relate to each other based on
// their positions and radii. It covers scenarios such as being disjoint, touching,
// overlapping, or one circle being contained within the other.
type CircleCircleRelationship uint8

// Valid values for CircleCircleRelationship
const (
	CCRMiss             CircleCircleRelationship = iota // Circles are disjoint
	CCRTouchingExternal                                 // Circles are externally tangent
	CCROverlapping                                      // Circles overlap and intersect at two points
	CCRTouchingInternal                                 // Circles are internally tangent
	CCRContained                                        // One circle is fully contained within the other
	CCREqual                                            // Circles are identical
)

// String returns the string representation of a [CircleCircleRelationship].
//
// The function converts the [CircleCircleRelationship] enum value into its
// corresponding string representation for readability and debugging purposes.
//
// Panics:
//   - If the [CircleCircleRelationship] has an unsupported or undefined value, the function will panic with a descriptive error.
//
// Returns:
//   - string: The name of the [CircleCircleRelationship] enum value.
func (r *CircleCircleRelationship) String() string {
	switch *r {
	case CCRMiss:
		return "CCRMiss"
	case CCRTouchingExternal:
		return "CCRTouchingExternal"
	case CCROverlapping:
		return "CCROverlapping"
	case CCRTouchingInternal:
		return "CCRTouchingInternal"
	case CCRContained:
		return "CCRContained"
	case CCREqual:
		return "CCREqual"
	default:
		panic(fmt.Errorf("unsupported CircleCircleRelationship"))
	}
}

// CircleLineSegmentRelationship defines the possible spatial relationships
// between a [Circle] and a [LineSegment] in a 2D plane.
//
// This enumeration categorizes how a [LineSegment] relates to a [Circle] based on
// its position, intersection, and tangency.
type CircleLineSegmentRelationship uint8

// Valid values for CircleLineSegmentRelationship
const (
	// CLROutside indicates that the line segment lies completely outside the circle,
	// with no intersection or tangency.
	CLROutside CircleLineSegmentRelationship = iota

	// CLRInside indicates that the line segment lies completely within the circle,
	// with both endpoints inside the circle's boundary.
	CLRInside

	// CLRIntersecting indicates that the line segment intersects the circle at exactly
	// two distinct points.
	CLRIntersecting

	// CLRTangent indicates that the line segment is tangent to the circle, touching it
	// at exactly one point with the tangent forming a 90-degree angle to the circle's radius.
	CLRTangent

	// CLROneEndOnCircumferenceOutside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies outside the circle.
	CLROneEndOnCircumferenceOutside

	// CLROneEndOnCircumferenceInside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies inside the circle.
	CLROneEndOnCircumferenceInside

	// CLRBothEndsOnCircumference indicates that both endpoints of the line segment lie
	// exactly on the circle's boundary. The line segment does not extend inside or outside
	// the circle.
	CLRBothEndsOnCircumference
)

// String returns the string representation of a [CircleLineSegmentRelationship].
//
// This function converts the [CircleLineSegmentRelationship] enum value into a
// corresponding string for improved readability and debugging.
//
// Panics:
//   - If the [CircleLineSegmentRelationship] has an unsupported or undefined value,
//     the function will panic.
//
// Returns:
//   - string: The name of the [CircleLineSegmentRelationship] enum value.
func (r *CircleLineSegmentRelationship) String() string {
	switch *r {
	case CLROutside:
		return "CLROutside"
	case CLRInside:
		return "CLRInside"
	case CLRIntersecting:
		return "CLRIntersecting"
	case CLRTangent:
		return "CLRTangent"
	case CLROneEndOnCircumferenceOutside:
		return "CLROneEndOnCircumferenceOutside"
	case CLROneEndOnCircumferenceInside:
		return "CLROneEndOnCircumferenceInside"
	case CLRBothEndsOnCircumference:
		return "CLRBothEndsOnCircumference"
	default:
		panic(fmt.Errorf("unsupported CircleLineSegmentRelationship"))
	}
}

// CirclePolyTreeRelationship defines the possible spatial relationships
// between a [Circle] and a [PolyTree], which is a hierarchical structure of polygons
// with holes and nested islands.
//
// This enumeration categorizes how a [Circle] relates to the [PolyTree], based on
// whether the [Circle] is inside, outside, intersecting, or touching the solid
// and hole regions of the [PolyTree].
type CirclePolyTreeRelationship uint8

// Valid values for CirclePolyTreeRelationship
const (
	CPTRMiss                 CirclePolyTreeRelationship = iota // Circle lies entirely outside the PolyTree
	CPTRTouchesSolidBoundary                                   // Circle touches the boundary of a solid polygon
	CPTRTouchesHoleBoundary                                    // Circle touches the boundary of a hole
	CPTRIntersectsSolid                                        // Circle intersects one or more solid polygons
	CPTRIntersectsHole                                         // Circle intersects one or more holes
	CPTRCircleInSolid                                          // Circle is fully contained within a solid polygon
	CPTRCircleInHole                                           // Circle is fully contained within a hole
	CPTRSolidInCircle                                          // A solid polygon is fully contained within the circle
	CPTRHoleInCircle                                           // A hole is fully contained within the circle
	CPTRSpansSolid                                             // Circle spans across multiple solid polygons
	CPTRSpansHole                                              // Circle spans across multiple holes
)

// CircleRectangleRelationship defines the possible spatial relationships
// between a [Circle] and a [Rectangle] in a 2D plane.
//
// This enumeration categorizes how a [Circle] relates to a [Rectangle] based on
// their positions and intersections. It can distinguish whether
// the [Circle] is fully inside, outside or intersecting the rectangle.
type CircleRectangleRelationship uint8

// Valid values for CircleRectangleRelationship
const (
	CRRMiss         CircleRectangleRelationship = iota // Circle and rectangle are disjoint
	CRRCircleInRect                                    // Circle is fully contained within the rectangle
	CRRRectInCircle                                    // Rectangle is fully contained within the circle
	CRRIntersection                                    // Circle and rectangle intersect but are not fully contained
)

// String returns the string representation of a [CircleRectangleRelationship].
//
// This function converts the [CircleRectangleRelationship] enum value into a
// corresponding string for improved readability and debugging.
//
// Panics:
//   - If the [CircleRectangleRelationship] has an unsupported or undefined value,
//     the function will panic.
//
// Returns:
//   - string: The name of the [CircleRectangleRelationship] enum value.
func (r *CircleRectangleRelationship) String() string {
	switch *r {
	case CRRMiss:
		return "CRRMiss"
	case CRRCircleInRect:
		return "CRRCircleInRect"
	case CRRRectInCircle:
		return "CRRRectInCircle"
	case CRRIntersection:
		return "CRRIntersection"
	default:
		panic(fmt.Errorf("unsupported CircleRectangleRelationship"))
	}
}

// LineSegmentLineSegmentRelationship defines the possible spatial relationships
// between two line segments in a 2D plane.
//
// This enumeration categorizes how two line segments relate to each other based on
// their positions, intersections, and collinearity.
//
// Values:
//   - LLRCollinearDisjoint (-1): The segments are collinear but do not intersect, overlap, or touch at any point.
//   - LLRMiss (0): The segments are not collinear and do not intersect, overlap, or touch at any point.
//   - LLRIntersects (1): The segments intersect at a unique point that is not an endpoint of either segment.
//   - LLRAeqC (2): The starting point (A) of segment AB coincides with the starting point (C) of segment CD.
//   - LLRAeqD (3): The starting point (A) of segment AB coincides with the ending point (D) of segment CD.
//   - LLRBeqC (4): The ending point (B) of segment AB coincides with the starting point (C) of segment CD.
//   - LLRBeqD (5): The ending point (B) of segment AB coincides with the ending point (D) of segment CD.
//   - LLRAonCD (6): The starting point (A) of segment AB lies on segment CD.
//   - LLRBonCD (7): The ending point (B) of segment AB lies on segment CD.
//   - LLRConAB (8): The starting point (C) of segment CD lies on segment AB.
//   - LLRDonAB (9): The ending point (D) of segment CD lies on segment AB.
//   - LLRCollinearAonCD (10): The starting point (A) of segment AB lies on segment CD, with the segments being collinear.
//   - LLRCollinearBonCD (11): The ending point (B) of segment AB lies on segment CD, with the segments being collinear.
//   - LLRCollinearABinCD (12): Segment AB is fully contained within segment CD.
//   - LLRCollinearCDinAB (13): Segment CD is fully contained within segment AB.
//   - LLRCollinearEqual (14): Segments AB and CD are collinear and exactly equal, sharing both endpoints.
type LineSegmentLineSegmentRelationship int8

// Valid values for LineSegmentLineSegmentRelationship
const (
	LLRCollinearDisjoint LineSegmentLineSegmentRelationship = iota - 1 // Segments are collinear and do not intersect, overlap, or touch at any point.
	LLRMiss                                                            // The segments are not collinear, disjoint and do not intersect, overlap, or touch at any point.
	LLRIntersects                                                      // The segments intersect at a unique point that is not an endpoint.
	LLRAeqC                                                            // Point A of segment AB coincides with Point C of segment CD
	LLRAeqD                                                            // Point A of segment AB coincides with Point D of segment CD
	LLRBeqC                                                            // Point End of segment AB coincides with Point C of segment CD
	LLRBeqD                                                            // Point End of segment AB coincides with Point D of segment CD
	LLRAonCD                                                           // Point A lies on LineSegment CD
	LLRBonCD                                                           // Point End lies on LineSegment CD
	LLRConAB                                                           // Point C lies on LineSegment AB
	LLRDonAB                                                           // Point D lies on LineSegment AB
	LLRCollinearAonCD                                                  // Point A lies on LineSegment CD (partial overlap), and line segments are collinear
	LLRCollinearBonCD                                                  // Point End lies on LineSegment CD (partial overlap), and line segments are collinear
	LLRCollinearABinCD                                                 // Segment AB is fully contained within segment CD
	LLRCollinearCDinAB                                                 // Segment CD is fully contained within segment AB
	LLRCollinearEqual                                                  // The segments AB and CD are exactly equal, sharing both endpoints in the same locations.
)

func (r *LineSegmentLineSegmentRelationship) String() string {
	switch *r {
	case LLRCollinearDisjoint:
		return "LLRCollinearDisjoint"
	case LLRMiss:
		return "LLRMiss"
	case LLRIntersects:
		return "LLRIntersects"
	case LLRAeqC:
		return "LLRAeqC"
	case LLRAeqD:
		return "LLRAeqD"
	case LLRBeqC:
		return "LLRBeqC"
	case LLRBeqD:
		return "LLRBeqD"
	case LLRAonCD:
		return "LLRAonCD"
	case LLRBonCD:
		return "LLRAonCD"
	case LLRConAB:
		return "LLRConAB"
	case LLRDonAB:
		return "LLRDonAB"
	case LLRCollinearAonCD:
		return "LLRCollinearAonCD"
	case LLRCollinearBonCD:
		return "LLRCollinearBonCD"
	case LLRCollinearABinCD:
		return "LLRCollinearABinCD"
	case LLRCollinearCDinAB:
		return "LLRCollinearCDinAB"
	case LLRCollinearEqual:
		return "LLRCollinearEqual"
	default:
		panic(fmt.Errorf("unsupported LineSegmentLineSegmentRelationship"))
	}
}

// PointCircleRelationship defines the possible spatial relationships between a point
// and a circle in a 2D plane. This type categorizes whether the point is inside,
// outside, or exactly on the circle's circumference.
//
// Values:
//   - PCROutside (0): The point lies outside the circle, meaning its distance
//     from the circle's center is greater than the circle's radius.
//   - PCROnCircumference (1): The point lies exactly on the circle's circumference,
//     meaning its distance from the circle's center is equal to the circle's radius
//     (within a specified tolerance, if applicable).
//   - PCRInside (2): The point is inside the circle, meaning its distance from
//     the circle's center is less than the circle's radius.
type PointCircleRelationship uint8

// Valid values for PointCircleRelationship
const (
	PCROutside         PointCircleRelationship = iota // Point is outside the circle
	PCROnCircumference                                // Point lies exactly on the circle's circumference
	PCRInside                                         // Point is inside the circle
)

func (r *PointCircleRelationship) String() string {
	switch *r {
	case PCROutside:
		return "PCROutside"
	case PCROnCircumference:
		return "PCROnCircumference"
	case PCRInside:
		return "PCRInside"
	default:
		panic(fmt.Errorf("unsupported PointCircleRelationship"))
	}
}

// PointLineSegmentRelationship defines the possible spatial relationships
// between a point and a line segment.
//
// This enumeration is used to classify how a point relates to a line segment,
// whether it lies on the infinite line, on the segment itself, or elsewhere.
//
// Values:
//   - PLSRPointOnLine (-1): The point lies on the infinite line defined by the
//     line segment but not within the segment's bounds.
//   - PLSRMiss (0): The point does not lie on the line segment or the infinite
//     line that extends through it.
//   - PLSRPointEqStart (1): The point coincides with the start of the line segment.
//   - PLSRPointEqEnd (2): The point coincides with the end of the line segment.
//   - PLSRPointOnLineSegment (3): The point lies on the segment itself but not
//     at either endpoint.
type PointLineSegmentRelationship int8

// Valid values for PointLineSegmentRelationship
const (
	PLRPointOnLine        PointLineSegmentRelationship = iota - 1 // Point lies on the infinite line but not the segment
	PLRMiss                                                       // Point misses the line segment entirely
	PLRPointEqStart                                               // Point coincides with the start of the segment
	PLRPointEqEnd                                                 // Point coincides with the end of the segment
	PLRPointOnLineSegment                                         // Point lies on the segment (not at an endpoint)
)

// PointPolyTreeRelationship defines the possible spatial relationships between a point
// and a polygon in a 2D plane. This type accounts for the presence of holes, solid regions,
// and nested islands within the polygon.
//
// This enumeration provides fine-grained distinctions for where a point lies relative
// to the polygon's structure, including its boundaries, holes, and nested regions.
//
// Values:
//   - PPTRPointInHole (-1): The point is inside a hole within the polygon. Holes are void regions
//     within the polygon that are not part of its solid area.
//   - PPTRPointOutside (0): The point lies outside the root polygon, including points outside
//     the boundary and not within any nested holes or islands.
//   - PPTRPointOnVertex (1): The point coincides with a vertex of the polygon, including vertices
//     of its holes or nested islands.
//   - PPTRPointOnEdge (2): The point lies exactly on an edge of the polygon. This includes edges of
//     the root polygon, its holes, or its nested islands.
//   - PPTRPointInside (3): The point is strictly inside the solid area of the polygon, excluding
//     any holes within the polygon.
//   - PPTRPointInsideIsland (4): The point lies within a nested island inside the polygon.
//     Islands are solid regions contained within holes of the polygon.
type PointPolyTreeRelationship int8

// Valid values for PointPolyTreeRelationship
const (
	// PPTRPointInHole indicates the point is inside a hole within the polygon.
	// Holes are void regions within the polygon that are not part of its solid area.
	PPTRPointInHole PointPolyTreeRelationship = iota - 1

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

// PointRectangleRelationship defines the possible spatial relationships between a point
// and a rectangle in a 2D plane. This type categorizes whether the point is inside,
// outside, on an edge, or on a vertex of the rectangle.
//
// Values:
//   - PRROutside (0): The point lies outside the rectangle, meaning it does not fall
//     within the rectangle's boundaries or on its edges or vertices.
//   - PRRInside (1): The point lies strictly inside the rectangle, meaning it falls
//     within the rectangle's boundaries but not on its edges or vertices.
//   - PRROnVertex (2): The point lies on one of the rectangle's vertices
//     (top-left, top-right, bottom-left, or bottom-right).
//   - PRROnEdge (3): The point lies on one of the rectangle's edges but not on a vertex.
type PointRectangleRelationship uint8

const (
	PRROutside  PointRectangleRelationship = iota // The point lies outside the rectangle.
	PRRInside                                     // The point lies strictly inside the rectangle.
	PRROnVertex                                   // The point lies on a vertex of the rectangle.
	PRROnEdge                                     // The point lies on an edge of the rectangle.
)

// PolyTreeLineSegmentRelationship defines the possible spatial relationships
// between a line segment and a PolyTree, which is a hierarchical structure
// of polygons with holes and nested islands.
//
// This enumeration categorizes how a line segment relates to the PolyTree
// based on its position, intersections, and interactions with the boundaries
// of polygons and holes.
//
// Values:
//   - PTLSMiss: The line segment is entirely outside the PolyTree and does not touch or intersect any of its polygons or holes.
//   - PTLSInsideSolid: The line segment is entirely within a solid polygon in the PolyTree, without touching any boundaries or holes.
//   - PTLSInsideHole: The line segment is entirely within a hole in the PolyTree, without touching any boundaries or solid polygons.
//   - PTLSIntersectsBoundary: The line segment crosses one or more boundaries in the PolyTree.
type PolyTreeLineSegmentRelationship uint8

// PolyTreeLineSegmentRelationship describes the relationship between a line segment
// and a PolyTree, capturing whether the segment is inside, outside, or intersects
// boundaries of the PolyTree's polygons.
const (
	// PTLRMiss indicates that the line segment lies entirely outside the PolyTree.
	PTLRMiss PolyTreeLineSegmentRelationship = iota

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

// RectangleLineSegmentRelationship defines the possible spatial relationships between
// a line segment and a rectangle in a 2D plane. This type categorizes whether the segment
// is inside, outside, touches edges or vertices, or intersects the rectangle's boundary.
//
// Values:
//   - RLROutside (0): The segment lies entirely outside the rectangle, with no intersection
//     or contact with its edges or vertices.
//   - RLROutsideEndTouchesEdge (1): The segment lies outside the rectangle, but one of its
//     endpoints touches an edge of the rectangle.
//   - RLROutsideEndTouchesVertex (2): The segment lies outside the rectangle, but one of its
//     endpoints touches a vertex of the rectangle.
//   - RLRInside (3): The segment lies entirely within the rectangle, without touching or crossing
//     its boundary.
//   - RLRInsideEndTouchesEdge (4): The segment lies within the rectangle, with one of its
//     endpoints touching an edge but not crossing the boundary.
//   - RLRInsideEndTouchesVertex (5): The segment lies within the rectangle, with one of its
//     endpoints touching a vertex but not crossing the boundary.
//   - RLROnEdge (6): The segment lies entirely on one of the rectangle's edges, without crossing
//     into its interior or exterior.
//   - RLROnEdgeEndTouchesVertex (7): The segment lies entirely on one of the rectangle's edges,
//     and one or both endpoints touch a vertex.
//   - RLRIntersects (8): The segment crosses one or more edges of the rectangle but does not
//     fully enter or exit through the boundaries.
//   - RLREntersAndExits (9): The segment enters the rectangle through one edge and exits through
//     another edge, crossing the interior.
type RectangleLineSegmentRelationship uint8

// Valid values for RectangleLineSegmentRelationship
const (
	RLROutside                 RectangleLineSegmentRelationship = iota // The segment lies entirely outside the rectangle.
	RLROutsideEndTouchesEdge                                           // The segment lies outside the rectangle and one end touches an edge.
	RLROutsideEndTouchesVertex                                         // The segment lies outside the rectangle and one end touches a vertex.
	RLRInside                                                          // The segment lies entirely within the rectangle.
	RLRInsideEndTouchesEdge                                            // The segment lies within the rectangle and one end touches an edge without crossing the boundary.
	RLRInsideEndTouchesVertex                                          // The segment lies within the rectangle and one end touches a vertex without crossing the boundary.
	RLROnEdge                                                          // The segment lies entirely on one edge of the rectangle.
	RLROnEdgeEndTouchesVertex                                          // The segment lies entirely on one edge of the rectangle and one or both ends touch a vertex.
	RLRIntersects                                                      // The segment crosses one or more edges of the rectangle.
	RLREntersAndExits                                                  // The segment enters through one edge and exits through another.
)

// RectangleRectangleRelationship defines the possible spatial relationships
// between two rectangles in a 2D plane.
//
// This enumeration categorizes how two rectangles relate to each other based on
// their positions, overlaps, and containment.
//
// Values:
//   - RRRMiss (0): The rectangles are disjoint, with no overlap or touching.
//   - RRRTouchingEdge (1): The rectangles share a complete edge but do not overlap.
//   - RRRTouchingVertex (2): The rectangles share a single vertex but do not overlap.
//   - RRRIntersecting (3): The rectangles overlap but neither is fully contained within the other.
//   - RRRContained (4): One rectangle is fully contained within the other without touching its edges.
//   - RRRTouchingContained (5): One rectangle is fully contained within the other and touches one or more edges.
//   - RRREqual (6): The rectangles are identical, sharing the same top-left and bottom-right coordinates.
type RectangleRectangleRelationship uint8

// Valid values for RectangleRectangleRelationship
const (
	RRRMiss              RectangleRectangleRelationship = iota // Rectangles are disjoint, with no overlap or touching.
	RRRTouchingEdge                                            // Rectangles share a complete edge but do not overlap.
	RRRTouchingVertex                                          // Rectangles share a single vertex but do not overlap.
	RRRIntersecting                                            // Rectangles overlap but neither is fully contained in the other.
	RRRContained                                               // One rectangle is fully contained within the other without touching edges.
	RRRTouchingContained                                       // One rectangle is fully contained within the other and touches edges.
	RRREqual                                                   // Rectangles are identical in position and size.
)
