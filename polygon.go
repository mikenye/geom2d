// File polygon.go defines the Polygon type and implements methods for operations on polygons in 2D space.
//
// This file includes the definition of the Polygon type, which supports complex structures such as
// polygons with holes and islands. It provides methods for common operations such as area calculations,
// determining point-in-polygon relationships, and managing hierarchical relationships between polygons.
//
// Functions and methods in this file support:
// - Calculating properties of polygons, such as area and centroid.
// - Determining spatial relationships of points to polygons.
// - Managing nested structures, such as holes and islands, in polygons.
//
// The Polygon type is a foundational structure in the geom2d package for 2D geometric computations.

package geom2d

import (
	"fmt"
	"slices"
)

// polyTraversalDirection defines the direction in which a polygon's vertices are traversed.
// This can either be clockwise or counterclockwise and is used to specify how to iterate
// through a polygon's vertices or edges during boolean operations and other processing.
type polyTraversalDirection int

const (
	// polyTraversalDirectionCounterClockwise specifies that the traversal proceeds
	// in a counterclockwise direction through the polygon's vertices or edges.
	polyTraversalDirectionCounterClockwise = polyTraversalDirection(iota)

	// polyTraversalDirectionClockwise specifies that the traversal proceeds
	// in a clockwise direction through the polygon's vertices or edges.
	polyTraversalDirectionClockwise
)

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

// polyPoint represents a point in a polygon, with additional metadata to support
// advanced polygon operations such as Boolean operations (union, intersection, subtraction)
// and traversal algorithms.
//
// This type is used internally to track details about polygon vertices, intersection points,
// and their relationships with other polygons during geometric computations.
type polyPoint[T SignedNumber] struct {
	// point defines the geometric coordinates of the point in 2D space.
	point Point[T]

	// pointType indicates the type of the point, which can represent a normal vertex,
	// an intersection point, or a midpoint between intersections.
	pointType polyPointType

	// entryExit specifies whether this point marks an entry to or exit from the "area of interest" during polygon traversal.
	// This is relevant for Boolean operations where traversal directions are critical.
	entryExit polyIntersectionType

	// visited tracks whether this point has already been visited during traversal algorithms.
	// This helps prevent redundant processing during operations such as polygon traversal.
	visited bool

	// otherPolyIndex stores the index of the corresponding point in another polygon when this point is an intersection.
	// This is useful for linking intersection points between polygons.
	otherPolyIndex int
}

// polyPointType defines the type of point in a polygon, used to distinguish between
// original vertices and additional points introduced during polygon operations.
//
// This type is essential for managing polygon data during Boolean operations (e.g., union,
// intersection, subtraction) and other algorithms that require distinguishing original
// points from dynamically added points.
type polyPointType int

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

// BooleanOperation defines the types of Boolean operations that can be performed on polygons.
// These operations are fundamental in computational geometry for combining or modifying shapes.
//
// The supported operations are:
// - Union: Combines two polygons into one, merging their areas.
// - Intersection: Finds the overlapping region between two polygons.
// - Subtraction: Subtracts one polygon's area from another.
type BooleanOperation int

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

// polyIntersectionType defines the type of intersection point in polygon operations,
// distinguishing between entry and exit (of area of interest) points during traversal.
//
// This type is primarily used in Boolean operations (e.g., union, intersection, subtraction)
// to identify transitions at intersection points between polygons.
type polyIntersectionType int

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
type PointPolygonRelationship int

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

// PolygonType (PT) defines whether the inside of the contour of a polygon represents either a solid region (island)
// or a void region (hole). This distinction is essential for operations involving polygons
// with complex structures, such as those containing holes or nested islands.
type PolygonType int

const (
	// PTSolid represents a solid region of the polygon, commonly referred to as an "island."
	// PTSolid polygons are the primary filled areas, excluding any void regions (holes).
	PTSolid PolygonType = iota

	// PTHole represents a void region of the polygon, often nested within a solid polygon.
	// Holes are not part of the filled area of the polygon and are treated as exclusions.
	PTHole
)

// Polygon represents a 2D polygon, which can contain nested structures such as holes and islands.
// The Polygon type supports operations such as determining point-in-polygon relationships,
// calculating properties like area and centroid, and managing hierarchical relationships.
//
// The generic parameter T must satisfy the SignedNumber constraint, allowing the Polygon
// to handle various numeric types (e.g., int, float32, float64).
type Polygon[T SignedNumber] struct {
	// points defines the full outline of the polygon, including all vertices and relevant metadata.
	// Each entry is a polyPoint, which tracks whether the point is a regular vertex, an intersection,
	// or a midpoint between intersections.
	points []polyPoint[T]

	// polygonType indicates whether the polygon is a solid region (PTSolid) or a hole (PTHole).
	// This classification is essential for distinguishing between filled and void areas of the polygon.
	polygonType PolygonType

	// children contains references to nested polygons, which may represent holes (if this polygon
	// is a solid region) or solid islands (if this polygon is a hole). These hierarchical relationships
	// allow for complex polygon structures.
	children []*Polygon[T]

	// parent points to the parent polygon of this polygon, if any. For example, a hole's parent
	// would be the solid polygon that contains it.
	parent *Polygon[T]

	// hull optionally stores the convex hull of the polygon, represented as a simpleConvexPolygon.
	// This can be used to optimize certain operations, such as point-in-polygon checks, by
	// quickly ruling out points that lie outside the convex hull.
	hull *simpleConvexPolygon[T]

	// maxX stores the maximum X-coordinate value among the polygon's vertices. This is used
	// for ray-casting operations to determine point-in-polygon relationships.
	maxX T
}

// Points returns a slice of all the vertex points that define the polygon's outline.
// This method extracts points marked as either normal vertices or intersection points
// from the polygon's internal representation and converts them to their geometric coordinates.
//
// The method performs the following steps:
//  1. Iterates over the internal polyPoint representation of the polygon.
//  2. Filters out points that are not of type pointTypeOriginal or pointTypeOriginalAndIntersection.
//  3. Converts the valid points into Point objects with adjusted coordinates (halving x and y values).
//  4. Returns the resulting slice of Point objects.
//
// Returns:
//   - []Point[T]: A slice of Points representing the polygon's outline.
func (p Polygon[T]) Points() []Point[T] {
	points := make([]Point[T], 0, len(p.points))
	for _, point := range p.points {
		if !(point.pointType == pointTypeOriginal || point.pointType == pointTypeOriginalAndIntersection) {
			continue
		}
		points = append(points, Point[T]{
			x: point.point.x / 2,
			y: point.point.y / 2,
		})
	}
	return points
}

// RelationshipToPoint determines the spatial relationship of a given point to the polygon.
// This includes checking if the point is outside, inside, on a vertex, or on an edge of the polygon.
// It also accounts for the polygon's hierarchical structure, such as holes and nested islands.
//
// The method uses the following steps:
//  1. Preliminary Convex Hull Check: If the polygon has a convex hull defined, the point is first checked
//     against the hull. If the point is outside the hull, it is immediately classified as outside the polygon.
//  2. Ray-Casting Algorithm: A horizontal ray is cast to the right from the given point, and intersections
//     with the polygon's edges are counted to determine if the point is inside or outside.
//  3. Edge and Vertex Checks: The method explicitly checks for cases where the point lies on a vertex
//     or on an edge of the polygon.
//  4. Recursive Relationship with Children: If the polygon has nested children (holes or islands),
//     the method recursively determines the most specific relationship with the lowest-level child polygon.
//
// Parameters:
//
//   - point Point[T]: The point to check the relationship for.
//
// Returns a PointPolygonRelationship (PPR): An enumerated value representing the relationship of the point to the polygon:
//   - PPRPointOutside: The point is outside the polygon.
//   - PPRPointInside: The point is inside the polygon.
//   - PPRPointOnVertex: The point coincides with a vertex of the polygon.
//   - PPRPointOnEdge: The point lies on an edge of the polygon.
//   - PPRPointInHole: The point is inside a hole within the polygon.
//   - PPRPointInsideIsland: The point is inside an island nested within the polygon.
//
// Example:
//
//	polygon := NewPolygon(points)
//	relationship := polygon.RelationshipToPoint(NewPoint(1, 2))
//	fmt.Println(relationship) // Outputs the spatial relationship of the point to the polygon.
func (p Polygon[T]) RelationshipToPoint(point Point[T]) PointPolygonRelationship {
	// Use the convex hull as a preliminary check for outside points.
	if p.hull != nil && !p.hull.ContainsPoint(point) {
		return PPRPointOutside
	}

	// Get the points from the polygon
	ppoints := p.Points()

	// Set up the ray for the ray-casting algorithm.
	ray := NewLineSegment(point, Point[T]{x: p.maxX + 2, y: point.y}) // Horizontal ray to the right

	intersections := 0

	for i := 0; i < len(ppoints); i++ {
		start := ppoints[i]
		end := ppoints[(i+1)%len(ppoints)]
		edge := NewLineSegment(start, end)

		// Ray-casting intersection check using RelationshipToLineSegment.
		switch edge.RelationshipToLineSegment(ray) {
		case LSRAeqC, LSRAeqD, LSRBeqC, LSRBeqD:
			return PPRPointOnVertex
		case LSRAonCD, LSRBonCD, LSRConAB, LSRDonAB, LSRIntersects:
			intersections++
		case LSRCollinearAonCD, LSRCollinearBonCD, LSRCollinearEqual:
			// If the ray and edge are collinear and overlap, the point is on the edge.
			return PPRPointOnEdge
		}
	}

	// Determine if the point is inside or outside based on intersections
	isInside := intersections%2 == 1
	if !isInside {
		return PPRPointOutside
	}

	// Recursive check for the "lowest" relationship with children
	relationship := PPRPointInside // Start with the main polygon relationship

	for _, child := range p.children {
		childRelationship := child.RelationshipToPoint(point)

		switch childRelationship {
		case PPRPointInside, PPRPointOnEdge, PPRPointOnVertex:
			// If the child polygon is a hole, mark the relationship as PPRPointInHole
			// If the child polygon is a solid (island), mark the relationship as PPRPointInsideIsland
			switch child.polygonType {
			case PTHole:
				relationship = PPRPointInHole
			case PTSolid:
				relationship = PPRPointInsideIsland
			}
		case PPRPointInHole, PPRPointInsideIsland:
			// If the child relationship is already the innermost context, return immediately
			return childRelationship
		}
	}

	return relationship
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

// edgesFromPolyPoints converts a slice of polyPoint objects into a slice of polyEdge objects,
// effectively creating edges from consecutive points in a polygon.
//
// This function is used internally to represent the edges of a polygon, where each edge is
// defined as a line segment connecting two consecutive points.
//
// Parameters:
//   - polygon []polyPoint[T]: A slice of polyPoint objects representing the vertices of the polygon.
//     The points are assumed to be ordered sequentially, forming a closed loop.
//
// Returns:
//   - []polyEdge[T]: A slice of polyEdge objects, where each edge corresponds to a line segment
//     between two consecutive points in the input slice. The last edge connects
//     the last point back to the first point, ensuring the polygon is closed.
//
// Example:
//
//	polygon := []polyPoint{
//	    {point: Point{X: 0, Y: 0}},
//	    {point: Point{X: 1, Y: 0}},
//	    {point: Point{X: 1, Y: 1}},
//	}
//	edges := edgesFromPolyPoints(polygon)
//	// edges will contain:
//	// [
//	//   polyEdge{lineSegment: LineSegment{start: (0,0), end: (1,0)}},
//	//   polyEdge{lineSegment: LineSegment{start: (1,0), end: (1,1)}},
//	//   polyEdge{lineSegment: LineSegment{start: (1,1), end: (0,0)}}
//	// ]
//
// Notes:
//   - This function assumes the input slice represents a valid polygon outline.
//   - The returned edges form a closed loop by connecting the last point to the first point.
func edgesFromPolyPoints[T SignedNumber](polygon []polyPoint[T]) []polyEdge[T] {
	edges := make([]polyEdge[T], 0, len(polygon))
	for i := range polygon {
		j := (i + 1) % len(polygon)
		edges = append(edges, polyEdge[T]{
			lineSegment: NewLineSegment(polygon[i].point, polygon[j].point),
		})
	}
	return edges
}

// findIntersectionsBetweenPolys identifies and inserts intersection points between the edges
// of two polygons. The function modifies the input polygons by adding intersection points
// as new vertices, ensuring proper ordering along their respective edges.
//
// Parameters:
//
//   - poly1 []polyPoint[T]: A slice of polyPoint objects representing the first polygon.
//     The points are assumed to be ordered sequentially, forming a closed polygon.
//   - poly2 []polyPoint[T]: A slice of polyPoint objects representing the second polygon.
//     The points are assumed to be ordered sequentially, forming a closed polygon.
//
// Returns:
//
//   - ([]polyPoint[T], []polyPoint[T]): The updated slices of polyPoint objects for both polygons,
//     with intersection points added as vertices.
//
// Behavior:
//
//   - The function performs sanity checks to ensure both polygons have at least 3 points.
//   - It iterates over all edges of poly1 and poly2 to check for intersections.
//   - For each intersection found, the function firstly checks if the intersection point is already present
//     in either polygon to avoid duplicates; secondly, it converts the intersection point into a polyPoint
//     object with the type `pointTypeAddedIntersection`; lastly, it inserts the intersection into the appropriate
//     position in both polygons, maintaining vertex order.
//   - The function increments indices after inserting an intersection to avoid re-processing it.
//   - The resulting polygons include all original points and any newly added intersection points.
//
// Notes:
//   - This function assumes that all original polyPoint coordinates are doubled during the construction
//     of polygons (e.g., in NewPolygon). This ensures that when casting intersection points to the generic
//     type T (which could be an integer), any fractional values are dropped without precision loss.
//   - This function uses NewLineSegment to define polygon edges and relies on the IntersectionPoint
//     method to calculate intersections between line segments.
//   - Intersection points are added in sorted order along their respective edges to maintain polygon structure.
//   - If either polygon has fewer than 3 points, the function panics.
//
// Example:
//
//	poly1 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 4}},
//	}
//	poly2 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 2, Y: -1}},
//	    {point: Point[float64]{X: 2, Y: 5}},
//	    {point: Point[float64]{X: 6, Y: 2}},
//	}
//	updatedPoly1, updatedPoly2 := findIntersectionsBetweenPolys(poly1, poly2)
//	// Both updatedPoly1 and updatedPoly2 will include intersection points as new vertices.
//
// Panics:
//   - If either poly1 or poly2 has fewer than 3 points.
//
// Dependencies:
//   - This function uses helper functions such as isPointInPolyPoints to avoid duplicate intersections
//     and insertPolyIntersectionSorted to correctly place intersection points in the polygon's vertex list.
//
// todo: show the output of updatedPoly1 & updatedPoly2 in the example above
func findIntersectionsBetweenPolys[T SignedNumber](poly1, poly2 []polyPoint[T]) ([]polyPoint[T], []polyPoint[T]) {
	// Sanity checks
	if len(poly1) < 3 || len(poly2) < 3 {
		panic("Polygons must have at least 3 points")
	}

	// Iterate through each edge in poly1
	for i1 := 0; i1 < len(poly1); i1++ {
		j1 := (i1 + 1) % len(poly1) // Wrap around to form a closed polygon
		segment1 := NewLineSegment(poly1[i1].point, poly1[j1].point)

		// Iterate through each edge in poly2
		for i2 := 0; i2 < len(poly2); i2++ {
			j2 := (i2 + 1) % len(poly2)
			segment2 := NewLineSegment(poly2[i2].point, poly2[j2].point)

			// Check for intersection between the segments
			intersectionPoint, intersects := segment1.IntersectionPoint(segment2)
			intersectionPointT := NewPoint(T(intersectionPoint.x), T(intersectionPoint.y))
			if intersects {
				// Check if the intersection point already exists in poly1 or poly2
				if isPointInPolyPoints(intersectionPointT, poly1) || isPointInPolyPoints(intersectionPointT, poly2) {
					continue // Skip adding duplicate intersections
				}

				// Convert the intersection point to a polyPoint
				intersection := polyPoint[T]{
					point:     NewPoint(T(intersectionPoint.x), T(intersectionPoint.y)),
					pointType: pointTypeAddedIntersection, // Mark as intersection
				}

				// Insert intersection into both polygons
				poly1 = insertPolyIntersectionSorted(poly1, i1, j1, intersection)
				poly2 = insertPolyIntersectionSorted(poly2, i2, j2, intersection)

				// Increment indices to avoid re-processing this intersection
				i1++
				i2++
			}
		}
	}

	return poly1, poly2
}

// findTraversalStartingPoint locates the first unvisited entry point in a collection of polygons,
// which serves as the starting point for a traversal algorithm.
//
// This function is typically used in Boolean operations (e.g., union, intersection, subtraction)
// to find the starting point for traversing polygon edges and intersections.
//
// Parameters:
//
//   - polys [][]polyPoint[T]: A slice of polygons, where each polygon is represented as a slice of polyPoint objects.
//     Each polyPoint may include metadata, such as whether it is an entry or exit point.
//
// Returns (int, int): A tuple representing the indices of the polygon and the point within that polygon:
//   - The first index corresponds to the polygon in the polys slice.
//   - The second index corresponds to the point within that polygon.
//   - If no unvisited entry points are found, the function returns (-1, -1).
//
// Behavior:
//   - The function iterates over each polygon in the input slice and over each point in the polygon.
//   - For each point, it checks if the point's entryExit field is intersectionTypeEntry and if the
//     point's `visited` field is false (i.e., the point has not been processed during a traversal).
//   - If a matching point is found, the function immediately returns its indices.
//   - If no unvisited entry points are found after scanning all polygons, the function returns (-1, -1).
//
// Notes:
//   - This function assumes that the entryExit and visited fields of polyPoint objects are correctly set
//     before calling it.
//   - It is the caller's responsibility to handle the case where the function returns (-1, -1).
//
// Example:
//
//	polys := [][]polyPoint[float64]{
//	    {
//	        {entryExit: intersectionTypeEntry, visited: false}, // Entry point
//	        {entryExit: intersectionTypeExit, visited: false},
//	    },
//	    {
//	        {entryExit: intersectionTypeEntry, visited: true},  // Already visited
//	        {entryExit: intersectionTypeExit, visited: false},
//	    },
//	}
//	polyIndex, pointIndex := findTraversalStartingPoint(polys)
//	// polyIndex == 0, pointIndex == 0
//
//	polyIndex, pointIndex = findTraversalStartingPoint([][]polyPoint[float64]{})
//	// polyIndex == -1, pointIndex == -1
func findTraversalStartingPoint[T SignedNumber](polys [][]polyPoint[T]) (int, int) {
	for polyIndex := range polys {
		for pointIndex := range polys[polyIndex] {
			if polys[polyIndex][pointIndex].entryExit == intersectionTypeEntry && !polys[polyIndex][pointIndex].visited {
				return polyIndex, pointIndex
			}
		}
	}
	// Return -1 if no Entry points are found
	return -1, -1
}

// insertPolyIntersectionSorted inserts an intersection point into a polygon's vertex list
// at the correct position, maintaining the order of points along the affected edge.
//
// This function ensures that the intersection point is placed between the start and end
// points in the polygon's edge, preserving the logical order of the polygon.
//
// Parameters:
//
//   - poly []polyPoint[T]: The slice of polyPoint objects representing the polygon's vertices.
//     The vertices are assumed to form a closed polygon.
//   - start int: The index of the starting vertex of the edge where the intersection occurred.
//   - end int: The index of the ending vertex of the edge where the intersection occurred.
//   - intersection polyPoint[T]: The intersection point to be inserted. This point should be
//     a valid polyPoint with appropriate metadata.
//
// Returns:
//
//   - []polyPoint[T]: A new slice of polyPoint objects with the intersection point inserted
//     at the correct position between the `start` and `end` vertices.
//
// Behavior:
//   - The function identifies the correct position for the intersection point by measuring
//     distances along the edge defined by start and end.
//   - It iterates over intermediate points between start and end to find the insertion position
//     such that the order of points along the edge remains consistent.
//   - The intersection point is inserted into the slice using slices.Insert,
//     and the updated polygon is returned.
//
// Notes:
//   - The function assumes that poly represents a valid polygon with at least 3 vertices.
//   - The intersection point must already be calculated and passed as input with relevant metadata.
//   - If multiple intersection points exist along the same edge, this function must be called
//     for each intersection separately.
//
// Example:
//
//	poly := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 4}},
//	}
//	intersection := polyPoint[float64]{
//	    point: Point[float64]{X: 2, Y: 0},
//	    pointType: PointTypeAddedIntersection,
//	}
//	updatedPoly := insertPolyIntersectionSorted(poly, 0, 1, intersection)
//	// updatedPoly will now include the intersection point between (0,0) and (4,0).
//
//	fmt.Println(updatedPoly)
//	// Output:
//	// [
//	//   {point: (0, 0), ...},
//	//   {point: (2, 0), pointType: PointTypeAddedIntersection, ...},
//	//   {point: (4, 0), ...},
//	//   {point: (4, 4), ...}
//	// ]
func insertPolyIntersectionSorted[T SignedNumber](poly []polyPoint[T], start, end int, intersection polyPoint[T]) []polyPoint[T] {
	segment := NewLineSegment(poly[start].point, poly[end].point)

	// Find the correct position to insert the intersection
	insertPos := end
	for i := start + 1; i < end; i++ {
		existingSegment := NewLineSegment(poly[start].point, poly[i].point)
		if segment.DistanceToPoint(intersection.point) < existingSegment.DistanceToPoint(poly[i].point) {
			insertPos = i
			break
		}
	}

	// Insert the intersection at the calculated position
	return slices.Insert(poly, insertPos, intersection)
}

// isInsidePolygon determines whether a given point lies inside a polygon using the ray-casting algorithm.
// The algorithm checks how many times a ray extending to the right from the given point crosses the edges
// of the polygon. An odd number of crossings indicates the point is inside, while an even number
// indicates the point is outside.
//
// Parameters:
//
//   - polygon []polyPoint[T]: A slice of polyPoint objects representing the vertices of the polygon.
//     The points are assumed to form a closed polygon.
//   - point Point[T]: The point to check.
//
// Returns:
//
//   - bool: True if the point lies inside or on the edge of the polygon; false otherwise.
//
// Behavior:
//   - A horizontal ray is cast to the right from the given point.
//   - The function calculates the ray's relationship with each edge of the polygon.
//   - The function handles special cases where the point lies exactly on an edge, treating the point as inside.
//   - The number of ray crossings with polygon edges is counted, with certain configurations (e.g., collinear points)
//     contributing additional crossings as necessary.
//   - If the total number of crossings is odd, the point is inside; otherwise, it is outside.
//
// Notes:
//   - The function assumes the input polygon is valid and forms a closed loop.
//   - The edges of the polygon are constructed using the edgesFromPolyPoints function.
//   - Special cases are handled to ensure robustness, such as when the point lies exactly on an edge.
//
// Example:
//
//	polygon := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 0}},
//	    {point: Point[float64]{X: 4, Y: 4}},
//	    {point: Point[float64]{X: 0, Y: 4}},
//	}
//	inside := isInsidePolygon(polygon, Point[float64]{X: 2, Y: 2}) // Returns true
//	outside := isInsidePolygon(polygon, Point[float64]{X: 5, Y: 2}) // Returns false
//
// Algorithm:
//   - A ray is cast to the right from the point.
//   - Each edge of the polygon is checked for intersection with the ray, using the LineSegment.RelationshipToLineSegment function.
//   - Special cases, such as collinear overlaps and edge endpoints, are handled explicitly to ensure correctness.
//   - The total number of ray crossings is used to determine whether the point is inside the polygon.
//
// Dependencies:
//   - This function relies on helper functions such as `edgesFromPolyPoints` and `RelationshipToLineSegment`
//     for constructing edges and determining ray-edge relationships.
//   - The `inOrder` function is used to verify vertical ordering when handling edge overlap cases.
//
// Example Special Cases:
//   - If the point lies exactly on a polygon edge, it is considered inside.
//   - Collinear points with the ray contribute additional crossings.
//
// Panics:
//   - This function does not validate the input polygon and assumes it is properly closed.
//
// For more information, see [Inclusion of a Point in a Polygon], specifically the "The Crossing Number (cn) method".
//
// [Inclusion of a Point in a Polygon]: https://web.archive.org/web/20130126163405/http://geomalgorithms.com/a03-_inclusion.html
func isInsidePolygon[T SignedNumber](polygon []polyPoint[T], point Point[T]) bool {
	crosses := 0

	// Cast ray from point to right
	maxX := point.x
	for _, p := range polygon {
		maxX = max(maxX, p.point.x)
	}
	maxX++
	ray := NewLineSegment(point, NewPoint(maxX, point.y))

	// Get edges
	edges := edgesFromPolyPoints(polygon)

	// Determine relationship with ray for each edge
	for i := range edges {

		// If point is directly on the edge, it's inside
		if point.IsOnLineSegment(edges[i].lineSegment) {
			return true
		}

		// Store relationship
		edges[i].rel = ray.RelationshipToLineSegment(edges[i].lineSegment)
	}

	// Check for crosses
	for i := range edges {
		iNext := (i + 1) % len(edges)

		switch edges[i].rel {
		case LSRIntersects:
			crosses++

		case LSRCollinearCDinAB:
			crosses += 2

		case LSRConAB:
			crosses++
			if edges[iNext].rel == LSRDonAB {
				if inOrder(edges[i].lineSegment.start.y, point.y, edges[iNext].lineSegment.end.y) {
					crosses++
				}
			}

		case LSRDonAB:
			crosses++
			if edges[iNext].rel == LSRConAB {
				if inOrder(edges[i].lineSegment.start.y, point.y, edges[iNext].lineSegment.end.y) {
					crosses++
				}
			}
		}
	}

	return crosses%2 == 1 // Odd crossings mean the point is inside
}

// isPointInPolyPoints checks if a given point exists in a slice of polyPoint objects.
//
// The function compares the coordinates of the input point with the coordinates of each polyPoint
// in the provided slice. It returns true if the point matches any polyPoint, and false otherwise.
//
// Parameters:
//
//   - point Point[T]: The point to search for. The coordinates of this point are compared directly
//     with the coordinates of each polyPoint in the slice.
//   - poly []polyPoint[T]: A slice of polyPoint objects representing the vertices of a polygon.
//
// Returns:
//
//   - bool: True if the point exists in the slice of polyPoint objects; false otherwise.
//
// Behavior:
//   - The function iterates through each polyPoint in the slice.
//   - For each polyPoint, it compares the x and y coordinates with the input point.
//   - If a match is found, the function immediately returns true.
//   - If no match is found after iterating through all polyPoints, the function returns false.
//
// Notes:
//   - This function uses the generic type T for consistency with the polyPoint and Point types.
//   - The comparison assumes exact equality of coordinates, which may be sensitive to floating-point precision.
//
// Example:
//
//	poly := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}},
//	    {point: Point[float64]{X: 1, Y: 1}},
//	}
//	exists := isPointInPolyPoints(Point[float64]{X: 1, Y: 1}, poly) // Returns true
//	notExists := isPointInPolyPoints(Point[float64]{X: 2, Y: 2}, poly) // Returns false
//
// Performance:
//   - This function performs a linear search, so its complexity is O(n), where n is the number of polyPoints.
//
// todo: need optional epsilon to handle floating-point imprecision
func isPointInPolyPoints[T SignedNumber](point Point[T], poly []polyPoint[T]) bool {
	for _, p := range poly {
		if p.point.x == point.x && p.point.y == point.y {
			return true
		}
	}
	return false
}

// markEntryExitPoints marks entry and exit points between two polygons based on the specified Boolean operation.
//
// This function assigns the appropriate entryExit type (e.g., intersectionTypeEntry or intersectionTypeExit)
// to intersection points in both polygons. It determines whether the traversal at each intersection enters or exits
// the other polygon, based on the relationship between the two polygons and the specified Boolean operation.
//
// Parameters:
//   - poly1 []polyPoint[T]: The first polygon, represented as a slice of polyPoint objects.
//   - poly2 []polyPoint[T]: The second polygon, represented as a slice of polyPoint objects.
//   - operation BooleanOperation: The Boolean operation to perform. Determines how entry and exit points are marked.
//
// operation can be one of:
//   - BooleanUnion: Combines the areas of the two polygons.
//   - BooleanIntersection: Retains only the overlapping region.
//   - BooleanSubtraction: Subtracts the second polygon from the first.
//
// Behavior:
//   - The function iterates through each intersection point in poly1 and poly2.
//   - For each intersection, it determines whether traversal at the intersection point enters or exits
//     the other polygon by checking if the midpoint of the edge following the intersection is inside the other polygon.
//   - Based on the Boolean operation, it marks the entryExit field of the corresponding intersection points
//     in both polygons.
//   - It also updates the otherPolyIndex field to link corresponding intersection points in poly1 and poly2.
//
// Notes:
//   - This function assumes that all original polyPoint coordinates are doubled during the construction
//     of polygons (e.g., in NewPolygon). This ensures that when casting points to the generic type T
//     (which could be an integer), any fractional values are dropped without precision loss.
//   - The function uses helper functions such as isInsidePolygon to determine the containment relationship
//     for midpoints of edges.
//
// Example:
//
//	poly1 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 4, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 4, Y: 4}, pointType: PointTypeAddedIntersection},
//	}
//	poly2 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 2, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 6, Y: 2}, pointType: PointTypeAddedIntersection},
//	    {point: Point[float64]{X: 2, Y: 5}, pointType: PointTypeOriginal},
//	}
//	markEntryExitPoints(poly1, poly2, BooleanUnion)
//	// Intersection points in poly1 and poly2 will have their entryExit fields updated based on the operation.
//
// Dependencies:
//   - This function uses helper functions such as isInsidePolygon and NewLineSegment
//     to calculate relationships between polygons and edges.
//
// Panics:
//   - The function assumes that polygons are valid and that intersection points have been precomputed.
func markEntryExitPoints[T SignedNumber](poly1, poly2 []polyPoint[T], operation BooleanOperation) {
	for poly1Point1Index, poly1Point1 := range poly1 {
		poly1Point2Index := (poly1Point1Index + 1) % len(poly1)
		if poly1Point1.pointType == pointTypeAddedIntersection {
			for poly2PointIndex, poly2Point := range poly2 {
				if poly2Point.pointType == pointTypeAddedIntersection && poly1Point1.point.Eq(poly2Point.point) {

					// determine if poly1 traversal is poly1EnteringPoly2 or existing poly2
					mid := NewLineSegment(
						poly1Point1.point,
						poly1[poly1Point2Index].point).Midpoint()
					midT := NewPoint[T](T(mid.x), T(mid.y))
					poly1EnteringPoly2 := isInsidePolygon(poly2, midT)

					switch operation {
					case BooleanUnion:

						if poly1EnteringPoly2 {
							poly1[poly1Point1Index].entryExit = intersectionTypeExit
							poly2[poly2PointIndex].entryExit = intersectionTypeEntry
						} else {
							poly1[poly1Point1Index].entryExit = intersectionTypeEntry
							poly2[poly2PointIndex].entryExit = intersectionTypeExit
						}

					case BooleanIntersection:

						if poly1EnteringPoly2 {
							poly1[poly1Point1Index].entryExit = intersectionTypeEntry
							poly2[poly2PointIndex].entryExit = intersectionTypeExit
						} else {
							poly1[poly1Point1Index].entryExit = intersectionTypeExit
							poly2[poly2PointIndex].entryExit = intersectionTypeEntry
						}

					case BooleanSubtraction:

						if poly1EnteringPoly2 {
							poly1[poly1Point1Index].entryExit = intersectionTypeExit
							poly2[poly2PointIndex].entryExit = intersectionTypeExit
						} else {
							poly1[poly1Point1Index].entryExit = intersectionTypeEntry
							poly2[poly2PointIndex].entryExit = intersectionTypeEntry
						}

					}
					poly1[poly1Point1Index].otherPolyIndex = poly2PointIndex
					poly2[poly2PointIndex].otherPolyIndex = poly1Point1Index
				}
			}
		}
	}
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
func newSimpleConvexPolygon[T SignedNumber](points []Point[T]) *simpleConvexPolygon[T] {
	// Assume `points` is already ordered to form a convex polygon
	return &simpleConvexPolygon[T]{Points: points}
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

//   - If the input is polyTraversalDirectionClockwise, the output
//     will be polyTraversalDirectionCounterClockwise.
//   - If the input is polyTraversalDirectionCounterClockwise, the output
//     will be polyTraversalDirectionClockwise.
//
// Example:
//
//	currentDirection := polyTraversalDirectionClockwise
//	newDirection := togglePolyTraversalDirection(currentDirection)
//	// newDirection == polyTraversalDirectionCounterClockwise
//
// Notes:
//   - This function assumes that the input is a valid polyTraversalDirection value.
//   - If an invalid value is provided, the function may not behave as expected.
//
// Dependencies:
//   - Relies on the polyTraversalDirection type and its constants for traversal direction.
func togglePolyTraversalDirection(direction polyTraversalDirection) polyTraversalDirection {
	if direction == polyTraversalDirectionClockwise {
		return polyTraversalDirectionCounterClockwise
	}
	return polyTraversalDirectionClockwise
}

// traverse performs a traversal of two polygons (`poly1` and `poly2`) to generate the resulting
// polygons based on the specified Boolean operation. This function handles both nested polygon
// containment scenarios and general intersection-based traversal.
//
// Parameters:
//   - poly1 []polyPoint[T]: The first polygon, represented as a slice of polyPoint objects.
//   - poly2 []polyPoint[T]: The second polygon, represented as a slice of polyPoint objects.
//   - operation BooleanOperation: The Boolean operation to perform.
//
// operation can be one of:
//   - BooleanUnion: Combines the areas of the two polygons.
//   - BooleanIntersection: Retains only the overlapping region.
//   - BooleanSubtraction: Subtracts the second polygon from the first.
//
// Returns:
//
//   - [][]polyPoint[T]: A slice of slices, where each inner slice represents a resulting polygon
//     from the traversal. Each polygon is represented as a slice of polyPoint objects.
//
// # Behavior:
//
// 1. Containment Check:
//   - If one polygon is entirely contained within the other, the function handles the containment case
//     directly without performing a full traversal, returning the appropriate result for the given operation.
//
// 2. Traversal Initialization:
//   - Initializes the traversal direction as counterclockwise and prepares an empty results list.
//
// 3. Traversal Loop:
//   - Finds the starting point for traversal using findTraversalStartingPoint.
//   - Iteratively follows edges of the polygons, switching polygons at intersection points,
//     and optionally reversing traversal direction for subtraction operations.
//   - Stops when a complete loop is formed, appending the resulting path to the results list.
//
// 4. Handle Empty Results:
//   - If no results are generated (e.g., no intersections), the function handles this case
//     appropriately for the given Boolean operation.
//
// Notes:
//   - This function assumes that intersection points between polygons have already been computed
//     and appropriately marked with entryExit values (e.g., using markEntryExitPoints).
//   - The function assumes that all original polyPoint coordinates are doubled during the construction
//     of polygons (e.g., in NewPolygon). This ensures that when casting points to the generic type T
//     (which could be an integer), any fractional values are dropped without precision loss.
//   - Handles scenarios where polygons are disjoint, overlapping, or nested.
//
// Example:
//
//	poly1 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 0, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 4, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 4, Y: 4}, pointType: PointTypeOriginal},
//	}
//	poly2 := []polyPoint[float64]{
//	    {point: Point[float64]{X: 2, Y: 0}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 6, Y: 2}, pointType: PointTypeOriginal},
//	    {point: Point[float64]{X: 2, Y: 5}, pointType: PointTypeOriginal},
//	}
//	results := traverse(poly1, poly2, BooleanUnion)
//	// results will contain the resulting polygons after the union operation.
//
// Panics:
//   - This function assumes valid input polygons and does not validate them for correctness.
//
// Dependencies:
//   - Relies on helper functions such as findTraversalStartingPoint, togglePolyTraversalDirection,
//     and isInsidePolygon to perform its operations.
//
// Performance:
//   - The function performs a combination of edge traversal and intersection checks, resulting
//     in a complexity proportional to the number of edges and intersection points.
func traverse[T SignedNumber](poly1, poly2 []polyPoint[T], operation BooleanOperation) [][]polyPoint[T] {

	// Step 1: Check for nested polygons if no intersections
	if len(poly1) > 0 && len(poly2) > 0 {
		allPoly2InsidePoly1 := true
		for _, p := range poly2 {
			if !isInsidePolygon(poly1, p.point) {
				allPoly2InsidePoly1 = false
				break
			}
		}

		allPoly1InsidePoly2 := true
		for _, p := range poly1 {
			if !isInsidePolygon(poly2, p.point) {
				allPoly1InsidePoly2 = false
				break
			}
		}

		// Handle containment scenarios
		if allPoly2InsidePoly1 {
			switch operation {
			case BooleanUnion:
				return [][]polyPoint[T]{poly1}
			case BooleanIntersection:
				return [][]polyPoint[T]{poly2}
			case BooleanSubtraction:
				return [][]polyPoint[T]{poly1}
			}
		}

		if allPoly1InsidePoly2 {
			switch operation {
			case BooleanUnion:
				return [][]polyPoint[T]{poly2}
			case BooleanIntersection:
				return [][]polyPoint[T]{poly1}
			case BooleanSubtraction:
				return [][]polyPoint[T]{}
			}
		}
	}

	// Step 2: Normal traversal logic
	direction := polyTraversalDirectionCounterClockwise
	polys := [][]polyPoint[T]{poly1, poly2}
	results := make([][]polyPoint[T], 0)

	for {
		// Find the starting point for traversal
		polyIndex, pointIndex := findTraversalStartingPoint(polys)
		if polyIndex == -1 || pointIndex == -1 {
			break // No unvisited entry points
		}

		// Initialize result path
		result := make([]polyPoint[T], 0, len(poly1)+len(poly2))

		// loop (combined):
		//   - add current point to result & mark as visited
		//   - if direction == CCW { increment point index } else { decrement point index }
		//   - if point is exit:
		//       - swap poly
		//       - if operation is subtraction, reverse direction
		//   - if point matches start point, loop completed

		for {
			// Append current point to result path
			result = append(result, polys[polyIndex][pointIndex])

			// Mark the current point as visited
			polys[polyIndex][pointIndex].visited = true

			// Move to the next point in the current polygon, depending on direction
			if direction == polyTraversalDirectionCounterClockwise {
				pointIndex = (pointIndex + 1) % len(polys[polyIndex])
			} else {
				pointIndex = pointIndex - 1
				if pointIndex == -1 {
					pointIndex = len(polys[polyIndex]) - 1
				}
			}

			// if we've come to an exit point, swap polys & maybe change direction
			if polys[polyIndex][pointIndex].entryExit == intersectionTypeExit {

				// swap poly
				pointIndex = polys[polyIndex][pointIndex].otherPolyIndex
				polyIndex = (polyIndex + 1) % 2

				// swap direction if operation is subtraction
				if operation == BooleanSubtraction {
					direction = togglePolyTraversalDirection(direction)
				}
			}

			// Stop if we loop back to the starting point in the same polygon
			if polys[polyIndex][pointIndex].point.Eq(result[0].point) {
				break
			}
		}

		results = append(results, result)
	}

	// Handle no results for specific operations
	if len(results) == 0 && operation == BooleanUnion {
		return [][]polyPoint[T]{poly1, poly2}
	}

	return results
}

// NewPolygon creates a new polygon with the specified points, type, and optional child polygons.
//
// This function constructs a Polygon object, ensuring the points are correctly oriented based on
// the polygon type and doubling the points to allow integer-safe operations in subsequent calculations.
//
// Parameters:
//
//   - points []Point[T]: A slice of points defining the vertices of the polygon.
//     The points must form a valid polygon with at least three vertices.
//   - polygonType PolygonType: Specifies whether the polygon is solid (PTSolid) or a hole (PTHole).
//     The orientation of the points is adjusted based on this type:
//   - Solid polygons are counterclockwise.
//   - Hole polygons are clockwise.
//   - children ...*Polygon[T]: Optional child polygons (e.g., holes or nested islands) associated
//     with the polygon.
//
// Returns (*Polygon[T], error): A pointer to the newly created Polygon object and an error, if any:
//   - Returns an error if the input polygon has fewer than three points
//     or if the area of the polygon is zero.
//
// Behavior:
//   - The function checks if the input polygon has at least three points and calculates the signed area
//     to ensure the polygon is valid (non-zero area).
//   - Adjusts the orientation of the input points based on the polygon type (solid or hole).
//   - Doubles the coordinates of the points to ensure integer-safe division when casting to the generic type T.
//   - Creates a convex hull of the polygon for optimization purposes in point-in-polygon checks.
//   - Initializes the optional child polygons and prepares the polyPoint representation for the vertices.
//
// Notes:
//   - This function assumes all child polygons have already been created and provided as input.
//   - Points are doubled during initialization to simplify integer division when casting to T.
//     Doubling ensures that any fractional values are dropped without precision loss.
//
// Example:
//
//	points := []Point[float64]{
//	    {X: 0, Y: 0},
//	    {X: 4, Y: 0},
//	    {X: 4, Y: 4},
//	    {X: 0, Y: 4},
//	}
//	polygon, err := NewPolygon(points, PTSolid)
//	if err != nil {
//	    fmt.Println("Error:", err)
//	} else {
//	    fmt.Println("Polygon created:", polygon)
//	}
//
// Errors:
//   - This function does not panic but returns an error for invalid input, such as fewer than three points
//     or a polygon with zero area.
//
// Dependencies:
//   - This function relies on helper functions such as SignedArea2X, EnsureCounterClockwise,
//     EnsureClockwise, and ConvexHull to validate and preprocess the input.
//
// Performance:
//   - The function involves computing the signed area and convex hull, making its complexity proportional
//     to the number of points in the input.
func NewPolygon[T SignedNumber](points []Point[T], polygonType PolygonType, children ...*Polygon[T]) (*Polygon[T], error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("a polygon must have at least three points")
	}

	// Calculate twice the signed area to check area.
	signedArea2X := SignedArea2X(points)
	if signedArea2X == 0 {
		return nil, fmt.Errorf("polygon area must be greater than zero")
	}

	// Ensure the points are in the correct orientation.
	if polygonType == PTSolid {
		points = EnsureCounterClockwise(points)
	} else if polygonType == PTHole {
		points = EnsureClockwise(points)
	}

	// Create new polygon.
	p := new(Polygon[T])
	p.points = make([]polyPoint[T], len(points))
	p.polygonType = polygonType
	p.children = children
	p.maxX = points[0].x

	// Create polyPoints. Double the points.
	for i, point := range points {
		p.points[i] = polyPoint[T]{
			point: Point[T]{
				x: point.x * 2,
				y: point.y * 2,
			},
			pointType:      pointTypeOriginal,
			otherPolyIndex: -1,
		}
		if point.x > p.maxX {
			p.maxX = point.x
		}
	}

	// Create convex hull
	hull := ConvexHull(points...)
	hull = EnsureCounterClockwise(hull)
	p.hull = newSimpleConvexPolygon(hull)

	// Add any self-intersection

	return p, nil
}
