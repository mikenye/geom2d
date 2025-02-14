package types

import "fmt"

// PointOrientation represents the relative orientation of three points in a two-dimensional plane.
// It describes whether the points are collinear, form a clockwise turn, or form a counterclockwise turn.
// This type is commonly used in computational geometry algorithms to determine the spatial relationship
// between points in relation to each other.
type PointOrientation uint8

// Valid values for PointOrientation.
const (
	// PointsCollinear indicates that the points are collinear, meaning they lie on a single straight line.
	PointsCollinear PointOrientation = iota

	// PointsClockwise indicates that the points are arranged in a clockwise orientation.
	PointsClockwise

	// PointsCounterClockwise indicates that the points are arranged in a counterclockwise orientation.
	PointsCounterClockwise
)

// String converts a [PointOrientation] constant into its string representation.
//
// This method is used to provide a human-readable description of the [PointOrientation] value.
// It is particularly useful for debugging and logging, as it outputs the specific orientation
// type (e.g., [PointsCollinear], [PointsClockwise], or [PointsCounterClockwise]).
//
// Behavior:
//   - If the value corresponds to a defined [PointOrientation] constant, the method returns its name.
//   - If the value is unsupported or invalid, the method panics with an error.
//
// Returns:
//   - string: The string representation of the [PointOrientation] value.
//
// Panics:
//   - If the [PointOrientation] value is invalid or not one of the defined constants
//     ([PointsCollinear], [PointsClockwise], [PointsCounterClockwise]), the function panics
//     with a descriptive error message.
func (o PointOrientation) String() string {
	switch o {
	case PointsCollinear:
		return "PointsCollinear"
	case PointsClockwise:
		return "PointsClockwise"
	case PointsCounterClockwise:
		return "PointsCounterClockwise"
	default:
		panic(fmt.Errorf("unsupported PointOrientation: %d", o))
	}
}
