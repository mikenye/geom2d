package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestLineSegment_Slope(t *testing.T) {
	plane, err := New()
	require.NoError(t, err, "unexpected error creating plane")

	tests := map[string]struct {
		x1, y1, x2, y2 float64
		expectedSlope  float64
	}{
		"positive slope": {0, 0, 10, 10, 1.0},        // y = x
		"negative slope": {0, 10, 10, 0, -1.0},       // y = -x + 10
		"zero slope":     {0, 5, 10, 5, 0.0},         // horizontal
		"vertical slope": {5, 0, 5, 10, math.Inf(1)}, // vertical
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ls, err := plane.NewLineSegment(tc.x1, tc.y1, tc.x2, tc.y2)
			require.NoError(t, err, "unexpected error creating line segment")
			assert.Equal(t, tc.expectedSlope, ls.Slope(), "Incorrect slope for %s", name)
		})
	}
}

func TestLineSegment_SlopeCmp(t *testing.T) {
	plane, err := New()
	require.NoError(t, err, "unexpected error creating plane")

	tests := map[string]struct {
		aX1, aY1, aX2, aY2 float64
		bX1, bY1, bX2, bY2 float64
		expectedCmp        int
	}{
		"equal slopes":             {0, 0, 10, 10, 5, 5, 15, 15, 0},  // y = x vs y = x
		"equal vertical":           {0, 0, 0, 10, 5, 5, 5, 15, 0},    // y = x vs y = x
		"positive vs negative":     {0, 0, 10, 10, 0, 10, 10, 0, 1},  // y = x vs y = -x
		"negative vs positive":     {0, 10, 10, 0, 0, 0, 10, 10, -1}, // y = -x vs y = x
		"horizontal vs diagonal":   {0, 5, 10, 5, 0, 0, 10, 10, -1},  // flat vs y=x
		"diagonal vs horizontal":   {0, 0, 10, 10, 0, 5, 10, 5, 1},   // y=x vs flat
		"vertical vs diagonal":     {5, 0, 5, 10, 0, 0, 10, 10, -1},  // vertical is steepest
		"diagonal vs vertical":     {0, 0, 10, 10, 5, 0, 5, 10, 1},   // vertical is steepest
		"steeper negative slope":   {0, 10, 10, 0, 0, 10, 20, 5, -1}, // -1 vs -0.5
		"shallower negative slope": {0, 10, 20, 5, 0, 10, 10, 0, 1},  // -0.5 vs -1
		"steeper positive slope":   {0, 0, 10, 10, 0, 0, 10, 5, -1},  // 1 vs 0.5
		"shallower positive slope": {0, 0, 10, 5, 0, 0, 10, 10, 1},   // 0.5 vs 1
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			segA, err := plane.NewLineSegment(tc.aX1, tc.aY1, tc.aX2, tc.aY2)
			require.NoError(t, err, "unexpected error creating first line segment")

			segB, err := plane.NewLineSegment(tc.bX1, tc.bY1, tc.bX2, tc.bY2)
			require.NoError(t, err, "unexpected error creating second line segment")

			assert.Equal(t, tc.expectedCmp, segA.SlopeCmp(segB), "Incorrect slope comparison for %s", name)
		})
	}
}

func TestLineSegment_XAtY(t *testing.T) {
	plane, err := New()
	require.NoError(t, err, "unexpected error creating plane")

	// Test a diagonal segment (y = x)
	ls, err := plane.NewLineSegment(0, 0, 10, 10)
	require.NoError(t, err, "unexpected error creating line segment")

	assert.Equal(t, 5.0, ls.XAtY(5), "Expected XAtY(5) to return 5")
	assert.Equal(t, 10.0, ls.XAtY(10), "Expected XAtY(10) to return 10")
	assert.Equal(t, 0.0, ls.XAtY(0), "Expected XAtY(0) to return 0")

	// Test a vertical segment (X should be constant)
	lsVertical, err := plane.NewLineSegment(5, 0, 5, 10)
	require.NoError(t, err, "unexpected error creating line segment")

	assert.Equal(t, 5.0, lsVertical.XAtY(0), "Expected vertical segment to return 5 at all Y")
	assert.Equal(t, 5.0, lsVertical.XAtY(5), "Expected vertical segment to return 5 at all Y")
	assert.Equal(t, 5.0, lsVertical.XAtY(10), "Expected vertical segment to return 5 at all Y")

	// Test out-of-bounds Y (should return NaN)
	assert.True(t, math.IsNaN(ls.XAtY(-1)), "Expected XAtY(-1) to return NaN for out-of-bounds Y")
	assert.True(t, math.IsNaN(ls.XAtY(11)), "Expected XAtY(11) to return NaN for out-of-bounds Y")
}

func TestLineSegment_YAtX(t *testing.T) {
	plane, err := New()
	require.NoError(t, err, "unexpected error creating plane")

	// Test a diagonal segment (y = x)
	ls, err := plane.NewLineSegment(0, 0, 10, 10)
	require.NoError(t, err, "unexpected error creating line segment")

	assert.Equal(t, 5.0, ls.YAtX(5), "Expected YAtX(5) to return 5")
	assert.Equal(t, 10.0, ls.YAtX(10), "Expected YAtX(10) to return 10")
	assert.Equal(t, 0.0, ls.YAtX(0), "Expected YAtX(0) to return 0")

	// Test a horizontal segment (Y should be constant)
	lsHorizontal, err := plane.NewLineSegment(0, 5, 10, 5)
	require.NoError(t, err, "unexpected error creating line segment")

	assert.Equal(t, 5.0, lsHorizontal.YAtX(0), "Expected horizontal segment to return 5 at all X")
	assert.Equal(t, 5.0, lsHorizontal.YAtX(5), "Expected horizontal segment to return 5 at all X")
	assert.Equal(t, 5.0, lsHorizontal.YAtX(10), "Expected horizontal segment to return 5 at all X")

	// Test out-of-bounds X (should return NaN)
	assert.True(t, math.IsNaN(ls.YAtX(-1)), "Expected YAtX(-1) to return NaN for out-of-bounds X")
	assert.True(t, math.IsNaN(ls.YAtX(11)), "Expected YAtX(11) to return NaN for out-of-bounds X")
}
