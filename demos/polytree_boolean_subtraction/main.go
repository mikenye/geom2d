package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mikenye/geom2d"
	"golang.org/x/image/colornames"
	"log"
)

// App represents the state of the demo application.
type App struct {
	mouseX, mouseY int // Current mouse position.

	polyTreeA       *geom2d.PolyTree[int]
	polyTreeAOffset geom2d.Point[int]
	polyTreeB       *geom2d.PolyTree[int]
	polyTreeResult  *geom2d.PolyTree[int]
}

// Draw renders the current state of the application to the screen.
func (app *App) Draw(screen *ebiten.Image) {

	// Render the PolyTree A
	for poly := range app.polyTreeA.Nodes {
		for _, edge := range poly.Edges() {
			for p := range edge.Bresenham {
				screen.Set(p.X(), p.Y(), colornames.Darkred) // Render each edge of the polygon.
			}
		}
	}

	// Render the PolyTree B
	for poly := range app.polyTreeB.Nodes {
		for _, edge := range poly.Edges() {
			for p := range edge.Bresenham {
				screen.Set(p.X(), p.Y(), colornames.Darkblue) // Render each edge of the polygon.
			}
		}
	}

	// Render the Resulting PolyTree
	for poly := range app.polyTreeResult.Nodes {
		for _, edge := range poly.Edges() {
			for p := range edge.Bresenham {
				screen.Set(p.X(), p.Y(), colornames.White) // Render each edge of the polygon.
			}
		}
	}

	// Display mouse position, FPS, and TPS.
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Mouse: (%d, %d)", app.mouseX, app.mouseY), 5, 224)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), 200, 224)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()), 280, 224)
}

// Layout specifies the screen size for the application.
func (app *App) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 350, 250 // Fixed screen size.
}

// Update computes the relationships between the mouse point and other geometries.
func (app *App) Update() error {
	var err error

	// Update mouse position and recreate mouse geometry.
	fmt.Println("get cursor pos")
	app.mouseX, app.mouseY = ebiten.CursorPosition()

	//
	fmt.Println("app.polyTreeAOffset.Negate():", app.polyTreeAOffset.Negate())
	app.polyTreeB = app.polyTreeA.Translate(app.polyTreeAOffset.Negate())
	fmt.Println("geom2d.NewPoint(app.mouseX, app.mouseY):", geom2d.NewPoint(app.mouseX, app.mouseY))
	app.polyTreeB = app.polyTreeB.Translate(
		geom2d.NewPoint(app.mouseX, app.mouseY),
	)

	fmt.Println("PolyTreeA:", app.polyTreeA)
	fmt.Println("PolyTreeB:", app.polyTreeB)

	fmt.Println("perform boolean op")
	app.polyTreeResult, err = app.polyTreeA.BooleanOperation(
		app.polyTreeB,
		geom2d.BooleanSubtraction,
	)

	fmt.Println("polyTreeResult:", app.polyTreeResult)

	return err
}

// initGeometry initializes all geometric objects and handles potential errors.
func (app *App) initGeometry() error {
	var err error

	// Initialize PolyTreeA.
	holeA, err := geom2d.NewPolyTree[int]([]geom2d.Point[int]{
		geom2d.NewPoint(25, 25),
		geom2d.NewPoint(25, 75),
		geom2d.NewPoint(75, 75),
		geom2d.NewPoint(75, 25),
	}, geom2d.PTHole)
	if err != nil {
		return err
	}
	app.polyTreeA, err = geom2d.NewPolyTree[int]([]geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(100, 0),
		geom2d.NewPoint(100, 100),
		geom2d.NewPoint(0, 100),
	}, geom2d.PTSolid, geom2d.WithChildren(holeA))
	if err != nil {
		return err
	}
	app.polyTreeAOffset = geom2d.NewPoint(125, 75)
	app.polyTreeA = app.polyTreeA.Translate(app.polyTreeAOffset)
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
	ebiten.SetWindowTitle("geom2d demo: polytree boolean subtraction")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	err = ebiten.RunGame(app)
	if err != nil {
		log.Fatal(err)
	}
}
