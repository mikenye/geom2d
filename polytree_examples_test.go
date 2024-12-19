package geom2d_test

import (
	"fmt"
	"geom2d"
	"log"
	"math"
)

func ExampleNewPolyTree() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create hole polygon - slightly smaller square
	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Create island polygon - even slightly smaller square
	island, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(40, 40),
		geom2d.NewPoint(60, 40),
		geom2d.NewPoint(60, 60),
		geom2d.NewPoint(40, 60),
	}, geom2d.PTSolid)

	// Set up polygon relationships
	_ = hole.AddChild(island)
	_ = root.AddChild(hole)

	// Print polytree
	fmt.Println(root)

	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (100, 0), (100, 100), (0, 100)]
	//   PolyTree: PTHole
	//   Contour Points: [(20, 20), (20, 80), (80, 80), (80, 20)]
	//     PolyTree: PTSolid
	//     Contour Points: [(40, 40), (60, 40), (60, 60), (40, 60)]
}

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
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
	//   PolyTree: PTHole
	//   Contour Points: [(3, 3), (3, 7), (7, 7), (7, 3)]
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
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
	//   PolyTree: PTHole
	//   Contour Points: [(3, 3), (3, 7), (7, 7), (7, 3)]
}

func ExampleWithSiblings() {
	// Create the first PolyTree
	root1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create the second PolyTree
	root2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(200, 200),
		geom2d.NewPoint(300, 200),
		geom2d.NewPoint(300, 300),
		geom2d.NewPoint(200, 300),
	}, geom2d.PTSolid)

	// Create a third PolyTree and assign siblings
	root3, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(400, 400),
		geom2d.NewPoint(500, 400),
		geom2d.NewPoint(500, 500),
		geom2d.NewPoint(400, 500),
	}, geom2d.PTSolid, geom2d.WithSiblings(root1, root2))

	// Check siblings of root3
	fmt.Println("Siblings of root3:")
	for _, sibling := range root3.Siblings() {
		fmt.Println(sibling.Contour())
	}

	// Output:
	// Siblings of root3:
	// [Point[(0, 0)] Point[(100, 0)] Point[(100, 100)] Point[(0, 100)]]
	// [Point[(200, 200)] Point[(300, 200)] Point[(300, 300)] Point[(200, 300)]]
}

func ExamplePolyTree_Area() {
	// Create a square PolyTree
	contour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
	}
	poly, err := geom2d.NewPolyTree(contour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Calculate the area
	area := poly.Area()
	fmt.Printf("Area of the square: %.2f\n", area)
	// Output:
	// Area of the square: 100.00
}

func ExamplePolyTree_AsFloat32() {
	// Create an example PolyTree[int]
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	}
	polyTree, _ := geom2d.NewPolyTree(points, geom2d.PTSolid)

	fmt.Println("Original:")
	fmt.Println(polyTree)
	fmt.Println("AsFloat32:")
	fmt.Println(polyTree.AsFloat32())

	// Output:
	// Original:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
	//
	// AsFloat32:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
}

func ExamplePolyTree_AsFloat64() {
	// Create an example PolyTree[int]
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(10, 0),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(0, 10),
	}
	polyTree, _ := geom2d.NewPolyTree(points, geom2d.PTSolid)

	fmt.Println("Original:")
	fmt.Println(polyTree)
	fmt.Println("AsFloat64:")
	fmt.Println(polyTree.AsFloat64())

	// Output:
	// Original:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
	//
	// AsFloat64:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
}

func ExamplePolyTree_AsInt() {
	// Create an example PolyTree[float64]
	points := []geom2d.Point[float64]{
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10.5, 0),
		geom2d.NewPoint[float64](10.5, 10.5),
		geom2d.NewPoint[float64](0, 10.5),
	}
	polyTree, _ := geom2d.NewPolyTree(points, geom2d.PTSolid)

	fmt.Println("Original:")
	fmt.Println(polyTree)
	fmt.Println("AsInt:")
	fmt.Println(polyTree.AsInt())

	// Output:
	// Original:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10.5, 0), (10.5, 10.5), (0, 10.5)]
	//
	// AsInt:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
}

func ExamplePolyTree_AsIntRounded() {
	// Create an example PolyTree[float64]
	points := []geom2d.Point[float64]{
		geom2d.NewPoint[float64](0, 0),
		geom2d.NewPoint[float64](10.5, 0),
		geom2d.NewPoint[float64](10.5, 10.5),
		geom2d.NewPoint[float64](0, 10.5),
	}
	polyTree, _ := geom2d.NewPolyTree(points, geom2d.PTSolid)

	fmt.Println("Original:")
	fmt.Println(polyTree)
	fmt.Println("AsIntRounded:")
	fmt.Println(polyTree.AsIntRounded())

	// Output:
	// Original:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (10.5, 0), (10.5, 10.5), (0, 10.5)]
	//
	// AsIntRounded:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (11, 0), (11, 11), (0, 11)]
}

func ExamplePolyTree_BooleanOperation_intersection() {

	// define root contour
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(20, 0),
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(0, 20),
	}

	// define hole contour
	holeContour := []geom2d.Point[int]{
		geom2d.NewPoint(5, 5),
		geom2d.NewPoint(15, 5),
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(5, 15),
	}

	// create hole polytree
	pt1Hole, err := geom2d.NewPolyTree(holeContour, geom2d.PTHole)
	if err != nil {
		log.Fatal(err)
	}

	// create root polytree with hole as child
	pt1, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid, geom2d.WithChildren(pt1Hole))
	if err != nil {
		log.Fatal(err)
	}

	// create a new polytree from pt1, translated by (7, 7)
	pt2 := pt1.Translate(geom2d.NewPoint(7, 7))

	// perform Union operation
	pt3, err := pt1.BooleanOperation(pt2, geom2d.BooleanIntersection)
	if err != nil {
		log.Fatal(err)
	}

	// print pt3
	fmt.Println(pt3)

	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(15, 7), (20, 7), (20, 12), (15, 12)]
	// PolyTree: PTSolid
	// Contour Points: [(7, 15), (12, 15), (12, 20), (7, 20)]
}

func ExamplePolyTree_BooleanOperation_subtract() {

	// define root contour
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(20, 0),
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(0, 20),
	}

	// define hole contour
	holeContour := []geom2d.Point[int]{
		geom2d.NewPoint(5, 5),
		geom2d.NewPoint(15, 5),
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(5, 15),
	}

	// create hole polytree
	pt1Hole, err := geom2d.NewPolyTree(holeContour, geom2d.PTHole)
	if err != nil {
		log.Fatal(err)
	}

	// create root polytree with hole as child
	pt1, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid, geom2d.WithChildren(pt1Hole))
	if err != nil {
		log.Fatal(err)
	}

	// create a new polytree from pt1, translated by (7, 7)
	pt2 := pt1.Translate(geom2d.NewPoint(7, 7))

	// perform Union operation
	pt3, err := pt1.BooleanOperation(pt2, geom2d.BooleanSubtraction)
	if err != nil {
		log.Fatal(err)
	}

	// print pt3
	fmt.Println(pt3)

	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (20, 0), (20, 7), (15, 7), (15, 5), (5, 5), (5, 15), (7, 15), (7, 20), (0, 20)]
	// PolyTree: PTSolid
	// Contour Points: [(15, 12), (20, 12), (20, 20), (12, 20), (12, 15), (15, 15)]
}

func ExamplePolyTree_BooleanOperation_union() {

	// define root contour
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(20, 0),
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(0, 20),
	}

	// define hole contour
	holeContour := []geom2d.Point[int]{
		geom2d.NewPoint(5, 5),
		geom2d.NewPoint(15, 5),
		geom2d.NewPoint(15, 15),
		geom2d.NewPoint(5, 15),
	}

	// create hole polytree
	pt1Hole, err := geom2d.NewPolyTree(holeContour, geom2d.PTHole)
	if err != nil {
		log.Fatal(err)
	}

	// create root polytree with hole as child
	pt1, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid, geom2d.WithChildren(pt1Hole))
	if err != nil {
		log.Fatal(err)
	}

	// create a new polytree from pt1, translated by (7, 7)
	pt2 := pt1.Translate(geom2d.NewPoint(7, 7))

	// perform Union operation
	pt3, err := pt1.BooleanOperation(pt2, geom2d.BooleanUnion)
	if err != nil {
		log.Fatal(err)
	}

	// print pt3
	fmt.Println(pt3)

	// Output:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (20, 0), (20, 7), (27, 7), (27, 27), (7, 27), (7, 20), (0, 20)]
	//   PolyTree: PTHole
	//   Contour Points: [(5, 15), (7, 15), (7, 7), (15, 7), (15, 5), (5, 5)]
	//   PolyTree: PTHole
	//   Contour Points: [(12, 15), (15, 15), (15, 12), (12, 12)]
	//   PolyTree: PTHole
	//   Contour Points: [(20, 20), (12, 20), (12, 22), (22, 22), (22, 12), (20, 12)]
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

func ExamplePolyTree_Children() {
	// Create a root PolyTree
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}
	root, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating root: %v\n", err)
		return
	}

	// Create a child PolyTree
	childContour := []geom2d.Point[int]{
		geom2d.NewPoint(25, 25),
		geom2d.NewPoint(25, 75),
		geom2d.NewPoint(75, 75),
		geom2d.NewPoint(75, 25),
	}
	child, err := geom2d.NewPolyTree(childContour, geom2d.PTHole)
	if err != nil {
		fmt.Printf("Error creating child: %v\n", err)
		return
	}

	// Add the child to the root
	if err := root.AddChild(child); err != nil {
		fmt.Printf("Error adding child: %v\n", err)
		return
	}

	// Retrieve and print the children of the root
	children := root.Children()
	fmt.Printf("Number of children: %d\n", len(children))
	for i, c := range children {
		fmt.Printf("Child %d area: %.2f\n", i+1, c.Area())
	}
	// Output:
	// Number of children: 1
	// Child 1 area: 2500.00
}

func ExamplePolyTree_Contour() {
	// Create a PolyTree representing a square
	pt, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	// Get the contour of the PolyTree
	contour := pt.Contour()

	// Print the contour points
	fmt.Println("Contour Points:")
	for _, point := range contour {
		fmt.Println(point)
	}
	// Output:
	// Contour Points:
	// Point[(0, 0)]
	// Point[(100, 0)]
	// Point[(100, 100)]
	// Point[(0, 100)]
}

func ExamplePolyTree_Edges() {
	// Create a PolyTree
	contour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}
	poly, err := geom2d.NewPolyTree(contour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating PolyTree: %v\n", err)
		return
	}

	// Retrieve and print edges
	edges := poly.Edges()
	fmt.Printf("Number of edges: %d\n", len(edges))
	for i, edge := range edges {
		fmt.Printf("Edge %d: Start %v, End %v\n", i+1, edge.Start(), edge.End())
	}
	// Output:
	// Number of edges: 4
	// Edge 1: Start Point[(0, 0)], End Point[(100, 0)]
	// Edge 2: Start Point[(100, 0)], End Point[(100, 100)]
	// Edge 3: Start Point[(100, 100)], End Point[(0, 100)]
	// Edge 4: Start Point[(0, 100)], End Point[(0, 0)]
}

func ExamplePolyTree_Eq() {
	// Create the first PolyTree (root with one child and one sibling)
	root1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	child1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	sibling1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(200, 200),
		geom2d.NewPoint(300, 200),
		geom2d.NewPoint(300, 300),
		geom2d.NewPoint(200, 300),
	}, geom2d.PTSolid)

	_ = root1.AddChild(child1)
	_ = root1.AddSibling(sibling1)

	// Create the second PolyTree with identical structure and geometry
	root2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	child2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	sibling2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(200, 200),
		geom2d.NewPoint(300, 200),
		geom2d.NewPoint(300, 300),
		geom2d.NewPoint(200, 300),
	}, geom2d.PTSolid)

	_ = root2.AddChild(child2)
	_ = root2.AddSibling(sibling2)

	// Compare the two PolyTrees
	equal, mismatches := root1.Eq(root2)

	// Print the results
	fmt.Printf("Are the PolyTrees equal? %v\n", equal)
	fmt.Printf("Mismatch bitmask: %v\n", mismatches)

	// Output:
	// Are the PolyTrees equal? true
	// Mismatch bitmask: 0
}

func ExamplePolyTree_Hull() {
	// Create a PolyTree with a triangular contour
	contour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(50, 100),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(50, 50),
	}
	poly, err := geom2d.NewPolyTree(contour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating PolyTree: %v\n", err)
		return
	}

	// Retrieve and print the hull
	hull := poly.Hull()
	fmt.Printf("Convex Hull:\n")
	for _, point := range hull {
		fmt.Printf("Point: %v\n", point)
	}
	// Output:
	// Convex Hull:
	// Point: Point[(0, 0)]
	// Point: Point[(100, 0)]
	// Point: Point[(50, 100)]
}

func ExamplePolyTree_Intersects() {
	// Create the first PolyTree - a square
	root1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create the second PolyTree - a smaller square inside the first
	root2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTSolid)

	// Create a third PolyTree - disjoint from the first two
	root3, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(200, 200),
		geom2d.NewPoint(300, 200),
		geom2d.NewPoint(300, 300),
		geom2d.NewPoint(200, 300),
	}, geom2d.PTSolid)

	// Check intersections
	fmt.Printf("Root1 intersects Root2: %v\n", root1.Intersects(root2)) // Expect true
	fmt.Printf("Root1 intersects Root3: %v\n", root1.Intersects(root3)) // Expect false
	fmt.Printf("Root2 intersects Root3: %v\n", root2.Intersects(root3)) // Expect false

	// Output:
	// Root1 intersects Root2: true
	// Root1 intersects Root3: false
	// Root2 intersects Root3: false
}

func ExamplePolyTree_IsRoot() {
	// Create a root PolyTree
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}
	root, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating root PolyTree: %v\n", err)
		return
	}

	// Create a child PolyTree
	childContour := []geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}
	child, err := geom2d.NewPolyTree(childContour, geom2d.PTHole)
	if err != nil {
		fmt.Printf("Error creating child PolyTree: %v\n", err)
		return
	}
	err = root.AddChild(child)
	if err != nil {
		fmt.Printf("Error adding child to root: %v\n", err)
		return
	}

	// Print the root status
	fmt.Printf("Is the root a root? %v\n", root.IsRoot())
	fmt.Printf("Is the child a root? %v\n", child.IsRoot())
	// Output:
	// Is the root a root? true
	// Is the child a root? false
}

func ExamplePolyTree_Len() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create hole polygon - slightly smaller square
	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Create island polygon - even slightly smaller square
	island, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(40, 40),
		geom2d.NewPoint(60, 40),
		geom2d.NewPoint(60, 60),
		geom2d.NewPoint(40, 60),
	}, geom2d.PTSolid)

	// Establish relationships
	_ = hole.AddChild(island)
	_ = root.AddChild(hole)

	// Output the total number of PolyTree nodes
	fmt.Printf("Total nodes in root: %d\n", root.Len())
	fmt.Printf("Total nodes in hole: %d\n", hole.Len())
	fmt.Printf("Total nodes in island: %d\n", island.Len())

	// Output:
	// Total nodes in root: 3
	// Total nodes in hole: 2
	// Total nodes in island: 1
}

func ExamplePolyTree_Nodes() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create hole polygon - slightly smaller square
	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Create island polygon - even slightly smaller square
	island, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(40, 40),
		geom2d.NewPoint(60, 40),
		geom2d.NewPoint(60, 60),
		geom2d.NewPoint(40, 60),
	}, geom2d.PTSolid)

	// Establish relationships
	_ = hole.AddChild(island)
	_ = root.AddChild(hole)

	// Use Nodes to iterate over all polygons in the PolyTree
	fmt.Println("Iterating over all nodes in the PolyTree:")
	for node := range root.Nodes {
		fmt.Printf("Polygon Type: %s\n", node.PolygonType())
		fmt.Printf("Contour: %v\n", node.Contour())
	}

	// Output:
	// Iterating over all nodes in the PolyTree:
	// Polygon Type: PTSolid
	// Contour: [Point[(0, 0)] Point[(100, 0)] Point[(100, 100)] Point[(0, 100)]]
	// Polygon Type: PTHole
	// Contour: [Point[(20, 20)] Point[(20, 80)] Point[(80, 80)] Point[(80, 20)]]
	// Polygon Type: PTSolid
	// Contour: [Point[(40, 40)] Point[(60, 40)] Point[(60, 60)] Point[(40, 60)]]
}

func ExamplePolyTree_Parent() {
	// Create a root PolyTree
	rootContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}
	root, err := geom2d.NewPolyTree(rootContour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating root PolyTree: %v\n", err)
		return
	}

	// Create a child PolyTree
	childContour := []geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(20, 80),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(80, 20),
	}
	child, err := geom2d.NewPolyTree(childContour, geom2d.PTHole)
	if err != nil {
		fmt.Printf("Error creating child PolyTree: %v\n", err)
		return
	}
	err = root.AddChild(child)
	if err != nil {
		fmt.Printf("Error adding child to root: %v\n", err)
		return
	}

	// Print the parent of the child node
	if parent := child.Parent(); parent != nil {
		fmt.Println("Child's parent exists.")
	} else {
		fmt.Println("Child's parent is nil.")
	}

	// Print the parent of the root node
	if parent := root.Parent(); parent != nil {
		fmt.Println("Root's parent exists.")
	} else {
		fmt.Println("Root's parent is nil.")
	}
	// Output:
	// Child's parent exists.
	// Root's parent is nil.
}

func ExamplePolyTree_Perimeter() {
	// Create a PolyTree representing a triangle
	triangleContour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 0),
	}
	triangle, err := geom2d.NewPolyTree(triangleContour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating PolyTree: %v\n", err)
		return
	}

	// Calculate the perimeter
	perimeter := triangle.Perimeter()

	fmt.Printf("The perimeter of the triangle is: %.2f\n", perimeter)
	// Output:
	// The perimeter of the triangle is: 34.14
}

func ExamplePolyTree_PolygonType() {
	// Define a polygon contour
	contour := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 10),
		geom2d.NewPoint(10, 10),
		geom2d.NewPoint(10, 0),
	}

	// Create a solid polygon
	polyTree, err := geom2d.NewPolyTree(contour, geom2d.PTSolid)
	if err != nil {
		fmt.Printf("Error creating PolyTree: %v\n", err)
		return
	}

	// Get the PolygonType
	polygonType := polyTree.PolygonType()

	// Print the PolygonType
	fmt.Printf("PolygonType: %v\n", polygonType)
	// Output:
	// PolygonType: PTSolid
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
	circle := geom2d.NewCircle(geom2d.NewPoint(50, 50), 40)

	// Determine relationships
	relationships := root.RelationshipToCircle(circle, geom2d.WithEpsilon(1e-10))

	// Output the relationships
	fmt.Printf("Root polygon relationship: %v\n", relationships[root])
	fmt.Printf("Hole polygon relationship: %v\n", relationships[hole])
	// Output:
	// Root polygon relationship: RelationshipContains
	// Hole polygon relationship: RelationshipIntersection
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

func ExamplePolyTree_Root() {
	// Create a root/parent polygon
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create a child polygon
	child, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Create a grandchild polygon
	grandchild, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(40, 40),
		geom2d.NewPoint(60, 40),
		geom2d.NewPoint(60, 60),
		geom2d.NewPoint(40, 60),
	}, geom2d.PTSolid)

	// Establish relationships
	_ = child.AddChild(grandchild)
	_ = root.AddChild(child)

	// Retrieve the root node from the grandchild
	fmt.Println(grandchild.Root() == root) // true

	// Output:
	// true
}

func ExamplePolyTree_Rotate() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Define pivot point (0, 0) and rotation angle (90° counterclockwise)
	pivot := geom2d.NewPoint(0, 0)
	angle := math.Pi / 2

	// Perform rotation
	rotated := root.Rotate(pivot, angle)

	// Print before and after rotation
	fmt.Println("Before rotation:")
	fmt.Println(root)
	fmt.Println("After 90° counterclockwise rotation:")
	fmt.Println(rotated.AsIntRounded())

	// Output:
	// Before rotation:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (100, 0), (100, 100), (0, 100)]
	//
	// After 90° counterclockwise rotation:
	// PolyTree: PTSolid
	// Contour Points: [(-100, 0), (0, 0), (0, 100), (-100, 100)]
}

func ExamplePolyTree_Scale() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create hole polygon - slightly smaller square
	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Create island polygon - even slightly smaller square
	island, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(40, 40),
		geom2d.NewPoint(60, 40),
		geom2d.NewPoint(60, 60),
		geom2d.NewPoint(40, 60),
	}, geom2d.PTSolid)

	// Set up polygon relationships
	_ = hole.AddChild(island)
	_ = root.AddChild(hole)

	// Scale by a factor of 2 with origin 0,0
	scaled := root.Scale(geom2d.NewPoint(0, 0), 2)

	// Print output
	fmt.Println("Before scaling:")
	fmt.Println(root)
	fmt.Println("After scaling:")
	fmt.Println(scaled)

	// Output:
	// Before scaling:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (100, 0), (100, 100), (0, 100)]
	//   PolyTree: PTHole
	//   Contour Points: [(20, 20), (20, 80), (80, 80), (80, 20)]
	//     PolyTree: PTSolid
	//     Contour Points: [(40, 40), (60, 40), (60, 60), (40, 60)]
	//
	// After scaling:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (200, 0), (200, 200), (0, 200)]
	//   PolyTree: PTHole
	//   Contour Points: [(40, 160), (160, 160), (160, 40), (40, 40)]
	//     PolyTree: PTSolid
	//     Contour Points: [(80, 80), (120, 80), (120, 120), (80, 120)]
}

func ExamplePolyTree_Siblings() {
	// Create a root polygon
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(0, 100),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(100, 0),
	}, geom2d.PTSolid)

	// Create sibling polygons
	sibling1, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(150, 150),
		geom2d.NewPoint(150, 250),
		geom2d.NewPoint(250, 250),
		geom2d.NewPoint(250, 150),
	}, geom2d.PTSolid)

	sibling2, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(300, 300),
		geom2d.NewPoint(300, 400),
		geom2d.NewPoint(400, 400),
		geom2d.NewPoint(400, 300),
	}, geom2d.PTSolid)

	// Add siblings
	_ = root.AddSibling(sibling1)
	_ = root.AddSibling(sibling2)

	// Note: While errors are ignored in this example for simplicity, it is important to handle errors properly in
	// production code to ensure robustness and reliability.

	// Get siblings of root
	siblings := root.Siblings()
	for _, sibling := range siblings {
		fmt.Println(sibling.Contour())
	}
	// Output:
	// [Point[(150, 150)] Point[(250, 150)] Point[(250, 250)] Point[(150, 250)]]
	// [Point[(300, 300)] Point[(400, 300)] Point[(400, 400)] Point[(300, 400)]]
}

func ExamplePolyTree_Translate() {
	// Create root/parent polygon - large square
	root, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid)

	// Create hole polygon - smaller square
	hole, _ := geom2d.NewPolyTree([]geom2d.Point[int]{
		geom2d.NewPoint(20, 20),
		geom2d.NewPoint(80, 20),
		geom2d.NewPoint(80, 80),
		geom2d.NewPoint(20, 80),
	}, geom2d.PTHole)

	// Add hole to root
	_ = root.AddChild(hole)

	// Translate the entire PolyTree by (10, 10)
	translated := root.Translate(geom2d.NewPoint(10, 10))

	// Print before and after translation
	fmt.Println("Before translation:")
	fmt.Println(root)
	fmt.Println("After translation:")
	fmt.Println(translated)

	// Output:
	// Before translation:
	// PolyTree: PTSolid
	// Contour Points: [(0, 0), (100, 0), (100, 100), (0, 100)]
	//   PolyTree: PTHole
	//   Contour Points: [(20, 20), (20, 80), (80, 80), (80, 20)]
	//
	// After translation:
	// PolyTree: PTSolid
	// Contour Points: [(10, 10), (110, 10), (110, 110), (10, 110)]
	//   PolyTree: PTHole
	//   Contour Points: [(30, 90), (90, 90), (90, 30), (30, 30)]
}

func ExamplePolygonType_String() {
	fmt.Println(geom2d.PTSolid.String())
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
	// Contour Points: [(0, 0), (10, 0), (10, 10), (0, 10)]
	//   PolyTree: PTHole
	//   Contour Points: [(3, 3), (3, 7), (7, 7), (7, 3)]
	//     PolyTree: PTSolid
	//     Contour Points: [(4, 4), (6, 4), (6, 6), (4, 6)]
}
