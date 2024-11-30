package geom2d

import (
	"fmt"
	"slices"
)

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

func (c *contour[T]) isContourInside(c2 contour[T]) bool {
	for _, p := range c2 {
		if !c.isPointInside(p.point) {
			return false
		}
	}
	return true
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

// iterEdges iterates over all edges in the contour, calling the `yield` function for each edge.
// Each edge is represented as a LineSegment connecting two consecutive points in the contour.
// The contour is assumed to represent a closed loop, so the last point is connected back to the first.
//
// The `yield` function receives a LineSegment representing an edge and should return `true`
// to continue iteration or `false` to terminate early.
//
// If the contour has fewer than two points, the method returns without calling `yield`.
func (c *contour[T]) iterEdges(yield func(LineSegment[T]) bool) {
	if len(*c) < 2 {
		return // A contour with fewer than two points cannot form edges
	}

	for i := range *c {
		j := (i + 1) % len(*c) // Wrap around to connect the last point to the first
		if !yield(NewLineSegment((*c)[i].point, (*c)[j].point)) {
			return // Exit early if `yield` returns false
		}
	}
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

	// polygonType indicates whether the polygon is a solid region (PTSolid) or a hole (PTHole).
	// This classification is essential for distinguishing between filled and void areas of the polygon.
	polygonType PolygonType

	// siblings contains references to sibling polygons, which must not overlap this polygon, and this are the samw
	// polygonType as this polygon.
	siblings []*PolyTree[T]

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

func (p *PolyTree[T]) BooleanOperation(other *PolyTree[T], operation BooleanOperation) (*PolyTree[T], error) {
	// Edge Case: Check if polygons intersect
	if !p.Intersects(other) {
		switch operation {
		case BooleanUnion:
			if err := p.addSibling(other); err != nil {
				return nil, fmt.Errorf("failed to add sibling: %w", err)
			}
			return p, nil
		case BooleanIntersection:
			return nil, nil // No intersection
		case BooleanSubtraction:
			return p, nil // Original polygon remains unchanged
		default:
			return nil, fmt.Errorf("unknown operation: %v", operation)
		}
	}

	// find intersection points between all polys
	p.findIntersections(other)

	// mark points for Intersection
	p.markEntryExitPoints(other, operation)

	// perform traversal and nest resultant polygons
	return nestPointsToPolyTrees(p.booleanOperationTraversal(other, operation))
}

func (p *PolyTree[T]) Intersects(other *PolyTree[T]) bool {
	// Check if any point of "other" is inside "p"
	for _, otherPoint := range other.contour {
		if p.contour.isPointInside(otherPoint.point) {
			return true
		}
	}

	// Check if any point of "p" is inside "other"
	for _, point := range p.contour {
		if other.contour.isPointInside(point.point) {
			return true
		}
	}

	// Check for edge intersections
	for poly1Edge := range p.contour.iterEdges {
		for poly2Edge := range other.contour.iterEdges {
			if poly1Edge.IntersectsLineSegment(poly2Edge) {
				return true
			}
		}
	}

	// No intersections detected
	return false
}

// returns this poly and all children (including nested children of children etc)
func (p *PolyTree[T]) iterPolys(yield func(*PolyTree[T]) bool) {
	// yield self
	if !yield(p) {
		return
	}

	// yield siblings
	for _, sibling := range p.siblings {
		for s := range sibling.iterPolys {
			if !yield(s) {
				return
			}
		}
	}

	// yield children
	for _, child := range p.children {
		for c := range child.iterPolys {
			if !yield(c) {
				return
			}
		}
	}
}

func (p *PolyTree[T]) resetIntersectionMetadata() {
	// remove intersection data
	for poly := range p.iterPolys {
		for i := 0; i < len(poly.contour); i++ {
			if poly.contour[i].pointType == pointTypeAddedIntersection {
				poly.contour = slices.Delete(poly.contour, i, i+1)
				i--
				continue
			}
			poly.contour[i].pointType = pointTypeOriginal
			poly.contour[i].entryExit = intersectionTypeNotSet
			poly.contour[i].visited = false
			poly.contour[i].intersectionPartner = nil
			poly.contour[i].intersectionPartnerPointIndex = -1
		}
	}
}

func (p *PolyTree[T]) addChild(child *PolyTree[T]) error {
	if p.polygonType == child.polygonType {
		return fmt.Errorf("cannot add child: mismatched polygon types (parent: %v, child: %v)", p.polygonType, child.polygonType)
	}
	child.parent = p
	p.children = append(p.children, child)
	return nil
}

func (p *PolyTree[T]) addSibling(sibling *PolyTree[T]) error {
	if p.polygonType != sibling.polygonType {
		return fmt.Errorf("cannot add sibling as polygonType is mismatched")
	}

	// Add the new sibling to the existing siblings of p
	for _, existingSibling := range p.siblings {
		existingSibling.siblings = append(existingSibling.siblings, sibling)
		sibling.siblings = append(sibling.siblings, existingSibling)
	}

	// Add p to the sibling's sibling list
	sibling.siblings = append(sibling.siblings, p)

	// Add the sibling to p's sibling list
	p.siblings = append(p.siblings, sibling)

	return nil
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
		EnsureCounterClockwise(orderedPoints)

	case PTHole:
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

	// Create convex hull
	hull := ConvexHull(points)
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
	// Iterate through all combinations of polygons:
	poly1i := 0
	for poly1 := range p.iterPolys {

		poly2i := 0

		for poly2 := range other.iterPolys {

			for poly1Point1Index, poly1Point1 := range poly1.contour {
				poly1Point2Index := (poly1Point1Index + 1) % len(poly1.contour)

				if poly1Point1.pointType == pointTypeAddedIntersection {
					for poly2PointIndex, poly2Point := range poly2.contour {
						if poly2Point.pointType == pointTypeAddedIntersection && poly1Point1.point.Eq(poly2Point.point) {

							if poly1.contour[poly1Point1Index].entryExit != intersectionTypeNotSet || poly2.contour[poly2PointIndex].entryExit != intersectionTypeNotSet {
								panic(fmt.Errorf("found intersection metadata when none was expected"))
							}

							// Determine if poly1 traversal is poly1EnteringPoly2 or exiting poly2
							mid := NewLineSegment(
								poly1Point1.point,
								poly1.contour[poly1Point2Index].point).Midpoint()
							midT := NewPoint[T](T(mid.x), T(mid.y))
							poly1EnteringPoly2 := poly2.contour.isPointInside(midT)

							// Debug Logging: Midpoint Information
							if poly1i == 0 {
								fmt.Println("poly1 outer contour")
							} else {
								fmt.Println("poly1 hole contour")
							}
							if poly2i == 0 {
								fmt.Println("poly2 outer contour")
							} else {
								fmt.Println("poly2 hole contour")
							}
							fmt.Printf("Poly1 [%d]: %v -> Poly2 [%d]: %v\n",
								poly1Point1Index, poly1Point1.point, poly2PointIndex, poly2Point.point)
							fmt.Printf("Midpoint: %v, Inside Poly2: %t\n", midT, poly1EnteringPoly2)

							// use lookup table to determine entry/exit points
							poly1.contour[poly1Point1Index].entryExit = entryExitPointLookUpTable[operation][poly1.polygonType][poly2.polygonType][poly1EnteringPoly2].poly1PointType
							poly2.contour[poly2PointIndex].entryExit = entryExitPointLookUpTable[operation][poly1.polygonType][poly2.polygonType][poly1EnteringPoly2].poly2PointType

							// Debug Logging: Marked Entry/Exit
							fmt.Printf("Poly1 EntryExit: %s, Poly2 EntryExit: %s\n",
								poly1.contour[poly1Point1Index].entryExit.String(),
								poly2.contour[poly2PointIndex].entryExit.String())

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

func (p *PolyTree[T]) booleanOperationTraversal(other *PolyTree[T], operation BooleanOperation) [][]Point[T] {
	var direction polyTraversalDirection

	// todo: Step 1: handle edge cases like polygons not intersecting etc.

	// Step 2: Normal traversal logic
	resultContours := make([][]Point[T], 0)

	for {
		// Find the starting point for traversal
		currentPoly, currentPointIndex := p.findTraversalStartingPoint(other)
		if currentPoly == nil || currentPointIndex == -1 {
			fmt.Println("no unvisited entry points")
			break // No unvisited entry points
		}
		fmt.Println("starting with:", currentPoly.contour[currentPointIndex].point)

		// Initialize resultContour path
		resultContour := make([]Point[T], 0, len(p.contour)+len(other.contour))

		// loop (combined):
		//   - add current point to resultContour & mark as visited
		//   - if direction == CCW { increment point index } else { decrement point index }
		//   - if point is exit:
		//       - swap poly
		//       - if operation is subtraction, reverse direction
		//   - if point matches start point, loop completed

		// set initial traversal direction
		//if currentPoly.polygonType == PTSolid {
		direction = polyTraversalForward
		//} else {
		//	direction = polyTraversalReverse
		//}

		for {

			// Append current point to resultContour path
			// Do division here
			resultContour = append(resultContour, NewPoint(currentPoly.contour[currentPointIndex].point.x/2, currentPoly.contour[currentPointIndex].point.y/2))
			fmt.Println("added point:", currentPoly.contour[currentPointIndex].point)

			if currentPoly.contour[currentPointIndex].entryExit == intersectionTypeEntry {
				fmt.Println("point is entry")
			}
			if currentPoly.contour[currentPointIndex].entryExit == intersectionTypeExit {
				fmt.Println("point is exit")
			}

			// Mark the current point as visited
			currentPoly.contour[currentPointIndex].visited = true

			// Also mark the partner point as visited during polygon switches
			if currentPoly.contour[currentPointIndex].intersectionPartner != nil {
				partnerPoly := currentPoly.contour[currentPointIndex].intersectionPartner
				partnerPointIndex := currentPoly.contour[currentPointIndex].intersectionPartnerPointIndex
				partnerPoly.contour[partnerPointIndex].visited = true
			}

			// Move to the next point in the current polygon, depending on direction
			if direction == polyTraversalForward {
				currentPointIndex = (currentPointIndex + 1) % len(currentPoly.contour)
				fmt.Println("incrementing current point index to:", currentPointIndex)
			} else {
				currentPointIndex = currentPointIndex - 1
				if currentPointIndex == -1 {
					currentPointIndex = len(currentPoly.contour) - 1
				}
				fmt.Println("decrementing current point index to:", currentPointIndex)
			}

			// if we've come to an exit point, swap polys & maybe change direction
			if currentPoly.contour[currentPointIndex].entryExit == intersectionTypeExit ||
				(operation == BooleanSubtraction && currentPoly.polygonType == PTHole && currentPoly.contour[currentPointIndex].entryExit == intersectionTypeEntry) ||
				(operation == BooleanSubtraction && currentPoly.polygonType == PTSolid && currentPoly.contour[currentPointIndex].entryExit == intersectionTypeEntry && direction == polyTraversalReverse) {

				fmt.Println("point type exit or entry for hole in subtraction found")

				// Swap polygons
				newCurrentPointIndex := currentPoly.contour[currentPointIndex].intersectionPartnerPointIndex
				currentPoly = currentPoly.contour[currentPointIndex].intersectionPartner
				currentPointIndex = newCurrentPointIndex

				fmt.Println("swapped polys, and changed point index to:", currentPointIndex)

				// Adjust direction if needed
				if operation == BooleanSubtraction {
					direction = togglePolyTraversalDirection(direction)
					fmt.Println("reversed traversal direction due to subtraction")
				}

				if direction == polyTraversalForward {
					fmt.Println("direction is forward")
				}
				if direction == polyTraversalReverse {
					fmt.Println("direction is reverse")
				}
			}

			// Stop if we loop back to the starting point in the same polygon
			if currentPoly.contour[currentPointIndex].point.x/2 == resultContour[0].x && currentPoly.contour[currentPointIndex].point.y/2 == resultContour[0].y {
				fmt.Println("finished loop")
				break
			}
		}

		resultContours = append(resultContours, resultContour)
	}

	// TODO: Handle no resultContours for specific operations
	//if len(resultContours) == 0 && operation == BooleanUnion {
	//	return [][]polyPoint[T]{poly1, poly2}
	//}

	return resultContours
}

func nestPointsToPolyTrees[T SignedNumber](contours [][]Point[T]) (*PolyTree[T], error) {
	// Sanity check: ensure contours exist
	if len(contours) == 0 {
		return nil, fmt.Errorf("no contours provided")
	}

	// Sort polygons by area (ascending order, largest last)
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

	// Create the root PolyTree from the largest polygon
	rootTree, err := NewPolyTree(contours[len(contours)-1], PTSolid)
	if err != nil {
		return nil, fmt.Errorf("failed to create root PolyTree: %w", err)
	}

	// Process the remaining polygons
	for i := len(contours) - 2; i >= 0; i-- {
		polyToNest, err := NewPolyTree(contours[i], PTSolid)
		if err != nil {
			return nil, fmt.Errorf("failed to create PolyTree for contour %d: %w", i, err)
		}

		// Try to find the correct parent polygon
		parent := rootTree.findParentPolygon(polyToNest)
		if parent == nil {
			// If no parent is found, add as a sibling to the root
			if err := rootTree.addSibling(polyToNest); err != nil {
				return nil, fmt.Errorf("failed to add sibling: %w", err)
			}
		} else {
			// Add as a child of the parent
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

	return rootTree, nil
}

func (p *PolyTree[T]) findParentPolygon(polyToNest *PolyTree[T]) *PolyTree[T] {
	// Check if polyToNest is inside the current polygon
	if p.contour.isContourInside(polyToNest.contour) {
		// Check recursively if it fits inside any child
		for _, child := range p.children {
			if nestedParent := child.findParentPolygon(polyToNest); nestedParent != nil {
				return nestedParent
			}
		}
		// If no child contains it, the current polygon is the parent
		return p
	}
	// If not inside the current polygon, return nil
	return nil
}
