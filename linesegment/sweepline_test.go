package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

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

func TestFindIntersections(t *testing.T) {
	tests := map[string]struct {
		segments              []LineSegment[float64]
		expectedIntersections []IntersectionResult[float64]
	}{
		"parallel non-intersecting segments": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 5, 5),
				New[float64](0, 1, 5, 6),
			},
			expectedIntersections: []IntersectionResult[float64]{}, // No intersections
		},
		"simple X-shape intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 5, 0),
				New[float64](0, 0, 5, 5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](2.5, 2.5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 5, 5, 0),
						New[float64](5, 5, 0, 0),
					},
				},
			},
		},
		"horizontal and vertical lines": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // Horizontal segment
				New[float64](5, -5, 5, 5), // Vertical segment
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0), // Horizontal segment
						New[float64](5, 5, 5, -5), // Vertical segment
					},
				},
			},
		},
		"diagonal and horizontal lines": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 4, 4), // Diagonal line
				New[float64](2, 4, 6, 4), // Horizontal line
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](4, 4), // Intersection point
					InputLineSegments: []LineSegment[float64]{
						New[float64](2, 4, 6, 4), // Horizontal line
						New[float64](4, 4, 0, 0), // Diagonal line
					},
				},
			},
		},
		"overlapping diagonal segments": {
			segments: []LineSegment[float64]{
				New[float64](1, 1, 5, 5), // Segment 1
				New[float64](3, 3, 7, 7), // Segment 2
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](5, 5, 3, 3),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 5, 1, 1), // Segment 1
						New[float64](7, 7, 3, 3), // Segment 2,
					},
				},
			},
		},
		"overlapping horizontal segments": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // Segment 1
				New[float64](2, 0, 8, 0),  // Segment 2
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](2, 0, 8, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0), // Segment 1
						New[float64](2, 0, 8, 0),  // Segment 2
					},
				},
			},
		},
		"overlapping vertical segments": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 0, 10), // Segment 1
				New[float64](0, 2, 0, 8),  // Segment 2
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](0, 8, 0, 2),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 8, 0, 2),  // Segment 2
						New[float64](0, 10, 0, 0), // Segment 1
					},
				},
			},
		},
		"coincident segments": { // full overlap
			segments: []LineSegment[float64]{
				New[float64](1, 1, 5, 5),
				New[float64](1, 1, 5, 5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](5, 5, 1, 1),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 5, 1, 1),
						New[float64](5, 5, 1, 1),
					},
				},
			},
		},
		"shared endpoint": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 5, 5),
				New[float64](5, 5, 10, 0),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 5, 0, 0),
						New[float64](5, 5, 10, 0),
					},
				},
			},
		},
		"diamond shape": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 5, 10),
				New[float64](5, 10, 10, 5),
				New[float64](10, 5, 5, 0),
				New[float64](5, 0, 0, 5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 10),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 10, 0, 5),
						New[float64](5, 10, 10, 5),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](10, 5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](10, 5, 5, 0),
						New[float64](5, 10, 10, 5),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 5, 5, 0),
						New[float64](10, 5, 5, 0),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](0, 5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 5, 5, 0),
						New[float64](5, 10, 0, 5),
					},
				},
			},
		},
		"t-intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // Horizontal segment
				New[float64](5, -5, 5, 0), // Vertical segment terminating at (5, 0)
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0),
						New[float64](5, 0, 5, -5), // Normalized
					},
				},
			},
		},
		"t-intersection, rotated 90 deg": {
			segments: []LineSegment[float64]{
				New[float64](5, 0, 10, 0),
				New[float64](5, 5, 5, -5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 0, 10, 0),
						New[float64](5, 5, 5, -5), // Normalized
					},
				},
			},
		},
		"t-intersection, rotated 180 deg": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 10, 0),
				New[float64](5, 0, 5, 5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 10, 0),
						New[float64](5, 5, 5, 0),
					},
				},
			},
		},
		"t-intersection, rotated 270 deg": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 5, 0),
				New[float64](5, 5, 5, -5),
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 0, 5, 0),
						New[float64](5, 5, 5, -5),
					},
				},
			},
		},
		"x-shape with overlap": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 10, 10), // Diagonal segment 1
				New[float64](0, 10, 10, 0), // Diagonal segment 2 (intersects at (5, 5))
				New[float64](3, 3, 7, 7),   // Overlaps diagonal segment 1
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](7, 7, 3, 3),
					InputLineSegments: []LineSegment[float64]{
						New[float64](7, 7, 3, 3),
						New[float64](10, 10, 0, 0),
					},
				},
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](5, 5),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 10, 10, 0),             // Diagonal segment 2 (intersects at (5, 5))
						New[float64](0, 0, 10, 10).normalize(), // Diagonal segment 1
						New[float64](3, 3, 7, 7).normalize(),
					},
				},
			},
		},
		"multiple overlapping segments": {
			segments: []LineSegment[float64]{
				New[float64](1, 1, 6, 6), // Segment 1
				New[float64](2, 2, 7, 7), // Segment 2
				New[float64](3, 3, 5, 5), // Segment 3 (completely inside)
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](6, 6, 2, 2),
					InputLineSegments: []LineSegment[float64]{
						New[float64](7, 7, 2, 2),
						New[float64](6, 6, 1, 1),
					},
				},
				{
					IntersectionType:   IntersectionOverlappingSegment,
					OverlappingSegment: New[float64](5, 5, 3, 3),
					InputLineSegments: []LineSegment[float64]{
						New[float64](5, 5, 3, 3),
						New[float64](6, 6, 1, 1),
						New[float64](7, 7, 2, 2),
					},
				},
			},
		},
		"vertical and horizontal overlap": {
			segments: []LineSegment[float64]{
				New[float64](0, 0, 0, 5), // Vertical segment
				New[float64](0, 0, 5, 0), // Horizontal segment
			},
			expectedIntersections: []IntersectionResult[float64]{
				{
					IntersectionType:  IntersectionPoint,
					IntersectionPoint: point.New[float64](0, 0),
					InputLineSegments: []LineSegment[float64]{
						New[float64](0, 5, 0, 0), // Vertical segment
						New[float64](0, 0, 5, 0), // Horizontal segment
					},
				},
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

					actualIntersections := FindIntersections(tc.segments, options.WithEpsilon(epsilon))
					actualIntersectionsFromSlow := FindIntersectionsSlow(tc.segments, options.WithEpsilon(epsilon))

					fmt.Println("expected:       ", tc.expectedIntersections)
					fmt.Println("From sweep line:", actualIntersections)
					fmt.Println("From naive algo:", actualIntersectionsFromSlow)

					assert.ElementsMatch(t, actualIntersections, actualIntersectionsFromSlow, "mismatch between FindIntersections and FindIntersectionsSlow")

					// Assert number of intersections
					assert.Equal(t, len(tc.expectedIntersections), len(actualIntersections), "Number of intersections mismatch")

					// Assert each intersection point and its associated segments
					for _, expected := range tc.expectedIntersections {
						found := false
						switch expected.IntersectionType {
						case IntersectionNone:
							t.Fatal("unexpected IntersectionNone")
						case IntersectionPoint:
							for _, actual := range actualIntersections {
								if expected.IntersectionPoint.Eq(actual.IntersectionPoint) {
									assert.ElementsMatch(t, expected.InputLineSegments, actual.InputLineSegments, "Segments mismatch at intersection point")
									found = true
									break
								}
							}
						case IntersectionOverlappingSegment:
							for _, actual := range actualIntersections {
								if expected.OverlappingSegment.Eq(actual.OverlappingSegment) {
									assert.ElementsMatch(t, expected.InputLineSegments, actual.InputLineSegments, "Segments mismatch at intersection point")
									found = true
									break
								}
							}
						}
						assert.True(t, found, "Expected intersection not found")
					}
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
		"leftmost and rightmost with both neighbors": {
			point: point.New[float64](3, 3),
			UofP: []LineSegment[float64]{
				New[float64](2, 2, 3, 3), // Part of U(p)
			},
			CofP: []LineSegment[float64]{
				New[float64](4, 4, 3, 3), // Part of C(p)
			},
			statusItems: []statusItem{
				{segment: New[float64](1, 1, 3, 3), sweepY: 3}, // Left neighbor
				{segment: New[float64](2, 2, 3, 3), sweepY: 3}, // U(p)
				{segment: New[float64](4, 4, 3, 3), sweepY: 3}, // C(p)
				{segment: New[float64](5, 5, 3, 3), sweepY: 3}, // Right neighbor
			},
			expectedSL:           &statusItem{segment: New[float64](1, 1, 3, 3), sweepY: 3},
			expectedSR:           &statusItem{segment: New[float64](5, 5, 3, 3), sweepY: 3},
			expectedSPrime:       linSegPtr(2, 2, 3, 3),
			expectedSDoublePrime: linSegPtr(4, 4, 3, 3),
		},
		"no neighbors with single segment": {
			point: point.New[float64](3, 3),
			UofP: []LineSegment[float64]{
				New[float64](2, 2, 3, 3),
			},
			CofP:                 nil,
			statusItems:          []statusItem{{segment: New[float64](2, 2, 3, 3), sweepY: 3}},
			expectedSL:           nil,
			expectedSR:           nil,
			expectedSPrime:       linSegPtr(2, 2, 3, 3),
			expectedSDoublePrime: linSegPtr(2, 2, 3, 3),
		},
		"left neighbor only": {
			point: point.New[float64](3, 3),
			UofP: []LineSegment[float64]{
				New[float64](2, 2, 3, 3),
			},
			CofP: []LineSegment[float64]{
				New[float64](4, 4, 3, 3),
			},
			statusItems: []statusItem{
				{segment: New[float64](1, 1, 3, 3), sweepY: 3}, // Left neighbor
				{segment: New[float64](2, 2, 3, 3), sweepY: 3}, // U(p)
				{segment: New[float64](4, 4, 3, 3), sweepY: 3}, // C(p)
			},
			expectedSL:           &statusItem{segment: New[float64](1, 1, 3, 3), sweepY: 3},
			expectedSR:           nil,
			expectedSPrime:       linSegPtr(2, 2, 3, 3),
			expectedSDoublePrime: linSegPtr(4, 4, 3, 3),
		},
		"right neighbor only": {
			point: point.New[float64](3, 3),
			UofP:  nil,
			CofP: []LineSegment[float64]{
				New[float64](4, 4, 3, 3),
			},
			statusItems: []statusItem{
				{segment: New[float64](4, 4, 3, 3), sweepY: 3}, // C(p)
				{segment: New[float64](5, 5, 3, 3), sweepY: 3}, // Right neighbor
			},
			expectedSL:           nil,
			expectedSR:           &statusItem{segment: New[float64](5, 5, 3, 3), sweepY: 3},
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
				{segment: New[float64](2, 4, 1, 1), sweepY: 3}, // Left neighbor
				{segment: New[float64](4, 6, 3, 3), sweepY: 3}, // Matching segment
				{segment: New[float64](4, 4, 5, 2), sweepY: 3}, // Right neighbor
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  &statusItem{segment: New[float64](2, 4, 1, 1), sweepY: 3},
			expectedRight: &statusItem{segment: New[float64](4, 4, 5, 2), sweepY: 3},
		},
		"point with no left neighbor": {
			statusItems: []statusItem{
				{segment: New[float64](4, 6, 3, 3), sweepY: 3}, // Matching segment
				{segment: New[float64](4, 4, 5, 2), sweepY: 3}, // Right neighbor
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  nil,
			expectedRight: &statusItem{segment: New[float64](4, 4, 5, 2), sweepY: 3},
		},
		"point with no right neighbor": {
			statusItems: []statusItem{
				{segment: New[float64](2, 4, 1, 1), sweepY: 3}, // Left neighbor
				{segment: New[float64](4, 6, 3, 3), sweepY: 3}, // Matching segment
			},
			point:         qItem{point: point.New[float64](3, 3)},
			expectedLeft:  &statusItem{segment: New[float64](2, 4, 1, 1), sweepY: 3},
			expectedRight: nil,
		},
		"point with no neighbors (single item)": {
			statusItems: []statusItem{
				{segment: New[float64](1, 1, 3, 3), sweepY: 3},
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
				{segment: New[float64](3, 1, 5, 5), sweepY: 0}, // Will have XAtY = 4 at sweepY = 3
				{segment: New[float64](1, 1, 3, 3), sweepY: 0}, // Will have XAtY = 3 at sweepY = 3
				{segment: New[float64](2, 4, 1, 1), sweepY: 0}, // Will have XAtY = 1.5 at sweepY = 3
			},
			sweepY: 3,
			expected: []statusItem{
				{segment: New[float64](2, 4, 1, 1), sweepY: 3}, // XAtY = 1.5
				{segment: New[float64](1, 1, 3, 3), sweepY: 3}, // XAtY = 3
				{segment: New[float64](3, 1, 5, 5), sweepY: 3}, // XAtY = 4
			},
		},
		"horizontal segments": {
			input: []statusItem{
				{segment: New[float64](1, 3, 5, 3), sweepY: 0}, // Horizontal at y = 3
				{segment: New[float64](3, 1, 3, 5), sweepY: 0}, // Vertical at x = 3
			},
			sweepY: 3,
			expected: []statusItem{
				{segment: New[float64](3, 1, 3, 5), sweepY: 3}, // Vertical at x = 3
				{segment: New[float64](1, 3, 5, 3), sweepY: 3}, // Horizontal at y = 3
			},
		},
		"already sorted": {
			input: []statusItem{
				{segment: New[float64](2, 4, 1, 1), sweepY: 0},
				{segment: New[float64](1, 1, 3, 3), sweepY: 0},
				{segment: New[float64](3, 1, 5, 5), sweepY: 0},
			},
			sweepY: 3,
			expected: []statusItem{
				{segment: New[float64](2, 4, 1, 1), sweepY: 3},
				{segment: New[float64](1, 1, 3, 3), sweepY: 3},
				{segment: New[float64](3, 1, 5, 5), sweepY: 3},
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
			sortStatusBySweepLine(tc.input, tc.sweepY)

			for i := range tc.input {
				assert.Equal(t, tc.expected[i].segment, tc.input[i].segment, "Segment mismatch at index %d", i)
				assert.Equal(t, tc.sweepY, tc.input[i].sweepY, "SweepY mismatch at index %d", i)
			}
		})
	}
}
