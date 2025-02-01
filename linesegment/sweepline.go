package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"maps"
	"math"
	"slices"
	"strings"
)

type qItem struct {
	point    point.Point[float64]
	segments []LineSegment[float64]
}

func (qi qItem) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Queue Item: %s\n", qi.point.String()))
	for _, seg := range qi.segments {
		builder.WriteString(fmt.Sprintf("      - Segment: %s\n", seg.String()))
	}
	return builder.String()
}

type statusItem struct {
	segment LineSegment[float64]
	sweepY  float64 // Pointer to the current sweep line's y-coordinate
}

func debugPrintQueue(Q *btree.BTreeG[qItem]) {
	Qcopy := Q.Clone()
	fmt.Println("Event queue (Q):")
	for Qcopy.Len() > 0 {
		item, _ := Qcopy.DeleteMin()
		fmt.Printf("  - %s", item.String())
	}
}

func debugPrintStatus(S []statusItem) {
	fmt.Println("Status structure:")
	for _, s := range S {
		xaty := s.segment.XAtY(s.sweepY)
		fmt.Printf("  - %s (x=%f @ y=%f)\n", s.segment.String(), xaty, s.sweepY)
	}
}

func dedupeSegments[T types.SignedNumber](segments []LineSegment[T]) []LineSegment[T] {
	tmpMap := make(map[LineSegment[T]]bool)
	for _, seg := range segments {
		tmpMap[seg.normalize()] = false
	}
	deduped := make([]LineSegment[T], 0, len(segments))
	for k := range maps.Keys(tmpMap) {
		deduped = append(deduped, k)
	}
	return deduped
}

// Helper to delete segments from S
func deleteSegmentsFromStatus(
	S []statusItem,
	segments []LineSegment[float64],
	opts ...options.GeometryOptionsFunc,
) []statusItem {
	return slices.DeleteFunc(S, func(item statusItem) bool {
		for _, seg := range segments {
			if item.segment.Eq(seg, opts...) {
				// Optional: Log the removal for debugging
				fmt.Printf("Removing segment: %v\n", item.segment)
				return true
			}
		}
		return false
	})
}

func FindIntersections[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) []IntersectionResult[float64] {

	fmt.Println("FindIntersections started")

	// dedupe input
	segments = dedupeSegments(segments)

	// Initialize results
	R := newIntersectionResults[float64]()

	// Initialize an empty event queue
	Q := btree.NewG[qItem](2, qItemLess)

	// Insert the segment endpoints into Q.
	// When an upper endpoint is inserted, the corresponding segment should be stored with it.
	for i := range segments {
		// skip degenerate line segments
		if segments[i].Start().Eq(segments[i].End(), opts...) {
			fmt.Println("skipping degenerate line segment:", segments[i].String())
			continue
		}
		insertSegmentIntoQueue(segments[i].AsFloat64(), Q)
	}

	// Initialize an empty status structure S
	// (in the book they use T, but that would clobber the generic type T).
	S := make([]statusItem, 0)

	// while Q is not empty
	iternum := 1
	for Q.Len() > 0 {

		fmt.Printf("\n\n\n--- ITERATION %d ---\n\n\n", iternum)
		iternum++

		// DEBUGGING: show queue
		debugPrintQueue(Q)

		// Determine the next event point p in Q and delete it
		p, ok := Q.DeleteMin()
		if !ok {
			panic(fmt.Errorf("unexpected empty queue"))
		}

		// DEBUGGING: show popped event
		fmt.Printf("Popped: %s", p.String())

		S = handleEventPoint(p, Q, S, R, opts...)
	}

	fmt.Println("FindIntersections finished")
	return R.Results()
}

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
	//slices.SortFunc(UCofP, func(a, b LineSegment[float64]) int {
	//	xa := a.XAtY(sweepY)
	//	xb := b.XAtY(sweepY)
	//
	//	// Identify horizontal and vertical segments
	//	aIsHorizontal := a.Start().Y() == a.End().Y()
	//	bIsHorizontal := b.Start().Y() == b.End().Y()
	//	aIsVertical := a.Start().X() == a.End().X()
	//	bIsVertical := b.Start().X() == b.End().X()
	//
	//	// **1. Vertical segments should always come first**
	//	if aIsVertical && !bIsVertical && a.Start().X() == p.X() {
	//		return -1
	//	}
	//	if bIsVertical && !aIsVertical && b.Start().X() == p.X() {
	//		return 1
	//	}
	//
	//	//// **2. If XAtY returns NaN (horizontal segment), use leftmost X-coordinate instead**
	//	// 2. If XAtY returns NaN (horizontal segment), use p X-coordinate instead
	//	if math.IsNaN(xa) {
	//		//xa = math.Min(a.Start().X(), a.End().X())
	//		xa = p.X()
	//	}
	//	if math.IsNaN(xb) {
	//		//xb = math.Min(b.Start().X(), b.End().X())
	//		xb = p.X()
	//	}
	//
	//	// **3. Primary Sort: By X-coordinates at sweepY**
	//	if xa < xb {
	//		return -1
	//	}
	//	if xa > xb {
	//		return 1
	//	}
	//
	//	// **4. Horizontal segments should come last if tied by X**
	//	if aIsHorizontal && !bIsHorizontal {
	//		return 1
	//	}
	//	if bIsHorizontal && !aIsHorizontal {
	//		return -1
	//	}
	//
	//	// **5. Higher start Y should come first**
	//	if a.Start().Y() > b.Start().Y() {
	//		return -1
	//	}
	//	if a.Start().Y() < b.Start().Y() {
	//		return 1
	//	}
	//
	//	// **6. Higher end Y should come first**
	//	if a.End().Y() > b.End().Y() {
	//		return -1
	//	}
	//	if a.End().Y() < b.End().Y() {
	//		return 1
	//	}
	//
	//	// **7. Lower X-start should come first**
	//	if a.Start().X() < b.Start().X() {
	//		return -1
	//	}
	//	if a.Start().X() > b.Start().X() {
	//		return 1
	//	}
	//
	//	// **8. Lower X-end should come first**
	//	if a.End().X() < b.End().X() {
	//		return -1
	//	}
	//	if a.End().X() > b.End().X() {
	//		return 1
	//	}
	//
	//	// **9. If completely equal, return 0**
	//	return 0
	//})

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

	if sl != nil {
		fmt.Printf("Left neighbor: %v\n", *sl)
	}
	if sr != nil {
		fmt.Printf("Right neighbor: %v\n", *sr)
	}

	return sl, sr
}

// // Helper to find neighbors in S
//func findNeighbors(
//	S []statusItem,
//	p qItem,
//	opts ...options.GeometryOptionsFunc,
//) (*statusItem, *statusItem) {
//	var left, right *statusItem
//	index := -1
//	for i, item := range S {
//		fmt.Printf("findNeighbors: %s == %s ?: ", item.segment.End().String(), p.point.String())
//		if item.segment.End().Eq(p.point, opts...) {
//			fmt.Println("true")
//			index = i
//			break
//		}
//		fmt.Println("false")
//	}
//	leftIndex := index - 1
//	rightIndex := index + 1
//	if leftIndex >= 0 && leftIndex < len(S) {
//		left = &S[leftIndex]
//	}
//	if rightIndex >= 0 && rightIndex < len(S) {
//		right = &S[rightIndex]
//	}
//	fmt.Printf("Left neighbor: %v\n", left)
//	fmt.Printf("Right neighbor: %v\n", right)
//	return left, right
//}

func findNewEvent(
	sl, sr LineSegment[float64],
	p point.Point[float64],
	Q *btree.BTreeG[qItem],
	R *intersectionResults[float64],
	opts ...options.GeometryOptionsFunc,
) {

	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	fmt.Println("ENTERING findNewEvent")
	defer fmt.Println("EXITING findNewEvent")

	fmt.Printf("finding new events between %s and %s: ", sl.String(), sr.String())

	// Find the intersection between segments sl and sr.
	intersection := sl.Intersection(sr, opts...)
	fmt.Println(intersection.String())

	if intersection.IntersectionType == IntersectionOverlappingSegment {
		fmt.Println("IntersectionOverlappingSegment", intersection.OverlappingSegment.String())
		R.Add(IntersectionResult[float64]{
			IntersectionType:   IntersectionOverlappingSegment,
			OverlappingSegment: intersection.OverlappingSegment,
			InputLineSegments:  []LineSegment[float64]{sl, sr},
		}, opts...)
	}

	if intersection.IntersectionType != IntersectionPoint {
		fmt.Println("no intersection")
		return // No intersection, so nothing to do.
	}

	// Extract the intersection point.
	newPoint := intersection.IntersectionPoint

	// if sl and sr intersect below the sweep line, or on it and to the right of the
	// current event point p, and the intersection is not yet present as an
	// event in Q then Insert the intersection point as an event into Q.

	if numeric.FloatGreaterThan(newPoint.Y(), p.Y(), geoOpts.Epsilon) || // skip point above sweep line
		(numeric.FloatEquals(newPoint.Y(), p.Y(), geoOpts.Epsilon) && numeric.FloatLessThanOrEqualTo(newPoint.X(), p.X(), geoOpts.Epsilon)) { // skip point on swwp line and to the left of or equal to current event point p

		fmt.Printf("The point is above or equal to the current event point, so skip it: %s\n", newPoint.String())
		return // The point is above or equal to the current event point, so skip it.
	}

	// Check if the intersection point is already in Q.
	exists := Q.Has(qItem{point: newPoint})
	if exists {
		fmt.Println("Point is already in Q, so skip insertion.")
		return // Point is already in Q, so skip insertion.
	}

	// Insert the intersection point into Q, associating both segments with it.
	qi := qItem{
		point: newPoint,
		//segments: []LineSegment[float64]{sl, sr},
	}
	fmt.Printf("inserting item into event Q: %s\n", qi.String())
	Q.ReplaceOrInsert(qi)

	return
}

func handleEventPoint(p qItem, Q *btree.BTreeG[qItem], S []statusItem, R *intersectionResults[float64], opts ...options.GeometryOptionsFunc) []statusItem {

	fmt.Println("ENTERING handleEventPoint")
	defer fmt.Println("EXITING handleEventPoint")

	// Let U(p) be the set of segments whose upper endpoint is p;
	// these segments are stored with the event point p.
	// (For horizontal segments, the upper endpoint is by definition the left endpoint.)
	UofP := p.segments
	fmt.Println("U(p):", UofP)

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
	fmt.Println("Find all segments stored in S that contain p:", segments)

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
	fmt.Println("L(p):", LofP)
	fmt.Println("C(p):", CofP)

	// if L(p) ∪ U(p) ∪ C(p) contains more than one segment...
	if len(LofP)+len(UofP)+len(CofP) > 1 {

		// then Report p as an intersection, together with L(p), U(p), and C(p).
		fmt.Println("L(p):", LofP)
		fmt.Println("U(p):", UofP)
		fmt.Println("C(p):", CofP)

		for _, result := range FindIntersectionsSlow(append(LofP, append(UofP, CofP...)...), opts...) {
			fmt.Println("INTERSECTION: ", result)
			R.Add(result, opts...)
		}
	}

	// Delete segments in L(p) ∪ C(p) from S
	S = deleteSegmentsFromStatus(S, LofP, opts...)
	sortStatusBySweepLine(S, p, opts...) // Re-sort to account for new sweep line position

	// DEBUGGING: show status of event queue
	debugPrintStatus(S)

	// Insert the segments in U(p) ∪ C(p) into S.
	// The order of the segments in S should correspond to the order in which they are
	// intersected by a sweep line just below p. If there is a horizontal segment, it comes
	// last among all segments containing p.
	fmt.Println("Insert the segments in U(p) ∪ C(p) into S:")
	for _, seg := range append(UofP, CofP...) {
		// Ensure segment is not already in S
		alreadyInS := slices.ContainsFunc(S, func(item statusItem) bool {
			return item.segment.Eq(seg, opts...)
		})
		if !alreadyInS {
			fmt.Printf("  - ADDING: %s\n", seg.String())
			S = append(S, statusItem{
				segment: seg,
				sweepY:  p.point.Y(),
			})
		} else {
			fmt.Printf("  - SKIPPING DUPLICATE: %s\n", seg.String())
		}

	}
	sortStatusBySweepLine(S, p, opts...) // Re-sort S after insertion

	// DEBUGGING: show status of event queue
	debugPrintStatus(S)

	// If U(p) ∪ C(p) = 0, find neighbors in S and call FINDNEWEVENT
	if len(UofP)+len(CofP) == 0 {
		fmt.Println("U(p) ∪ C(p) = 0")
		sL, sR := findNeighbors(S, p, opts...)
		if sL != nil {
			fmt.Println("sl: ", sL.segment.String())
		} else {
			fmt.Println("sl: nil")
		}
		if sR != nil {
			fmt.Println("sr: ", sR.segment.String())
		} else {
			fmt.Println("sr: nil")
		}
		if sL != nil && sR != nil {
			// findNewEvent(*sl, *sr, p) (from book)
			fmt.Println("running findNewEvent(sl,sr, p)")
			findNewEvent(sL.segment, sR.segment, p.point, Q, R, opts...)
		}

	} else {
		fmt.Println("U(p) ∪ C(p) != 0")
		// Let s' be the leftmost segment of U(p) ∪ C(p) in S.
		// Let sl be the left neighbor of s' in S.
		// Let s'' be the rightmost segment of U(p) ∪ C(p) in S.
		// Let sr be the right neighbor of s'' in S.
		sPrime, sDoublePrime, sL, sR := findLeftmostAndRightmostSegmentAndNeighbors(p.point, UofP, CofP, S, opts...)

		UCofP := append(UofP, CofP...)
		fmt.Println("U(p) ∪ C(p):", UCofP)

		if sPrime != nil {
			fmt.Println("s': ", sPrime.String())
		} else {
			fmt.Println("s': nil")
		}
		if sDoublePrime != nil {
			fmt.Println("s'': ", sDoublePrime.String())
		} else {
			fmt.Println("s'': nil")
		}
		if sL != nil {
			fmt.Println("sl: ", sL.segment.String())
		} else {
			fmt.Println("sl: nil")
		}
		if sR != nil {
			fmt.Println("sr: ", sR.segment.String())
		} else {
			fmt.Println("sr: nil")
		}

		// findNewEvent(sl,s', p) (from book)
		if sPrime != nil && sL != nil {
			fmt.Println("running findNewEvent(sl,s', p)")
			findNewEvent(sL.segment, *sPrime, p.point, Q, R, opts...)
		}

		// findNewEvent(s'',sr, p) (from book)
		if sDoublePrime != nil && sR != nil {
			fmt.Println("running findNewEvent(s'',sr, p)")
			findNewEvent(*sDoublePrime, sR.segment, p.point, Q, R, opts...)
		}
	}

	return S
}

func insertSegmentIntoQueue(seg LineSegment[float64], Q *btree.BTreeG[qItem]) {

	//fmt.Println("insertSegmentIntoQueue called for:", seg.String())

	// Ensure correct ordering
	seg = seg.normalize()

	//fmt.Println("normalised to:", seg.String())

	// Check if segment is degenerate (single point)
	if seg.Start().Eq(seg.End()) {
		//fmt.Println("Degenerate segment detected, treating as a point:", seg.Start())

		// Insert the degenerate point **without associating a segment**
		if !Q.Has(qItem{point: seg.Start()}) {
			Q.ReplaceOrInsert(qItem{point: seg.Start()})
		}
		return // Don't process as a segment
	}

	// Retrieve or update the upper endpoint
	existingUpper, exists := Q.Get(qItem{point: seg.Start()})
	if exists {
		//fmt.Println("upper point already exists in Q, appending segment")
		// Append the segment to the existing qItem
		existingUpper.segments = append(existingUpper.segments, seg)
		Q.ReplaceOrInsert(existingUpper) // Re-insert the updated item back into Q
	} else {
		// Insert a new qItem for the upper endpoint
		//fmt.Println("upper point does not exist in Q, adding new queue entry")
		Q.ReplaceOrInsert(qItem{point: seg.Start(), segments: []LineSegment[float64]{seg}})
	}

	// Insert the lower endpoint as a new qItem (no associated segment)
	if !Q.Has(qItem{point: seg.End()}) {
		//fmt.Println("adding lower point queue entry")
		Q.ReplaceOrInsert(qItem{point: seg.End()})
	} else {
		//fmt.Println("lower point queue entry already exists, not adding")
	}

	//fmt.Println("state of queue:")
	//debugPrintQueue(Q)
	//fmt.Print("\n\n\n")
}

func qItemLess(p, q qItem) bool {
	// Compare based on the sweep line event order:
	// 1. Higher y-coordinates are "smaller" (processed first).
	// 2. For equal y-coordinates, smaller x-coordinates are "smaller".
	if p.point.Y() > q.point.Y() {
		return true
	}
	if p.point.Y() == q.point.Y() && p.point.X() < q.point.X() {
		return true
	}
	return false
}

// ordered left to right at sweep line y.
// If there is a horizontal segment, it comes last among all segments containing p.
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

	fmt.Printf(
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
		fmt.Println(bSlope < 0, "via slope & intersection with vertical (a)")
		return bSlope < 0
	}
	if bIsVertical && bContainsP && numeric.FloatEquals(bX, p.X(), geoOpts.Epsilon) && !aIsVertical && !bIsHorizontal && aContainsP {
		fmt.Println(aSlope > 0, "via slope & intersection with vertical (b)")
		return aSlope > 0
	}

	//// If both are vertical, sort by Y **to maintain correct order in S**
	//if aIsVertical && bIsVertical {
	//	fmt.Println(a.Start().Y() < b.Start().Y(), "via sorting two vertical segments by Y")
	//	return a.Start().Y() < b.Start().Y()
	//}

	// Horizontal lines still come last if they contain p
	if aIsHorizontal && b.ContainsPoint(p, opts...) {
		fmt.Println("false via horizontal handling (a is horizontal, b contains p)")
		return false
	}
	if bIsHorizontal && a.ContainsPoint(p, opts...) {
		fmt.Println("true via horizontal handling (b is horizontal, a contains p)")
		return true
	}

	// If XAtY matches
	if numeric.FloatEquals(aX, bX, geoOpts.Epsilon) {

		// order by slope
		if (aSlope < 0 && bSlope > 0) || (aSlope > 0 && bSlope < 0) {
			fmt.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes opposite")
			return aSlope > bSlope
		} else if aSlope < 0 && bSlope < 0 {
			fmt.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both negative")
			return aSlope < bSlope
		} else {
			fmt.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both positive")
			return aSlope < bSlope
		}
	}

	fmt.Println(aX < bX, "via default XAtY comparison")
	return aX < bX // Default XAtY comparison

}

// Helper to sort S by the current sweep line
func sortStatusBySweepLine(S []statusItem, p qItem, opts ...options.GeometryOptionsFunc) {
	// Assign the sweepY value to each status item
	for i := range S {
		S[i].sweepY = p.point.Y()
	}

	// Sort using a custom comparison function
	slices.SortFunc(S, func(a, b statusItem) int {
		if segmentSortLess(a.segment, b.segment, p.point, opts...) {
			return -1
		}
		return 1
	})
	//slices.SortFunc(S, func(a, b statusItem) int {
	//
	//	//fmt.Printf("\n\nchecking %s vs %s for y=%f\n", a.segment.String(), b.segment.String(), a.sweepY)
	//	if a.sweepY != b.sweepY {
	//		panic(fmt.Errorf("unexpected sweepY"))
	//	}
	//
	//	xa := a.segment.XAtY(p.point.Y())
	//	xb := b.segment.XAtY(p.point.Y())
	//
	//	//fmt.Printf("xa: %f\n", xa)
	//	//fmt.Printf("xb: %f\n", xb)
	//
	//	// Identify horizontal and vertical segments
	//	aIsHorizontal := a.segment.Start().Y() == a.segment.End().Y()
	//	bIsHorizontal := b.segment.Start().Y() == b.segment.End().Y()
	//	aIsVertical := a.segment.Start().X() == a.segment.End().X()
	//	bIsVertical := b.segment.Start().X() == b.segment.End().X()
	//
	//	//fmt.Printf("aIsHorizontal: %t, bIsHorizontal: %t\n", aIsHorizontal, bIsHorizontal)
	//	//fmt.Printf("aIsVertical: %t, bIsVertical: %t\n", aIsVertical, bIsVertical)
	//
	//	// **1. Vertical segments should always come first**
	//	if aIsVertical && !bIsVertical && a.segment.Start().X() == p.point.X() {
	//		//fmt.Println("Vertical segments come first: aIsVertical && !bIsVertical")
	//		return -1
	//	}
	//	if bIsVertical && !aIsVertical && b.segment.Start().X() == p.point.X() {
	//		//fmt.Println("Vertical segments come first: bIsVertical && !aIsVertical")
	//		return 1
	//	}
	//
	//	////**2. If XAtY returns NaN (horizontal segment), use leftmost X-coordinate instead**
	//	// 2. If XAtY returns NaN (horizontal segment), use p X-coordinate instead**
	//	if math.IsNaN(xa) {
	//		//xa = math.Min(a.segment.Start().X(), a.segment.End().X())
	//		xa = p.point.X()
	//		//fmt.Printf("as xa was NaN, using leftmost x value of: %f\n", xa)
	//	}
	//	if math.IsNaN(xb) {
	//		//xb = math.Min(b.segment.Start().X(), b.segment.End().X())
	//		xb = p.point.X()
	//		//fmt.Printf("as xb was NaN, using leftmost x value of: %f\n", xb)
	//	}
	//
	//	// **3. Primary Sort: By X-coordinates at sweepY**
	//	if xa < xb {
	//		//fmt.Println("Primary: Sort by X-coordinates at sweepY: xa < xb")
	//		return -1
	//	}
	//	if xa > xb {
	//		//fmt.Println("Primary: Sort by X-coordinates at sweepY: xa > xb")
	//		return 1
	//	}
	//
	//	// **4. Horizontal segments should come last if tied by X**
	//	if aIsHorizontal && !bIsHorizontal {
	//		//fmt.Println("Horizontal segments should come last: aIsHorizontal && !bIsHorizontal")
	//		return 1
	//	}
	//	if bIsHorizontal && !aIsHorizontal {
	//		//fmt.Println("Horizontal segments should come last: bIsHorizontal && !aIsHorizontal")
	//		return -1
	//	}
	//
	//	// **5. Higher start Y should come first**
	//	if a.segment.Start().Y() > b.segment.Start().Y() {
	//		//fmt.Println("Higher Y-start should come first: a.segment.Start().Y() > b.segment.Start().Y()")
	//		return -1
	//	}
	//	if a.segment.Start().Y() < b.segment.Start().Y() {
	//		//fmt.Println("Higher Y-start should come first: a.segment.Start().Y() < b.segment.Start().Y()")
	//		return 1
	//	}
	//
	//	// **6. Higher end Y should come first**
	//	if a.segment.End().Y() > b.segment.End().Y() {
	//		//fmt.Println("Higher Y-end should come first: a.segment.End().Y() > b.segment.End().Y()")
	//		return -1
	//	}
	//	if a.segment.End().Y() < b.segment.End().Y() {
	//		//fmt.Println("Higher Y-end should come first: a.segment.End().Y() < b.segment.End().Y()")
	//		return 1
	//	}
	//
	//	// **7. Lower X-start should come first**
	//	if a.segment.Start().X() < b.segment.Start().X() {
	//		//fmt.Println("Lower X-start should come first: a.segment.Start().X() < b.segment.Start().X()")
	//		return -1
	//	}
	//	if a.segment.Start().X() > b.segment.Start().X() {
	//		//fmt.Println("Lower X-start should come first: a.segment.Start().X() > b.segment.Start().X()")
	//		return 1
	//	}
	//
	//	// **8. Lower X-end should come first**
	//	if a.segment.End().X() < b.segment.End().X() {
	//		//fmt.Println("Lower X-end should come first: a.segment.End().X() < b.segment.End().X()")
	//		return -1
	//	}
	//	if a.segment.End().X() > b.segment.End().X() {
	//		//fmt.Println("Lower X-end should come first: a.segment.End().X() > b.segment.End().X()")
	//		return 1
	//	}
	//
	//	// **9. If completely equal, return 0**
	//	//fmt.Println("completely equal")
	//	return 0
	//})
}
