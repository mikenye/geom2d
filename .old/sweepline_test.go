package _old

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand/v2"
	"testing"
)

func TestSweeplineEventQueue_AddIntersectionEvents(t *testing.T) {
	// define line segments in "X" orientation
	lsA := NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](10, 10))
	lsB := NewLineSegment[float64](NewPoint[float64](0, 10), NewPoint[float64](10, 0))

	// define line segment that does not intersect
	lsC := NewLineSegment[float64](NewPoint[float64](0, 20), NewPoint[float64](10, 20))

	// prepare an event queue
	pq := newSweeplineEventQueue([]LineSegment[float64]{lsA, lsB, lsC})

	// hack in intersections
	pq.AddIntersectionEvent(lsA, lsB)
	pq.AddIntersectionEvent(lsA, lsC)
	pq.AddIntersectionEvent(lsB, lsC)

	// length of queue should be 5 at this point
	assert.Len(t, pq.events, 7, "length of queue should be 7")

	// check the queue at this point
	for i := range pq.events {
		pq.events[i].index = 0          // make index 0 for ease of checking...
		pq.events[i].insertionOrder = 0 // make insertion order 0 for ease of checking
	}
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA},
			point:        lsA.start,
			eventType:    sweeplineEventStart,
			index:        0,
		},
		"missing event for lsA start point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA},
			point:        lsA.end,
			eventType:    sweeplineEventEnd,
			index:        0,
		},
		"missing event for lsA end point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsB},
			point:        lsB.start,
			eventType:    sweeplineEventStart,
			index:        0,
		},
		"missing event for lsB start point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsB},
			point:        lsB.end,
			eventType:    sweeplineEventEnd,
			index:        0,
		},
		"missing event for lsB end point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsC},
			point:        lsC.start,
			eventType:    sweeplineEventStart,
			index:        0,
		},
		"missing event for lsB start point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsC},
			point:        lsC.end,
			eventType:    sweeplineEventEnd,
			index:        0,
		},
		"missing event for lsB end point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA, lsB},
			point:        NewPoint[float64](5, 5),
			eventType:    sweeplineEventIntersection,
			index:        0,
		},
		"missing event for intersection between lsA and lsB",
	)

	// re-init the heap due to resetting indexes earlier
	heap.Init(pq)
}

func TestSweeplineEventQueue_Less(t *testing.T) {

	tests := map[string]struct {
		lineSegments   []LineSegment[float64]
		expectedOrder  []sweeplineEvent
		randomizeInput bool
	}{
		"order by X": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 0)),
				NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](7, 1)),
				NewLineSegment(NewPoint[float64](2, 2), NewPoint[float64](8, 2)),
				NewLineSegment(NewPoint[float64](3, 3), NewPoint[float64](6, 3)),
			},
			randomizeInput: true,
			expectedOrder: []sweeplineEvent{
				{point: NewPoint[float64](0, 0), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 0))}},
				{point: NewPoint[float64](1, 1), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](7, 1))}},
				{point: NewPoint[float64](2, 2), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](2, 2), NewPoint[float64](8, 2))}},
				{point: NewPoint[float64](3, 3), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](3, 3), NewPoint[float64](6, 3))}},
				{point: NewPoint[float64](6, 3), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](3, 3), NewPoint[float64](6, 3))}},
				{point: NewPoint[float64](7, 1), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](7, 1))}},
				{point: NewPoint[float64](8, 2), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](2, 2), NewPoint[float64](8, 2))}},
				{point: NewPoint[float64](10, 0), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 0))}},
			},
		},
		"order by Y": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 13), NewPoint[float64](10, 13)),
				NewLineSegment(NewPoint[float64](0, 14), NewPoint[float64](10, 14)),
				NewLineSegment(NewPoint[float64](0, 15), NewPoint[float64](10, 15)),
				NewLineSegment(NewPoint[float64](0, 16), NewPoint[float64](10, 16)),
			},
			randomizeInput: true,
			expectedOrder: []sweeplineEvent{
				{point: NewPoint[float64](0, 13), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 13), NewPoint[float64](10, 13))}},
				{point: NewPoint[float64](0, 14), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 14), NewPoint[float64](10, 14))}},
				{point: NewPoint[float64](0, 15), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 15), NewPoint[float64](10, 15))}},
				{point: NewPoint[float64](0, 16), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 16), NewPoint[float64](10, 16))}},
				{point: NewPoint[float64](10, 13), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 13), NewPoint[float64](10, 13))}},
				{point: NewPoint[float64](10, 14), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 14), NewPoint[float64](10, 14))}},
				{point: NewPoint[float64](10, 15), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 15), NewPoint[float64](10, 15))}},
				{point: NewPoint[float64](10, 16), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 16), NewPoint[float64](10, 16))}},
			},
		},
		"order by eventType": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](5, 5)),
				NewLineSegment(NewPoint[float64](5, 5), NewPoint[float64](10, 0)),
			},
			randomizeInput: true,
			expectedOrder: []sweeplineEvent{
				{point: NewPoint[float64](0, 0), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](5, 5))}},
				{point: NewPoint[float64](5, 5), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](5, 5), NewPoint[float64](10, 0))}},
				{point: NewPoint[float64](5, 5), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](5, 5))}},
				{point: NewPoint[float64](10, 0), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](5, 5), NewPoint[float64](10, 0))}},
			},
		},
		"order by insertion order": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 10)),
				NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 0)),
			},
			randomizeInput: false,
			expectedOrder: []sweeplineEvent{
				{point: NewPoint[float64](0, 5), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 10))}},
				{point: NewPoint[float64](0, 5), eventType: sweeplineEventStart, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 0))}},
				{point: NewPoint[float64](10, 0), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 0))}},
				{point: NewPoint[float64](10, 10), eventType: sweeplineEventEnd, lineSegments: []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 5), NewPoint[float64](10, 10))}},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// randomize input
			if tc.randomizeInput {
				rand.Shuffle(len(tc.lineSegments), func(i, j int) { tc.lineSegments[i], tc.lineSegments[j] = tc.lineSegments[j], tc.lineSegments[i] })
			}

			// show input
			t.Log("input line segments:")
			for i := range tc.lineSegments {
				t.Logf("\t%s", tc.lineSegments[i].String())
			}

			// create priority queue
			pq := newSweeplineEventQueue(tc.lineSegments)

			// check order of queue
			t.Log("queue order:")
			for i := range tc.expectedOrder {
				event := heap.Pop(pq).(*sweeplineEvent)
				t.Logf("\t%s", event.String())
				assert.Equal(t, tc.expectedOrder[i].point, event.point, "out of order item encountered")
				assert.Equal(t, tc.expectedOrder[i].eventType, event.eventType, "item type mismatch")
				assert.Equal(t, tc.expectedOrder[i].lineSegments, event.lineSegments, "linesegment mismatch")
			}
		})
	}
}

func TestSweeplineEventQueue_RemoveIntersection(t *testing.T) {
	// define line segments in "X" orientation
	lsA := NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](10, 10))
	lsB := NewLineSegment[float64](NewPoint[float64](0, 10), NewPoint[float64](10, 0))

	// prepare an event queue
	pq := newSweeplineEventQueue([]LineSegment[float64]{lsA, lsB})

	// hack in the intersection
	pq.AddIntersectionEvent(lsA, lsB)

	// check the queue at this point
	for i := range pq.events {
		pq.events[i].index = 0          // make index 0 for ease of checking...
		pq.events[i].insertionOrder = 0 // make insertionOrder 0 for ease of checking...
	}

	// length of queue should be 4 at this point
	assert.Len(t, pq.events, 5, "length of queue should be 5")

	// ensure intersection is present before removing in
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA, lsB},
			point:        NewPoint[float64](5, 5),
			eventType:    sweeplineEventIntersection,
			index:        0,
		},
		"missing event for intersection betweem lsA & lsB",
	)

	// re-init the heap due to resetting indexes earlier
	heap.Init(pq)

	// remove end point of lsA
	pq.RemoveIntersection(lsA, lsB)

	// length of queue should be 4 at this point
	assert.Len(t, pq.events, 4, "length of queue should be 5")

	// check the queue at this point
	for i := range pq.events {
		pq.events[i].index = 0          // make index 0 for ease of checking...
		pq.events[i].insertionOrder = 0 // make insertionOrder 0 for ease of checking...
	}
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA},
			point:        lsA.start,
			eventType:    sweeplineEventStart,
			index:        0,
		},
		"missing event for lsA start point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA},
			point:        lsA.end,
			eventType:    sweeplineEventEnd,
			index:        0,
		},
		"missing event for lsA end point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsB},
			point:        lsB.start,
			eventType:    sweeplineEventStart,
			index:        0,
		},
		"missing event for lsB start point",
	)
	assert.Contains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsB},
			point:        lsB.end,
			eventType:    sweeplineEventEnd,
			index:        0,
		},
		"missing event for lsB end point",
	)
	assert.NotContains(t, pq.events,
		&sweeplineEvent{
			lineSegments: []LineSegment[float64]{lsA, lsB},
			point:        NewPoint[float64](5, 5),
			eventType:    sweeplineEventIntersection,
			index:        0,
		},
		"event for intersection between lsA and lsB should be removed",
	)
}

func TestSweepLine_Intersections(t *testing.T) {
	tests := map[string]struct {
		lineSegments          []LineSegment[float64]
		expectedIntersections []Point[float64]
	}{
		"no intersection": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 0)),
				NewLineSegment(NewPoint[float64](0, 10), NewPoint[float64](10, 10)),
			},
			expectedIntersections: []Point[float64]{},
		},
		"one intersection ('X' type shape)": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
				NewLineSegment(NewPoint[float64](0, 10), NewPoint[float64](10, 0)),
			},
			expectedIntersections: []Point[float64]{
				NewPoint[float64](5, 5),
			},
		},
		"two intersections ('≠' type shape)": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 4)),
				NewLineSegment(NewPoint[float64](0, 2), NewPoint[float64](10, 0)),
				NewLineSegment(NewPoint[float64](0, 4), NewPoint[float64](10, 2)),
			},
			expectedIntersections: []Point[float64]{
				NewPoint[float64](3.33330, 1.33330),
				NewPoint[float64](6.66667, 2.66667),
			},
		},
		"three intersections": {
			lineSegments: []LineSegment[float64]{
				NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 4)),
				NewLineSegment(NewPoint[float64](0, 2), NewPoint[float64](10, 0)),
				NewLineSegment(NewPoint[float64](0, 4), NewPoint[float64](10, 2)),
				NewLineSegment(NewPoint[float64](1, -1), NewPoint[float64](10, 1)),
			},
			expectedIntersections: []Point[float64]{
				NewPoint[float64](3.3333, 1.3333),
				NewPoint[float64](6.6667, 2.6667),
				NewPoint[float64](7.6316, 0.4737),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// run sweepline
			result := SweepLine(tc.lineSegments, WithEpsilon(1e-4))

			// check length
			require.Len(t, result.IntersectionPoints, len(tc.expectedIntersections), "result len mismatch")

			// ensure points match
			for i := range tc.expectedIntersections {
				foundPoint := false
				for j := range result.IntersectionPoints {
					if result.IntersectionPoints[j].Eq(tc.expectedIntersections[i], WithEpsilon(1e-4)) {
						foundPoint = true
						break
					}
				}
				assert.Truef(t, foundPoint, "expected %s in result.IntersectionPoints", tc.expectedIntersections[i].String())
			}
		})
	}
}
