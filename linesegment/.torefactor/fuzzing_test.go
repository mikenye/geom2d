package _torefactor

import (
	"github.com/mikenye/geom2d/options"
	"testing"
)

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
