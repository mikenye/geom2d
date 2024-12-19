# **Geom2D**

Geom2D is a computational geometry library for Go, designed for 2D polygon operations and other fundamental geometric types. The library is currently in its **initial development phase** and is not yet ready for production use.

## **Project Goals**

Geom2D aims to provide a robust, flexible, and efficient implementation of 2D geometric operations, featuring:

- **Geometry Types**:
  - **Point**: Basic 2D point representation.
  - **LineSegment**: Represents a line segment and supports operations such as intersection and reflection.
  - **Circle**: Support for operations like circumference, area, and intersection checks.
  - **Rectangle**: Axis-aligned bounding box with methods for containment, intersection, and transformation.
  - **Polygon (PolyTree)**: Supports polygons with holes and nested structures, with methods for orientation, correction, and Boolean operations.
- **Polygon Boolean Operations**: Union, intersection and subtraction.
- **Geometry to Geometry Relationships**: Fast and reliable algorithms for determining geometric relationships.
- **Generics**: The library leverages Go's generics to allow users to work with both integers (`int`) and floating-point (`float64`) types, offering precision and flexibility depending on the application's requirements.

## Getting Started

To install Geom2D, use go get:

```bash
go get github.com/mikenye/geom2d
```

### Example

```go
package main

import (
    "fmt"
    "github.com/mikenye/geom2d"
)

func main() {
    // Create a new point
    p := geom2d.NewPoint(3, 4)

    // Create a new circle with center p and radius 5
    c := geom2d.NewCircle(p, 5)

    // Calculate the area of the circle
    area := c.Area()

    fmt.Println("Circle Area:", area)
}
```

## Documentation

Comprehensive documentation exists in the [repository's wiki](https://github.com/mikenye/geom2d/wiki).

## Geometric Relationships

# Geometric Relationships Table

This table describes the **relationship of the left-side type (column) to the top-side type (row)**.  
Each cell indicates the valid relationship types.

| **Left ↓, Right →**         | Point                                                                    | Line Segment                                                             | Circle                                                                                                                   | Rectangle                                                                                                                | Polygon within PolyTree                                                                                                  |
|-----------------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|
| **Point**                   | RelationshipDisjoint<br>RelationshipEqual                                | RelationshipDisjoint<br>RelationshipIntersection                         | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              |
| **Line Segment**            | RelationshipDisjoint<br>RelationshipIntersection                         | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipEqual    | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy                                              |
| **Circle**                  | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains<br>RelationshipEqual | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      |
| **Rectangle**               | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains<br>RelationshipEqual | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      |
| **Polygon within PolyTree** | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContains | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains                      | RelationshipDisjoint<br>RelationshipIntersection<br>RelationshipContainedBy<br>RelationshipContains<br>RelationshipEqual |

## Acknowledgments

Geom2D builds upon the work of others and is grateful for the foundations they have laid. Specifically:

- **Martínez et al.**: Their paper on Boolean operations on polygons has been instrumental in the implementation of the Martínez algorithm in this library. See [A simple algorithm for Boolean operations on polygons](https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124).

- **Tom Wright**: The inspiration for starting this library came from Tom Wright’s repository [Provably Correct Polygon Algorithms](https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms) and his accompanying paper. While Geom2D follows its own approach, certain ideas have been influenced by his work.

- This project is a collaborative effort, with significant assistance from [OpenAI's Assistant](https://openai.com/) for brainstorming, debugging, and refining implementations.

## License

See the LICENSE file for details.
