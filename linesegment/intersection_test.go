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
		expected []IntersectionResult[int]
	}{
		"no intersections": {
			segments: []LineSegment[int]{
				New(0, 0, 1, 1),
				New(2, 2, 3, 3),
			},
			expected: []IntersectionResult[int]{},
		},
		"single intersection": {
			segments: []LineSegment[int]{
				New(0, 0, 2, 2),
				New(0, 2, 2, 0),
			},
			expected: []IntersectionResult[int]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 1),
					InputLineSegments: []LineSegment[int]{
						New(0, 0, 2, 2),
						New(0, 2, 2, 0),
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
			expected: []IntersectionResult[int]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 1),
					InputLineSegments: []LineSegment[int]{
						New(0, 0, 3, 3),
						New(1, 0, 1, 3),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1, 2),
					InputLineSegments: []LineSegment[int]{
						New(0, 3, 3, 0),
						New(1, 0, 1, 3),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](1.5, 1.5),
					InputLineSegments: []LineSegment[int]{
						New(0, 0, 3, 3),
						New(0, 3, 3, 0),
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindIntersectionsSlow(tc.segments, tc.opts...)
			assert.ElementsMatch(t, tc.expected, actual, "unexpected intersections")
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
			expectedIntersectionType: IntersectionSegment,
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
			case IntersectionSegment:
				assert.Equal(t, IntersectionSegment, actual.IntersectionType)
				assert.Equal(t, tc.expectedResult, actual.IntersectionSegment)
			}
		})
	}
}
