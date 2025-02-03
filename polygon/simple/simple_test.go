package simple

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsWellFormed(t *testing.T) {
	tests := []struct {
		name      string
		points    []point.Point[int]
		expected  bool
		errSubstr string // Substring expected in the error message
	}{
		{
			name: "Valid triangle",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(10, 0),
				point.New(5, 5),
			},
			expected: true,
		},
		{
			name: "Too few points",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(10, 0),
			},
			expected:  false,
			errSubstr: "at least 3 points",
		},
		{
			name: "Zero area (collinear points)",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(5, 0),
				point.New(10, 0),
			},
			expected:  false,
			errSubstr: "zero area",
		},
		{
			name: "Self-intersecting polygon",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(10, 10),
				point.New(10, 0),
				point.New(0, 2),
			},
			expected:  false,
			errSubstr: "self-intersecting",
		},
		{
			name: "Valid large polygon",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(10, 0),
				point.New(10, 10),
				point.New(0, 10),
			},
			expected: true,
		},
		{
			name: "Polygon with duplicate points",
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(5, 0),
				point.New(5, 5),
				point.New(0, 5),
				point.New(0, 0), // Duplicate
			},
			expected: true, // duplicate point will be ignored
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := IsWellFormed(tc.points, options.WithEpsilon(1e-8))

			assert.Equal(t, tc.expected, result)

			if tc.errSubstr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
