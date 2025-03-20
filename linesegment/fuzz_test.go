package linesegment

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func FuzzFindIntersections_2segments(f *testing.F) {

	f.Add(0.0, 0.0, 10.0, 10.0, 10.0, 0.0, 20.0, 10.0)
	f.Add(0.0, 0.0, 10.0, 10.0, 10.0, 10.0, 20.0, 0.0)
	f.Add(0.0, 10.0, 10.0, 0.0, 10.0, 0.0, 20.0, 10.0)
	f.Add(0.0, 10.0, 10.0, 20.0, 0.0, 10.0, 10.0, 0.0)
	f.Add(0.0, 20.0, 10.0, 10.0, 10.0, 10.0, 0.0, 0.0)
	f.Add(0.0, 0.0, 10.0, 10.0, 10.0, 0.0, 0.0, 10.0)
	f.Add(10.0, 20.0, 10.0, 0.0, 0.0, 20.0, 20.0, 0.0)
	f.Add(10.0, 20.0, 10.0, 0.0, 20.0, 20.0, 0.0, 0.0)
	f.Add(0.0, 10.0, 20.0, 10.0, 20.0, 20.0, 0.0, 0.0)
	f.Add(0.0, 10.0, 20.0, 10.0, 0.0, 20.0, 20.0, 0.0)
	f.Add(10.0, 20.0, 10.0, 0.0, 0.0, 10.0, 20.0, 10.0)
	f.Add(20.0, 20.0, 0.0, 0.0, 0.0, 20.0, 20.0, 0.0)

	f.Fuzz(func(t *testing.T, ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 float64) {
		segA := New(ax1, ay1, ax2, ay2)
		t.Logf("Input segment A: %s", segA)
		segB := New(bx1, by1, bx2, by2)
		t.Logf("Input segment B: %s", segB)
		input := []LineSegment{segA, segB}

		resultsSweepLine := FindIntersections(input)
		t.Logf("sweep line:\n%s", resultsSweepLine)

		resultsBruteForce := FindIntersectionsBruteForce(input)
		t.Logf("brute force:\n%s", resultsBruteForce)

		require.Equal(t, resultsSweepLine.Size(), resultsBruteForce.Size(), "results size mismatch")

		bfNode := resultsBruteForce.Min(resultsBruteForce.Root())
		for slNode := resultsSweepLine.Min(resultsSweepLine.Root()); !resultsSweepLine.IsNil(slNode); slNode = resultsSweepLine.Successor(slNode) {
			assert.True(t, resultsSweepLine.Key(slNode).Eq(resultsBruteForce.Key(bfNode)), "point mismatch")
			assert.Equal(t, resultsSweepLine.Value(slNode), resultsBruteForce.Value(bfNode), "line segment mismatch")
			bfNode = resultsBruteForce.Successor(bfNode)
		}

	})
}
