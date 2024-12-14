package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleNewPolyTreeOption() {
	// Create a new PolyTree with a child
	child, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(3, 3),
		geom2d.NewPoint(7, 3),
		geom2d.NewPoint(7, 7),
		geom2d.NewPoint(3, 7),
	}, geom2d.PTHole)
	parent, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	}, geom2d.PTSolid, geom2d.WithChildren(child))
	fmt.Println(parent.String())
	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (5, 0), (5, 5), (0, 5)]
	//   PolyTree: PTHole
	//   Contour Points: [(1, 1), (1, 3), (3, 3), (3, 1)]
}

func ExampleWithChildren() {
	// Create a new PolyTree with a child
	child, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(3, 3),
		geom2d.NewPoint(7, 3),
		geom2d.NewPoint(7, 7),
		geom2d.NewPoint(3, 7),
	}, geom2d.PTHole)
	parent, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	}, geom2d.PTSolid, geom2d.WithChildren(child))
	fmt.Println(parent.String())
	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (5, 0), (5, 5), (0, 5)]
	//   PolyTree: PTHole
	//   Contour Points: [(1, 1), (1, 3), (3, 3), (3, 1)]
}

func ExamplePolyTree_BoundingBox() {
	// Create a PolyTree with a single polygon
	polyTree, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 20),
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 0),
	}, geom2d.PTSolid)

	// While errors are ignored in this example, please handle them appropriately in production code.

	// Calculate the bounding box of the PolyTree
	boundingBox := polyTree.BoundingBox()

	// Output the bounding box
	fmt.Println("Bounding box:", boundingBox)
	// Output:
	// Bounding box: Rectangle[(0, 0), (20, 0), (20, 20), (0, 20)]
}

func ExamplePolyTree_RelationshipToCircle() {
	// Create a PolyTree
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}, geom2d.PTHole)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define a Circle
	circle := geom2d.NewCircle(geom2d.NewPoint(50, 50), 10)

	// Determine relationships
	relationships := root.RelationshipToCircle(circle, geom2d.WithEpsilon(1e-10))

	// Output the relationships
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContains
	// Hole polygon relationship: RelationshipContains
}

func ExamplePolyTree_RelationshipToLineSegment() {
	// Create a PolyTree
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}, geom2d.PTHole)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define a LineSegment
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(10, 10), geom2d.NewPoint(90, 90))

	// Determine relationships
	relationships := root.RelationshipToLineSegment(lineSegment, geom2d.WithEpsilon(1e-10))

	// Output the relationships
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContains
	// Hole polygon relationship: RelationshipIntersection
}

func ExamplePolyTree_RelationshipToPoint() {
	// Create a PolyTree
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}, geom2d.PTHole)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define a point
	point := geom2d.NewPoint(50, 50)

	// Determine relationships
	relationships := root.RelationshipToPoint(point, geom2d.WithEpsilon(1e-10))

	// Output the relationships
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContains
	// Hole polygon relationship: RelationshipContains
}

func ExamplePolyTree_RelationshipToPolyTree() {
	// Create the first PolyTree
	pt1, err := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
	}, geom2d.PTSolid)
	if err != nil {
		fmt.Println("Error creating PolyTree 1:", err)
		return
	}

	// Create the second PolyTree
	pt2, err := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(5, 5),
		geom2d.NewPoint(5, 15),
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(15, 5),
	}, geom2d.PTSolid)
	if err != nil {
		fmt.Println("Error creating PolyTree 2:", err)
		return
	}

	// Calculate the relationships
	relationships := pt1.RelationshipToPolyTree(pt2)

	fmt.Printf("pt1's relationship to pt2: %v\n", relationships[pt1][pt2])
	// Output:
	// pt1's relationship to pt2: RelationshipIntersection
}

func ExamplePolyTree_RelationshipToRectangle() {
	// Create a PolyTree with a root polygon and a hole
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}, geom2d.PTHole)
	_ = root.AddChild(hole)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Define a rectangle
	rect := geom2d.NewRectangle(
		[]geom2d.Point[int]{
			geom2d.NewPoint(10, 10),
			geom2d.NewPoint(90, 10),
			geom2d.NewPoint(90, 90),
			geom2d.NewPoint(10, 90),
		},
	)

	// Calculate relationships
	relationships := root.RelationshipToRectangle(rect)

	// Output the relationships
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContains
	// Hole polygon relationship: RelationshipContainedBy
}

func ExamplePolygonType_String() {
	pt := geom2d.PTSolid
	fmt.Println(pt.String())
	// Output:
	// PTSolid
}

func ExamplePolyTree_String() {
	// Create a root PolyTree with a solid contour
	root, _ := geom2d.NewPolyTree(
		[]geom2d.Point[int]{
			geom2d.NewPoint(0, 0),
			geom2d.NewPoint(10, 0),
			geom2d.NewPoint(10, 10),
			geom2d.NewPoint(0, 10),
		},
		geom2d.PTSolid,
	)

	// Create a hole PolyTree inside the root
	hole, _ := geom2d.NewPolyTree(
		[]geom2d.Point[int]{
			geom2d.NewPoint(3, 3),
			geom2d.NewPoint(7, 3),
			geom2d.NewPoint(7, 7),
			geom2d.NewPoint(3, 7),
		},
		geom2d.PTHole,
	)

	// Add the hole as a child of the root
	_ = root.AddChild(hole)

	// Create an island PolyTree inside the hole
	island, _ := geom2d.NewPolyTree(
		[]geom2d.Point[int]{
			geom2d.NewPoint(4, 4),
			geom2d.NewPoint(6, 4),
			geom2d.NewPoint(6, 6),
			geom2d.NewPoint(4, 6),
		},
		geom2d.PTSolid,
	)

	// Add the island as a child of the hole
	_ = hole.AddChild(island)

	// Print the PolyTree's string representation
	fmt.Println(root.String())
	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (5, 0), (5, 5), (0, 5)]
	//   PolyTree: PTHole
	//   Contour Points: [(1, 1), (1, 3), (3, 3), (3, 1)]
	//     PolyTree: PTSolid
	//     Contour Points: [(2, 2), (3, 2), (3, 3), (2, 3)]
}
