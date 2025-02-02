package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/assert"
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
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive: %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
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
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive: %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
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
		if x7 == x7 && y8 == y8 {
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
			t.Fatalf("Mismatch found!\nSegments: %v\nNaive: %v\nSweep Line: %v\nReason: %v\n", segments, naiveResults, sweepResults, reason)
		}
	})
}

func TestDeleteSegmentsFromStatus(t *testing.T) {
	// Define the status structure as a sorted slice
	status := []statusItem{
		{segment: New[float64](1, 5, 3, 1)},
		{segment: New[float64](3, 1, 5, 5)},
		{segment: New[float64](5, 5, 1, 1)},
	}

	// Segments to delete
	toDelete := []LineSegment[float64]{
		New[float64](3, 1, 5, 5),
		New[float64](5, 5, 1, 1),
	}

	// Call deleteSegmentsFromStatus
	newStatus := deleteSegmentsFromStatus(status, toDelete, options.WithEpsilon(1e-6))

	// Expected status structure
	expected := []statusItem{
		{segment: New[float64](1, 5, 3, 1)},
	}

	// Assert the result
	assert.Equal(t, expected, newStatus, "Remaining status structure should match expected")
}

// TestFindIntersections ensures that the output of the sweep line algorithm (FindIntersectionsFast) matches
// the output of the naïve algorithm (FindIntersectionsSlow).
func TestFindIntersections(t *testing.T) {
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
				New[int](10, 0, 0, 0),
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

					fmt.Printf("From sweep line: %#v\n", actualIntersections)
					fmt.Printf("From naive algo: %#v\n", actualIntersectionsFromSlow)

					ok, reason := compareIntersectionResults(actualIntersections, actualIntersectionsFromSlow, epsilon)
					require.True(t, ok, reason)
				})
			}
		})
	}
}

func TestFindLeftmostAndRightmostSegmentAndNeighbors(t *testing.T) {

	// Helper function to create a pointer to a point
	linSegPtr := func(x1, y1, x2, y2 float64) *LineSegment[float64] {
		l := New[float64](x1, y1, x2, y2)
		return &l
	}

	tests := map[string]struct {
		point                point.Point[float64]
		UofP                 []LineSegment[float64]
		CofP                 []LineSegment[float64]
		statusItems          []statusItem
		expectedSL           *statusItem
		expectedSR           *statusItem
		expectedSPrime       *LineSegment[float64]
		expectedSDoublePrime *LineSegment[float64]
	}{
		"no neighbors with single segment": {
			point: point.New[float64](3, 3),
			UofP: []LineSegment[float64]{
				New[float64](2, 2, 3, 3),
			},
			CofP:                 nil,
			statusItems:          []statusItem{{segment: New[float64](2, 2, 3, 3)}},
			expectedSL:           nil,
			expectedSR:           nil,
			expectedSPrime:       linSegPtr(2, 2, 3, 3),
			expectedSDoublePrime: linSegPtr(2, 2, 3, 3),
		},
		"right neighbor only": {
			point: point.New[float64](3, 3),
			UofP:  nil,
			CofP: []LineSegment[float64]{
				New[float64](4, 4, 3, 3),
			},
			statusItems: []statusItem{
				{segment: New[float64](4, 4, 3, 3)}, // C(p)
				{segment: New[float64](5, 5, 3, 3)}, // Right neighbor
			},
			expectedSL:           nil,
			expectedSR:           &statusItem{segment: New[float64](5, 5, 3, 3)},
			expectedSPrime:       linSegPtr(4, 4, 3, 3),
			expectedSDoublePrime: linSegPtr(4, 4, 3, 3),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sPrime, sDoublePrime, sL, sR := findLeftmostAndRightmostSegmentAndNeighbors(tc.point, tc.UofP, tc.CofP, tc.statusItems)

			if tc.expectedSPrime == nil {
				assert.Nil(t, sPrime, "Expected sPrime to be nil")
			} else {
				assert.NotNil(t, sPrime, "Expected sPrime to be non-nil")
				assert.Equal(t, *tc.expectedSPrime, *sPrime, "sPrime neighbor mismatch")
			}

			if tc.expectedSDoublePrime == nil {
				assert.Nil(t, sDoublePrime, "Expected sPrime to be nil")
			} else {
				assert.NotNil(t, sDoublePrime, "Expected sPrime to be non-nil")
				assert.Equal(t, *tc.expectedSDoublePrime, *sDoublePrime, "sPrime neighbor mismatch")
			}

			// Check sL neighbor
			if tc.expectedSL == nil {
				assert.Nil(t, sL, "Expected sL neighbor to be nil")
			} else {
				assert.NotNil(t, sL, "Expected sL neighbor to be non-nil")
				assert.Equal(t, *tc.expectedSL, *sL, "Left neighbor mismatch")
			}

			// Check sR neighbor
			if tc.expectedSR == nil {
				assert.Nil(t, sR, "Expected sR neighbor to be nil")
			} else {
				assert.NotNil(t, sR, "Expected sR neighbor to be non-nil")
				assert.Equal(t, *tc.expectedSR, *sR, "Right neighbor mismatch")
			}
		})
	}
}

func TestFindNeighbors(t *testing.T) {
	tests := map[string]struct {
		statusItems   []statusItem // The sorted status structure (slice of statusItem)
		point         qItem        // The event point to find neighbors for
		expectedLeft  *statusItem  // Expected left neighbor (nil if none exists)
		expectedRight *statusItem  // Expected right neighbor (nil if none exists)
	}{
		"point with both neighbors": {
			statusItems: []statusItem{
				{segment: New[float64](2, 4, 1, 1)}, // Left neighbor
				{segment: New[float64](4, 6, 3, 3)}, // Matching segment
				{segment: New[float64](4, 4, 5, 2)}, // Right neighbor
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  &statusItem{segment: New[float64](2, 4, 1, 1)},
			expectedRight: &statusItem{segment: New[float64](4, 4, 5, 2)},
		},
		"point with no left neighbor": {
			statusItems: []statusItem{
				{segment: New[float64](4, 6, 3, 3)}, // Matching segment
				{segment: New[float64](4, 4, 5, 2)}, // Right neighbor
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  nil,
			expectedRight: &statusItem{segment: New[float64](4, 4, 5, 2)},
		},
		"point with no right neighbor": {
			statusItems: []statusItem{
				{segment: New[float64](2, 4, 1, 1)}, // Left neighbor
				{segment: New[float64](4, 6, 3, 3)}, // Matching segment
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  &statusItem{segment: New[float64](2, 4, 1, 1)},
			expectedRight: nil,
		},
		"point with no neighbors (single item)": {
			statusItems: []statusItem{
				{segment: New[float64](1, 1, 3, 3)},
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  nil,
			expectedRight: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			t.Logf("Status Items: %v", tc.statusItems)
			left, right := findNeighbors(tc.statusItems, tc.point)

			// Check left neighbor
			if tc.expectedLeft == nil {
				assert.Nil(t, left, "Expected left neighbor to be nil")
			} else {
				assert.NotNil(t, left, "Expected left neighbor to be non-nil")
				assert.Equal(t, *tc.expectedLeft, *left, "Left neighbor mismatch")
			}

			// Check right neighbor
			if tc.expectedRight == nil {
				assert.Nil(t, right, "Expected right neighbor to be nil")
			} else {
				assert.NotNil(t, right, "Expected right neighbor to be non-nil")
				assert.Equal(t, *tc.expectedRight, *right, "Right neighbor mismatch")
			}
		})
	}
}

func TestFindNewEvent(t *testing.T) {
	// Define a small epsilon for geometric calculations
	epsilon := 1e-9
	opts := []options.GeometryOptionsFunc{options.WithEpsilon(epsilon)}

	// Helper function to create a pointer to a point
	pointPtr := func(x, y float64) *point.Point[float64] {
		p := point.New[float64](x, y)
		return &p
	}

	tests := map[string]struct {
		segment1         LineSegment[float64]
		segment2         LineSegment[float64]
		currentPoint     point.Point[float64]
		expectedPoint    *point.Point[float64] // Nil if no intersection is expected
		expectedInQueue  bool                  // Whether the point should be in the queue
		expectedCountInQ int                   // The number of matching points in the queue
	}{
		"valid intersection below current point": {
			segment1:         New[float64](0, 5, 5, 0),
			segment2:         New[float64](0, 0, 5, 5),
			currentPoint:     point.New[float64](3, 4),
			expectedPoint:    pointPtr(2.5, 2.5),
			expectedInQueue:  true,
			expectedCountInQ: 1,
		},
		"intersection above current point": {
			segment1:         New[float64](0, 5, 5, 0),
			segment2:         New[float64](0, 0, 5, 5),
			currentPoint:     point.New[float64](2, 2),
			expectedPoint:    nil,
			expectedInQueue:  false,
			expectedCountInQ: 0,
		},
		"no intersection": {
			segment1:         New[float64](0, 5, 5, 5),
			segment2:         New[float64](0, 0, 5, 0),
			currentPoint:     point.New[float64](3, 4),
			expectedPoint:    nil,
			expectedInQueue:  false,
			expectedCountInQ: 0,
		},
		"intersection already in queue": {
			segment1:         New[float64](0, 5, 5, 0),
			segment2:         New[float64](0, 0, 5, 5),
			currentPoint:     point.New[float64](3, 4),
			expectedPoint:    pointPtr(2.5, 2.5),
			expectedInQueue:  true,
			expectedCountInQ: 1, // Ensure no duplicates are added
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize the event queue
			Q := btree.NewG[qItem](2, qItemLess)

			// Initialise results
			R := newIntersectionResults[float64]()

			// Prepopulate the queue for cases where intersection is already in queue
			if tc.expectedPoint != nil && tc.expectedCountInQ > 0 {
				Q.ReplaceOrInsert(qItem{
					point:    *tc.expectedPoint,
					segments: []LineSegment[float64]{tc.segment1, tc.segment2},
				})
			}

			// Call findNewEvent
			findNewEvent(tc.segment1, tc.segment2, tc.currentPoint, Q, R, opts...)

			// Verify the count of matching points in the queue
			matchingCount := 0
			Q.Ascend(func(item qItem) bool {
				if tc.expectedPoint != nil && item.point.Eq(*tc.expectedPoint, opts...) {
					matchingCount++
				}
				return true
			})
			assert.Equal(t, tc.expectedCountInQ, matchingCount,
				"Mismatch in the number of matching intersection points in the queue")

			// Check if the expected point is in the queue
			if tc.expectedPoint != nil {
				found := Q.Has(qItem{point: *tc.expectedPoint})
				assert.Equal(t, tc.expectedInQueue, found,
					"Intersection point presence in queue mismatch")
			} else {
				// If no expected point, ensure no new event was added
				assert.Equal(t, 0, Q.Len(), "No new intersection point should be added")
			}
		})
	}
}

func TestInsertSegmentIntoQueue(t *testing.T) {
	// Create a new btree for the event queue
	Q := btree.NewG[qItem](2, qItemLess)

	// Define test segments
	seg1 := New[float64](1, 5, 3, 1) // Upper endpoint (1,5), lower endpoint (3,1)
	seg2 := New[float64](3, 1, 5, 5) // Upper endpoint (5,5), lower endpoint (3,1)
	seg3 := New[float64](1, 5, 2, 2) // Upper endpoint (1,5), lower endpoint (2,2)

	// Insert the first segment
	insertSegmentIntoQueue(seg1, Q)

	// Verify upper and lower points for seg1
	upper1 := qItem{point: point.New[float64](1, 5), segments: []LineSegment[float64]{seg1.normalize()}}
	lower1 := qItem{point: point.New[float64](3, 1), segments: nil}
	upper1FromQ, ok := Q.Get(upper1)
	require.True(t, ok, "Upper endpoint should exist in Q")
	assert.Equal(t, upper1, upper1FromQ, "Upper endpoint should match expected")
	assert.True(t, Q.Has(lower1), "Lower endpoint should exist in Q without a segment")

	// Insert the second segment
	insertSegmentIntoQueue(seg2, Q)

	// Verify upper and lower points for seg2
	upper2 := qItem{point: point.New[float64](5, 5), segments: []LineSegment[float64]{seg2.normalize()}}
	lower2 := qItem{point: point.New[float64](3, 1), segments: nil}
	upper2FromQ, ok := Q.Get(upper2)
	require.True(t, ok, "Upper endpoint should exist in Q")
	assert.Equal(t, upper2, upper2FromQ, "Upper endpoint should match expected")
	assert.True(t, Q.Has(lower2), "Lower endpoint should exist in Q without a segment")

	// Insert the third segment
	insertSegmentIntoQueue(seg3, Q)

	// Verify upper and lower points for seg3
	// Note: upper1's segment list should now include seg3
	upper3 := qItem{point: point.New[float64](1, 5), segments: []LineSegment[float64]{seg1.normalize(), seg3.normalize()}}
	lower3 := qItem{point: point.New[float64](2, 2), segments: nil}
	upper3FromQ, ok := Q.Get(upper3)
	require.True(t, ok, "Upper endpoint should exist in Q")
	assert.Equal(t, upper3, upper3FromQ, "Upper endpoint should match expected")
	assert.True(t, Q.Has(lower3), "Lower endpoint should exist in Q without a segment")
}

func TestSortStatusBySweepLine(t *testing.T) {
	tests := map[string]struct {
		input    []statusItem
		sweepY   float64
		expected []statusItem
	}{
		"basic case": {
			input: []statusItem{
				{segment: New[float64](3, 1, 5, 5)}, // Will have XAtY = 4 at sweepY = 3
				{segment: New[float64](1, 1, 3, 3)}, // Will have XAtY = 3 at sweepY = 3
				{segment: New[float64](2, 4, 1, 1)}, // Will have XAtY = 1.5 at sweepY = 3
			},
			sweepY: 3,
			expected: []statusItem{
				{segment: New[float64](2, 4, 1, 1)}, // XAtY = 1.5
				{segment: New[float64](1, 1, 3, 3)}, // XAtY = 3
				{segment: New[float64](3, 1, 5, 5)}, // XAtY = 4
			},
		},
		"already sorted": {
			input: []statusItem{
				{segment: New[float64](2, 4, 1, 1)},
				{segment: New[float64](1, 1, 3, 3)},
				{segment: New[float64](3, 1, 5, 5)},
			},
			sweepY: 3,
			expected: []statusItem{
				{segment: New[float64](2, 4, 1, 1)},
				{segment: New[float64](1, 1, 3, 3)},
				{segment: New[float64](3, 1, 5, 5)},
			},
		},
		"empty input": {
			input:    []statusItem{},
			sweepY:   3,
			expected: []statusItem{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			sortStatusBySweepLine(tc.input, qItem{point: point.New(0, tc.sweepY)})

			debugPrintStatus(tc.input, tc.sweepY)

			fmt.Println("input: ", tc.input)
			fmt.Println("expect:", tc.expected)

			for i := range tc.input {
				assert.Equal(t, tc.expected[i].segment, tc.input[i].segment, "Segment mismatch at index %d", i)
			}
		})
	}
}
