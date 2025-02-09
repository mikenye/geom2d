package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"log"
	"math"
	"slices"
)

// statusItem represents an entry in the sweep line status structure.
//
// The status structure maintains active line segments as the sweep line
// progresses through event points in the plane-sweep algorithm. Each
// statusItem corresponds to a single line segment that is currently
// intersected by the sweep line.
//
// The sweep line status helps efficiently determine segment relationships,
// such as finding neighboring segments for intersection testing.
//
// In the context of the sweep line algorithm, statusItem instances are
// stored in a balanced search tree or sorted list to maintain order based
// on their relative position to the sweep line.
type statusItem struct {
	segment LineSegment[float64]
}

func debugPrintStatus(S []statusItem, y float64) {
	log.Printf("Status structure at y=%f:\n", y)
	for _, s := range S {
		xaty := s.segment.XAtY(y)
		log.Printf("  - %s (x=%f @ y=%f)\n", s.segment.String(), xaty, y)
	}
}

// deleteSegmentsFromStatus removes the specified line segments from the status structure S.
//
// This function takes a slice of status items (S) representing active segments in the sweep line algorithm
// and removes any segments that match those in the segments slice. It uses the Eq method for comparison.
//
// Parameters:
//   - S: The current status structure, which holds active segments in the sweep line algorithm.
//   - segments: A slice of line segments to be removed from S.
//   - opts: Optional geometry configuration options (such as epsilon for floating-point comparisons).
//
// Returns:
//   - A new slice of status items with the specified segments removed.
//
// Example usage:
//
//	segments := []linesegment.LineSegment[float64]{
//		linesegment.New(1, 2, 3, 4),
//		linesegment.New(5, 6, 7, 8),
//	}
//
//	status := []statusItem{
//		{segment: linesegment.New(1, 2, 3, 4)},
//		{segment: linesegment.New(5, 6, 7, 8)},
//		{segment: linesegment.New(9, 10, 11, 12)},
//	}
//
//	updatedStatus := deleteSegmentsFromStatus(status, segments)
//	log.Println(updatedStatus) // Output: [{(9,10)(11,12)}]
//
// This function is used in the sweep line algorithm to efficiently remove segments that are no longer active.
func deleteSegmentsFromStatus(
	S []statusItem,
	segments []LineSegment[float64],
	opts ...options.GeometryOptionsFunc,
) []statusItem {
	return slices.DeleteFunc(S, func(item statusItem) bool {
		for _, seg := range segments {
			if item.segment.Eq(seg, opts...) {
				return true
			}
		}
		return false
	})
}

// FindIntersectionsFast computes all intersection points and overlapping segments among a set of line segments
// using the sweep line algorithm, outlined in Section 2.1 of [Computational Geometry: Algorithms and Applications].
//
// This function efficiently finds intersections in O((n + k) log n) time, where n is the number of input segments
// and k is the number of intersections found. It is significantly faster than the naive O(n²) method when dealing
// with large input sets.
//
// The function removes duplicate segments before processing and ignores degenerate segments (where Start == End).
// It uses a balanced search tree (B-tree) as the event queue to maintain efficient event processing.
//
// Parameters:
//   - segments: A slice of line segments to check for intersections.
//   - opts: Optional geometry configuration options (e.g., epsilon for floating-point precision).
//   - If [options.WithEpsilon] is provided, the function performs an approximate equality check,
//     considering the points equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - A slice of IntersectionResult[float64] containing all intersection points and overlapping segments found.
//
// Algorithm Overview:
//  1. Deduplicate Input: Removes exact duplicate segments.
//  2. Event Queue Initialization: Inserts segment endpoints into a B-tree event queue (Q).
//     Upper endpoints store the corresponding segments.
//     Degenerate segments are skipped.
//  3. Status Structure Initialization: Uses a balanced search tree (S) to track active segments.
//  4. Sweep Line Execution: Iterates through event points in Q, updating S and computing intersections.
//  5. Result Collection: Returns a slice of intersection results.
//
// Performance Considerations:
//   - If the number of segments is small, the naive O(n²) method (FindIntersectionsSlow) may be faster
//     due to lower constant overhead.
//
// See Also:
//   - FindIntersectionsSlow for the naive O(n²) intersection detection method.
//
// [Computational Geometry: Algorithms and Applications]: https://www.springer.com/gp/book/9783540779735
func FindIntersectionsFast[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) []IntersectionResult[float64] {

	// Initialize an empty event queue
	Q := newEventQueue(segments, opts...)

	// Initialize results
	R := newIntersectionResults[float64](opts...)

	// Initialize an empty status structure S
	// (in the book they use T, but that would clobber the generic type T).
	S := make([]statusItem, 0)

	// while Q is not empty
	iterCount := 0
	for Q.Len() > 0 {
		iterCount++

		log.Printf("\n\n\n---ITERATION %d---\n\n\n", iterCount)

		// DEBUGGING: show queue
		debugPrintQueue(Q)

		// Determine the next event point p in Q and delete it
		p, ok := Q.DeleteMin()
		if !ok {
			panic(fmt.Errorf("unexpected empty queue"))
		}

		log.Printf("\nPopped: %s\n", p)

		// Handle the event
		S = handleEventPoint(p, Q, S, R, opts...)
	}

	return R.Results()
}

// findLeftmostAndRightmostSegmentAndNeighbors identifies the leftmost and rightmost segments among U(p) ∪ C(p),
// as well as their immediate neighbors in the status structure S.
//
// This function is a key part of the sweep line algorithm, helping to determine which segments
// should be tested for new intersections.
//
// Parameters:
//   - p: The current event point being processed.
//   - UofP: The set of segments whose upper endpoint is at p.
//   - CofP: The set of segments that contain p but are neither upper nor lower endpoints.
//   - S: The status structure (sweep line status), which maintains the order of active segments.
//   - opts: Optional geometry configuration options (e.g., epsilon for floating-point precision).
//
// Returns:
//   - sPrime: The leftmost segment in U(p) ∪ C(p), or nil if no segments exist.
//   - sDoublePrime: The rightmost segment in U(p) ∪ C(p), or nil if no segments exist.
//   - sL: The left neighbor of sPrime in S, or nil if no left neighbor exists.
//   - sR: The right neighbor of sDoublePrime in S, or nil if no right neighbor exists.
//
// Algorithm Overview:
//  1. Combine U(p) and C(p): The function first merges U(p) and C(p) into a single slice.
//  2. Sort UCofP: The segments are sorted using segmentSortLess, which ensures proper ordering below p.
//  3. Identify Leftmost & Rightmost Segments: The first and last elements in the sorted list are selected.
//  4. Find the Segments in S: The function locates sPrime and sDoublePrime in the sweep line status.
//  5. Determine Neighboring Segments: If sPrime or sDoublePrime are found in S, their immediate neighbors are returned.
//
// Performance Notes:
// - Sorting is performed using slices.SortStableFunc, ensuring segments are ordered correctly relative to p.
// - The function assumes the status structure S maintains a valid order of segments as the sweep line progresses.
// - Future TODO: we should be able to get better performance by switching from a slice to a binary search tree, as outlined in Section 2.1 of the book [Computational Geometry: Algorithms and Applications].
//
// [Computational Geometry: Algorithms and Applications]: https://www.springer.com/gp/book/9783540779735
func findLeftmostAndRightmostSegmentAndNeighbors(
	p point.Point[float64],
	UofP, CofP []LineSegment[float64],
	S []statusItem,
	opts ...options.GeometryOptionsFunc,
) (
	sPrime, sDoublePrime *LineSegment[float64],
	sL, sR *statusItem,
) {

	// Step 1: Combine U(p) and C(p)
	UCofP := append([]LineSegment[float64]{}, UofP...)
	UCofP = append(UCofP, CofP...)

	// Step 2: Sort UCofP by custom rules
	slices.SortStableFunc(UCofP, func(a, b LineSegment[float64]) int {
		if segmentSortLess(a, b, p, opts...) {
			return -1
		}
		return 1
	})

	log.Println("\nUC(p):", UCofP)

	// Step 3: Get the leftmost & rightmost segment (first & ladt elements in sorted UCofP)
	if len(UCofP) == 0 {
		return nil, nil, nil, nil // No segments to process
	}
	leftmost := UCofP[0]
	rightmost := UCofP[len(UCofP)-1]

	// Step 4: Find leftmost & rightmost segment in S
	var leftmostIndex int
	var leftmostFound bool
	for i, item := range S {
		if item.segment.Eq(leftmost, opts...) {
			leftmostIndex = i - 1
			leftmostFound = true
			sPrime = &leftmost
			break
		}
	}

	var rightmostIndex int
	var rightmostFound bool
	for i, item := range S {
		if item.segment.Eq(rightmost, opts...) {
			rightmostIndex = i + 1
			rightmostFound = true
			sDoublePrime = &rightmost
			break
		}
	}

	// Step 5: Determine neighbors
	if leftmostIndex >= 0 && leftmostIndex < len(S) && leftmostFound {
		sL = &S[leftmostIndex]
	}
	if rightmostIndex >= 0 && rightmostIndex < len(S) && rightmostFound {
		sR = &S[rightmostIndex]
	}

	return sPrime, sDoublePrime, sL, sR
}

// findNeighbors identifies the left and right neighbors of a given event point p in the status structure S.
//
// This function is a key component of the sweep line algorithm, used to determine which segments
// are adjacent to the event point in the active sweep line status.
//
// Parameters:
//   - S: The status structure, which maintains the order of active segments as the sweep line progresses.
//   - p: The event point with associated segments that are being processed.
//   - opts: Optional geometry configuration options (e.g., epsilon for floating-point precision).
//
// Returns:
//   - sl: Pointer to the left neighbor segment in S, or nil if no left neighbor exists.
//   - sr: Pointer to the right neighbor segment in S, or nil if no right neighbor exists.
//
// Algorithm Overview:
//  1. Locate p in the status structure S:
//     - If a segment contains p, its index is recorded.
//     - If no segment contains p, the function finds the first segment whose x-coordinate at p's y-value is greater than p.x.
//  2. Determine left and right neighbors:
//     - If p was found in S, the left neighbor is S[index-1] (if valid).
//     - The right neighbor depends on whether p was found exactly or was approximated.
//
// Performance Notes:
// - The function performs a linear scan through S, making it O(n) in complexity.
// - The logic ensures robustness when dealing with horizontal and near-vertical segments.
func findNeighbors(
	S []statusItem,
	p qItem,
	opts ...options.GeometryOptionsFunc,
) (sl, sr *statusItem) {
	if len(S) == 0 {
		return nil, nil
	}

	// Step 1: Locate p in the status structure
	index := -1
	found := false
	for i, item := range S {
		segment := item.segment
		if segment.ContainsPoint(p.point, opts...) { // Check if p lies on this segment
			index = i
			found = true
			break
		} else if segment.XAtY(p.point.Y()) > p.point.X() { // If we've passed p, break early
			index = i
			break
		}
	}

	if index == -1 {
		return nil, nil
	}

	// Step 2: Find the left and right neighbors
	if index > 0 && index-1 >= 0 {
		sl = &S[index-1] // Left neighbor
	}
	if index <= len(S)-1 {
		if found {
			if index+1 <= len(S)-1 {
				sr = &S[index+1] // Right neighbor
			}
		} else {
			sr = &S[index] // closest segment to the right of p
		}
	}

	return sl, sr
}

// findNewEvent determines if two adjacent segments in the status structure intersect below the current event point
// and inserts the intersection as a future event in the event queue Q.
//
// This function is a critical part of the Bentley-Ottmann sweep line algorithm, ensuring that all intersections are
// detected efficiently.
//
// Parameters:
//   - sl: The left segment in the status structure.
//   - sr: The right segment in the status structure.
//   - p: The current event point being processed.
//   - Q: The event queue, storing upcoming intersection events.
//   - R: The intersection results accumulator.
//   - opts: Optional geometry configuration options (e.g., epsilon for floating-point precision).
//
// Algorithm Overview:
//  1. Compute the intersection between sl and sr.
//  2. If the segments overlap, record an IntersectionOverlappingSegment in the results.
//  3. If they do not intersect at a single point, exit early.
//  4. Extract the intersection point and check its position relative to p.
//  5. If the intersection is below the sweep line, or exactly on it but to the right of p, proceed.
//  6. If the intersection is already in the event queue, exit early to avoid duplicates.
//  7. Otherwise, insert the intersection point into the event queue.
//
// Conditions for Skipping an Intersection:
//   - If the intersection lies above the current event point, it is ignored.
//   - If the intersection is at the same Y-level as the current event but lies to its left, it is ignored.
//   - If the intersection is already queued, no action is taken.
//
// Performance Considerations:
// - This function runs in **O(1)** in most cases, with the intersection calculation itself being the most expensive step.
// - It avoids redundant work by preventing duplicate intersection events from entering Q.
func findNewEvent(
	sl, sr LineSegment[float64],
	p point.Point[float64],
	Q *btree.BTreeG[qItem],
	R *intersectionResults[float64],
	opts ...options.GeometryOptionsFunc,
) {

	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	// Find the intersection between segments sl and sr.
	intersection := sl.Intersection(sr, opts...)
	log.Printf("\nIntersection between %s & %s: %s\n", sl, sr, intersection)

	if intersection.IntersectionType == IntersectionOverlappingSegment {
		R.Add(IntersectionResult[float64]{
			IntersectionType:   IntersectionOverlappingSegment,
			OverlappingSegment: intersection.OverlappingSegment,
			InputLineSegments:  []LineSegment[float64]{sl, sr},
		}, opts...)
	}

	if intersection.IntersectionType != IntersectionPoint {
		return // No intersection, so nothing to do.
	}

	// Extract the intersection point.
	newPoint := intersection.IntersectionPoint

	// if sl and sr intersect below the sweep line, or on it and to the right of the
	// current event point p, and the intersection is not yet present as an
	// event in Q then Insert the intersection point as an event into Q.

	if numeric.FloatGreaterThan(newPoint.Y(), p.Y(), geoOpts.Epsilon) || // skip point above sweep line
		(numeric.FloatEquals(newPoint.Y(), p.Y(), geoOpts.Epsilon) && numeric.FloatLessThanOrEqualTo(newPoint.X(), p.X(), geoOpts.Epsilon)) { // skip point on swwp line and to the left of or equal to current event point p
		return // The point is above or equal to the current event point, so skip it.
	}

	// Check if the intersection point is already in Q.
	exists := Q.Has(qItem{point: newPoint})
	if exists {
		return // Point is already in Q, so skip insertion.
	}

	// Insert the intersection point into Q
	qi := qItem{
		point: newPoint,
	}
	Q.ReplaceOrInsert(qi)

	return
}

// handleEventPoint processes an event point p in the sweep line algorithm, updating the event queue (Q),
// the status structure (S), and the intersection results (R). This function is a core part of the
// sweep line algorithm for finding intersections among a set of line segments.
//
// Parameters:
//   - p: The current event point being processed, which may be an endpoint or an intersection.
//   - Q: The event queue storing future event points.
//   - S: The status structure maintaining the active segments intersecting the sweep line at p.
//   - R: The intersection results accumulator, collecting detected intersections.
//   - opts: Optional geometry configuration settings, including precision tolerance.
//
// Algorithm Overview:
//
// Identify three key sets of segments related to p:
//   - U(p): Segments whose upper endpoint is p (start at p).
//   - L(p): Segments whose lower endpoint is p (end at p).
//   - C(p): Segments that contain p in their interior.
//
// If multiple segments pass through p, report it as an intersection.
//
// Remove L(p) ∪ C(p) from the status structure S ("∪" is "unioned with").
//
// Insert U(p) ∪ C(p) into S, ensuring correct order in the status structure.
//
// Determine neighboring segments in S and use findNewEvent to check for potential intersections.
//
// Conditions for Reporting an Intersection:
//
//   - If L(p) ∪ U(p) ∪ C(p) contains more than one segment, p is reported as an intersection.
//   - The slow intersection method (FindIntersectionsSlow) is used to confirm and merge intersections.
//
// Sorting & Order Maintenance:
//
//   - S must be ordered such that it corresponds to the order in which segments would be intersected
//     by a sweep line just below p.
//   - Horizontal segments always come last among all segments containing p.
//
// Handling New Events:
//
//   - If U(p) ∪ C(p) = 0, the left and right neighbors of p in S are checked for new intersections.
//   - Otherwise, the leftmost (s') and rightmost (s”) segments of U(p) ∪ C(p) are found in S.
//   - The function findNewEvent is used to detect new intersections between adjacent segments.
//
// Performance Considerations:
//   - This function operates in O(log n) on average due to event queue operations and status structure updates.
//   - Sorting of U(p) ∪ C(p) ensures the sweep line maintains correct order efficiently.
func handleEventPoint(p qItem, Q *btree.BTreeG[qItem], S []statusItem, R *intersectionResults[float64], opts ...options.GeometryOptionsFunc) []statusItem {

	// Let U(p) be the set of segments whose upper endpoint is p;
	// these segments are stored with the event point p.
	// (For horizontal segments, the upper endpoint is by definition the left endpoint.)
	UofP := p.segments

	// Find all segments stored in S that contain p;
	// they are adjacent in S.
	// Binary search for the closest item
	segments := make([]LineSegment[float64], 0)
	for _, item := range S {
		containsPoint := item.segment.ContainsPoint(p.point, opts...)
		if containsPoint {
			segments = append(segments, item.segment)
		}
	}

	// Let L(p) denote the subset of segments found whose lower endpoint is p.
	// Let C(p) denote the subset of segments found that contain p in their interior.
	LofP := make([]LineSegment[float64], 0)
	CofP := make([]LineSegment[float64], 0)
	for _, seg := range segments {
		if seg.End().Eq(p.point, opts...) {
			LofP = append(LofP, seg)
		} else if !seg.Start().Eq(p.point, opts...) {
			CofP = append(CofP, seg)
		}
	}

	log.Printf("\nU(p): ")
	for _, seg := range UofP {
		log.Printf("%s ", seg)
	}
	log.Printf("\nC(p): ")
	for _, seg := range CofP {
		log.Printf("%s ", seg)
	}
	log.Printf("\nL(p): ")
	for _, seg := range LofP {
		log.Printf("%s ", seg)
	}
	log.Print("\n")

	// if L(p) ∪ U(p) ∪ C(p) contains more than one segment...
	if len(LofP)+len(UofP)+len(CofP) > 1 {

		log.Println("\nL(p) ∪ U(p) ∪ C(p) contains more than one segment")

		// then Report p as an intersection, together with L(p), U(p), and C(p).
		for _, result := range FindIntersectionsSlow(append(LofP, append(UofP, CofP...)...), opts...) {
			R.Add(result, opts...)
		}
	}

	// Delete segments in L(p) ∪ C(p) from S
	S = deleteSegmentsFromStatus(S, LofP, opts...)
	sortStatusBySweepLine(S, p, opts...) // Re-sort to account for new sweep line position

	// DEBUGGING: show status of status structure
	log.Print("\n")
	debugPrintStatus(S, p.point.Y())

	// Insert the segments in U(p) ∪ C(p) into S.
	// The order of the segments in S should correspond to the order in which they are
	// intersected by a sweep line just below p. If there is a horizontal segment, it comes
	// last among all segments containing p.
	for _, seg := range append(UofP, CofP...) {
		// Ensure segment is not already in S
		alreadyInS := slices.ContainsFunc(S, func(item statusItem) bool {
			return item.segment.Eq(seg, opts...)
		})
		if !alreadyInS {
			S = append(S, statusItem{
				segment: seg,
			})
		} //else {
		// skip duplicate
		//}

	}
	sortStatusBySweepLine(S, p, opts...) // Re-sort S after insertion

	// DEBUGGING: show status of status structure
	log.Print("\n")
	debugPrintStatus(S, p.point.Y())

	// If U(p) ∪ C(p) = 0, find neighbors in S and call FINDNEWEVENT
	if len(UofP)+len(CofP) == 0 {

		log.Println("\nU(p) ∪ C(p) = 0")

		sL, sR := findNeighbors(S, p, opts...)

		log.Println("sL:", sL)
		log.Println("sR:", sR)

		if sL != nil && sR != nil {
			// findNewEvent(*sl, *sr, p) (from book)
			findNewEvent(sL.segment, sR.segment, p.point, Q, R, opts...)
		}

	} else {
		// Let s' be the leftmost segment of U(p) ∪ C(p) in S.
		// Let sl be the left neighbor of s' in S.
		// Let s'' be the rightmost segment of U(p) ∪ C(p) in S.
		// Let sr be the right neighbor of s'' in S.
		sPrime, sDoublePrime, sL, sR := findLeftmostAndRightmostSegmentAndNeighbors(p.point, UofP, CofP, S, opts...)

		log.Println("\ns': ", sPrime)
		log.Println("s'':", sDoublePrime)
		log.Println("sL: ", sL)
		log.Println("sR: ", sR)

		// findNewEvent(sl,s', p) (from book)
		if sPrime != nil && sL != nil {
			findNewEvent(sL.segment, *sPrime, p.point, Q, R, opts...)
		}

		// findNewEvent(s'',sr, p) (from book)
		if sDoublePrime != nil && sR != nil {
			findNewEvent(*sDoublePrime, sR.segment, p.point, Q, R, opts...)
		}
	}

	return S
}

// sortStatusBySweepLine sorts the status structure (S) based on the ordering of line segments
// at the current sweep line position (p). This ensures that segments are processed in the correct
// order as the sweep line progresses through the plane.
//
// Sorting Logic:
//   - Uses segmentSortLess to determine the relative order of segments at the current event point p.
//   - Ensures the order in S matches the order in which segments are intersected by a horizontal
//     sweep line positioned just below p.
//
// Parameters:
//   - S: The status structure containing active line segments ([]statusItem).
//   - p: The current event point (qItem) where sorting occurs.
//   - opts: Additional geometry options, such as epsilon, for floating-point comparisons.
func sortStatusBySweepLine(S []statusItem, p qItem, opts ...options.GeometryOptionsFunc) {

	// Sort using a custom comparison function
	slices.SortFunc(S, func(a, b statusItem) int {
		if segmentSortLess(a.segment, b.segment, p.point, opts...) {
			return -1
		}
		return 1
	})
}

// segmentSortLess determines the relative ordering of two line segments (a and b)
// at a given sweep line position (p). It is used in the sweep line algorithm
// to maintain the correct ordering of segments in the status structure.
//
// Sorting Rules:
//
// Primary Order: XAtY Comparison
//   - The function first compares where each segment intersects the horizontal sweep line at p.Y().
//   - If the X-coordinates differ, the segment with the smaller X-coordinate is considered "less."
//
// Handling Vertical Segments: If one segment is vertical and intersects the other at p, the diagonal segment's slope determines order:
//   - A negative slope (`\`) means it should be after the vertical segment in the ordering.
//   - A positive slope (`/`) means it should be before the vertical segment in the ordering.
//
// Handling Horizontal Segments:
//   - Horizontal segments should always be considered *after* any non-horizontal segments that pass through p.
//
// Tie-breaking with Slope: If both segments intersect at the same X-coordinate, slopes are compared to break ties:
//   - Segments sloping upward (positive slope) come after those sloping downward.
//   - If both have negative slopes, the one with a steeper slope comes first.
//   - If both have positive slopes, the one with a less steep slope comes first.
//
// Parameters:
//   - a: The first line segment (LineSegment[float64]).
//   - b: The second line segment (LineSegment[float64]).
//   - p: The current event point (point.Point[float64]) where the comparison occurs.
//   - opts: Additional geometry options, including epsilon for floating-point comparisons.
//
// Returns:
//   - true if a should be ordered before b in the status structure, otherwise false.
func segmentSortLess(a, b LineSegment[float64], p point.Point[float64], opts ...options.GeometryOptionsFunc) bool {

	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	aX := a.XAtY(p.Y())
	aSlope := a.Slope()
	aIsHorizontal := math.IsNaN(aX)
	aIsVertical := math.IsNaN(aSlope)
	aContainsP := a.ContainsPoint(p, opts...)

	bX := b.XAtY(p.Y())
	bSlope := b.Slope()
	bIsHorizontal := math.IsNaN(bX)
	bIsVertical := math.IsNaN(bSlope)
	bContainsP := b.ContainsPoint(p, opts...)

	// for horizontal lines, artificially truncate start position to point,
	// since we don't care about anything to the left, as that is considered above the sweep line
	if math.IsNaN(aX) {
		aX = p.X()
	}
	if math.IsNaN(bX) {
		bX = p.X()
	}

	log.Printf(
		"is %s (x=%f, s=%f) to the left of %s (x=%f, s=%f) at %s: ",
		a.String(),
		aX,
		aSlope,
		b.String(),
		bX,
		bSlope,
		p.String(),
	)

	// Vertical segment ordering logic: Handle cases where a vertical segment intersects a diagonal one.
	if aIsVertical && aContainsP && numeric.FloatEquals(aX, p.X(), geoOpts.Epsilon) && !bIsVertical && !bIsHorizontal && bContainsP {
		log.Println(bSlope < 0, "via slope & intersection with vertical (a)")
		return bSlope < 0
	}
	if bIsVertical && bContainsP && numeric.FloatEquals(bX, p.X(), geoOpts.Epsilon) && !aIsVertical && !bIsHorizontal && aContainsP {
		log.Println(aSlope > 0, "via slope & intersection with vertical (b)")
		return aSlope > 0
	}

	// Horizontal lines still come last if they contain p
	if aIsHorizontal && b.ContainsPoint(p, opts...) && !bIsHorizontal {
		log.Println("false via horizontal handling (a is horizontal, b contains p)")
		return false
	}
	if bIsHorizontal && a.ContainsPoint(p, opts...) && !aIsHorizontal {
		log.Println("true via horizontal handling (b is horizontal, a contains p)")
		return true
	}

	// if both horizontal, order by end x
	if aIsHorizontal && bIsHorizontal {
		log.Println("true via horizontal handling (both horizontal)")
		return a.End().X() > b.End().X()
	}

	// If XAtY matches
	if numeric.FloatEquals(aX, bX, geoOpts.Epsilon) {

		// order by slope
		if (aSlope < 0 && bSlope > 0) || (aSlope > 0 && bSlope < 0) {
			log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes opposite")
			return aSlope > bSlope
		} else if aSlope < 0 && bSlope < 0 {
			log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both negative")
			return aSlope < bSlope
		} else {
			log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both positive")
			return aSlope < bSlope
		}
	}

	log.Println(aX < bX, "via default XAtY comparison")
	return aX < bX // Default XAtY comparison

}
