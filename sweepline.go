package geom2d

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
)

type sweeplineEventType uint8

const (
	sweeplineEventStart sweeplineEventType = iota
	sweeplineEventEnd
	sweeplineEventIntersection
)

func (t *sweeplineEventType) String() string {
	switch *t {
	case sweeplineEventStart:
		return "sweeplineEventStart"
	case sweeplineEventEnd:
		return "sweeplineEventEnd"
	case sweeplineEventIntersection:
		return "sweeplineEventIntersection"
	default:
		panic(fmt.Errorf("unsupported sweeplineEventType"))
	}
}

type sweeplineEvent struct {
	lineSegments []LineSegment[float64]
	point        Point[float64]
	eventType    sweeplineEventType

	insertionOrder uint64

	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func (e *sweeplineEvent) String() string {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("event at: %s, type: %s, line segments:\n", e.point.String(), e.eventType.String()))
	for _, l := range e.lineSegments {
		out.WriteString(fmt.Sprintf("\t%s\n", l.String()))
	}
	return out.String()
}

type sweeplineEventQueue struct {
	events     []*sweeplineEvent
	insertions uint64
}

// AddIntersectionEvent will determine the intersection between line segments A & B, and if one exists,
// an intersection event will be added to the queue.
func (pq *sweeplineEventQueue) AddIntersectionEvent(A, B LineSegment[float64], opts ...Option) {
	intersection := A.IntersectionGeometry(B, opts...)

	switch intersection.IntersectionType {
	case IntersectionNone:
		// do nothing
	case IntersectionPoint:
		// add the intersection point to the queue
		event := &sweeplineEvent{
			lineSegments: []LineSegment[float64]{A, B},
			point:        intersection.IntersectionPoint,
			eventType:    sweeplineEventIntersection,
		}
		heap.Push(pq, event)

	case IntersectionSegment:
		panic(fmt.Errorf("not yet implemented"))
	}
}

func (pq *sweeplineEventQueue) RemoveIntersection(A, B LineSegment[float64], opts ...Option) {
	queueNeedsSorting := false
	for i := 0; i < len(pq.events); i++ {
		if pq.events[i].eventType != sweeplineEventIntersection {
			continue
		}
		foundA := slices.ContainsFunc(pq.events[i].lineSegments, func(l LineSegment[float64]) bool {
			return l.Eq(A, opts...)
		})
		foundB := slices.ContainsFunc(pq.events[i].lineSegments, func(l LineSegment[float64]) bool {
			return l.Eq(B, opts...)
		})
		if !foundA || !foundB {
			continue
		}
		// remove item from underlying slice, and mark that queue will need sorting
		log.Printf("[queue] removing event: %s", pq.events[i].String())
		(*pq).events[i], (*pq).events[len(pq.events)-1] = (*pq).events[len(pq.events)-1], (*pq).events[i]
		pq.events = (*pq).events[:len(pq.events)-1]
		i--
		queueNeedsSorting = true
	}
	// if queue needs sorting (as items were removed), then do it
	if queueNeedsSorting {
		heap.Init(pq)
	}
}

func (pq *sweeplineEventQueue) Len() int {
	return len(pq.events)
}

func (pq *sweeplineEventQueue) Less(i, j int) bool {

	// compare X
	if (*pq).events[i].point.x != (*pq).events[j].point.x {
		return (*pq).events[i].point.x < (*pq).events[j].point.x
	}

	// compare Y
	if (*pq).events[i].point.y != (*pq).events[j].point.y {
		return (*pq).events[i].point.y < (*pq).events[j].point.y
	}

	// compare type
	if (*pq).events[i].eventType != (*pq).events[j].eventType {
		return (*pq).events[i].eventType < (*pq).events[j].eventType
	}

	// tiebreaker: insertion order
	return (*pq).events[i].insertionOrder < (*pq).events[j].insertionOrder

}

func (pq *sweeplineEventQueue) Push(x any) {
	item := x.(*sweeplineEvent)
	item.index = len(pq.events)
	pq.insertions++
	item.insertionOrder = pq.insertions
	log.Printf("[queue] pushing event: %s", item.String())
	pq.events = append(pq.events, item)
}

func (pq *sweeplineEventQueue) Pop() any {
	old := pq.events
	n := len(old)
	item := old[n-1]
	pq.events = old[0 : n-1]
	log.Printf("[queue] popped event: %s", item.String())
	return item
}

func (pq *sweeplineEventQueue) Swap(i int, j int) {
	(*pq).events[i], (*pq).events[j] = (*pq).events[j], (*pq).events[i]
	(*pq).events[i].index = i
	(*pq).events[j].index = j
}

func newSweeplineEventQueue(lineSegments []LineSegment[float64]) *sweeplineEventQueue {
	pq := &sweeplineEventQueue{
		events: make([]*sweeplineEvent, len(lineSegments)*2), // pre-allocate memory for start & end points of each line segment
	}
	for _, l := range lineSegments {
		// "normalize" line segment to ensure leftmost-lowest point is start point
		l = l.Normalize()

		// add each point to queue
		for _, p := range l.Points() {
			pq.events[pq.insertions] = &sweeplineEvent{
				lineSegments: []LineSegment[float64]{l},
				point:        p,
				eventType:    sweeplineEventType(pq.insertions % 2), // 0 (start) for first point, 1 (end) for next point, then back to 0....
			}
			pq.insertions++
		}
	}
	heap.Init(pq)
	return pq
}

type SweepLineResult struct {
	LineSegments       []LineSegment[float64]
	IntersectionPoints []Point[float64]
}

type sweepLineStatusEntry struct {
	lineSegment LineSegment[float64]
}

type sweepLineStatus struct {
	entries  []sweepLineStatusEntry
	currentX float64
}

// Index returns the index of line segment `segment` in sweepLineStatus
func (sls *sweepLineStatus) Index(segment LineSegment[float64]) int {
	for i, l := range sls.entries {
		if l.lineSegment.Eq(segment) {
			return i
		}
	}
	return -1
}

func (sls *sweepLineStatus) Insert(segment LineSegment[float64]) {
	entry := sweepLineStatusEntry{
		lineSegment: segment,
	}
	log.Printf("[status] inserting segment: %s", entry.lineSegment.String())
	sls.entries = append(sls.entries, entry)
	sls.Sort()
}

func (sls *sweepLineStatus) Neighbors(segment LineSegment[float64]) (above, below *sweepLineStatusEntry) {
	i := sls.Index(segment)
	if i == -1 {
		return nil, nil
	}
	if i > 0 {
		below = &sls.entries[i-1]
	}
	if i < len(sls.entries)-1 {
		above = &sls.entries[i+1]
	}
	return above, below
}

func (sls *sweepLineStatus) Remove(segment LineSegment[float64]) {
	i := sls.Index(segment)
	if i == -1 {
		return
	}
	log.Printf("[status] removing segment: %s", segment.String())
	sls.entries = append(sls.entries[:i], sls.entries[i+1:]...)
}

func (sls *sweepLineStatus) Sort() {
	slices.SortStableFunc(sls.entries, func(a, b sweepLineStatusEntry) int {
		ay, aOk := a.lineSegment.YAtX(sls.currentX)
		by, bOk := b.lineSegment.YAtX(sls.currentX)
		if (aOk && bOk) && (ay != by) {
			switch {
			case ay < by:
				return -1
			case ay > by:
				return 1
			}
		}
		if a.lineSegment.start.x != b.lineSegment.start.x {
			switch {
			case a.lineSegment.start.x < b.lineSegment.start.x:
				return -1
			case a.lineSegment.start.x > b.lineSegment.start.x:
				return 1
			}
		}
		return 0
	})
}

func (sls *sweepLineStatus) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("sweep line at x=%0.4f:\n", sls.currentX))
	for i := range sls.entries {
		s.WriteString(fmt.Sprintf("\t%s\n", sls.entries[i].lineSegment.String()))
	}
	return s.String()
}

func (sls *sweepLineStatus) Swap(A, B LineSegment[float64]) {
	iA := sls.Index(A)
	iB := sls.Index(B)
	if iA == -1 || iB == -1 {
		return
	}
	sls.entries[iA], sls.entries[iB] = sls.entries[iB], sls.entries[iA]
}

func SweepLine[T SignedNumber](linesegments []LineSegment[T], opts ...Option) SweepLineResult {

	// convert input to float64
	inputSegments := make([]LineSegment[float64], len(linesegments))
	for i := range linesegments {
		inputSegments[i] = linesegments[i].AsFloat64()
	}

	// Initialize a priority queue `pq` of potential future events, each associated with a point
	// in the plane and prioritized by the x-coordinate of the point.
	// So, initially, `pq` contains an event for each of the endpoints of the input segments.
	pq := newSweeplineEventQueue(inputSegments)

	// prepare result
	result := SweepLineResult{
		LineSegments:       make([]LineSegment[float64], 0, pq.Len()),
		IntersectionPoints: make([]Point[float64], 0),
	}

	// perform Bentley-Ottmann sweep line

	// initialise the "sweep line status tree"
	status := sweepLineStatus{
		entries:  make([]sweepLineStatusEntry, 0, len(pq.events)),
		currentX: -math.MaxFloat64,
	}

	// for each event in the queue
	for pq.Len() > 0 {

		// pop event from heap
		event := heap.Pop(pq).(*sweeplineEvent)

		// set scan line X
		status.currentX = event.point.x

		switch event.eventType {

		case sweeplineEventStart:
			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// If p is the left endpoint of a line segment s, insert s into T.
			status.Insert(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// Find the line-segments r and t that are respectively immediately above and below s in T (if they exist)
			above, below := status.Neighbors(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// if the crossing of r and t (the neighbours of s in the status data structure) forms a potential future
			// event in the event queue, remove this possible future event from the event queue
			if above != nil && below != nil {
				pq.RemoveIntersection(above.lineSegment, below.lineSegment, opts...)
			}

			if above != nil {
				// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
				// If s crosses r or t, add those crossing points as potential future events in the event queue.
				pq.AddIntersectionEvent(event.lineSegments[0], above.lineSegment, opts...)
			}

			if below != nil {
				// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
				// If s crosses r or t, add those crossing points as potential future events in the event queue.
				pq.AddIntersectionEvent(event.lineSegments[0], below.lineSegment, opts...)
			}

		case sweeplineEventEnd:

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// Find the segments r and t that (prior to the removal of s) were respectively immediately
			// above and below it in T (if they exist).
			above, below := status.Neighbors(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// If p is the right endpoint of a line segment s, remove s from T.
			status.Remove(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// If r and t cross, add that crossing point as a potential future event in the event queue.
			if above != nil && below != nil {
				pq.AddIntersectionEvent(above.lineSegment, below.lineSegment, opts...)
			}

		case sweeplineEventIntersection:

			// add intersection to result
			result.IntersectionPoints = append(result.IntersectionPoints, event.point)
			log.Printf("[result] added intersection: %s", event.point.String())

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm
			// If p is the crossing point of two segments s and t (with s below t to the left of the crossing),
			// swap the positions of s and t in T.
			if event.lineSegments[0].start.y < event.lineSegments[1].start.y {
				s := event.lineSegments[0]
				t := event.lineSegments[1]
				status.Swap(s, t)

				// After the swap, find the segments r and u (if they exist) that are immediately
				// below and above t and s, respectively.
				_, r := status.Neighbors(t)
				u, _ := status.Neighbors(s)

				// Remove any crossing points rs (i.e. a crossing point between r and s) and tu
				// (i.e. a crossing point between t and u) from the event queue and,
				// if r and t cross or s and u cross, add those crossing points to the event queue.
				if r != nil {
					pq.RemoveIntersection(r.lineSegment, s)
					pq.AddIntersectionEvent(r.lineSegment, t)
				}
				if u != nil {
					pq.RemoveIntersection(t, u.lineSegment)
					pq.AddIntersectionEvent(s, u.lineSegment)
				}
			}

			if event.lineSegments[1].start.y < event.lineSegments[0].start.y {
				s := event.lineSegments[1]
				t := event.lineSegments[0]
				status.Swap(s, t)

				// After the swap, find the segments r and u (if they exist) that are immediately
				// below and above t and s, respectively.
				_, r := status.Neighbors(t)
				u, _ := status.Neighbors(s)

				// Remove any crossing points rs (i.e. a crossing point between r and s) and tu
				// (i.e. a crossing point between t and u) from the event queue and,
				// if r and t cross or s and u cross, add those crossing points to the event queue.
				if r != nil {
					pq.RemoveIntersection(r.lineSegment, s)
					pq.AddIntersectionEvent(r.lineSegment, t)
				}
				if u != nil {
					pq.RemoveIntersection(t, u.lineSegment)
					pq.AddIntersectionEvent(s, u.lineSegment)
				}
			}
		}
		log.Printf("[status] %s", status.String())
	}

	return result
}
