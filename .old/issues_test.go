package _old

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIssue7(t *testing.T) {

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
