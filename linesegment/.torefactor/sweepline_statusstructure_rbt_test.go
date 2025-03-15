package _torefactor

import (
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"slices"
	"testing"
)

func TestStatusStructureRBT(t *testing.T) {
	tests := map[string]struct {
		// input line segments.
		// must have upper p as start
		segments []LineSegment[float64]

		// sweep line p (should be start/end p of segment, or intersection).
		// should be a valid sweep line p given the segments in statys structure.
		// all segments above must be on or to the right or below of this p.
		p point.Point[float64]

		// line segments as-per their order in the status structure at p
		expected []LineSegment[float64]
	}{
		"cross, sweep line above intersection": {
			segments: []LineSegment[float64]{
				New[float64](5, 5, 0, 0),
				New[float64](0, 5, 5, 0),
			},
			p: point.New[float64](5, 5), // upper p of 5,5,0,0 segment
			expected: []LineSegment[float64]{
				New[float64](0, 5, 5, 0),
				New[float64](5, 5, 0, 0),
			},
		},
		"cross, sweep line on intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 5, 0),
				New[float64](5, 5, 0, 0),
			},
			p: point.New[float64](2.5, 2.5), // intersection
			expected: []LineSegment[float64]{
				New[float64](5, 5, 0, 0),
				New[float64](0, 5, 5, 0),
			},
		},
		"cross, sweep line below intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 5, 0),
				New[float64](5, 5, 0, 0),
			},
			p: point.New[float64](0, 0), // end p of leftmost line
			expected: []LineSegment[float64]{
				New[float64](5, 5, 0, 0),
				New[float64](0, 5, 5, 0),
			},
		},
		"vertical & negative slope diagonal, sweep line above intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 20, 20, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](10, 20),
			expected: []LineSegment[float64]{
				New[float64](0, 20, 20, 0), // diagonal first
				New[float64](10, 20, 10, 0),
			},
		},
		"vertical & negative slope diagonal, sweep line at intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 20, 20, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](10, 10),
			expected: []LineSegment[float64]{
				New[float64](10, 20, 10, 0), // vertical first
				New[float64](0, 20, 20, 0),
			},
		},
		"vertical & negative slope diagonal, sweep line below intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 20, 20, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](10, 0),
			expected: []LineSegment[float64]{
				New[float64](10, 20, 10, 0), // vertical first
				New[float64](0, 20, 20, 0),
			},
		},
		"vertical & positive slope diagonal, sweep line above intersection": {
			segments: []LineSegment[float64]{
				New[float64](20, 20, 0, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](10, 20),
			expected: []LineSegment[float64]{
				New[float64](10, 20, 10, 0), // vertical first
				New[float64](20, 20, 0, 0),
			},
		},
		"vertical & positive slope diagonal, sweep line at intersection": {
			segments: []LineSegment[float64]{
				New[float64](20, 20, 0, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](10, 10),
			expected: []LineSegment[float64]{
				New[float64](20, 20, 0, 0), // diagonal first
				New[float64](10, 20, 10, 0),
			},
		},
		"vertical & positive slope diagonal, sweep line below intersection": {
			segments: []LineSegment[float64]{
				New[float64](20, 20, 0, 0),
				New[float64](10, 20, 10, 0),
			},
			p: point.New[float64](0, 0),
			expected: []LineSegment[float64]{
				New[float64](20, 20, 0, 0), // diagonal first
				New[float64](10, 20, 10, 0),
			},
		},
		"horizontal & negative slope diagonal, sweep line above (to the left of) intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](0, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // horizontal first
				New[float64](0, 5, 10, -5),
			},
		},
		"horizontal & negative slope diagonal, sweep line at intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](5, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // diagonal first
				New[float64](0, 0, 10, 0),
			},
		},
		"horizontal & negative slope diagonal, sweep line below (to the right of) intersection": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](10, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // diagonal first
				New[float64](0, 0, 10, 0),
			},
		},
		"horizontal & positive slope diagonal, sweep line above (to the left of) intersection": {
			segments: []LineSegment[float64]{
				New[float64](10, 5, 0, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](0, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // horizontal first
				New[float64](10, 5, 0, -5),
			},
		},
		"horizontal & positive slope diagonal, sweep line at intersection": {
			segments: []LineSegment[float64]{
				New[float64](10, 5, 0, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](5, 0),
			expected: []LineSegment[float64]{
				New[float64](10, 5, 0, -5),
				New[float64](0, 0, 10, 0),
			},
		},
		"horizontal & positive slope diagonal, sweep line below (to the right of) intersection": {
			segments: []LineSegment[float64]{
				New[float64](10, 5, 0, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](10, 0),
			expected: []LineSegment[float64]{
				New[float64](10, 5, 0, -5),
				New[float64](0, 0, 10, 0),
			},
		},
		"horizontal & vertical, sweep line above (to the left of) intersection": {
			segments: []LineSegment[float64]{
				New[float64](5, 5, 5, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](0, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 0, 10, 0), // horizontal first
				New[float64](5, 5, 5, -5),
			},
		},
		"horizontal & vertical, sweep line at intersection": {
			segments: []LineSegment[float64]{
				New[float64](5, 5, 5, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](5, 0),
			expected: []LineSegment[float64]{
				New[float64](5, 5, 5, -5),
				New[float64](0, 0, 10, 0), // horizontal last
			},
		},
		"horizontal & vertical, sweep line after intersection": {
			segments: []LineSegment[float64]{
				New[float64](5, 5, 5, -5),
				New[float64](0, 0, 10, 0),
			},
			p: point.New[float64](5, 0),
			expected: []LineSegment[float64]{
				New[float64](5, 5, 5, -5),
				New[float64](0, 0, 10, 0), // horizontal last
			},
		},
		"asterisk, sweep line before intersection, above horizontal segment": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
				//New[float64](0, 0, 10, 0), // Horiz.
			},
			p: point.New[float64](10, 5),
			expected: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
				//New[float64](0, 0, 10, 0), // Horiz.
			},
		},
		"asterisk, sweep line before intersection, start of horizontal segment": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
				New[float64](0, 0, 10, 0),  // Horiz.
			},
			p: point.New[float64](0, 0),
			expected: []LineSegment[float64]{
				New[float64](0, 0, 10, 0),  // Horiz.
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
			},
		},
		"asterisk, sweep line after intersection, intersection p on horizontal segment": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
				New[float64](0, 0, 10, 0),  // Horiz.
			},
			p: point.New[float64](5, 0),
			expected: []LineSegment[float64]{
				New[float64](10, 5, 0, -5), // TR to BL
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](0, 0, 10, 0),  // Horiz.
			},
		},
		"asterisk, sweep line after intersection, end of horizontal segment": {
			segments: []LineSegment[float64]{
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](10, 5, 0, -5), // TR to BL
				New[float64](0, 0, 10, 0),  // Horiz.
			},
			p: point.New[float64](10, 0),
			expected: []LineSegment[float64]{
				New[float64](10, 5, 0, -5), // TR to BL
				New[float64](5, 5, 5, -5),  // Vert.
				New[float64](0, 5, 10, -5), // TL to BR
				New[float64](0, 0, 10, 0),  // Horiz.
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			mode := []string{"normal", "input segment order reversed"}

			for i, subname := range mode {

				t.Run(subname, func(t *testing.T) {

					if i == 1 {
						slices.Reverse(tc.segments)
					}

					// set up status structure
					S := newStatusStructureRBT(tc.p)

					// insert points into status structure
					for _, seg := range tc.segments {
						S.Insert(seg)
					}
					log.Printf("Status Structure:\n%s", S.String())

					// ensure lengths match
					assert.Equal(t, len(tc.expected), S.structure.Size(), "length mismatch")

					// check output matches
					iter := S.structure.Iterator()
					for n, expected := range tc.expected {
						t.Logf("checking index %d", n)
						found := iter.Next()
						require.True(t, found, "could not pop from S")
						actual := iter.Key().(statusStructureEntry).segment
						t.Log("popped:", actual)
						t.Log("expect:", expected)
						assert.Equal(t, expected, actual, "segment mismatch")
					}
				})
			}
		})
	}
}

//func TestStatusStructureRBT_FindCofPAndLofP(t *testing.T) {
//	tests := map[string]struct {
//		// input line segments.
//		// must have upper p as start
//		segments []LineSegment[float64]
//
//		// sweep line p (should be start/end p of segment, or intersection).
//		// should be a valid sweep line p given the segments in statys structure.
//		// all segments above must be on or to the right or below of this p.
//		p point.Point[float64]
//
//		expectedCofP, expectedLofP []LineSegment[float64]
//	}{
//		"simple X-shape intersection #1": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 5, 5, 0),
//				New[float64](5, 5, 0, 0),
//			},
//			p:            point.New[float64](5, 5),
//			expectedCofP: nil,
//			expectedLofP: nil,
//		},
//		"simple X-shape intersection #2": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 5, 5, 0),
//				New[float64](5, 5, 0, 0),
//			},
//			p: point.New[float64](2.5, 2.5),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](0, 5, 5, 0),
//				New[float64](5, 5, 0, 0),
//			},
//			expectedLofP: nil,
//		},
//		"simple X-shape intersection #3": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 5, 5, 0),
//				New[float64](5, 5, 0, 0),
//			},
//			p:            point.New[float64](0, 0),
//			expectedCofP: nil,
//			expectedLofP: []LineSegment[float64]{
//				New[float64](5, 5, 0, 0),
//			},
//		},
//		"diagonal and horizontal lines #1": {
//			segments: []LineSegment[float64]{
//				New[float64](4, 4, 0, 0), // Diagonal line
//				New[float64](2, 4, 6, 4), // Horizontal line
//			},
//			p: point.New[float64](4, 4),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](2, 4, 6, 4),
//			},
//			expectedLofP: nil,
//		},
//		"diagonal and horizontal lines #2": {
//			segments: []LineSegment[float64]{
//				New[float64](4, 4, 0, 0), // Diagonal line
//				New[float64](2, 4, 6, 4), // Horizontal line
//			},
//			p:            point.New[float64](6, 4),
//			expectedCofP: nil,
//			expectedLofP: []LineSegment[float64]{
//				New[float64](2, 4, 6, 4),
//			},
//		},
//		"overlapping diagonal segments #1": {
//			segments: []LineSegment[float64]{
//				New[float64](5, 5, 1, 1), // Segment 1
//				New[float64](7, 7, 3, 3), // Segment 2
//			},
//			p: point.New[float64](5, 5),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](7, 7, 3, 3),
//			},
//			expectedLofP: nil,
//		},
//		"overlapping diagonal segments #2": {
//			segments: []LineSegment[float64]{
//				New[float64](5, 5, 1, 1), // Segment 1
//				New[float64](7, 7, 3, 3), // Segment 2
//			},
//			p: point.New[float64](3, 3),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](5, 5, 1, 1),
//			},
//			expectedLofP: []LineSegment[float64]{
//				New[float64](7, 7, 3, 3),
//			},
//		},
//		"overlapping horizontal segments #1": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 0, 10, 0), // Segment 1
//				New[float64](2, 0, 8, 0),  // Segment 2
//			},
//			p: point.New[float64](2, 0),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](0, 0, 10, 0),
//			},
//			expectedLofP: nil,
//		},
//		"overlapping horizontal segments #2": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 0, 10, 0), // Segment 1
//				New[float64](2, 0, 8, 0),  // Segment 2
//			},
//			p: point.New[float64](8, 0),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](0, 0, 10, 0),
//			},
//			expectedLofP: []LineSegment[float64]{
//				New[float64](2, 0, 8, 0),
//			},
//		},
//		"overlapping vertical segments #1": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 10, 0, 0), // Segment 1
//				New[float64](0, 8, 0, 2),  // Segment 2
//			},
//			p: point.New[float64](0, 8),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](0, 10, 0, 0),
//			},
//			expectedLofP: nil,
//		},
//		"overlapping vertical segments #2": {
//			segments: []LineSegment[float64]{
//				New[float64](0, 10, 0, 0), // Segment 1
//				New[float64](0, 8, 0, 2),  // Segment 2
//			},
//			p: point.New[float64](0, 2),
//			expectedCofP: []LineSegment[float64]{
//				New[float64](0, 10, 0, 0),
//			},
//			expectedLofP: []LineSegment[float64]{
//				New[float64](0, 8, 0, 2),
//			},
//		},
//	}
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//
//			mode := []string{"normal", "input segment order reversed"}
//
//			for i, subname := range mode {
//				t.Run(subname, func(t *testing.T) {
//
//					// if input segment order reversed, then reverse the input segments :)
//					if i == 1 {
//						slices.Reverse(tc.segments)
//					}
//
//					// set up status structure
//					S := newStatusStructureRBT(tc.p)
//
//					// insert points into status structure
//					for _, seg := range tc.segments {
//						S.Insert(seg)
//					}
//					log.Printf("Status Structure:\n%s", S.String())
//
//					CofP, LofP := S.FindCofPAndLofP(tc.p)
//
//					log.Println("CofP:", CofP)
//					log.Println("LofP:", LofP)
//
//					assert.ElementsMatch(t, tc.expectedCofP, CofP, "CofP element mismatch")
//					assert.ElementsMatch(t, tc.expectedLofP, LofP, "LofP element mismatch")
//
//				})
//			}
//		})
//	}
//}

func TestStatusStructureRBT_FindNeighborsOfPoint(t *testing.T) {
	tests := map[string]struct {
		segments                          []LineSegment[float64]
		segmentToDelete                   LineSegment[float64]
		p                                 point.Point[float64]
		expectedLeft, expectedRight       LineSegment[float64]
		expectedLeftNil, expectedRightNil bool
	}{
		"left and right neighbors": { // "\\//" shape
			segments: []LineSegment[float64]{
				New[float64](-4, 10, 1, 0),
				New[float64](0, 10, 5, 0),
				New[float64](8, 10, 10, 5),
				New[float64](20, 10, 15, 0),
				New[float64](24, 10, 19, 0),
			},
			segmentToDelete:  New[float64](8, 10, 10, 5),
			p:                point.New[float64](10, 5),
			expectedLeft:     New[float64](0, 10, 5, 0),
			expectedRight:    New[float64](20, 10, 15, 0),
			expectedLeftNil:  false,
			expectedRightNil: false,
		},
		"right only neighbor": { // "\\//" shape
			segments: []LineSegment[float64]{
				New[float64](8, 10, 10, 5),
				New[float64](20, 10, 15, 0),
				New[float64](24, 10, 19, 0),
			},
			segmentToDelete:  New[float64](8, 10, 10, 5),
			p:                point.New[float64](10, 5),
			expectedRight:    New[float64](20, 10, 15, 0),
			expectedLeftNil:  true,
			expectedRightNil: false,
		},
		"left only neighbor": { // "\\//" shape
			segments: []LineSegment[float64]{
				New[float64](-4, 10, 1, 0),
				New[float64](0, 10, 5, 0),
				New[float64](8, 10, 10, 5),
			},
			segmentToDelete:  New[float64](8, 10, 10, 5),
			p:                point.New[float64](10, 5),
			expectedLeft:     New[float64](0, 10, 5, 0),
			expectedLeftNil:  false,
			expectedRightNil: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// set up status structure
			S := newStatusStructureRBT(tc.p)

			// insert points into status structure
			for _, seg := range tc.segments {
				S.Insert(seg)
			}
			log.Printf("Status Structure before delete:\n%s", S.String())

			// delete segment
			S.Remove(tc.segmentToDelete)
			log.Printf("Status Structure after delete:\n%s", S.String())

			actualLeft, actualRight := S.FindNeighborsOfPoint(tc.p)

			// check left
			if tc.expectedLeftNil {
				assert.Nil(t, actualLeft, "expected left neighbor to be nil")
			} else {
				assert.Equal(t, tc.expectedLeft, *actualLeft, "left neighbor mismatch")
			}

			// check right
			if tc.expectedRightNil {
				assert.Nil(t, actualRight, "expected right neighbor to be nil")
			} else {
				assert.Equal(t, tc.expectedRight, *actualRight, "right neighbor mismatch")
			}
		})
	}
}

func TestStatusStructureRBT_FindNeighborsOfUofPAndCofP(t *testing.T) {
	tests := map[string]struct {
		segments                []LineSegment[float64]
		segmentToDelete         LineSegment[float64]
		UofP, CofP              []LineSegment[float64]
		p                       point.Point[float64]
		expectedSPrime          LineSegment[float64]
		expectedSPrimeNil       bool
		expectedSL              LineSegment[float64]
		expectedSLNil           bool
		expectedSDoublePrime    LineSegment[float64]
		expectedSDoublePrimeNil bool
		expectedSR              LineSegment[float64]
		expectedSRNil           bool
	}{
		"left and right neighbors": { // "\\//" shape
			segments: []LineSegment[float64]{
				New[float64](-4, 10, 1, 0),
				New[float64](0, 10, 5, 0),
				New[float64](8, 10, 10, 5),
				New[float64](20, 10, 15, 0),
				New[float64](24, 10, 19, 0),
				New[float64](12, 7, 8, 3),
				New[float64](12, 6, 8, 4),
			},
			segmentToDelete: New[float64](8, 10, 10, 5),
			UofP:            nil,
			CofP: []LineSegment[float64]{
				New[float64](12, 7, 8, 3),
				New[float64](12, 6, 8, 4),
			},
			p:                       point.New[float64](10, 5),
			expectedSPrime:          New[float64](12, 6, 8, 4),
			expectedSPrimeNil:       false,
			expectedSL:              New[float64](0, 10, 5, 0),
			expectedSLNil:           false,
			expectedSDoublePrime:    New[float64](12, 7, 8, 3),
			expectedSDoublePrimeNil: false,
			expectedSR:              New[float64](20, 10, 15, 0),
			expectedSRNil:           false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// set up status structure
			S := newStatusStructureRBT(tc.p)

			// insert points into status structure
			for _, seg := range tc.segments {
				S.Insert(seg)
			}
			log.Printf("Status Structure before delete:\n%s", S.String())

			// delete point
			S.Remove(tc.segmentToDelete)

			log.Printf("Status Structure after delete:\n%s", S.String())

			sPrime, sL, sDoublePrime, sR := S.FindNeighborsOfUofPAndCofP(tc.UofP, tc.CofP)

			t.Log("sPrime:", sPrime)
			t.Log("sL:", sL)
			t.Log("sDoublePrime:", sDoublePrime)
			t.Log("sR:", sR)

			if tc.expectedSPrimeNil {
				assert.Nil(t, sPrime, "expected sPrime to be nil")
			} else {
				assert.Equal(t, tc.expectedSPrime, *sPrime, "sPrime mismatch")
			}

			if tc.expectedSLNil {
				assert.Nil(t, sL, "expected sL to be nil")
			} else {
				assert.Equal(t, tc.expectedSL, *sL, "sL mismatch")
			}

			if tc.expectedSDoublePrimeNil {
				assert.Nil(t, sDoublePrime, "expected sDoublePrime to be nil")
			} else {
				assert.Equal(t, tc.expectedSDoublePrime, *sDoublePrime, "sDoublePrime mismatch")
			}

			if tc.expectedSRNil {
				assert.Nil(t, sR, "expected sR to be nil")
			} else {
				assert.Equal(t, tc.expectedSR, *sR, "sR mismatch")
			}
		})
	}
}

func TestStatusStructureRBT_Remove(t *testing.T) {
	tests := map[string]struct {
		segments        []LineSegment[float64]
		segmentToDelete LineSegment[float64]
		p               point.Point[float64]
		expected        []LineSegment[float64]
	}{
		"W type shape": {
			segments: []LineSegment[float64]{
				New[float64](0, 10, 11, -1),
				New[float64](10, 10, 4, 4),
				New[float64](20, 10, 9, -1),
			},
			segmentToDelete: New[float64](10, 10, 4, 4),
			p:               point.New[float64](4, 4),
			expected: []LineSegment[float64]{
				New[float64](0, 10, 11, -1),
				New[float64](20, 10, 9, -1),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// set up status structure
			S := newStatusStructureRBT(tc.p)

			// insert points into status structure
			for _, seg := range tc.segments {
				S.Insert(seg)
			}
			log.Printf("Status Structure before delete:\n%s", S.String())

			// delete point
			S.Remove(tc.segmentToDelete)

			log.Printf("Status Structure after delete:\n%s", S.String())

			// check output matches
			iter := S.structure.Iterator()
			for n, expected := range tc.expected {
				t.Logf("checking index %d", n)
				found := iter.Next()
				require.True(t, found, "could not pop from S")
				actual := iter.Key().(statusStructureEntry).segment
				t.Log("popped:", actual)
				t.Log("expect:", expected)
				assert.Equal(t, expected, actual, "segment mismatch")
			}
		})
	}
}
