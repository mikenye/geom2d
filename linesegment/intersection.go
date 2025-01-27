package linesegment

import (
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
)

// FindIntersectionsSlow performs a naive O(n^2) check to find all intersections
// between the given line segments, considering the provided geometry options.
//
// Parameters:
//   - segments: A slice of [LineSegment][T] instances to check for intersections.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] to customize the intersection behavior.
//
// Returns:
//   - []IntersectionResult[T]: A slice of intersection results, including points and segments
//     where intersections occur. Duplicate intersection points or segments are not included.
//
// Behavior:
//   - This function compares every pair of segments in the input slice and uses the
//     Intersection method to check for intersections, applying the provided options.
//     It avoids redundant checks; each pair is only checked once.
//
// Note:
//   - This is a naive implementation and should be used for small input sizes or as a baseline
//     for benchmarking more efficient algorithms.
func FindIntersectionsSlow[T types.SignedNumber](segments []LineSegment[T], opts ...options.GeometryOptionsFunc) []IntersectionResult[float64] {
	R := newIntersectionResults[float64]()

	// Compare each segment with every other segment
	for i := 0; i < len(segments); i++ {
		for j := i + 1; j < len(segments); j++ { // Start at i+1 to avoid duplicate checks
			// Check for intersection
			R.Add(segments[i].AsFloat64().Intersection(segments[j].AsFloat64(), opts...))
		}
	}

	return R.Results()
}

// Intersection calculates the intersection between two [LineSegment] instances.
//
// This method determines whether the current LineSegment and the given `other` LineSegment
// intersect within their boundaries. If an intersection exists, it returns details about
// the intersection, including whether it is a single point or a segment (in the case of
// collinear overlapping segments). If no intersection exists, it returns an appropriate
// IntersectionResult indicating no intersection.
//
// Parameters:
//   - other (LineSegment[T]): The second line segment to check for intersection.
//   - opts (options.GeometryOptionsFunc): Optional parameters, such as the epsilon value for
//     numerical precision adjustments.
//
// Returns:
//   - IntersectionResult[T]: A structure containing information about the type of intersection
//     (point, segment, or none) and any relevant intersection geometry.
//
// Behavior:
//   - If the segments are collinear and overlap, the function returns an IntersectionOverlappingSegment type
//     with the overlapping segment.
//   - If the segments intersect at a single point, the function returns an IntersectionPoint type
//     with the intersection coordinates.
//   - If the segments are parallel but not collinear, or if the intersection lies outside the segment
//     bounds, the function returns an IntersectionNone type.
func (l LineSegment[T]) Intersection(other LineSegment[T], opts ...options.GeometryOptionsFunc) IntersectionResult[T] {

	// Apply geometry options with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	// Define segment endpoints for AB (l) and CD (other)
	A, B := l.start.AsFloat64(), l.end.AsFloat64()
	C, D := other.start.AsFloat64(), other.end.AsFloat64()

	// Calculate the direction vectors
	dir1 := B.Translate(A.Negate())
	dir2 := D.Translate(C.Negate())

	// Calculate the determinants
	denominator := dir1.CrossProduct(dir2)

	// Handle collinear case (denominator == 0)
	if denominator == 0 {
		// Check if the segments are collinear
		AC := C.Translate(A.Negate())
		if AC.CrossProduct(dir1) != 0 {

			// Parallel but not collinear
			return IntersectionResult[T]{
				IntersectionType:  IntersectionNone,
				InputLineSegments: []LineSegment[T]{l, other},
			}
		}

		// Check overlap by projecting points onto the line
		tStart := (C.Translate(A.Negate())).DotProduct(dir1) / dir1.DotProduct(dir1)
		tEnd := (D.Translate(A.Negate())).DotProduct(dir1) / dir1.DotProduct(dir1)

		// Ensure tStart < tEnd for consistency
		if tStart > tEnd {
			tStart, tEnd = tEnd, tStart
		}

		// Check for overlap
		tOverlapStart := math.Max(0.0, tStart)
		tOverlapEnd := math.Min(1.0, tEnd)

		if tOverlapStart > tOverlapEnd {

			// No overlap
			return IntersectionResult[T]{
				IntersectionType:  IntersectionNone,
				InputLineSegments: []LineSegment[T]{l, other},
			}
		}

		// Calculate the overlapping segment
		overlapStart := point.New(
			numeric.SnapToEpsilon(A.X()+tOverlapStart*dir1.X(), geoOpts.Epsilon),
			numeric.SnapToEpsilon(A.Y()+tOverlapStart*dir1.Y(), geoOpts.Epsilon),
		)

		overlapEnd := point.New(
			numeric.SnapToEpsilon(A.X()+tOverlapEnd*dir1.X(), geoOpts.Epsilon),
			numeric.SnapToEpsilon(A.Y()+tOverlapEnd*dir1.Y(), geoOpts.Epsilon),
		)

		return IntersectionResult[T]{
			IntersectionType:   IntersectionOverlappingSegment,
			OverlappingSegment: NewFromPoints(overlapStart, overlapEnd),
			InputLineSegments:  []LineSegment[T]{l, other},
		}
	}

	// Calculate parameters t and u for non-collinear case
	AC := C.Translate(A.Negate())
	tNumerator := AC.CrossProduct(dir2)
	uNumerator := AC.CrossProduct(dir1)

	// It uses the parametric form of the line segments to solve for intersection parameters t and u.
	// If t and u are both in the range [0, 1], the intersection point lies within the bounds of
	// both segments.
	t := tNumerator / denominator
	u := uNumerator / denominator

	// Check if intersection occurs within the segment bounds
	if t < 0 || t > 1 || u < 0 || u > 1 {

		// Intersection is outside the segments
		return IntersectionResult[T]{
			IntersectionType:  IntersectionNone,
			InputLineSegments: []LineSegment[T]{l, other},
		}
	}

	// Calculate the intersection point
	intersection := point.New(
		numeric.SnapToEpsilon(A.X()+t*dir1.X(), geoOpts.Epsilon),
		numeric.SnapToEpsilon(A.Y()+t*dir1.Y(), geoOpts.Epsilon),
	)

	return IntersectionResult[T]{
		IntersectionType:  IntersectionPoint,
		IntersectionPoint: intersection,
		InputLineSegments: []LineSegment[T]{l, other},
	}
}
