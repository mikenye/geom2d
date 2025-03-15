package _torefactor

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewEventQueueRBT(t *testing.T) {
	type result struct {
		point    point.Point[float64]
		segments []LineSegment[float64]
	}
	tests := map[string]struct {
		segments []LineSegment[int]
		dupe     point.Point[float64] // duplicate point, should be an upper point
		unique   point.Point[float64] // unique point, maybe an intersection
		expected []result             // should contain the correctly-ordered results of the queue
	}{
		"zigzag": {
			segments: []LineSegment[int]{
				New[int](0, 0, 2, 2),
				New[int](2, 2, 4, 0),
				New[int](4, 0, 6, 2),
				New[int](6, 2, 8, 0),
				New[int](1, 1, 7, 1),
			},
			dupe:   point.New[float64](2, 2),
			unique: point.New[float64](3, 1),
			expected: []result{
				{
					point: point.New[float64](2, 2),
					segments: []LineSegment[float64]{
						New[float64](0, 0, 2, 2),
						New[float64](2, 2, 4, 0),
					},
				},
				{
					point: point.New[float64](6, 2),
					segments: []LineSegment[float64]{
						New[float64](4, 0, 6, 2),
						New[float64](6, 2, 8, 0),
					},
				},
				{
					point: point.New[float64](1, 1),
					segments: []LineSegment[float64]{
						New[float64](1, 1, 7, 1),
					},
				},
				{
					point: point.New[float64](3, 1), // unique, intersection
					//segments: []LineSegment[float64]{},
				},
				{
					point:    point.New[float64](7, 1),
					segments: nil,
				},
				{
					point:    point.New[float64](0, 0),
					segments: nil,
				},
				{
					point:    point.New[float64](4, 0),
					segments: nil,
				},
				{
					point:    point.New[float64](8, 0),
					segments: nil,
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// build event queue
			Q := newEventQueueRBT(tc.segments, options.WithEpsilon(1e-8))

			// log contents of queue for troubleshooting
			t.Logf("Initial state of queue:\n%s", Q.String())

			// attempt to insert dupe point
			Q.InsertPoint(tc.dupe)

			// log contents of queue for troubleshooting
			t.Logf("After attempt to add dupe point %s:\n%s", tc.dupe, Q.String())

			// add unique point
			Q.InsertPoint(tc.unique)

			// log contents of queue for troubleshooting
			t.Logf("After add unique point %s:\n%s", tc.unique, Q.String())

			// dump event queue into slice
			actual := make([]result, 0, len(tc.segments)*2)
			for !Q.IsEmpty() {
				p, ls := Q.Pop()
				actual = append(actual, result{
					point:    p,
					segments: ls,
				})
			}

			// check contents are expected
			t.Log("Actual:  ", actual)
			t.Log("Expected:", tc.expected)

			require.Len(t, actual, len(tc.expected), "actual vs expected len mismatch")
			for i := range tc.expected {
				require.Equalf(t, tc.expected[i].point, actual[i].point, "point mismatch at index %d", i)
				require.Len(t, actual[i].segments, len(tc.expected[i].segments), "segment len mismatch")
				for _, expectedSeg := range tc.expected[i].segments {
					foundSeg := false
					for _, actualSeg := range actual[i].segments {
						if expectedSeg.Eq(actualSeg) {
							foundSeg = true
							break
						}
					}
					if !foundSeg {
						assert.Fail(t, "segment mismatch")
					}
				}
			}
		})
	}
}
