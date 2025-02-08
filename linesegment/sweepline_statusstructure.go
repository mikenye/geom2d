package linesegment

import (
	"fmt"
	"github.com/google/btree"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"log"
	"math"
	"strings"
)

// todo: proper doc comments
// segment contains the possibly modified line segment that the status structure is sorted via
// original is the slice of the line segments that are part of the segment
type sItem struct {
	segment   LineSegment[float64]
	originals []LineSegment[float64]
}

func (si *sItem) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s (originally:", si.segment))
	for _, seg := range si.originals {
		builder.WriteString(fmt.Sprintf(" %s", seg))
	}
	builder.WriteString(")")
	return builder.String()
}

func debugStatusStructure(S *btree.BTreeG[sItem]) {
	Scopy := S.Clone()
	for Scopy.Len() > 0 {
		item, _ := Scopy.DeleteMin()
		log.Printf("  - %s\n", item.String())
	}
}

func updateStatusStructure(S *btree.BTreeG[sItem], p qItem, opts ...options.GeometryOptionsFunc) *btree.BTreeG[sItem] {
	// if the status structure doesn't exist, create it
	if S == nil {
		return btree.NewG[sItem](2, segmentSortLessHigherOrder(p.point, opts...))
	}
	// otherwise, re-create
	newS := btree.NewG[sItem](2, segmentSortLessHigherOrder(p.point, opts...))
	S.Ascend(func(item sItem) bool {
		newS.ReplaceOrInsert(item)
		return true
	})
	return newS
}

func segmentSortLessHigherOrder(p point.Point[float64], opts ...options.GeometryOptionsFunc) btree.LessFunc[sItem] {
	return func(a, b sItem) bool {
		geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

		aX := a.segment.XAtY(p.Y())
		aSlope := a.segment.Slope()
		aIsHorizontal := math.IsNaN(aX)
		aIsVertical := math.IsNaN(aSlope)
		aContainsP := a.segment.ContainsPoint(p, opts...)

		bX := b.segment.XAtY(p.Y())
		bSlope := b.segment.Slope()
		bIsHorizontal := math.IsNaN(bX)
		bIsVertical := math.IsNaN(bSlope)
		bContainsP := b.segment.ContainsPoint(p, opts...)

		// for horizontal lines, artificially truncate start position to point,
		// since we don't care about anything to the left, as that is considered above the sweep line
		if math.IsNaN(aX) {
			aX = p.X()
		}
		if math.IsNaN(bX) {
			bX = p.X()
		}

		//log.Printf(
		//	"is %s (x=%f, s=%f) to the left of %s (x=%f, s=%f) at %s: ",
		//	a.segment.String(),
		//	aX,
		//	aSlope,
		//	b.segment.String(),
		//	bX,
		//	bSlope,
		//	p.String(),
		//)

		// Vertical segment ordering logic: Handle cases where a vertical segment intersects a diagonal one.
		if aIsVertical && aContainsP && numeric.FloatEquals(aX, p.X(), geoOpts.Epsilon) && !bIsVertical && !bIsHorizontal && bContainsP {
			//log.Println(bSlope < 0, "via slope & intersection with vertical (a)")
			return bSlope < 0
		}
		if bIsVertical && bContainsP && numeric.FloatEquals(bX, p.X(), geoOpts.Epsilon) && !aIsVertical && !bIsHorizontal && aContainsP {
			//log.Println(aSlope > 0, "via slope & intersection with vertical (b)")
			return aSlope > 0
		}

		// Horizontal lines still come last if they contain p
		if aIsHorizontal && b.segment.ContainsPoint(p, opts...) && !bIsHorizontal {
			//log.Println("false via horizontal handling (a is horizontal, b contains p)")
			return false
		}
		if bIsHorizontal && a.segment.ContainsPoint(p, opts...) && !aIsHorizontal {
			//log.Println("true via horizontal handling (b is horizontal, a contains p)")
			return true
		}

		// if both horizontal, order by end x
		if aIsHorizontal && bIsHorizontal {
			//log.Println("true via horizontal handling (both horizontal)")
			return a.segment.End().X() > b.segment.End().X()
		}

		// If XAtY matches
		if numeric.FloatEquals(aX, bX, geoOpts.Epsilon) {

			// if slopes are equal, the line segments are collinear
			if aSlope == bSlope {
				if a.segment.Start().Y() != b.segment.Start().Y() {
					return a.segment.Start().Y() > b.segment.Start().Y() // order by start y
				}
				if a.segment.Start().X() != b.segment.Start().X() {
					return a.segment.Start().X() < b.segment.Start().X() // order by start x
				}
				if a.segment.End().Y() != b.segment.End().Y() {
					return a.segment.End().Y() > b.segment.End().Y() // order by end y
				}
				if a.segment.End().X() != b.segment.End().X() {
					return a.segment.End().X() < b.segment.End().X() // order by end y
				}
			}

			// order by slope
			if (aSlope < 0 && bSlope > 0) || (aSlope > 0 && bSlope < 0) {
				//log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes opposite")
				return aSlope > bSlope
			} else if aSlope < 0 && bSlope < 0 {
				//log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both negative")
				return aSlope < bSlope
			} else {
				//log.Println(aSlope > bSlope, "via slope as XAtY was equal & slopes both positive")
				return aSlope < bSlope
			}
		}

		//log.Println(aX < bX, "via default XAtY comparison")
		return aX < bX // Default XAtY comparison
	}
}
