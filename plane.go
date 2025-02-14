package geom2d

import (
	"errors"
	"fmt"
	"math"
)

const DefaultPrecision = float64(1e-10)

var (
	errPrecisionTooSmall = errors.New("precision must be greater than zero")
	errPrecisionTooLarge = errors.New("precision must be less than one")
)

type Plane struct {
	precision          float64
	scaleFactor        float64 // 1/precision
	minBound, maxBound float64 // bounds
}

type Option func(*Plane)

// WithPrecision sets a custom precision to override the default of DefaultPrecision
func WithPrecision(precision float64) Option {
	return func(p *Plane) {
		p.precision = precision
		p.scaleFactor = 1 / precision
	}
}

func New(opts ...Option) (*Plane, error) {
	p := new(Plane)

	// defaults
	p.precision = DefaultPrecision

	// Apply user options
	for _, opt := range opts {
		opt(p)
	}

	// Sanity checks after user options
	if p.precision <= 0 {
		return nil, errPrecisionTooSmall
	}
	if p.precision >= 1 {
		return nil, errPrecisionTooLarge
	}

	// calculate dynamic fields
	p.scaleFactor = 1 / p.precision
	p.maxBound = float64(math.MaxInt64) / p.scaleFactor
	p.minBound = -p.maxBound

	return p, nil
}

func (p *Plane) NewPoint(x, y float64) (*Point, error) {

	// Bounds check
	if x < p.minBound || x > p.maxBound || y < p.minBound || y > p.maxBound {
		return nil, errPointOutOfBounds
	}

	// Create point
	return &Point{
		scaledX: int64(math.Round(x * p.scaleFactor)), // Scale the X coordinate to integer
		scaledY: int64(math.Round(y * p.scaleFactor)), // Scale the Y coordinate to integer
		plane:   p,
	}, nil
}

func (p *Plane) NewLineSegment(x1, y1, x2, y2 float64) (*LineSegment, error) {

	// Ensure correct ordering: pointUpper should have a higher y, or if equal, a lower x
	upperX, upperY, lowerX, lowerY := x1, y1, x2, y2
	if y2 > y1 || (y1 == y2 && x2 > x1) {
		upperX, upperY, lowerX, lowerY = x2, y2, x1, y1
	}

	// Create upper point
	pointUpper, err := p.NewPoint(upperX, upperY)
	if err != nil {
		return nil, fmt.Errorf("invalid line segment: %w", err)
	}

	// Create lower point
	pointLower, err := p.NewPoint(lowerX, lowerY)
	if err != nil {
		return nil, fmt.Errorf("invalid line segment: %w", err)
	}

	// Return the created line segment
	return &LineSegment{pointUpper: pointUpper, pointLower: pointLower}, nil
}
