# **Geom2D**

Geom2D is a computational geometry library for Go, designed for 2D polygon operations and other fundamental geometric types. The library is currently in its **initial development phase** and is not yet ready for production use.

## **Project Goals**

Geom2D aims to provide a robust, flexible, and efficient implementation of 2D geometric operations, featuring:

- **Boolean Operations**: Union, intersection, subtraction, and more for polygons.
- **Point-in-Polygon Testing**: Fast and reliable algorithms for determining point relationships to polygons.
- **Geometry Types**:
    - **Point**: Basic 2D point representation.
    - **LineSegment**: Represents a line segment and supports operations such as intersection and reflection.
    - **Circle**: Support for operations like circumference, area, and intersection checks.
    - **Rectangle**: Axis-aligned bounding box with methods for containment, intersection, and transformation.
    - **Polygon**: Supports polygons with holes and nested structures, with methods for orientation, correction, and Boolean operations.
- **Generics**: The library leverages Go's generics to allow users to work with both integers (`int`) and floating-point (`float64`) types, offering precision and flexibility depending on the application's requirements.

## **Getting Started**

Detailed installation instructions, usage examples, and API documentation will be provided as the library approaches stability. For now, explore the codebase and experiment with the provided types and methods.

### **Example: Working with Generics**

```go
package main

import (
	"fmt"
	"github.com/mikenye/geom2d"
)

func main() {
	// Create points using int
	p1 := geom2d.NewPoint[int](2, 3)
	p2 := geom2d.NewPoint[int](5, 7)

	// Create points using float64
	p3 := geom2d.NewPoint[float64](1.2, 3.4)

	fmt.Println(p1, p2, p3)
}
```

## Acknowledgments

Geom2D builds upon the work of others and is grateful for the foundations they have laid. Specifically:

- **Martínez et al.**: Their paper on Boolean operations on polygons has been instrumental in the implementation of the Martínez algorithm in this library. See [A simple algorithm for Boolean operations on polygons](https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124).

- **Tom Wright**: The inspiration for starting this library came from Tom Wright’s repository [Provably Correct Polygon Algorithms](https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms) and his accompanying paper. While Geom2D follows its own approach, certain ideas have been influenced by his work.

- This project is a collaborative effort, with significant assistance from [OpenAI's Assistant](https://openai.com/) for brainstorming, debugging, and refining implementations.

## License

See the LICENSE file for details.
