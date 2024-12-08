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
