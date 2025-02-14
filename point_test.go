package geom2d

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint_Eq(t *testing.T) {
	plane1, err := New()
	require.NoError(t, err, "unexpected error creating plane1")

	plane2, err := New()
	require.NoError(t, err, "unexpected error creating plane2")

	plane3, err := New(WithPrecision(1e-9))
	require.NoError(t, err, "unexpected error creating plane3")

	// Same point, same plane, should be equal
	p1, err := plane1.NewPoint(1.23456789, -4.56789)
	require.NoError(t, err, "unexpected error creating p1")

	p2, err := plane1.NewPoint(1.23456789, -4.56789)
	require.NoError(t, err, "unexpected error creating p2")
	assert.True(t, p1.Eq(p2), "Points with same coordinates on the same Plane should be equal")

	// Different coordinates, should NOT be equal
	p3, err := plane1.NewPoint(2.0, -4.56789)
	require.NoError(t, err, "unexpected error creating p3")
	assert.False(t, p1.Eq(p3), "Points with different coordinates should not be equal")

	// Same coordinates, different plane, same precision, should be equal
	p4, err := plane2.NewPoint(1.23456789, -4.56789)
	require.NoError(t, err, "unexpected error creating p4")
	assert.True(t, p1.Eq(p4), "Points on different planes with same precision should be equal if coordinates match")

	// Same coordinates, different plane, different precision, should NOT be equal
	p5, err := plane3.NewPoint(1.23456789, -4.56789)
	require.NoError(t, err, "unexpected error creating p5")
	assert.False(t, p1.Eq(p5), "Points on planes with different precision should not be equal even if coordinates match")
}
