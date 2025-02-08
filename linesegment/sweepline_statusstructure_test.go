package linesegment

import (
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSegmentSortLessHigherOrder(t *testing.T) {
	tests := map[string]struct {
		setupStatusStructure func() *btree.BTreeG[sItem]
		segments             []LineSegment[float64]
		expected             []LineSegment[float64]
	}{
		"ordering of collinear segments": {
			setupStatusStructure: func() *btree.BTreeG[sItem] {
				// Initialize status structure
				var StatusStructure *btree.BTreeG[sItem]
				StatusStructure = nil

				// Event where failure happens
				event := qItem{
					point:    point.New[float64](7, 7),
					segments: []LineSegment[float64]{New[float64](7, 7, 3, 3)},
				}

				// Update the status structure based on new sweepline position
				return updateStatusStructure(StatusStructure, event, options.WithEpsilon(1e-8))
			},
			segments: []LineSegment[float64]{
				New[float64](7, 7, 3, 3),
				New[float64](10, 10, 0, 0),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// setup
			StatusStructure := tc.setupStatusStructure()
			debugStatusStructure(StatusStructure)

			// add segments
			for _, seg := range tc.segments {
				StatusStructure.ReplaceOrInsert(sItem{segment: seg})
			}
			debugStatusStructure(StatusStructure)

			// check order matches
			for _, expected := range tc.expected {
				actual, exists := StatusStructure.DeleteMin()
				require.True(t, exists, "StatusStructure unexpectedly empty")
				assert.Equal(t, expected, actual.segment)
			}
		})
	}
}
