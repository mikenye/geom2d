// Package geom2d provides a comprehensive set of tools for computational geometry in two dimensions.
//
// The geom2d package is built around core types like point.Point, linesegment.LineSegment, circle.Circle, rectangle.Rectangle,
// supporting a wide range of operations including transformations, boolean geometry, and spatial relationships.
// TODO: add polygon core type when implemented.
//
// Designed for both performance and clarity, geom2d provides intuitive APIs for working with 2D geometric data
// using float64 coordinates with configurable epsilon for precision control.
//
// # Coordinate System
//
// This library assumes a standard Cartesian coordinate system where the x-axis increases to the right and the y-axis
// increases upward. This system is commonly referred to as a right-handed Cartesian coordinate system.
// All geometric operations and relationships (e.g., clockwise or counterclockwise points) are based on this convention.
//
// # Core Geometric Types
//
// The library includes support for the following 2D geometric types:
//
//   - point.Point: Represents a single coordinate in 2D space.
//   - linesegment.LineSegment: Represents a straight line segment defined by two endpoints.
//   - rectangle.Rectangle: Represents an axis-aligned rectangle, defined by its corners.
//   - circle.Circle: Represents a circle defined by a center point and radius.
//
// TODO: add polygon core type when implemented.
//
// # Precision and Floating-Point Handling
//
// geom2d uses float64 coordinates throughout, providing a good balance of precision and performance.
// The library includes robust handling of floating-point precision issues through configurable epsilon
// values and dedicated numeric utility functions.
//
// # Precision Control with Epsilon
//
// geom2d incorporates an epsilon parameter in many of its relationship methods to handle floating-point
// precision issues. This allows you to control the tolerance for comparisons, making the library robust
// for real-world applications where precision errors can occur.
//
// # Relationships Between Geometric Types
//
// This library provides methods to compute relationships between geometric types using a standardized set of relationships:
// [types.RelationshipDisjoint], [types.RelationshipIntersection], [types.RelationshipContainedBy], [types.RelationshipContains], and [types.RelationshipEqual].
//
// TODO: requires re-implementing after implementation of polygon core types
//
// # Acknowledgments
//
// geom2d builds upon the work of others and is grateful for the foundations they have laid. Specifically:
//
//   - Mark de Berg, Otfried Cheong, Marc van Kreveld, and Mark Overmars: Their book [Computational Geometry: Algorithms and Applications]
//     has been an invaluable resource, providing rigorous explanations and guiding the implementation of fundamental geometric algorithms.
//   - Martínez et al.: Their paper on Boolean operations on polygons has been instrumental in the implementation of
//     the Martínez algorithm in this library. See [A simple algorithm for Boolean operations on polygons].
//   - Tom Wright: The inspiration for starting this library came from Tom Wright’s repository
//     [Provably Correct Polygon Algorithms] and his accompanying paper. While geom2d follows its own approach,
//     certain ideas have been influenced by his work.
//   - Jack Bresenham: The Bresenham's Line Algorithm and Bresenham's Circle Algorithm implemented in this library are
//     inspired by Jack Bresenham's pioneering work. These algorithms are efficient methods for rasterizing lines and
//     circles in computer graphics. For more details, see Bresenham's original paper
//     ["Algorithm for computer control of a digital plotter." IBM Systems Journal, 1965.]
//   - This project is a collaborative effort, with assistance from [OpenAI's Assistant] for brainstorming, debugging,
//     and refining implementations.
//
// [Computational Geometry: Algorithms and Applications]: https://www.springer.com/gp/book/9783540779735
// [A simple algorithm for Boolean operations on polygons]: https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124
// [Provably Correct Polygon Algorithms]: https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms
// ["Algorithm for computer control of a digital plotter." IBM Systems Journal, 1965.]: https://dl.acm.org/doi/10.1147/sj.41.025
// [OpenAI's Assistant]: https://openai.com/
package geom2d

var epsilon float64 = 1e-12

// GetEpsilon returns the current epsilon value used for floating-point comparisons across the library.
//
// Epsilon is a small positive value used for comparing floating-point numbers while accounting
// for precision errors. This function retrieves the global epsilon value that applies to all
// geometric operations in the library that require approximate floating-point comparisons.
//
// Returns:
//   - float64: The current epsilon value, initially set to 1e-12 by default.
func GetEpsilon() float64 {
	return epsilon
}

// SetEpsilon changes the global epsilon value used for floating-point comparisons.
//
// This function allows users to customize the tolerance level for floating-point
// comparisons throughout the library. Setting a larger epsilon makes comparisons more
// lenient, while a smaller value makes them more strict.
//
// Parameters:
//   - e (float64): The new epsilon value to use. Should be a small positive number.
//
// Usage:
//   - Increase epsilon when working with values that may have accumulated more
//     floating-point errors, such as after many transformations.
//   - Decrease epsilon when higher precision is required for specific operations.
func SetEpsilon(e float64) {
	epsilon = e
}
