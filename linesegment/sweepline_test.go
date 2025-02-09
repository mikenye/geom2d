package linesegment

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func compareIntersectionResult[T types.SignedNumber](resultA, resultB IntersectionResult[T], epsilon float64) (bool, string) {

	// compare IntersectionType
	if resultA.IntersectionType != resultB.IntersectionType {
		return false, "IntersectionType mismatch"
	}

	// compare point
	if resultA.IntersectionType == IntersectionPoint {
		if !resultA.IntersectionPoint.Eq(resultB.IntersectionPoint, options.WithEpsilon(epsilon)) {
			return false, "IntersectionPoint mismatch"
		}
	}

	// compare overlapping segment
	if resultA.IntersectionType == IntersectionOverlappingSegment {
		if !resultA.OverlappingSegment.Eq(resultB.OverlappingSegment, options.WithEpsilon(epsilon)) {
			return false, "OverlappingSegment mismatch"
		}
	}

	// compare input segments
	for _, segA := range resultA.InputLineSegments {
		found := false
		for _, segB := range resultB.InputLineSegments {
			if segA.Eq(segB, options.WithEpsilon(epsilon)) {
				found = true
				break
			}
		}
		if !found {
			return false, "InputLineSegments mismatch"
		}
	}

	return true, ""
}

func compareIntersectionResults[T types.SignedNumber](A, B []IntersectionResult[T], epsilon float64) (bool, string) {
	// length check
	if len(A) != len(B) {
		return false, "result slice length mismatch"
	}

	// look for matches
	for _, resultA := range A {
		found := false
		for _, resultB := range B {
			found, _ = compareIntersectionResult[T](resultA, resultB, epsilon)
			if found {
				break
			}
		}
		if !found {
			return false, "missing result"
		}
	}
	return true, ""
}

func FuzzFindIntersections_Int_2Segments(f *testing.F) {
	// Seed with sample inputs
	f.Add(0, 0, 10, 10, 5, 5, 15, 15) // Diagonal overlap
	f.Add(0, 0, 10, 0, 5, 0, 15, 0)   // Horizontal overlap
	f.Add(0, 0, 0, 10, 5, 0, 15, 0)   // Vertical overlap
	f.Add(0, 5, 10, 5, 5, 0, 5, 10)   // "+" shape
	f.Add(0, 0, 10, 10, 0, 10, 10, 0) // "X" shape
	f.Add(0, 10, 0, 0, 0, 0, 10, 0)   // "L" shape
	f.Add(4, 7, 5, 5, 5, 10, 4, 0)    // Lines cross, steep slope
	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, x3, y3, x4, y4 int) {
		// Ensure valid segments
		if x1 == x2 && y1 == y2 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x3 == x4 && y3 == y4 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}

		segments := []LineSegment[int]{
			New(x1, y1, x2, y2),
			New(x3, y3, x4, y4),
		}

		naiveResults := FindIntersectionsSlow(segments, options.WithEpsilon(1e-8))
		sweepResults := FindIntersectionsFast(segments, options.WithEpsilon(1e-8))

		ok, reason := compareIntersectionResults(naiveResults, sweepResults, 1e-8)

		if !ok {
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive:      %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
		}
	})
}

func FuzzFindIntersections_Int_3Segments(f *testing.F) {
	// Seed with sample inputs
	f.Add(0, 0, 5, 10, 5, 10, 10, 0, 10, 0, 0, 0)   // triangle
	f.Add(0, 8, 10, 8, 0, 3, 10, 3, 1, 0, 9, 10)    // not equals shape ("≠")
	f.Add(3, 6, 7, 6, 3, 8, 7, 8, 5, 10, 5, 6)      // plus-minus shape ("±")
	f.Add(0, 10, 10, 10, 10, 10, 0, 0, 0, 0, 10, 0) // "Z" shape
	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, x3, y3, x4, y4, x5, y5, x6, y6 int) {
		// Ensure valid segments
		if x1 == x2 && y1 == y2 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x3 == x4 && y3 == y4 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x5 == x6 && y5 == y6 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}

		segments := []LineSegment[int]{
			New(x1, y1, x2, y2),
			New(x3, y3, x4, y4),
			New(x5, y5, x6, y6),
		}

		naiveResults := FindIntersectionsSlow(segments, options.WithEpsilon(1e-8))
		sweepResults := FindIntersectionsFast(segments, options.WithEpsilon(1e-8))

		ok, reason := compareIntersectionResults(naiveResults, sweepResults, 1e-8)

		if !ok {
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive:      %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
		}
	})
}

func FuzzFindIntersections_Int_4Segments(f *testing.F) {
	// Seed with sample inputs
	f.Add(0, 5, 5, 10, 5, 10, 10, 5, 10, 5, 5, 0, 5, 0, 0, 5)     // diamond shape
	f.Add(0, 0, 10, 0, 10, 0, 10, 10, 10, 10, 0, 10, 10, 0, 0, 0) // square shape
	f.Add(0, 0, 5, 10, 5, 10, 10, 0, 0, 10, 5, 0, 5, 0, 10, 10)   // crisscrossing W shape
	f.Add(0, 7, 10, 7, 0, 3, 10, 3, 3, 10, 3, 0, 7, 10, 7, 0)     // octothorpe shape ("#")
	f.Fuzz(func(t *testing.T, x1, y1, x2, y2, x3, y3, x4, y4, x5, y5, x6, y6, x7, y7, x8, y8 int) {
		// Ensure valid segments
		if x1 == x2 && y1 == y2 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x3 == x4 && y3 == y4 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x5 == x6 && y5 == y6 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}
		if x7 == x8 && y7 == y8 {
			return // skip degenerate (don't use t.Skip() or fuzz will store the test
		}

		segments := []LineSegment[int]{
			New(x1, y1, x2, y2),
			New(x3, y3, x4, y4),
			New(x5, y5, x6, y6),
			New(x7, y7, x8, y8),
		}

		naiveResults := FindIntersectionsSlow(segments, options.WithEpsilon(1e-8))
		sweepResults := FindIntersectionsFast(segments, options.WithEpsilon(1e-8))

		ok, reason := compareIntersectionResults(naiveResults, sweepResults, 1e-8)

		if !ok {
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive:      %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
		}
	})
}

// TestFindIntersectionsFast ensures that the output of the sweep line algorithm (FindIntersectionsFast) matches
// the output of the naïve algorithm (FindIntersectionsSlow).
func TestFindIntersectionsFast(t *testing.T) {
	tests := map[string]struct {
		segments []LineSegment[int]
	}{
		"parallel non-intersecting segments": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 5),
				New[int](0, 1, 5, 6),
			},
		},
		"simple X-shape intersection": {
			segments: []LineSegment[int]{
				New[int](0, 5, 5, 0),
				New[int](0, 0, 5, 5),
			},
		},
		"horizontal and vertical lines": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0), // Horizontal segment
				New[int](5, -5, 5, 5), // Vertical segment
			},
		},
		"diagonal and horizontal lines": {
			segments: []LineSegment[int]{
				New[int](0, 0, 4, 4), // Diagonal line
				New[int](2, 4, 6, 4), // Horizontal line
			},
		},
		"overlapping diagonal segments": {
			segments: []LineSegment[int]{
				New[int](1, 1, 5, 5), // Segment 1
				New[int](3, 3, 7, 7), // Segment 2
			},
		},
		"overlapping horizontal segments": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0), // Segment 1
				New[int](2, 0, 8, 0),  // Segment 2
			},
		},
		"overlapping vertical segments": {
			segments: []LineSegment[int]{
				New[int](0, 0, 0, 10), // Segment 1
				New[int](0, 2, 0, 8),  // Segment 2
			},
		},
		"coincident segments": { // full overlap
			segments: []LineSegment[int]{
				New[int](1, 1, 5, 5),
				New[int](1, 1, 5, 5),
			},
		},
		"shared endpoint": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 5),
				New[int](5, 5, 10, 0),
			},
		},
		"square shape": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0),
				New[int](10, 0, 10, 10),
				New[int](10, 10, 0, 10),
				New[int](0, 10, 0, 0),
			},
		},
		"diamond shape": {
			segments: []LineSegment[int]{
				New[int](0, 5, 5, 10),
				New[int](5, 10, 10, 5),
				New[int](10, 5, 5, 0),
				New[int](5, 0, 0, 5),
			},
		},
		"t-intersection": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0), // Horizontal segment
				New[int](5, -5, 5, 0), // Vertical segment terminating at (5, 0)
			},
		},
		"t-intersection, rotated 90 deg": {
			segments: []LineSegment[int]{
				New[int](5, 0, 10, 0),
				New[int](5, 5, 5, -5),
			},
		},
		"t-intersection, rotated 180 deg": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0),
				New[int](5, 0, 5, 5),
			},
		},
		"t-intersection, rotated 270 deg": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 0),
				New[int](5, 5, 5, -5),
			},
		},
		"x-shape with overlap": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 10), // Diagonal segment 1
				New[int](0, 10, 10, 0), // Diagonal segment 2 (intersects at (5, 5))
				New[int](3, 3, 7, 7),   // Overlaps diagonal segment 1
			},
		},
		"multiple overlapping segments": {
			segments: []LineSegment[int]{
				New[int](1, 1, 6, 6), // Segment 1
				New[int](2, 2, 7, 7), // Segment 2
				New[int](3, 3, 5, 5), // Segment 3 (completely inside)
			},
		},
		"vertical and horizontal overlap": {
			segments: []LineSegment[int]{
				New[int](0, 0, 0, 5), // Vertical segment
				New[int](0, 0, 5, 0), // Horizontal segment
			},
		},
		"three-way intersection": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 5),  // Segment 1
				New[int](10, 0, 5, 5), // Segment 2
				New[int](5, 5, 5, 10), // Vertical Segment 3 (crosses both at (5,5))
			},
		},
		"single-point (degenerate) overlap": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 5),
				New[int](5, 5, 10, 10),
				New[int](5, 5, 5, 5), // A degenerate single-point segment
			},
		},
		"crisscrossing W shape": {
			segments: []LineSegment[int]{
				New[int](0, 0, 5, 10),
				New[int](5, 10, 10, 0),
				New[int](0, 10, 5, 0),
				New[int](5, 0, 10, 10),
			},
		},
		"zigzag": {
			segments: []LineSegment[int]{
				New[int](0, 0, 2, 2),
				New[int](2, 2, 4, 0),
				New[int](4, 0, 6, 2),
				New[int](6, 2, 8, 0),
				New[int](1, 1, 7, 1), // Horizontal line intersecting all segments
			},
		},
		"multiple_collinear_overlaps": {
			segments: []LineSegment[int]{
				New[int](0, 0, 10, 0), // Full segment
				New[int](2, 0, 8, 0),  // Inside segment
				New[int](4, 0, 6, 0),  // Inside segment
			},
		},
		"octothorpe": {
			segments: []LineSegment[int]{
				// Horizontal lines
				New[int](0, 7, 10, 7),
				New[int](0, 3, 10, 3),
				// Vertical lines
				New[int](3, 10, 3, 0),
				New[int](7, 10, 7, 0),
			},
		},
		"steep vertical slopes": {
			segments: []LineSegment[int]{
				New[int](4, 0, 5, 10),
				New[int](4, 7, 5, 5),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			for i := 0; i <= 1; i++ {
				subName := "normal"
				if i == 1 {
					subName = "input segments flipped"
					for j := range tc.segments {
						tc.segments[j] = tc.segments[j].Flip()
					}
				}

				t.Run(subName, func(t *testing.T) {
					epsilon := 1e-8
					actualIntersections := FindIntersectionsFast(tc.segments, options.WithEpsilon(epsilon))
					actualIntersectionsFromSlow := FindIntersectionsSlow(tc.segments, options.WithEpsilon(epsilon))

					t.Log("From sweep line:", actualIntersections)
					t.Log("From naive algo:", actualIntersectionsFromSlow)

					require.True(t, InterSectionResultsEq(actualIntersections, actualIntersectionsFromSlow, options.WithEpsilon(epsilon)))
				})
			}
		})
	}
}
