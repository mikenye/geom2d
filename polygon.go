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

type polyDirection int

const (
	polyDirectionCounterClockwise = polyDirection(iota)
	polyDirectionClockwise
)

type polyEdge[T SignedNumber] struct {
	lineSegment LineSegment[T]           // line segment representing the edge
	rel         LineSegmentsRelationship // relationship with ray for isInsidePolygon
}

type polyPoint[T SignedNumber] struct {
	point          Point[T]
	pointType      polyPointType
	entryExit      polyIntersectionType
	visited        bool
	otherPolyIndex int
}

type polyPointType int

const (
	PointTypeNormal                polyPointType = iota // original, unmodified point
	PointTypeIntersection                               // normal point that is also an intersection
	PointTypeAddedIntersection                          // intersection point added, not part of original poly
	PointTypeAddedSelfIntersection                      // intersection point added, not part of original poly
)

type BooleanOperation int

const (
	BooleanUnion BooleanOperation = iota
	BooleanIntersection
	BooleanSubtraction
)

type polyIntersectionType int

const (
	intersectionTypeNotSet polyIntersectionType = iota // Default, not set
	intersectionTypeEntry                              // Indicates this is an entry point
	intersectionTypeExit                               // Indicates this is an exit point
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

// PolygonType (PT) defines the type of a polygon, categorizing it as either a solid region (island)
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
//  2. Filters out points that are not of type PointTypeNormal or PointTypeIntersection.
//  3. Converts the valid points into Point objects with adjusted coordinates (halving x and y values).
//  4. Returns the resulting slice of Point objects.
//
// Returns:
//   - []Point[T]: A slice of Points representing the polygon's outline.
func (p Polygon[T]) Points() []Point[T] {
	points := make([]Point[T], 0, len(p.points))
	for _, point := range p.points {
		if !(point.pointType == PointTypeNormal || point.pointType == PointTypeIntersection) {
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
//  1. **Preliminary Convex Hull Check**: If the polygon has a convex hull defined, the point is first checked
//     against the hull. If the point is outside the hull, it is immediately classified as outside the polygon.
//  2. **Ray-Casting Algorithm**: A horizontal ray is cast to the right from the given point, and intersections
//     with the polygon's edges are counted to determine if the point is inside or outside.
//  3. **Edge and Vertex Checks**: The method explicitly checks for cases where the point lies on a vertex
//     or on an edge of the polygon.
//  4. **Recursive Relationship with Children**: If the polygon has nested children (holes or islands),
//     the method recursively determines the most specific relationship with the lowest-level child polygon.
//
// Parameters:
//
//   - point Point[T]: The point to check the relationship for.
//
// Returns a PointPolygonRelationship: An enumerated value representing the relationship of the point to the polygon:
//   - PPRPointOutside: The point is outside the polygon.
//   - PPRPointInside: The point is inside the polygon.
//   - PPRPointOnVertex: The point coincides with a vertex of the polygon.
//   - PPRPointOnEdge: The point lies on an edge of the polygon.
//   - PPRPointInHole: The point is inside a hole within the polygon.
//   - PPRPointInsideIsland: The point is inside an island nested within the polygon.
//
// Notes:
//   - If the point is outside the convex hull (when present), the relationship is immediately classified as PPRPointOutside.
//   - The ray-casting algorithm ensures robustness by handling edge cases such as collinear points and overlapping segments.
//   - The recursive check ensures that relationships with child polygons (holes and islands) are prioritized.
//
// Example:
//
//	polygon := NewPolygon(points)
//	relationship := polygon.RelationshipToPoint(Point{X: 1, Y: 2})
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

// simpleConvexPolygon represents a convex polygon defined by a set of ordered points.
type simpleConvexPolygon[T SignedNumber] struct {
	Points []Point[T] // Ordered points that form the convex hull
}

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
			if intersects {
				// Check if the intersection point already exists in poly1 or poly2
				if isPointInPolyPoints(intersectionPoint, poly1) || isPointInPolyPoints(intersectionPoint, poly2) {
					continue // Skip adding duplicate intersections
				}

				// Convert the intersection point to a polyPoint
				intersection := polyPoint[T]{
					point:     NewPoint(T(intersectionPoint.x), T(intersectionPoint.y)),
					pointType: PointTypeAddedIntersection, // Mark as intersection
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

func isInsidePolygon[T SignedNumber](polygon []polyPoint[T], point Point[T]) bool {

	crosses := 0

	// cast ray from point to right
	maxX := point.x
	for _, p := range polygon {
		maxX = max(maxX, p.point.x)
	}
	maxX++
	ray := NewLineSegment(point, NewPoint(maxX, point.y))

	// get edges
	edges := edgesFromPolyPoints(polygon)

	// determine relationship with ray for each edge
	for i := range edges {

		// If point is directly on the edge, it's inside
		if point.IsOnLineSegment(edges[i].lineSegment) {
			return true
		}

		// store relationship
		edges[i].rel = ray.RelationshipToLineSegment(edges[i].lineSegment)
	}

	// check for crosses
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

func isPointInPolyPoints[T SignedNumber](point Point[float64], poly []polyPoint[T]) bool {
	for _, p := range poly {
		if float64(p.point.x) == point.x && float64(p.point.y) == point.y {
			return true
		}
	}
	return false
}

// integer values must be multiples of two
func markEntryExitPoints[T SignedNumber](poly1, poly2 []polyPoint[T], operation BooleanOperation) {
	for poly1Point1Index, poly1Point1 := range poly1 {
		poly1Point2Index := (poly1Point1Index + 1) % len(poly1)
		if poly1Point1.pointType == PointTypeAddedIntersection {
			for poly2PointIndex, poly2Point := range poly2 {
				if poly2Point.pointType == PointTypeAddedIntersection && poly1Point1.point.Eq(poly2Point.point) {

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

func newSimpleConvexPolygon[T SignedNumber](points []Point[T]) *simpleConvexPolygon[T] {
	// Assume `points` is already ordered to form a convex polygon
	return &simpleConvexPolygon[T]{Points: points}
}

func togglePolyDirection(direction polyDirection) polyDirection {
	return (direction + 1) % 2
}

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
	direction := polyDirectionCounterClockwise
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
			if direction == polyDirectionCounterClockwise {
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
					direction = togglePolyDirection(direction)
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

// NewPolygon creates a new Polygon instance, ensuring it has at least three points,
// a non-zero area, and the correct point orientation based on the polygon type.
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
			pointType:      PointTypeNormal,
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
