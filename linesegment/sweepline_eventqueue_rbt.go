package linesegment

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"strings"
)

type eventQueueRBT struct {
	queue *rbt.Tree
}

//func (Q *eventQueueRBT) Get(p point.Point[float64]) ([]LineSegment[float64], bool) {
//	v, exists := Q.queue.Get(p)
//	if !exists {
//		return nil, exists
//	}
//	return v.([]LineSegment[float64]), exists
//}

func (Q *eventQueueRBT) IsEmpty() bool {
	return Q.queue.Empty()
}

// Pop returns the event point & U(p) for the event point
func (Q *eventQueueRBT) Pop() (point.Point[float64], []LineSegment[float64]) {
	node := Q.queue.Left()
	if node == nil {
		panic(fmt.Errorf("tried to pop from empty queue"))
	}
	Q.queue.Remove(node.Key)
	return node.Key.(point.Point[float64]), node.Value.([]LineSegment[float64])
}

func (Q *eventQueueRBT) InsertPoint(p point.Point[float64]) {
	// does the point exist?
	_, exists := Q.queue.Get(p)
	if exists {
		return
	}
	Q.queue.Put(p, []LineSegment[float64]{})
}

func (Q *eventQueueRBT) String() string {
	out := strings.Builder{}
	iter := Q.queue.Iterator()
	i := 0
	for iter.Next() {
		k := iter.Key().(point.Point[float64])
		v := iter.Value().([]LineSegment[float64])
		out.WriteString(fmt.Sprintf("Event Queue Item %d: %s [U(p):", i, k.String()))
		if len(v) > 0 {
			for _, s := range v {
				out.WriteString(" ")
				out.WriteString(s.String())
			}
		} else {
			out.WriteString(" <empty>")
		}
		out.WriteString("]\n")
		i++
	}
	return out.String()
}

//type eventQueueItem struct {
//	// point is the point.Point of the event
//	point point.Point[float64]
//
//	// segments is the list of LineSegments where point is the upper point of the segment
//	segments []LineSegment[float64]
//}

//func (qi eventQueueItem) String() string {
//	builder := strings.Builder{}
//	builder.WriteString(fmt.Sprintf("Queue Item: %s, U(p): ", qi.point.String()))
//	first := true
//	for _, seg := range qi.segments {
//		if first {
//			builder.WriteString(seg.String())
//			first = false
//			continue
//		}
//		builder.WriteString(fmt.Sprintf(", %s", seg.String()))
//	}
//	return builder.String()
//}

func eventQueueComparator(a interface{}, b interface{}) int {
	// Should return a number:
	// negative , if a < b
	// zero     , if a == b
	// positive , if a > b

	// from the book:
	// "If p and q are two event points then we have p≺q if and only if py > qy holds or py = qy and px < qx holds."
	p := a.(point.Point[float64])
	q := b.(point.Point[float64])
	if p.Y() > q.Y() || (p.Y() == q.Y() && p.X() < q.X()) {
		return -1 // less
	}
	if p.Eq(q) {
		return 0 // equal
	}
	return 1
}

func newEventQueueRBT[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) *eventQueueRBT {

	Q := new(eventQueueRBT)
	Q.queue = rbt.NewWith(eventQueueComparator)

	// Insert the segment endpoints into Q.
	// When an upper endpoint is inserted, the corresponding segment should be stored with it.
	for _, seg := range segments {

		segf := seg.AsFloat64()
		upper, upperDegenerate := segf.sweeplineUpperPoint()
		lower, lowerDegenerate := segf.sweeplineLowerPoint()
		upperSegments := []LineSegment[float64]{NewFromPoints(upper, lower)}

		// skip degenerate points
		if upperDegenerate || lowerDegenerate {
			continue
		}

		// check if upper point is in the queue already
		existingUpperSegments, upperExists := Q.queue.Get(upper)

		// if upper point does exist, delete it, merge line segments
		// else, add to queue
		if upperExists {
			Q.queue.Remove(upper)
			upperSegments = mergeSegments(existingUpperSegments.([]LineSegment[float64]), upperSegments, opts...)
		}
		Q.queue.Put(upper, upperSegments)

		// check if upper point is in the queue already
		_, lowerExists := Q.queue.Get(lower)

		// if lower point does exist, do nothing
		// else, add to queue
		if !lowerExists {
			Q.queue.Put(lower, []LineSegment[float64]{})
		}
	}
	return Q
}

// it needs:
//  - creator
//  - len or isEmpty method
//  - pop method (pop lowest)
//  - "has" method (lookup)
//  - "insert" method
