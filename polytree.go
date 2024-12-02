// The polytree.go file defines the PolyTree type, a hierarchical representation
// of polygons that supports complex geometric operations such as unions, intersections,
// and subtractions. A PolyTree can represent solid polygons, holes, and nested structures,
// making it suitable for a wide range of computational geometry tasks.
//
// Key Features:
// - Boolean Operations: Perform union, intersection, and subtraction between polygons.
// - Hierarchical Representation: Support for nested polygons, where solid polygons
//   can contain hole polygons and vice versa.
// - Traversal Utilities: Methods for efficiently traversing and manipulating the hierarchy.
// - Intersection Detection: Robust handling of polygon edge intersections with metadata
//   for entry/exit relationships.
//
// The PolyTree type and its associated methods are designed to be flexible and performant,
// leveraging generic types to support various numeric types while maintaining precision.
//
// Usage Example:
// A PolyTree can be constructed, manipulated, and queried to perform advanced operations:
//
//     poly1 := []Point[int]{
//    		NewPoint(0, 0),
//   		NewPoint(10, 0),
//  		NewPoint(10, 10),
// 			NewPoint(0, 10),
//	   }
//     poly2 := []Point[int]{
//    		NewPoint(5, 5),
//   		NewPoint(15, 5),
//  		NewPoint(15, 15),
// 			NewPoint(5, 15),
//	   }
//     tree1, _ := NewPolyTree(poly1, PTSolid)
//     tree2, _ := NewPolyTree(poly2, PTSolid)
//     result, _ := tree1.BooleanOperation(tree2, BooleanUnion)
//
// This file also includes helper functions for sorting, relationship checking,
// and traversal to support the primary functionality of the PolyTree type.

package geom2d

import (
	"fmt"
	"slices"
)

// Valid values for BooleanOperation
const (
	// BooleanUnion represents the union operation, which combines two polygons into a single polygon
	// that encompasses the areas of both input polygons. Overlapping regions are merged.
	BooleanUnion BooleanOperation = iota

	// BooleanIntersection represents the intersection operation, which computes the region(s)
	// where two polygons overlap. The result is one or more polygons that covers only the shared area.
	BooleanIntersection

	// BooleanSubtraction represents the subtraction operation, which removes the area of one polygon
	// from another. The result is one or more polygons representing the area of the first polygon excluding
	// the overlapping region with the second polygon.
	BooleanSubtraction
)

// Valid values for PointPolygonRelationship
const (
	// PPRPointInHole indicates the point is inside a hole within the polygon.
	// Holes are void regions within the polygon that are not part of its solid area.
	PPRPointInHole PointPolygonRelationship = iota - 1

	// PPRPointOutside indicates the point lies outside the root polygon.
	// This includes points outside the boundary and not within any nested holes or islands.
	PPRPointOutside

	// PPRPointOnVertex indicates the point coincides with a vertex of the polygon,
	// including vertices of its holes or nested islands.
	PPRPointOnVertex

	// PPRPointOnEdge indicates the point lies exactly on an edge of the polygon.
	// This includes edges of the root polygon, its holes, or its nested islands.
	PPRPointOnEdge

	// PPRPointInside indicates the point is strictly inside the solid area of the polygon,
	// excluding any holes within the polygon.
	PPRPointInside

	// PPRPointInsideIsland indicates the point lies within a nested island inside the polygon.
	// Islands are solid regions contained within holes of the polygon.
	PPRPointInsideIsland
)

// Valid values for PolygonType
const (
	// PTSolid represents a solid region of the polygon, commonly referred to as an "island."
	// PTSolid polygons are the primary filled areas, excluding any void regions (holes).
	PTSolid PolygonType = iota

	// PTHole represents a void region of the polygon, often nested within a solid polygon.
	// Holes are not part of the filled area of the polygon and are treated as exclusions.
	PTHole
)

// Valid values for polyIntersectionType
const (
	// intersectionTypeNotSet indicates that the intersection type has not been set.
	// This is the default value for uninitialized points or non-intersection points.
	intersectionTypeNotSet polyIntersectionType = iota

	// intersectionTypeEntry indicates that the point serves as an entry point to the area of interest
	// when traversing the polygon during an operation.
	intersectionTypeEntry

	// intersectionTypeExit indicates that the point serves as an exit point from the area of interest
	// when traversing the polygon during an operation.
	intersectionTypeExit
)

// Valid values for polyPointType
const (
	// pointTypeOriginal indicates that the point is an original, unmodified vertex of the polygon.
	// These points are part of the polygon's initial definition.
	pointTypeOriginal polyPointType = iota

	// pointTypeOriginalAndIntersection indicates that the point is an original vertex that also serves
	// as an intersection point between polygons during operations such as union or intersection.
	pointTypeOriginalAndIntersection

	// pointTypeAddedIntersection indicates that the point is an intersection point that was added
	// during a polygon operation. These points are not part of the polygon's original definition
	// but are dynamically introduced for computational purposes.
	pointTypeAddedIntersection
)

// Valid values for polyTraversalDirection
const (
	// polyTraversalForward specifies that the traversal proceeds
	// in a counterclockwise direction through the polygon's vertices or edges.
	polyTraversalForward = polyTraversalDirection(iota)

	// polyTraversalReverse specifies that the traversal proceeds
	// in a clockwise direction through the polygon's vertices or edges.
	polyTraversalReverse
)

// Valid values for PolyTreeMismatch
const (
	// PTMNoMismatch indicates that there is no mismatch between the compared PolyTree structures.
	PTMNoMismatch PolyTreeMismatch = 0

	// PTMNilPolygonMismatch indicates that one of the PolyTree structures is nil while the other is not.
	PTMNilPolygonMismatch PolyTreeMismatch = 1 << iota

	// PTMContourMismatch indicates that the contours of the compared PolyTree structures are not identical.
	// This could mean differences in the points, their order, or their overall shape.
	PTMContourMismatch

	// PTMSiblingMismatch indicates that the siblings of the compared PolyTree structures do not match.
	// This could mean differences in the number of siblings, their contours, or their structure.
	PTMSiblingMismatch

	// PTMChildMismatch indicates that the children of the compared PolyTree structures do not match.
	// This could mean differences in the number of children, their contours, or their structure.
	PTMChildMismatch
)

// entryExitPointLookUpTable is a lookup table that determines the entry/exit type of intersection points
// between two polygons based on their types and the Boolean operation being performed.
//
// The structure of the lookup table is as follows:
//
//	map[BooleanOperation]map[PolygonType]map[PolygonType]map[bool]struct{
//	    poly1PointType polyIntersectionType
//	    poly2PointType polyIntersectionType
//	}
//
// Keys:
//   - BooleanOperation: The type of Boolean operation (e.g., Union, Intersection, Subtraction).
//   - PolygonType: The type of the polygons being compared (PTSolid or PTHole).
//   - bool: Whether the traversal direction indicates that poly1 is entering poly2 (true) or exiting (false).
//
// Values:
//   - poly1PointType: The intersection type for the current point in poly1.
//   - poly2PointType: The intersection type for the current point in poly2.
//
// Example Usage:
// When determining how to mark intersection points during a Boolean operation,
// this table provides the correct entry and exit types for each polygon.
//
// Structure:
//
// For BooleanUnion:
//   - Solid-Solid: Determines entry/exit depending on whether poly1 is entering or exiting poly2.
//   - Solid-Hole: Adjusts the entry/exit direction for the hole.
//
// For BooleanIntersection:
//   - Adjusts entry/exit depending on whether the traversal is into or out of holes and solids.
//
// For BooleanSubtraction:
//   - Marks both intersection points as exits or entries, depending on traversal.
var entryExitPointLookUpTable = map[BooleanOperation]map[PolygonType]map[PolygonType]map[bool]struct {
	poly1PointType polyIntersectionType
	poly2PointType polyIntersectionType
}{
	BooleanUnion: {
		PTSolid: {
			PTSolid: {
				true: {
					poly1PointType: intersectionTypeExit,
					poly2PointType: intersectionTypeEntry,
				},
				false: {
					poly1PointType: intersectionTypeEntry,
					poly2PointType: intersectionTypeExit,
				},
			},
			PTHole: {
				true: {
					poly1PointType: intersectionTypeEntry,
					poly2PointType: intersectionTypeExit,
				},
				false: {
					poly1PointType: intersectionTypeExit,
					poly2PointType: intersectionTypeEntry,
				},
			},
		},
		PTHole: {
			PTSolid: {
				true: {
					poly1PointType: intersectionTypeExit,
					poly2PointType: intersectionTypeEntry,
				},
				false: {
					poly1PointType: intersectionTypeEntry,
					poly2PointType: intersectionTypeExit,
				},
			},
			PTHole: {
				true: {
					poly1PointType: intersectionTypeEntry,
					poly2PointType: intersectionTypeExit,
				},
				false: {
					poly1PointType: intersectionTypeExit,
					poly2PointType: intersectionTypeEntry,
				},
			},
		},
	},
	BooleanIntersection: {
		PTSolid: {
			PTSolid: {
				true:  {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeExit},
				false: {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeEntry},
			},
			PTHole: {
				true:  {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeEntry},
				false: {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeExit},
			},
		},
		PTHole: {
			PTSolid: {
				true:  {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeExit},
				false: {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeEntry},
			},
			PTHole: {
				true:  {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeEntry},
				false: {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeExit},
			},
		},
	},
	BooleanSubtraction: {
		PTSolid: {
			PTSolid: {
				true:  {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeExit},
				false: {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeEntry},
			},
			PTHole: {
				true:  {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeEntry},
				false: {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeExit},
			},
		},
		PTHole: {
			PTSolid: {
				true:  {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeExit},
				false: {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeEntry},
			},
			PTHole: {
				true:  {poly1PointType: intersectionTypeEntry, poly2PointType: intersectionTypeEntry},
				false: {poly1PointType: intersectionTypeExit, poly2PointType: intersectionTypeExit},
			},
		},
	},
}

// BooleanOperation defines the types of Boolean operations that can be performed on polygons.
// These operations are fundamental in computational geometry for combining or modifying shapes.
//
// The supported operations are:
// - Union: Combines two polygons into one, merging their areas.
// - Intersection: Finds the overlapping region between two polygons.
// - Subtraction: Subtracts one polygon's area from another.
type BooleanOperation uint8

// NewPolyTreeOption defines a functional option type for configuring a new PolyTree during creation.
//
// This type allows for flexible and extensible initialization of PolyTree objects by applying optional
// configurations after the core properties have been set.
//
// Parameters:
//   - T: The numeric type of the coordinates in the PolyTree, constrained by the SignedNumber interface.
//
// Example Usage:
//
//	// Define an option that adds children to a PolyTree
//	func WithChildren[T SignedNumber](children ...*PolyTree[T]) NewPolyTreeOption[T] {
//	    return func(pt *PolyTree[T]) {
//	        pt.children = append(pt.children, children...)
//	    }
//	}
//
//	// Create a new PolyTree with a child
//	child, _ := NewPolyTree([]Point[int]{{1, 1}, {2, 1}, {2, 2}, {1, 2}}, PTSolid)
//	parent, _ := NewPolyTree([]Point[int]{{0, 0}, {3, 0}, {3, 3}, {0, 3}}, PTSolid, WithChildren(child))
//
// This pattern makes it easy to add optional properties to a PolyTree without requiring an extensive list
// of parameters in the NewPolyTree function.
type NewPolyTreeOption[T SignedNumber] func(*PolyTree[T])

// PointPolygonRelationship (PPR) defines the possible spatial relationships between a point
// and a polygon, accounting for structures such as holes and nested islands.
//
// The relationships are enumerated as follows:
//   - PPRPointInside: The point lies strictly within the boundaries of the polygon.
//   - PPRPointOutside: The point lies outside the outermost boundary of the polygon.
//   - PPRPointOnVertex: The point coincides with a vertex of the polygon or one of its holes/islands.
//   - PPRPointOnEdge: The point lies exactly on an edge of the polygon or one of its holes/islands.
//   - PPRPointInHole: The point lies inside a hole within the polygon.
//   - PPRPointInsideIsland: The point lies inside an island nested within the polygon.
type PointPolygonRelationship int8

// PolygonType (PT) defines whether the inside of the contour of a polygon represents either a solid region (island)
// or a void region (hole). This distinction is essential for operations involving polygons
// with complex structures, such as those containing holes or nested islands.
type PolygonType uint8

// PolyTree represents a polygon that can be part of a hierarchical structure. It supports
// complex polygons with holes and islands and enables operations like unions, intersections,
// and subtractions.
//
// A PolyTree has the following characteristics:
//   - A contour: The outer boundary of the polygon, represented as a sequence of vertices.
//     This includes all the original points, intersection points, and added midpoints as needed
//     for polygon operations.
//   - A PolygonType: Indicates whether the polygon is a solid region (PTSolid) or a hole (PTHole).
//     This classification is essential for understanding the relationship between the polygon
//     and its children.
//   - A parent: Points to the parent polygon in the hierarchy. For example, a hole's parent
//     would be the solid polygon that contains it. If a polygon is the root polygon in the PolyTree, its parent is nil.
//   - Zero or more siblings: A list of sibling polygons that are not nested within each other but share
//     the same parent. Siblings must be of the same PolygonType.
//   - Zero or more children: A list of child polygons nested within this polygon. If the PolygonType is
//     PTSolid, the children are holes (PTHole). If the PolygonType is PTHole, the children
//     are solid islands (PTSolid).
//
// Hierarchy Rules:
//   - A solid polygon (PTSolid) can contain holes (PTHole) as its children.
//   - A hole (PTHole) can contain solid polygons (PTSolid) as its children.
//   - Siblings are polygons of the same PolygonType that do not overlap.
//
// These relationships form a tree structure, where each node is a PolyTree, allowing for
// complex geometric modeling and operations.
//
// Note: Internal optimizations, such as convex hull caching and point indexing, are abstracted
// away from the user.
type PolyTree[T SignedNumber] struct {

	// contour defines the complete outline of the polygon, including all vertices, intersection points,
	// and midpoints added during boolean operations. The points in the contour are doubled to avoid
	// precision issues (e.g., rounding errors or loss of accuracy) when calculating midpoints, such as
	// during entry/exit point determination for boolean operations. This ensures accurate and consistent
	// results, especially when working with integer-based coordinates.
	contour contour[T]

	// polygonType specifies whether the polygon is solid (PTSolid) or a hole (PTHole).
	// This determines the interpretation of its children: solids contain holes, and holes contain solids.
	polygonType PolygonType

	// siblings stores polygons of the same type that share the same parent. This maintains
	// adjacency relationships between polygons without nesting.
	siblings []*PolyTree[T]

	// children contains polygons nested within this polygon. These are either holes (for solids)
	// or solid islands (for holes), forming a hierarchical structure.
	children []*PolyTree[T]

	// parent refers to the polygon that contains this one, maintaining the tree structure.
	// A nil parent indicates that this polygon is the root of the hierarchy.
	parent *PolyTree[T]

	// hull caches the convex hull of the polygon for faster point-in-polygon checks.
	hull simpleConvexPolygon[T]

	// maxX stores the maximum X-coordinate among the polygon's vertices. Used for ray-casting
	// in point-in-polygon checks and other spatial queries.
	maxX T
}

// PolyTreeMismatch represents a bitmask of potential mismatches between two PolyTree structures.
type PolyTreeMismatch uint8

// contour represents the outline of a polygon as a slice of polyTreePoint entries.
// Each entry contains metadata about the point, such as whether it is a normal vertex,
// an intersection point, or a midpoint between intersections.
//
// The contour is used to define the polygon's shape and is processed during boolean
// operations. Points within the contour are typically doubled to facilitate calculations
// involving midpoints and to avoid precision issues when working with integer-based
// coordinates.
type contour[T SignedNumber] []polyTreePoint[T]

// polyEdge represents an edge of a polygon, storing the geometric line segment
// and additional metadata for polygon operations.
//
// This type is used internally for operations such as determining point-in-polygon
// relationships and handling ray intersection tests. It provides both the edge's
// geometric representation and its relationship with a ray used in algorithms.
type polyEdge[T SignedNumber] struct {
	// lineSegment represents the geometric edge of the polygon as a line segment.
	// This field is used for geometric operations such as intersection checks and edge traversal.
	lineSegment LineSegment[T]

	// rel specifies the relationship of this edge with a ray during point-in-polygon tests.
	// This field is primarily used for algorithms like ray-casting to determine whether
	// a point is inside or outside the polygon.
	rel LineSegmentsRelationship
}

// polyIntersectionType defines the type of intersection point in polygon operations,
// distinguishing between entry and exit (of area of interest) points during traversal.
//
// This type is primarily used in Boolean operations (e.g., union, intersection, subtraction)
// to identify transitions at intersection points between polygons.
type polyIntersectionType int

// polyPointType defines the type of point in a polygon, used to distinguish between
// original vertices and additional points introduced during polygon operations.
//
// This type is essential for managing polygon data during Boolean operations (e.g., union,
// intersection, subtraction) and other algorithms that require distinguishing original
// points from dynamically added points.
type polyPointType uint8

// polyTraversalDirection defines the direction in which a polygon's vertices are traversed.
// This can either be clockwise or counterclockwise and is used to specify how to iterate
// through a polygon's vertices or edges during boolean operations and other processing.
type polyTraversalDirection uint8

// polyTreeEqOption defines a functional option for configuring the behaviour of the Eq method
// on a PolyTree. These options are used to customize the comparison logic, such as whether to
// track visited nodes to prevent infinite recursion during comparisons of siblings and children.
//
// For example, `WithVisited` can be used to provide a map of visited nodes to avoid reprocessing
// already-checked polygons in recursive structures.
type polyTreeEqOption[T SignedNumber] func(*polyTreeEqConfig[T])

// polyTreeEqConfig is a configuration struct used to control the behaviour of the Eq method
// on a PolyTree. It provides additional context and state to support complex comparisons,
// such as tracking visited nodes to prevent infinite recursion when comparing nested or cyclic
// structures.
//
// Fields:
//   - visited: A map used to track which PolyTree nodes have already been visited during the
//     comparison. This prevents infinite loops when comparing siblings or children that may
//     reference each other in complex hierarchies.
type polyTreeEqConfig[T SignedNumber] struct {
	visited map[*PolyTree[T]]bool // Tracks visited nodes to avoid infinite recursion.
}

// polyTreePoint represents a point in a polygon, with additional metadata used
// for polygon operations such as Boolean operations and traversal.
//
// Fields:
//   - point: The geometric coordinates of the point in 2D space.
//   - pointType: Indicates the type of the point, which can be a normal vertex,
//     an intersection point, or a midpoint between intersections.
//   - entryExit: Specifies whether this point is an entry or exit point during
//     polygon traversal in Boolean operations. This field is critical for determining
//     traversal directions and relationships between polygons.
//   - visited: Tracks whether this point has been visited during traversal algorithms,
//     helping to avoid redundant processing.
//   - intersectionPartner: A reference to the partner polygon involved in an intersection,
//     if this point is an intersection point. This field is nil for normal vertices.
//   - intersectionPartnerPointIndex: The index of the corresponding intersection point in
//     the partner polygon's contour, if this point is an intersection. A value of -1
//     indicates no partner exists.
//
// Usage:
//
// This struct is primarily used in [PolyTree]'s contour to represent the points
// and their metadata, enabling advanced polygon operations such as union, intersection,
// and subtraction.
type polyTreePoint[T SignedNumber] struct {
	// The geometric coordinates of the point in 2D space.
	point Point[T]

	// The type of the point, such as a normal vertex, intersection point, or midpoint.
	pointType polyPointType

	// Indicates whether this point is an entry or exit point during traversal.
	entryExit polyIntersectionType

	// Tracks whether this point has been visited during traversal algorithms.
	visited bool

	// Reference to the partner polygon for intersection points.
	intersectionPartner *PolyTree[T]

	// Index of the corresponding intersection point in the partner polygon.
	intersectionPartnerPointIndex int
}

// simpleConvexPolygon represents a convex polygon, which is a polygon where all interior angles
// are less than 180 degrees, and no line segment between two points on the boundary extends outside the polygon.
//
// This type is used internally to optimize geometric operations such as point-in-polygon checks,
// as convex polygons allow for faster algorithms compared to general polygons.
//
// As this is used internally, no checks are in place to enforce convexity.
// The ConvexHull function returns this type.
type simpleConvexPolygon[T SignedNumber] struct {
	// Points contains the ordered vertices of the convex polygon. The points are arranged
	// sequentially in either clockwise or counterclockwise order, forming the boundary of the convex hull.
	Points []Point[T]
}

// NewPolyTree creates a new [PolyTree] object from a set of points defining the polygon's contour.
// The function also allows for optional configuration using [NewPolyTreeOption].
//
// Parameters:
//   - points: A slice of Point[T] representing the vertices of the polygon.
//   - t: The type of polygon, either PTSolid or PTHole.
//   - opts: Optional function configurations applied to the resulting PolyTree.
//
// Returns:
//   - A pointer to the newly created PolyTree.
//   - An error if the input points are invalid (e.g., less than three points or zero area).
//
// Notes:
//   - The function ensures that the polygon's points are oriented correctly based on its type.
//   - Points are doubled internally to avoid integer division/precision issues during midpoint calculations.
//   - The polygon's convex hull is computed and stored for potential optimisations.
//   - Child polygons must have the opposite PolygonType (e.g., holes for a solid polygon and solids for a hole polygon).
func NewPolyTree[T SignedNumber](points []Point[T], t PolygonType, opts ...NewPolyTreeOption[T]) (*PolyTree[T], error) {

	// Sanity check: A polygon must have at least three points.
	if len(points) < 3 {
		return nil, fmt.Errorf("new polytree must have at least 3 points")
	}

	// Sanity check: A polygon must have a non-zero area.
	if SignedArea2X(points) == 0 {
		return nil, fmt.Errorf("new polytree must have non-zero area")
	}

	// Create a new, zero-initialised PolyTree.
	p := new(PolyTree[T])

	// Set the polygon type (e.g., solid or hole).
	p.polygonType = t

	// Ensure the points are oriented correctly based on the polygon type.
	orderedPoints := make([]Point[T], len(points))
	copy(orderedPoints, points)
	switch p.polygonType {
	case PTSolid:
		// Solid polygons must be counterclockwise.
		EnsureCounterClockwise(orderedPoints)
	case PTHole:
		// Hole polygons must be clockwise.
		EnsureClockwise(orderedPoints)
	}

	// Assign and double the points for precision handling.
	p.maxX = points[0].x * 2
	p.contour = make([]polyTreePoint[T], len(orderedPoints))
	for i, point := range orderedPoints {
		p.contour[i] = polyTreePoint[T]{
			point: NewPoint(point.x*2, point.y*2),
		}
		p.maxX = max(p.maxX, p.contour[i].point.x)
	}

	// Reorder contour points and reset intersection metadata.
	p.resetIntersectionMetadataAndReorder()

	// Increment maxX to ensure boundary checks in ray-casting operations.
	p.maxX++

	// Compute and store the convex hull of the polygon for optimisation.
	hull := ConvexHull(points)
	EnsureCounterClockwise(hull) // Convex hulls are always counterclockwise.
	p.hull = newSimpleConvexPolygon(hull)

	// Apply optional configurations.
	for _, opt := range opts {
		opt(p)
	}

	// Sanity check: Ensure all children have the opposite PolygonType.
	for _, c := range p.children {
		switch p.polygonType {
		case PTSolid:
			// Solid polygons can only have hole children.
			if c.polygonType != PTHole {
				return nil, fmt.Errorf("expected all children to have PolygonType PTHole")
			}
		case PTHole:
			// Hole polygons can only have solid children.
			if c.polygonType != PTSolid {
				return nil, fmt.Errorf("expected all children to have PolygonType PTSolid")
			}
		}
	}

	// Return the constructed PolyTree.
	return p, nil
}

// newSimpleConvexPolygon creates a new simpleConvexPolygon from a given slice of points.
// The input points are assumed to be ordered to form a convex polygon.
//
// This function is primarily used internally to construct a convex polygon representation
// for optimization purposes, such as in point-in-polygon checks.
//
// Parameters:
//
//   - points []Point[T]: A slice of points representing the vertices of the convex polygon.
//     The points are assumed to be ordered sequentially, either clockwise
//     or counterclockwise, to define the boundary of a convex polygon.
//
// Returns:
//
//   - *simpleConvexPolygon[T]: A pointer to a new simpleConvexPolygon containing the provided points.
//
// Notes:
//   - This function assumes that the input points are already ordered and form a valid convex polygon.
//     No validation is performed to verify the convexity of the polygon or the order of the points.
//   - As this function is used internally, it expects the caller to ensure the validity of the input.
//
// Example:
//
//	points := []Point[float64]{
//	    {X: 0, Y: 0},
//	    {X: 4, Y: 0},
//	    {X: 4, Y: 4},
//	    {X: 0, Y: 4},
//	}
//	scp := newSimpleConvexPolygon(points)
//	// scp represents a convex polygon with the given points.
func newSimpleConvexPolygon[T SignedNumber](points []Point[T]) simpleConvexPolygon[T] {
	// Assume `points` is already ordered to form a convex polygon
	return simpleConvexPolygon[T]{Points: points}
}

// WithChildren is an option for the NewPolyTree function that assigns child polygons to the created PolyTree.
// It also sets up parent-child relationships and orders the children for consistency.
//
// Parameters:
//   - children: A variadic list of pointers to PolyTree objects representing the child polygons.
//
// Behavior:
//   - The function assigns the provided children to the PolyTree being created.
//   - It establishes the parent-child relationship by setting the parent of each child to the newly created PolyTree.
//
// Returns:
//   - A NewPolyTreeOption that can be passed to the NewPolyTree function.
func WithChildren[T SignedNumber](children ...*PolyTree[T]) NewPolyTreeOption[T] {
	return func(p *PolyTree[T]) {

		// Assign the provided children to the parent polygon.
		p.children = children

		// Order the children for consistency in traversal and comparison.
		p.orderSiblingsAndChildren()

		// Set the parent field of each child to the current polygon.
		for i := range p.children {
			p.children[i].parent = p
		}
	}
}

// contains checks whether a given point is present in the contour.
// It iterates through the contour to find a matching point by comparing
// their x and y coordinates.
//
// Parameters:
//   - point: The Point to check for existence in the contour.
//
// Returns:
//   - true if the given point exists in the contour, otherwise false.
func (c *contour[T]) contains(point Point[T]) bool {
	// Use slices.ContainsFunc to determine if the given point exists in the contour.
	return slices.ContainsFunc(*c, func(p polyTreePoint[T]) bool {
		// Compare the x and y coordinates of the current point in the contour with the given point.
		// Return true if they match, indicating the point exists in the contour.
		if p.point.x == point.x && p.point.y == point.y {
			return true
		}
		// Otherwise, return false to continue searching.
		return false
	})
}

// eq compares two contours to determine if they are equivalent.
// Contours are considered equal if they contain the same points in the same relative order,
// accounting for possible rotations and ensuring consistent orientation (e.g., clockwise or counterclockwise).
//
// Parameters:
//   - other: The contour to compare with the current contour.
//
// Returns:
//   - true if the contours are equivalent, otherwise false.
func (c *contour[T]) eq(other contour[T]) bool {

	// Check if both contours are empty; two empty contours are considered equal.
	if len(*c) == 0 && len(other) == 0 {
		return true
	}

	// If the lengths of the contours differ, they cannot be equal.
	if len(*c) != len(other) {
		return false
	}

	// Create copies of the contours to avoid modifying the originals during comparison.
	copyC := c.toPoints()
	copyOther := other.toPoints()

	// Ensure both contours have the same orientation.
	// Reverse the order of points if the signed area indicates inconsistent orientation.
	if SignedArea2X(copyC) > 0 {
		slices.Reverse(copyC)
	}
	if SignedArea2X(copyOther) > 0 {
		slices.Reverse(copyOther)
	}

	// Attempt to match contours by rotating the starting point of the second contour.
	for start := 0; start < len(copyOther); start++ {
		matches := true

		// Check if this rotation aligns all points in both contours.
		for i := 0; i < len(copyC); i++ {
			if copyC[i] != copyOther[(i+start)%len(copyOther)] {
				// If a mismatch is found, mark as not matching and break out of the loop.
				matches = false
				break
			}
		}

		// If all points match for this rotation, the contours are equivalent.
		if matches {
			return true
		}
	}

	// No rotation of the second contour resulted in a match, so the contours are not equivalent.
	return false
}

// findLowestLeftmost identifies the lowest, leftmost point in a given contour.
//
// This method iterates through all points in the contour to find the point
// with the smallest y coordinate. If multiple points share the same y
// coordinate, it selects the one with the smallest x coordinate.
//
// Parameters:
//   - c: A pointer to the contour, which is a slice of polyTreePoint.
//
// Returns:
//   - The Point with the lowest y coordinate and, in case of ties, the
//     smallest x coordinate.
func (c *contour[T]) findLowestLeftmost() Point[T] {
	// Initialize the minimum point as the first point in the contour.
	minimum := (*c)[0].point

	// Iterate through each point in the contour to find the lowest, leftmost point.
	for _, pt := range *c {
		// Update the minimum point if:
		// - The current point's y coordinate is smaller.
		// - Or, if the y coordinates are equal, the current point's x coordinate is smaller.
		if pt.point.y < minimum.y || (pt.point.y == minimum.y && pt.point.x < minimum.x) {
			minimum = pt.point // Update the minimum to the current point.
		}
	}

	// Return the identified lowest, leftmost point.
	return minimum
}

// insertIntersectionPoint inserts an intersection point into a contour between two specified indices.
// The insertion position is determined based on the proximity of the intersection point to the
// corresponding line segment, ensuring correct ordering for polygon operations.
//
// Parameters:
//   - start: The index of the starting point of the line segment in the contour.
//   - end: The index of the ending point of the line segment in the contour.
//   - intersection: The intersection point to insert, represented as a polyTreePoint.
//
// Notes:
//   - The function assumes that start and end indices are valid and start < end.
//   - The insertion logic ensures that the intersection is placed in the correct position relative to
//     other intermediate points, preserving the geometric consistency of the contour.
//   - Simply inserting at the start index is not sufficient because the contour may already contain
//     intermediate points between start and end. Proper ordering is necessary to maintain the
//     validity of polygon operations such as traversals and Boolean operations.
func (c *contour[T]) insertIntersectionPoint(start, end int, intersection polyTreePoint[T]) {
	// Define the line segment between the start and end points.
	segment := NewLineSegment((*c)[start].point, (*c)[end].point)

	// Initialize the insertion position to the end index.
	insertPos := end

	// Iterate through the intermediate points in the contour between start and end.
	// Find the correct position to insert the intersection point.
	for i := start + 1; i < end; i++ {
		// Define the segment from the start point to the current point.
		existingSegment := NewLineSegment((*c)[start].point, (*c)[i].point)

		// Compare the distance of the intersection to the original segment
		// with its distance to the intermediate segment.
		if segment.DistanceToPoint(intersection.point) < existingSegment.DistanceToPoint((*c)[i].point) {
			// Update the insertion position if the intersection is closer to the original segment.
			insertPos = i
			break
		}
	}

	// Insert the intersection point into the contour at the calculated position.
	*c = slices.Insert(*c, insertPos, intersection)
}

// isContourInside checks if all points of one contour (c2) are inside another contour (c).
//
// Parameters:
//   - c2: The contour to check against the current contour (c).
//
// Returns:
//   - true: If all points of c2 are inside c.
//   - false: If any point of c2 is outside c.
//
// Notes:
//   - This function iterates over each point in c2 and checks whether it is inside c using isPointInside.
//   - If any point of c2 is found to be outside c, the function returns false immediately.
//   - The function assumes that the orientation of c and c2 is correct (i.e., counter-clockwise for solids and clockwise for holes).
func (c *contour[T]) isContourInside(other contour[T]) bool {
	// Iterate over each point in other.
	for _, p := range other {
		// Check if the point is not inside c.
		if !c.isPointInside(p.point) {
			// If any point is outside, return false immediately.
			return false
		}
	}

	// If all points are inside, return true.
	return true
}

// isPointInside determines if a given point is inside the contour.
//
// Parameters:
//   - point: The point to check.
//
// Returns:
//   - true: If the point is inside the contour.
//   - false: If the point is outside the contour.
//
// Notes:
//   - This function uses the ray-casting algorithm to determine the point's relationship to the contour.
//   - A horizontal ray is cast from the point to the right, and the number of times it crosses contour edges is counted.
//   - If the number of crossings is odd, the point is inside; if even, the point is outside.
//   - Points lying directly on the contour edges are considered inside.
func (c *contour[T]) isPointInside(point Point[T]) bool {
	crosses := 0

	// Calculate the farthest x-coordinate to cast a ray.
	// todo: is this required? possibly redundant code...
	maxX := point.x
	for _, p := range *c {
		maxX = max(maxX, p.point.x)
	}
	// Ensure the ray extends beyond the contour's bounds.
	maxX++
	ray := NewLineSegment(point, NewPoint(maxX, point.y))

	// Convert the contour to edges (line segments).
	edges := c.toEdges()

	// Determine the relationship of each edge to the ray.
	for i := range edges {
		// If the point lies directly on the edge, consider it inside.
		if point.IsOnLineSegment(edges[i].lineSegment) {
			return true
		}

		// Determine the relationship of the edge to the ray.
		edges[i].rel = ray.RelationshipToLineSegment(edges[i].lineSegment)
	}

	// Analyze the relationships and count ray crossings.
	for i := range edges {
		// Look at the next edge in the contour.
		iNext := (i + 1) % len(edges)

		switch edges[i].rel {
		case LSRIntersects: // Ray intersects the edge.
			crosses++

		case LSRCollinearCDinAB: // Ray is collinear with the edge and overlaps it.
			crosses += 2

		case LSRConAB: // Ray starts on the edge.
			crosses++

			// Handle potential overlaps with the next edge.
			if edges[iNext].rel == LSRDonAB {
				if inOrder(edges[i].lineSegment.start.y, point.y, edges[iNext].lineSegment.end.y) {
					crosses++
				}
			}

		case LSRDonAB: // Ray ends on the edge.
			crosses++

			// Handle potential overlaps with the next edge.
			if edges[iNext].rel == LSRConAB {
				if inOrder(edges[i].lineSegment.start.y, point.y, edges[iNext].lineSegment.end.y) {
					crosses++
				}
			}

		default:
			// No action for edges that don't interact with the ray.
		}
	}

	// A point is inside if the number of crossings is odd.
	return crosses%2 == 1
}

// iterEdges iterates over all edges (line segments) in the contour, invoking the provided yield function for each edge.
//
// Parameters:
//
// yield: A callback function that takes a LineSegment and returns a boolean.
//   - If yield returns true, iteration continues to the next edge.
//   - If yield returns false, iteration stops immediately.
//
// Behavior:
//   - Each edge is formed by connecting consecutive points in the contour.
//   - The contour is treated as a closed loop, so the last point is connected back to the first.
//   - If the contour contains fewer than two points, no edges can be formed, and the function exits without calling `yield`.
//
// Notes:
//   - This function avoids modifying the contour during iteration.
//   - It allows early termination of iteration.
//
// Example:
//
//	c := contour[int]{
//		{point: Point[int]{0, 0}},
//		{point: Point[int]{10, 0}},
//		{point: Point[int]{10, 10}},
//		{point: Point[int]{0, 10}},
//	}
//	fmt.Println("Edges in the contour:")
//	for edge := range c.iterEdges {
//		fmt.Printf("Edge: %v -> %v\n", edge.start, edge.end)
//	}
//
// Returns:
//
//	Edges in the contour:
//	Edge: Point[(0, 0)] -> Point[(10, 0)]
//	Edge: Point[(10, 0)] -> Point[(10, 10)]
//	Edge: Point[(10, 10)] -> Point[(0, 10)]
//	Edge: Point[(0, 10)] -> Point[(0, 0)]
func (c *contour[T]) iterEdges(yield func(LineSegment[T]) bool) {
	// A contour with fewer than two points cannot form edges.
	if len(*c) < 2 {
		return
	}

	// Iterate over all points in the contour.
	for i := range *c {
		j := (i + 1) % len(*c) // Wrap around to connect the last point to the first.

		// Construct a line segment for the current edge and pass it to the yield function.
		if !yield(NewLineSegment((*c)[i].point, (*c)[j].point)) {
			return // Exit early if `yield` signals to stop iteration.
		}
	}
}

// reorder adjusts the order of points in the contour such that the lowest, leftmost point
// becomes the starting point. This ensures a consistent and predictable point order for
// comparisons and processing. The contour is modified in place.
//
// The "lowest, leftmost point" is the point with the smallest y-coordinate. If there are
// multiple points with the same y-coordinate, the one with the smallest x-coordinate is chosen.
//
// Example:
//
//	Before: [(3, 4), (1, 1), (2, 2), (0, 1)]
//	After:  [(0, 1), (3, 4), (1, 1), (2, 2)]
func (c *contour[T]) reorder() {
	// Find the index of the lowest, leftmost point
	minIndex := 0
	for i := 1; i < len(*c); i++ {
		// Compare the current point with the current minimum
		if (*c)[i].point.y < (*c)[minIndex].point.y || // Lower y-coordinate
			((*c)[i].point.y == (*c)[minIndex].point.y && // Same y-coordinate, check x-coordinate
				(*c)[i].point.x < (*c)[minIndex].point.x) {
			minIndex = i // Update the index of the minimum point
		}
	}

	// Rotate the contour to make the lowest, leftmost point the starting point
	// The rotation is achieved by slicing the contour and appending the two slices.
	rotated := append((*c)[minIndex:], (*c)[:minIndex]...)

	// Copy the rotated points back to the original contour
	copy(*c, rotated)
}

// toEdges converts the contour into a slice of polyEdge objects, where each edge
// represents a line segment connecting two consecutive points in the contour.
// The contour is treated as a closed loop, so the last point connects back to the first.
//
// This function leverages the iterEdges method to iterate through the edges of the contour,
// making it concise and consistent with edge iteration logic.
//
// Returns:
//   - A slice of polyEdge objects, with each edge containing a LineSegment.
//     The slice will contain the same number of edges as the number of points in the contour.
func (c *contour[T]) toEdges() []polyEdge[T] {
	// Preallocate the slice of edges to match the number of points in the contour
	edges := make([]polyEdge[T], 0, len(*c))

	// Use iterEdges to process each edge in the contour
	c.iterEdges(func(edge LineSegment[T]) bool {
		edges = append(edges, polyEdge[T]{lineSegment: edge})
		return true // Continue processing all edges
	})

	return edges
}

// toPoints extracts the original points from the contour, excluding any added points
// (e.g., intersection or midpoint points). The points are halved to reverse the doubling
// applied during the contour's creation, restoring their original values.
//
// This function ensures that only points marked as pointTypeOriginal are included
// in the resulting slice, which represents the true vertices of the polygon.
//
// Returns:
//   - A slice of Point[T] containing the original, halved points.
func (c *contour[T]) toPoints() []Point[T] {
	// Preallocate the slice to avoid unnecessary reallocations
	originalPoints := make([]Point[T], 0, len(*c))

	// Iterate over the contour to extract only the original points
	for _, p := range *c {
		// Include only points marked as pointTypeOriginal
		if p.pointType == pointTypeOriginal {
			// Halve the x and y coordinates to reverse the earlier doubling
			originalPoints = append(originalPoints, NewPoint[T](
				p.point.x/2,
				p.point.y/2,
			))
		}
	}

	return originalPoints
}

func (p *polyIntersectionType) String() string {
	switch *p {
	case intersectionTypeNotSet:
		return "not set"
	case intersectionTypeEntry:
		return "entry"
	case intersectionTypeExit:
		return "exit"
	}
	return ""
}

// BooleanOperation performs a Boolean operation (union, intersection, or subtraction)
// between the current polygon (p) and another polygon (other). The result is
// returned as a new PolyTree, or an error is returned if the operation fails.
//
// Parameters:
//   - other: The polygon to perform the operation with.
//   - operation: The type of Boolean operation to perform (e.g., union, intersection, subtraction).
//
// Returns:
//   - A new PolyTree resulting from the operation, or an error if the operation fails.
//
// Behavior:
//
// For non-intersecting polygons:
//   - Union: Adds the other polygon as a sibling of p.
//   - Intersection: Returns nil as there is no overlapping area.
//   - Subtraction: Returns p unchanged, as other does not overlap with p.
//
// For intersecting polygons:
//   - Computes intersection points and marks entry/exit points for traversal.
//   - Traverses the polygons to construct the result of the operation.
//   - Returns a nested PolyTree structure representing the operation result.
func (p *PolyTree[T]) BooleanOperation(other *PolyTree[T], operation BooleanOperation) (*PolyTree[T], error) {
	// Edge Case: Check if the polygons intersect
	if !p.Intersects(other) {
		switch operation {
		case BooleanUnion:
			// Non-intersecting polygons: Add other as a sibling
			if err := p.addSibling(other); err != nil {
				return nil, fmt.Errorf("failed to add sibling: %w", err)
			}
			return p, nil

		case BooleanIntersection:
			// Non-intersecting polygons: No intersection, return nil
			return nil, nil

		case BooleanSubtraction:
			// Non-intersecting polygons: No change to p
			return p, nil

		default:
			// Invalid or unsupported operation
			return nil, fmt.Errorf("unknown operation: %v", operation)
		}
	}

	// Step 1: Find intersection points between all polygons
	p.findIntersections(other)

	// Step 2: Mark entry/exit points for traversal based on the operation
	p.markEntryExitPoints(other, operation)

	// Step 3: Perform traversal to construct the result of the Boolean operation
	return nestPointsToPolyTrees(p.booleanOperationTraversal(other, operation))
}

// Eq compares two PolyTree objects (p and other) for structural and content equality.
// It identifies mismatches in contours, siblings, and children while avoiding infinite recursion
// by tracking visited nodes. The comparison results are represented as a boolean and a bitmask.
//
// Parameters:
//   - other: The PolyTree to compare with the current PolyTree.
//   - opts: Optional configurations for the comparison, such as tracking visited nodes.
//
// Returns:
//   - A boolean indicating whether the two PolyTree objects are equal.
//   - A PolyTreeMismatch bitmask that specifies which aspects (if any) differ.
//
// Behavior:
//   - Handles nil cases for the p and other PolyTree objects.
//   - Recursively compares contours, siblings, and children, considering different sibling/child orders.
//   - Uses the visited map to prevent infinite recursion in cyclic structures.
func (p *PolyTree[T]) Eq(other *PolyTree[T], opts ...polyTreeEqOption[T]) (bool, PolyTreeMismatch) {
	var mismatches PolyTreeMismatch

	// Handle nil cases
	switch {
	case p == nil && other == nil:
		return true, PTMNoMismatch // Both are nil, no mismatch
	case p == nil && other != nil:
		return false, PTMNilPolygonMismatch // One is nil, the other is not
	case p != nil && other == nil:
		return false, PTMNilPolygonMismatch // One is nil, the other is not
	}

	// Initialize comparison configuration
	config := &polyTreeEqConfig[T]{
		visited: make(map[*PolyTree[T]]bool),
	}
	for _, opt := range opts {
		opt(config) // Apply configuration options
	}

	// Avoid comparing the same PolyTree objects multiple times
	if config.visited[p] || config.visited[other] {
		return true, PTMNoMismatch // Already compared, no mismatch
	}
	config.visited[p] = true
	config.visited[other] = true

	// Compare contours
	if !p.contour.eq(other.contour) {
		mismatches |= PTMContourMismatch
		fmt.Println("Mismatch in contours detected")
	}

	// Compare siblings
	if len(p.siblings) != len(other.siblings) {
		mismatches |= PTMSiblingMismatch
		fmt.Println("Mismatch in siblings count detected")
	} else {
		// Siblings order doesn't matter, check if each sibling matches
		for _, sibling := range p.siblings {
			found := false
			for _, otherSibling := range other.siblings {
				if eq, _ := sibling.Eq(otherSibling, withVisited(config.visited)); eq {
					found = true
					break
				}
			}
			if !found {
				mismatches |= PTMSiblingMismatch
				fmt.Println("Mismatch in siblings content detected")
				break
			}
		}
	}

	// Compare children
	if len(p.children) != len(other.children) {
		mismatches |= PTMChildMismatch
		fmt.Println("Mismatch in children count detected")
	} else {
		// Children order doesn't matter, check if each child matches
		for _, child := range p.children {
			found := false
			for _, otherChild := range other.children {
				if eq, _ := child.Eq(otherChild, withVisited(config.visited)); eq {
					found = true
					break
				}
			}
			if !found {
				mismatches |= PTMChildMismatch
				fmt.Println("Mismatch in children content detected")
				break
			}
		}
	}

	// Return whether there are no mismatches and the mismatch bitmask
	return mismatches == PTMNoMismatch, mismatches
}

// Intersects checks whether the current PolyTree intersects with another PolyTree.
//
// Parameters:
//   - other: The PolyTree to compare against the current PolyTree.
//
// Returns:
//   - true if the two PolyTree objects intersect, either by containment or edge overlap.
//   - false if there is no intersection.
//
// Behavior:
//   - Checks if any point from one PolyTree lies inside the contour of the other PolyTree.
//   - Verifies if any edges of the two PolyTree objects intersect.
//
// This method accounts for all potential intersection cases, including:
//   - One polygon entirely inside the other.
//   - Overlapping polygons with shared edges.
//   - Polygons touching at vertices or along edges.
func (p *PolyTree[T]) Intersects(other *PolyTree[T]) bool {
	// Check if any point of "other" is inside the contour of "p"
	for _, otherPoint := range other.contour {
		if p.contour.isPointInside(otherPoint.point) {
			return true // Intersection found via point containment
		}
	}

	// Check if any point of "p" is inside the contour of "other"
	for _, point := range p.contour {
		if other.contour.isPointInside(point.point) {
			return true // Intersection found via point containment
		}
	}

	// Check for edge intersections between "p" and "other"
	for poly1Edge := range p.contour.iterEdges {
		for poly2Edge := range other.contour.iterEdges {
			if poly1Edge.IntersectsLineSegment(poly2Edge) {
				return true // Intersection found via edge overlap
			}
		}
	}

	// No intersections detected
	return false
}

// addChild adds a child PolyTree to the current PolyTree.
//
// Parameters:
// - child: A pointer to the PolyTree to be added as a child.
//
// Returns:
// - error: An error if the operation fails. Possible error scenarios include:
//   - The child is nil.
//   - The child has the same polygonType as the parent.
//
// Behavior:
//   - Validates that the child is not nil.
//   - Ensures that the PolygonType of the child is the opposite of the parent's. A PTSolid parent can only have PTHole children. A PTHole parent can only have PTSolid children.
//   - Sets the current PolyTree as the parent of the child.
//   - Adds the child to the children slice and ensures proper ordering of siblings and children.
//
// todo: make public?
func (p *PolyTree[T]) addChild(child *PolyTree[T]) error {
	// Check if the child is nil
	if child == nil {
		return fmt.Errorf("attempt to add nil child")
	}

	// Ensure the polygon types are compatible
	if p.polygonType == child.polygonType {
		return fmt.Errorf(
			"cannot add child: mismatched polygon types (parent: %v, child: %v)",
			p.polygonType,
			child.polygonType,
		)
	}

	// Set the parent of the child
	child.parent = p

	// Append the child to the children slice
	p.children = append(p.children, child)

	// Order siblings and children for consistency
	p.orderSiblingsAndChildren()

	// Successfully added the child
	return nil
}

// addSibling adds a sibling PolyTree to the current PolyTree.
//
// Parameters:
//   - sibling: A pointer to the PolyTree to be added as a sibling.
//
// Returns:
//
// error: An error if the operation fails. Possible error scenarios include:
//   - The sibling is nil.
//   - The sibling has a different PolygonType than the current PolyTree.
//
// Behavior:
//   - Validates that the sibling is not nil.
//   - Ensures that the PolygonType of the sibling matches the current PolyTree.
//   - Establishes sibling relationships between the current PolyTree and the new sibling,
//     as well as among existing siblings, ensuring a consistent sibling list across all related PolyTree nodes.
//   - Calls orderSiblingsAndChildren on all affected nodes to maintain consistent ordering.
func (p *PolyTree[T]) addSibling(sibling *PolyTree[T]) error {
	// Check if the sibling is nil
	if sibling == nil {
		return fmt.Errorf("attempt to add nil sibling")
	}

	// Ensure the polygon types match
	if p.polygonType != sibling.polygonType {
		return fmt.Errorf("cannot add sibling as polygonType is mismatched")
	}

	// Add the new sibling to the sibling lists of existing siblings
	for _, existingSibling := range p.siblings {
		// Update sibling relationships
		existingSibling.siblings = append(existingSibling.siblings, sibling)

		// Maintain consistent ordering
		existingSibling.orderSiblingsAndChildren()

		// Add the existing sibling to the new sibling's sibling list
		sibling.siblings = append(sibling.siblings, existingSibling)
	}

	// Add the current `PolyTree` to the new sibling's sibling list
	sibling.siblings = append(sibling.siblings, p)

	// Maintain consistent ordering
	sibling.orderSiblingsAndChildren()

	// Add the new sibling to the current `PolyTree`'s sibling list
	p.siblings = append(p.siblings, sibling)

	// Maintain consistent ordering
	p.orderSiblingsAndChildren()

	return nil
}

// booleanOperationTraversal performs a Boolean operation (Union, Intersection, Subtraction) between
// two PolyTrees, traversing their contours to generate the resulting contours.
//
// Parameters:
//   - other: The PolyTree to operate with.
//   - operation: The type of Boolean operation (Union, Intersection, or Subtraction).
//
// Returns:
//   - A slice of slices of Point[T], where each inner slice represents a resulting contour from the operation.
//
// Behavior:
//   - Handles edge cases where polygons do not intersect.
//   - Traverses the contours of both PolyTrees, marking entry and exit points based on the Boolean operation.
//   - Constructs the resulting contours by switching between polygons at intersection points.
//
// Note:
//   - Polygons are assumed to have properly marked entry and exit points before calling this function.
//   - Points are halved during traversal to revert the earlier doubling for precision handling.
func (p *PolyTree[T]) booleanOperationTraversal(other *PolyTree[T], operation BooleanOperation) [][]Point[T] {
	var direction polyTraversalDirection

	// Initialize the resulting contours
	resultContours := make([][]Point[T], 0)

	for {
		// Find the starting point for traversal
		currentPoly, currentPointIndex := p.findTraversalStartingPoint(other)
		if currentPoly == nil || currentPointIndex == -1 {
			// No unvisited entry points, traversal is complete
			break
		}

		// Initialize a new contour for the result
		resultContour := make([]Point[T], 0, len(p.contour)+len(other.contour))

		// Set the initial traversal direction
		direction = polyTraversalForward

		for {
			// Append the current point to the resultContour (halving coordinates to undo doubling)
			resultContour = append(resultContour, NewPoint(
				currentPoly.contour[currentPointIndex].point.x/2,
				currentPoly.contour[currentPointIndex].point.y/2,
			))

			// Mark the current point and its intersection partner (if any) as visited
			currentPoly.contour[currentPointIndex].visited = true
			if currentPoly.contour[currentPointIndex].intersectionPartner != nil {
				partnerPoly := currentPoly.contour[currentPointIndex].intersectionPartner
				partnerPointIndex := currentPoly.contour[currentPointIndex].intersectionPartnerPointIndex
				partnerPoly.contour[partnerPointIndex].visited = true
			}

			// Move to the next point in the current polygon based on the direction
			if direction == polyTraversalForward {
				currentPointIndex = (currentPointIndex + 1) % len(currentPoly.contour)
			} else {
				currentPointIndex = (currentPointIndex - 1 + len(currentPoly.contour)) % len(currentPoly.contour)
			}

			// Handle polygon switching at entry/exit points and adjust traversal direction if needed
			if currentPoly.contour[currentPointIndex].entryExit == intersectionTypeExit ||
				(operation == BooleanSubtraction &&
					((currentPoly.polygonType == PTHole && currentPoly.contour[currentPointIndex].entryExit == intersectionTypeEntry) ||
						(currentPoly.polygonType == PTSolid && currentPoly.contour[currentPointIndex].entryExit == intersectionTypeEntry && direction == polyTraversalReverse))) {

				// Swap to the partner polygon and update the point index
				newCurrentPointIndex := currentPoly.contour[currentPointIndex].intersectionPartnerPointIndex
				currentPoly = currentPoly.contour[currentPointIndex].intersectionPartner
				currentPointIndex = newCurrentPointIndex

				// Reverse traversal direction for subtraction operations
				if operation == BooleanSubtraction {
					direction = togglePolyTraversalDirection(direction)
				}
			}

			// Terminate the loop when we return to the starting point in the same polygon
			if currentPoly.contour[currentPointIndex].point.x/2 == resultContour[0].x &&
				currentPoly.contour[currentPointIndex].point.y/2 == resultContour[0].y {
				break
			}
		}

		// Append the completed contour to the results
		resultContours = append(resultContours, resultContour)
	}

	// Handle edge cases: No intersections found
	if len(resultContours) == 0 {
		switch operation {
		case BooleanUnion:
			return [][]Point[T]{
				p.contour.toPoints(),
				other.contour.toPoints(),
			}
		case BooleanIntersection:
			return nil // No intersection
		case BooleanSubtraction:
			return [][]Point[T]{
				p.contour.toPoints(),
			}
		default:
			panic(fmt.Errorf("unknown BooleanOperation: %v", operation))
		}
	}

	return resultContours
}

// findIntersections identifies intersection points between the contours of two PolyTrees
// and adds these intersection points to the respective contours of both polygons.
//
// Parameters:
// - other: The PolyTree to check for intersections with the current PolyTree.
//
// Behavior:
//   - Resets intersection metadata and reorders both PolyTrees to ensure consistent traversal.
//   - Iterates through all combinations of polygons and edges between the two PolyTrees.
//   - For each intersecting edge pair, calculates the intersection point, ensures it is unique,
//     and inserts it into the contours of both PolyTrees.
//
// Note:
//   - This function modifies the contours of both PolyTrees by inserting intersection points.
//   - The inserted intersection points are marked as pointTypeAddedIntersection.
//
// Assumptions:
//   - The input PolyTrees represent closed polygons with properly ordered contours.
func (p *PolyTree[T]) findIntersections(other *PolyTree[T]) {

	// Step 1: Reset intersection metadata and reorder both PolyTrees
	p.resetIntersectionMetadataAndReorder()
	other.resetIntersectionMetadataAndReorder()

	// Step 2: Iterate through all combinations of polygons
	for poly1 := range p.iterPolys {
		for poly2 := range other.iterPolys {

			// Step 3: Check all edge combinations between poly1 and poly2
			for i1 := 0; i1 < len(poly1.contour); i1++ {
				// Get the next index to form an edge in poly1
				j1 := (i1 + 1) % len(poly1.contour)
				segment1 := NewLineSegment(poly1.contour[i1].point, poly1.contour[j1].point)

				for i2 := 0; i2 < len(poly2.contour); i2++ {
					// Get the next index to form an edge in poly2
					j2 := (i2 + 1) % len(poly2.contour)
					segment2 := NewLineSegment(poly2.contour[i2].point, poly2.contour[j2].point)

					// Step 4: Check for intersection between segment1 and segment2
					intersectionPoint, intersects := segment1.IntersectionPoint(segment2)
					if intersects {
						intersectionPointT := NewPoint(T(intersectionPoint.x), T(intersectionPoint.y))

						// Step 5: Ensure the intersection point is unique
						if poly1.contour.contains(intersectionPointT) || poly2.contour.contains(intersectionPointT) {
							continue // Skip duplicate intersections
						}

						// Step 6: Convert the intersection point to a polyTreePoint
						intersection := polyTreePoint[T]{
							point:                         intersectionPointT,
							pointType:                     pointTypeAddedIntersection, // Mark as intersection
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						}

						// Step 7: Insert the intersection point into both polygons
						poly1.contour.insertIntersectionPoint(i1, j1, intersection)
						poly2.contour.insertIntersectionPoint(i2, j2, intersection)

						// Step 8: Increment indices to avoid re-processing the same intersection
						i1++
						i2++
					}
				}
			}
		}
	}
}

// findParentPolygon determines the most suitable parent polygon for a given polygon to nest inside.
// This is part of the hierarchical structure of polygons, where polygons may contain holes or islands.
//
// Parameters:
//   - polyToNest: The polygon that needs to be nested.
//
// Returns:
//   - A pointer to the parent polygon (*PolyTree[T]) in which polyToNest should be nested.
//   - nil if polyToNest does not belong inside the current polygon.
//
// Behavior:
//   - Checks if polyToNest is completely contained within the current polygon.
//   - Recursively evaluates children polygons to find the most specific parent polygon.
//   - If no child contains polyToNest, the current polygon (p) is considered the parent.
//
// Assumptions:
//   - polyToNest and p are valid, non-nil polygons with correctly ordered contours.
//
// Complexity:
//   - The complexity depends on the number of children polygons and the depth of the polygon hierarchy.
func (p *PolyTree[T]) findParentPolygon(polyToNest *PolyTree[T]) *PolyTree[T] {
	// Step 1: Check if polyToNest is entirely inside the current polygon
	if p.contour.isContourInside(polyToNest.contour) {
		// Step 2: Recursively check children for a more specific parent
		for _, child := range p.children {
			if nestedParent := child.findParentPolygon(polyToNest); nestedParent != nil {
				return nestedParent // Found a more specific parent
			}
		}
		// Step 3: If no child contains polyToNest, the current polygon is the parent
		return p
	}
	// Step 4: If polyToNest is not inside the current polygon, return nil
	return nil
}

// findTraversalStartingPoint identifies the first unvisited entry point for polygon traversal.
// This is used during Boolean operations to start traversing a polygon's contours from a valid entry point.
//
// Parameters:
//   - other: The other polygon involved in the traversal.
//
// Returns:
//   - A pointer to the PolyTree[T] containing the starting point.
//   - The index of the starting point within the contour of the identified PolyTree[T].
//   - Returns "(nil, -1)" if no unvisited entry points are found.
//
// Behavior:
//   - Iterates through the current polygon (p) and the other polygon (other).
//   - Searches all polygons in both trees for the first unvisited entry point (marked as intersectionTypeEntry).
//   - Skips any points that are already visited.
func (p *PolyTree[T]) findTraversalStartingPoint(other *PolyTree[T]) (*PolyTree[T], int) {
	// Iterate through both polygons (`p` and `other`) to check their contours for unvisited entry points.
	for _, ptOuter := range []*PolyTree[T]{p, other} {
		// Iterate through all polygons in the current tree
		for polyTree := range ptOuter.iterPolys {
			// Iterate through all points in the current polygon's contour
			for pointIndex := range polyTree.contour {
				// Check if the point is an entry point and has not been visited
				if polyTree.contour[pointIndex].entryExit == intersectionTypeEntry && !polyTree.contour[pointIndex].visited {
					return polyTree, pointIndex // Return the polygon and the point index
				}
			}
		}
	}

	// If no unvisited entry points are found, return nil and -1
	return nil, -1
}

// iterPolys iterates over this PolyTree and all its nested polygons, including siblings and children at all levels.
// It calls the yield function for each polygon, and the iteration continues as long as yield returns true.
//
// Parameters:
//   - yield: A callback function that receives a pointer to the current PolyTree being iterated over.
//     The iteration stops early if yield returns false.
//
// Behavior:
//   - Starts with the current polygon (p) and yields it to the yield function.
//   - Recursively iterates over all siblings and their nested polygons.
//   - Recursively iterates over all children and their nested polygons.
//
// Example:
//
//	p := ... // Root PolyTree
//	for poly := range p.iterPolys {
//	    fmt.Printf("Polygon ID: %v, Type: %v\n", poly.ID, poly.polygonType)
//	}
func (p *PolyTree[T]) iterPolys(yield func(*PolyTree[T]) bool) {
	// Yield the current polygon (p)
	if !yield(p) {
		return // Stop if yield returns false
	}

	// Yield all siblings and their nested polygons
	for _, sibling := range p.siblings {
		for s := range sibling.iterPolys {
			if !yield(s) {
				return // Stop if yield returns false
			}
		}
	}

	// Yield all children and their nested polygons
	for _, child := range p.children {
		for c := range child.iterPolys {
			if !yield(c) {
				return // Stop if yield returns false
			}
		}
	}
}

// markEntryExitPoints assigns entry and exit metadata to intersection points for Boolean operations
// between two PolyTree objects.
//
// Preconditions:
//   - The findIntersections method must be called on the same PolyTree and with the same other parameter
//     before this function to ensure intersection points are identified.
//
// Parameters:
//   - other: The other PolyTree involved in the Boolean operation.
//   - operation: The Boolean operation being performed (e.g., union, intersection, subtraction).
//
// Behavior:
//
// Iterates over all contours of p and other.
//
// For each intersection point:
//   - Determines if the point is an entry or exit for traversal using a midpoint test.
//   - Uses a lookup table (entryExitPointLookUpTable) to assign entryExit types to both p and other.
//
// Links intersection points in p and other to ensure traversal can switch between polygons.
func (p *PolyTree[T]) markEntryExitPoints(other *PolyTree[T], operation BooleanOperation) {
	// Iterate through all combinations of polygons in `p` and `other`
	poly1i := 0
	for poly1 := range p.iterPolys {
		poly2i := 0
		for poly2 := range other.iterPolys {

			// Iterate through each edge in `poly1`
			for poly1Point1Index, poly1Point1 := range poly1.contour {
				poly1Point2Index := (poly1Point1Index + 1) % len(poly1.contour)

				// Process only intersection points
				if poly1Point1.pointType == pointTypeAddedIntersection {
					for poly2PointIndex, poly2Point := range poly2.contour {

						// Match intersection points in `poly1` and `poly2`
						if poly2Point.pointType == pointTypeAddedIntersection && poly1Point1.point.Eq(poly2Point.point) {

							// Sanity check: Ensure intersection metadata is not already set
							if poly1.contour[poly1Point1Index].entryExit != intersectionTypeNotSet ||
								poly2.contour[poly2PointIndex].entryExit != intersectionTypeNotSet {
								panic(fmt.Errorf("found intersection metadata when none was expected"))
							}

							// Determine if `poly1` is entering or exiting `poly2` using a midpoint test
							mid := NewLineSegment(
								poly1Point1.point,
								poly1.contour[poly1Point2Index].point).Midpoint()
							midT := NewPoint[T](T(mid.x), T(mid.y))
							poly1EnteringPoly2 := poly2.contour.isPointInside(midT)

							//// Debug information about the polygons and midpoint test
							//if poly1i == 0 {
							//	fmt.Println("poly1 outer contour")
							//} else {
							//	fmt.Println("poly1 hole contour")
							//}
							//if poly2i == 0 {
							//	fmt.Println("poly2 outer contour")
							//} else {
							//	fmt.Println("poly2 hole contour")
							//}
							//fmt.Printf("Poly1 [%d]: %v -> Poly2 [%d]: %v\n",
							//	poly1Point1Index, poly1Point1.point, poly2PointIndex, poly2Point.point)
							//fmt.Printf("Midpoint: %v, Inside Poly2: %t\n", midT, poly1EnteringPoly2)

							// Use the lookup table to determine and set entry/exit points
							poly1.contour[poly1Point1Index].entryExit =
								entryExitPointLookUpTable[operation][poly1.polygonType][poly2.polygonType][poly1EnteringPoly2].poly1PointType
							poly2.contour[poly2PointIndex].entryExit =
								entryExitPointLookUpTable[operation][poly1.polygonType][poly2.polygonType][poly1EnteringPoly2].poly2PointType

							//// Debug information about the marked entry/exit points
							//fmt.Printf("Poly1 EntryExit: %s, Poly2 EntryExit: %s\n",
							//	poly1.contour[poly1Point1Index].entryExit.String(),
							//	poly2.contour[poly2PointIndex].entryExit.String())

							// Link the intersection points for traversal
							poly1.contour[poly1Point1Index].intersectionPartner = poly2
							poly1.contour[poly1Point1Index].intersectionPartnerPointIndex = poly2PointIndex

							poly2.contour[poly2PointIndex].intersectionPartner = poly1
							poly2.contour[poly2PointIndex].intersectionPartnerPointIndex = poly1Point1Index
						}
					}
				}
			}
			poly2i++
		}
		poly1i++
	}
}

// orderSiblingsAndChildren ensures that the siblings and children of the PolyTree
// are ordered consistently. The sorting is based on the lowest, leftmost point
// of their contours.
//
// This function is particularly useful for maintaining consistent ordering of siblings
// and children across operations such as Boolean operations, where the order of polygons
// may affect the output or comparisons.
//
// Behavior:
//   - Siblings are sorted by the lowest, leftmost point in their contours.
//   - Children are similarly sorted by the lowest, leftmost point.
//
// Sorting Criteria:
//   - The compareLowestLeftmost function is used to compare contours based on their
//     lowest, leftmost points.
func (p *PolyTree[T]) orderSiblingsAndChildren() {
	// Sort siblings by the lowest, leftmost point in their contours
	slices.SortFunc(p.siblings, func(a, b *PolyTree[T]) int {
		// Compare contours using their lowest, leftmost points
		return compareLowestLeftmost(a.contour, b.contour)
	})

	// Sort children by the lowest, leftmost point in their contours
	slices.SortFunc(p.children, func(a, b *PolyTree[T]) int {
		// Compare contours using their lowest, leftmost points
		return compareLowestLeftmost(a.contour, b.contour)
	})
}

// resetIntersectionMetadataAndReorder removes intersection-related metadata from the PolyTree
// and reorders the contour points to ensure a consistent starting point.
//
// This function performs the following actions:
//  1. Iterates through all polygons in the PolyTree, including nested polygons (children).
//  2. Removes any points that were added as intersection points during Boolean operations.
//  3. Resets metadata fields such as pointType, entryExit, visited, and intersection partner information.
//  4. Reorders the contour to start at the lowest, leftmost point.
//
// This function is typically used as a preparation step before performing a new Boolean operation or traversal.
func (p *PolyTree[T]) resetIntersectionMetadataAndReorder() {
	// Iterate over all polygons in the PolyTree, including nested ones
	for poly := range p.iterPolys {

		// Iterate through the contour points of the polygon
		for i := 0; i < len(poly.contour); i++ {

			// Remove points that were added as intersections
			if poly.contour[i].pointType == pointTypeAddedIntersection {
				poly.contour = slices.Delete(poly.contour, i, i+1) // Remove the intersection point
				i--                                                // Adjust index to account for the removed point
				continue
			}

			// Reset point metadata to its default state
			poly.contour[i].pointType = pointTypeOriginal      // Mark as an original point
			poly.contour[i].entryExit = intersectionTypeNotSet // Reset entry/exit status
			poly.contour[i].visited = false                    // Mark as unvisited
			poly.contour[i].intersectionPartner = nil          // Remove intersection partner reference
			poly.contour[i].intersectionPartnerPointIndex = -1 // Reset partner point index
		}
	}

	// Reorder the contour to ensure a consistent starting point
	p.contour.reorder()
}

// ContainsPoint determines whether a given point lies inside the convex polygon.
//
// This method uses a clockwise orientation check for each edge of the polygon to determine
// if the point is on the "correct" side of all edges. For convex polygons, this is sufficient
// to verify containment.
//
// Parameters:
//
//   - point Point[T]: The point to check.
//
// Returns:
//
//   - bool: True if the point lies inside or on the boundary of the convex polygon; false otherwise.
//
// Algorithm:
//   - Iterate over each edge of the convex polygon, defined by consecutive points in the Points slice.
//   - For each edge, check the orientation of the given point relative to the edge using the Orientation function.
//   - If the point is found to be on the "outside" of any edge (i.e., the orientation is clockwise),
//     it is determined to be outside the polygon, and the method returns false.
//   - If the point passes all edge checks, it is inside the polygon, and the method returns true.
//
// Notes:
//   - This method assumes that the `simpleConvexPolygon` is indeed convex. No validation of convexity
//     is performed, as this type is intended for internal use and relies on being constructed correctly.
//
// Example:
//
//	scp := simpleConvexPolygon{Points: []Point{
//	    {X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 0, Y: 4},
//	}}
//	inside := scp.ContainsPoint(Point{X: 2, Y: 2}) // Returns true
//	outside := scp.ContainsPoint(Point{X: 5, Y: 5}) // Returns false
//
// todo: this example only makes sense in the context of this module - users of the module won't be able to do this as simpleConvexPolygon is private. Consider how we address this. Potentially add a ConvexHull method of Polygon?
func (scp *simpleConvexPolygon[T]) ContainsPoint(point Point[T]) bool {
	// Loop over each edge, defined by consecutive points
	for i := 0; i < len(scp.Points); i++ {
		a := scp.Points[i]
		b := scp.Points[(i+1)%len(scp.Points)] // Wrap to form a closed polygon

		// Check if the point is on the correct side of the edge
		if Orientation(a, b, point) == PointsClockwise {
			return false // Point is outside
		}
	}
	return true // Point is inside
}

// compareLowestLeftmost compares two contours and determines their relative order
// based on the lowest, leftmost point in each contour.
//
// It is intended for use with slices.SortFunc.
//
// The comparison is performed as follows:
//   - The contour with the point that has the smallest y coordinate is considered "smaller".
//   - If the y coordinates are equal, the contour with the smallest x coordinate is considered "smaller".
//   - Returns -1 if a should come before b, 1 if b should come before a, or 0 if both contours have the same lowest, leftmost point.
//
// This function is typically used for sorting contours in a consistent order.
//
// Parameters:
//   - a: The first contour to compare.
//   - b: The second contour to compare.
//
// Returns:
//   - An integer indicating the relative order of the contours.
func compareLowestLeftmost[T SignedNumber](a, b contour[T]) int {
	// Find the lowest, leftmost point in contour `a`
	aMin := a.findLowestLeftmost()

	// Find the lowest, leftmost point in contour `b`
	bMin := b.findLowestLeftmost()

	// Compare based on the `y` coordinate first, then the `x` coordinate
	if aMin.y < bMin.y || (aMin.y == bMin.y && aMin.x < bMin.x) {
		return -1 // Contour `a` comes before contour `b`
	}
	if aMin.y > bMin.y || (aMin.y == bMin.y && aMin.x > bMin.x) {
		return 1 // Contour `b` comes before contour `a`
	}

	return 0 // Both contours are equal in order
}

// nestPointsToPolyTrees organizes a set of contours into a hierarchical PolyTree structure.
//
// This function processes a list of polygon contours and builds a tree where:
//   - The largest contour becomes the root of the tree (assumed to be solid).
//   - Smaller contours are nested hierarchically based on spatial relationships.
//   - Polygons within solid polygons are classified as holes, and polygons within holes
//     are classified as solid islands.
//
// Parameters:
//   - contours: A slice of slices of Points, where each inner slice represents a closed polygon contour.
//
// Returns:
//   - A pointer to the root PolyTree if the nesting is successful.
//   - An error if no contours are provided or if any operation fails.
//
// Behavior:
//   - The function sorts the contours by area in ascending order to ensure the largest polygon is processed last.
//   - It determines nesting relationships between polygons to establish a parent-child hierarchy.
//
// Example:
//
//	contours := [][]Point[int]{
//	    {{0, 0}, {10, 0}, {10, 10}, {0, 10}}, // Outer square
//	    {{2, 2}, {4, 2}, {4, 4}, {2, 4}},     // Hole
//	}
//	rootTree, err := nestPointsToPolyTrees(contours)
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	}
func nestPointsToPolyTrees[T SignedNumber](contours [][]Point[T]) (*PolyTree[T], error) {
	// Sanity check: ensure contours exist
	if len(contours) == 0 {
		return nil, fmt.Errorf("no contours provided")
	}

	// Step 1: Sort polygons by area (ascending order, largest last)
	slices.SortFunc(contours, func(a, b []Point[T]) int {
		areaA := SignedArea2X(a)
		areaB := SignedArea2X(b)
		switch {
		case areaA < areaB:
			return -1
		case areaA > areaB:
			return 1
		default:
			return 0
		}
	})

	// Step 2: Create the root PolyTree from the largest polygon
	rootTree, err := NewPolyTree(contours[len(contours)-1], PTSolid)
	if err != nil {
		return nil, fmt.Errorf("failed to create root PolyTree: %w", err)
	}

	// Step 3: Process the remaining polygons in reverse order (smallest first)
	for i := len(contours) - 2; i >= 0; i-- {
		// Create a PolyTree for the current contour
		polyToNest, err := NewPolyTree(contours[i], PTSolid)
		if err != nil {
			return nil, fmt.Errorf("failed to create PolyTree for contour %d: %w", i, err)
		}

		// Step 4: Find the correct parent polygon
		parent := rootTree.findParentPolygon(polyToNest)
		if parent == nil {
			// No parent found: Add as a sibling to the root
			if err := rootTree.addSibling(polyToNest); err != nil {
				return nil, fmt.Errorf("failed to add sibling: %w", err)
			}
		} else {
			// Parent found: Add as a child of the parent
			// Adjust polygon type based on the parent's type
			if parent.polygonType == PTSolid {
				polyToNest.polygonType = PTHole
			} else {
				polyToNest.polygonType = PTSolid
			}
			if err := parent.addChild(polyToNest); err != nil {
				return nil, fmt.Errorf("failed to add child to parent polygon: %w", err)
			}
		}
	}

	// Step 5: Return the root PolyTree
	return rootTree, nil
}

// togglePolyTraversalDirection toggles the traversal direction of a polygon between clockwise
// and counterclockwise.
//
// This function is used in algorithms that require switching the traversal direction of a polygon's
// vertices, such as during Boolean operations or polygon manipulations.
//
// Parameters:
//
//   - direction polyTraversalDirection: The current traversal direction, either clockwise
//     or counterclockwise.
//
// Returns polyTraversalDirection: The opposite traversal direction.
//   - If the input is polyTraversalReverse, the output
//     will be polyTraversalForward.
//   - If the input is polyTraversalForward, the output
//     will be polyTraversalReverse.
//
// Example:
//
//	currentDirection := polyTraversalReverse
//	newDirection := togglePolyTraversalDirection(currentDirection)
//	// newDirection == polyTraversalForward
//
// Notes:
//   - This function assumes that the input is a valid polyTraversalDirection value.
//   - If an invalid value is provided, the function may not behave as expected.
//
// Dependencies:
//   - Relies on the polyTraversalDirection type and its constants for traversal direction.
func togglePolyTraversalDirection(direction polyTraversalDirection) polyTraversalDirection {
	if direction == polyTraversalReverse {
		return polyTraversalForward
	}
	return polyTraversalReverse
}

// withVisited creates an option to configure a PolyTree equality check with a predefined
// set of visited nodes.
//
// This is useful for avoiding infinite recursion when comparing deeply nested or circular
// PolyTree structures. By passing a shared `visited` map, multiple equality checks can
// track which nodes have already been processed.
//
// Parameters:
//   - visited: A map where keys are PolyTree pointers that have already been visited.
//
// Returns:
//   - A polyTreeEqOption function that updates the equality check configuration.
//
// Example:
//
//	visitedNodes := make(map[*PolyTree[int]]bool)
//	eqOption := withVisited(visitedNodes)
//	match, mismatches := polyTree1.Eq(polyTree2, eqOption)
func withVisited[T SignedNumber](visited map[*PolyTree[T]]bool) polyTreeEqOption[T] {
	return func(cfg *polyTreeEqConfig[T]) {
		// Update the visited map in the equality configuration
		cfg.visited = visited
	}
}
