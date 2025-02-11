package linesegment

//func TestSegmentSortLessHigherOrder(t *testing.T) {
//	tests := map[string]struct {
//		setupStatusStructure func() *btree.BTreeG[sItem]
//		segments             []LineSegment[float64]
//		expected             []LineSegment[float64]
//	}{
//		"vertical & negative slope": { // vertical should come first
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point:    point.New[float64](0, 0),
//					segments: []LineSegment[float64]{New[float64](0, 0, 29, -11)},
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](-81, 89, 10, 0),
//				New[float64](10, 109, 10, -66),
//			},
//			// The order of the segments should correspond to the order in which they are intersected
//			// by a sweep line just below the event point.
//			// If there is a horizontal segment, it comes last among all segments containing the event point.
//			// Thus, in this case:
//			//   - The vertical line should come first
//			//   - The diagnonal line should be second, as it proceeds through the point to be slightly to the right
//			//     of the point just below the event point.
//			expected: []LineSegment[float64]{
//				New[float64](10, 109, 10, -66),
//				New[float64](-81, 89, 10, 0),
//			},
//		},
//		"vertical & positive slope": { // vertical should come last
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](0, 0),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](0, 10, 0, -10),
//				New[float64](10, 10, -10, -10),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](10, 10, -10, -10),
//				New[float64](0, 10, 0, -10),
//			},
//		},
//		"vertical & positive slope #2": { // vertical should come last
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](19, 10),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](10, 10, 10, -56),
//				New[float64](10, 10, 0, -99),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](10, 10, 0, -99),
//				New[float64](10, 10, 10, -56),
//			},
//		},
//		"vertical & horizontal slope": { // horizontal should come last
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](0, 0),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](-10, 0, 10, 0),
//				New[float64](0, 10, 0, -10),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](0, 10, 0, -10),
//				New[float64](-10, 0, 10, 0),
//			},
//		},
//		"both diagonal, negative slopes": { // steeper slope should come first
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](0, 0),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](-10, 10, 10, -10),
//				New[float64](-5, 10, 5, -10),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](-5, 10, 5, -10),
//				New[float64](-10, 10, 10, -10),
//			},
//		},
//		"both diagonal, positive slopes": { // steeper slope should come first
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](0, 0),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](10, 10, -10, -10),
//				New[float64](5, 10, -5, -10),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](10, 10, -10, -10),
//				New[float64](5, 10, -5, -10),
//			},
//		},
//		"both diagonal, opposing slopes": { // negative slope should come first
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point: point.New[float64](0, 0),
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](-10, 10, 10, -10),
//				New[float64](10, 10, -10, -10),
//			},
//			expected: []LineSegment[float64]{
//				New[float64](10, 10, -10, -10),
//				New[float64](-10, 10, 10, -10),
//			},
//		},
//		"ordering of collinear segments": {
//			setupStatusStructure: func() *btree.BTreeG[sItem] {
//				// Initialize status structure
//				var StatusStructure *btree.BTreeG[sItem]
//
//				// Event where failure happens
//				event := qItem{
//					point:    point.New[float64](7, 7),
//					segments: []LineSegment[float64]{New[float64](7, 7, 3, 3)},
//				}
//
//				// Update the status structure based on new sweepline position
//				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
//			},
//			segments: []LineSegment[float64]{
//				New[float64](7, 7, 3, 3),
//				New[float64](10, 10, 0, 0),
//			},
//		},
//	}
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//
//			subTests := []string{"normal", "input segment order reversed"}
//			for i, subName := range subTests {
//				t.Run(subName, func(t *testing.T) {
//					// setup
//					StatusStructure := tc.setupStatusStructure()
//					debugStatusStructure(StatusStructure)
//
//					// reverse slices
//					if i == 1 {
//						slices.Reverse(tc.segments)
//					}
//
//					// add segments
//					for _, seg := range tc.segments {
//						StatusStructure.ReplaceOrInsert(sItem{segment: seg})
//					}
//
//					log.Println("Status Structure:")
//					debugStatusStructure(StatusStructure)
//
//					// check order matches
//					for _, expected := range tc.expected {
//						actual, exists := StatusStructure.DeleteMin()
//						require.True(t, exists, "StatusStructure unexpectedly empty")
//						assert.Equal(t, expected, actual.segment)
//					}
//				})
//			}
//		})
//	}
//}
