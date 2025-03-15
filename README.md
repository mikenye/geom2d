# Geom2D

[![Go Reference](https://pkg.go.dev/badge/github.com/mikenye/geom2d.svg)](https://pkg.go.dev/github.com/mikenye/geom2d)
[![codecov](https://codecov.io/gh/mikenye/geom2d/graph/badge.svg?token=Z73XKN7N8D)](https://codecov.io/gh/mikenye/geom2d)

Geom2D is a computational geometry library for Go, designed for 2D polygon operations and other fundamental geometric types, and is currently reaching its release candidate phase, nearing production readiness.

**As of March 2025, thid package is actively in development.**

## Table of Contents
- [Geom2D](#geom2d)
- [Project Goals](#project-goals)
- [Getting Started](#getting-started)
- [Documentation](#documentation)
- [Geometric Relationships](#geometric-relationships)
- [Acknowledgments](#acknowledgments)
- [License](#license)

## Project Goals

Geom2D aims to provide a robust, flexible, and efficient implementation of 2D geometric operations, featuring:

- **Geometry Types**:
  - **Point**: Basic 2D point representation.
  - **LineSegment**: Represents a line segment and supports operations such as intersection and reflection.
  - **Circle**: Support for operations like circumference, area, and intersection checks.
  - **Rectangle**: Axis-aligned bounding box with methods for containment, intersection, and transformation.
  - **Polygon**: Support polygons with holes and nested structures, with methods for orientation, correction, and Boolean operations.
- **Polygon Boolean Operations**: Union, intersection, and subtraction.
- **Geometry to Geometry Relationships**: Fast and reliable algorithms for determining geometric relationships.

## Getting Started

To install Geom2D, use go get:

```bash
go get github.com/mikenye/geom2d
```

### Examples

For detailed examples, please see the [package documentation](https://pkg.go.dev/github.com/mikenye/geom2d), where almost every public function has an example.

## Documentation

For detailed API documentation and usage examples, visit the [geom2d documentation at the Go Package Discovery and Documentation site](https://pkg.go.dev/github.com/mikenye/geom2d).

## Geometric Relationships

This section describes the geometric relationships between different types of geometric objects supported by the library. Relationships can include concepts like disjoint, intersection, containment, and equality. The relationships are determined using efficient algorithms for each pair of types.

This table describes the **relationship of the left-side type (column) to the top-side type (row)**.  
Each cell indicates the valid relationship types.

| **Left ↓, Right →**         | Point                                                                    | Line Segment                                                             | Circle                                                                                                                   | Rectangle                                                                                                                | Polygon within PolyTree                                                                                                  |
|-----------------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|
| **Point**                   | `RelationshipDisjoint`<br>`RelationshipEqual`                               | `RelationshipDisjoint`<br>`RelationshipIntersection`                         | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              |
| **Line Segment**            | `RelationshipDisjoint`<br>`RelationshipIntersection`                        | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipEqual`    | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`                                              |
| **Circle**                  | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`<br>`RelationshipEqual` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      |
| **Rectangle**               | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`<br>`RelationshipEqual` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      |
| **Polygon within PolyTree** | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContains` | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`                      | `RelationshipDisjoint`<br>`RelationshipIntersection`<br>`RelationshipContainedBy`<br>`RelationshipContains`<br>`RelationshipEqual` |

### Relationship Definitions

 - `RelationshipDisjoint`: The two geometric objects do not overlap.
 - `RelationshipIntersection`: The two geometric objects overlap partially or at boundaries.
 - `RelationshipContains`: The left-side geometric object fully contains the right-side object.
 - `RelationshipContainedBy`: The left-side geometric object is fully contained by the right-side object.
 - `RelationshipEqual`: The two geometric objects are identical.

## Acknowledgments

Geom2D builds upon the work of others and is grateful for the foundations they have laid. Specifically:

- **Computational Geometry: Algorithms and Applications**: The implementation of geometric algorithms in Geom2D follows the approach presented in [Computational Geometry: Algorithms and Applications (3rd edition)](https://link.springer.com/book/10.1007/978-3-540-77974-2) by Mark de Berg, Otfried Cheong, Marc van Kreveld, and Mark Overmars. This book serves as the primary reference for methods such as the sweep line algorithm for line segment intersections, polygon operations, and spatial relationships.

- **Martínez et al.**: the work of Martínez et al. on Boolean operations for polygons provided valuable insight into polygonal computational geometry. See [A simple algorithm for Boolean operations on polygons](https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124).

- **Bentley-Ottmann Algorithm**: The Bentley-Ottmann algorithm for efficiently reporting line segment intersections served as the foundation for the sweep line approach used in this library. While Geom2D follows the refined method outlined in Computational Geometry: Algorithms and Applications, the Bentley-Ottmann technique remains an influential cornerstone in computational geometry. See Bentley, J. L., & Ottmann, T. A. ["Algorithms for reporting and counting geometric intersections."](https://doi.org/10.1145/361002.361007) Communications of the ACM, 1979, or [the "Bentley–Ottmann algorithm" article on Wikipedia](https://en.wikipedia.org/wiki/Bentley–Ottmann_algorithm).

- **Tom Wright**: The inspiration for starting this library came from Tom Wright’s repository [Provably Correct Polygon Algorithms](https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms) and his accompanying paper. While Geom2D follows its own approach, certain ideas have been influenced by his work.

- **Jack Bresenham**: The Bresenham's Line Algorithm and Bresenham's Circle Algorithm implemented in this library are inspired by Jack Bresenham's pioneering work. These algorithms are efficient methods for rasterizing lines and circles in computer graphics. For more details, see Bresenham's original paper ["Algorithm for computer control of a digital plotter." IBM Systems Journal, 1965.](https://dl.acm.org/doi/10.1147/sj.41.025)

- This project is a collaborative effort, with significant assistance from [OpenAI's Assistant](https://openai.com/) for brainstorming, debugging, and refining implementations.

To learn more about the work that inspired this library, visit the linked papers and repositories.

## Future Features

- Support for additional geometric types (e.g., ellipses, splines).
- Enhanced visualisation tools using libraries like Ebitengine.
- Performance optimizations for large datasets.
- Pathfinding within geometric types: Efficient algorithms for finding shortest paths constrained by geometric boundaries.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](https://github.com/mikenye/geom2d/blob/main/CONTRIBUTING.md) for details.

## License

See the [LICENSE](https://github.com/mikenye/geom2d/blob/main/LICENSE) file for details.
