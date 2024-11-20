package geom2d

import (
	"fmt"
	"slices"
)

type PointPolygonRelationship int

const (
	PointInHole       PointPolygonRelationship = iota - 1 // point is inside a hole
	PointOutside                                          // point is outside the root polygon
	PointOnVertex                                         // point is on the vertex of a polygon/hole/island
	PointOnEdge                                           // point is on the edge of a polygon/hole/island
	PointInside                                           // point is inside the polygon
	PointInsideIsland                                     // point is inside an island nested within the polygon
)

type PolygonType int

const (
	Solid PolygonType = iota // Solid region ("island")
	Hole                     // Void region ("hole")
)

type Polygon[T SignedNumber] struct {
	points      []polyPoint[T]          // Full outline of the polygon
	polygonType PolygonType             // Whether the polygon is solid or a hole
	children    []*Polygon[T]           // Nested polygons (holes or islands)
	parent      *Polygon[T]             // parent polygon (if any)
	hull        *simpleConvexPolygon[T] // Optional convex hull for fast point-in-polygon checks
	maxX        T                       // maximum X value (for ray casting)
}

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

func (p Polygon[T]) RelationshipToPoint(point Point[T]) PointPolygonRelationship {
	// Use the convex hull as a preliminary check for outside points.
	if p.hull != nil && !p.hull.ContainsPoint(point) {
		return PointOutside
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
			return PointOnVertex
		case LSRAonCD, LSRBonCD, LSRConAB, LSRDonAB, LSRIntersects:
			intersections++
		case LSRCollinearAonCD, LSRCollinearBonCD, LSRCollinearEqual:
			// If the ray and edge are collinear and overlap, the point is on the edge.
			return PointOnEdge
		}
	}

	// Determine if the point is inside or outside based on intersections
	isInside := intersections%2 == 1
	if !isInside {
		return PointOutside
	}

	// Recursive check for the "lowest" relationship with children
	relationship := PointInside // Start with the main polygon relationship

	for _, child := range p.children {
		childRelationship := child.RelationshipToPoint(point)

		switch childRelationship {
		case PointInside, PointOnEdge, PointOnVertex:
			// If the child polygon is a hole, mark the relationship as PointInHole
			// If the child polygon is a solid (island), mark the relationship as PointInsideIsland
			switch child.polygonType {
			case Hole:
				relationship = PointInHole
			case Solid:
				relationship = PointInsideIsland
			}
		case PointInHole, PointInsideIsland:
			// If the child relationship is already the innermost context, return immediately
			return childRelationship
		}
	}

	return relationship
}

type polyPoint[T SignedNumber] struct {
	point          Point[T]
	pointType      polyPointType
	entryExit      EntryExitType
	visited        bool
	otherPolyIndex int
}

type EntryExitType int

const (
	NotSet EntryExitType = iota // Default, not set
	Entry                       // Indicates this is an entry point
	Exit                        // Indicates this is an exit point
)

type polyPointType int

const (
	PointTypeNormal                polyPointType = iota // original, unmodified point
	PointTypeIntersection                               // normal point that is also an intersection
	PointTypeAddedIntersection                          // intersection point added, not part of original poly
	PointTypeAddedSelfIntersection                      // intersection point added, not part of original poly
)

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
		if Orientation(a, b, point) == Clockwise {
			return false // Point is outside
		}
	}
	return true // Point is inside
}

func newSimpleConvexPolygon[T SignedNumber](points []Point[T]) *simpleConvexPolygon[T] {
	// Assume `points` is already ordered to form a convex polygon
	return &simpleConvexPolygon[T]{Points: points}
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
	if polygonType == Solid {
		points = EnsureCounterClockwise(points)
	} else if polygonType == Hole {
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

func isPointInPolyPoints[T SignedNumber](point Point[float64], poly []polyPoint[T]) bool {
	for _, p := range poly {
		if float64(p.point.x) == point.x && float64(p.point.y) == point.y {
			return true
		}
	}
	return false
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

type BooleanOperation int

const (
	BooleanUnion BooleanOperation = iota
	BooleanIntersection
	BooleanSubtraction
)

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
							poly1[poly1Point1Index].entryExit = Exit
							poly2[poly2PointIndex].entryExit = Entry
						} else {
							poly1[poly1Point1Index].entryExit = Entry
							poly2[poly2PointIndex].entryExit = Exit
						}

					case BooleanIntersection:

						if poly1EnteringPoly2 {
							poly1[poly1Point1Index].entryExit = Entry
							poly2[poly2PointIndex].entryExit = Exit
						} else {
							poly1[poly1Point1Index].entryExit = Exit
							poly2[poly2PointIndex].entryExit = Entry
						}

					case BooleanSubtraction:

						if poly1EnteringPoly2 {
							poly1[poly1Point1Index].entryExit = Exit
							poly2[poly2PointIndex].entryExit = Exit
						} else {
							poly1[poly1Point1Index].entryExit = Entry
							poly2[poly2PointIndex].entryExit = Entry
						}

					}
					poly1[poly1Point1Index].otherPolyIndex = poly2PointIndex
					poly2[poly2PointIndex].otherPolyIndex = poly1Point1Index
				}
			}
		}
	}
}

type polyEdge[T SignedNumber] struct {
	lineSegment LineSegment[T]       // line segment representing the edge
	rel         TwoLinesRelationship // relationship with ray for isInsidePolygon
}

type polyDirection int

const (
	polyDirectionCounterClockwise = polyDirection(iota)
	polyDirectionClockwise
)

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

// inOrder returns true if b lies between a and c
func inOrder[T SignedNumber](a, b, c T) bool {
	return (a-b)*(b-c) > 0
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

func isLeftOfEdge[T SignedNumber](point, edgeStart, edgeEnd Point[T]) bool {
	// Compute the cross product of the vector from edgeStart to edgeEnd and
	// the vector from edgeStart to the point.
	// If the cross product is > 0, the point is to the left of the edge.
	return edgeEnd.Sub(edgeStart).CrossProduct(point.Sub(edgeStart)) > 0
}

func findTraversalStartingPoint[T SignedNumber](polys [][]polyPoint[T]) (int, int) {
	for polyIndex := range polys {
		for pointIndex := range polys[polyIndex] {
			if polys[polyIndex][pointIndex].entryExit == Entry && !polys[polyIndex][pointIndex].visited {
				return polyIndex, pointIndex
			}
		}
	}
	// Return -1 if no Entry points are found
	return -1, -1
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
			if polys[polyIndex][pointIndex].entryExit == Exit {

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

func findNextIntersection[T SignedNumber](poly []polyPoint[T], point Point[T]) int {
	for i, p := range poly {
		if p.pointType == PointTypeAddedIntersection && p.point == point {
			return i
		}
	}
	return -1 // No next intersection found
}

// Helper to convert []polyPoint[T] to []Point[T] for area calculation
func getPointsFromPolyPoints[T SignedNumber](poly []polyPoint[T]) []Point[T] {
	points := make([]Point[T], len(poly))
	for i, p := range poly {
		points[i] = p.point
	}
	return points
}
