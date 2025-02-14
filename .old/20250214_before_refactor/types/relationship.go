package types

import "fmt"

// Relationship defines the spatial or geometric relationship between two shapes or entities
// in 2D space. This type is used across various functions to represent how one geometric
// entity relates to another.
//
// The possible relationships include:
//   - Disjoint: The entities do not intersect, overlap, or touch in any way.
//   - Intersection: The entities share some boundary or overlap.
//   - Contained: One entity is fully within the boundary of another.
//   - ContainedBy: One entity is fully enclosed by another.
//   - Equal: The entities are identical in shape, size, and position.
//
// Each relationship type is represented as a constant value of the Relationship type.
// Functions that evaluate relationships between geometric entities typically return one
// of these constants to describe the spatial relationship between them.
//
// See the individual constant definitions for more details.
type Relationship uint8

// Valid values for Relationship:
const (
	// RelationshipDisjoint indicates that the two entities are completely separate,
	// with no overlap, touching, or intersection.
	RelationshipDisjoint Relationship = iota

	// RelationshipIntersection indicates that the two entities overlap or share
	// a boundary. This includes cases where the entities partially intersect
	// or where they touch at one or more points.
	RelationshipIntersection

	// RelationshipContainedBy indicates that the first entity is fully enclosed
	// within the second entity. The boundary of the first entity does not extend
	// outside the boundary of the second entity.
	RelationshipContainedBy

	// RelationshipContains indicates that the first entity fully encloses the
	// second entity. The boundary of the second entity does not extend outside
	// the boundary of the first entity.
	RelationshipContains

	// RelationshipEqual indicates that the two entities are identical in shape,
	// size, and position. This includes cases where their boundaries coincide exactly.
	RelationshipEqual
)

// FlipContainment reverses containment relationships for a [Relationship].
//
// This method is used to swap the roles of containment when interpreting
// relationships. Specifically:
//   - If the [Relationship] is RelationshipContainedBy, it is flipped to RelationshipContains.
//   - If the [Relationship] is RelationshipContains, it is flipped to RelationshipContainedBy.
//   - All other [Relationship] values are returned unchanged.
//
// Returns:
//   - [Relationship]: The flipped or unchanged [Relationship].
//
// Example:
//
//	rel := RelationshipContainedBy
//	flipped := rel.flipContainment()
//	fmt.Println(flipped) // Output: RelationshipContains
func (r Relationship) FlipContainment() Relationship {
	switch r {
	case RelationshipContainedBy:
		return RelationshipContains
	case RelationshipContains:
		return RelationshipContainedBy
	default:
		return r
	}
}

// String converts a [Relationship] value to its string representation.
//
// This method provides a human-readable string corresponding to the [Relationship]
// constant, such as RelationshipDisjoint or RelationshipContainedBy. It is useful
// for debugging and logging purposes.
//
// Supported [Relationship] values:
//   - [RelationshipDisjoint]: The objects are disjoint and do not touch or intersect.
//   - [RelationshipIntersection]: The objects intersect or overlap at some point.
//   - [RelationshipContainedBy]: The object is fully contained within another object.
//   - [RelationshipContains]: The object fully contains another object.
//   - [RelationshipEqual]: The objects are identical in size, shape, and position.
//
// Returns:
//   - string: The string representation of the [Relationship].
//
// Panics:
//   - If the [Relationship] value is not supported, this method panics with an error message.
func (r Relationship) String() string {
	switch r {
	case RelationshipDisjoint:
		return "RelationshipDisjoint"
	case RelationshipIntersection:
		return "RelationshipIntersection"
	case RelationshipContainedBy:
		return "RelationshipContainedBy"
	case RelationshipContains:
		return "RelationshipContains"
	case RelationshipEqual:
		return "RelationshipEqual"
	default:
		panic(fmt.Errorf("unsupported relationship type: %d", r))
	}
}
