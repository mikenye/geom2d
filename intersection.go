package geom2d

import "log"

// intersects determines if two geometric objects intersect.
// This function handles combinations of Point, LineSegment, Rectangle, Circle, and PolyTree.
//
// Parameters:
//   - a: The first geometric object.
//   - b: The second geometric object.
//
// Returns:
//   - bool: True if the objects intersect, false otherwise.
//
// todo: use the Relationship* functions
func intersects[T SignedNumber](a, b any, opts ...Option) bool {

	switch a := a.(type) {
	case Point[T]:
		switch b := b.(type) {
		case Point[T]:
			return a.Eq(b, opts...)
		case LineSegment[T]:
			return b.ContainsPoint(a) // todo: implement epsilon?
		case Rectangle[T]:
			return intersectsPointRectangle(a, b) // todo: implement epsilon?
		case Circle[T]:
			return intersectsPointCircle(a, b) // todo: implement epsilon?
		case PolyTree[T]:
			// todo: implement
			log.Printf("Intersection between Point and PolyTree not implemented")
			return false
		}

	case LineSegment[T]:
		switch b := b.(type) {
		case Point[T]:
			return a.ContainsPoint(b) // todo: implement epsilon?
		case LineSegment[T]:
			return a.IntersectsLineSegment(b) // todo: implement epsilon?
		case Rectangle[T]:
			return intersectsLineSegmentRectangle(a, b) // todo: implement epsilon?
		case Circle[T]:
			return intersectsLineSegmentCircle(a, b)
		case PolyTree[T]:
			log.Printf("Intersection between LineSegment and PolyTree not implemented")
			return false
		}

	case Rectangle[T]:
		switch b := b.(type) {
		case Point[T]:
			return intersectsPointRectangle(b, a) // todo: implement epsilon?
		case LineSegment[T]:
			return intersectsLineSegmentRectangle(b, a) // todo: implement epsilon?
		case Rectangle[T]:
			log.Printf("Intersection between Rectangle and Rectangle not implemented")
			return false
		case Circle[T]:
			log.Printf("Intersection between Rectangle and Circle not implemented")
			return false
		case PolyTree[T]:
			log.Printf("Intersection between Rectangle and PolyTree not implemented")
			return false
		}

	case Circle[T]:
		switch b := b.(type) {
		case Point[T]:
			return intersectsPointCircle(b, a) // todo: implement epsilon?
		case LineSegment[T]:
			return intersectsLineSegmentCircle(b, a) // todo: implement epsilon?
		case Rectangle[T]:
			log.Printf("Intersection between Circle and Rectangle not implemented")
			return false
		case Circle[T]:
			log.Printf("Intersection between Circle and Circle not implemented")
			return false
		case PolyTree[T]:
			log.Printf("Intersection between Circle and PolyTree not implemented")
			return false
		}

		//case PolyTree[T]:
		//	switch _ = b.(type) {
		//	case Point[T]:
		//		log.Printf("Intersection between PolyTree and Point not implemented")
		//		return false
		//	case LineSegment[T]:
		//		log.Printf("Intersection between PolyTree and LineSegment not implemented")
		//		return false
		//	case Rectangle[T]:
		//		log.Printf("Intersection between PolyTree and Rectangle not implemented")
		//		return false
		//	case Circle[T]:
		//		log.Printf("Intersection between PolyTree and Circle not implemented")
		//		return false
		//	case PolyTree[T]:
		//		log.Printf("Intersection between PolyTree and PolyTree not implemented")
		//		return false
		//	}
	}

	log.Printf("Unknown geometric types: %T, %T", a, b)
	return false
}

func intersectsPointRectangle[T SignedNumber](p Point[T], r Rectangle[T]) bool {
	rel := r.RelationshipToPoint(p) // todo: implement epsilon?
	if rel == RelationshipPointRectanglePointOnEdge || rel == RelationshipPointRectanglePointOnVertex {
		return true
	}
	return false
}

func intersectsPointCircle[T SignedNumber](p Point[T], c Circle[T], opts ...Option) bool {
	rel := c.RelationshipToPoint(p, opts...)
	if rel == RelationshipPointCircleOnCircumference {
		return true
	}
	return false
}

func intersectsLineSegmentRectangle[T SignedNumber](l LineSegment[T], r Rectangle[T], opts ...Option) bool {
	rel := r.RelationshipToLineSegment(l)
	switch rel {
	case RelationshipLineSegmentRectangleEndTouchesEdgeExternally, RelationshipLineSegmentRectangleEndTouchesVertexExternally, RelationshipLineSegmentRectangleEndTouchesEdgeInternally, RelationshipLineSegmentRectangleEndTouchesVertexInternally, RelationshipLineSegmentRectangleEdgeCollinear, RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex, RelationshipLineSegmentRectangleIntersects, RelationshipLineSegmentRectangleEntersAndExits:
		return true
	default:
		return false
	}
}

func intersectsLineSegmentCircle[T SignedNumber](l LineSegment[T], c Circle[T], opts ...Option) bool {
	rel := c.RelationshipToLineSegment(l, opts...)
	switch rel {
	case RelationshipLineSegmentCircleIntersecting, RelationshipLineSegmentCircleTangentToCircle, RelationshipLineSegmentCircleEndOnCircumferenceOutside, RelationshipLineSegmentCircleEndOnCircumferenceInside, RelationshipLineSegmentCircleBothEndsOnCircumference:
		return true
	default:
		return false
	}
}
