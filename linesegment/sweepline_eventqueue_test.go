package linesegment

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestNewEventQueue(t *testing.T) {
	tests := map[string]struct {
		segments []LineSegment[int]
		expected []qItem
	}{
		"zigzag": {
			segments: []LineSegment[int]{
				New[int](0, 0, 2, 2),
				New[int](2, 2, 4, 0),
				New[int](4, 0, 6, 2),
				New[int](6, 2, 8, 0),
				New[int](1, 1, 7, 1),
			},
			expected: []qItem{
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
			Q := newEventQueue(tc.segments, options.WithEpsilon(1e-8))

			// dump event queue into slice
			actual := make([]qItem, 0, len(tc.segments)*2)
			Q.Ascend(func(item qItem) bool {
				actual = append(actual, item)
				return true
			})

			// check contents are expected
			log.Println("Actual:  ", actual)
			log.Println("Expected:", tc.expected)

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
