# Geom2D

[![Go Reference](https://pkg.go.dev/badge/github.com/mikenye/geom2d.svg)](https://pkg.go.dev/github.com/mikenye/geom2d)
[![codecov](https://codecov.io/gh/mikenye/geom2d/graph/badge.svg?token=Z73XKN7N8D)](https://codecov.io/gh/mikenye/geom2d)

Geom2D is a computational geometry library for Go, designed for 2D polygon operations and other fundamental geometric types, and is currently reaching its release candidate phase, nearing production readiness.

**As of April 2025, this package is actively in development.** The base geometric types (`circle`, `linesegment`, `point` and `rectangle`) are well documented and functional. The sweep line method of finding line segment intersections (`FindIntersections`) now supports line segments in general position, including handling edge cases like collinear segments and segment endpoints on other segments. For simpler use cases or verification, the `FindIntersectionsBruteForce` method (O(n²)) is also available.

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
  - ✅ **Point**: Complete 2D point representation with vector operations, transformations, and distance calculations.
  - ✅ **LineSegment**: Fully functional line segments with intersection detection (both brute force and sweep line methods), containment checks, and transformations.
  - ✅ **Circle**: Complete circle operations including area, circumference, relationships with points, and transformations.
  - ✅ **Rectangle**: Axis-aligned bounding box with methods for containment, intersection, and transformation.
  - 🚧 **Polygons**: Support polygons with holes and nested structures, with methods for orientation, correction, and Boolean operations.
- 🚧**Polygon Boolean Operations**: Union, intersection, and subtraction.
- ✅**Geometry to Geometry Relationships**: Well-documented, consistent relationship types (Disjoint, Intersection, ContainedBy, Contains, Equal) between geometric entities.

## Getting Started

To install Geom2D, use go get:

```bash
go get github.com/mikenye/geom2d
```

## Documentation & Examples

The library is thoroughly documented with Go doc comments. For comprehensive API documentation and usage examples, visit the [geom2d documentation at the Go Package Discovery and Documentation site](https://pkg.go.dev/github.com/mikenye/geom2d).

### Key Packages

- **geom2d** - Core package with global settings like epsilon for floating-point precision
- **point** - 2D point representation with vector operations
- **linesegment** - Line segment operations including intersection detection
- **circle** - Circle operations
- **rectangle** - Axis-aligned rectangle operations
- **types** - Common types and relationships between geometric entities
- **numeric** - Utilities for handling floating-point precision

## Geometric Relationships

Geom2D provides a consistent system for expressing spatial relationships between geometric entities, using standardized relationship types:

- **Disjoint**: The geometries do not intersect or overlap.
- **Intersection**: The geometries share common points (without either containing the other).
- **ContainedBy**: The first geometry is fully contained by the second.
- **Contains**: The first geometry fully contains the second.
- **Equal**: The geometries are equivalent.

Currently implemented relationship methods:
- Point to Point: Equality or disjoint
- Point to shapes: All geometric types implement `RelationshipToPoint` to determine how a point relates to them

Future development will expand the relationship system to include full inter-type relationships (e.g., Circle to LineSegment, Rectangle to Circle, etc.).

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](https://github.com/mikenye/geom2d/blob/main/CONTRIBUTING.md) for details.

## License

See the [LICENSE](https://github.com/mikenye/geom2d/blob/main/LICENSE) file for details.
