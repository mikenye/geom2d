# Geom2D

[![Go Reference](https://pkg.go.dev/badge/github.com/mikenye/geom2d.svg)](https://pkg.go.dev/github.com/mikenye/geom2d)

Geom2D is a computational geometry library for Go, designed for 2D polygon operations and other fundamental geometric types, and is currently reaching its release candidate phase, nearing production readiness.

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
  - **Polygon (PolyTree)**: Supports polygons with holes and nested structures, with methods for orientation, correction, and Boolean operations.
- **Polygon Boolean Operations**: Union, intersection, and subtraction.
- **Geometry to Geometry Relationships**: Fast and reliable algorithms for determining geometric relationships.
- **Generics**: The library leverages Go's generics to allow users to work with both integers (`int`) and floating-point (`float64`) types, offering precision and flexibility depending on the application's requirements.

## Getting Started

To install Geom2D, use go get:

```bash
go get github.com/mikenye/geom2d
```

### Example

For detailed examples, please see the [repository's wiki](https://github.com/mikenye/geom2d/wiki), where almost every public function has an example.

```go
package main

import (
    "fmt"
    "github.com/mikenye/geom2d"
)

func areaOfCircle() {
    // Create a new point
    p := geom2d.NewPoint(3, 4)
    
    // Create a new circle with center p and radius 5
    c := geom2d.NewCircle(p, 5)
    
    // Calculate the area of the circle
    area := c.Area()
    
    // Print area
    fmt.Printf("Circle Area: %.2f\n", area)
    
    // Output:
    // Circle Area: 78.54
}

func polyIntersection() {
    // Define root contour.
    rootContour := []geom2d.Point[int]{
      geom2d.NewPoint(0, 0),
      geom2d.NewPoint(20, 0),
      geom2d.NewPoint(20, 20),
      geom2d.NewPoint(0, 20),
    }
  
    // Define hole contour within root contour.
    holeContour := []geom2d.Point[int]{
      geom2d.NewPoint(5, 5),
      geom2d.NewPoint(15, 5),
      geom2d.NewPoint(15, 15),
      geom2d.NewPoint(5, 15),
    }
  
    // Create hole polytree.
    pt1Hole, err := geom2d.NewPolyTree(holeContour, geom2d.PTHole)
    if err != nil {
	  // log.Fatal is used in the examples for simplicity and should be replaced with proper error handling in production applications.
      log.Fatal(err)
    }
  
    // Create root polytree with hole as child.
    pt1, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid, geom2d.WithChildren(pt1Hole))
    if err != nil {
      log.Fatal(err)
    }
  
    // Create a new polytree from pt1, translated by (7, 7), so there is overlap of the solid and hole regions.
    pt2 := pt1.Translate(geom2d.NewPoint(7, 7))
  
    // Perform Intersection operation, returning a PolyTree that contains only the overlapping areas
    pt3, err := pt1.BooleanOperation(pt2, geom2d.BooleanIntersection)
    if err != nil {
      log.Fatal(err)
    }
  
    // print pt3
    fmt.Println(pt3)
  
    // Output is overlapping areas, given as two solid, sibling polygons:
    // PolyTree: PTSolid
    // Contour Points: [(15, 7), (20, 7), (20, 12), (15, 12)]
    // PolyTree: PTSolid
    // Contour Points: [(7, 15), (12, 15), (12, 20), (7, 20)]
}

func main() {
	areaOfCircle()
	polyIntersection()
}
```

## Documentation

Comprehensive documentation, including detailed examples, API references, and advanced usage, is available in the [repository's wiki](https://github.com/mikenye/geom2d/wiki).

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

- **Martínez et al.**: Their paper on Boolean operations on polygons has been instrumental in the implementation of Martínez's algorithm in this library. See [A simple algorithm for Boolean operations on polygons](https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124).

- **Tom Wright**: The inspiration for starting this library came from Tom Wright’s repository [Provably Correct Polygon Algorithms](https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms) and his accompanying paper. While Geom2D follows its own approach, certain ideas have been influenced by his work.

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
