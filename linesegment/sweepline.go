package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
	"slices"
	"sort"
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
		builder.WriteString(fmt.Sprintf("  - Segment: %s\n", seg.String()))
	}
	return builder.String()
}

type statusItem struct {
	segment LineSegment[float64]
	sweepY  float64 // Pointer to the current sweep line's y-coordinate
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

	// Step 2: Sort UCofP by XAtY
	sweepY := p.Y()
	sort.Slice(UCofP, func(i, j int) bool {
		return UCofP[i].XAtY(sweepY) < UCofP[j].XAtY(sweepY)
	})

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

// Helper to find neighbors in S
func findNeighbors(
	S []statusItem,
	p qItem,
	opts ...options.GeometryOptionsFunc,
) (*statusItem, *statusItem) {
	var left, right *statusItem
	index := -1
	for i, item := range S {
		if item.segment.End().Eq(p.point, opts...) {
			index = i
			break
		}
	}
	leftIndex := index - 1
	rightIndex := index + 1
	if leftIndex >= 0 && leftIndex < len(S) {
		left = &S[leftIndex]
	}
	if rightIndex >= 0 && rightIndex < len(S) {
		right = &S[rightIndex]
	}
	fmt.Printf("Left neighbor: %v\n", left)
	fmt.Printf("Right neighbor: %v\n", right)
	return left, right
}

func findNewEvent(
	s, t LineSegment[float64],
	currentPoint point.Point[float64],
	Q *btree.BTreeG[qItem],
	R *intersectionResults[float64],
	opts ...options.GeometryOptionsFunc,
) {

	fmt.Printf("finding new events between %s and %s: ", s.String(), t.String())

	// Find the intersection between segments s and t.
	intersection := s.Intersection(t, opts...)

	if intersection.IntersectionType == IntersectionOverlappingSegment {
		fmt.Println("IntersectionOverlappingSegment", intersection.OverlappingSegment.String())
		R.Add(IntersectionResult[float64]{
			IntersectionType:   IntersectionOverlappingSegment,
			OverlappingSegment: intersection.OverlappingSegment,
			InputLineSegments:  []LineSegment[float64]{s, t},
		}, opts...)
	}

	if intersection.IntersectionType != IntersectionPoint {
		fmt.Println("no intersection")
		return // No intersection, so nothing to do.
	}

	// Extract the intersection point.
	newPoint := intersection.IntersectionPoint

	// Check if the intersection point lies strictly below the current event point.
	if newPoint.Y() > currentPoint.Y() ||
		(newPoint.Y() == currentPoint.Y() && newPoint.X() <= currentPoint.X()) {
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
		//segments: []LineSegment[float64]{s, t},
	}
	fmt.Printf("inserting item into event Q: %s\n", qi.String())
	Q.ReplaceOrInsert(qi)

	return
}

func insertSegmentIntoQueue(seg LineSegment[float64], Q *btree.BTreeG[qItem]) {

	// Ensure correct ordering
	seg = seg.normalize()

	// Retrieve or update the upper endpoint
	existingUpper, exists := Q.Get(qItem{point: seg.Start()})
	if exists {
		// Append the segment to the existing qItem
		existingUpper.segments = append(existingUpper.segments, seg)
		Q.ReplaceOrInsert(existingUpper) // Re-insert the updated item back into Q
	} else {
		// Insert a new qItem for the upper endpoint
		Q.ReplaceOrInsert(qItem{point: seg.Start(), segments: []LineSegment[float64]{seg}})
	}

	// Insert the lower endpoint as a new qItem (no associated segment)
	Q.ReplaceOrInsert(qItem{point: seg.End(), segments: nil})
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

// Helper to sort S by the current sweep line
func sortStatusBySweepLine(S []statusItem, sweepY float64) {
	for i := range S {
		S[i].sweepY = sweepY
	}
	slices.SortFunc(S, func(a, b statusItem) int {
		if statusItemLess(a, b) {
			return -1
		}
		return 1
	})
}

func statusItemLess(a, b statusItem) bool {

	xa := a.segment.XAtY(a.sweepY)
	xb := b.segment.XAtY(b.sweepY)

	// Check if one segment is vertical
	aIsVertical := math.IsNaN(a.segment.Slope())
	bIsVertical := math.IsNaN(b.segment.Slope())

	// Vertical segments take precedence
	if aIsVertical && !bIsVertical {
		return true
	}
	if bIsVertical && !aIsVertical {
		return false
	}

	// Compare x-coordinates at the current sweep line
	if xa < xb {
		return true
	}
	if xa > xb {
		return false
	}

	// Final tie-breaking by segment start points
	if a.segment.start.Y() != b.segment.start.Y() {
		return a.segment.start.Y() > b.segment.start.Y()
	}
	if a.segment.start.X() != b.segment.start.X() {
		return a.segment.start.X() < b.segment.start.X()
	}

	// Compare end points as a final fallback
	if a.segment.end.Y() != b.segment.end.Y() {
		return a.segment.end.Y() > b.segment.end.Y()
	}
	return a.segment.end.X() < b.segment.end.X()
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
		fmt.Printf("  - %s\n", s.segment.String())
	}
}

func FindIntersections[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) []IntersectionResult[float64] {

	// Initialize results
	R := newIntersectionResults[float64]()

	// Initialize an empty event queue
	Q := btree.NewG[qItem](2, qItemLess)

	// Insert the segment endpoints into Q.
	// When an upper endpoint is inserted, the corresponding segment should be stored with it.
	for i := range segments {
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
	return R.Results()
}

func handleEventPoint(p qItem, Q *btree.BTreeG[qItem], S []statusItem, R *intersectionResults[float64], opts ...options.GeometryOptionsFunc) []statusItem {

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
		if item.segment.ContainsPoint(p.point, opts...) {
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
		fmt.Printf("Intersection(s) at: %s\n", p.point.String())
		fmt.Println("  - L(p):", LofP)
		fmt.Println("  - U(p):", UofP)
		fmt.Println("  - C(p):", CofP)

		for _, result := range FindIntersectionsSlow(append(LofP, append(UofP, CofP...)...), opts...) {
			R.Add(result)
		}
	}

	// Delete segments in L(p) ∪ C(p) from S
	S = deleteSegmentsFromStatus(S, LofP, opts...)
	sortStatusBySweepLine(S, p.point.Y()) // Re-sort to account for new sweep line position

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
	sortStatusBySweepLine(S, p.point.Y()) // Re-sort S after insertion

	// DEBUGGING: show status of event queue
	debugPrintStatus(S)

	// If U(p) ∪ C(p) = 0, find neighbors in S and call FINDNEWEVENT
	if len(UofP)+len(CofP) == 0 {
		sL, sR := findNeighbors(S, p)
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
