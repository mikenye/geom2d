package linesegment

import (
	"fmt"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/gotrees/rbtree"
	"log"
	"strings"
)

// eventData stores information about segments associated with an event point.
//
// Each event point in the sweep line algorithm can be associated with:
// - Segments whose upper endpoint is the event point (entering the sweep line)
// - Segments whose lower endpoint is the event point (leaving the sweep line)
// - Segments that contain the event point as an interior point
type eventData struct {
	UofP []LineSegment // Segments whose upper endpoint is at this event point
	LofP []LineSegment // Segments whose lower endpoint is at this event point
	CofP []LineSegment // Segments that have the event point on their interior
}

// addEventPoint adds an event point to the event queue with its associated segments.
//
// If the event point already exists in the queue, it merges the new segments with the
// existing ones. Otherwise, it creates a new event with the provided segments.
//
// Parameters:
//   - Q: The event queue tree
//   - p: The event point to add
//   - UofP: Segments whose upper endpoint is p (entering the sweep line)
//   - LofP: Segments whose lower endpoint is p (leaving the sweep line)
//
// Note: In the sweep line algorithm, a point can be both an upper endpoint of some segments
// and a lower endpoint of others, in which case both UofP and LofP will contain segments.
func addEventPoint(
	Q *rbtree.Tree[point.Point, eventData],
	p point.Point,
	UofP []LineSegment,
	LofP []LineSegment,
	CofP []LineSegment,
) {
	// Check if the event point already exists in the queue
	existingP, found := Q.Search(p)
	if found {
		// If the event point exists, merge the new segments with existing ones
		ed := Q.Value(existingP)
		ed.UofP = append(ed.UofP, UofP...)
		ed.LofP = append(ed.LofP, LofP...)
		ed.CofP = append(ed.CofP, CofP...)
		Q.Insert(p, ed)
	} else {
		// Otherwise, add a new event with the provided segments
		ed := eventData{
			UofP: UofP,
			LofP: LofP,
			CofP: CofP,
		}
		Q.Insert(p, ed)
	}
}

// lexicographicOrder is a comparison function that orders points for the event queue.
//
// In the sweep line algorithm, events are processed from top to bottom. This function
// orders points first by decreasing Y-coordinate (higher points come first), and if the
// Y-coordinates are equal, then by increasing X-coordinate (leftmost points come first).
//
// Parameters:
//   - a, b: The two points to compare
//
// Returns:
//   - true if point a should come before point b in the ordering, false otherwise
func lexicographicOrder(a, b point.Point) bool {
	debugLog := strings.Builder{}
	defer func() {
		log.Println(debugLog.String())
	}()

	debugLog.WriteString(fmt.Sprintf("event queue order: a: %s, b: %s,  ", a, b))

	// If Y-coordinates are equal (within epsilon)
	if numeric.FloatEquals(a.Y(), b.Y(), geom2d.GetEpsilon()) {
		// Order by X-coordinate (leftmost first)
		debugLog.WriteString(fmt.Sprintf("y equal, order by x: %t", a.X() < b.X()))
		return numeric.FloatLessThan(a.X(), b.X(), geom2d.GetEpsilon())
	}

	// Otherwise, order by Y-coordinate (highest first)
	debugLog.WriteString(fmt.Sprintf("order by y: %t", a.Y() > b.Y()))
	return numeric.FloatGreaterThan(a.Y(), b.Y(), geom2d.GetEpsilon())
}

// statusStructureOrder is a comparison function that orders line segments in the status structure.
//
// The status structure must maintain the left-to-right ordering of segments as they
// intersect with the current sweep line position. This function orders segments by:
// 1. X-coordinate at the sweep line's Y position
// 2. Orientation of points if X-coordinates are equal
// 3. Various tie-breaking rules for special cases
//
// Parameters:
//   - p: Pointer to the current sweep line position
//   - a, b: The two line segments to compare
//
// Returns:
//   - true if segment a should come before segment b in the ordering, false otherwise
func statusStructureOrder(p *point.Point, a, b LineSegment) bool {
	debugLog := strings.Builder{}
	defer func() {
		log.Println(debugLog.String())
	}()

	// Get slopes of both segments
	aSlope := a.Slope()
	bSlope := b.Slope()

	// Check if segments are horizontal (special case)
	aIsHorizontal := aSlope == 0
	bIsHorizontal := bSlope == 0

	// Calculate X-coordinates where each segment intersects the sweep line
	var aX, bX float64
	if aIsHorizontal {
		aX = p.X() // For horizontal segments at sweep line, use sweep line's X
	} else {
		aX = a.XAtY(p.Y()) // Otherwise calculate intersection X at sweep line's Y
	}

	if bIsHorizontal {
		bX = p.X()
	} else {
		bX = b.XAtY(p.Y())
	}

	debugLog.WriteString(fmt.Sprintf("p: %s, a: %s, b: %s, aX: %f, bX: %f, sA: %f, sB: %f, ", p, a, b, aX, bX, aSlope, bSlope))

	// Primary ordering: by X-coordinate at sweep line position
	if !numeric.FloatEquals(aX, bX, geom2d.GetEpsilon()) {
		debugLog.WriteString(fmt.Sprintf("order by x: %t", aX < bX))
		return aX < bX
	}

	// If X-coordinates are equal, use orientation to break ties
	o := point.Orientation(*p, a.lower, b.lower)
	if o == point.Counterclockwise {
		debugLog.WriteString("p -> a.lower -> b.lower: counterclockwise, a is less than b")
		return true
	}
	if o == point.Clockwise {
		debugLog.WriteString("p -> a.lower -> b.lower: clockwise, a is not less than b")
		return false
	}

	// Handle cases where segments share endpoints or are collinear
	if p.Eq(a.lower) {
		debugLog.WriteString("p == a.lower: a is not less than b")
		return false
	}
	if p.Eq(b.lower) {
		debugLog.WriteString("p == b.lower: a less than b")
		return true
	}
	if p.Eq(a.lower) && p.Eq(b.lower) {
		// If both segments start at p, order by slope (steeper first)
		debugLog.WriteString(fmt.Sprintf("p == a.lower == b.lower: order by slope above the line: %t", aSlope < bSlope))
		return aSlope < bSlope
	}

	// Additional tie-breaking rules based on endpoints
	if !numeric.FloatEquals(a.lower.Y(), b.lower.Y(), geom2d.GetEpsilon()) {
		// Order by lower endpoint Y (higher first)
		debugLog.WriteString(fmt.Sprintf("order by lower y: %t", a.lower.Y() > b.lower.Y()))
		return a.lower.Y() > b.lower.Y()
	}

	if !numeric.FloatEquals(a.lower.X(), b.lower.X(), geom2d.GetEpsilon()) {
		// Order by lower endpoint X (leftmost first)
		debugLog.WriteString(fmt.Sprintf("order by lower x: %t", a.lower.X() < b.lower.X()))
		return a.lower.X() < b.lower.X()
	}

	if !numeric.FloatEquals(a.upper.Y(), b.upper.Y(), geom2d.GetEpsilon()) {
		// Order by upper endpoint Y (higher first)
		debugLog.WriteString(fmt.Sprintf("order by upper y: %t", a.upper.Y() > b.upper.Y()))
		return a.upper.Y() > b.upper.Y()
	}

	// Final tie-break by upper endpoint X
	return a.upper.X() < b.upper.X()
}

// FindIntersections uses the Bentley-Ottmann sweep line algorithm to efficiently find
// all intersection points among a set of line segments.
//
// The sweep line algorithm works by moving an imaginary horizontal line downward across
// the plane, stopping at event points (segment endpoints and intersections). At each stop,
// it updates a data structure that keeps track of which segments intersect the sweep line
// and checks for new intersections.
//
// Parameters:
//   - S: A slice of LineSegment objects to find intersections among
//
// Returns:
//   - A red-black tree mapping intersection points to the set of segments involved in each intersection
//
// Algorithm Overview:
// 1. Initialize data structures:
//   - An event queue Q ordered by Y-coordinate (highest first), then X-coordinate
//   - A status structure T tracking segments that intersect the current sweep line
//   - A results structure R storing intersections
//
// 2. Process events in order, handling three types:
//   - Upper endpoint events: Add segments to status structure
//   - Lower endpoint events: Remove segments from status structure
//   - Intersection events: Record intersection and update segment order
//
// 3. At each event, check for new intersections between segments that become adjacent
//
// Time Complexity: O((n+k)log(n)) where n is the number of segments and k is the number of intersections
// Space Complexity: O(n+k)
func FindIntersections(S []LineSegment) *rbtree.Tree[point.Point, map[LineSegment]struct{}] {
	// Pointer to the current sweep line position, used by the status structure ordering function
	var sweepLinePoint *point.Point

	// Sanitize input segments (remove duplicates and degenerate segments)
	S = sanitizeInput(S)

	// Initialize the event queue Q with lexicographic ordering (highest Y first, then smallest X)
	Q := rbtree.New[point.Point, eventData](lexicographicOrder)

	// Initialize the status structure T that tracks segments intersecting the sweep line
	// The comparison function orders segments by their X-coordinate at the sweep line position
	T := rbtree.New[LineSegment, struct{}](func(a, b LineSegment) bool {
		return statusStructureOrder(sweepLinePoint, a, b)
	})

	// Initialize the results structure to store intersection points
	R := rbtree.New[point.Point, map[LineSegment]struct{}](lexicographicOrder)

	// Insert all segment endpoints into the event queue
	for _, l := range S {
		// For upper endpoints, store the associated segment
		addEventPoint(Q, l.upper, []LineSegment{l}, nil, nil)
		// For lower endpoints, mark as a segment that will leave the sweep line
		addEventPoint(Q, l.lower, nil, []LineSegment{l}, nil)
	}

	// Debug counter
	i := 0

	// Main sweep line loop
	for {
		// Debug logging
		i++
		log.Printf("\n\n======== ITERATION %d ========\n", i)
		log.Printf("Q:\n%s", Q.String())
		log.Printf("T:\n%s", T.String())

		// Get the next event point (highest Y-coordinate, then leftmost X-coordinate)
		event := Q.Min(Q.Root())

		// If no more events, we're done
		if Q.IsNil(event) {
			break
		}

		// Remove the current event from the queue
		Q.Delete(event)

		// Extract the event point and associated segments
		eventPoint := Q.Key(event)
		// Initialize sweep line position if this is the first event
		if sweepLinePoint == nil {
			sweepLinePoint = &eventPoint
		}
		ed := Q.Value(event)

		log.Printf("popped event: %s\n", event.String())

		// Process the event point (handle segment updates and check for new intersections)
		handleEventPoint(Q, T, R, eventPoint, ed.UofP, ed.LofP, ed.CofP, sweepLinePoint)
	}

	// Debug output of final results
	log.Println("======== RESULTS ========")
	log.Printf("results: %v\n", R)

	// At the end of the algorithm, the status structure should be empty
	// (all segments should have been processed)
	if T.Size() != 0 {
		panic("status not empty at end of FindIntersections")
	}

	return R
}

// handleEventPoint processes an event point in the sweep line algorithm.
//
// An event point occurs when the sweep line encounters:
// - The upper endpoint of a segment (segment enters the sweep line)
// - The lower endpoint of a segment (segment leaves the sweep line)
// - An intersection between segments
//
// Parameters:
//   - Q: The event queue tree
//   - T: The status structure tree containing segments currently intersecting the sweep line
//   - R: The results tree storing discovered intersections
//   - p: The current event point being processed
//   - UofP: Segments whose upper endpoint is p (entering the sweep line)
//   - LofP: Segments whose lower endpoint is p (leaving the sweep line)
//   - sweepLinePoint: Pointer to the current position of the sweep line
//
// The algorithm follows these steps:
// 1. Find segments in the status structure that contain point p (not endpoints)
// 2. Report p as an intersection if it's a true intersection (more than 1 segment)
// 3. Remove segments that are leaving the sweep line
// 4. Add/reinsert segments that are entering/continuing on the sweep line
// 5. Check for new intersections between segments that have become adjacent
func handleEventPoint(
	Q *rbtree.Tree[point.Point, eventData],
	T *rbtree.Tree[LineSegment, struct{}],
	R *rbtree.Tree[point.Point, map[LineSegment]struct{}],
	p point.Point,
	UofP []LineSegment,
	LofP []LineSegment,
	CofP []LineSegment,
	sweepLinePoint *point.Point,
) {
	// Find all segments in the status structure that contain point p (interior points, not endpoints)
	log.Printf("Status structure sweep line point at:%s\n", sweepLinePoint)
	CofP = append(CofP, findSegmentsContainingPoint(T, p)...)

	// Debug logging
	log.Printf("UofP: %v\n", UofP)
	log.Printf("LofP: %v\n", LofP)
	log.Printf("CofP: %v\n", CofP)

	// If point p is incident to more than one segment, it's an intersection
	if len(LofP)+len(UofP)+len(CofP) > 1 {
		log.Println("LofP ∪ UofP ∪ CofP contains more than one segment, thus is an intersection")

		// Report p as an intersection with all segments involved
		segs := append(LofP, append(UofP, CofP...)...)
		addResult(R, p, segs)

		log.Printf("intersection result: %v %v\n", p, segs)
	}

	// Remove segments that are leaving the sweep line (lower endpoints and containing segments)
	log.Println("Delete the segments in LofP ∪ CofP from T")

	segmentsToRemove := append(LofP, CofP...)
	for _, l := range segmentsToRemove {
		node, found := T.Search(l)
		if !found {
			panic(fmt.Errorf("could not find node to delete: %v\n", l))
		}
		log.Printf("deleting: %v\n", node)
		T.Delete(node)
	}
	log.Printf("T after deletion:\n%s", T.String())

	// Update sweep line position to the current event point
	*sweepLinePoint = p
	log.Printf("Status structure sweep line point at:%s\n", sweepLinePoint)

	// Insert segments that are entering or continuing through the sweep line
	log.Println("Insert the segments in UofP ∪ CofP into T")

	segmentsToInsert := append(UofP, CofP...)
	for _, l := range segmentsToInsert {
		node, _ := T.Insert(l, struct{}{})
		log.Printf("inserted: %v\n", node)
	}
	log.Printf("T after insertion:\n%s", T.String())

	// Check for new potential intersections based on the updated status structure
	log.Printf("UofP ∪ CofP: %d\n", len(UofP)+len(CofP))

	if len(UofP)+len(CofP) == 0 {
		// When a segment is completely removed, check if its left and right neighbors
		// in the status structure now intersect
		log.Println("Let sl and sr be the left and right neighbors of p in T")
		sl, sr := findNeighborsOfPoint(T, p)
		log.Printf("sl: %v\n", sl)
		log.Printf("sr: %v\n", sr)

		if sl != nil && sr != nil {
			findNewEvent(Q, R, *sl, *sr, p)
		}

		// To handle collinear segments, also check the LofP segments against their old neighbors
		for _, seg := range LofP {
			log.Println("collinear segment handling")
			leftNeighbor, leftNeighborFound := T.Ceiling(seg)
			if leftNeighborFound && !T.IsNil(leftNeighbor) {
				log.Printf("left neighbor of %s: %s", seg, leftNeighbor)
				findNewEvent(Q, R, T.Key(leftNeighbor), seg, p)
			}
			rightNeighbor, rightNeighborFound := T.Ceiling(seg)
			if rightNeighborFound && !T.IsNil(rightNeighbor) {
				log.Printf("right neighbor of %s: %s", seg, rightNeighbor)
				findNewEvent(Q, R, T.Key(rightNeighbor), seg, p)
			}
		}

	} else {
		// When segments are added or updated, find the leftmost and rightmost segments
		// from those updated and check for intersections with their new neighbors
		log.Println("find the leftmost & rightmost segments of UofP ∪ CofP")

		// Create temporary tree to order the segments at the current point
		tempT := rbtree.New[LineSegment, struct{}](func(a, b LineSegment) bool {
			return statusStructureOrder(&p, a, b)
		})

		// Insert all updated segments into temporary tree
		activeSegments := append(UofP, CofP...)
		for _, l := range activeSegments {
			tempT.Insert(l, struct{}{})
		}

		// Find leftmost and rightmost segments
		leftmost := tempT.Min(tempT.Root())
		rightmost := tempT.Max(tempT.Root())
		log.Printf("leftmost:  %s\n", leftmost)
		log.Printf("rightmost: %s\n", rightmost)

		// Check for intersection between leftmost segment and its left neighbor
		sPrime, sPrimeFound := T.Search(tempT.Key(leftmost))
		if sPrimeFound {
			log.Printf("sPrime: %s\n", sPrime)
			sl := T.Predecessor(sPrime)
			if !T.IsNil(sl) {
				log.Printf("sl: %v\n", sl)
				findNewEvent(Q, R, T.Key(sl), T.Key(sPrime), p)
			}
		}

		// Check for intersection between rightmost segment and its right neighbor
		sDoublePrime, sDoublePrimeFound := T.Search(tempT.Key(rightmost))
		if sDoublePrimeFound {
			log.Printf("sDoublePrime: %s\n", sDoublePrime)
			sr := T.Successor(sDoublePrime)
			if !T.IsNil(sr) {
				log.Printf("sr: %v\n", sr)
				// BUG FIX: We should check sDoublePrime against sr, not sPrime against sr
				findNewEvent(Q, R, T.Key(sDoublePrime), T.Key(sr), p)
			}
		}
	}
}

// findSegmentsContainingPoint locates all line segments in the status structure T
// that contain the point p at the current sweep line position.
//
// Using the Floor and Ceiling methods on the rbtree, this function creates a degenerate
// line segment at the point p and then searches for segments containing p by traversing
// the tree in both directions from that point.
//
// Parameters:
//   - T: The status structure, a red-black tree containing line segments ordered by statusStructureOrder
//   - p: The point to test for containment
//
// Returns:
//   - []LineSegment: All line segments from T that contain the point p
func findSegmentsContainingPoint(
	T *rbtree.Tree[LineSegment, struct{}],
	p point.Point,
) []LineSegment {
	// If the tree is empty, return an empty result
	if T.IsNil(T.Root()) {
		return nil
	}

	// Create a degenerate line segment at point p to use as a search key
	degenerateSegment := NewFromPoints(p, p)

	// Collection of segments containing point p
	CofP := make([]LineSegment, 0)

	// Find a segment at or near the point's position in the ordering
	floorNode, floorFound := T.Floor(degenerateSegment)
	if floorFound {
		// Start with the floor node and check if it contains p
		currentNode := floorNode
		for !T.IsNil(currentNode) {
			seg := T.Key(currentNode)
			if seg.ContainsPoint(p) {
				// Only check if p is an endpoint of THIS segment
				if !seg.upper.Eq(p) && !seg.lower.Eq(p) {
					CofP = append(CofP, seg)
				}
				// Continue checking predecessors
				currentNode = T.Predecessor(currentNode)
			} else {
				// If a segment doesn't contain p, no earlier segments will either
				break
			}
		}
	}

	// Also check ceiling node and successors
	ceilingNode, ceilingFound := T.Ceiling(degenerateSegment)
	if ceilingFound && (floorNode != ceilingNode || !floorFound) {
		// Skip the floor node if it was already processed
		currentNode := ceilingNode
		for !T.IsNil(currentNode) {
			seg := T.Key(currentNode)
			if seg.ContainsPoint(p) {
				// Only check if p is an endpoint of THIS segment
				if !seg.upper.Eq(p) && !seg.lower.Eq(p) {
					CofP = append(CofP, seg)
				}
				// Continue checking successors
				currentNode = T.Successor(currentNode)
			} else {
				// If a segment doesn't contain p, no later segments will either
				break
			}
		}
	}

	return CofP
}

// findNeighborsOfPoint locates the left and right neighbors of a point p in the status structure T.
// In the context of the sweep line algorithm, these are the segments immediately to the left and right
// of point p in the ordering defined by statusStructureOrder.
//
// The function works by:
// 1. Creating a degenerate line segment at point p
// 2. Using Floor to find the largest segment less than or equal to the degenerate segment
// 3. Using Ceiling to find the smallest segment greater than or equal to the degenerate segment
// 4. If either neighbor contains point p, continuing to search until finding segments that don't
//
// This approach elegantly handles cases where point p lies on one or more line segments in the
// status structure, ensuring we find the true left and right neighbors that don't contain p.
//
// Parameters:
//   - T: The status structure, a red-black tree containing line segments ordered by statusStructureOrder
//   - p: The point for which to find neighbors
//
// Returns:
//   - sl: Pointer to the left neighbor segment (nil if none exists)
//   - sr: Pointer to the right neighbor segment (nil if none exists)
//
// findNeighborsOfPoint locates the left and right neighbors of a point p in the status structure T.
func findNeighborsOfPoint(
	T *rbtree.Tree[LineSegment, struct{}],
	p point.Point,
) (sl, sr *LineSegment) {
	// Create a degenerate line segment at point p
	degenerateSegment := NewFromPoints(p, p)

	// Find the floor (greatest segment less than or equal to the degenerate segment)
	floorNode, floorFound := T.Floor(degenerateSegment)

	// Keep looking for a predecessor until we find one that doesn't contain p
	for floorFound && !T.IsNil(floorNode) {
		floorSeg := T.Key(floorNode)
		if !floorSeg.ContainsPoint(p) {
			sl = &floorSeg
			break
		}
		floorNode = T.Predecessor(floorNode)
	}

	// Find the ceiling (smallest segment greater than or equal to the degenerate segment)
	ceilingNode, ceilingFound := T.Ceiling(degenerateSegment)

	// Keep looking for a successor until we find one that doesn't contain p
	for ceilingFound && !T.IsNil(ceilingNode) {
		ceilingSeg := T.Key(ceilingNode)
		if !ceilingSeg.ContainsPoint(p) {
			sr = &ceilingSeg
			break
		}
		ceilingNode = T.Successor(ceilingNode)
	}

	return sl, sr
}

// findNewEvent checks for intersection between two line segments sl and sr,
// adding any valid intersection points to both the event queue Q and the results R.
//
// This function is called when the status structure changes and there's a potential
// for new intersections between segments that have become adjacent in the status structure.
//
// Parameters:
//   - Q: The event queue tree storing upcoming event points
//   - R: The results tree storing discovered intersections
//   - sl: The left segment to check for intersection
//   - sr: The right segment to check for intersection
//   - p: The current event point being processed (i.e., the sweep line position)
//
// The algorithm:
// 1. Calculates intersection points between sl and sr
// 2. Filters intersections to only include those below the sweep line or on the sweep line to the right of p
// 3. Adds valid intersections to both the event queue and the results collection
// 4. For collinear segments, ensures that all overlapping endpoints are properly reported
// findNewEvent checks for intersection between two line segments sl and sr,
// adding any valid intersection points to both the event queue Q and the results R.
//
// This function is called when the status structure changes and there's a potential
// for new intersections between segments that have become adjacent in the status structure.
//
// Parameters:
//   - Q: The event queue tree storing upcoming event points
//   - R: The results tree storing discovered intersections
//   - sl: The left segment to check for intersection
//   - sr: The right segment to check for intersection
//   - p: The current event point being processed (i.e., the sweep line position)
func findNewEvent(
	Q *rbtree.Tree[point.Point, eventData],
	R *rbtree.Tree[point.Point, map[LineSegment]struct{}],
	sl, sr LineSegment,
	p point.Point,
) {
	// Find intersection points between segments
	intersectionPoints, intersects := sl.IntersectionPoints(sr)
	if !intersects {
		return
	}

	// Handle collinear segments special case
	// When segments are collinear, IntersectionPoints returns the endpoints of the overlap
	isCollinear := len(intersectionPoints) > 1

	// Add all intersection points to both the event queue and results
	for _, intersectionPoint := range intersectionPoints {

		// Also add to results (with both segments associated with this intersection point)
		addResult(R, intersectionPoint, []LineSegment{sl, sr})

		// For regular intersections (or the first point of collinear segments),
		// filter to only include those that:
		// - Are below the current sweep line (i.e., have lower y-coordinate), or
		// - Are on the sweep line but to the right of the current event point
		if numeric.FloatLessThan(intersectionPoint.Y(), p.Y(), geom2d.GetEpsilon()) ||
			(numeric.FloatEquals(intersectionPoint.Y(), p.Y(), geom2d.GetEpsilon()) &&
				numeric.FloatGreaterThan(intersectionPoint.X(), p.X(), geom2d.GetEpsilon())) {

			if isCollinear {
				// Add to event queue to ensure we process this intersection point
				addEventPoint(Q, intersectionPoint, nil, nil, []LineSegment{sl, sr})
			} else {
				// Add to event queue to ensure we process this intersection point
				addEventPoint(Q, intersectionPoint, nil, nil, nil)
			}

			// Log for debugging
			log.Printf("findNewEvent added event point: %s\n", intersectionPoint)
		}
	}
}

// FindIntersectionsBruteForce finds all intersection points among a set of line segments using a naive approach.
//
// This function implements a simple brute-force algorithm that compares every pair of line segments
// to check for intersections. While less efficient than sweep-line for large datasets, it provides
// a reliable reference implementation with O(n²) time complexity.
//
// Parameters:
//   - S: A slice of LineSegment objects to find intersections among
//
// Returns:
//   - A red-black tree mapping intersection points to the set of segments involved in each intersection
//
// Behavior:
//   - First sanitizes the input to remove duplicate segments
//   - Iterates through all pairs of segments, checking each pair for intersection
//   - Records all intersection points and their associated segments in the result tree
//   - Uses the same result format as FindIntersections for consistency
//
// Note: While less efficient than the sweep line algorithm, this method is useful for:
//   - Testing correctness of the more complex algorithm
//   - Small datasets where simplicity is preferred over performance
func FindIntersectionsBruteForce(S []LineSegment) *rbtree.Tree[point.Point, map[LineSegment]struct{}] {
	S = sanitizeInput(S)
	R := rbtree.New[point.Point, map[LineSegment]struct{}](lexicographicOrder)
	for i, segA := range S {
		for j, segB := range S {
			if j <= i {
				continue
			}

			intersections, _ := segA.IntersectionPoints(segB)
			segs := []LineSegment{segA, segB}
			for _, intersection := range intersections {
				addResult(R, intersection, segs)
			}
		}
	}
	return R
}

// addResult adds an intersection point to the results tree and associates it with all segments that intersect at that point.
//
// This helper function maintains the map of segments that intersect at each point in the result tree.
// For each intersection point, it either creates a new entry or updates an existing one by adding
// all the segments that participate in the intersection.
//
// Parameters:
//   - R: The results tree mapping intersection points to sets of segments
//   - p: The intersection point to add
//   - segs: A slice of segments that intersect at point p
//
// Behavior:
//   - If the point already exists in the results, it adds the new segments to the existing set
//   - If the point is new, it creates a new map containing all the provided segments
//   - The function ensures that each segment appears only once in the set for a given point
func addResult(R *rbtree.Tree[point.Point, map[LineSegment]struct{}], p point.Point, segs []LineSegment) {
	log.Printf("adding result: %v, %v", p, segs)
	n, exists := R.Search(p)
	if exists {
		segsMap := R.Value(n)
		for _, seg := range segs {
			segsMap[seg] = struct{}{}
		}
		R.Insert(p, segsMap)
		return
	}

	segsMap := make(map[LineSegment]struct{})
	for _, seg := range segs {
		segsMap[seg] = struct{}{}
	}
	R.Insert(p, segsMap)
}

// sanitizeInput preprocesses the input line segments to remove duplicates and degenerate segments.
//
// This function prepares the input data for intersection algorithms by:
// 1. Removing duplicate segments (those with identical endpoints)
// 2. Filtering out degenerate segments (where both endpoints are the same point)
//
// Parameters:
//   - S: A slice of LineSegment objects to sanitize
//
// Returns:
//   - A new slice containing the unique, non-degenerate line segments
//
// Behavior:
//   - Uses a map to efficiently detect and eliminate duplicate segments
//   - Skips any segment where the upper and lower points are equal (zero-length segment)
//   - Preserves the original segments' orientations and properties
func sanitizeInput(S []LineSegment) []LineSegment {
	tmp := make(map[LineSegment]struct{}, len(S))

	// dedupe
	for _, seg := range S {

		// skip degenerate segments
		if seg.lower.Eq(seg.upper) {
			continue
		}
		tmp[seg] = struct{}{}
	}
	out := make([]LineSegment, 0, len(tmp))
	for k := range tmp {
		out = append(out, k)
	}
	return out
}
