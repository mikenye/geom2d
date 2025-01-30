package linesegment

import (
	"fmt"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"slices"
	"strings"
)

type IntersectionType uint8

const (
	IntersectionNone IntersectionType = iota
	IntersectionPoint
	IntersectionOverlappingSegment
)

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

type IntersectionResult[T types.SignedNumber] struct {
	IntersectionType   IntersectionType     // Type of intersection
	IntersectionPoint  point.Point[float64] // Valid if Type == IntersectionPoint
	OverlappingSegment LineSegment[float64] // Valid if Type == OverlappingSegment
	InputLineSegments  []LineSegment[T]     // Input line segments
}

func (ir IntersectionResult[T]) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Intersection type: %s", ir.IntersectionType.String()))
	switch ir.IntersectionType {
	case IntersectionNone:
	case IntersectionPoint:
		builder.WriteString(fmt.Sprintf(": %s from segments:\n", ir.IntersectionPoint.String()))
	case IntersectionOverlappingSegment:
		builder.WriteString(fmt.Sprintf(": %s from segments:\n", ir.OverlappingSegment.String()))
	}
	for _, seg := range ir.InputLineSegments {
		builder.WriteString(fmt.Sprintf("  - %s\n", seg.String()))
	}
	return builder.String()
}

type intersectionResults[T types.SignedNumber] struct {
	results []IntersectionResult[T]
}

func newIntersectionResults[T types.SignedNumber]() *intersectionResults[T] {
	return &intersectionResults[T]{
		results: make([]IntersectionResult[T], 0),
	}
}

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
		R.results = append(R.results, result)
		return
	}

	// Else, merge input line segments
	for _, seg := range result.InputLineSegments {
		if !slices.Contains(R.results[index].InputLineSegments, seg) {
			R.results[index].InputLineSegments = append(R.results[index].InputLineSegments, seg)
		}
	}
}

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
	}
}

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

func (R *intersectionResults[T]) Results() []IntersectionResult[T] {
	R.sortInputSegments()
	return R.results
}
