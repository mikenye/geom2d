package _torefactor

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIntersectionResults_Point_MergeResults(t *testing.T) {

	R := newIntersectionResults[float64]()

	A := IntersectionResult[float64]{
		IntersectionType:   IntersectionPoint,
		IntersectionPoint:  point.New[float64](1, 10),
		OverlappingSegment: LineSegment[float64]{},
		InputLineSegments: []LineSegment[float64]{
			New[float64](1, 15, 1, 5),
			New[float64](0, 10, 10, 10),
		},
	}

	B := IntersectionResult[float64]{
		IntersectionType:   IntersectionPoint,
		IntersectionPoint:  point.New[float64](1, 10),
		OverlappingSegment: LineSegment[float64]{},
		InputLineSegments: []LineSegment[float64]{
			New[float64](2, 11, 0, 9),
			New[float64](0, 10, 10, 10),
		},
	}

	expected := IntersectionResult[float64]{
		IntersectionType:   IntersectionPoint,
		IntersectionPoint:  point.New[float64](1, 10),
		OverlappingSegment: LineSegment[float64]{},
		InputLineSegments: []LineSegment[float64]{
			New[float64](1, 15, 1, 5),
			New[float64](0, 10, 10, 10),
			New[float64](2, 11, 0, 9),
		},
	}

	R.Add(A)
	R.Add(B)

	results := R.Results()

	require.Len(t, results, 1)
	assert.Equal(t, results[0].IntersectionType, expected.IntersectionType)
	assert.Equal(t, results[0].IntersectionPoint, expected.IntersectionPoint)
	assert.Equal(t, results[0].OverlappingSegment, expected.OverlappingSegment)
	assert.ElementsMatch(t, results[0].InputLineSegments, expected.InputLineSegments)

}

func TestIntersectionResults_Eq(t *testing.T) {
	tests := map[string]struct {
		A, B     IntersectionResult[float64]
		epsilon  float64
		expected bool
	}{
		"equal point, input segments reversed and different order": {
			A: IntersectionResult[float64]{
				IntersectionType:   IntersectionPoint,
				IntersectionPoint:  point.New[float64](1, 10),
				OverlappingSegment: LineSegment[float64]{},
				InputLineSegments: []LineSegment[float64]{
					New[float64](1, 15, 1, 5),
					New[float64](0, 10, 10, 10),
				},
			},
			B: IntersectionResult[float64]{
				IntersectionType:   IntersectionPoint,
				IntersectionPoint:  point.New[float64](1, 10),
				OverlappingSegment: LineSegment[float64]{},
				InputLineSegments: []LineSegment[float64]{
					New[float64](10, 10, 0, 10),
					New[float64](1, 5, 1, 15),
				},
			},
			epsilon:  1e-8,
			expected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.A.Eq(tc.B, options.WithEpsilon(tc.epsilon))
			assert.Equal(t, tc.expected, actual)
		})
	}
}
