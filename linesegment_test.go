package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func BenchmarkLineSegment_DistanceToLineSegment(b *testing.B) {
	segment1 := NewLineSegment(NewPoint(0, 0), NewPoint(10, 10))
	segment2 := NewLineSegment(NewPoint(5, 5), NewPoint(15, 15))

	for i := 0; i < b.N; i++ {
		segment1.DistanceToLineSegment(segment2)
	}
}

func BenchmarkLineSegment_ProjectOntoLineSegment(b *testing.B) {
	segment := NewLineSegment(NewPoint(0, 0), NewPoint(10, 10))
	point := NewPoint(5, 3)

	for i := 0; i < b.N; i++ {
		point.ProjectOntoLineSegment(segment)
	}
}

func TestDetailedLineSegmentRelationship_String(t *testing.T) {
	tests := map[string]struct {
		input       detailedLineSegmentRelationship
		expected    string
		shouldPanic bool
	}{
		"lsrCollinearDisjoint": {
			input:       lsrCollinearDisjoint,
			expected:    "lsrCollinearDisjoint",
			shouldPanic: false,
		},
		"lsrMiss": {
			input:       lsrMiss,
			expected:    "lsrMiss",
			shouldPanic: false,
		},
		"lsrIntersects": {
			input:       lsrIntersects,
			expected:    "lsrIntersects",
			shouldPanic: false,
		},
		"lsrAeqC": {
			input:       lsrAeqC,
			expected:    "lsrAeqC",
			shouldPanic: false,
		},
		"lsrAeqD": {
			input:       lsrAeqD,
			expected:    "lsrAeqD",
			shouldPanic: false,
		},
		"lsrBeqC": {
			input:       lsrBeqC,
			expected:    "lsrBeqC",
			shouldPanic: false,
		},
		"lsrBeqD": {
			input:       lsrBeqD,
			expected:    "lsrBeqD",
			shouldPanic: false,
		},
		"lsrAonCD": {
			input:       lsrAonCD,
			expected:    "lsrAonCD",
			shouldPanic: false,
		},
		"lsrBonCD": {
			input:       lsrBonCD,
			expected:    "lsrBonCD",
			shouldPanic: false,
		},
		"lsrConAB": {
			input:       lsrConAB,
			expected:    "lsrConAB",
			shouldPanic: false,
		},
		"lsrDonAB": {
			input:       lsrDonAB,
			expected:    "lsrDonAB",
			shouldPanic: false,
		},
		"lsrCollinearAonCD": {
			input:       lsrCollinearAonCD,
			expected:    "lsrCollinearAonCD",
			shouldPanic: false,
		},
		"lsrCollinearBonCD": {
			input:       lsrCollinearBonCD,
			expected:    "lsrCollinearBonCD",
			shouldPanic: false,
		},
		"lsrCollinearABinCD": {
			input:       lsrCollinearABinCD,
			expected:    "lsrCollinearABinCD",
			shouldPanic: false,
		},
		"lsrCollinearCDinAB": {
			input:       lsrCollinearCDinAB,
			expected:    "lsrCollinearCDinAB",
			shouldPanic: false,
		},
		"lsrCollinearEqual": {
			input:       lsrCollinearEqual,
			expected:    "lsrCollinearEqual",
			shouldPanic: false,
		},
		"UnsupportedRelationship": {
			input:       detailedLineSegmentRelationship(100), // An unsupported relationship
			shouldPanic: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.shouldPanic {
				require.Panics(t, func() {
					_ = tt.input.String()
				}, "Expected panic for unsupported relationship")
			} else {
				require.NotPanics(t, func() {
					output := tt.input.String()
					assert.Equal(t, tt.expected, output, "Unexpected string for relationship")
				}, "Did not expect panic for supported relationship")
			}
		})
	}
}

func TestLineSegment_AddLineSegment(t *testing.T) {
	tests := map[string]struct {
		segment1 any                  // First line segment (can be int or float64)
		segment2 any                  // Second line segment to add (can be int or float64)
		expected LineSegment[float64] // Expected resulting line segment (float64 type)
	}{
		// Integer line segment test cases
		"int: Add line segment to segment": {
			segment1: NewLineSegment[int](NewPoint(1, 1), NewPoint(4, 5)),
			segment2: NewLineSegment[int](NewPoint(2, 2), NewPoint(3, 3)),
			expected: NewLineSegment[float64](NewPoint[float64](3, 3), NewPoint[float64](7, 8)),
		},

		// Float64 line segment test cases
		"float64: Add line segment to segment": {
			segment1: NewLineSegment[float64](NewPoint(1.5, 2.5), NewPoint(4.0, 5.5)),
			segment2: NewLineSegment[float64](NewPoint(2.0, 1.5), NewPoint(1.0, 3.0)),
			expected: NewLineSegment[float64](NewPoint[float64](3.5, 4.0), NewPoint[float64](5.0, 8.5)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch segment1 := tt.segment1.(type) {
			case LineSegment[int]:
				segment2 := tt.segment2.(LineSegment[int])
				result := segment1.AddLineSegment(segment2)
				assert.InDelta(t, tt.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tt.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tt.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tt.expected.end.y, result.end.y, 0.001)

			case LineSegment[float64]:
				segment2 := tt.segment2.(LineSegment[float64])
				result := segment1.AddLineSegment(segment2)
				assert.InDelta(t, tt.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tt.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tt.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tt.expected.end.y, result.end.y, 0.001)
			}
		})
	}
}

func TestLineSegment_Translate(t *testing.T) {
	tests := map[string]struct {
		segment  any                  // Original line segment (can be int or float64)
		vector   any                  // Vector to add (can be int or float64)
		expected LineSegment[float64] // Expected resulting line segment (float64 type)
	}{
		// Integer vector test cases
		"int: Add vector to segment": {
			segment:  NewLineSegment[int](NewPoint(1, 1), NewPoint(4, 5)),
			vector:   NewPoint[int](3, 3),
			expected: NewLineSegment[float64](NewPoint[float64](4, 4), NewPoint[float64](7, 8)),
		},

		// Float64 vector test cases
		"float64: Add vector to segment": {
			segment:  NewLineSegment[float64](NewPoint(1.5, 2.5), NewPoint(4.0, 5.5)),
			vector:   NewPoint[float64](2.0, 3.0),
			expected: NewLineSegment[float64](NewPoint[float64](3.5, 5.5), NewPoint[float64](6.0, 8.5)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			switch segment := tc.segment.(type) {
			case LineSegment[int]:
				vec := tc.vector.(Point[int])
				result := segment.Translate(vec)
				assert.InDelta(t, tc.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tc.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tc.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tc.expected.end.y, result.end.y, 0.001)

			case LineSegment[float64]:
				vec := tc.vector.(Point[float64])
				result := segment.Translate(vec)
				assert.InDelta(t, tc.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tc.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tc.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tc.expected.end.y, result.end.y, 0.001)
			}
		})
	}
}

func TestLineSegment_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		point    any
		segment  any
		expected bool
	}{
		"Point on line segment (float64)": {
			point:    NewPoint[float64](1, 1),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point at endpoint A (float64)": {
			point:    NewPoint[float64](0, 0),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point at endpoint End (float64)": {
			point:    NewPoint[float64](2, 2),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point collinear but outside bounding box (float64)": {
			point:    NewPoint[float64](3, 3),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: false,
		},
		"Point not collinear (float64)": {
			point:    NewPoint[float64](1, 2),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: false,
		},
		"Point on line segment (int)": {
			point:    NewPoint[int](1, 1),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point at endpoint A (int)": {
			point:    NewPoint[int](0, 0),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point at endpoint End (int)": {
			point:    NewPoint[int](2, 2),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point collinear but outside bounding box (int)": {
			point:    NewPoint[int](3, 3),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: false,
		},
		"Point not collinear (int)": {
			point:    NewPoint[int](1, 2),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				segment := tt.segment.(LineSegment[int])
				result := segment.ContainsPoint(point)
				assert.Equal(t, tt.expected, result, "Test %s failed", name)
			case Point[float64]:
				segment := tt.segment.(LineSegment[float64])
				result := segment.ContainsPoint(point)
				assert.Equal(t, tt.expected, result, "Test %s failed", name)
			default:
				t.Errorf("Unsupported point type in test %s", name)
			}
		})
	}
}

func TestLineSegment_detailedRelationshipToLineSegment(t *testing.T) {
	tests := map[string]struct {
		AB, CD   any                             // Supports both LineSegment[int] and LineSegment[float64]
		expected detailedLineSegmentRelationship // Expected result
	}{
		// Disjoint cases
		"Disjoint non-collinear (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(3, -3), NewPoint(5, -5)),
			expected: lsrMiss,
		},
		"Disjoint collinear (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(3, 3), NewPoint(4, 4)),
			expected: lsrCollinearDisjoint,
		},

		// Intersection cases
		"Intersecting at unique point (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(4, 4)),
			CD:       NewLineSegment(NewPoint(0, 4), NewPoint(4, 0)),
			expected: lsrIntersects,
		},

		// Endpoint coincidences
		"Endpoint A equals C (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(0, 0), NewPoint(2, -2)),
			expected: lsrAeqC,
		},
		"Endpoint A equals D (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(2, -2), NewPoint(0, 0)),
			expected: lsrAeqD,
		},
		"Endpoint End equals C (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(3, 3)),
			CD:       NewLineSegment(NewPoint(3, 3), NewPoint(2, 0)),
			expected: lsrBeqC,
		},
		"Endpoint End equals D (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(3, 3)),
			CD:       NewLineSegment(NewPoint(2, 0), NewPoint(3, 3)),
			expected: lsrBeqD,
		},

		// Endpoint-on-segment cases (non-collinear)
		"A on CD without collinearity (int)": {
			AB:       NewLineSegment(NewPoint(0, 10), NewPoint(0, 0)),
			CD:       NewLineSegment(NewPoint(-10, 10), NewPoint(10, 10)),
			expected: lsrAonCD,
		},
		"End on CD without collinearity (int)": {
			AB:       NewLineSegment(NewPoint(2, 2), NewPoint(3, 1)),
			CD:       NewLineSegment(NewPoint(1, 1), NewPoint(4, 1)),
			expected: lsrBonCD,
		},
		"C on AB without collinearity (int)": {
			AB:       NewLineSegment(NewPoint(-10, 10), NewPoint(10, 10)),
			CD:       NewLineSegment(NewPoint(0, 10), NewPoint(0, 0)),
			expected: lsrConAB,
		},
		"D on AB without collinearity (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(4, 1)),
			CD:       NewLineSegment(NewPoint(2, 2), NewPoint(3, 1)),
			expected: lsrDonAB,
		},

		// Collinear partial overlaps
		"A on CD with collinearity (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(4, 4)),
			CD:       NewLineSegment(NewPoint(0, 0), NewPoint(3, 3)),
			expected: lsrCollinearAonCD,
		},
		"End on CD with collinearity (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(3, 3)),
			CD:       NewLineSegment(NewPoint(1, 1), NewPoint(4, 4)),
			expected: lsrCollinearBonCD,
		},

		// Full containment
		"AB fully within CD (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(0, 0), NewPoint(3, 3)),
			expected: lsrCollinearABinCD,
		},
		"CD fully within AB (int)": {
			AB:       NewLineSegment(NewPoint(0, 0), NewPoint(4, 4)),
			CD:       NewLineSegment(NewPoint(1, 1), NewPoint(2, 2)),
			expected: lsrCollinearCDinAB,
		},

		// Exact equality
		"Segments are exactly equal (int)": {
			AB:       NewLineSegment(NewPoint(1, 1), NewPoint(2, 2)),
			CD:       NewLineSegment(NewPoint(1, 1), NewPoint(2, 2)),
			expected: lsrCollinearEqual,
		},

		// Disjoint cases
		"Disjoint non-collinear (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			CD:       NewLineSegment(NewPoint[float64](3, -3), NewPoint[float64](5, -5)),
			expected: lsrMiss,
		},
		"Disjoint collinear (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			CD:       NewLineSegment(NewPoint[float64](3, 3), NewPoint[float64](4, 4)),
			expected: lsrCollinearDisjoint,
		},

		// Intersection cases
		"Intersecting at unique point (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](4, 4)),
			CD:       NewLineSegment(NewPoint[float64](0, 4), NewPoint[float64](4, 0)),
			expected: lsrIntersects,
		},

		// Endpoint coincidences
		"Endpoint A equals C (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			CD:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, -2)),
			expected: lsrAeqC,
		},
		"Endpoint End equals D (float64)": {
			AB:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](3, 3)),
			CD:       NewLineSegment(NewPoint[float64](2, 0), NewPoint[float64](3, 3)),
			expected: lsrBeqD,
		},

		// Endpoint-on-segment cases (non-collinear)
		"A on CD without collinearity (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 10), NewPoint[float64](0, 0)),
			CD:       NewLineSegment(NewPoint[float64](-10, 10), NewPoint[float64](10, 10)),
			expected: lsrAonCD,
		},
		"End on CD without collinearity (float64)": {
			AB:       NewLineSegment(NewPoint[float64](2, 2), NewPoint[float64](3, 1)),
			CD:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](4, 1)),
			expected: lsrBonCD,
		},

		// Collinear partial overlaps
		"A on CD with collinearity (float64)": {
			AB:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](4, 4)),
			CD:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](3, 3)),
			expected: lsrCollinearAonCD,
		},
		"End on CD with collinearity (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](3, 3)),
			CD:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](4, 4)),
			expected: lsrCollinearBonCD,
		},

		// Full containment
		"AB fully within CD (float64)": {
			AB:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](2, 2)),
			CD:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](3, 3)),
			expected: lsrCollinearABinCD,
		},
		"CD fully within AB (float64)": {
			AB:       NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](4, 4)),
			CD:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](2, 2)),
			expected: lsrCollinearCDinAB,
		},

		// Exact equality
		"Segments are exactly equal (float64)": {
			AB:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](2, 2)),
			CD:       NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](2, 2)),
			expected: lsrCollinearEqual,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch ab := tt.AB.(type) {
			case LineSegment[int]:
				cd := tt.CD.(LineSegment[int])
				result := ab.detailedRelationshipToLineSegment(cd)
				assert.Equal(t, tt.expected, result, "Test %s failed for LineSegment[int]", name)
			case LineSegment[float64]:
				cd := tt.CD.(LineSegment[float64])
				result := ab.detailedRelationshipToLineSegment(cd)
				assert.Equal(t, tt.expected, result, "Test %s failed for LineSegment[float64]", name)
			default:
				t.Errorf("Unsupported segment type in test %s", name)
			}
		})
	}
}

func TestLineSegment_DistanceToLineSegment(t *testing.T) {
	tests := map[string]struct {
		AB, CD           any // line segments
		expectedDistance float64
	}{
		"parallel vertical lines (int)": {
			AB:               NewLineSegment(NewPoint[int](10, 20), NewPoint[int](10, 30)),
			CD:               NewLineSegment(NewPoint[int](20, 20), NewPoint[int](20, 30)),
			expectedDistance: 10,
		},
		"parallel vertical lines (float64)": {
			AB:               NewLineSegment(NewPoint[float64](10, 20), NewPoint[float64](10, 30)),
			CD:               NewLineSegment(NewPoint[float64](20, 20), NewPoint[float64](20, 30)),
			expectedDistance: 10,
		},
		"parallel horizontal lines (int)": {
			AB:               NewLineSegment(NewPoint[int](-20, 10), NewPoint[int](-10, -10)),
			CD:               NewLineSegment(NewPoint[int](-20, -20), NewPoint[int](-10, -20)),
			expectedDistance: 10,
		},
		"parallel horizontal lines (float64)": {
			AB:               NewLineSegment(NewPoint[float64](-20, 10), NewPoint[float64](-10, -10)),
			CD:               NewLineSegment(NewPoint[float64](-20, -20), NewPoint[float64](-10, -20)),
			expectedDistance: 10,
		},
		"parallel diagonal lines (int)": {
			AB:               NewLineSegment(NewPoint[int](-10, 0), NewPoint[int](0, -10)),
			CD:               NewLineSegment(NewPoint[int](0, 10), NewPoint[int](10, 0)),
			expectedDistance: 14.1421,
		},
		"parallel diagonal lines (float64)": {
			AB:               NewLineSegment(NewPoint[float64](-10, 0), NewPoint[float64](0, -10)),
			CD:               NewLineSegment(NewPoint[float64](0, 10), NewPoint[float64](10, 0)),
			expectedDistance: 14.1421,
		},
		"perpendicular lines (int)": {
			AB:               NewLineSegment(NewPoint[int](0, 0), NewPoint[int](0, 10)),
			CD:               NewLineSegment(NewPoint[int](-10, 20), NewPoint[int](10, 20)),
			expectedDistance: 10,
		},
		"perpendicular lines (float64)": {
			AB:               NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](0, 10)),
			CD:               NewLineSegment(NewPoint[float64](-10, 20), NewPoint[float64](10, 20)),
			expectedDistance: 10,
		},
		"non-intersecting, non-parallel, non-perpendicular (int)": {
			AB:               NewLineSegment(NewPoint[int](-15, 28), NewPoint[int](16, 18)),
			CD:               NewLineSegment(NewPoint[int](-7, 21), NewPoint[int](0, 8)),
			expectedDistance: 4.2059,
		},
		"non-intersecting, non-parallel, non-perpendicular (float64)": {
			AB:               NewLineSegment(NewPoint[float64](-15, 28), NewPoint[float64](16, 18)),
			CD:               NewLineSegment(NewPoint[float64](-7, 21), NewPoint[float64](0, 8)),
			expectedDistance: 4.2059,
		},
		"intersecting, oblique (float64)": {
			AB:               NewLineSegment(NewPoint[float64](-13, 19), NewPoint[float64](12, 23)),
			CD:               NewLineSegment(NewPoint[float64](-12, 9), NewPoint[float64](7, 26)),
			expectedDistance: 0,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch AB := tt.AB.(type) {
			case LineSegment[int]:
				CD := tt.CD.(LineSegment[int])
				actual := AB.DistanceToLineSegment(CD)
				assert.InDelta(t, tt.expectedDistance, actual, 0.0001)
			case LineSegment[float64]:
				CD := tt.CD.(LineSegment[float64])
				actual := AB.DistanceToLineSegment(CD)
				assert.InDelta(t, tt.expectedDistance, actual, 0.0001)
			}
		})
	}
}

func TestLineSegment_DistanceToPoint(t *testing.T) {
	tests := map[string]struct {
		point    any     // Point to be projected (can be int or float64)
		segment  any     // Line segment for projection (can be int or float64)
		expected float64 // Expected distance
	}{
		// Integer points test cases
		"Project onto segment from inside (int)": {
			point:    NewPoint(5, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 0.0, // Point is on the segment
		},
		"Project onto segment from above (int)": {
			point:    NewPoint(4, 6),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.386,
		},
		"Project onto segment from below (int)": {
			point:    NewPoint(4, 2),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.941,
		},
		"Project off the start of segment (int)": {
			point:    NewPoint(0, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},
		"Project off the end of segment (int)": {
			point:    NewPoint(10, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},

		// Float64 points test cases
		"Project onto segment from inside (float64)": {
			point:    NewPoint(5.5, 5.5),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 0.1387,
		},
		"Project onto segment from above (float64)": {
			point:    NewPoint(4.0, 6.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.386,
		},
		"Project onto segment from below (float64)": {
			point:    NewPoint(4.0, 2.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.941,
		},
		"Project off the start of segment (float64)": {
			point:    NewPoint(0.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
		"Project off the end of segment (float64)": {
			point:    NewPoint(10.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := seg.DistanceToPoint(p)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for int points")

			case Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := seg.DistanceToPoint(p)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for float64 points")
			}
		})
	}
}

func TestLineSegment_Length(t *testing.T) {
	tests := map[string]struct {
		start, end     any
		expectedLength float64
	}{
		"int: start: (0,0), end: (10,0)": {
			start:          NewPoint[int](0, 0),
			end:            NewPoint[int](10, 0),
			expectedLength: 10,
		},
		"float64: start: (0,0), end: (0,10)": {
			start:          NewPoint[float64](0, 0),
			end:            NewPoint[float64](0, 10),
			expectedLength: 10,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case Point[int]:
				end := tt.end.(Point[int])
				ls := NewLineSegment(start, end)
				assert.Equal(t, tt.expectedLength, ls.Length())
			case Point[float64]:
				end := tt.end.(Point[float64])
				ls := NewLineSegment(start, end)
				assert.Equal(t, tt.expectedLength, ls.Length())
			}
		})
	}
}

func TestLineSegment_Midpoint(t *testing.T) {
	tests := map[string]struct {
		start, end       any
		expectedMidpoint Point[float64]
	}{
		"int: start: (1,2), end: (3,4)": {
			start:            NewPoint[int](0, 0),
			end:              NewPoint[int](10, 10),
			expectedMidpoint: NewPoint[float64](5, 5),
		},
		"float64: start: (1,2), end: (3,4)": {
			start:            NewPoint[float64](0, 0),
			end:              NewPoint[float64](10, 10),
			expectedMidpoint: NewPoint[float64](5, 5),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case Point[int]:
				end := tt.end.(Point[int])
				ls := NewLineSegment(start, end)
				assert.Equal(t, tt.expectedMidpoint, ls.Center())
			case Point[float64]:
				end := tt.end.(Point[float64])
				ls := NewLineSegment(start, end)
				assert.Equal(t, tt.expectedMidpoint, ls.Center())
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	cases := map[string]struct {
		input    LineSegment[int]
		expected LineSegment[int]
	}{
		"already normalized": {
			input:    NewLineSegment(NewPoint(1, 2), NewPoint(3, 5)),
			expected: NewLineSegment(NewPoint(1, 2), NewPoint(3, 5)),
		},
		"needs normalization": {
			input:    NewLineSegment(NewPoint(3, 5), NewPoint(1, 2)),
			expected: NewLineSegment(NewPoint(1, 2), NewPoint(3, 5)),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := tc.input.Normalize()
			if !result.Eq(tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestLineSegment_Reflect(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment[float64]   // Original line segment
		axis     ReflectionAxis         // Axis or line for reflection
		line     []LineSegment[float64] // Custom line segment (optional, can be empty)
		expected LineSegment[float64]   // Expected reflected line segment
	}{
		"reflect across x-axis": {
			segment:  NewLineSegment(NewPoint[float64](2, 3), NewPoint[float64](4, 5)),
			axis:     ReflectAcrossXAxis,
			expected: NewLineSegment(NewPoint[float64](2, -3), NewPoint[float64](4, -5)),
		},
		"reflect across y-axis": {
			segment:  NewLineSegment(NewPoint[float64](2, 3), NewPoint[float64](4, 5)),
			axis:     ReflectAcrossYAxis,
			expected: NewLineSegment(NewPoint[float64](-2, 3), NewPoint[float64](-4, 5)),
		},
		"reflect across y = x line (ReflectAcrossCustomLine)": {
			segment:  NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),
			axis:     ReflectAcrossCustomLine,
			line:     []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](1, 1))},
			expected: NewLineSegment(NewPoint[float64](4, 3), NewPoint[float64](7, 6)),
		},
		"reflect across y = -x line (ReflectAcrossCustomLine)": {
			segment:  NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),
			axis:     ReflectAcrossCustomLine,
			line:     []LineSegment[float64]{NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](-1, 1))},
			expected: NewLineSegment(NewPoint[float64](-4, -3), NewPoint[float64](-7, -6)),
		},
		"reflect across degenerate line (ReflectAcrossCustomLine)": {
			segment:  NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),
			axis:     ReflectAcrossCustomLine,
			line:     []LineSegment[float64]{NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](1, 1))}, // Degenerate line
			expected: NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),                         // Expect unchanged segment
		},
		"no custom line provided (ReflectAcrossCustomLine)": {
			segment:  NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),
			axis:     ReflectAcrossCustomLine,
			line:     nil,                                                              // No custom line provided
			expected: NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)), // Expect unchanged segment
		},
		"invalid axis": {
			segment:  NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)),
			axis:     ReflectionAxis(999),                                              // Invalid axis
			expected: NewLineSegment(NewPoint[float64](3, 4), NewPoint[float64](6, 7)), // Expect unchanged segment
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var result LineSegment[float64]
			if tt.axis == ReflectAcrossCustomLine && len(tt.line) > 0 {
				result = tt.segment.Reflect(ReflectAcrossCustomLine, tt.line[0])
			} else if tt.axis == ReflectAcrossCustomLine {
				result = tt.segment.Reflect(ReflectAcrossCustomLine)
			} else {
				result = tt.segment.Reflect(tt.axis)
			}

			assert.InDelta(t, tt.expected.start.x, result.start.x, 0.001)
			assert.InDelta(t, tt.expected.start.y, result.start.y, 0.001)
			assert.InDelta(t, tt.expected.end.x, result.end.x, 0.001)
			assert.InDelta(t, tt.expected.end.y, result.end.y, 0.001)
		})
	}
}

func TestLineSegment_RelationshipToCircle(t *testing.T) {
	// Define a circle
	circle := NewCircle(NewPoint[float64](5, 5), 5.0)

	// Test cases
	tests := map[string]struct {
		line     LineSegment[float64]
		expected Relationship
	}{
		"Disjoint": {
			line:     NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](-2, -2)),
			expected: RelationshipDisjoint,
		},
		"Intersecting at one endpoint": {
			line:     NewLineSegment(NewPoint[float64](5, 0), NewPoint[float64](10, 5)),
			expected: RelationshipIntersection,
		},
		"Intersecting along segment": {
			line:     NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			expected: RelationshipIntersection,
		},
		"Contained": {
			line:     NewLineSegment(NewPoint[float64](5, 6), NewPoint[float64](5, 4)),
			expected: RelationshipContainedBy,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.line.RelationshipToCircle(circle)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLineSegment_RelationshipToLineSegment(t *testing.T) {
	// Test cases
	tests := map[string]struct {
		line1    LineSegment[float64]
		line2    LineSegment[float64]
		expected Relationship
	}{
		"Disjoint segments": {
			line1:    NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			line2:    NewLineSegment(NewPoint[float64](20, 20), NewPoint[float64](30, 30)),
			expected: RelationshipDisjoint,
		},
		"Intersecting segments": {
			line1:    NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			line2:    NewLineSegment(NewPoint[float64](5, 5), NewPoint[float64](15, 15)),
			expected: RelationshipIntersection,
		},
		"Equal segments": {
			line1:    NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			line2:    NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			expected: RelationshipEqual,
		},
		"Crossing segments": {
			line1:    NewLineSegment(NewPoint[float64](0, 10), NewPoint[float64](10, 0)),
			line2:    NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](10, 10)),
			expected: RelationshipIntersection,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.line1.RelationshipToLineSegment(tc.line2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLineSegment_RelationshipToPoint(t *testing.T) {
	// Define a line segment
	segment := NewLineSegment(
		NewPoint[float64](0, 0),
		NewPoint[float64](10, 10),
	)

	// Test cases
	tests := map[string]struct {
		point    Point[float64]
		expected Relationship
	}{
		"Point on the segment": {
			point:    NewPoint[float64](5, 5),
			expected: RelationshipIntersection,
		},
		"Point disjoint from the segment": {
			point:    NewPoint[float64](10, 0),
			expected: RelationshipDisjoint,
		},
		"Point coinciding with an endpoint": {
			point:    NewPoint[float64](0, 0),
			expected: RelationshipIntersection,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := segment.RelationshipToPoint(tc.point)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLineSegment_RelationshipToRectangle(t *testing.T) {
	// Define a rectangle
	rect := NewRectangle([]Point[float64]{
		NewPoint[float64](0, 0),
		NewPoint[float64](10, 0),
		NewPoint[float64](10, 10),
		NewPoint[float64](0, 10),
	})

	// Test cases
	tests := map[string]struct {
		line     LineSegment[float64]
		expected Relationship
	}{
		"Intersects": {
			line:     NewLineSegment(NewPoint[float64](5, 5), NewPoint[float64](15, 15)),
			expected: RelationshipIntersection,
		},
		"Contained": {
			line:     NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](9, 9)),
			expected: RelationshipContainedBy,
		},
		"Disjoint": {
			line:     NewLineSegment(NewPoint[float64](20, 20), NewPoint[float64](30, 30)),
			expected: RelationshipDisjoint,
		},
		"Touches edge but does not intersect": {
			line:     NewLineSegment(NewPoint[float64](10, 10), NewPoint[float64](15, 15)),
			expected: RelationshipIntersection,
		},
		"Touches vertex only": {
			line:     NewLineSegment(NewPoint[float64](10, 10), NewPoint[float64](10, 10)),
			expected: RelationshipIntersection,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.line.RelationshipToRectangle(rect)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLineSegment_RoundToEpsilon(t *testing.T) {
	tests := map[string]struct {
		input    LineSegment[float64]
		epsilon  float64
		expected LineSegment[float64]
	}{
		"standard rounding": {
			input:    NewLineSegment(NewPoint[float64](1.23, 4.56), NewPoint[float64](7.89, 3.21)),
			epsilon:  0.1,
			expected: NewLineSegment(NewPoint[float64](1.2, 4.6), NewPoint[float64](7.9, 3.2)),
		},
		"no change for exact multiples": {
			input:    NewLineSegment(NewPoint[float64](1.0, 4.0), NewPoint[float64](7.0, 3.0)),
			epsilon:  0.5,
			expected: NewLineSegment(NewPoint[float64](1.0, 4.0), NewPoint[float64](7.0, 3.0)),
		},
		"small epsilon": {
			input:    NewLineSegment(NewPoint[float64](1.2345, 4.5678), NewPoint[float64](7.8912, 3.2109)),
			epsilon:  0.01,
			expected: NewLineSegment(NewPoint[float64](1.23, 4.57), NewPoint[float64](7.89, 3.21)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.input.RoundToEpsilon(tc.epsilon)
			assert.InDelta(t, tc.expected.start.x, result.start.x, tc.epsilon, "unexpected start x value")
			assert.InDelta(t, tc.expected.start.y, result.start.y, tc.epsilon, "unexpected start y value")
			assert.InDelta(t, tc.expected.end.x, result.end.x, tc.epsilon, "unexpected start x value")
			assert.InDelta(t, tc.expected.end.y, result.end.y, tc.epsilon, "unexpected start y value")
		})
	}
}

func TestLineSegment_SubLineSegment(t *testing.T) {
	tests := map[string]struct {
		segment1 any                  // First line segment (can be int or float64)
		segment2 any                  // Second line segment to subtract (can be int or float64)
		expected LineSegment[float64] // Expected resulting line segment (float64 type)
	}{
		// Integer line segment test cases
		"int: Subtract line segment from segment": {
			segment1: NewLineSegment[int](NewPoint(5, 5), NewPoint(8, 10)),
			segment2: NewLineSegment[int](NewPoint(2, 2), NewPoint(4, 3)),
			expected: NewLineSegment[float64](NewPoint[float64](3, 3), NewPoint[float64](4, 7)),
		},

		// Float64 line segment test cases
		"float64: Subtract line segment from segment": {
			segment1: NewLineSegment[float64](NewPoint(5.5, 6.5), NewPoint(8.0, 10.5)),
			segment2: NewLineSegment[float64](NewPoint(1.5, 2.5), NewPoint(2.0, 3.0)),
			expected: NewLineSegment[float64](NewPoint[float64](4.0, 4.0), NewPoint[float64](6.0, 7.5)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch segment1 := tt.segment1.(type) {
			case LineSegment[int]:
				segment2 := tt.segment2.(LineSegment[int])
				result := segment1.SubLineSegment(segment2)
				assert.InDelta(t, tt.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tt.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tt.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tt.expected.end.y, result.end.y, 0.001)

			case LineSegment[float64]:
				segment2 := tt.segment2.(LineSegment[float64])
				result := segment1.SubLineSegment(segment2)
				assert.InDelta(t, tt.expected.start.x, result.start.x, 0.001)
				assert.InDelta(t, tt.expected.start.y, result.start.y, 0.001)
				assert.InDelta(t, tt.expected.end.x, result.end.x, 0.001)
				assert.InDelta(t, tt.expected.end.y, result.end.y, 0.001)
			}
		})
	}
}

func TestLineSegment_XAtY(t *testing.T) {
	seg := NewLineSegment(NewPoint(2, 2), NewPoint(8, 5))

	// Test diagonal line
	x, ok := seg.XAtY(3)
	require.True(t, ok, "expected valid X for Y=3")
	assert.Equal(t, 4.0, x)

	// Test out-of-bounds Y
	_, ok = seg.XAtY(6)
	require.False(t, ok, "expected out-of-bounds for Y=6")

	// Test vertical line
	seg = NewLineSegment(NewPoint(5, 1), NewPoint(5, 4))
	x, ok = seg.XAtY(3)
	require.True(t, ok, "expected valid X for vertical line")
	assert.Equal(t, 5.0, x)

	// Test horizontal line
	seg = NewLineSegment(NewPoint(2, 3), NewPoint(7, 3))
	x, ok = seg.XAtY(3)
	require.True(t, ok, "expected valid X for horizontal line")
	assert.Equal(t, 2.0, x) // Can test any X in the range
}

func TestLineSegment_YAtX(t *testing.T) {
	seg := NewLineSegment(NewPoint(2, 2), NewPoint(8, 5))

	// Test diagonal line
	y, ok := seg.YAtX(4)
	require.True(t, ok, "expected valid Y for X=4")
	assert.Equal(t, 3.0, y)

	// Test out-of-bounds X
	_, ok = seg.YAtX(9)
	require.False(t, ok, "expected out-of-bounds for X=9")

	// Test vertical line
	seg = NewLineSegment(NewPoint(5, 1), NewPoint(5, 4))
	y, ok = seg.YAtX(5)
	require.True(t, ok, "expected valid Y for vertical line")
	assert.Equal(t, 1.0, y) // Can test for other Y values in the range

	// Test horizontal line
	seg = NewLineSegment(NewPoint(2, 3), NewPoint(7, 3))
	y, ok = seg.YAtX(5)
	require.True(t, ok, "expected valid Y for horizontal line")
	assert.Equal(t, 3.0, y)
}

func TestNewLineSegment(t *testing.T) {
	tests := map[string]struct {
		start, end any
		expected   any
	}{
		"int: start: (1,2), end: (3,4)": {
			start: NewPoint[int](1, 2),
			end:   NewPoint[int](3, 4),
			expected: LineSegment[int]{
				start: Point[int]{x: 1, y: 2},
				end:   Point[int]{x: 3, y: 4},
			},
		},
		"float64: start: (1,2), end: (3,4)": {
			start: NewPoint[float64](1, 2),
			end:   NewPoint[float64](3, 4),
			expected: LineSegment[float64]{
				start: Point[float64]{x: 1, y: 2},
				end:   Point[float64]{x: 3, y: 4},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case Point[int]:
				end := tt.end.(Point[int])
				expected := tt.expected.(LineSegment[int])
				actual := NewLineSegment(start, end)
				assert.Equal(t, expected, actual)
			case Point[float64]:
				end := tt.end.(Point[float64])
				expected := tt.expected.(LineSegment[float64])
				actual := NewLineSegment(start, end)
				assert.Equal(t, expected, actual)
			}
		})
	}
}
