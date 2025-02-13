package linesegment

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindIntersectionsSlow(t *testing.T) {
	tests := map[string]struct {
		segments []LineSegment[int]
		opts     []options.GeometryOptionsFunc
		expected []IntersectionResult[float64]
	}{
		"no intersections": {
			segments: []LineSegment[int]{
				New(0, 0, 1, 1),
				New(2, 2, 3, 3),
			},
			expected: []IntersectionResult[float64]{},
		},
		"single intersection": {
			segments: []LineSegment[int]{
				New(0, 0, 2, 2),
				New(0, 2, 2, 0),
			},
			expected: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 1),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 2, 2),
						New[float64](0, 2, 2, 0),
					},
				},
			},
		},
		"multiple intersections": {
			segments: []LineSegment[int]{
				New(0, 0, 3, 3),
				New(0, 3, 3, 0),
				New(1, 0, 1, 3),
			},
			expected: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 1),
					InputLineSegments: []LineSegment[float64]{
						New[float64](1, 3, 1, 0),
						New[float64](0, 0, 3, 3),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1.5, 1.5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 3, 3, 0),
						New[float64](3, 3, 0, 0),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 2),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 3, 3, 0),
						New[float64](1, 3, 1, 0),
					},
				},
			},
		},
		"square shape": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0),
				New[int](10, 0, 10, 10),
				New[int](10, 10, 0, 10),
				New[int](0, 10, 0, 0),
			},
			expected: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](0, 10),
					InputLineSegments: []LineSegment[float64]{
						New[float64](10, 10, 0, 10),
						New[float64](0, 10, 0, 0),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](10, 10),
					InputLineSegments: []LineSegment[float64]{
						New[float64](10, 0, 10, 10),
						New[float64](10, 10, 0, 10),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](10, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0),
						New[float64](10, 0, 10, 10),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](0, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0),
						New[float64](0, 10, 0, 0),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindIntersectionsSlow(tc.segments, tc.opts...)
			t.Log("Expected:", tc.expected)
			t.Log("Actual:  ", actual)
			assert.True(t, InterSectionResultsEq(actual, tc.expected, options.WithEpsilon(1e-8)))
		})
	}
}

func TestLineSegment_Intersection(t *testing.T) {
	tests := map[string]struct {
		AB, CD                   LineSegment[int]
		expectedIntersectionType IntersectionType
		expectedResult           any
	}{
		"Intersecting segments returning point": {
			AB:                       New(0, 0, 10, 10),
			CD:                       New(0, 10, 10, 0),
			expectedIntersectionType: IntersectionPoint,
			expectedResult:           point.New(5.0, 5.0),
		},
		"Intersecting collinear segments returning line segment": {
			AB:                       New(0, 0, 10, 0),
			CD:                       New(-5, 0, 5, 0),
			expectedIntersectionType: IntersectionOverlappingSegment,
			expectedResult:           New(0.0, 0.0, 5.0, 0.0),
		},
		"Non-intersecting segments": {
			AB:                       New(0, 0, 5, 5),
			CD:                       New(6, 6, 10, 10),
			expectedIntersectionType: IntersectionNone,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.AB.Intersection(tc.CD)
			assert.ElementsMatch(t, []LineSegment[int]{tc.AB, tc.CD}, actual.InputLineSegments)
			switch tc.expectedIntersectionType {
			case IntersectionNone:
				assert.Equal(t, IntersectionNone, actual.IntersectionType)
			case IntersectionPoint:
				assert.Equal(t, IntersectionPoint, actual.IntersectionType)
				assert.Equal(t, tc.expectedResult, actual.IntersectionPoint)
			case IntersectionOverlappingSegment:
				assert.Equal(t, IntersectionOverlappingSegment, actual.IntersectionType)
				assert.Equal(t, tc.expectedResult, actual.OverlappingSegment)
			}
		})
	}
}
