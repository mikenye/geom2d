package geom2d

import (
	"fmt"
	"slices"
)

type contour[T SignedNumber] []polyTreePoint[T]

func (c *contour[T]) contains(point Point[T]) bool {
	return slices.ContainsFunc(*c, func(p polyTreePoint[T]) bool {
		if p.point.x == point.x && p.point.y == point.y {
			return true
		}
		return false
	})
}

func (c *contour[T]) isPointInside(point Point[T]) bool {
	crosses := 0

	// Cast ray from point to right
	maxX := point.x
	for _, p := range *c {
		maxX = max(maxX, p.point.x)
	}
	maxX++
	ray := NewLineSegment(point, NewPoint(maxX, point.y))

	// Get edges
	edges := c.toEdges()

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

func (c *contour[T]) insertPoint(start, end int, intersection polyTreePoint[T]) {
	segment := NewLineSegment((*c)[start].point, (*c)[end].point)

	// Find the correct position to insert the intersection
	insertPos := end
	for i := start + 1; i < end; i++ {
		existingSegment := NewLineSegment((*c)[start].point, (*c)[i].point)
		if segment.DistanceToPoint(intersection.point) < existingSegment.DistanceToPoint((*c)[i].point) {
			insertPos = i
			break
		}
	}

	// Insert the intersection at the calculated position
	*c = slices.Insert(*c, insertPos, intersection)
}

func (c *contour[T]) toEdges() []polyEdge[T] {
	edges := make([]polyEdge[T], 0, len(*c))
	for i := range *c {
		j := (i + 1) % len(*c)
		edges = append(edges, polyEdge[T]{
			lineSegment: NewLineSegment((*c)[i].point, (*c)[j].point),
		})
	}
	return edges
}

type PolyTree[T SignedNumber] struct {
	// points defines the full outline of the polygon, including all vertices and relevant metadata.
	// Each entry is a polyPoint, which tracks whether the point is a regular vertex, an intersection,
	// or a midpoint between intersections.
	contour contour[T]

	pointIndex int

	traversalDirection polyTraversalDirection

	// polygonType indicates whether the polygon is a solid region (PTSolid) or a hole (PTHole).
	// This classification is essential for distinguishing between filled and void areas of the polygon.
	polygonType PolygonType

	// children contains references to nested polygons, which may represent holes (if this polygon
	// is a solid region) or solid islands (if this polygon is a hole). These hierarchical relationships
	// allow for complex polygon structures.
	children []*PolyTree[T]

	// parent points to the parent polygon of this polygon, if any. For example, a hole's parent
	// would be the solid polygon that contains it.
	parent *PolyTree[T]

	// hull optionally stores the convex hull of the polygon, represented as a simpleConvexPolygon.
	// This can be used to optimize certain operations, such as point-in-polygon checks, by
	// quickly ruling out points that lie outside the convex hull.
	hull simpleConvexPolygon[T]

	// maxX stores the maximum X-coordinate value among the polygon's vertices. This is used
	// for ray-casting operations to determine point-in-polygon relationships.
	maxX T
}

func (p *PolyTree[T]) currentPolyPoint() polyTreePoint[T] {
	return p.contour[p.pointIndex]
}

// returns this poly and all children (including nested children of children etc)
func (p *PolyTree[T]) iterPolys(yield func(*PolyTree[T]) bool) {
	// return the first polygon p
	if !yield(p) {
		return
	}

	// return the children
	for _, child := range p.children {
		for c := range child.iterPolys {
			if !yield(c) {
				return
			}
		}
	}
}

func (p *PolyTree[T]) nextPolyPoint() polyTreePoint[T] {
	p.pointIndex = (p.pointIndex + 1) % len(p.contour)
	return p.currentPolyPoint()
}

func (p *PolyTree[T]) prevPolyPoint() polyTreePoint[T] {
	p.pointIndex--
	if p.pointIndex < 0 {
		p.pointIndex = len(p.contour) - 1
	}
	return p.currentPolyPoint()
}

func (p *PolyTree[T]) resetIntersectionMetadata() {
	// reset pointIndex so no invalid references
	p.pointIndex = 0

	// remove intersection data
	for poly := range p.iterPolys {
		for i := 0; i < len(poly.contour); i++ {
			if poly.contour[i].pointType == pointTypeAddedIntersection {
				poly.contour = slices.Delete(poly.contour, i, i+1)
				i--
			}
			poly.contour[i].pointType = pointTypeOriginal
			poly.contour[i].entryExit = intersectionTypeNotSet
			poly.contour[i].visited = false
			poly.contour[i].intersectionPartner = nil
			poly.contour[i].intersectionPartnerPointIndex = -1
		}
	}
}

func (p *PolyTree[T]) setPointIndex(i int) {
	if i < 0 || i > (len((*p).contour)-1) {
		panic(fmt.Errorf("invalid point index"))
	}
	p.pointIndex = i
}

type polyTreePoint[T SignedNumber] struct {
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

	intersectionPartner           *PolyTree[T]
	intersectionPartnerPointIndex int
}

type NewPolyTreeOption[T SignedNumber] func(*PolyTree[T])

func NewPolyTree[T SignedNumber](points []Point[T], t PolygonType, opts ...NewPolyTreeOption[T]) (*PolyTree[T], error) {

	// sanity check: must be at least three points
	if len(points) < 3 {
		return nil, fmt.Errorf("new polytree must have at least 3 points")
	}

	// sanity check: must have non-zero area
	if SignedArea2X(points) == 0 {
		return nil, fmt.Errorf("new polytree must have non-zero area")
	}

	// create newly allocated zero value Polygon
	p := new(PolyTree[T])

	// set polygon type
	p.polygonType = t

	// ensure points in correct orientation
	orderedPoints := make([]Point[T], len(points))
	copy(orderedPoints, points)
	switch p.polygonType {
	case PTSolid:
		p.traversalDirection = polyTraversalDirectionCounterClockwise
		EnsureCounterClockwise(orderedPoints)

	case PTHole:
		p.traversalDirection = polyTraversalDirectionClockwise
		EnsureClockwise(orderedPoints)
	}

	// set polygon points to points given
	p.maxX = points[0].x * 2
	p.contour = make([]polyTreePoint[T], len(orderedPoints))
	for i, point := range orderedPoints {
		p.contour[i] = polyTreePoint[T]{
			point: NewPoint(point.x*2, point.y*2),
		}
		p.maxX = max(p.maxX, p.contour[i].point.x)
	}
	p.resetIntersectionMetadata()
	p.maxX++
	p.pointIndex = 0

	// Create convex hull
	hull := ConvexHull(points...) // todo: just use the slice?
	EnsureCounterClockwise(hull)
	p.hull = newSimpleConvexPolygon(hull)

	// handle function options
	for _, opt := range opts {
		opt(p)
	}

	// sanity check: all children should be opposite PolygonType
	for _, c := range p.children {
		switch p.polygonType {
		case PTSolid:
			if c.polygonType != PTHole {
				return nil, fmt.Errorf("expected all children to have PolygonType PTHole")
			}
		case PTHole:
			if c.polygonType != PTSolid {
				return nil, fmt.Errorf("expected all children to have PolygonType PTSolid")
			}
		}
	}

	return p, nil
}

func WithChildren[T SignedNumber](children ...*PolyTree[T]) NewPolyTreeOption[T] {
	return func(p *PolyTree[T]) {

		// set children of p
		p.children = children

		// set parent of children to p
		for i := range p.children {
			p.children[i].parent = p
		}
	}
}

func (p *PolyTree[T]) findIntersections(other *PolyTree[T]) {

	// reset intersection metadata
	p.resetIntersectionMetadata()
	other.resetIntersectionMetadata()

	// iterate through all combinations of polys:
	for poly1 := range p.iterPolys {
		for poly2 := range other.iterPolys {

			// Iterate through each edge in poly1
			for i1 := 0; i1 < len(poly1.contour); i1++ {
				j1 := (i1 + 1) % len(poly1.contour) // Wrap around to form a closed polygon
				segment1 := NewLineSegment(poly1.contour[i1].point, poly1.contour[j1].point)

				// Iterate through each edge in poly2
				for i2 := 0; i2 < len(poly2.contour); i2++ {
					j2 := (i2 + 1) % len(poly2.contour)
					segment2 := NewLineSegment(poly2.contour[i2].point, poly2.contour[j2].point)

					// Check for intersection between the segments
					intersectionPoint, intersects := segment1.IntersectionPoint(segment2)
					intersectionPointT := NewPoint(T(intersectionPoint.x), T(intersectionPoint.y))
					if intersects {
						// Check if the intersection point already exists in poly1 or poly2
						if poly1.contour.contains(intersectionPointT) || poly2.contour.contains(intersectionPointT) {
							continue // Skip adding duplicate intersections
						}

						// Convert the intersection point to a polyPoint
						intersection := polyTreePoint[T]{
							point:                         NewPoint(T(intersectionPoint.x), T(intersectionPoint.y)),
							pointType:                     pointTypeAddedIntersection, // Mark as intersection
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						}

						// Insert intersection into both polygons
						poly1.contour.insertPoint(i1, j1, intersection)
						poly2.contour.insertPoint(i2, j2, intersection)

						// Increment indices to avoid re-processing this intersection
						i1++
						i2++
					}
				}
			}
		}
	}
}

func (p *PolyTree[T]) findTraversalStartingPoint(other *PolyTree[T]) (*PolyTree[T], int) {
	for _, ptOuter := range []*PolyTree[T]{p, other} {
		for polyTree := range ptOuter.iterPolys {
			for pointIndex := range polyTree.contour {
				if polyTree.contour[pointIndex].entryExit == intersectionTypeEntry && !polyTree.contour[pointIndex].visited {
					return polyTree, pointIndex
				}
			}
		}
	}

	// Return -1 if no Entry points are found
	return nil, -1
}

// PolyTree.findIntersections must be run prior!
func (p *PolyTree[T]) markEntryExitPoints(other *PolyTree[T], operation BooleanOperation) {
	// iterate through all combinations of polys:
	for poly1 := range p.iterPolys {
		for poly2 := range other.iterPolys {
			for poly1Point1Index, poly1Point1 := range poly1.contour {
				poly1Point2Index := (poly1Point1Index + 1) % len(poly1.contour)

				if poly1Point1.pointType == pointTypeAddedIntersection {
					for poly2PointIndex, poly2Point := range poly2.contour {
						if poly2Point.pointType == pointTypeAddedIntersection && poly1Point1.point.Eq(poly2Point.point) {

							if poly1.contour[poly1Point1Index].entryExit != intersectionTypeNotSet || poly2.contour[poly2PointIndex].entryExit != intersectionTypeNotSet {
								panic(fmt.Errorf("found intersection metadata when none was expected"))
							}

							// determine if poly1 traversal is poly1EnteringPoly2 or existing poly2
							mid := NewLineSegment(
								poly1Point1.point,
								poly1.contour[poly1Point2Index].point).Midpoint()
							midT := NewPoint[T](T(mid.x), T(mid.y))
							poly1EnteringPoly2 := poly2.contour.isPointInside(midT)

							switch operation {
							case BooleanUnion:

								if poly1EnteringPoly2 {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeExit
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeEntry
								} else {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeEntry
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeExit
								}

							case BooleanIntersection:

								if poly1EnteringPoly2 {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeEntry
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeExit
								} else {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeExit
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeEntry
								}

							case BooleanSubtraction:

								if poly1EnteringPoly2 {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeExit
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeExit
								} else {
									poly1.contour[poly1Point1Index].entryExit = intersectionTypeEntry
									poly2.contour[poly2PointIndex].entryExit = intersectionTypeEntry
								}

							}

							poly1.contour[poly1Point1Index].intersectionPartner = poly2
							poly1.contour[poly1Point1Index].intersectionPartnerPointIndex = poly2PointIndex

							poly2.contour[poly2PointIndex].intersectionPartner = poly1
							poly2.contour[poly2PointIndex].intersectionPartnerPointIndex = poly1Point1Index
						}
					}
				}
			}
		}
	}
}

func (p *PolyTree[T]) traverse(other *PolyTree[T], operation BooleanOperation) [][]polyTreePoint[T] {

	// todo: Step 1: handle edge cases like polygons not intersecting etc.

	// Step 2: Normal traversal logic
	direction := polyTraversalDirectionCounterClockwise
	results := make([][]polyTreePoint[T], 0)

	for {
		// Find the starting point for traversal
		currentPoly, currentPointIndex := p.findTraversalStartingPoint(other)
		if currentPoly == nil || currentPointIndex == -1 {
			fmt.Println("no unvisited entry points")
			break // No unvisited entry points
		}
		fmt.Println("starting with:", currentPoly, currentPointIndex)

		// Initialize result path
		result := make([]polyTreePoint[T], 0, len(p.contour)+len(other.contour))

		// loop (combined):
		//   - add current point to result & mark as visited
		//   - if direction == CCW { increment point index } else { decrement point index }
		//   - if point is exit:
		//       - swap poly
		//       - if operation is subtraction, reverse direction
		//   - if point matches start point, loop completed

		for {
			// Append current point to result path
			result = append(result, currentPoly.contour[currentPointIndex])
			fmt.Println("added point:", currentPoly.contour[currentPointIndex])

			// Mark the current point as visited
			currentPoly.contour[currentPointIndex].visited = true

			// Move to the next point in the current polygon, depending on direction
			if direction == polyTraversalDirectionCounterClockwise {

				currentPointIndex = (currentPointIndex + 1) % len(currentPoly.contour)
			} else {
				currentPointIndex = currentPointIndex - 1
				if currentPointIndex == -1 {
					currentPointIndex = len(currentPoly.contour) - 1
				}
			}

			// if we've come to an exit point, swap polys & maybe change direction
			if currentPoly.contour[currentPointIndex].entryExit == intersectionTypeExit {

				// swap poly
				currentPointIndex = currentPoly.contour[currentPointIndex].intersectionPartnerPointIndex
				currentPoly = currentPoly.contour[currentPointIndex].intersectionPartner

				fmt.Println("swapped polys")

				// swap direction if operation is subtraction
				if operation == BooleanSubtraction {
					direction = togglePolyTraversalDirection(direction)
				}
			}

			// Stop if we loop back to the starting point in the same polygon
			if currentPoly.contour[currentPointIndex].point.Eq(result[0].point) {
				fmt.Println("finished loop")
				break
			}
		}

		results = append(results, result)
	}

	// TODO: Handle no results for specific operations
	//if len(results) == 0 && operation == BooleanUnion {
	//	return [][]polyPoint[T]{poly1, poly2}
	//}

	return results
}
