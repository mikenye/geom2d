package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestNewLineSegment(t *testing.T) {

	tests := map[string]struct {
		x1, y1, x2, y2                 float64
		expectError                    bool
		expectedXUpper, expectedYUpper float64
		expectedXLower, expectedYLower float64
	}{
		"different Y": {
			x1:             1,
			y1:             1,
			x2:             2,
			y2:             2,
			expectError:    false,
			expectedXUpper: 2,
			expectedYUpper: 2,
			expectedXLower: 1,
			expectedYLower: 1,
		},
		"equal Y, different X": {
			x1:             1,
			y1:             2,
			x2:             2,
			y2:             2,
			expectError:    false,
			expectedXUpper: 2,
			expectedYUpper: 2,
			expectedXLower: 1,
			expectedYLower: 2,
		},
		"out of bounds x1,y1": {
			x1:          math.MaxFloat64,
			y1:          math.MaxFloat64,
			x2:          0,
			y2:          0,
			expectError: true,
		},
		"out of bounds x2,y2": {
			x1:          0,
			y1:          0,
			x2:          -math.MaxFloat64,
			y2:          -math.MaxFloat64,
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			// Crete Plane
			plane, err := New()
			require.NoError(t, err, "unexpected error creating Plane")

			// Basic line segment creation
			ls, err := plane.NewLineSegment(tc.x1, tc.y1, tc.x2, tc.expectedXUpper)
			if tc.expectError {
				require.Error(t, err, "expected error creating LineSegment")
			} else {
				require.NoError(t, err, "unexpected error creating LineSegment")
				require.NotNil(t, ls, "expected a valid line segment")

				upper, lower := ls.Endpoints()

				assert.InDelta(t, tc.expectedXUpper, upper.X(), DefaultPrecision, "upper X mismatch")
				assert.InDelta(t, tc.expectedYUpper, upper.Y(), DefaultPrecision, "upper Y mismatch")
				assert.InDelta(t, tc.expectedXLower, lower.X(), DefaultPrecision, "lower X mismatch")
				assert.InDelta(t, tc.expectedYLower, lower.Y(), DefaultPrecision, "lower Y mismatch")
			}

		})
	}
}

func TestNewPlane_Defaults(t *testing.T) {
	p, err := New()
	require.NoError(t, err)
	assert.Equal(t, DefaultPrecision, p.precision)
	assert.Equal(t, 1/DefaultPrecision, p.scaleFactor)
	assert.LessOrEqual(t, p.minBound, math.MaxFloat64)
	assert.GreaterOrEqual(t, p.maxBound, -math.MaxFloat64)
}

func TestNewPlane_WithPrecision(t *testing.T) {
	p, err := New(WithPrecision(1e-6))
	require.NoError(t, err)
	assert.Equal(t, 1e-6, p.precision)
	assert.Equal(t, 1/1e-6, p.scaleFactor)
	assert.LessOrEqual(t, p.minBound, math.MaxFloat64)
	assert.GreaterOrEqual(t, p.maxBound, -math.MaxFloat64)
}

func TestNewPlane_WithPrecision_TooSmall(t *testing.T) {
	p, err := New(WithPrecision(0))
	require.Nil(t, p)
	assert.ErrorIs(t, err, errPrecisionTooSmall)
}

func TestNewPlane_WithPrecision_TooLarge(t *testing.T) {
	p, err := New(WithPrecision(1))
	require.Nil(t, p)
	assert.ErrorIs(t, err, errPrecisionTooLarge)
}

func TestNewPoint(t *testing.T) {
	plane, err := New()
	require.NoError(t, err, "unexpected error when creating Plane")

	// Basic point creation
	pt, err := plane.NewPoint(1.23456789, -4.56789)
	require.NoError(t, err, "Expected no error when creating point")
	require.NotNil(t, pt, "Expected a valid point")
	assert.InDelta(t, 1.23456789, pt.X(), DefaultPrecision, "X coordinate should match")
	assert.InDelta(t, -4.56789, pt.Y(), DefaultPrecision, "Y coordinate should match")

	// Duplicate point should return the same instance
	samePt, err := plane.NewPoint(1.23456789, -4.56789)
	assert.NoError(t, err)
	assert.Equal(t, pt, samePt, "Should return the same instance for identical coordinates")

	// Bounds checking
	_, err = plane.NewPoint(math.MaxFloat64, math.MaxFloat64)
	assert.Error(t, err, "Should return an error for out-of-bounds coordinates")
}
