package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIssue7(t *testing.T) {
	// https://github.com/mikenye/geom2d/issues/7

	t.Run("issue as-per github", func(t *testing.T) {
		holeInPolyTree, err := NewPolyTree[int]([]Point[int]{
			NewPoint(299, 191),
			NewPoint(329, 195),
			NewPoint(325, 210),
			NewPoint(298, 211),
		}, PTHole)
		require.NoError(t, err, "unexpected error creating holeInPolyTree")

		polyTree, err := NewPolyTree[int]([]Point[int]{
			NewPoint(333, 218),
			NewPoint(345, 195),
			NewPoint(324, 181),
			NewPoint(341, 164),
			NewPoint(307, 169),
			NewPoint(270, 163),
			NewPoint(254, 180),
			NewPoint(263, 193),
			NewPoint(253, 210),
			NewPoint(290, 181),
			NewPoint(288, 218),
		}, PTSolid, WithChildren(holeInPolyTree))
		require.NoError(t, err, "unexpected error creating polyTree")

		// test every y value
		x := 250
		for y := 260; y >= 222; y-- {
			point := NewPoint(x, y)
			rel := point.RelationshipToPolyTree(polyTree, WithEpsilon(1e-9))
			require.Equal(t, RelationshipDisjoint, rel[polyTree])
			require.Equal(t, RelationshipDisjoint, rel[holeInPolyTree])
		}
	})

	t.Run("", func(t *testing.T) {
		pt, err := NewPolyTree([]Point[int]{
			NewPoint(0, 0),
			NewPoint(3, 0),
			NewPoint(5, 3),
			NewPoint(7, 0),
			NewPoint(10, 0),
			NewPoint(10, 4),
			NewPoint(13, 4),
			NewPoint(13, 5),
			NewPoint(16, 5),
			NewPoint(16, 4),
			NewPoint(19, 4),
			NewPoint(19, 3),
			NewPoint(22, 3),
			NewPoint(22, 4),
			NewPoint(25, 4),
			NewPoint(26, 2),
			NewPoint(27, 4),
			NewPoint(27, 7),
			NewPoint(25, 8),
			NewPoint(27, 9),
			NewPoint(28, 10),
			NewPoint(27, 11),
			NewPoint(25, 12),
			NewPoint(27, 13),
			NewPoint(27, 16),
			NewPoint(26, 18),
			NewPoint(25, 16),
			NewPoint(22, 16),
			NewPoint(22, 17),
			NewPoint(19, 17),
			NewPoint(19, 16),
			NewPoint(16, 16),
			NewPoint(16, 15),
			NewPoint(13, 15),
			NewPoint(13, 16),
			NewPoint(10, 16),
			NewPoint(10, 20),
			NewPoint(7, 20),
			NewPoint(5, 17),
			NewPoint(3, 20),
			NewPoint(0, 20),
			NewPoint(-1, 19),
			NewPoint(0, 18),
			NewPoint(0, 2),
			NewPoint(-1, 1),
		}, PTSolid)

		require.NoError(t, err, "unexpected error creating pt")

		// test every y value
		x := -2
		for y := 21; y >= -1; y-- {
			point := NewPoint(x, y)
			rel := point.RelationshipToPolyTree(pt, WithEpsilon(1e-9))
			require.Equal(t, RelationshipDisjoint, rel[pt])
		}
	})
}

func TestIssue15(t *testing.T) {

	// mimmick the issue conditions for issue #15 (https://github.com/mikenye/geom2d/issues/15)

	ptAHole, err := NewPolyTree([]Point[int]{
		NewPoint(150, 150),
		NewPoint(200, 150),
		NewPoint(200, 100),
		NewPoint(150, 100),
	}, PTHole)
	require.NoError(t, err, "unexpected error creating ptAHole")

	ptA, err := NewPolyTree([]Point[int]{
		NewPoint(125, 75),
		NewPoint(225, 75),
		NewPoint(225, 175),
		NewPoint(125, 175),
	}, PTSolid, WithChildren(ptAHole))
	require.NoError(t, err, "unexpected error creating ptA")

	ptBHole, err := NewPolyTree([]Point[int]{
		NewPoint(35, 135),
		NewPoint(35, 185),
		NewPoint(85, 185),
		NewPoint(85, 135),
	}, PTHole)
	require.NoError(t, err, "unexpected error creating ptBHole")

	ptB, err := NewPolyTree([]Point[int]{
		NewPoint(10, 110),
		NewPoint(110, 110),
		NewPoint(110, 210),
		NewPoint(10, 210),
	}, PTSolid, WithChildren(ptBHole))
	require.NoError(t, err, "unexpected error creating ptB")

	// mimmick user dragging ptB across ptA
	for i := 0; i < 230; i++ {

		// perform boolean operation in goroutine
		completed := false
		go func() {
			_, err := ptA.BooleanOperation(
				ptB,
				BooleanSubtraction,
			)
			require.NoError(t, err, "unexpected error creating ptResult")
			completed = true
		}()

		// Require that goroutine finishes within 3 seconds.
		// If not, then the function must be stuck in an infinite loop.
		require.EventuallyWithTf(t, func(c *assert.CollectT) {
			require.True(c, completed, "expected 'completed' to be true")
		}, 3*time.Second, 100*time.Millisecond, "external state has not changed to 'true'; still false")

		// mimmick dragging by moving ptB one pixel to the right
		ptB = ptB.Translate(NewPoint(1, 0))
	}
}
