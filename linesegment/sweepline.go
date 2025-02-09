package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"log"
)

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
	EventQueue := newEventQueue(segments, opts...)

	// Initialize status structure
	var StatusStructure *btree.BTreeG[sItem]
	StatusStructure = nil

	// Initialize results
	Results := newIntersectionResults[float64]()

	// while EventQueue is not empty
	iterCount := 0
	for EventQueue.Len() > 0 {
		iterCount++

		log.Printf("\n\n\n---ITERATION %d---\n\n\n", iterCount)

		// DEBUGGING: show queue
		debugPrintQueue(EventQueue)

		// Determine the next event point p in Q and delete it
		event, ok := EventQueue.DeleteMin()
		if !ok {
			panic(fmt.Errorf("unexpected empty queue"))
		}

		log.Printf("Popped event: %s\n", event)

		// Update the status structure based on new sweepline position
		StatusStructure = updateStatusStructure(StatusStructure, event, opts...)

		// DEBUGGING: show status structure
		log.Println("Status structure (S):")
		debugStatusStructure(StatusStructure)

		// Handle the event
		handleEventPointNew(event, EventQueue, StatusStructure, Results, opts...)
	}
	return Results.Results()
}

func handleEventPointNew(
	event qItem,
	EventQueue *btree.BTreeG[qItem],
	StatusStructure *btree.BTreeG[sItem],
	Results *intersectionResults[float64],
	opts ...options.GeometryOptionsFunc,
) {

	// Let U(p) be the set of segments whose upper endpoint is event.point
	UofP := event.segments

	// Find all segments stored in StatusStructure that contain event

	// DEBUGGING
	log.Println("Find all segments stored in StatusStructure that contain event:")
	// Let L(p) denote the subset of segments found whose lower endpoint is event
	// Let C(p) denote the subset of segments found that contain event in their interior
	CofP := make([]LineSegment[float64], 0)
	LofP := make([]LineSegment[float64], 0)

	// todo: can this be optimised?
	StatusStructure.Ascend(func(item sItem) bool {

		// skip if upper matches
		upper, _ := item.segment.sweeplineUpperPoint()
		if upper.Eq(event.point, opts...) {

			// DEBUGGING
			log.Printf("Ignoring %s", item.segment)
			return true
		}

		// check lower endpoint
		lower, _ := item.segment.sweeplineLowerPoint()
		if lower.Eq(event.point, opts...) {

			// DEBUGGING
			log.Printf("Adding %s to L(p)", item.segment)
			LofP = append(LofP, item.segment)
			return true
		}

		// check interior
		if item.segment.ContainsPoint(event.point, opts...) {

			// DEBUGGING
			log.Printf("Adding %s to C(p)", item.segment)
			CofP = append(CofP, item.segment)
			return true
		}

		// DEBUGGING
		log.Printf("Ignoring %s", item.segment)
		return true
	})

	log.Println("U(p):", UofP)
	log.Println("C(p):", CofP)
	log.Println("L(p):", LofP)

	// if L(p) ∪ U(p) ∪ C(p) contains more than one segment
	// then Report event as an intersection, together with L(p), U(p), and C(p).
	if len(UofP)+len(CofP)+len(LofP) > 1 {
		log.Printf("L(p) ∪ U(p) ∪ C(p) contains more than one segment, so event '%s' is an intersection", event)

		for _, result := range FindIntersectionsSlow(append(UofP, append(CofP, LofP...)...), opts...) {
			Results.Add(result)

			if result.IntersectionType == IntersectionOverlappingSegment {
				start, end := result.OverlappingSegment.Points()
				for _, p := range []point.Point[float64]{start, end} {
					pointValid := p.Y() < event.point.Y() ||
						(p.Y() == event.point.Y() && p.X() > event.point.X())
					if pointValid {
						newQItem := qItem{point: p}
						if !EventQueue.Has(newQItem) {
							log.Printf("Inserting overlapping segment endpoint to EventQueue: %s", newQItem.point)
							EventQueue.ReplaceOrInsert(newQItem)
						} else {
							log.Printf("Overlapping segment endpoint already exists in EventQueue: %s", newQItem.point)
						}

					}
				}
			}
		}
	}

	// Delete the segments in L(p) ∪ C(p) from StatusStructure.

	// DEBUGGING
	log.Println("Delete the segments in L(p) ∪ C(p) from StatusStructure:")
	for _, seg := range append(LofP, CofP...) {
		item, deleted := StatusStructure.Delete(sItem{
			segment: seg,
		})

		// DEBUGGING
		if deleted {
			log.Printf("Deleted: %s", item.String())
		} else {
			log.Printf("Attemted to delete but was not in StatusStructure: %s", item.String())
		}
	}

	// Insert the segments in U(p) ∪ C(p) into StatusStructure.
	// The order of the segments in StatusStructure should correspond to the order in which they are
	// intersected by a sweep line just below event.
	// If there is a horizontal segment, it comes last among all segments containing p.

	// DEBUGGING
	log.Println("Insert the segments in U(p) ∪ C(p) into StatusStructure:")

	var UCofP *btree.BTreeG[sItem]
	UCofP = updateStatusStructure(UCofP, event, opts...)
	for _, seg := range append(UofP, CofP...) {
		item, replaced := UCofP.ReplaceOrInsert(sItem{segment: seg})
		if replaced {
			log.Printf("Replaced: %s with %s", item.String(), seg.String())
		} else {
			log.Printf("Inserted: %s", seg)
		}

	}
	log.Println("U(p) ∪ C(p):")
	debugStatusStructure(UCofP)

	UCofP.Ascend(func(item sItem) bool {

		// get upper & lower points of segment being added
		upper, _ := item.segment.sweeplineUpperPoint()
		lower, _ := item.segment.sweeplineLowerPoint()

		// create new sItem
		newItem := sItem{
			segment:   NewFromPoints(upper, lower),
			originals: []LineSegment[float64]{item.segment},
		}

		// insert or replace item
		replacedItem, replaced := StatusStructure.ReplaceOrInsert(item)

		// DEBUGGING
		if replaced {
			log.Printf("Replaced %s with %s", replacedItem.String(), newItem.String())
		} else {
			log.Printf("Inserted %s", newItem.String())
		}

		return true
	})

	// DEBUGGING: show status structure
	log.Println("Status structure (S):")
	debugStatusStructure(StatusStructure)

	// if U(p) ∪ C(p) = 0
	// ...then Let sl and sr be the left and right neighbors of event in StatusStructure.
	if UCofP.Len() == 0 {

		// DEBUGGING
		log.Println("Let sl and sr be the left and right neighbors of event in StatusStructure:")

		// find neighbors
		sL, sR, sLFound, sRFound := findNighborsByPoint(StatusStructure, event.point, opts...)

		// DEBUGGING
		if sLFound {
			log.Printf("sL: %s", sL.String())
		} else {
			log.Println("sL: not found")
		}

		// DEBUGGING
		if sRFound {
			log.Printf("sR: %s", sR.String())
		} else {
			log.Println("sR: not found")
		}

		if sRFound && sLFound {
			log.Println("Find new event between sL & sR:")
			findNewEventNew(EventQueue, sL, sR, event.point, opts...)
		}

	} else {

		// Let sPrime be the leftmost segment of U(p) ∪ C(p) in StatusStructure.
		// todo: can optimise with a lemgth check. if StatusStructure len == 1 then return the only entry
		log.Println("Let sPrime be the leftmost segment of U(p) ∪ C(p) in StatusStructure:")
		var sPrime LineSegment[float64]
		sPrimeFound := false
		UCofP.Ascend(func(item sItem) bool {
			if StatusStructure.Has(item) {
				sPrimeFound = true
				sPrime = item.segment
				return false
			}
			return true
		})

		// DEBUGGING
		if sPrimeFound {
			log.Printf("sPrime: %s", sPrime.String())
		} else {
			log.Println("sPrime: not found")
		}

		// Let sL be the left neighbor of sPrime in StatusStructure.
		log.Println("Let sL be the left neighbor of sPrime in StatusStructure:")
		if sPrimeFound {
			sL, _, sLFound, _ := findNighborsByLineSegment(StatusStructure, sPrime, opts...)

			// DEBUGGING
			if sLFound {
				log.Printf("sL: %s", sL.String())
			} else {
				log.Println("sL: not found")
			}

			if sLFound {
				log.Println("Find new event between sL & sPrime:")
				findNewEventNew(EventQueue, sL, sPrime, event.point, opts...)
			}
		}

		// Let sDoublePrime be the rightmost segment of U(p) ∪ C(p) in StatusStructure.
		// todo: can optimise with a lemgth check. if StatusStructure len == 1 then return the only entry
		log.Println("Let sDoublePrime be the rightmost segment of U(p) ∪ C(p) in StatusStructure:")
		var sDoublePrime LineSegment[float64]
		sDoublePrimeFound := false
		UCofP.Descend(func(item sItem) bool {
			if StatusStructure.Has(item) {
				sDoublePrimeFound = true
				sDoublePrime = item.segment
				return false
			}
			return true
		})

		// DEBUGGING
		if sDoublePrimeFound {
			log.Printf("sDoublePrime: %s", sDoublePrime.String())
		} else {
			log.Println("sDoublePrime: not found")
		}

		// Let sR be the right neighbor of sDoublePrime in StatusStructure.
		log.Println("Let sR be the right neighbor of sDoublePrime in StatusStructure:")
		if sDoublePrimeFound {
			_, sR, _, sRFound := findNighborsByLineSegment(StatusStructure, sDoublePrime, opts...)

			// DEBUGGING
			if sRFound {
				log.Printf("sR: %s", sR.String())
			} else {
				log.Println("sR: not found")
			}

			if sRFound {
				log.Println("Find new event between sDoublePrime & sR:")
				findNewEventNew(EventQueue, sDoublePrime, sR, event.point, opts...)
			}
		}
	}
}

func findNewEventNew(
	EventQueue *btree.BTreeG[qItem],
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
		if intersection.IntersectionPoint.Y() < event.Y() ||
			(intersection.IntersectionPoint.Y() == event.Y() && intersection.IntersectionPoint.X() > event.X()) {
			newQItem := qItem{
				point:    intersection.IntersectionPoint,
				segments: nil,
			}
			if !EventQueue.Has(newQItem) {
				log.Printf("Inserting intersection to EventQueue: %s", newQItem.point)
				EventQueue.ReplaceOrInsert(newQItem)
			} else {
				log.Printf("Intersection already exists in EventQueue: %s", newQItem.point)
			}
		}
	case IntersectionOverlappingSegment:
		log.Fatalln("overlapping segment")
	}
}

func findNighborsByLineSegment(
	StatusStructure *btree.BTreeG[sItem],
	event LineSegment[float64],
	opts ...options.GeometryOptionsFunc,
) (
	sL, sR LineSegment[float64],
	sLFound, sRFound bool,
) {
	// find point or right neighbor
	StatusStructure.DescendLessOrEqual(sItem{segment: event}, func(item sItem) bool {
		if item.segment.Eq(event, opts...) {
			return true
		}
		sLFound = true
		sL = item.segment
		return false
	})
	StatusStructure.AscendGreaterOrEqual(sItem{segment: event}, func(item sItem) bool {
		if item.segment.Eq(event, opts...) {
			return true
		}
		sRFound = true
		sR = item.segment
		return false
	})
	return sL, sR, sLFound, sRFound
}

func findNighborsByPoint(
	StatusStructure *btree.BTreeG[sItem],
	event point.Point[float64],
	opts ...options.GeometryOptionsFunc,
) (
	sL, sR LineSegment[float64],
	sLFound, sRFound bool,
) {
	// find point or right neighbor
	pivotFound := false
	var pivot sItem
	StatusStructure.Ascend(func(item sItem) bool {
		if item.segment.ContainsPoint(event, opts...) {
			pivot = item
			pivotFound = true
			return false
		} else if item.segment.XAtY(event.Y()) > event.X() { // If we've passed event, break
			pivot = item
			return false
		}
		return true
	})

	StatusStructure.DescendLessOrEqual(pivot, func(item sItem) bool {
		if item.segment.Eq(pivot.segment, opts...) {
			return true
		}
		sL = item.segment
		sLFound = true
		return false
	})

	if pivotFound {
		StatusStructure.AscendGreaterOrEqual(pivot, func(item sItem) bool {
			if item.segment.Eq(pivot.segment, opts...) {
				return true
			}
			sR = item.segment
			sRFound = true
			return false
		})
	} else {
		sR = pivot.segment
		sRFound = true
	}

	return sL, sR, sLFound, sRFound
}
