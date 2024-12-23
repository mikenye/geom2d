package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mikenye/geom2d"
	"golang.org/x/image/colornames"
	"log"
)

// App represents the state of the demo application.
type App struct {
	mouseX, mouseY   int                     // Current mouse position.
	mouseLineSegment geom2d.LineSegment[int] // Current mouse position as a geom2d LineSegment.

	// Single point and its relationship to the mouse point.
	point             geom2d.Point[int]
	pointRelationship geom2d.Relationship

	// Line segment and its relationship to the mouse point.
	lineSegment             geom2d.LineSegment[int]
	lineSegmentRelationship geom2d.Relationship

	// Rectangle and its relationship to the mouse point.
	rectangle             geom2d.Rectangle[int]
	rectangleRelationship geom2d.Relationship

	// Circle and its relationship to the mouse point.
	circle             geom2d.Circle[int]
	circleRelationship geom2d.Relationship

	// PolyTree and its relationship to the mouse point.
	polyTree             *geom2d.PolyTree[int]
	polyTreeRelationship map[*geom2d.PolyTree[int]]geom2d.Relationship
}

// Draw renders the current state of the application to the screen.
func (app *App) Draw(screen *ebiten.Image) {
	// Render the single point.
	ebitenutil.DebugPrintAt(screen, app.point.String(), 5, 5)              // Display point details.
	ebitenutil.DebugPrintAt(screen, app.pointRelationship.String(), 5, 20) // Display relationship details.
	screen.Set(app.point.X(), app.point.Y(), colornames.Lightgreen)

	// Separator for visual clarity.
	vector.StrokeLine(screen, 0, 40, 800, 40, 1, colornames.White, false)

	// Render the line segment.
	ebitenutil.DebugPrintAt(screen, app.lineSegment.String(), 5, 45)
	ebitenutil.DebugPrintAt(screen, app.lineSegmentRelationship.String(), 5, 60)
	for p := range app.lineSegment.Bresenham {
		screen.Set(p.X(), p.Y(), colornames.Cyan) // Render each pixel of the line segment.
	}

	// Separator for visual clarity.
	vector.StrokeLine(screen, 0, 80, 800, 80, 1, colornames.White, false)

	// Render the rectangle by iterating through its edges.
	ebitenutil.DebugPrintAt(screen, app.rectangle.String(), 5, 85)
	ebitenutil.DebugPrintAt(screen, app.rectangleRelationship.String(), 5, 100)
	for _, edge := range app.rectangle.Edges() {
		for p := range edge.Bresenham {
			screen.Set(p.X(), p.Y(), colornames.Yellow) // Render each edge pixel.
		}
	}

	// Separator for visual clarity.
	vector.StrokeLine(screen, 0, 120, 800, 120, 1, colornames.White, false)

	// Render the circle.
	ebitenutil.DebugPrintAt(screen, app.circle.String(), 5, 125)
	ebitenutil.DebugPrintAt(screen, app.circleRelationship.String(), 5, 140)
	for p := range app.circle.Bresenham {
		screen.Set(p.X(), p.Y(), colornames.Pink) // Render the circle.
	}

	// Separator for visual clarity.
	vector.StrokeLine(screen, 0, 160, 800, 160, 1, colornames.White, false)

	// Render the PolyTree (polygon with a hole).
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Polygon: %s", app.polyTreeRelationship[app.polyTree].String()), 5, 165)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hole: %s", app.polyTreeRelationship[app.polyTree.Children()[0]].String()), 5, 180)
	for poly := range app.polyTree.Nodes {
		for _, edge := range poly.Edges() {
			for p := range edge.Bresenham {
				screen.Set(p.X(), p.Y(), colornames.Orange) // Render each edge of the polygon.
			}
		}
	}

	// Separator for visual clarity.
	vector.StrokeLine(screen, 0, 224, 800, 224, 1, colornames.White, false)

	// Render the mouse point as a red pixel.
	for p := range app.mouseLineSegment.Bresenham {
		screen.Set(p.X(), p.Y(), colornames.Red) // Render the circle.
	}

	// Display mouse position, FPS, and TPS.
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: %s", app.mouseLineSegment.Start().String()), 5, 224)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 200, 224)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()), 280, 224)
}

// Layout specifies the screen size for the application.
func (app *App) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 350, 250 // Fixed screen size.
}

// Update computes the relationships between the mouse point and other geometries.
func (app *App) Update() error {
	// Update mouse position and convert to a geom2d LineSegment.
	app.mouseX, app.mouseY = ebiten.CursorPosition()
	app.mouseLineSegment = geom2d.NewLineSegment[int](geom2d.NewPoint(app.mouseX, app.mouseY), geom2d.NewPoint(app.mouseX+6, app.mouseY+6))

	// Update relationships for each geometric object.
	app.pointRelationship = app.mouseLineSegment.RelationshipToPoint(app.point, geom2d.WithEpsilon(1))
	app.lineSegmentRelationship = app.mouseLineSegment.RelationshipToLineSegment(app.lineSegment, geom2d.WithEpsilon(1e-9))
	app.rectangleRelationship = app.mouseLineSegment.RelationshipToRectangle(app.rectangle, geom2d.WithEpsilon(1e-9))
	app.circleRelationship = app.mouseLineSegment.RelationshipToCircle(app.circle, geom2d.WithEpsilon(1))
	app.polyTreeRelationship = app.mouseLineSegment.RelationshipToPolyTree(app.polyTree, geom2d.WithEpsilon(1))

	return nil
}

// initGeometry initializes all geometric objects and handles potential errors.
func (app *App) initGeometry() error {
	// Initialize the mouse line segment.
	app.mouseLineSegment = geom2d.NewLineSegment[int](geom2d.NewPoint(0, 0), geom2d.NewPoint(6, 6))

	// Initialize the single point.
	app.point = geom2d.NewPoint[int](240, 20)

	// Initialize the line segment.
	app.lineSegment = geom2d.NewLineSegment[int](geom2d.NewPoint[int](175, 68), geom2d.NewPoint[int](310, 68))

	// Initialize the rectangle.
	app.rectangle = geom2d.NewRectangle[int]([]geom2d.Point[int]{
		geom2d.NewPoint[int](175, 115),
		geom2d.NewPoint[int](310, 115),
		geom2d.NewPoint[int](310, 100),
		geom2d.NewPoint[int](175, 100),
	})

	// Initialize the circle.
	app.circle = geom2d.NewCircle[int](geom2d.NewPoint[int](280, 139), 16)

	// Initialize the PolyTree.
	holeInPolyTree, err := geom2d.NewPolyTree[int]([]geom2d.Point[int]{
		geom2d.NewPoint(299, 191),
		geom2d.NewPoint(329, 195),
		geom2d.NewPoint(325, 210),
		geom2d.NewPoint(298, 211),
	}, geom2d.PTHole)
	if err != nil {
		return err
	}
	app.polyTree, err = geom2d.NewPolyTree[int]([]geom2d.Point[int]{
		geom2d.NewPoint(333, 218),
		geom2d.NewPoint(345, 195),
		geom2d.NewPoint(324, 181),
		geom2d.NewPoint(341, 164),
		geom2d.NewPoint(307, 169),
		geom2d.NewPoint(270, 163),
		geom2d.NewPoint(254, 180),
		geom2d.NewPoint(263, 193),
		geom2d.NewPoint(253, 210),
		geom2d.NewPoint(290, 181),
		geom2d.NewPoint(288, 218),
	}, geom2d.PTSolid, geom2d.WithChildren(holeInPolyTree))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Initialize the application and geometries.
	app := new(App)
	err := app.initGeometry()
	if err != nil {
		log.Fatal(err)
	}

	// Configure Ebitengine settings and start the game.
	ebiten.SetWindowSize(700, 500)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("geom2d demo: line segment relationship to...")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	err = ebiten.RunGame(app)
	if err != nil {
		log.Fatal(err)
	}
}
