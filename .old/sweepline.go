package _old

import (
	"container/heap"
	"fmt"
	"github.com/mikenye/geom2d/types"
	"log"
	"math"
	"slices"
	"strings"
)

// sweeplineEventType represents the type of event in the sweep line algorithm.
type sweeplineEventType uint8

// Valid values of sweeplineEventType
const (
	// sweeplineEventStart indicates the event where a line segment begins.
	// This event is triggered when the sweep line encounters the starting point of a line segment.
	sweeplineEventStart sweeplineEventType = iota

	// sweeplineEventEnd indicates the event where a line segment ends.
	// This event is triggered when the sweep line encounters the ending point of a line segment.
	sweeplineEventEnd

	// sweeplineEventIntersection indicates the event where two or more line segments intersect.
	// This event is triggered when the sweep line detects an intersection point between line segments.
	sweeplineEventIntersection
)

// String returns the string representation of a sweeplineEventType.
// This method is used primarily for debugging and logging purposes to
// provide a human-readable description of the event type.
//
// Returns:
//   - "sweeplineEventStart" for sweeplineEventStart.
//   - "sweeplineEventEnd" for sweeplineEventEnd.
//   - "sweeplineEventIntersection" for sweeplineEventIntersection.
//
// Panics:
//   - If the sweeplineEventType has an unsupported or unknown value, this function
//     will panic with an error message indicating the invalid type.
func (t *sweeplineEventType) String() string {
	switch *t {
	case sweeplineEventStart:
		return "sweeplineEventStart"
	case sweeplineEventEnd:
		return "sweeplineEventEnd"
	case sweeplineEventIntersection:
		return "sweeplineEventIntersection"
	default:
		panic(fmt.Errorf("unsupported sweeplineEventType"))
	}
}

// sweeplineEvent represents an event in the sweep line algorithm used to process line segments
// for geometric computations such as finding intersections. Each event corresponds to a specific
// point on the sweep line and a set of associated line segments.
type sweeplineEvent struct {

	// The line segments associated with this event. The specific segments
	// depend on the eventType. For example:
	//   - For a start (sweeplineEventStart) event, it contains the segment starting at the point.
	//   - For an end (sweeplineEventEnd) event, it contains the segment ending at the point.
	//   - For an intersection (sweeplineEventIntersection) event, it contains the intersecting segments.
	lineSegments []LineSegment[float64]

	// The point on the plane where the event occurs. This represents the X-Y coordinates
	// associated with the event.
	point Point[float64]

	// The type of event, which determines the role of this point in the sweep line
	// process. Possible values are:
	//   - sweeplineEventStart: A start point of a line segment.
	//   - sweeplineEventEnd: An end point of a line segment.
	//   - sweeplineEventIntersection: An intersection point between two or more segments.
	eventType sweeplineEventType

	// A unique identifier for the insertion order of the event. This is used to
	// break ties during sorting or comparisons in the event queue to ensure a deterministic order
	// when multiple events share the same coordinates and type.
	insertionOrder uint64

	// The index of the event in the priority queue. This is managed internally by the
	// heap.Interface methods and is used for efficient updates and removals from the queue.
	index int
}

// String returns a human-readable representation of a sweeplineEvent.
// The output includes the point where the event occurs, the event type,
// and the associated line segments. This is useful for debugging and logging.
//
// The output format is as follows:
//
//	 event at: [Point], type: [EventType], line segments:
//		[LineSegment1]
//		[LineSegment2] (if it exists)
//		...
//
// Example:
//
// If the event is an intersection at Point[(3, 4)] between two line segments:
// LineSegment[(1, 1) -> (5, 5)] and LineSegment[(1, 5) -> (5, 1)],
// the output will look like:
//
//	 event at: Point[(3, 4)], type: sweeplineEventIntersection, line segments:
//		LineSegment[(1, 1) -> (5, 5)]
//		LineSegment[(1, 5) -> (5, 1)]
func (e *sweeplineEvent) String() string {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("event at: %s, type: %s, line segments:\n", e.point.String(), e.eventType.String()))
	for _, l := range e.lineSegments {
		out.WriteString(fmt.Sprintf("\t%s\n", l.String()))
	}
	return out.String()
}

// sweeplineEventQueue represents a priority queue of sweepline events used in the
// Bentley-Ottmann algorithm. It maintains the order of events based on their
// position, type, and insertion order for tie-breaking.
//
// The priority queue is implemented as a min-heap, with events ordered by their
// x-coordinate, y-coordinate, event type, and insertion order. This ensures that
// events are processed in the correct sequence.
//
// Usage:
//
// The sweeplineEventQueue is designed to be used with the `heap` package, which
// provides efficient heap operations such as insertion, removal, and reordering.
//
// Example:
//
//	q := &sweeplineEventQueue{}
//	heap.Init(q)
//	heap.Push(q, &sweeplineEvent{...})
//	nextEvent := heap.Pop(q).(*sweeplineEvent)
type sweeplineEventQueue struct {

	// A slice of pointers to sweeplineEvent objects that make up the heap.
	events []*sweeplineEvent

	// A counter for tracking the insertion order of events. This is used
	// to break ties when events have the same coordinates and type.
	insertions uint64
}

// AddIntersectionEvent determines the intersection between two line segments A and B.
// If an intersection exists, it creates an appropriate intersection event and adds it to the sweepline event queue.
//
// Parameters:
//   - A, B ([LineSegment][float64]): Line segments to test for intersection.
//   - opts: Optional parameters that may influence the behaviour of the intersection calculation.
//
// Behaviour:
//   - If there is no intersection ([IntersectionNone]), the function does nothing.
//   - If the intersection is a single point ([IntersectionPoint]), an intersection event is created and added to the queue.
//   - If the intersection is a line segment ([IntersectionSegment]), the function will panic, as this case is not yet implemented.
//
// Notes:
//   - The intersection calculation is delegated to the [LineSegment.IntersectionGeometry] method of the [LineSegment] type.
//   - The event is pushed onto the priority queue, maintaining the correct event order for the Bentley-Ottmann algorithm.
//
// Panics:
// - This function panics if the intersection type is [IntersectionSegment], as that scenario is not yet implemented.
func (pq *sweeplineEventQueue) AddIntersectionEvent(A, B LineSegment[float64], opts ...Option) {
	intersection := A.IntersectionGeometry(B, opts...)

	switch intersection.IntersectionType {
	case IntersectionNone:
		// do nothing
	case IntersectionPoint:
		// add the intersection point to the queue
		event := &sweeplineEvent{
			lineSegments: []LineSegment[float64]{A, B},
			point:        intersection.IntersectionPoint,
			eventType:    sweeplineEventIntersection,
		}
		heap.Push(pq, event)

	case IntersectionSegment:
		panic(fmt.Errorf("not yet implemented"))
	}
}

// RemoveIntersection removes an intersection event involving the specified line segments A and B from the sweepline event queue.
//
// Parameters:
//   - A, B ([LineSegment][float64]): The line segments forming the intersection to be removed.
//   - opts: Optional parameters that may influence the behaviour of equality checks between line segments.
//
// Behaviour:
//   - The function iterates through the events in the queue, looking for intersection events involving both line segments A and B.
//   - If a matching intersection event is found, it is removed from the queue.
//   - After all matching events are removed, the queue is re-sorted to maintain the correct order.
//
// Notes:
//
//   - This function is specifically designed to work with intersection events (sweeplineEventIntersection).
//   - If no matching intersection is found, the function does nothing.
//   - Uses slices.ContainsFunc to check for matching line segments in the event.
//   - Ensure the sweeplineEventQueue is initialized with the heap package: heap.Init(pq)
func (pq *sweeplineEventQueue) RemoveIntersection(A, B LineSegment[float64], opts ...Option) {
	queueNeedsSorting := false // Flag to indicate if the queue needs to be re-sorted after removal.

	// Iterate through all events in the queue.
	for i := 0; i < len(pq.events); i++ {
		// Skip non-intersection events.
		if pq.events[i].eventType != sweeplineEventIntersection {
			continue
		}

		// Check if the current event contains line segment A.
		foundA := slices.ContainsFunc(pq.events[i].lineSegments, func(l LineSegment[float64]) bool {
			return l.Eq(A, opts...)
		})

		// Check if the current event contains line segment B.
		foundB := slices.ContainsFunc(pq.events[i].lineSegments, func(l LineSegment[float64]) bool {
			return l.Eq(B, opts...)
		})

		// If either segment is not found in the current event, skip it.
		if !foundA || !foundB {
			continue
		}

		// Log the removal of the intersection event.
		log.Printf("[queue] removing event: %s", pq.events[i].String())

		// Swap the matching event with the last event in the slice and shrink the slice by one.
		(*pq).events[i], (*pq).events[len(pq.events)-1] = (*pq).events[len(pq.events)-1], (*pq).events[i]
		pq.events = (*pq).events[:len(pq.events)-1]

		// Adjust the index to account for the removed event.
		i--

		// Mark that the queue will need re-sorting after removal.
		queueNeedsSorting = true
	}

	// If any events were removed, re-initialize the heap to restore correct order.
	if queueNeedsSorting {
		heap.Init(pq)
	}
}

// Len returns the number of events currently in the sweepline event queue.
//
// This method implements the [heap.Interface] requirement for determining
// the length of the heap, which is crucial for maintaining the heap's structure.
//
// Returns:
// - The number of events in the queue as an integer.
func (pq *sweeplineEventQueue) Len() int {
	return len(pq.events)
}

// Less determines the ordering of two events in the sweepline event queue.
//
// This method implements the [heap.Interface] requirement for ordering elements
// in the heap. It uses a series of comparisons to prioritize events based on
// their X-coordinate, Y-coordinate, type, and insertion order, in that sequence.
//
// Parameters:
//   - i (int): The index of the first event in the queue to compare.
//   - j (int): The index of the second event in the queue to compare.
//
// Returns:
//   - A boolean indicating whether the event at index i should be ordered before
//     the event at index j.
//
// The comparison logic follows these rules:
//  1. X-coordinate: Events are ordered by their X-coordinate (ascending).
//  2. Y-coordinate: If X-coordinates are equal, events are ordered by their Y-coordinate (ascending).
//  3. Event Type: If both coordinates are equal, events are ordered by their type
//     (sweeplineEventStart < sweeplineEventEnd < sweeplineEventIntersection).
//  4. Insertion Order: As a final tiebreaker, events are ordered by their insertion order.
func (pq *sweeplineEventQueue) Less(i, j int) bool {

	// Compare X-coordinates
	if pq.events[i].point.x != pq.events[j].point.x {
		return pq.events[i].point.x < pq.events[j].point.x
	}

	// Compare Y-coordinates
	if pq.events[i].point.y != pq.events[j].point.y {
		return pq.events[i].point.y < pq.events[j].point.y
	}

	// Compare event types
	if pq.events[i].eventType != pq.events[j].eventType {
		return pq.events[i].eventType < pq.events[j].eventType
	}

	// Tiebreaker: Compare insertion order
	return pq.events[i].insertionOrder < pq.events[j].insertionOrder
}

// Push adds a new event to the sweepline event queue.
//
// This method implements the [heap.Interface] requirement for pushing an element
// onto the heap. It sets the index of the event, updates the insertion count,
// assigns an insertion order to maintain tie-breaking stability, and appends the
// event to the internal slice of events.
//
// Parameters:
//   - x (*sweeplineEvent): The event to push onto the queue.
func (pq *sweeplineEventQueue) Push(x any) {
	// Cast the input to *sweeplineEvent
	item := x.(*sweeplineEvent)

	// Assign the index of the new event to its current position in the events slice
	item.index = len(pq.events)

	// Increment the insertion counter for stable ordering of tie-breaking cases
	pq.insertions++

	// Assign the current insertion count to the event as its insertion order
	item.insertionOrder = pq.insertions

	// Log the addition of the new event
	log.Printf("[queue] pushing event: %s", item.String())

	// Append the event to the slice of events
	pq.events = append(pq.events, item)
}

// Pop removes and returns the last event from the sweepline event queue.
//
// This method implements the [heap.Interface] requirement for popping an element
// from the heap. It removes the last element in the queue, reduces the slice length,
// logs the operation, and returns the removed event.
//
// Returns:
//   - any: The event that was removed from the queue, of type *sweeplineEvent.
//
// Example Usage:
//
//	event := pq.Pop().(*sweeplineEvent)
func (pq *sweeplineEventQueue) Pop() any {
	// Store the current events slice
	old := pq.events

	// Get the number of events in the queue
	n := len(old)

	// Get the last item from the slice
	item := old[n-1]

	// Reduce the slice to exclude the last item
	pq.events = old[0 : n-1]

	// Log the event being popped
	log.Printf("[queue] popped event: %s", item.String())

	// Return the popped event
	return item
}

// Swap exchanges two events in the sweepline event queue by their indices.
//
// This method implements the [heap.Interface] requirement for swapping two elements
// in the heap. It updates the positions of the events in the underlying slice and
// ensures their index fields are correctly updated to maintain consistency.
//
// Parameters:
//   - i (int): The index of the first event to swap.
//   - j (int): The index of the second event to swap.
func (pq *sweeplineEventQueue) Swap(i, j int) {
	// Swap the events at indices i and j in the slice.
	(*pq).events[i], (*pq).events[j] = (*pq).events[j], (*pq).events[i]

	// Update the index fields of the swapped events to reflect their new positions.
	(*pq).events[i].index = i
	(*pq).events[j].index = j
}

// newSweeplineEventQueue initializes a new sweepline event queue for the given line segments.
//
// This function processes an array of line segments, ensuring each segment is normalized
// so that the leftmost-lowest point is always the starting point. For each line segment,
// two events (start and end) are created and added to the queue. The event queue is then
// heap-initialized to ensure correct sorting for the Bentley-Ottmann algorithm.
//
// Parameters:
//   - lineSegments ([][LineSegment][float64]): The line segments for which the event queue is created.
//
// Returns:
//   - *sweeplineEventQueue: A pointer to the initialized sweepline event queue.
func newSweeplineEventQueue(lineSegments []LineSegment[float64]) *sweeplineEventQueue {

	// Initialize the sweepline event queue with pre-allocated memory.
	// Each line segment contributes two events: one for its start point and one for its end point.
	pq := &sweeplineEventQueue{
		events: make([]*sweeplineEvent, len(lineSegments)*2), // Pre-allocate for efficiency.
	}

	// Iterate through the given line segments.
	for _, l := range lineSegments {

		// Normalize the line segment to ensure the start point is the leftmost-lowest.
		l = l.Normalize()

		// Add the start and end points of the line segment as events in the queue.
		for _, p := range l.Points() {
			pq.events[pq.insertions] = &sweeplineEvent{
				lineSegments: []LineSegment[float64]{l},             // Associate the line segment with the event.
				point:        p,                                     // The event's point.
				eventType:    sweeplineEventType(pq.insertions % 2), // Alternates between start (0) and end (1).
			}
			pq.insertions++ // update insertion order so tie breaking works
		}
	}

	// Initialize the heap to ensure proper ordering of events.
	heap.Init(pq)

	// Return the initialized event queue.
	return pq
}

// SweepLineResult represents the outcome of processing a set of line segments
// using the sweep line algorithm.
//
// This structure captures the resulting line segments and any intersection points
// identified during the processing. It is typically returned as the final result
// of the sweep line algorithm.
type SweepLineResult struct {

	// The line segments that were processed.
	// These may include the original segments as well as new segments created from splitting
	// at intersection points.
	LineSegments []LineSegment[float64]

	// IntersectionPoints holds the points where intersections occurred between the line segments.
	IntersectionPoints []Point[float64]
}

// sweepLineStatusEntry represents an entry in the sweep line status data structure.
//
// Each entry corresponds to a line segment that is currently active in the sweep line algorithm.
// The sweep line status maintains the active line segments at a specific position along the x-axis,
// sorted by their y-coordinate at the current x-position. This struct serves as a wrapper for a
// line segment within the sweep line status.
type sweepLineStatusEntry struct {

	// lineSegment is the line segment represented by this entry.
	lineSegment LineSegment[float64]
}

// sweepLineStatus represents the active line segments in the sweep line algorithm.
//
// The sweep line status keeps track of the line segments that are currently intersected by the
// vertical sweep line at a specific x-coordinate. It provides functionality for managing these
// active segments, including insertion, removal, and neighbor queries.
type sweepLineStatus struct {

	// entries is the list of active line segments in the sweep line status, sorted by their
	// y-coordinates at the current x-coordinate.
	entries []sweepLineStatusEntry

	// currentX represents the x-coordinate of the sweep line, used to determine the y-coordinates
	// of the active line segments and maintain their order.
	currentX float64
}

// Index searches for the specified line segment in the sweepLineStatus and returns its index.
//
// This method iterates through the entries in the sweepLineStatus to find the given line segment.
// If the segment is found, it returns its index; otherwise, it returns -1.
//
// Parameters:
//   - segment ([LineSegment][float64]): The line segment to search for in the sweepLineStatus.
//
// Returns:
//   - (int): The index of the line segment in the `entries` slice if found, or -1 if not found.
func (sls *sweepLineStatus) Index(segment LineSegment[float64]) int {
	for i, l := range sls.entries {
		if l.lineSegment.Eq(segment) {
			return i
		}
	}
	return -1
}

// Insert adds a new line segment to the sweepLineStatus and ensures the entries remain sorted.
//
// This method appends the provided line segment to the list of entries and then sorts the list
// to maintain the correct order for sweep line processing.
//
// Parameters:
//   - segment ([LineSegment][float64]): The line segment to add to the sweepLineStatus.
func (sls *sweepLineStatus) Insert(segment LineSegment[float64]) {
	entry := sweepLineStatusEntry{
		lineSegment: segment,
	}
	log.Printf("[status] inserting segment: %s", entry.lineSegment.String())
	sls.entries = append(sls.entries, entry)
	sls.Sort()
}

// Neighbors returns the entries immediately above and below the given line segment in the sweepLineStatus.
//
// The method finds the index of the specified line segment in the `entries` slice and retrieves its neighbors,
// if they exist. Neighbors are determined by their position in the sorted slice of entries.
//
// Parameters:
//   - segment ([LineSegment][float64]): The line segment for which neighbors are to be found.
//
// Returns:
//   - above (*sweepLineStatusEntry): The entry directly above the given line segment, or nil if none exists.
//   - below (*sweepLineStatusEntry): The entry directly below the given line segment, or nil if none exists.
//
// Behavior:
//   - If the specified segment is not in the status, both `above` and `below` will be nil.
//   - If the segment is the first in the list, `below` will be nil.
//   - If the segment is the last in the list, `above` will be nil.
func (sls *sweepLineStatus) Neighbors(segment LineSegment[float64]) (above, below *sweepLineStatusEntry) {
	i := sls.Index(segment)
	if i == -1 {
		return nil, nil
	}
	if i > 0 {
		below = &sls.entries[i-1]
	}
	if i < len(sls.entries)-1 {
		above = &sls.entries[i+1]
	}
	return above, below
}

// Remove removes a given line segment from the sweepLineStatus.
//
// The method searches for the specified line segment in the entries slice.
// If the segment is found, it is removed.
// If the segment is not found, the method returns without making any changes.
//
// Parameters:
// - segment ([LineSegment][float64]): The line segment to be removed.
func (sls *sweepLineStatus) Remove(segment LineSegment[float64]) {
	i := sls.Index(segment)
	if i == -1 {
		return
	}
	log.Printf("[status] removing segment: %s", segment.String())
	sls.entries = append(sls.entries[:i], sls.entries[i+1:]...)
}

// Sort ensures that the entries slice in the sweepLineStatus is sorted in a stable order.
//
// The sorting criteria are as follows:
//  1. The y-coordinate at the current sweep line x-position (currentX) is compared first.
//     If both line segments are defined at currentX (i.e., YAtX succeeds),
//     they are ordered by their y values.
//  2. If the y-coordinates are equal or undefined for both line segments,
//     the starting x-coordinate of the segments is compared.
//     Segments with smaller starting x values are sorted first.
//  3. If both y and starting x-coordinates are equal, the segments are considered equal.
//
// Behavior:
//   - Uses slices.SortStableFunc to sort the entries slice based on the above criteria.
//   - Stability ensures that ties (segments considered equal) retain their relative order.
func (sls *sweepLineStatus) Sort() {
	slices.SortStableFunc(sls.entries, func(a, b sweepLineStatusEntry) int {
		ay, aOk := a.lineSegment.YAtX(sls.currentX)
		by, bOk := b.lineSegment.YAtX(sls.currentX)

		// Compare by y-coordinate at sweep line's current X if both segments are defined
		if aOk && bOk && ay != by {
			switch {
			case ay < by:
				return -1
			case ay > by:
				return 1
			}
		}

		// Compare by start x-coordinate
		if a.lineSegment.start.x != b.lineSegment.start.x {
			switch {
			case a.lineSegment.start.x < b.lineSegment.start.x:
				return -1
			case a.lineSegment.start.x > b.lineSegment.start.x:
				return 1
			}
		}

		// If all comparisons result in equality
		return 0
	})
}

// String generates a human-readable representation of the current state of the sweepLineStatus.
//
// The output includes:
// 1. The current x-position of the sweep line (currentX).
// 2. The list of active line segments (entries), ordered as they are in the sweep line status.
//
// Each line segment is represented in its string format, indented for better readability.
//
// Example:
//
// If the sweepLineStatus contains the following data:
//
//	currentX: 2.5
//	entries: [
//	    LineSegment[(-1, 3) -> (4, 5)],
//	    LineSegment[(0, 2) -> (3, 6)],
//	]
//
// Then the output will look like:
//
//	sweep line at x=2.5000:
//	    LineSegment[(-1, 3) -> (4, 5)]
//	    LineSegment[(0, 2) -> (3, 6)]
func (sls *sweepLineStatus) String() string {
	s := strings.Builder{}

	// Write the current sweep line position
	s.WriteString(fmt.Sprintf("sweep line at x=%0.4f:\n", sls.currentX))

	// Write each line segment in the status
	for i := range sls.entries {
		s.WriteString(fmt.Sprintf("\t%s\n", sls.entries[i].lineSegment.String()))
	}

	return s.String()
}

// Swap exchanges the positions of two line segments, A and B, in the sweepLineStatus.
//
// This method searches for the indices of A and B within the entries slice.
// If both segments are found, their positions in the slice are swapped.
//
// Parameters:
//
//   - A (LineSegment[float64]): The first line segment to swap.
//   - B (LineSegment[float64]): The second line segment to swap.
//
// Behavior:
//   - If either A or B is not found in the entries, the method returns without making changes.
//   - If both segments are found, their positions in the entries slice are exchanged.
//
// Example:
//
// Given the following `sweepLineStatus`:
//
//	entries: [L1, L2, L3]
//
// Calling `Swap(L1, L2)` will result in:
//
//	entries: [L2, L1, L3]
func (sls *sweepLineStatus) Swap(A, B LineSegment[float64]) {
	// Find the indices of A and B in the entries slice
	iA := sls.Index(A)
	iB := sls.Index(B)

	// If either segment is not found, return without changes
	if iA == -1 || iB == -1 {
		return
	}

	// Swap the two entries
	sls.entries[iA], sls.entries[iB] = sls.entries[iB], sls.entries[iA]
}

// SweepLine performs the Bentley-Ottmann algorithm to find intersections between a set of line segments.
//
// This function detects all intersection points between the provided line segments efficiently using a sweep line
// approach. The algorithm operates by sorting and processing events (start points, end points, and intersections)
// and maintaining a dynamic status of active line segments intersecting the current sweep line.
//
// Parameters:
//   - linesegments: A slice of [LineSegment][T], where T is a [SignedNumber] type.
//   - opts: Optional configuration parameters.
//
// Returns:
//
// A SweepLineResult containing:
//   - The intersection points detected between line segments.
//   - A copy of the line segments (as float64).
//
// Algorithm Steps:
//
//  1. Convert all input line segments to [LineSegment][float64].
//  2. Initialize a priority queue (event queue) with start and end events for all line segments.
//  3. Process each event in the event queue in ascending order of x-coordinate (with ties resolved by y-coordinate).
//  4. Maintain a status structure (`sweepLineStatus`) to store active line segments intersecting the current sweep line.
//  5. For each event:
//     Start: Insert the segment into the sweep line status and check for intersections with neighbors.
//     End: Remove the segment from the sweep line status and check for new intersections between its former neighbors.
//     Intersection: Record the intersection and swap the two intersecting segments in the sweep line status.
func SweepLine[T types.SignedNumber](linesegments []LineSegment[T], opts ...Option) SweepLineResult {

	// Convert input slices to float64
	inputSegments := make([]LineSegment[float64], len(linesegments))
	for i := range linesegments {
		inputSegments[i] = linesegments[i].AsFloat64()
	}

	// Initialize a priority queue `pq` of potential future events, each associated with a point
	// in the plane and prioritized by the x-coordinate of the point.
	// So, initially, `pq` contains an event for each of the endpoints of the input segments.
	pq := newSweeplineEventQueue(inputSegments)

	// Prepare result struct.
	// We pre-allocate LineSegments with sufficient capacity to hold the number of line segments given.
	// InsertionPoints is unknown, so we let Go handle capacity management.
	result := SweepLineResult{
		LineSegments:       make([]LineSegment[float64], 0, len(linesegments)),
		IntersectionPoints: make([]Point[float64], 0),
	}

	// Perform Bentley-Ottmann sweep line below.
	// Code below has comments directly from wikipedia, explaining the algorithm.
	// After the original variable names, I've included this implementation's variable names in parentheses.

	// Initialise the "sweep line status tree".
	// Currently, this is just a sorted slice. The actual algorithm calls for a balanced binary search tree.
	// We may implement a BST in the future.
	status := sweepLineStatus{
		entries:  make([]sweepLineStatusEntry, 0, len(pq.events)),
		currentX: -math.MaxFloat64,
	}

	// Prepare for guarding against maximum iterations, see below.
	maxIterations := uint64(len(inputSegments) * len(inputSegments))
	iterationCount := uint64(0)

	// for each event in the queue...
	for pq.Len() > 0 {

		// Sanity check: guard against infinite loop.
		//
		// The Bentley-Ottmann algorithm processes a finite number of events: start, end,
		// and intersection events.
		//
		// The total number of events is at most O(n^2), where n is the number of input line segments.
		// This includes:
		//  - 2n 'start' and 'end' events (one for each endpoint of each line segment).
		//  - Up to n(n-1)/2 intersection events in the worst case, where every segment intersects
		//    with every other segment.
		//
		// By setting a limit of n^2 iterations, we ensure the algorithm terminates even in
		// worst-case scenarios involving dense intersections or overlapping inputs.
		//
		// This limit acts as a safeguard against infinite loops caused by unexpected
		// implementation bugs or edge cases not yet catered for.
		//
		// A warning is logged if the iteration count nears the maximum limit to assist with debugging.
		if iterationCount > maxIterations {
			panic(fmt.Errorf("exceeded maximum expected iterations: %d", maxIterations))
		}
		if iterationCount > uint64(0.9*float64(maxIterations)) {
			log.Printf("warning: sweep line iterations nearing maximum expected value (%d/%d)", iterationCount, maxIterations)
		}
		iterationCount++

		// Pop the event from the event queue
		event := heap.Pop(pq).(*sweeplineEvent)

		// Sanity check: ensure the scan line only moves rightwards.
		// This panic should never happen but is here to guard against errors in the code.
		if event.point.x < status.currentX {
			panic(fmt.Errorf("sweep line moves backwards from x=%f to x=%f",
				status.currentX, event.point.x))
		}

		// Set the sweep line's current X position to be that of the point
		status.currentX = event.point.x

		// Perform Bentley-Ottmann depending on event type
		switch event.eventType {

		case sweeplineEventStart:
			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// If p (event.point) is the left endpoint of a line segment s (event.lineSegments[0]),
			// insert s (event.lineSegments[0]) into T (status).
			status.Insert(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// Find the line-segments r (above) and t (below) that are respectively
			// immediately above and below s (event.lineSegments[0]) in T (status) (if they exist)
			above, below := status.Neighbors(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// if the crossing of r (above) and t (below) (the neighbours of s (event.lineSegments[0])
			// in the status data structure) forms a potential future event in the event queue,
			// remove this possible future event from the event queue
			if above != nil && below != nil {
				pq.RemoveIntersection(above.lineSegment, below.lineSegment, opts...)
			}

			if above != nil {
				// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
				// If s (event.lineSegments[0]) crosses r (above) or t (below),
				// add those crossing points as potential future events in the event queue.
				pq.AddIntersectionEvent(event.lineSegments[0], above.lineSegment, opts...)
			}

			if below != nil {
				// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
				// If s (event.lineSegments[0]) crosses r (above) or t (below),
				// add those crossing points as potential future events in the event queue.
				pq.AddIntersectionEvent(event.lineSegments[0], below.lineSegment, opts...)
			}

		case sweeplineEventEnd:

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// Find the segments r (above) and t (below) that (prior to the removal of s (event.lineSegments[0]))
			// were respectively immediately above and below it in T (status) (if they exist).
			above, below := status.Neighbors(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// If p (event.point) is the right endpoint of a line segment s (event.lineSegments[0]),
			// remove s (event.lineSegments[0]) from T (status).
			status.Remove(event.lineSegments[0])

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// If r (above) and t (below) cross, add that crossing point as a potential future event in the event queue.
			if above != nil && below != nil {
				pq.AddIntersectionEvent(above.lineSegment, below.lineSegment, opts...)
			}

		case sweeplineEventIntersection:

			// Add intersection point to SweepLineResult
			result.IntersectionPoints = append(result.IntersectionPoints, event.point)
			log.Printf("[result] added intersection: %s", event.point.String())

			// From https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm:
			// If p (event.point) is the crossing point of two segments s and t (with s below t to the left of the crossing),
			// swap the positions of s and t in T.
			//
			// Note:
			//   - There are two `if` blocks to handle both orders of s and t in event.lineSegments.
			//   - The variable names assigned match the descriptions above.
			//   - Repetition is kept for code clarity.
			if event.lineSegments[0].start.y < event.lineSegments[1].start.y {
				s := event.lineSegments[0]
				t := event.lineSegments[1]
				status.Swap(s, t)

				// After the swap, find the segments r and u (if they exist) that are immediately
				// below and above t and s, respectively.
				_, r := status.Neighbors(t)
				u, _ := status.Neighbors(s)

				// Remove any crossing points rs (i.e. a crossing point between r and s) and tu
				// (i.e. a crossing point between t and u) from the event queue and,
				// if r and t cross or s and u cross, add those crossing points to the event queue.
				if r != nil {
					pq.RemoveIntersection(r.lineSegment, s)
					pq.AddIntersectionEvent(r.lineSegment, t)
				}
				if u != nil {
					pq.RemoveIntersection(t, u.lineSegment)
					pq.AddIntersectionEvent(s, u.lineSegment)
				}
			}

			if event.lineSegments[1].start.y < event.lineSegments[0].start.y {
				s := event.lineSegments[1]
				t := event.lineSegments[0]
				status.Swap(s, t)

				// After the swap, find the segments r and u (if they exist) that are immediately
				// below and above t and s, respectively.
				_, r := status.Neighbors(t)
				u, _ := status.Neighbors(s)

				// Remove any crossing points rs (i.e. a crossing point between r and s) and tu
				// (i.e. a crossing point between t and u) from the event queue and,
				// if r and t cross or s and u cross, add those crossing points to the event queue.
				if r != nil {
					pq.RemoveIntersection(r.lineSegment, s)
					pq.AddIntersectionEvent(r.lineSegment, t)
				}
				if u != nil {
					pq.RemoveIntersection(t, u.lineSegment)
					pq.AddIntersectionEvent(s, u.lineSegment)
				}
			}
		}

		log.Printf("[status] %s", status.String())
	}

	return result
}
