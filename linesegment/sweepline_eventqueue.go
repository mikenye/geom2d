package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"log"
	"slices"
	"strings"
)

// qItem represents an entry in the event queue used in the sweep line algorithm for finding line segment intersections.
//
// Each qItem consists of:
//   - point: The event point where an intersection or segment endpoint occurs.
//   - segments: A slice of LineSegment[float64] that are associated with this event point. These could be:
//   - Segments that start at this point (U(p))
//   - Segments that end at this point (L(p))
//   - Segments that contain this point as an intersection (C(p))
//
// Purpose:
//   - qItem structures the event queue, ensuring each event is processed correctly in handleEventPoint.
//   - It allows multiple segments to be grouped under the same event point, preventing redundant queue entries.
//
// Notes:
//   - qItem is used internally by the sweep line algorithm and is not exposed publicly.
//   - The event queue processes qItem entries in lexicographic order, breaking ties as needed.
//
// This ensures that events are processed in the correct order during the sweep.
type qItem struct {
	point    point.Point[float64]
	segments []LineSegment[float64]
}

// String returns a human-readable representation of the event queue item (qItem).
//
// The format follows:
//
//	Queue Item: (x, y), U(p): (x1,y1)(x2,y2), (x3,y3)(x4,y4), ...
//
// This output helps visualize event processing in the sweep line algorithm.
// The event point p is displayed first, followed by all associated segments "U(p)" where the
// event point is the upper point of the line segment.
//
// Example Output:
//
//	Queue Item: (5,10), U(p): (5,10)(15,20), (5,10)(7,14)
//
// This function is particularly useful for debugging and logging the event queue.
func (qi qItem) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Queue Item: %s, U(p): ", qi.point.String()))
	first := true
	for _, seg := range qi.segments {
		if first {
			builder.WriteString(seg.String())
			first = false
			continue
		}
		builder.WriteString(fmt.Sprintf(", %s", seg.String()))
	}
	return builder.String()
}

func debugPrintQueue(Q *btree.BTreeG[qItem]) {
	Qcopy := Q.Clone()
	log.Println("Event queue (Q):")
	for Qcopy.Len() > 0 {
		item, _ := Qcopy.DeleteMin()
		log.Printf("  - %s\n", item.String())
	}
}

func mergeSegments(a, b []LineSegment[float64], opts ...options.GeometryOptionsFunc) []LineSegment[float64] {
	input := append(a, b...)
	output := make([]LineSegment[float64], 0, len(a)+len(b))
	for _, seg := range input {
		if !slices.ContainsFunc(output, func(l LineSegment[float64]) bool {
			return l.Eq(seg, opts...)
		}) {
			output = append(output, seg)
		}
	}
	return output
}

func newEventQueue[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) *btree.BTreeG[qItem] {

	// Initialize an empty event queue
	Q := btree.NewG[qItem](2, qItemLess)

	// Insert the segment endpoints into Q.
	// When an upper endpoint is inserted, the corresponding segment should be stored with it.
	for _, seg := range segments {

		segf := seg.AsFloat64()
		upper, upperDegenerate := segf.sweeplineUpperPoint()
		lower, lowerDegenerate := segf.sweeplineLowerPoint()

		// skip degenerate points
		if upperDegenerate || lowerDegenerate {
			continue
		}

		// create new items
		upperQItem := qItem{
			point:    upper,
			segments: []LineSegment[float64]{NewFromPoints(upper, lower)},
		}
		lowerQItem := qItem{
			point:    lower,
			segments: nil,
		}

		// add the events to the event queue
		for _, newItem := range []qItem{upperQItem, lowerQItem} {

			// check if upper point is in the queue already
			existingItem, exists := Q.Get(newItem)

			// if it does, merge line segments
			if exists {
				newItem.segments = mergeSegments(newItem.segments, existingItem.segments, opts...)
			}

			// add the new event
			Q.ReplaceOrInsert(newItem)
		}
	}

	return Q
}

// qItemLess defines the ordering of event queue items (qItem) for use in a balanced B-tree (btree.BTreeG[qItem]).
// It ensures that events are processed in the correct order by the sweep line algorithm.
//
// Ordering Rules:
//  1. Higher Y-coordinates come first (processed earlier by the sweep line).
//  2. For equal Y-coordinates, smaller X-coordinates come first (ensuring left-to-right processing).
//
// Parameters:
//   - p: The first event queue item (qItem) to compare.
//   - q: The second event queue item (qItem) to compare.
//
// Returns:
//   - true if p should be processed before q, otherwise false.
//
// Usage in the Sweep Line Algorithm:
//   - The event queue (Q) is implemented as a balanced B-tree (btree.BTreeG[qItem]).
//   - This function defines the comparison rule for inserting and retrieving events efficiently.
func qItemLess(p, q qItem) bool {
	// Compare based on the sweep line event order:
	// 1. Higher y-coordinates are "smaller" (processed first).
	// 2. For equal y-coordinates, smaller x-coordinates are "smaller".
	if p.point.Y() > q.point.Y() {
		return true
	}
	if p.point.Y() == q.point.Y() && p.point.X() < q.point.X() {
		return true
	}
	return false
}
