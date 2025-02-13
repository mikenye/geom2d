package linesegment

import (
	"cmp"
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"math"
	"strings"
)

type statusStructureRBT struct {
	structure  *rbt.Tree
	sweepPoint point.Point[float64]
	epsilon    float64
}

func (S *statusStructureRBT) containsPoint(
	p point.Point[float64],
	seg LineSegment[float64],
	LofPset map[LineSegment[float64]]bool,
	CofPset map[LineSegment[float64]]bool,
) bool {

	//log.Println("checking:", seg)

	// if we've found an upper p, do nothing as these are already stored with the event p
	if seg.Start().Eq(p) {
		// do nothing
		//log.Println("start point matches, do nothing")

		// if we've found a lower p, add to LofP
	} else if seg.End().Eq(p) {
		//log.Println("end point matches, add to L(p)")
		LofPset[seg] = false

		// if we've found a segment that contains p on its interior, then add to CofP
	} else if seg.ContainsPoint(p, options.WithEpsilon(S.epsilon)) {
		//log.Println("contains point, add to C(p)")
		CofPset[seg] = false

		// if we haven't found any of these, then we're outside of the segments containing the p and we should
		// break out of the loop
	} else {
		//log.Println("does not contain point, breaking")
		return false
	}

	return true
}

func (S *statusStructureRBT) FindCofPAndLofP(p point.Point[float64]) (CofP, LofP []LineSegment[float64]) {

	// DEBUGGING:
	//debugLog := strings.Builder{}
	////log.Println("[FindCofPAndLofP] Entered")
	//defer func() {
	//	log.Println(debugLog.String())
	//	log.Println("[FindCofPAndLofP] Exited")
	//}()

	// create map-based sets to remove duplication
	LofPset := make(map[LineSegment[float64]]bool)
	CofPset := make(map[LineSegment[float64]]bool)

	// Prepare "dummy" key to search on.
	// We need to set entryType to statusStructureEntryFindPointNeighbors so that the comparison function knows
	// to compare the line segments in the status structure to a p, instead of via the normal process.
	// The only fields needed should be entryType and sweepLinePoint.
	k := statusStructureEntry{
		entryType:      statusStructureEntryFindPointNeighbors,
		sweepLinePoint: p,
	}

	// find floor
	floor, floorFound := S.structure.Floor(k)
	//log.Println("floorFound:", floorFound)
	//log.Println("floor:", floor)
	if floorFound {

		// add floor
		_ = S.containsPoint(p, floor.Key.(statusStructureEntry).segment, LofPset, CofPset)

		// add segments below floor
		iter := S.structure.IteratorAt(floor)
		for iter.Prev() {
			seg := iter.Key().(statusStructureEntry).segment
			found := S.containsPoint(p, seg, LofPset, CofPset)
			if !found {
				break
			}
		}
	}

	// find ceiling
	ceil, ceilFound := S.structure.Ceiling(k)
	//log.Println("ceilFound:", ceilFound)
	//log.Println("ceil:", ceil)
	if ceilFound {

		// add ceiling
		_ = S.containsPoint(p, ceil.Key.(statusStructureEntry).segment, LofPset, CofPset)

		iter := S.structure.IteratorAt(ceil)
		for iter.Next() {
			seg := iter.Key().(statusStructureEntry).segment
			found := S.containsPoint(p, seg, LofPset, CofPset)
			if !found {
				break
			}
		}
	}

	// convert maps to slices
	for l := range CofPset {
		CofP = append(CofP, l)
	}
	for l := range LofPset {
		LofP = append(LofP, l)
	}

	return CofP, LofP
}

func (S *statusStructureRBT) FindNeighborsOfPoint(p point.Point[float64]) (left, right *LineSegment[float64]) {

	// Prepare "dummy" key to search on.
	// We need to set entryType to statusStructureEntryFindPointNeighbors so that the comparison function knows
	// to compare the line segments in the status structure to a p, instead of via the normal process.
	// The only fields needed should be entryType and sweepLinePoint.
	k := statusStructureEntry{
		entryType:      statusStructureEntryFindPointNeighbors,
		sweepLinePoint: p,
	}

	// find floor and ceiling
	floor, _ := S.structure.Floor(k)
	ceil, _ := S.structure.Ceiling(k)

	if floor != nil {
		l := floor.Key.(statusStructureEntry).segment
		left = &l
	}
	if ceil != nil {
		r := ceil.Key.(statusStructureEntry).segment
		right = &r
	}

	return left, right
}

func (S *statusStructureRBT) FindNeighborsOfUofPAndCofP(
	UofP, CofP []LineSegment[float64],
) (
	sPrime, sL, sDoublePrime, sR *LineSegment[float64],
) {

	// Build and sort UofP ∪ CofP
	UCofP := newStatusStructureRBT(S.sweepPoint, options.WithEpsilon(S.epsilon))
	for _, seg := range append(UofP, CofP...) {
		UCofP.Insert(seg)
	}

	// firstly find the leftmost & rightmost segment of UofP ∪ CofP
	sPrimeNode := UCofP.structure.Left()
	sDoublePrimeNode := UCofP.structure.Right()

	// extract sPrime from sPrimeNode if not nil
	if sPrimeNode != nil {
		sPrimeTmp := sPrimeNode.Key.(statusStructureEntry).segment
		sPrime = &sPrimeTmp

		// find sL, the left neighbor of sPrime
		iter := S.structure.IteratorAt(S.structure.GetNode(sPrimeNode.Key))
		if iter.Node() != nil && iter.Prev() {
			sLTmp := iter.Key().(statusStructureEntry).segment
			sL = &sLTmp
		}
	}

	// extract sDoublePrime from sDoublePrimeNode if not nil
	if sDoublePrimeNode != nil {
		sDoublePrimeTmp := sDoublePrimeNode.Key.(statusStructureEntry).segment
		sDoublePrime = &sDoublePrimeTmp

		// find sR, the right neighbor of sDoublePrime
		iter := S.structure.IteratorAt(S.structure.GetNode(sDoublePrimeNode.Key))
		if iter.Node() != nil && iter.Next() {
			sRTmp := iter.Key().(statusStructureEntry).segment
			sR = &sRTmp
		}
	}

	return sPrime, sL, sDoublePrime, sR
}

func (S *statusStructureRBT) Insert(seg LineSegment[float64]) {

	e := statusStructureEntry{}

	e.entryType = statusStructureEntryNormal

	// static fields
	e.segment = seg
	e.slope = seg.Slope()
	e.isHorizontal = e.slope == 0
	e.isVertical = math.IsNaN(e.slope)
	e.epsilon = S.epsilon

	// dynamic fields
	e.update(S.sweepPoint)

	// insert
	S.structure.Put(e, nil)
}

func (S *statusStructureRBT) Remove(seg LineSegment[float64]) {
	e := statusStructureEntry{}

	// static fields
	e.segment = seg
	e.slope = seg.Slope()
	e.isHorizontal = e.slope == 0
	e.isVertical = math.IsNaN(e.slope)

	// todo: move epsilon to S
	//e.epsilon = epsilon

	// dynamic fields
	e.update(S.sweepPoint)

	// remove
	S.structure.Remove(e)
}

func (S *statusStructureRBT) String() string {
	out := strings.Builder{}
	iter := S.structure.Iterator()
	i := 0
	for iter.Next() {
		k := iter.Key().(statusStructureEntry)
		out.WriteString(fmt.Sprintf(
			"Status Structure Item %d at %s: %s (s=%f, v=%t, h=%t)\n",
			i,
			S.sweepPoint,
			k.segment,
			k.slope,
			k.isVertical,
			k.isHorizontal,
		))
		i++
	}
	return out.String()
}

func (S *statusStructureRBT) Update(p point.Point[float64]) sweepLineStatusStructure {
	newS := newStatusStructureRBT(p, options.WithEpsilon(S.epsilon))
	iter := S.structure.Iterator()
	for iter.Next() {
		// todo: this will re-calculate slope, which isnt needed
		newS.Insert(iter.Key().(statusStructureEntry).segment)
	}
	return newS
}

type statusStructureEntryType uint8

const (
	statusStructureEntryNormal statusStructureEntryType = iota
	statusStructureEntryFindPointNeighbors
)

type statusStructureEntry struct {

	// segment is the actual line segment stored in the status structure.
	// This is a static field.
	segment LineSegment[float64]

	// sweepLinePoint is the p where the sweep line is currently.
	// It is a dynamic field, set during insert and update.
	sweepLinePoint point.Point[float64]

	// slope is the slope of the line segment.
	// This is a static field, and is calculated and set on insert.
	slope float64

	// epsilon is used during the comparison calculations.
	// This is a static field, and is set on insert.
	epsilon float64

	// currX is the X value where the sweep line crosses the segment.
	// It is a dynamic field, set during insert and update.
	currX float64

	// entryType allows the comparison function to determine whether we are finding a p or
	// line segment in the status structure. This is a static field.
	entryType statusStructureEntryType

	// containsEvent is set to true if the line segment contains the current event p (the exact sweep line p).
	// It is a dynamic field, set during insert and update.
	containsEvent bool

	// isHorizontal is set to true if the line segment is horizontal.
	// This is a static field, and is calculated and set on insert.
	isHorizontal bool

	// isVertical is set to true if the line segment is vertical.
	// This is a static field, and is calculated and set on insert.
	isVertical bool
}

func (e *statusStructureEntry) update(sweepLinePoint point.Point[float64]) {
	e.sweepLinePoint = sweepLinePoint
	e.currX = e.segment.XAtY(e.sweepLinePoint.Y())
	e.containsEvent = e.segment.ContainsPoint(e.sweepLinePoint, options.WithEpsilon(e.epsilon))
}

func statusStructureComparator(sweepPointPtr *point.Point[float64], epsilonPtr *float64) func(a, b interface{}) int {
	return func(a, b interface{}) int {

		//debugLog := strings.Builder{}
		//debugLog.WriteString("[statusStructureComparator] ")
		//defer func() {
		//	log.Println(debugLog.String())
		//}()
		//debugPrintOrder := func(i int) string {
		//	if i < 0 {
		//		return "A before B"
		//	}
		//	if i > 0 {
		//		return "B before A"
		//	}
		//	return "A and B are Equal"
		//}

		// todo: is this neccessary?
		sweepPointPtr := sweepPointPtr
		epsilonPtr := epsilonPtr

		A := a.(statusStructureEntry)
		B := b.(statusStructureEntry)
		p := *sweepPointPtr

		if A.entryType == statusStructureEntryFindPointNeighbors {
			//debugLog.WriteString(fmt.Sprintf("mode: statusStructureEntryFindPointNeighbors, %s vs %s: ", A.sweepLinePoint, B.segment))
			if B.segment.ContainsPoint(A.sweepLinePoint, options.WithEpsilon(*epsilonPtr)) {
				//debugLog.WriteString(debugPrintOrder(0))
				return 0
			}
			if B.currX < A.sweepLinePoint.X() {
				//debugLog.WriteString(debugPrintOrder(1))
				return 1
			}
			//debugLog.WriteString(debugPrintOrder(-1))
			return -1
		}
		if B.entryType == statusStructureEntryFindPointNeighbors {
			//debugLog.WriteString(fmt.Sprintf("mode: statusStructureEntryFindPointNeighbors, %s vs %s: ", A.segment, B.sweepLinePoint))
			if A.segment.ContainsPoint(B.sweepLinePoint, options.WithEpsilon(*epsilonPtr)) {
				//debugLog.WriteString(debugPrintOrder(0))
				return 0
			}
			if A.currX < B.sweepLinePoint.X() {
				//debugLog.WriteString(debugPrintOrder(-1))
				return -1
			}
			//debugLog.WriteString(debugPrintOrder(1))
			return 1
		}

		if A.entryType == statusStructureEntryNormal && B.entryType == statusStructureEntryNormal {

			// todo: try running with this commented out
			// do we need to update any of a's dynamic fields
			if !A.sweepLinePoint.Eq(p) {
				//debugLog.WriteString("Updated A. ")
				A.update(p)
			}

			// do we need to update any of b's dynamic fields
			if !B.sweepLinePoint.Eq(p) {
				//debugLog.WriteString("Updated B. ")
				B.update(p)
			}

			aX := A.currX
			bX := B.currX

			// for horizontal lines, artificially truncate start position to p,
			// since we don't care about anything to the left, as that is considered above the sweep line
			if A.isHorizontal {
				aX = p.X()
			}
			if B.isHorizontal {
				bX = p.X()
			}

			//debugLog.WriteString(fmt.Sprintf(
			//	"Comparing A (p=%s x=%f, s=%f, v=%t, h=%t, c=%t) to B (p=%s, x=%f, s=%f, v=%t, h=%t, c=%t) at %s: ",
			//	A.segment.String(),
			//	aX,
			//	A.slope,
			//	A.isVertical,
			//	A.isHorizontal,
			//	A.containsEvent,
			//	B.segment.String(),
			//	bX,
			//	B.slope,
			//	B.isVertical,
			//	B.isHorizontal,
			//	B.containsEvent,
			//	p.String(),
			//))

			if A.segment.Eq(B.segment, options.WithEpsilon(A.epsilon)) {
				//debugLog.WriteString("segments are equal")
				return 0
			}

			// Vertical segment ordering logic: Handle cases where a vertical segment intersects a diagonal one.
			// if A is vertical *AND* A contains the event p *AND* B is diagonal *AND* B contains the event p:
			if A.isVertical && A.containsEvent && !B.isVertical && !B.isHorizontal && B.containsEvent {
				//debugLog.WriteString("A is vertical *AND* A contains the event p *AND* B is diagonal *AND* B contains the event p, ")
				// if B slope is negative, A should come before B as B will be after A slightly below the event p
				if B.slope < 0 {
					//debugLog.WriteString("B.slope < 0: A before B")
					return -1
				} else {
					//debugLog.WriteString("B.slope >= 0: B before A")
					return 1
				}
			}
			// if B is vertical *AND* B contains the event p *AND* A is diagonal *AND* A contains the event p:
			if B.isVertical && B.containsEvent && !A.isVertical && !A.isHorizontal && A.containsEvent {
				//debugLog.WriteString("B is vertical *AND* B contains the event p *AND* A is diagonal *AND* A contains the event p, ")
				// if A slope is negative, B should come before A as A will be after B slightly below the event p
				if A.slope < 0 {
					//debugLog.WriteString("A.slope < 0: B before A")
					return 1
				} else {
					//debugLog.WriteString("A.slope >= 0: A before B")
					return -1
				}
			}

			// Horizontal lines still come last if they contain p
			if A.isHorizontal && A.containsEvent && B.containsEvent && !B.isHorizontal {
				//debugLog.WriteString("A is horizontal and contains the event, B is not: B before A")
				return 1
			}
			// if B is horizontal and A is not, then B comes last
			if B.isHorizontal && B.containsEvent && A.containsEvent && !A.isHorizontal {
				//debugLog.WriteString("B is horizontal and contains the event, A is not: A before B")
				return -1
			}

			// If XAtY matches
			if numeric.FloatEquals(aX, bX, A.epsilon) {
				//debugLog.WriteString("A and B have equal X at event p, ")

				// Handle collinearity
				if A.slope == B.slope || A.isVertical && B.isVertical || A.isHorizontal && B.isHorizontal {
					//debugLog.WriteString("A and B are collinear, ")

					// if start Y points differ, order by start y
					if A.segment.Start().Y() != B.segment.Start().Y() {
						//debugLog.WriteString("order by start y: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(B.segment.Start().Y(), A.segment.Start().Y())))
						return cmp.Compare(B.segment.Start().Y(), A.segment.Start().Y()) // order by start y, highest first
					}

					// if start Y are equal, if start X differ, order by start x
					if A.segment.Start().X() != B.segment.Start().X() {
						//debugLog.WriteString("start y equal, so order by start x: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(A.segment.Start().X(), B.segment.Start().X())))
						return cmp.Compare(A.segment.Start().X(), B.segment.Start().X()) // order by start x, lowest first
					}

					// if start points are equal, if end Y differ, order by end y
					if A.segment.End().Y() != B.segment.End().Y() {
						//debugLog.WriteString("start points equal, so order by end y: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(B.segment.End().Y(), A.segment.End().Y())))
						return cmp.Compare(B.segment.End().Y(), A.segment.End().Y()) // / order by end y, highest first
					}

					// if start points are equal and end y is equal, if end x differ, order by end x
					if A.segment.End().X() != B.segment.End().X() {
						//debugLog.WriteString("start points equal, end y equal, so order by end x: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(A.segment.End().X(), B.segment.End().X())))
						return cmp.Compare(A.segment.End().X(), B.segment.End().X()) // order by end x, lowest first
					}
				}

				// one segment is vertical & one horizontal
				if A.isVertical && B.isHorizontal {
					if A.containsEvent == B.containsEvent {
						//debugLog.WriteString("A vertical, B horizontal, both contain event, horizontal last: ")
						//debugLog.WriteString(debugPrintOrder(1))
						return 1
					}
					//debugLog.WriteString("A vertical, B horizontal: ")
					//debugLog.WriteString(debugPrintOrder(cmp.Compare(p.X(), B.currX)))
					cmp.Compare(p.X(), B.currX)
				}
				if B.isVertical && A.isHorizontal {
					if B.containsEvent == A.containsEvent {
						//debugLog.WriteString("B vertical, A horizontal, both contain event, horizontal last")
						return -1
					}
					//debugLog.WriteString("B vertical, A horizontal: ")
					//debugLog.WriteString(debugPrintOrder(cmp.Compare(A.currX, p.X())))
					cmp.Compare(A.currX, p.X())
				}

				// order by slope: both slopes negative or both slopes positive
				if A.slope < 0 && B.slope < 0 || A.slope > 0 && B.slope > 0 {
					//debugLog.WriteString("order by slope, both slopes positive or negative: ")
					//debugLog.WriteString(debugPrintOrder(cmp.Compare(A.slope, B.slope)))
					return cmp.Compare(A.slope, B.slope)
				}

				// order by slope: opposing slopes
				if (A.slope < 0 && B.slope > 0) || (A.slope > 0 && B.slope < 0) {
					//debugLog.WriteString("order by slope, opposing slopes, ")

					// changed this from greater than with points swapped
					if numeric.FloatLessThan(p.X(), A.currX, A.epsilon) {
						//debugLog.WriteString("intersection is after event p: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(B.slope, A.slope)))
						return cmp.Compare(A.slope, B.slope)
					} else {
						//debugLog.WriteString("intersection is before event p: ")
						//debugLog.WriteString(debugPrintOrder(cmp.Compare(B.slope, A.slope)))
						return cmp.Compare(B.slope, A.slope)
					}
				}

				if A.isVertical && B.slope < 0 {
					//debugLog.WriteString("B slope negative, A is vertical ")
					if numeric.FloatLessThan(p.X(), A.currX, A.epsilon) {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(1))
						return 1
					} else {
						//debugLog.WriteString("and to the left of the event p: ")
						//debugLog.WriteString(debugPrintOrder(-1))
						return -1
					}
				}
				if B.isVertical && A.slope < 0 {
					//debugLog.WriteString("A slope negative, B is vertical ")
					if p.X() < B.currX {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(-1))
						return -1
					} else {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(1))
						return 1
					}
				}

				if A.isVertical && B.slope > 0 {
					//debugLog.WriteString("B slope positive, A is vertical ")
					if p.X() < A.currX {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(-1))
						return -1
					} else {
						//debugLog.WriteString("and to the left of the event p: ")
						//debugLog.WriteString(debugPrintOrder(1))
						return 1
					}
				}
				if B.isVertical && A.slope > 0 {
					//debugLog.WriteString("A slope positive, B is vertical ")
					if p.X() < B.currX {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(1))
						return 1
					} else {
						//debugLog.WriteString("and to the right of the event p: ")
						//debugLog.WriteString(debugPrintOrder(-1))
						return -1
					}
				}
			}

			//log.Printf("  - %t via default XAtY comparison", aX < bX)
			//debugLog.WriteString("default XAtY comparison: ")
			//debugLog.WriteString(debugPrintOrder(cmp.Compare(aX, bX)))
			return cmp.Compare(aX, bX)
		}
		panic(fmt.Errorf("unexpected comparison"))
	}
}

func newStatusStructureRBT(p point.Point[float64], opts ...options.GeometryOptionsFunc) *statusStructureRBT {
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)
	S := new(statusStructureRBT)
	S.sweepPoint = p
	S.structure = rbt.NewWith(statusStructureComparator(&S.sweepPoint, &S.epsilon))
	S.epsilon = geoOpts.Epsilon
	return S
}
