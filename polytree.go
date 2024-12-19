package geom2d

import (
	"fmt"
	"slices"
	"strings"
)

// BooleanOperation defines the types of Boolean operations that can be performed on polygons.
// These operations are fundamental in computational geometry for combining or modifying shapes.
type BooleanOperation uint8

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

// NewPolyTreeOption defines a functional option type for configuring a new [PolyTree] during creation.
//
// This type allows for flexible and extensible initialization of [PolyTree] objects by applying optional
// configurations after the core properties have been set.
//
// Parameters:
//   - T: The numeric type of the coordinates in the [PolyTree], constrained by the [SignedNumber] interface.
//
// This pattern makes it easy to add optional properties to a PolyTree without requiring an extensive list
// of parameters in the NewPolyTree function.
type NewPolyTreeOption[T SignedNumber] func(*PolyTree[T]) error

// WithChildren is an option for the [NewPolyTree] function that assigns child polygons to the created [PolyTree].
// It also sets up parent-child relationships and orders the children for consistency.
//
// Parameters:
//   - children: A variadic list of pointers to [PolyTree] objects representing the child polygons.
//
// Behavior:
//   - The function assigns the provided children to the [PolyTree] being created.
//   - It establishes the parent-child relationship by setting the parent of each child to the newly created [PolyTree].
//
// Returns:
//   - A [NewPolyTreeOption] that can be passed to the [NewPolyTree] function.
func WithChildren[T SignedNumber](children ...*PolyTree[T]) NewPolyTreeOption[T] {
	return func(p *PolyTree[T]) error {
		for _, child := range children {
			err := p.AddChild(child)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// WithSiblings is an option for the [NewPolyTree] function that assigns sibling polygons to the created [PolyTree].
//
// Parameters:
//   - siblings: A variadic list of pointers to [PolyTree] objects representing the sibling polygons.
//
// Behavior:
//   - The function assigns the provided children to the [PolyTree] being created.
//
// Returns:
//   - A [NewPolyTreeOption] that can be passed to the [NewPolyTree] function.
func WithSiblings[T SignedNumber](siblings ...*PolyTree[T]) NewPolyTreeOption[T] {
	return func(p *PolyTree[T]) error {
		for _, sibling := range siblings {
			err := p.AddSibling(sibling)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Valid values for PolygonType
const (
	// PTSolid represents a solid region of the polygon, commonly referred to as an "island."
	// PTSolid polygons are the primary filled areas, excluding any void regions (holes).
	PTSolid PolygonType = iota

	// PTHole represents a void region of the polygon, often nested within a solid polygon.
	// Holes are not part of the filled area of the polygon and are treated as exclusions.
	PTHole
)

// String returns a string representation of the [PolygonType].
//
// This method converts the enum value of a [PolygonType] into a human-readable string,
// making it useful for debugging, logging, or providing textual representations.
//
// Returns:
//   - string: A string representation of the [PolygonType], such as [PTSolid] or [PTHole].
//
// Panics:
//   - If the [PolygonType] value is unsupported, the method will panic with an appropriate error message.
//
// Notes:
//   - This method assumes that the [PolygonType] is valid. If new polygon types are added
//     in the future, ensure this method is updated to include them to avoid panics.
func (t PolygonType) String() string {
	switch t {
	case PTSolid:
		return "PTSolid"
	case PTHole:
		return "PTHole"
	default:
		panic(fmt.Errorf("unsupported PolygonType"))
	}
}

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

// PolyTreeMismatch represents a bitmask of potential mismatches between two PolyTree structures.
type PolyTreeMismatch uint8

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
//   - A [PolygonType]: Indicates whether the polygon is a solid region ([PTSolid]) or a hole ([PTHole]).
//     This classification is essential for understanding the relationship between the polygon
//     and its children.
//   - A parent: a pointer to the parent polygon in the hierarchy. For example, a hole's parent
//     would be the solid polygon that contains it. If a polygon is the root polygon in the PolyTree, its parent is nil.
//   - Zero or more siblings: A list of sibling polygons that are not nested within each other but share
//     the same parent. Siblings must be of the same [PolygonType].
//   - Zero or more children: A list of child polygons nested within this polygon. If the [PolygonType] is
//     [PTSolid], the children are holes ([PTHole]). If the [PolygonType] is [PTHole], the children
//     are solid islands ([PTSolid]).
//
// Hierarchy Rules:
//   - A solid polygon ([PTSolid]) can contain holes ([PTHole]) as its children.
//   - A hole ([PTHole]) can contain solid polygons ([PTSolid]) as its children.
//   - Siblings are polygons of the same [PolygonType] that do not overlap.
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
	// in point-in-polygon checks and other spatial queries. This value is doubled as-per contour.
	maxX T
}

// contour represents the outline of a polygon as a slice of polyTreePoint entries.
// Each entry contains metadata about the point, such as whether it is a normal vertex,
// an intersection point, or a midpoint between intersections.
//
// The contour is used to define the polygon's shape and is processed during boolean
// operations. Contour within the contour are typically doubled to facilitate calculations
// involving midpoints and to avoid precision issues when working with integer-based
// coordinates.
type contour[T SignedNumber] []polyTreePoint[T]

// String returns a string representation of the contour, listing all its points in order.
//
// This method iterates through all polyTreePoints in the contour and appends their coordinates
// to a human-readable string.
//
// Returns:
//   - string: A formatted string showing the list of points in the contour.
func (c *contour[T]) String() string {
	var builder strings.Builder
	builder.WriteString("Contour Points: [")

	first := true
	for _, pt := range c.toPoints() {
		if first {
			builder.WriteString(fmt.Sprintf("(%v, %v)", pt.x, pt.y))
			first = false
			continue
		}
		builder.WriteString(fmt.Sprintf(", (%v, %v)", pt.x, pt.y))
	}

	builder.WriteString("]")
	return builder.String()
}

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
	rel detailedLineSegmentRelationship
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
//   - points ([][Point][T]): A slice of [Point][T] representing the vertices of the polygon.
//   - t ([PolygonType]): The type of polygon, either [PTSolid] or [PTHole].
//   - opts: A variadic slice of [NewPolyTreeOption] functions, see below.
//
// [NewPolyTreeOption] functions:
//   - [WithChildren]: is an option for the [NewPolyTree] function that assigns child polygons to the
//     created [PolyTree]. This can also be done later with the [PolyTree.AddChild] method.
//   - [WithSiblings]: is an option for the [NewPolyTree] function that assigns sibling polygons to the
//     created [PolyTree]. This can also be done later with the [PolyTree.AddSibling] method.
//
// Returns:
//   - A pointer to the newly created PolyTree.
//   - An error if the input points are invalid (e.g., less than three points or zero area).
//
// Notes:
//   - The function ensures that the polygon's points are oriented correctly based on its type.
//   - Contour are doubled internally to avoid integer division/precision issues during midpoint calculations.
//   - The polygon's convex hull is computed and stored for potential optimisations.
//   - Child polygons must have the opposite [PolygonType] (e.g., holes for a solid polygon and solids for a hole polygon).
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
		err := opt(p)
		if err != nil {
			return nil, err
		}
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

// ensureClockwise ensures that the contour points are ordered in a clockwise direction.
//
// This function calculates the signed area of the contour. If the area is positive, it indicates
// that the points are in a counterclockwise order. In such cases, the function reverses the order
// of the points to make the contour clockwise. If the area is already negative (indicating clockwise
// order), no action is taken.
//
// This is particularly useful for consistency when dealing with hole polygons or when ensuring
// correct orientation for geometric operations like Boolean polygon operations.
//
// Notes:
//   - A positive signed area indicates counterclockwise orientation.
//   - A negative signed area indicates clockwise orientation.
func (c *contour[T]) ensureClockwise() {
	area := SignedArea2X(c.toPoints())
	if area < 0 {
		return // Already clockwise
	}
	slices.Reverse(*c)
}

// ensureCounterClockwise ensures that the contour points are ordered in a counterclockwise direction.
//
// This function calculates the signed area of the contour. If the area is negative, it indicates
// that the points are in a clockwise order. In such cases, the function reverses the order of
// the points to make the contour counterclockwise. If the area is already positive (indicating
// counterclockwise order), no action is taken.
//
// This is particularly useful for consistency when dealing with solid polygons or for geometric
// operations that require specific point orientations.
//
// Notes:
//   - A positive signed area indicates counterclockwise orientation.
//   - A negative signed area indicates clockwise orientation.
func (c *contour[T]) ensureCounterClockwise() {
	area := SignedArea2X(c.toPoints())
	if area > 0 {
		return // Already counterclockwise
	}
	slices.Reverse(*c)
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
	// Iterate over each point in "other".
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
//   - Contour lying directly on the contour edges are considered inside.
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
		if edges[i].lineSegment.ContainsPoint(point) {
			return true
		}

		// Determine the relationship of the edge to the ray.
		edges[i].rel = ray.detailedRelationshipToLineSegment(edges[i].lineSegment)
	}

	// Analyze the relationships and count ray crossings.
	for i := range edges {
		// Look at the next edge in the contour.
		iNext := (i + 1) % len(edges)

		switch edges[i].rel {
		case lsrIntersects: // Ray intersects the edge.
			crosses++

		case lsrCollinearCDinAB: // Ray is collinear with the edge and overlaps it.
			crosses += 1

		case lsrConAB: // Ray starts on the edge.
			crosses++

			// Handle potential overlaps with the next edge.
			if edges[iNext].rel == lsrDonAB {
				if inOrder(edges[i].lineSegment.start.y, point.y, edges[iNext].lineSegment.end.y) {
					crosses++
				}
			}

		case lsrDonAB: // Ray ends on the edge.
			crosses++

			// Handle potential overlaps with the next edge.
			if edges[iNext].rel == lsrConAB {
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

// String converts a polyIntersectionType value to its corresponding string representation.
//
// This method provides a human-readable string for the polyIntersectionType value, such as
// "not set", "entry", or "exit". If the value is unsupported, the method panics.
//
// Returns:
//   - string: A string representation of the polyIntersectionType.
//
// Panics:
//   - If the polyIntersectionType value is not one of the defined constants.
func (p *polyIntersectionType) String() string {
	switch *p {
	case intersectionTypeNotSet:
		return "not set"
	case intersectionTypeEntry:
		return "entry"
	case intersectionTypeExit:
		return "exit"
	default:
		panic("unsupported polyIntersectionType")
	}
}

// AddChild adds a child PolyTree to the current PolyTree.
//
// Parameters:
//   - child (*PolyTree[T]): A pointer to the PolyTree to be added as a child.
//
// Returns:
//
// error: An error if the operation fails. Possible error scenarios include:
//   - The child is nil.
//   - The child has the same [PolygonType] as the parent.
//   - The child polygon does not fit entirely within the parent's contour.
//   - The child polygon overlaps with existing children.
func (pt *PolyTree[T]) AddChild(child *PolyTree[T]) error {
	// Check if the child is nil
	if child == nil {
		return fmt.Errorf("attempt to add nil child")
	}

	// Ensure the polygon types are compatible
	if pt.polygonType == child.polygonType {
		return fmt.Errorf(
			"cannot add child: mismatched polygon types (parent: %v, child: %v)",
			pt.polygonType,
			child.polygonType,
		)
	}

	// Check if the child fits within the parent's contour
	if !pt.contour.isContourInside(child.contour) {
		return fmt.Errorf("child polygon does not fit entirely within the parent polygon")
	}

	// Check if the child overlaps with existing children
	for _, sibling := range pt.children {
		for siblingEdge := range sibling.contour.iterEdges {
			for childEdge := range child.contour.iterEdges {
				if siblingEdge.detailedRelationshipToLineSegment(childEdge) > lsrMiss {
					return fmt.Errorf("child polygon overlaps with an existing sibling polygon")
				}
			}
		}
	}

	// Set point directionality
	switch child.polygonType {
	case PTSolid:
		child.contour.ensureCounterClockwise()
	case PTHole:
		child.contour.ensureClockwise()
	}

	// Set the parent of the child
	child.parent = pt

	// Append the child to the children slice
	pt.children = append(pt.children, child)

	// Order siblings and children for consistency
	pt.orderSiblingsAndChildren()

	return nil
}

// AddSibling adds a sibling PolyTree to the current PolyTree.
//
// Parameters:
//   - sibling (*PolyTree[T]): A pointer to the PolyTree to be added as a sibling.
//
// Returns:
//
// error: An error if the operation fails. Possible error scenarios include:
//   - The sibling is nil.
//   - The sibling has a different [PolygonType] than the current PolyTree.
//   - The sibling polygon overlaps with existing siblings.
func (pt *PolyTree[T]) AddSibling(sibling *PolyTree[T]) error {
	// Check if the sibling is nil
	if sibling == nil {
		return fmt.Errorf("attempt to add nil sibling")
	}

	// Ensure the polygon types match
	if pt.polygonType != sibling.polygonType {
		return fmt.Errorf("cannot add sibling: mismatched polygon types")
	}

	// Check if the sibling overlaps with existing siblings
	for _, existingSibling := range pt.siblings {
		for existingSiblingEdge := range existingSibling.contour.iterEdges {
			for siblingEdge := range sibling.contour.iterEdges {
				if siblingEdge.detailedRelationshipToLineSegment(existingSiblingEdge) > lsrMiss {
					return fmt.Errorf("child polygon overlaps with an existing sibling polygon")
				}
			}
		}
	}

	// Add the new sibling to the sibling lists of existing siblings
	for _, existingSibling := range pt.siblings {
		// Update sibling relationships
		existingSibling.siblings = append(existingSibling.siblings, sibling)

		// Maintain consistent ordering
		existingSibling.orderSiblingsAndChildren()

		// Add the existing sibling to the new sibling's sibling list
		sibling.siblings = append(sibling.siblings, existingSibling)
	}

	// Add the current `PolyTree` to the new sibling's sibling list
	sibling.siblings = append(sibling.siblings, pt)

	// Maintain consistent ordering
	sibling.orderSiblingsAndChildren()

	// Add the new sibling to the current `PolyTree`'s sibling list
	pt.siblings = append(pt.siblings, sibling)

	// Maintain consistent ordering
	pt.orderSiblingsAndChildren()

	return nil
}

// Area calculates the area of the polygon represented by the PolyTree.
// This method computes the absolute area of the PolyTree's contour, ensuring
// that the result is always positive, regardless of the orientation of the contour.
//
// Returns:
//   - float64: The absolute area of the PolyTree.
//
// Notes:
//   - The method operates only on the contour of the PolyTree and does not account for any children polygons.
func (pt *PolyTree[T]) Area() float64 {
	area := float64(SignedArea2X(pt.Contour()))
	if area < 0 {
		area *= -1
	}
	area /= 2
	return area
}

// AsFloat32 converts a [PolyTree] with generic numeric type T to a new [PolyTree] with points of type float32.
//
// This method iterates over the contours and nodes of the current [PolyTree], converting all points
// to float32 using the [Point.AsFloat32] method. It then rebuilds the [PolyTree] with the converted points.
//
// Returns:
//   - *PolyTree[float32]: A new [PolyTree] instance where all points are of type float32.
//
// Panics:
//   - If an error occurs during the nesting of contours into the new [PolyTree], the function panics.
func (pt *PolyTree[T]) AsFloat32() *PolyTree[float32] {
	output, err := ApplyPointTransform[T](pt, func(p Point[T]) (Point[float32], error) {
		return p.AsFloat32(), nil
	})

	if err != nil {
		panic(err)
	}

	return output
}

// AsFloat64 converts a [PolyTree] with generic numeric type T to a new [PolyTree] with points of type float64.
//
// This method iterates over the contours and nodes of the current [PolyTree], converting all points
// to float64 using the [Point.AsFloat64] method. It then rebuilds the [PolyTree] with the converted points.
//
// Returns:
//   - *PolyTree[float64]: A new [PolyTree] instance where all points are of type float64.
//
// Panics:
//   - If an error occurs during the nesting of contours into the new [PolyTree], the function panics.
func (pt *PolyTree[T]) AsFloat64() *PolyTree[float64] {
	output, err := ApplyPointTransform[T](pt, func(p Point[T]) (Point[float64], error) {
		return p.AsFloat64(), nil
	})

	if err != nil {
		panic(err)
	}

	return output
}

// AsInt converts a [PolyTree] with generic numeric type T to a new [PolyTree] with points of type int.
//
// This method iterates over the contours and nodes of the current [PolyTree], converting all points
// to int using the [Point.AsInt] method. It then rebuilds the [PolyTree] with the converted points.
//
// Returns:
//   - *PolyTree[int]: A new [PolyTree] instance where all points are of type int.
//
// Panics:
//   - If an error occurs during the nesting of contours into the new [PolyTree], the function panics.
func (pt *PolyTree[T]) AsInt() *PolyTree[int] {

	output, err := ApplyPointTransform[T](pt, func(p Point[T]) (Point[int], error) {
		return p.AsInt(), nil
	})

	if err != nil {
		panic(err)
	}

	return output
}

// AsIntRounded converts a [PolyTree] with generic numeric type T to a new [PolyTree] with points of type int,
// rounding the coordinates of each point to the nearest integer.
//
// This method iterates over the contours and nodes of the current [PolyTree], converting all points to int
// using the [Point.AsIntRounded] method. It then rebuilds the [PolyTree] with the rounded points.
//
// Returns:
//   - *PolyTree[int]: A new [PolyTree] instance where all points are of type int, with coordinates rounded.
//
// Panics:
//   - If an error occurs during the nesting of contours into the new [PolyTree], the function panics.
func (pt *PolyTree[T]) AsIntRounded() *PolyTree[int] {
	output, err := ApplyPointTransform[T](pt, func(p Point[T]) (Point[int], error) {
		return p.AsIntRounded(), nil
	})

	if err != nil {
		panic(err)
	}

	return output
}

// BooleanOperation performs a Boolean operation (union, intersection, or subtraction)
// between the current polygon (p) and another polygon (other). The result is
// returned as a new PolyTree, or an error is returned if the operation fails.
//
// Parameters:
//   - other (*PolyTree[T]): The polygon to perform the operation with.
//   - operation ([BooleanOperation]): The type of Boolean operation to perform (e.g., union, intersection, subtraction).
//
// Returns:
//   - A new PolyTree resulting from the operation, or an error if the operation fails.
//
// Supported Operations:
//   - [BooleanIntersection] returns the result of the intersection operation, combining two polytrees into a single
//     polytree that encompasses the areas of both input polytrees. Overlapping regions are merged.
//   - [BooleanSubtraction] returns the result of the subtraction operation, which removes the area of polytree "other"
//     from the calling PolyTree. The result is a polytree containing one or more polygons representing the area of the
//     calling polygon excluding the overlapping region with the "other" polygon.
//   - [BooleanUnion] returns the result of the union operation, which computes the region(s)
//     where two polygons overlap. The result is a polytree with overlapping regions merged.
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
func (pt *PolyTree[T]) BooleanOperation(other *PolyTree[T], operation BooleanOperation) (*PolyTree[T], error) {
	// Edge Case: Check if the polygons intersect
	if !pt.Intersects(other) {
		switch operation {
		case BooleanUnion:
			// Non-intersecting polygons: Add other as a sibling
			if err := pt.AddSibling(other); err != nil {
				return nil, fmt.Errorf("failed to add sibling: %w", err)
			}
			return pt, nil

		case BooleanIntersection:
			// Non-intersecting polygons: No intersection, return nil
			return nil, nil

		case BooleanSubtraction:
			// Non-intersecting polygons: No change to pt
			return pt, nil

		default:
			// Invalid or unsupported operation
			return nil, fmt.Errorf("unknown operation: %v", operation)
		}
	}

	// Step 1: Find intersection points between all polygons
	pt.findIntersections(other)

	// Step 2: Mark entry/exit points for traversal based on the operation
	pt.markEntryExitPoints(other, operation)

	// Step 3: Perform traversal to construct the result of the Boolean operation
	contours := pt.booleanOperationTraversal(other, operation)

	return nestPointsToPolyTrees(contours)
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
//   - Contour are halved during traversal to revert the earlier doubling for precision handling.
func (pt *PolyTree[T]) booleanOperationTraversal(other *PolyTree[T], operation BooleanOperation) [][]Point[T] {
	var direction polyTraversalDirection

	// Initialize the resulting contours
	resultContours := make([][]Point[T], 0)

	for {
		// Find the starting point for traversal
		currentPoly, currentPointIndex := pt.findTraversalStartingPoint(other)
		if currentPoly == nil || currentPointIndex == -1 {
			// No unvisited entry points, traversal is complete
			break
		}

		// Initialize a new contour for the result
		resultContour := make([]Point[T], 0, len(pt.contour)+len(other.contour))

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

			// Handle polygon switching at entry/exit points and adjust traversal direction if needed.
			// The if condition is true when:
			//   1. The current point is an exit point (intersectionTypeExit),
			//	  OR
			//   2. The operation is BooleanSubtraction AND:
			//        - The polygon is a hole and the current point is an entry point,
			//	       OR
			//        - The polygon is solid, the current point is an entry point, and the traversal direction is reverse.
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
				pt.contour.toPoints(),
				other.contour.toPoints(),
			}
		case BooleanIntersection:
			return nil // No intersection
		case BooleanSubtraction:
			return [][]Point[T]{
				pt.contour.toPoints(),
			}
		default:
			panic(fmt.Errorf("unknown BooleanOperation: %v", operation))
		}
	}

	return resultContours
}

// BoundingBox calculates the axis-aligned bounding box (AABB) of the PolyTree.
//
// The bounding box is the smallest rectangle, aligned with the coordinate axes, that completely encloses
// all polygons in the PolyTree. The calculation uses the convex hull of each polygon in the tree, ensuring
// efficiency and accuracy.
//
// Returns:
//   - [Rectangle][T]: The axis-aligned bounding box that encloses all polygons in the PolyTree.
//
// Notes:
//   - The bounding box is computed by iterating through the convex hull points of all polygons in the PolyTree.
//   - If the PolyTree is empty, the behavior is undefined and should be handled externally by the caller.
func (pt *PolyTree[T]) BoundingBox() Rectangle[T] {
	minX := pt.contour[0].point.x / 2
	maxX := pt.contour[0].point.x / 2
	minY := pt.contour[0].point.y / 2
	maxY := pt.contour[0].point.y / 2
	for poly := range pt.Nodes {
		for _, point := range poly.hull.Points {
			minX = min(minX, point.x)
			maxX = max(maxX, point.x)
			minY = min(minY, point.y)
			maxY = max(maxY, point.y)
		}
	}
	// output points are divided by 2 as points in polytree are doubled
	return NewRectangle([]Point[T]{
		NewPoint[T](minX, minY),
		NewPoint[T](maxX, minY),
		NewPoint[T](maxX, maxY),
		NewPoint[T](minX, maxY),
	})
}

// Children returns the immediate child polygons of the current PolyTree node.
//
// The children are represented as a slice of pointers to PolyTree instances.
// Each child polygon is nested directly within the current PolyTree node.
//
// Returns:
//   - []*PolyTree[T]: A slice containing the children of the current PolyTree.
//
// Notes:
//   - The method does not include grandchildren or deeper descendants in the returned slice.
//   - If the current PolyTree node has no children, an empty slice is returned.
func (pt *PolyTree[T]) Children() []*PolyTree[T] {
	return pt.children
}

// Contour returns a slice of points that make up the external contour of the current PolyTree node.
// The contour represents the boundary of the polygon associated with the node.
//
// The last [Point] in the slice is assumed to connect to the first [Point], forming a closed contour.
//
// Returns:
//   - []Point[T]: A slice of [Point] representing the vertexes of the PolyTree node's contour.
func (pt *PolyTree[T]) Contour() []Point[T] {
	return pt.contour.toPoints()
}

// Edges returns the edges of the current PolyTree node as a slice of [LineSegment].
//
// Each edge represents a segment of the current PolyTree node's contour.
//
// Returns:
//   - [][LineSegment][T]: A slice of [LineSegment] representing the edges of the PolyTree node's contour.
//
// Notes:
//   - If the PolyTree node has no contour, an empty slice is returned.
//   - This method does not include edges from child polygons or sibling polygons.
func (pt *PolyTree[T]) Edges() []LineSegment[T] {
	edges := make([]LineSegment[T], 0, len(pt.contour))
	for edge := range pt.contour.iterEdges {
		edges = append(edges, NewLineSegment(
			NewPoint(edge.start.x/2, edge.start.y/2),
			NewPoint(edge.end.x/2, edge.end.y/2),
		))
	}
	return edges
}

// Eq compares two PolyTree objects (pt and other) for structural and content equality.
// It identifies mismatches in contours, siblings, and children while avoiding infinite recursion
// by tracking visited nodes. The comparison results are represented as a boolean and a bitmask.
//
// Parameters:
//   - other (*PolyTree[T]): The PolyTree to compare with the current PolyTree.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the start and end
//     points of the line segments. If the absolute difference between the coordinates of
//     the points is less than epsilon, they are considered equal.
//
// Returns:
//   - A boolean indicating whether the two PolyTree objects are equal.
//   - A [PolyTreeMismatch] bitmask that specifies which aspects (if any) differ.
//
// Behavior:
//   - Handles nil cases for the "pt" and "other" PolyTree objects.
//   - Recursively compares contours, siblings, and children, considering different sibling/child orders.
//   - Uses the visited map to prevent infinite recursion in cyclic structures.
func (pt *PolyTree[T]) Eq(other *PolyTree[T], opts ...polyTreeEqOption[T]) (bool, PolyTreeMismatch) {
	var mismatches PolyTreeMismatch

	// Handle nil cases
	switch {
	case pt == nil && other == nil:
		return true, PTMNoMismatch // Both are nil, no mismatch
	case pt == nil:
		return false, PTMNilPolygonMismatch // Only pt is nil
	case other == nil:
		return false, PTMNilPolygonMismatch // Only other is nil
	}

	// Initialize comparison configuration
	config := &polyTreeEqConfig[T]{
		visited: make(map[*PolyTree[T]]bool),
	}
	for _, opt := range opts {
		opt(config) // Apply configuration options
	}

	// Avoid comparing the same PolyTree objects multiple times
	if config.visited[pt] || config.visited[other] {
		return true, PTMNoMismatch // Already compared, no mismatch
	}
	config.visited[pt] = true
	config.visited[other] = true

	// Compare contours
	if !pt.contour.eq(other.contour) {
		mismatches |= PTMContourMismatch
		fmt.Println("Mismatch in contours detected")
	}

	// Compare siblings
	if len(pt.siblings) != len(other.siblings) {
		mismatches |= PTMSiblingMismatch
		fmt.Println("Mismatch in siblings count detected")
	} else {
		// Siblings order doesn't matter, check if each sibling matches
		for _, sibling := range pt.siblings {
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
	if len(pt.children) != len(other.children) {
		mismatches |= PTMChildMismatch
		fmt.Println("Mismatch in children count detected")
	} else {
		// Children order doesn't matter, check if each child matches
		for _, child := range pt.children {
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
func (pt *PolyTree[T]) findIntersections(other *PolyTree[T]) {

	// Step 1: Reset intersection metadata and reorder both PolyTrees
	pt.resetIntersectionMetadataAndReorder()
	other.resetIntersectionMetadataAndReorder()

	// Step 2: Iterate through all combinations of polygons
	for poly1 := range pt.Nodes {
		for poly2 := range other.Nodes {

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
func (pt *PolyTree[T]) findTraversalStartingPoint(other *PolyTree[T]) (*PolyTree[T], int) {
	// Iterate through both polygons (`pt` and `other`) to check their contours for unvisited entry points.
	for _, ptOuter := range []*PolyTree[T]{pt, other} {
		// Iterate through all polygons in the current tree
		for polyTree := range ptOuter.Nodes {
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

// Hull returns the convex hull of the PolyTree node's contour as a slice of [Point].
//
// The convex hull represents the smallest convex polygon that can enclose the points
// of the PolyTree node's contour. This is useful for optimizations in geometric operations.
//
// Returns:
//   - [][Point][T]: A slice of points representing the convex hull of the PolyTree node.
//
// Notes:
//   - If the PolyTree node does not have a contour, an empty slice is returned.
//   - The points in the hull are ordered counterclockwise.
func (pt *PolyTree[T]) Hull() []Point[T] {
	return pt.hull.Points
}

// Intersects checks whether the current PolyTree intersects with another PolyTree.
//
// Parameters:
//   - other (*PolyTree[T]): The PolyTree to compare against the current PolyTree.
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
func (pt *PolyTree[T]) Intersects(other *PolyTree[T]) bool {
	// Check if any point of "other" is inside the contour of "pt"
	for _, otherPoint := range other.contour {
		if pt.contour.isPointInside(otherPoint.point) {
			return true // Intersection found via point containment
		}
	}

	// Check if any point of "pt" is inside the contour of "other"
	for _, point := range pt.contour {
		if other.contour.isPointInside(point.point) {
			return true // Intersection found via point containment
		}
	}

	// Check for edge intersections between "pt" and "other"
	for poly1Edge := range pt.contour.iterEdges {
		for poly2Edge := range other.contour.iterEdges {
			if poly1Edge.IntersectsLineSegment(poly2Edge) {
				return true // Intersection found via edge overlap
			}
		}
	}

	// No intersections detected
	return false
}

// IsRoot determines if the current PolyTree node is the root of the tree.
//
// A PolyTree node is considered root if it does not have a parent.
//
// Returns:
//   - bool: true if the PolyTree node is the root (has no parent),
//     false otherwise.
func (pt *PolyTree[T]) IsRoot() bool {
	if pt.parent == nil {
		return true
	}
	return false
}

// Len returns the total number of PolyTree nodes in the current PolyTree structure,
// including the root, its siblings, and all nested children.
//
// This method iterates through all nodes in the PolyTree and counts them.
//
// Returns:
//   - int: The total number of PolyTree nodes.
func (pt *PolyTree[T]) Len() int {
	i := 0
	for range pt.Nodes {
		i++
	}
	return i
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
func (pt *PolyTree[T]) markEntryExitPoints(other *PolyTree[T], operation BooleanOperation) {
	// Iterate through all combinations of polygons in `pt` and `other`
	poly1i := 0
	for poly1 := range pt.Nodes {
		poly2i := 0
		for poly2 := range other.Nodes {

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
								poly1.contour[poly1Point2Index].point).Center()
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
							//fmt.Printf("Center: %v, Inside Poly2: %t\n", midT, poly1EnteringPoly2)

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

// Nodes iterates over this PolyTree and all its nested polygons, including siblings and children at all levels equal to
// and below the current node's level. It does not traverse above the current level, if this is required, it is
// suggested to first obtain the root using [PolyTree.Root].
//
// This function is an iterator, designed for use in a for loop, for example:
//
//	for node := range polytree.Nodes {
//	    // do something
//	}
//
// Parameters:
//   - yield: A callback function that receives a pointer to the current PolyTree being iterated over.
//     The iteration stops early if yield returns false.
//
// Behavior:
//   - Starts with the current polygon (p) and yields it to the yield function.
//   - Recursively iterates over all siblings and their nested polygons.
//   - Recursively iterates over all children and their nested polygons.
func (pt *PolyTree[T]) Nodes(yield func(*PolyTree[T]) bool) {
	var currentPt *PolyTree[T]

	toYield := make([]*PolyTree[T], 0, 16)
	yielded := make([]*PolyTree[T], 0, 16)

	// add pt to toYield
	toYield = append(toYield, pt)

	// for loop: while toYield has values
	for len(toYield) > 0 {

		// pop poly from toYield
		currentPt, toYield = toYield[0], toYield[1:]

		// yield poly, return if needed
		if !yield(currentPt) {
			return
		}

		// add poly to yielded
		yielded = append(yielded, currentPt)

		// add all siblings and children of poly to toYield if not already added, and if not already yielded
		for _, sibling := range currentPt.siblings {
			if !slices.Contains(yielded, sibling) && !slices.Contains(toYield, sibling) {
				toYield = append(toYield, sibling)
			}
		}
		for _, child := range currentPt.children {
			if !slices.Contains(yielded, child) && !slices.Contains(toYield, child) {
				toYield = append(toYield, child)
			}
		}
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
func (pt *PolyTree[T]) orderSiblingsAndChildren() {
	// Sort siblings by the lowest, leftmost point in their contours
	slices.SortFunc(pt.siblings, func(a, b *PolyTree[T]) int {
		// Compare contours using their lowest, leftmost points
		return compareLowestLeftmost(a.contour, b.contour)
	})

	// Sort children by the lowest, leftmost point in their contfours
	slices.SortFunc(pt.children, func(a, b *PolyTree[T]) int {
		// Compare contours using their lowest, leftmost points
		return compareLowestLeftmost(a.contour, b.contour)
	})
}

// Parent retrieves the parent of the current PolyTree node.
//
// The parent node represents the polygon in the hierarchy that contains the current node.
// If the current node is the root, this function returns nil.
//
// Returns:
//   - *PolyTree[T]: The parent node of the current PolyTree, or nil if the current node is the root.
func (pt *PolyTree[T]) Parent() *PolyTree[T] {
	return pt.parent
}

// Perimeter calculates the total perimeter of the PolyTree contour.
//
// This function computes the sum of the lengths of all edges that make up the PolyTree's contour.
// The calculation can be customized using the optional parameters.
//
// Parameters:
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     For example, you can use [WithEpsilon] to adjust precision for floating-point comparisons.
//
// Returns:
//   - float64: The total perimeter of the PolyTree.
//
// Notes:
//   - The perimeter only considers the contour of the PolyTree itself and does not include any child polygons.
func (pt *PolyTree[T]) Perimeter(opts ...Option) float64 {
	var length float64
	for _, edge := range pt.Edges() {
		length += edge.Length(opts...)
	}
	return length
}

// PolygonType returns the type of the PolyTree's polygon.
//
// This function identifies whether the polygon represented by the PolyTree
// is a solid polygon or a hole, based on the [PolygonType] enumeration.
//
// Returns:
//   - [PolygonType]: The type of the polygon (e.g., [PTSolid] or [PTHole]).
//
// Notes:
//   - A solid polygon ([PTSolid]) represents a filled area.
//   - A hole polygon ([PTHole]) represents a void inside a parent polygon.
func (pt *PolyTree[T]) PolygonType() PolygonType {
	return pt.polygonType
}

// RelationshipToCircle determines the spatial relationship between a PolyTree and a [Circle].
//
// This method evaluates the relationship between the calling PolyTree (pt) and the specified [Circle] (c)
// for each polygon in the PolyTree. The relationships include containment, intersection, and disjoint.
//
// Parameters:
//   - c ([Circle][T]): The [Circle] to evaluate against the PolyTree.
//   - opts ([Option]): A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and distances,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the PolyTree, the function determines whether the [Circle] lies within, intersects, or
//     is disjoint from the polygon.
//   - The containment relationship is flipped to reflect the PolyTree's relationship to the [Circle].
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each key is a polygon in the PolyTree, and the value is the
//     relationship between the PolyTree and the [Circle].
//
// Notes:
//   - This method assumes that the PolyTree is valid and non-degenerate.
//   - The flipped containment ensures that the returned relationships describe the PolyTree's relationship
//     to the [Circle], rather than the Circle's relationship to the PolyTree.
func (pt *PolyTree[T]) RelationshipToCircle(c Circle[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := c.RelationshipToPolyTree(pt, opts...)
	for k := range output {
		output[k] = output[k].flipContainment()
	}
	return output
}

// RelationshipToLineSegment determines the spatial relationship between a PolyTree and a [LineSegment].
//
// This method evaluates the relationship between the calling PolyTree (pt) and the specified [LineSegment] (l)
// for each polygon in the PolyTree. The relationships include containment, intersection, and disjoint.
//
// Parameters:
//   - l ([LineSegment][T]): The [LineSegment] to evaluate against the PolyTree.
//   - opts ([Option]): A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the PolyTree, the function determines whether the [LineSegment] lies within, intersects, or
//     is disjoint from the polygon.
//   - The containment relationship is flipped to reflect the PolyTree's relationship to the [LineSegment].
//
// Returns:
//   - map[*PolyTree[T]][Relationship]: A map where each key is a polygon in the PolyTree, and the value is the
//     relationship between the PolyTree and the [LineSegment].
//
// Notes:
//   - This method assumes that the PolyTree is valid and non-degenerate.
//   - The flipped containment ensures that the returned relationships describe the PolyTree's relationship
//     to the [LineSegment], rather than the LineSegment's relationship to the PolyTree.
func (pt *PolyTree[T]) RelationshipToLineSegment(l LineSegment[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := l.RelationshipToPolyTree(pt, opts...)
	for k := range output {
		output[k] = output[k].flipContainment()
	}
	return output
}

// RelationshipToPoint determines the spatial relationship between a PolyTree and a [Point].
//
// This method evaluates the relationship between the calling PolyTree (pt) and the specified [Point] (p),
// for each polygon in the PolyTree. The relationships include containment, intersection (point on edge or vertex),
// and disjoint.
//
// Parameters:
//   - p ([Point][T]): The point to evaluate against the PolyTree.
//   - opts ([Option]): A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the PolyTree, the function determines whether the point lies within, intersects, or is
//     disjoint from the polygon.
//   - The containment relationship is flipped so that the relationship reflects the PolyTree's relationship to the point.
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each key is a polygon in the PolyTree, and the value is the
//     relationship between the PolyTree and the point.
//
// Notes:
//   - This method assumes that the PolyTree is valid and non-degenerate.
//   - The flipped containment ensures that the returned relationships describe the PolyTree's relationship
//     to the point, rather than the point's relationship to the PolyTree.
func (pt *PolyTree[T]) RelationshipToPoint(p Point[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := p.RelationshipToPolyTree(pt, opts...)
	for k := range output {
		output[k] = output[k].flipContainment()
	}
	return output
}

// RelationshipToPolyTree determines the spatial relationship between the polygons in one [PolyTree] (pt)
// and the polygons in another [PolyTree] (other).
//
// This function evaluates pairwise relationships between all polygons in the two PolyTrees
// and returns a map. Each key in the outer map corresponds to a polygon from the first [PolyTree] (pt),
// and its value is another map where the keys are polygons from the second [PolyTree] (other) and the values
// are the relationships between the two polygons.
//
// Relationships include:
//   - [RelationshipEqual]: The two polygons are identical.
//   - [RelationshipIntersection]: The polygons overlap partially.
//   - [RelationshipContains]: A polygon in the first [PolyTree] completely encloses a polygon
//     in the second [PolyTree].
//   - [RelationshipContainedBy]: A polygon in the first [PolyTree] is completely enclosed
//     by a polygon in the second [PolyTree].
//   - [RelationshipDisjoint]: The two polygons have no overlap or containment.
//
// Parameters:
//   - other (*PolyTree[T]): The other [PolyTree] to compare against [PolyTree] pt.
//   - opts ([Option]): A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - map[*[PolyTree][T]]map[*[PolyTree][T]][Relationship]: A nested map where the first key is a polygon from
//     the first [PolyTree] (pt), the second key is a polygon from the second [PolyTree] (other), and the value
//     represents their spatial relationship.
//
// Notes:
//   - For efficiency, the function first checks for equality and then evaluates intersection and containment.
//   - The function assumes that the input PolyTrees have properly formed contours and edges.
func (pt *PolyTree[T]) RelationshipToPolyTree(other *PolyTree[T], opts ...Option) map[*PolyTree[T]]map[*PolyTree[T]]Relationship {
	output := make(map[*PolyTree[T]]map[*PolyTree[T]]Relationship, pt.Len())

	for ptPoly := range pt.Nodes {
		output[ptPoly] = make(map[*PolyTree[T]]Relationship, other.Len())

	RelationshipToPolyTreeOtherIterPolys:
		for otherPoly := range other.Nodes {

			ptPolyInsideOtherPoly := true
			otherPolyInsidePtPoly := true

			// check equality
			if ptPoly.contour.eq(otherPoly.contour) {
				output[ptPoly][otherPoly] = RelationshipEqual
				continue RelationshipToPolyTreeOtherIterPolys
			}

			for ptPolyEdge := range ptPoly.contour.iterEdges {
				for otherPolyEdge := range otherPoly.contour.iterEdges {

					// check intersection
					rel := ptPolyEdge.RelationshipToLineSegment(otherPolyEdge, opts...)
					if rel == RelationshipIntersection || rel == RelationshipEqual {
						output[ptPoly][otherPoly] = RelationshipIntersection
						continue RelationshipToPolyTreeOtherIterPolys
					}

					// check containment: otherPoly inside ptPoly
					if !ptPoly.contour.isPointInside(otherPolyEdge.start) || !ptPoly.contour.isPointInside(otherPolyEdge.end) {
						otherPolyInsidePtPoly = false
					}
				}

				// check containment: ptPoly inside otherPoly
				if !otherPoly.contour.isPointInside(ptPolyEdge.start) || !otherPoly.contour.isPointInside(ptPolyEdge.end) {
					ptPolyInsideOtherPoly = false
				}
			}

			// check containment
			if ptPolyInsideOtherPoly {
				output[ptPoly][otherPoly] = RelationshipContainedBy
				continue RelationshipToPolyTreeOtherIterPolys
			}
			if otherPolyInsidePtPoly {
				output[ptPoly][otherPoly] = RelationshipContains
				continue RelationshipToPolyTreeOtherIterPolys
			}

			// otherwise disjoint
			output[ptPoly][otherPoly] = RelationshipDisjoint
		}
	}
	return output
}

// RelationshipToRectangle computes the relationship between the given [Rectangle] (r) and each polygon in the [PolyTree] (pt).
//
// The function evaluates whether the rectangle is disjoint from, intersects with, contains, or is contained by
// each polygon in the [PolyTree]. It uses [Rectangle.RelationshipToPolyTree] to determine the relationship and
// flips the containment relationships so that they are expressed from the [PolyTree]'s perspective.
//
// Parameters:
//   - r ([Rectangle][T]): The rectangle to evaluate against the polygons in the [PolyTree].
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for geometric calculations, improving robustness
//     against floating-point imprecision.
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where the keys are pointers to the polygons in the [PolyTree],
//     and the values are the relationships of the rectangle to those polygons, with containment relationships
//     flipped to reflect the PolyTree's perspective.
//
// Notes:
//   - This method assumes that the [PolyTree] contains valid and non-overlapping polygons.
//   - The returned relationships are expressed relative to the [PolyTree], meaning containment relationships
//     are flipped.
func (pt *PolyTree[T]) RelationshipToRectangle(r Rectangle[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := r.RelationshipToPolyTree(pt, opts...)
	for k := range output {
		output[k] = output[k].flipContainment()
	}
	return output
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
func (pt *PolyTree[T]) resetIntersectionMetadataAndReorder() {
	// Iterate over all polygons in the PolyTree, including nested ones
	for poly := range pt.Nodes {

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
	pt.contour.reorder()
}

// Root returns the topmost parent node of the current PolyTree.
// This is useful for obtaining the root node of a nested PolyTree structure.
//
// Behavior:
//   - Starts at the current node and traverses upward through the parent chain
//     until a node with no parent is found.
//   - If the current node is already the root, it returns the current node.
//
// Returns:
//   - *PolyTree[T]: The root node of the PolyTree.
func (pt *PolyTree[T]) Root() *PolyTree[T] {
	// Start with the current node
	currentPoly := pt

	// Traverse upwards until a node with no parent is found
	for currentPoly.parent != nil {
		currentPoly = currentPoly.parent
	}

	// Return the topmost node
	return currentPoly
}

// Rotate rotates all points in the [PolyTree] around a specified pivot [Point] by a given angle in radians,
// in counter-clockwise direction.
//
// This method applies a rotation transformation to all points in the [PolyTree], including its contours and
// nested child polygons. The rotation is performed relative to the specified pivot point and is expressed
// in radians.
//
// Parameters:
//   - pivot ([Point][T]): The [Point] around which all points in the [PolyTree] are rotated.
//   - radians (float64): The angle of rotation in radians. Positive values indicate counter-clockwise rotation.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for geometric calculations, improving robustness
//     against floating-point imprecision.
//
// Returns:
//   - *[PolyTree][float64]: A new [PolyTree] with points rotated as float64.
//
// Panics:
//   - This function panics if the internal point transformation logic fails.
func (pt *PolyTree[T]) Rotate(pivot Point[T], radians float64, opts ...Option) *PolyTree[float64] {
	out, err := ApplyPointTransform(pt, func(p Point[T]) (Point[float64], error) {
		return p.Rotate(pivot, radians, opts...), nil
	})
	if err != nil {
		panic(err)
	}
	return out
}

// Scale scales all points in the PolyTree relative to a reference point by a given scaling factor.
//
// This method scales all the points of the PolyTree's contours and its children relative to the reference [Point] ref
// by the scaling factor k. It maintains the relationships between parent and child polygons and adjusts their points
// proportionally.
//
// Parameters:
//   - ref ([Point][T]): The reference point used as the centre of scaling.
//   - k (T): The scaling factor. A value greater than 1 enlarges the [PolyTree], a value between 0 and 1 shrinks it,
//     and a negative value mirrors it.
//
// Returns:
//   - *[PolyTree][T]: A new [PolyTree] where all the points have been scaled relative to ref.
//
// Panics:
//   - This function panics if the internal point transformation logic fails.
func (pt *PolyTree[T]) Scale(ref Point[T], k T) *PolyTree[T] {
	out, err := ApplyPointTransform(pt, func(p Point[T]) (Point[T], error) {
		return p.Scale(ref, k), nil
	})
	if err != nil {
		panic(err)
	}
	return out
}

// Siblings returns a slice of sibling polygons of the current PolyTree.
//
// This function provides access to polygons at the same hierarchical level as the current PolyTree.
// The siblings are polygons that share the same parent as the current PolyTree.
//
// Returns:
//   - []*PolyTree[T]: A slice of pointers to sibling PolyTrees.
//
// Notes:
//   - If the current PolyTree has no siblings, an empty slice is returned.
//   - The slice does not include the current PolyTree itself.
func (pt *PolyTree[T]) Siblings() []*PolyTree[T] {
	return pt.siblings
}

// String returns a string representation of the PolyTree, displaying its hierarchy,
// polygon type, and contour points.
//
// This method uses the Nodes method to traverse the entire PolyTree
// and represent each polygon's type, contour points, and relationships.
//
// Returns:
//   - string: A human-readable representation of the PolyTree.
func (pt *PolyTree[T]) String() string {
	var builder strings.Builder

	for poly := range pt.Nodes {

		// Calculate indentation level based on hierarchy depth
		depth := 0
		for parent := poly.parent; parent != nil; parent = parent.parent {
			depth++
		}
		indent := strings.Repeat("  ", depth)

		// Write polygon information
		builder.WriteString(fmt.Sprintf(
			"%sPolyTree: %s\n%s%s\n",
			indent,
			poly.polygonType.String(),
			indent,
			poly.contour.String(),
		))
	}

	return builder.String()
}

// Translate moves all points in the [PolyTree] by a specified displacement vector (given as a [Point]).
//
// This method applies the given delta vector to all points in the [PolyTree], including its contours
// and any nested child polygons. The displacement vector is added to each point's coordinates.
//
// Parameters:
//   - delta ([Point][T]): The displacement vector to apply to all points.
//
// Returns:
//   - *[PolyTree][T]: A new [PolyTree] where all the points have been translated.
//
// Panics:
//   - This function panics if the internal point transformation logic fails.
func (pt *PolyTree[T]) Translate(delta Point[T]) *PolyTree[T] {
	out, err := ApplyPointTransform(pt, func(p Point[T]) (Point[T], error) {
		return p.Translate(delta), nil
	})
	if err != nil {
		panic(err)
	}
	return out
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
//	scp := simpleConvexPolygon{Contour: []Point{
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

// ApplyPointTransform applies a transformation function to each point in a [PolyTree] and produces a new [PolyTree]
// with points of a different type. It supports transformation between numeric types in the [SignedNumber] constraint.
//
// This function allows users to convert the points in a PolyTree using a custom transformation function, for example,
// rounding float64 coordinates to int or scaling coordinates to another type.
//
// Parameters:
//   - pt (*[PolyTree][T]): The input [PolyTree] with points of type T.
//   - transformFunc (func([Point][T]) ([Point][N], error)): A user-defined function that transforms a [Point][T] to a [Point][N].
//     It can return an error if the transformation fails.
//
// Returns:
//   - *[PolyTree][N]: A new [PolyTree] where all points have been transformed to type N.
//   - error: Returns an error if any point transformation fails or if the new [PolyTree] cannot be created.
func ApplyPointTransform[T, N SignedNumber](pt *PolyTree[T], transformFunc func(p Point[T]) (Point[N], error)) (*PolyTree[N], error) {
	var (
		err   error        // Stores any errors encountered during transformation.
		ptOut *PolyTree[N] // The resulting PolyTree after transformation.
	)

	// Prepare a slice of transformed contours (polygons), where each point will be converted.
	contours := make([][]Point[N], pt.Len())

	// Iterate over each polygon (node) in the PolyTree.
	i := 0
	for poly := range pt.Nodes {

		// Prepare a slice for the transformed points of the current contour.
		polyContour := make([]Point[N], len(poly.contour))

		// Transform each point in the contour using the provided transformation function.
		for j, point := range poly.contour.toPoints() {
			polyContour[j], err = transformFunc(point)
			if err != nil {
				// Return a descriptive error if the transformation fails for any point.
				return nil, fmt.Errorf("failed to transform point at poly %d, index %d: %w", i, j, err)
			}
		}

		// Store the transformed contour into the output contours slice.
		contours[i] = polyContour
		i++
	}

	// Reconstruct the PolyTree from the transformed contours.
	ptOut, err = nestPointsToPolyTrees[N](contours)
	if err != nil {
		// Return an error if the new PolyTree cannot be created.
		return nil, fmt.Errorf("failed to nest points into PolyTree: %w", err)
	}

	// Return the transformed PolyTree.
	return ptOut, nil
}

// compareLowestLeftmost compares two contours and determines their relative order
// based on the lowest, leftmost point in each contour.
//
// It is intended for use with slices.SortFunc.
//
// The comparison is performed as follows:
//   - The contour with the point that has the smallest y coordinate is considered "smaller".
//   - If the y coordinates are equal, the contour with the smallest x coordinate is considered "smaller".
//   - Returns -1 if a should come before b, 1 if b should come before "a", or 0 if both contours have the same lowest, leftmost point.
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
//   - contours: A slice of slices of Contour, where each inner slice represents a closed polygon contour.
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

	// Step 1: Sort polygons by area
	slices.SortFunc(contours, sortPointsByAreaDescending)

	// Step 2: Create the root PolyTree from the largest polygon
	rootTree, err := NewPolyTree(contours[0], PTSolid)
	if err != nil {
		return nil, fmt.Errorf("failed to create root PolyTree: %w", err)
	}

	// Step 3: Process the remaining polygons
	for i := 1; i < len(contours); i++ {

		// create new poly
		newPoly, err := NewPolyTree(contours[i], PTSolid)
		if err != nil {
			return nil, fmt.Errorf("failed to create PolyTree: %w", err)
		}

		// find where new poly fits
		var childToPoly *PolyTree[T]
		for existingPoly := range rootTree.Nodes {

			// if poly fits inside, then set child
			if existingPoly.contour.isContourInside(newPoly.contour) {
				childToPoly = existingPoly
			}
		}

		// place poly as child
		if childToPoly != nil {
			// set newPoly's polygonType
			switch childToPoly.polygonType {
			case PTSolid:
				newPoly.polygonType = PTHole
				newPoly.contour.ensureClockwise()
			case PTHole:
				newPoly.polygonType = PTSolid
				newPoly.contour.ensureCounterClockwise()
			}

			// add as sibling to root
			err := childToPoly.AddChild(newPoly)
			if err != nil {
				return nil, fmt.Errorf("failed to create PolyTree relationship: %w", err)
			}

		} else {
			// add poly as root sibling
			newPoly.polygonType = rootTree.polygonType
			err := rootTree.AddSibling(newPoly)
			if err != nil {
				return nil, fmt.Errorf("failed to create PolyTree relationship: %w", err)
			}
		}
	}

	// Step 4: Return the root PolyTree
	return rootTree, nil
}

func sortPointsByAreaDescending[T SignedNumber](a, b []Point[T]) int {

	// get signed areas
	areaA := SignedArea2X(a)
	areaB := SignedArea2X(b)

	// get absolute values
	if areaA < 0 {
		areaA *= -1
	}
	if areaB < 0 {
		areaB *= -1
	}

	// return expected values to sort by area, descending
	switch {
	case areaA > areaB:
		return -1
	case areaA < areaB:
		return 1
	default:
		return 0
	}
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
