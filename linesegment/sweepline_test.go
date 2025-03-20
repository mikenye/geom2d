package linesegment

import (
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusStructureOrder(t *testing.T) {
	tests := map[string]struct {
		segmentA, segmentB             LineSegment
		eventPoint                     point.Point
		expectedALessThanBAtSweepLineY bool
	}{
		"X shape, above intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 20, 20, 0),
			eventPoint:                     point.New(20, 20),
			expectedALessThanBAtSweepLineY: false,
		},
		"X shape, at intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 20, 20, 0),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"X shape, below intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 20, 20, 0),
			eventPoint:                     point.New(0, 0),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment horizontal, one diagonal sloped /, left of intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(0, 10),
			expectedALessThanBAtSweepLineY: false,
		},
		"one segment horizontal, one diagonal sloped /, at intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment horizontal, one diagonal sloped /, after intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(20, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment horizontal, one diagonal sloped \\, left of intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(0, 10),
			expectedALessThanBAtSweepLineY: false,
		},
		"one segment horizontal, one diagonal sloped \\, at intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment horizontal, one diagonal sloped \\, after intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(20, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"two segments horizontal": {
			segmentA:                       New(0, 10, 20, 10),
			segmentB:                       New(0, 10, 19, 10),
			eventPoint:                     point.New(0, 10),
			expectedALessThanBAtSweepLineY: false,
		},
		"one segment vertical, one diagonal sloped /, above intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 20),
			expectedALessThanBAtSweepLineY: false,
		},
		"one segment vertical, one diagonal sloped /, at intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment vertical, one diagonal sloped /, below intersection": {
			segmentA:                       New(20, 20, 0, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 0),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment vertical, one diagonal sloped \\, above intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 20),
			expectedALessThanBAtSweepLineY: true,
		},
		"one segment vertical, one diagonal sloped \\, at intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: false,
		},
		"one segment vertical, one diagonal sloped \\, below intersection": {
			segmentA:                       New(0, 20, 20, 0),
			segmentB:                       New(10, 20, 10, 0),
			eventPoint:                     point.New(10, 0),
			expectedALessThanBAtSweepLineY: false,
		},
		"both vertical": {
			segmentA:                       New(10, 20, 10, 0),
			segmentB:                       New(10, 20, 10, 1),
			eventPoint:                     point.New(10, 20),
			expectedALessThanBAtSweepLineY: false,
		},
		"one vertical, one horizontal": {
			segmentA:                       New(10, 20, 10, 0),
			segmentB:                       New(0, 10, 20, 10),
			eventPoint:                     point.New(10, 10),
			expectedALessThanBAtSweepLineY: true,
		},
		"FuzzFindIntersections_2segments/313a517efce970e0": {
			segmentA:                       New(10, 20, 10, 0),
			segmentB:                       New(96, 10, 10, 1),
			eventPoint:                     point.New(10, 1),
			expectedALessThanBAtSweepLineY: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Logf("event point: %s\n", tc.eventPoint)
			t.Run("a < b", func(t *testing.T) {
				actual := statusStructureOrder(&tc.eventPoint, tc.segmentA, tc.segmentB)
				assert.Equal(t, tc.expectedALessThanBAtSweepLineY, actual)
			})
			t.Run("b < a", func(t *testing.T) {
				actual := statusStructureOrder(&tc.eventPoint, tc.segmentB, tc.segmentA)
				assert.NotEqual(t, tc.expectedALessThanBAtSweepLineY, actual)
			})
		})
	}
}

func TestFindIntersections(t *testing.T) {
	type intersectionResult struct {
		p    point.Point
		segs map[LineSegment]struct{}
	}

	tests := map[string]struct {
		segments      []LineSegment
		intersections []intersectionResult
	}{
		"2 segments/no intersection": {
			segments: []LineSegment{
				New(0, 0, 10, 10),
				New(10, 0, 20, 10),
			},
		},
		"2 segments/^ shape, intersection is segment upper points": {
			segments: []LineSegment{
				New(0, 0, 10, 10),
				New(10, 10, 20, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(10, 10, 0, 0):  {},
						New(10, 10, 20, 0): {},
					},
				},
			},
		},
		"2 segments/V shape, intersection is segment lower points": {
			segments: []LineSegment{
				New(0, 10, 10, 0),
				New(10, 0, 20, 10),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 0),
					segs: map[LineSegment]struct{}{
						New(0, 10, 10, 0):  {},
						New(20, 10, 10, 0): {},
					},
				},
			},
		},
		"2 segments/< shape, intersection is segments are lower/upper points": {
			segments: []LineSegment{
				New(0, 10, 10, 20),
				New(0, 10, 10, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(0, 10),
					segs: map[LineSegment]struct{}{
						New(10, 20, 0, 10): {},
						New(0, 10, 10, 0):  {},
					},
				},
			},
		},
		"2 segments, > shape, intersection is segments are lower/upper points": {
			segments: []LineSegment{
				New(0, 20, 10, 10),
				New(10, 10, 0, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(0, 20, 10, 10): {},
						New(10, 10, 0, 0):  {},
					},
				},
			},
		},
		"2 segments/X shape, intersection is internal to both segments": {
			segments: []LineSegment{
				New(0, 0, 10, 10),
				New(10, 0, 0, 10),
			},
			intersections: []intersectionResult{
				{
					p: point.New(5, 5),
					segs: map[LineSegment]struct{}{
						New(0, 10, 10, 0): {},
						New(10, 10, 0, 0): {},
					},
				},
			},
		},
		"2 segments/left-slanted X shape, one vertical segment, intersection is internal to both segments": {
			segments: []LineSegment{
				New(10, 20, 10, 0),
				New(0, 20, 20, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(10, 20, 10, 0): {},
						New(0, 20, 20, 0):  {},
					},
				},
			},
		},
		"2 segments/right-slanted X shape, one vertical segment, intersection is internal to both segments": {
			segments: []LineSegment{
				New(10, 20, 10, 0),
				New(20, 20, 0, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(10, 20, 10, 0): {},
						New(20, 20, 0, 0):  {},
					},
				},
			},
		},
		"2 segments/left-slanted, left-rotated X shape, one horizontal segment, intersection is internal to both segments": {
			segments: []LineSegment{
				New(0, 10, 20, 10),
				New(20, 20, 0, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(0, 10, 20, 10): {},
						New(20, 20, 0, 0):  {},
					},
				},
			},
		},
		"2 segments/right-slanted, left-rotated X shape, one horizontal segment, intersection is internal to both segments": {
			segments: []LineSegment{
				New(0, 10, 20, 10),
				New(0, 20, 20, 0),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(0, 10, 20, 10): {},
						New(0, 20, 20, 0):  {},
					},
				},
			},
		},
		"2 segments/+ shape, intersection is internal to both segments": {
			segments: []LineSegment{
				New(10, 20, 10, 0),
				New(0, 10, 20, 10),
			},
			intersections: []intersectionResult{
				{
					p: point.New(10, 10),
					segs: map[LineSegment]struct{}{
						New(10, 20, 10, 0): {},
						New(0, 10, 20, 10): {},
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			results := FindIntersections(tc.segments)
			assert.Len(t, tc.intersections, results.Size(), "number of results mismatch")
			for _, r := range tc.intersections {
				segs, found := results.Search(r.p)
				assert.True(t, found, "intersection result not found")
				assert.Equal(t, results.Value(segs), r.segs, "intersection line segments do not match")
			}

			// compare against brute force
			resultsBruteForce := FindIntersectionsBruteForce(tc.segments)
			for _, r := range tc.intersections {
				segs, found := resultsBruteForce.Search(r.p)
				assert.True(t, found, "intersection result not found")
				assert.Equal(t, resultsBruteForce.Value(segs), r.segs, "intersection line segments do not match")
			}
		})
	}
}
