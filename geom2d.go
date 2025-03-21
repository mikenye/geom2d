// Package geom2d provides a comprehensive set of tools for computational geometry in two dimensions.
//
// The geom2d package is built around core types like point.Point, linesegment.LineSegment, circle.Circle, rectangle.Rectangle,
// supporting a wide range of operations including transformations, boolean geometry, and spatial relationships.
// TODO: add polygon core type when implemented.
//
// Designed for both performance and clarity, geom2d leverages Go generics to handle various numeric types
// and provides intuitive APIs for working with 2D geometric data.
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
// # Support for Generics
//
// geom2d leverages Go’s generics, allowing you to use the library with different numeric types
// (int, float32, float64, etc.). This flexibility ensures the library can adapt to various applications,
// from integer-based grids to floating-point precision computations.
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
// [RelationshipDisjoint], [RelationshipIntersection], [RelationshipContainedBy], [RelationshipContains], and [RelationshipEqual].
//
// TODO: requires re-implementing after implementation of polygone core types
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

func GetEpsilon() float64 {
	return epsilon
}

func SetEpsilon(e float64) {
	epsilon = e
}
