package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeClockwise(t *testing.T) {
	tests := map[string]struct {
		points []point.Point[int]
	}{
		"already clockwise": {
			points: []point.Point[int]{point.New(0, 0), point.New(2, 3), point.New(4, 0)},
		},
		"counterclockwise input": {
			points: []point.Point[int]{point.New(0, 0), point.New(4, 0), point.New(2, 3)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			MakeClockwise(tc.points...)
			assert.Equal(t, types.PointsClockwise, Orientation(tc.points...))
		})
	}
}

func TestMakeCounterClockwise(t *testing.T) {
	tests := map[string]struct {
		points []point.Point[int]
	}{
		"already counterclockwise": {
			points: []point.Point[int]{point.New(0, 0), point.New(4, 0), point.New(2, 3)},
		},
		"clockwise input": {
			points: []point.Point[int]{point.New(0, 0), point.New(2, 3), point.New(4, 0)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			MakeCounterClockwise(tc.points...)
			assert.Equal(t, types.PointsCounterClockwise, Orientation(tc.points...))
		})
	}
}

func TestOrientation(t *testing.T) {
	tests := map[string]struct {
		p0, p1, p2 any
		expected   types.PointOrientation
	}{
		"int: (0,0), (10,10), (10,0)": {
			p0:       point.New[int](0, 0),
			p1:       point.New[int](10, 10),
			p2:       point.New[int](10, 0),
			expected: types.PointsClockwise,
		},
		"int: (0,0), (10,0), (10,10)": {
			p0:       point.New[int](0, 0),
			p1:       point.New[int](10, 0),
			p2:       point.New[int](10, 10),
			expected: types.PointsCounterClockwise,
		},
		"int: (0,0), (10,0), (20,0)": {
			p0:       point.New[int](0, 0),
			p1:       point.New[int](10, 0),
			p2:       point.New[int](20, 0),
			expected: types.PointsCollinear,
		},
		"float64: (0,0), (10,10), (10,0)": {
			p0:       point.New[float64](0, 0),
			p1:       point.New[float64](10, 10),
			p2:       point.New[float64](10, 0),
			expected: types.PointsClockwise,
		},
		"float64: (0,0), (10,0), (10,10)": {
			p0:       point.New[float64](0, 0),
			p1:       point.New[float64](10, 0),
			p2:       point.New[float64](10, 10),
			expected: types.PointsCounterClockwise,
		},
		"float64: (0,0), (10,0), (20,0)": {
			p0:       point.New[float64](0, 0),
			p1:       point.New[float64](10, 0),
			p2:       point.New[float64](20, 0),
			expected: types.PointsCollinear,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p0 := tt.p0.(type) {
			case point.Point[int]:
				p1 := tt.p1.(point.Point[int])
				p2 := tt.p2.(point.Point[int])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)

			case point.Point[float64]:
				p1 := tt.p1.(point.Point[float64])
				p2 := tt.p2.(point.Point[float64])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
