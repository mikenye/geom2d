package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/numeric"
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

func (ir IntersectionResult[T]) Eq(other IntersectionResult[T], opts ...options.GeometryOptionsFunc) bool {
	// check IntersectionType equality
	if ir.IntersectionType != other.IntersectionType {
		return false
	}

	// check intersections equality
	switch ir.IntersectionType {
	case IntersectionNone:
		return true
	case IntersectionPoint:
		if !ir.IntersectionPoint.Eq(other.IntersectionPoint, opts...) {
			return false
		}
	case IntersectionOverlappingSegment:
		if !ir.OverlappingSegment.Eq(other.OverlappingSegment, opts...) {
			return false
		}
	}

	// check InputLineSegments equality
	if len(ir.InputLineSegments) != len(other.InputLineSegments) {
		return false
	}
	for _, segA := range ir.InputLineSegments {
		found := false
		for _, segB := range other.InputLineSegments {
			if segA.Eq(segB, opts...) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
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
	results *btree.BTreeG[IntersectionResult[T]]
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
func newIntersectionResults[T types.SignedNumber](opts ...options.GeometryOptionsFunc) *intersectionResults[T] {
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)
	return &intersectionResults[T]{
		results: btree.NewG[IntersectionResult[T]](2, intersectionResultLessHigherFunc[T](geoOpts.Epsilon)),
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

	// don't bother proceeding if no intersection
	if result.IntersectionType == IntersectionNone {
		return
	}

	existing, found := R.results.Get(result)

	if found {
		for _, seg := range existing.InputLineSegments {
			if !slices.Contains(result.InputLineSegments, seg) {
				result.InputLineSegments = append(result.InputLineSegments, seg)
				log.Println("updating intersection result:", result)
			}
		}
	} else {
		log.Println("inserting intersection result:", result)
	}

	R.results.ReplaceOrInsert(result)
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
	final := make([]IntersectionResult[T], 0, R.results.Len())
	R.results.Ascend(func(item IntersectionResult[T]) bool {
		final = append(final, item)
		return true
	})
	return final
}

func intersectionResultLessHigherFunc[T types.SignedNumber](epsilon float64) func(a, b IntersectionResult[T]) bool {
	return func(a, b IntersectionResult[T]) bool {
		var la, lb LineSegment[float64]

		// Convert IntersectionResult to LineSegments for comparison
		switch a.IntersectionType {
		case IntersectionNone:
			panic(fmt.Errorf("cannot compare against none"))
		case IntersectionPoint:
			la = NewFromPoints(a.IntersectionPoint, a.IntersectionPoint)
		case IntersectionOverlappingSegment:
			la = a.OverlappingSegment
		}
		switch b.IntersectionType {
		case IntersectionNone:
			panic(fmt.Errorf("cannot compare against none"))
		case IntersectionPoint:
			lb = NewFromPoints(b.IntersectionPoint, b.IntersectionPoint)
		case IntersectionOverlappingSegment:
			lb = b.OverlappingSegment
		}

		// **Sorting logic to ensure consistent ordering**
		// 1. Compare by lower point (Y first, then X)
		laLower, _ := la.sweeplineLowerPoint()
		lbLower, _ := lb.sweeplineLowerPoint()

		if numeric.FloatLessThan(laLower.Y(), lbLower.Y(), epsilon) {
			return true
		} else if numeric.FloatGreaterThan(laLower.Y(), lbLower.Y(), epsilon) {
			return false
		}

		if numeric.FloatLessThan(laLower.X(), lbLower.X(), epsilon) {
			return true
		} else if numeric.FloatGreaterThan(laLower.X(), lbLower.X(), epsilon) {
			return false
		}

		// 2. If lower points are the same, compare upper points
		laUpper, _ := la.sweeplineUpperPoint()
		lbUpper, _ := lb.sweeplineUpperPoint()

		if numeric.FloatLessThan(laUpper.Y(), lbUpper.Y(), epsilon) {
			return true
		} else if numeric.FloatGreaterThan(laUpper.Y(), lbUpper.Y(), epsilon) {
			return false
		}

		if numeric.FloatLessThan(laUpper.X(), lbUpper.X(), epsilon) {
			return true
		} else if numeric.FloatGreaterThan(laUpper.X(), lbUpper.X(), epsilon) {
			return false
		}

		// 3. If both endpoints are the same, ensure a consistent ordering for IntersectionType
		return a.IntersectionType < b.IntersectionType
	}
}

func InterSectionResultsEq[T types.SignedNumber](a, b []IntersectionResult[T], opts ...options.GeometryOptionsFunc) bool {
	// length check
	if len(a) != len(b) {
		return false
	}
	// elements check
	for _, resultA := range a {
		found := false
		for _, resultB := range b {
			if resultA.Eq(resultB, opts...) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
