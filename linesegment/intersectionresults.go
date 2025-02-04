package linesegment

import (
	"fmt"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"log"
	"slices"
	"strings"
)

// IntersectionType represents the type of intersection between two LineSegment.
// It is used to classify intersection results into:
//   - IntersectionNone: There is no intersection
//   - IntersectionPoint: There is an intersection at a given point
//   - IntersectionOverlappingSegment: The line segments are collinear and overlap
type IntersectionType uint8

// Valid values for IntersectionType
const (
	// IntersectionNone indicates that there is no intersection between the given line segments.
	IntersectionNone IntersectionType = iota

	// IntersectionPoint indicates that the intersection occurs at a single point.
	// This happens when two line segments, cross each other at a single coordinate.
	IntersectionPoint

	// IntersectionOverlappingSegment indicates that the intersection results in a continuous overlapping segment.
	// This occurs when two collinear segments partially or fully overlap.
	IntersectionOverlappingSegment
)

// String returns a human-readable representation of the IntersectionType.
//
// This method implements the fmt.Stringer interface, allowing an IntersectionType
// to be printed as a string when used with fmt.Printf, fmt.Println, and similar functions.
//
// If an unsupported IntersectionType value is encountered, this function panics.
func (t IntersectionType) String() string {
	switch t {
	case IntersectionNone:
		return "IntersectionNone"
	case IntersectionPoint:
		return "IntersectionPoint"
	case IntersectionOverlappingSegment:
		return "IntersectionOverlappingSegment"
	default:
		panic(fmt.Errorf("unsupported line segment intersection type"))
	}
}

// IntersectionResult represents the outcome of an intersection between two or more line segments.
//
// The result of an intersection can be one of the following types, as indicated by the
// IntersectionType field:
//   - IntersectionNone: No intersection occurred.
//   - IntersectionPoint: The line segments intersect at a single point, stored in IntersectionPoint.
//   - IntersectionOverlappingSegment: The line segments overlap collinearly along a segment, stored in OverlappingSegment.
//
// Fields:
//   - IntersectionType: Specifies the type of intersection (IntersectionNone, IntersectionPoint, or IntersectionOverlappingSegment).
//   - IntersectionPoint: Stores the point of intersection if IntersectionType == IntersectionPoint.
//   - OverlappingSegment: Stores the overlapping segment if IntersectionType == IntersectionOverlappingSegment.
//   - InputLineSegments: Stores the original line segments that were tested for intersection.
//
// The generic type T allows the input line segments (InputLineSegments) to be stored with
// their original numeric type (e.g., `int`, `float64`), while intersection results (IntersectionPoint
// and OverlappingSegment) are always stored as float64 for precision.
type IntersectionResult[T types.SignedNumber] struct {

	// IntersectionType specifies the type of intersection
	//  - IntersectionNone
	//  - IntersectionPoint
	//  - IntersectionOverlappingSegment
	IntersectionType IntersectionType

	// IntersectionPoint stores the point of intersection if IntersectionType == IntersectionPoint
	IntersectionPoint point.Point[float64]

	// OverlappingSegment stores the overlapping segment if IntersectionType == IntersectionOverlappingSegment
	OverlappingSegment LineSegment[float64]

	// InputLineSegments stores the original line segments that were tested for intersection.
	InputLineSegments []LineSegment[T]
}

// String returns a human-readable string representation of the intersection result.
//
// The output format varies depending on the IntersectionType:
//   - IntersectionNone: "Intersection type: IntersectionNone, from segments: <segments>"
//   - IntersectionPoint: "Intersection type: IntersectionPoint: <point> from segments: <segments>"
//   - IntersectionOverlappingSegment: "Intersection type: IntersectionOverlappingSegment: <overlapping segment> from segments: <segments>"
//
// The segments involved in the intersection are listed after the intersection details.
//
// Example Outputs:
//
//	Intersection type: IntersectionNone, from segments: (1,1)(5,5), (1,5)(5,1)
//
//	Intersection type: IntersectionPoint: (3,3) from segments: (1,1)(5,5), (1,5)(5,1)
//
//	Intersection type: IntersectionOverlappingSegment: (2,2)(4,4) from segments: (1,1)(5,5), (2,2)(4,4)
//
// Returns:
//   - A formatted string describing the intersection type and the segments involved.
func (ir IntersectionResult[T]) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Intersection type: %s", ir.IntersectionType.String()))
	switch ir.IntersectionType {
	case IntersectionNone:
		builder.WriteString(", from segments: ")
	case IntersectionPoint:
		builder.WriteString(fmt.Sprintf(": %s from segments: ", ir.IntersectionPoint.String()))
	case IntersectionOverlappingSegment:
		builder.WriteString(fmt.Sprintf(": %s from segments: ", ir.OverlappingSegment.String()))
	}
	first := true
	for _, seg := range ir.InputLineSegments {
		if first {
			builder.WriteString(seg.String())
			first = false
			continue
		}
		builder.WriteString(fmt.Sprintf(", %s", seg.String()))
	}
	return builder.String()
}

// intersectionResults is a private utility type that manages intersection results.
//
// This type is used within the FindIntersectionsSlow & FindIntersectionsFast methods to accumulate, deduplicate,
// and merge intersection results before returning them to the caller. It provides methods
// for adding new results while ensuring correctness by avoiding duplicate intersection
// points or overlapping segments.
//
// Fields:
//   - results: A slice of IntersectionResult[T] storing detected intersections.
//
// This type is not exposed publicly and is intended for internal use within the
// sweep line and naive intersection detection algorithms.
type intersectionResults[T types.SignedNumber] struct {
	results []IntersectionResult[T]
}

// newIntersectionResults creates and returns a new instance of intersectionResults.
//
// This function initializes an empty slice of IntersectionResult[T] and returns a pointer
// to a newly allocated intersectionResults[T] instance. It is used internally
// within the FindIntersectionsSlow & FindIntersectionsFast methods to manage detected intersections efficiently.
//
// Returns:
//   - A pointer to an empty intersectionResults[T] instance, ready to store and process
//     intersection results.
func newIntersectionResults[T types.SignedNumber]() *intersectionResults[T] {
	return &intersectionResults[T]{
		results: make([]IntersectionResult[T], 0),
	}
}

// Add inserts an intersection result into the intersectionResults set while ensuring
// normalization and deduplication.
//
// This method processes the provided IntersectionResult[T] by normalizing it (ensuring
// consistent segment orientation) and then handling it based on its type. Depending on the
// intersection type, it delegates to either addPoint (for point intersections) or
// addOverlappingSegment (for segment overlaps). If the intersection type is IntersectionNone,
// it is ignored.
//
// Parameters:
//   - result: The intersection result to be added. This result is normalized before storage.
//   - opts: Optional geometry settings, such as epsilon values for floating-point precision handling.
//
// Behavior:
//   - If result.IntersectionType is IntersectionNone, the method does nothing.
//   - If result.IntersectionType is IntersectionPoint, the intersection point is processed and stored.
//   - If result.IntersectionType is IntersectionOverlappingSegment, the overlapping segment is processed.
//
// Internally, the method ensures that all input segments in the result are normalized before
// being added to the collection, maintaining consistency in how intersections are stored.
func (R *intersectionResults[T]) Add(result IntersectionResult[T], opts ...options.GeometryOptionsFunc) {

	// "normalize" result
	resultNormalized := IntersectionResult[T]{
		IntersectionType:   result.IntersectionType,
		IntersectionPoint:  result.IntersectionPoint,
		OverlappingSegment: result.OverlappingSegment.normalize(),
		InputLineSegments:  make([]LineSegment[T], 0, len(result.InputLineSegments)),
	}
	for _, seg := range result.InputLineSegments {
		resultNormalized.InputLineSegments = append(resultNormalized.InputLineSegments, seg.normalize())
	}

	// handle result
	switch result.IntersectionType {
	case IntersectionNone:
		// do nothing
	case IntersectionPoint:
		R.addPoint(resultNormalized, opts...)
	case IntersectionOverlappingSegment:
		R.addOverlappingSegment(resultNormalized, opts...)
	}
}

// addOverlappingSegment adds an intersection result of type IntersectionOverlappingSegment
// to the intersectionResults set while ensuring deduplication.
//
// This method is responsible for handling cases where two line segments overlap completely or
// partially. If the overlapping segment is not already present in the results, it is added.
// If it already exists, the function merges the input line segments into the existing result
// to track all segments involved in the overlap.
//
// Parameters:
//   - result: The intersection result to be added. Must be of type IntersectionOverlappingSegment.
//   - opts: Optional geometry settings, such as epsilon values for floating-point precision handling.
//
// Behavior:
//   - If IntersectionType is not IntersectionOverlappingSegment, the function panics.
//   - The function searches for an existing overlapping segment in R.results.
//   - If no existing overlap is found, result is added as a new entry.
//   - If an overlap is found, result.InputLineSegments are merged with the existing entry to
//     ensure all contributing segments are recorded.
func (R *intersectionResults[T]) addOverlappingSegment(result IntersectionResult[T], opts ...options.GeometryOptionsFunc) {
	if result.IntersectionType != IntersectionOverlappingSegment {
		panic(fmt.Errorf("IntersectionResult is not type IntersectionOverlappingSegment"))
	}

	// Check for existing overlaps
	index := slices.IndexFunc(R.results, func(i IntersectionResult[T]) bool {
		// Skip if existing intersection is not an overlap
		if i.IntersectionType != IntersectionOverlappingSegment {
			return false
		}

		// Check if the overlapping segments are equal
		return i.OverlappingSegment.Eq(result.OverlappingSegment, opts...)
	})

	// If overlap doesn't exist, add it
	if index == -1 {
		log.Println("adding intersection result:", result)
		R.results = append(R.results, result)
		return
	}

	// Else, merge input line segments
	for _, seg := range result.InputLineSegments {
		if !slices.Contains(R.results[index].InputLineSegments, seg) {
			R.results[index].InputLineSegments = append(R.results[index].InputLineSegments, seg)
			log.Println("updated intersection result:", R.results[index])
		}
	}
}

// addPoint adds an intersection result of type IntersectionPoint to the intersectionResults set,
// ensuring that duplicate intersection points are not stored while tracking all contributing segments.
//
// This method is responsible for handling point-based intersections where two or more line segments
// cross at a single point. If the intersection point does not already exist in the results, it is added.
// If the point is already present, the function merges `InputLineSegments` to ensure all contributing
// segments are recorded.
//
// Parameters:
//   - result: The intersection result to be added. Must be of type `IntersectionPoint`.
//   - opts: Optional geometry settings, such as epsilon values for floating-point precision handling.
//
// Behavior:
//   - If `result.IntersectionType` is not `IntersectionPoint`, the function panics.
//   - The function searches for an existing intersection point in `R.results`.
//   - If no existing intersection point is found, `result` is added as a new entry.
//   - If the intersection point is found, `result.InputLineSegments` are merged with the existing entry to
//     ensure all contributing segments are recorded.
func (R *intersectionResults[T]) addPoint(result IntersectionResult[T], opts ...options.GeometryOptionsFunc) {
	if result.IntersectionType != IntersectionPoint {
		panic(fmt.Errorf("IntersectionResult is not type IntersectionPoint"))
	}

	// check for existing points
	index := slices.IndexFunc(R.results, func(i IntersectionResult[T]) bool {
		// skip if existing intersection is not a point
		if i.IntersectionType != IntersectionPoint {
			return false
		}

		// skip if existing intersection point does not match new point
		if !i.IntersectionPoint.Eq(result.IntersectionPoint, opts...) {
			return false
		}

		return true
	})

	// if intersection point doesn't exist, add
	if index == -1 {
		log.Println("adding intersection result:", result)
		R.results = append(R.results, result)
		return
	}

	// else, add input line segments to existing intersection
	for _, seg := range result.InputLineSegments {

		// skip if line segment exists
		if slices.Contains(R.results[index].InputLineSegments, seg) {
			continue
		}

		// else, merge
		R.results[index].InputLineSegments = append(R.results[index].InputLineSegments, seg)
		log.Println("updated intersection result:", R.results[index])
	}
}

// sortInputSegments ensures that the InputLineSegments within each IntersectionResult
// are sorted in a consistent order. This is primarily used for test and example output consistency,
// allowing for reliable comparison of results.
//
// Sorting Criteria:
//   - Segments are sorted first by their start point Y-coordinate (ascending).
//   - If Y-coordinates are equal, they are sorted by start point X-coordinate (ascending).
//   - If start points are identical, sorting continues based on end point Y-coordinate (ascending).
//   - If end Y-coordinates are also equal, sorting falls back to end point X-coordinate (ascending).
//   - If all values are equal, the segments are considered identical.
//
// Purpose:
//   - Helps maintain deterministic output in tests and examples by ensuring consistent ordering of
//     input segments in intersection results.
//   - Uses slices.SortFunc from the Go standard library for in-place sorting.
//
// Notes:
//   - This function does not modify the intersection type or computed intersection points—
//     it only affects the order of `InputLineSegments` within each result.
func (R *intersectionResults[T]) sortInputSegments() {
	for i := range R.results {
		slices.SortFunc(R.results[i].InputLineSegments, func(a, b LineSegment[T]) int {
			// Compare by start point Y
			if a.Start().Y() < b.Start().Y() {
				return -1
			}
			if a.Start().Y() > b.Start().Y() {
				return 1
			}
			// Compare by start point X
			if a.Start().X() < b.Start().X() {
				return -1
			}
			if a.Start().X() > b.Start().X() {
				return 1
			}

			// Compare by end point Y
			if a.End().Y() < b.End().Y() {
				return -1
			}
			if a.End().Y() > b.End().Y() {
				return 1
			}
			// Compare by end point X
			if a.End().X() < b.End().X() {
				return -1
			}
			if a.End().X() > b.End().X() {
				return 1
			}

			// Otherwise, they are equal
			return 0
		})
	}
}

// Results returns the list of IntersectionResult entries stored in the intersectionResults struct.
// Before returning, it ensures that the InputLineSegments within each result are consistently sorted
// using sortInputSegments. This guarantees a deterministic order, making results easier to compare
// in tests and examples.
//
// Purpose:
//   - Provides access to the computed intersection results.
//   - Ensures consistent ordering of InputLineSegments before returning, aiding test reproducibility.
//
// Returns:
//   - A slice of IntersectionResult[T], where T is a signed number type (e.g., int, float).
//
// Notes:
//   - This function does not modify the intersection points themselves, only the segment order within results.
//   - The caller receives a reference to the underlying slice, meaning modifications to the returned slice
//     will affect the internal state of intersectionResults.
func (R *intersectionResults[T]) Results() []IntersectionResult[T] {
	R.sortInputSegments()
	return R.results
}
