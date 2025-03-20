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

type eventData struct {
	UofP []LineSegment
	LofP []LineSegment
}

func addEventPoint(
	Q *rbtree.Tree[point.Point, eventData],
	p point.Point,
	UofP []LineSegment,
	LofP []LineSegment,
) {
	// check if upper point exists
	existingP, found := Q.Search(p)
	if found {
		// if upper point exists, add linesegment to value slice
		// (upper segments are stored with the event point p)
		ed := Q.Value(existingP)
		ed.UofP = append(ed.UofP, UofP...)
		ed.LofP = append(ed.LofP, LofP...)
		Q.Insert(p, ed)
	} else {
		// otherwise, add new upper event
		ed := eventData{
			UofP: UofP,
			LofP: LofP,
		}
		Q.Insert(p, ed)
	}
}

func lexicographicOrder(a, b point.Point) bool {

	debugLog := strings.Builder{}
	defer func() {
		log.Println(debugLog.String())
	}()

	debugLog.WriteString(fmt.Sprintf("event queue order: a: %s, b: %s,  ", a, b))

	if numeric.FloatEquals(a.Y(), b.Y(), geom2d.GetEpsilon()) {
		debugLog.WriteString(fmt.Sprintf("y equal, order by x: %t", a.X() < b.X()))
		return numeric.FloatLessThan(a.X(), b.X(), geom2d.GetEpsilon())
	}
	debugLog.WriteString(fmt.Sprintf("order by y: %t", a.Y() > b.Y()))
	return numeric.FloatGreaterThan(a.Y(), b.Y(), geom2d.GetEpsilon())
}

func statusStructureOrder(p *point.Point, a, b LineSegment) bool {

	debugLog := strings.Builder{}
	defer func() {
		log.Println(debugLog.String())
	}()

	aSlope := a.Slope()
	bSlope := b.Slope()

	aIsHorizontal := aSlope == 0
	bIsHorizontal := bSlope == 0

	//aIsVertical := math.IsNaN(aSlope)
	//bIsVertical := math.IsNaN(bSlope)
	//
	//aContainsPoint := a.ContainsPoint(*p)
	//bContainsPoint := b.ContainsPoint(*p)

	var aX, bX float64

	if aIsHorizontal {
		aX = p.X()
	} else {
		aX = a.XAtY(p.Y())
	}
	if bIsHorizontal {
		bX = p.X()
	} else {
		bX = b.XAtY(p.Y())
	}

	debugLog.WriteString(fmt.Sprintf("p: %s, a: %s, b: %s, aX: %f, bX: %f, sA: %f, sB: %f, ", p, a, b, aX, bX, aSlope, bSlope))

	// if x's aren't equal, order by x
	if !numeric.FloatEquals(aX, bX, geom2d.GetEpsilon()) {
		debugLog.WriteString(fmt.Sprintf("order by x: %t", aX < bX))
		return aX < bX
	}

	// order by orientation
	o := point.Orientation(*p, a.lower, b.lower)
	if o == point.Counterclockwise {
		debugLog.WriteString("p -> a.lower -> b.lower: counterclockwise, a is less than b")
		return true
	}
	if o == point.Clockwise {
		debugLog.WriteString("p -> a.lower -> b.lower: clockwise, a is not less than b")
		return false
	}

	// identical orientation
	// one endpoint is equal to p
	if p.Eq(a.lower) {
		debugLog.WriteString("p == a.lower: a is not less than b")
		return false
	}
	if p.Eq(b.lower) {
		debugLog.WriteString("p == a.lower: a less than b")
		return true
	}
	if p.Eq(a.lower) && p.Eq(b.lower) {
		// order by slope above the line
		debugLog.WriteString(fmt.Sprintf("p == a.lower == b.lower: order by slope above the line: %t", aSlope < bSlope))
		return aSlope < bSlope
	}
	if !numeric.FloatEquals(a.lower.Y(), b.lower.Y(), geom2d.GetEpsilon()) {
		// order by lower y
		debugLog.WriteString(fmt.Sprintf("order by lower y: %t", a.lower.Y() < b.lower.Y()))
		return a.lower.Y() > b.lower.Y()
	}

	// order by lower x
	if !numeric.FloatEquals(a.lower.X(), b.lower.X(), geom2d.GetEpsilon()) {
		debugLog.WriteString(fmt.Sprintf("order by lower x: %t", a.lower.X() < b.lower.X()))
		return a.lower.X() < b.lower.X()
	}

	// order by upper y
	if !numeric.FloatEquals(a.upper.Y(), b.upper.Y(), geom2d.GetEpsilon()) {
		debugLog.WriteString(fmt.Sprintf("order by lower x: %t", a.lower.X() < b.lower.X()))
		return a.upper.Y() > b.upper.Y()
	}

	return a.upper.X() < b.upper.X()

}

func FindIntersections(S []LineSegment) *rbtree.Tree[point.Point, map[LineSegment]struct{}] {
	var sweepLinePoint *point.Point

	S = sanitizeInput(S)

	// Initialize an empty event queue Q.
	// Next, insert the segment endpoints into Q; when an upper endpoint is inserted,
	// the corresponding segment should be stored with it.
	Q := rbtree.New[point.Point, eventData](lexicographicOrder)

	// Initialize an empty status structure T.
	// Hack in a dynamic comparator to find each segment's X value at sweep line's Y value
	T := rbtree.New[LineSegment, struct{}](func(a, b LineSegment) bool {
		return statusStructureOrder(sweepLinePoint, a, b)
	})

	// Initialise results
	R := rbtree.New[point.Point, map[LineSegment]struct{}](lexicographicOrder)

	// insert the upper & lower points from each line segment to the event queue
	for _, l := range S {
		addEventPoint(Q, l.upper, []LineSegment{l}, nil)
		addEventPoint(Q, l.lower, nil, []LineSegment{l})
	}

	// temp debugging
	i := 0

	for {

		// temp debugging
		i++
		log.Printf("\n\n======== ITERATION %d ========\n", i)

		// temp debugging
		log.Printf("Q:\n%s", Q.String())
		log.Printf("T:\n%s", T.String())

		// pop event
		event := Q.Min(Q.Root())

		// if event was nil, we're done
		if Q.IsNil(event) {
			break
		}

		// delete event from queue
		Q.Delete(event)

		eventPoint := Q.Key(event)
		if sweepLinePoint == nil {
			sweepLinePoint = &eventPoint
		}
		ed := Q.Value(event)

		// temp debugging
		log.Printf("popped event: %s\n", event.String())

		// handle the event
		handleEventPoint(Q, T, R, eventPoint, ed.UofP, ed.LofP, sweepLinePoint)
	}

	// temp debugging
	log.Println("======== RESULTS ========")
	log.Printf("results: %v\n", R)

	// todo: remove when stable and robust
	// The check below asserts that the status structure T should be empty at this point.
	// If not empty, something has gone wrong.
	if T.Size() != 0 {
		panic("status not empty at end of FindIntersections")
	}

	return R
}

func handleEventPoint(
	Q *rbtree.Tree[point.Point, eventData],
	T *rbtree.Tree[LineSegment, struct{}],
	R *rbtree.Tree[point.Point, map[LineSegment]struct{}],
	p point.Point,
	UofP []LineSegment,
	LofP []LineSegment,
	sweepLinePoint *point.Point,
) {

	//previousSweepLinePoint := *sweepLinePoint

	// Find all segments stored in T that contain p
	// update sweepLinePoint for search
	//*sweepLinePoint = p
	log.Printf("Status structure sweep line point at:%s\n", sweepLinePoint)
	CofP := findSegmentsContainingPoint(T, p)

	// temp debugging
	log.Printf("UofP: %v\n", UofP)
	log.Printf("LofP: %v\n", LofP)
	log.Printf("CofP: %v\n", CofP)

	// if LofP ∪ UofP ∪ CofP contains more than one segment...
	if len(LofP)+len(UofP)+len(CofP) > 1 {

		// temp debugging
		log.Println("LofP ∪ UofP ∪ CofP contains more than one segment, thus is an intersection")

		// then Report p as an intersection, together with L(p), U(p), and C(p).
		segs := append(LofP, append(UofP, CofP...)...)
		addResult(R, p, segs)

		// temp debugging
		log.Printf("intersection result: %v %v\n", p, segs)
	}

	// Delete the segments in LofP ∪ CofP from T.
	// temp debugging
	log.Println("Delete the segments in LofP ∪ CofP from T")
	// todo (performance): can remove CofP and instead swap directly in tree

	// update sweepLinePoint for deletion
	//*sweepLinePoint = previousSweepLinePoint
	log.Printf("Status structure sweep line point at:%s\n", sweepLinePoint)
	for _, l := range append(LofP, CofP...) {
		node, found := T.Search(l)
		if !found {
			panic(fmt.Errorf("could not find node to delete: %v\n", l))
			//continue
		}
		// temp debugging
		log.Printf("deleting: %v\n", node)

		// todo (performance): could we store the node pointers so we dont have to search for them?
		T.Delete(node)
	}
	// temp debugging
	log.Printf("T after deletion:\n%s", T.String())

	// update sweepLinePoint for insertion
	*sweepLinePoint = p
	log.Printf("Status structure sweep line point at:%s\n", sweepLinePoint)

	// Insert the segments in UofP ∪ CofP into T.
	// The order of the segments in T should correspond to the order in which
	// they are intersected by a sweep line just below p.
	// If there is a horizontal segment, it comes last among all segments containing p.
	// todo (performance): can remove CofP and instead swap directly in tree
	// temp debugging
	log.Println("Insert the segments in UofP ∪ CofP into T")
	for _, l := range append(UofP, CofP...) {
		node, _ := T.Insert(l, struct{}{})
		// temp debugging
		log.Printf("inserted: %v\n", node)
	}
	// temp debugging
	log.Printf("T after insertion:\n%s", T.String())

	// if UofP ∪ CofP = 0
	// temp debugging
	log.Printf("UofP ∪ CofP: %d\n", len(UofP)+len(CofP))
	if len(UofP)+len(CofP) == 0 {

		// Let sl and sr be the left and right neighbors of p in T.
		// temp debugging
		log.Println("Let sl and sr be the left and right neighbors of p in T")
		sl, sr := findNeighborsOfPoint(T, p)
		// temp debugging
		log.Printf("sl: %v\n", sl)
		log.Printf("sr: %v\n", sr)
		if sl != nil && sr != nil {
			// find new event
			findNewEvent(Q, R, *sl, *sr, p)
		}
	} else {

		// find the leftmost & rightmost segments of UofP ∪ CofP
		// temp debugging
		log.Println("find the leftmost & rightmost segments of UofP ∪ CofP")
		tempT := rbtree.New[LineSegment, struct{}](func(a, b LineSegment) bool {
			return statusStructureOrder(&p, a, b)
		})
		// todo (performance): `append(UofP, CofP...)` used multiple times, break into own variable
		for _, l := range append(UofP, CofP...) {
			tempT.Insert(l, struct{}{})
		}
		leftmost := tempT.Min(tempT.Root())
		rightmost := tempT.Max(tempT.Root())
		// temp debugging
		log.Printf("leftmost:  %s\n", leftmost)
		log.Printf("rightmost: %s\n", rightmost)

		// Let sPrime be the leftmost segment of UofP ∪ CofP in T.
		sPrime, sPrimeFound := T.Search(tempT.Key(leftmost))

		if sPrimeFound {
			// temp debugging
			log.Printf("sPrime: %s\n", sPrime)
			// Let sl be the left neighbor of sPrime in T.
			sl := T.Predecessor(sPrime)
			if !T.IsNil(sl) {
				// temp debugging
				log.Printf("sl: %v\n", sl)
				findNewEvent(Q, R, T.Key(sl), T.Key(sPrime), p)
			}
		}

		// Let sDoublePrime be the rightmost segment of UofP ∪ CofP in T.
		sDoublePrime, sDoublePrimeFound := T.Search(tempT.Key(rightmost))

		if sDoublePrimeFound {
			// temp debugging
			log.Printf("sDoublePrime: %s\n", sDoublePrime)
			// Let sr be the right neighbor of sDoublePrime in T.
			sr := T.Successor(sDoublePrime)
			if !T.IsNil(sr) {
				// temp debugging
				log.Printf("sr: %v\n", sr)
				findNewEvent(Q, R, T.Key(sPrime), T.Key(sr), p)
			}
		}
	}
}

func findSegmentsContainingPoint(
	T *rbtree.Tree[LineSegment, struct{}],
	p point.Point,
) []LineSegment {

	centerNode := T.Root()
	parentNode := T.Parent(centerNode)

	// navigate through the tree until we hit a nil leaf or a line segment containing the point
	for !T.IsNil(centerNode) {
		parentNode = centerNode
		seg := T.Key(centerNode)
		o := point.Orientation(seg.lower, seg.upper, p)
		if o == point.Clockwise {
			centerNode = T.Right(centerNode)
		} else if o == point.Counterclockwise {
			centerNode = T.Left(centerNode)
		} else {
			break
		}
	}

	if !T.IsNil(centerNode) {
		centerNode = parentNode
	}

	// if centerNode is a nil leaf at this point, then no segments contain p
	if T.IsNil(centerNode) {
		return nil
	}

	CofP := make([]LineSegment, 0, 1)

	// process center node
	seg := T.Key(centerNode)
	if !seg.upper.Eq(p) && !seg.lower.Eq(p) && seg.ContainsPoint(p) { // Ensures correct handling of floating-point cases
		CofP = append(CofP, seg)
	}

	// head left
	predecessorNode := T.Predecessor(centerNode)
	for !T.IsNil(predecessorNode) {
		seg = T.Key(predecessorNode)
		if !seg.upper.Eq(p) && !seg.lower.Eq(p) && seg.ContainsPoint(p) {
			CofP = append(CofP, seg)
		} else {
			break
		}
		predecessorNode = T.Predecessor(predecessorNode)
	}

	// head right
	successorNode := T.Successor(centerNode)
	for !T.IsNil(successorNode) {
		seg = T.Key(successorNode)
		if !seg.upper.Eq(p) && !seg.lower.Eq(p) && seg.ContainsPoint(p) {
			CofP = append(CofP, seg)
		} else {
			break
		}
		successorNode = T.Successor(successorNode)
	}

	return CofP
}

func findNeighborsOfPoint(
	T *rbtree.Tree[LineSegment, struct{}],
	p point.Point,
) (sl, sr *LineSegment) {

	centerNode := T.Root()
	parentNode := T.Parent(centerNode)

	// navigate through the tree until we hit a nil leaf or a line segment containing the point
	for !T.IsNil(centerNode) {
		parentNode = T.Parent(centerNode)
		seg := T.Key(centerNode)
		o := point.Orientation(seg.lower, seg.upper, p)
		if o == point.Clockwise {
			centerNode = T.Left(centerNode)
		} else if o == point.Counterclockwise {
			centerNode = T.Right(centerNode)
		} else {
			// todo: will this break with horizontal?
			break
		}
	}

	// if we have found a line segment that contains the point:
	if !T.IsNil(centerNode) {
		// find left neighbor
		leftNode := centerNode
		for !T.IsNil(leftNode) && T.Key(leftNode).ContainsPoint(p) {
			leftNode = T.Predecessor(leftNode)
		}
		if !T.IsNil(leftNode) {
			slKey := T.Key(leftNode)
			sl = &slKey
		}

		// find right neighbor
		rightNode := centerNode
		for !T.IsNil(rightNode) && T.Key(rightNode).ContainsPoint(p) {
			rightNode = T.Predecessor(rightNode)
		}
		if !T.IsNil(rightNode) {
			srKey := T.Key(rightNode)
			sr = &srKey
		}
	}

	// if we haven't found a segment that contains the point
	if !T.IsNil(parentNode) {

		seg := T.Key(parentNode)
		o := point.Orientation(seg.lower, seg.upper, p)
		if o == point.Clockwise {

			// point would be a left child of parent
			rightNode := parentNode
			leftNode := T.Predecessor(rightNode)
			if !T.IsNil(leftNode) {
				slKey := T.Key(leftNode)
				sl = &slKey
			}
			if !T.IsNil(rightNode) {
				srKey := T.Key(rightNode)
				sr = &srKey
			}

		} else {

			// point would be a right child of parent
			leftNode := parentNode
			rightNode := T.Successor(leftNode)
			if !T.IsNil(leftNode) {
				slKey := T.Key(leftNode)
				sl = &slKey
			}
			if !T.IsNil(rightNode) {
				srKey := T.Key(rightNode)
				sr = &srKey
			}

		}
	}

	return
}

func findNewEvent(
	Q *rbtree.Tree[point.Point, eventData],
	R *rbtree.Tree[point.Point, map[LineSegment]struct{}],
	sl, sr LineSegment,
	p point.Point,
) {

	intersectionPoints, intersects := sl.IntersectionPoints(sr)
	if !intersects {
		return
	}

	// intersectionPoint must be below sweep line or on sweep line and to the right of current event point
	for _, intersectionPoint := range intersectionPoints {
		if numeric.FloatLessThan(intersectionPoint.Y(), p.Y(), geom2d.GetEpsilon()) ||
			(numeric.FloatEquals(intersectionPoint.Y(), p.Y(), geom2d.GetEpsilon()) &&
				numeric.FloatGreaterThan(intersectionPoint.X(), p.X(), geom2d.GetEpsilon())) {

			addEventPoint(Q, intersectionPoint, nil, nil)

			addResult(R, intersectionPoint, []LineSegment{sl, sr})

			// temp debug
			log.Printf("findNewEvent added event point: %s\n", intersectionPoint)
		}
	}
}

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
