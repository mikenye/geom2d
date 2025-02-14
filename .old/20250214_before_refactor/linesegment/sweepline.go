package linesegment

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
)

// sweepLineEventQueue interface allows us to easily "swap out" the event queue implementation
// with different data structures to determine optimal performance.
// So far, I've found red-black tree to be faster (by 9-24%) than btree by BenchmarkFindIntersectionsFast
type sweepLineEventQueue interface {
	IsEmpty() bool
	Pop() (p point.Point[float64], UofP []LineSegment[float64])
	InsertPoint(p point.Point[float64])
	String() string
}

type sweepLineStatusStructure interface {
	FindCofPAndLofP(eventPoint point.Point[float64]) (CofP []LineSegment[float64], LofP []LineSegment[float64])
	FindNeighborsOfPoint(eventPoint point.Point[float64]) (sL *LineSegment[float64], sR *LineSegment[float64])
	FindNeighborsOfUofPAndCofP(UofP []LineSegment[float64], CofP []LineSegment[float64]) (sPrime, sL, sDoublePrime, sR *LineSegment[float64])
	Insert(seg LineSegment[float64])
	Remove(seg LineSegment[float64])
	String() string
	Update(eventPoint point.Point[float64]) sweepLineStatusStructure
}

// FindIntersectionsFast computes all intersection points and overlapping segments among a set of line segments
// using the sweep line algorithm, outlined in Section 2.1 of [Computational Geometry: Algorithms and Applications].
//
// This function efficiently finds intersections in O((n + k) log n) time, where n is the number of input segments
// and k is the number of intersections found. It is significantly faster than the naive O(n²) method when dealing
// with large input sets.
//
// The function removes duplicate segments before processing and ignores degenerate segments (where Start == End).
// It uses a balanced search tree (B-tree) as the event queue to maintain efficient event processing.
//
// Parameters:
//   - segments: A slice of line segments to check for intersections.
//   - opts: Optional geometry configuration options (e.g., epsilon for floating-point precision).
//   - If [options.WithEpsilon] is provided, the function performs an approximate equality check,
//     considering the points equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - A slice of IntersectionResult[float64] containing all intersection points and overlapping segments found.
//
// Algorithm Overview:
//  1. Deduplicate Input: Removes exact duplicate segments.
//  2. Event Queue Initialization: Inserts segment endpoints into a B-tree event queue (Q).
//     Upper endpoints store the corresponding segments.
//     Degenerate segments are skipped.
//  3. Status Structure Initialization: Uses a balanced search tree (S) to track active segments.
//  4. Sweep Line Execution: Iterates through event points in Q, updating S and computing intersections.
//  5. Result Collection: Returns a slice of intersection results.
//
// Performance Considerations:
//   - If the number of segments is small, the naive O(n²) method (FindIntersectionsSlow) may be faster
//     due to lower constant overhead.
//
// See Also:
//   - FindIntersectionsSlow for the naive O(n²) intersection detection method.
//
// [Computational Geometry: Algorithms and Applications]: https://www.springer.com/gp/book/9783540779735
// todo: update doc comments, example func
func FindIntersectionsFast[T types.SignedNumber](
	segments []LineSegment[T],
	opts ...options.GeometryOptionsFunc,
) []IntersectionResult[float64] {

	// Initialize an empty event queue
	EventQueue := newEventQueueRBT(segments, opts...)

	// Initialize status structure
	StatusStructure := newStatusStructureRBT(point.New(-math.MaxFloat64, math.MaxFloat64), opts...)

	// Initialize results
	Results := newIntersectionResults[float64]()

	// while EventQueue is not empty
	//iterCount := 0 // used for debugging
	for !EventQueue.IsEmpty() {
		//iterCount++

		// DEBUGGING
		//log.Printf("\n\n\n---ITERATION %d---\n\n\n", iterCount)

		// DEBUGGING: show event queue
		//log.Printf("contents of event queue:\n%s", EventQueue)

		// Determine the next event point p in Q and delete it
		eventPoint, UofP := EventQueue.Pop()

		// DEBUGGING: show popped event
		//log.Printf("Popped event: %s [U(p): %v]\n", eventPoint, UofP)

		// Update the status structure based on new sweepline position
		StatusStructure = StatusStructure.Update(eventPoint).(*statusStructureRBT)

		// DEBUGGING: show status structure
		//log.Printf("Status structure:\n%s", StatusStructure)

		// Handle the event
		handleEventPointNew(eventPoint, UofP, EventQueue, StatusStructure, Results, opts...)
	}
	return Results.Results()
}

func handleEventPointNew(
	eventPoint point.Point[float64],
	UofP []LineSegment[float64],
	EventQueue sweepLineEventQueue,
	StatusStructure sweepLineStatusStructure,
	Results *intersectionResults[float64],
	opts ...options.GeometryOptionsFunc,
) {

	// Let U(p) be the set of segments whose upper endpoint is event.point

	// Find all segments stored in StatusStructure that contain event

	// DEBUGGING
	//log.Printf("U(p): %v", UofP)

	// Find all segments stored in StatusStructure that contain event
	// Let L(p) denote the subset of segments found whose lower endpoint is event
	// Let C(p) denote the subset of segments found that contain event in their interior

	CofP, LofP := StatusStructure.FindCofPAndLofP(eventPoint)

	// DEBUGGING:
	//log.Printf("L(p): %v", LofP)
	//log.Printf("C(p): %v", CofP)

	// if L(p) ∪ U(p) ∪ C(p) contains more than one segment
	// then Report event as an intersection, together with L(p), U(p), and C(p).
	if len(UofP)+len(CofP)+len(LofP) > 1 {

		// DEBUGGING
		//log.Printf("L(p) ∪ U(p) ∪ C(p) contains more than one segment, so event '%s' is an intersection", eventPoint)

		for _, result := range FindIntersectionsSlow(append(UofP, append(CofP, LofP...)...), opts...) {
			Results.Add(result)

			// if the result is an overlapping segment, then the overlapping segment endpoint should
			// be added as an event, if it is below or to the right of the sweepline.
			if result.IntersectionType == IntersectionOverlappingSegment {
				start, end := result.OverlappingSegment.Points()
				for _, p := range []point.Point[float64]{start, end} {
					pointBelowOrRightOfSweepline := p.Y() < eventPoint.Y() ||
						(p.Y() == eventPoint.Y() && p.X() > eventPoint.X())
					if pointBelowOrRightOfSweepline {
						EventQueue.InsertPoint(p)
					}
				}
			}
		}
	}

	// todo: we should be able to skip the delete/insert as we're recreating S when event popped

	// Delete the segments in L(p) ∪ C(p) from StatusStructure.

	// DEBUGGING
	//log.Println("Delete the segments in L(p) ∪ C(p) from StatusStructure:")

	for _, seg := range append(LofP, CofP...) {
		StatusStructure.Remove(seg)
	}

	// DEBUGGING: show status structure
	//log.Printf("Status structure after delete:\n%s", StatusStructure)

	// Insert the segments in U(p) ∪ C(p) into StatusStructure.
	// The order of the segments in StatusStructure should correspond to the order in which they are
	// intersected by a sweep line just below event.
	// If there is a horizontal segment, it comes last among all segments containing p.

	// DEBUGGING
	//log.Println("Insert the segments in U(p) ∪ C(p) into StatusStructure:")

	for _, seg := range append(UofP, CofP...) {
		StatusStructure.Insert(seg)
	}

	// DEBUGGING: show status structure
	//log.Printf("Status structure after insert:\n%s", StatusStructure)

	// if U(p) ∪ C(p) = 0
	// ...then Let sl and sr be the left and right neighbors of event in StatusStructure.
	if len(UofP)+len(CofP) == 0 {

		// DEBUGGING
		//log.Println("U(p) ∪ C(p) = 0")
		//log.Println("Let sL and sR be the left and right neighbors of event in StatusStructure:")

		// find neighbors
		sL, sR := StatusStructure.FindNeighborsOfPoint(eventPoint)

		// DEBUGGING
		//log.Println("sL:", sL)
		//log.Println("sR:", sR)

		if sL != nil && sR != nil {

			// DEBUGGING:
			//log.Println("Find new events between sL & sR")

			findNewEventNew(EventQueue, *sL, *sR, eventPoint, opts...)

			// DEBUGGING: show event queue
			//log.Printf("contents of event queue after find new events:\n%s", EventQueue)
		}

	} else {

		//log.Println("U(p) ∪ C(p) != 0")

		// Let sPrime be the leftmost segment of U(p) ∪ C(p) in StatusStructure.
		// Let sL be the left neighbor of sPrime in StatusStructure.
		// Let sDoublePrime be the rightmost segment of U(p) ∪ C(p) in StatusStructure.
		// Let sR be the right neighbor of sDoublePrime in StatusStructure.
		sPrime, sL, sDoublePrime, sR := StatusStructure.FindNeighborsOfUofPAndCofP(UofP, CofP)

		// DEBUGGING:
		//log.Println("Let sPrime be the leftmost segment of U(p) ∪ C(p) in StatusStructure")
		//log.Println("Let sL be the left neighbor of sPrime in StatusStructure")
		//log.Println("sPrime:", sPrime)
		//log.Println("sL:", sL)

		if sPrime != nil && sL != nil {

			// DEBUGGING
			//log.Println("Find new events between sL & sPrime")

			findNewEventNew(EventQueue, *sL, *sPrime, eventPoint, opts...)

			// DEBUGGING: show event queue
			//log.Printf("contents of event queue after find new events:\n%s", EventQueue)

		}

		// DEBUGGING:
		//log.Println("Let sDoublePrime be the rightmost segment of U(p) ∪ C(p) in StatusStructure")
		//log.Println("Let sR be the right neighbor of sDoublePrime in StatusStructure")
		//log.Println("sDoublePrime:", sDoublePrime)
		//log.Println("sR:", sR)

		if sDoublePrime != nil && sR != nil {

			// DEBUGGING
			//log.Println("Find new events between sDoublePrime & sR")

			findNewEventNew(EventQueue, *sDoublePrime, *sR, eventPoint, opts...)

			// DEBUGGING: show event queue
			//log.Printf("contents of event queue after find new events:\n%s", EventQueue)
		}
	}
}

func findNewEventNew(
	EventQueue sweepLineEventQueue,
	a, b LineSegment[float64],
	event point.Point[float64],
	opts ...options.GeometryOptionsFunc,
) {
	// if sL and sR intersect below the sweep line, or on it and to the right of the
	// current event point event, and the intersection is not yet present as an
	// event in EventQueue, then Insert the intersection point as an event into EventQueue.
	intersection := a.Intersection(b, opts...)
	switch intersection.IntersectionType {
	case IntersectionNone:
		return
	case IntersectionPoint:

		// check to ensure the intersection point is below or to the right of the sweepline
		// if so, add the event
		if intersection.IntersectionPoint.Y() < event.Y() ||
			(intersection.IntersectionPoint.Y() == event.Y() && intersection.IntersectionPoint.X() > event.X()) {
			EventQueue.InsertPoint(intersection.IntersectionPoint)
		}
	case IntersectionOverlappingSegment:

		// check overlapping segment start & end points are below or to the right of the sweepline
		// if so, add them
		if intersection.OverlappingSegment.Start().Y() < event.Y() ||
			(intersection.OverlappingSegment.Start().Y() == event.Y() && intersection.OverlappingSegment.Start().X() > event.X()) {
			EventQueue.InsertPoint(intersection.OverlappingSegment.Start())
		}
		if intersection.OverlappingSegment.End().Y() < event.Y() ||
			(intersection.OverlappingSegment.End().Y() == event.Y() && intersection.OverlappingSegment.End().X() > event.X()) {
			EventQueue.InsertPoint(intersection.OverlappingSegment.End())
		}
	}
}
