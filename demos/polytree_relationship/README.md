
# Geom2D Demo: PolyTree Relationships

This demo showcases the capabilities of the [Geom2D library](https://github.com/mikenye/geom2d) by visualizing relationships between a mouse-controlled polygon (diamond) and various geometric shapes (point, line segment, rectangle, circle, and polygon) in real-time.

![Red polygon (diamond) representing the mouse cursor moves around and over various geometries.](screenshot.gif "Animated Screenshot of Demo")

## How to Use

1. **Mouse Cursor**:
   - The mouse cursor is represented by a **red diamond** (a PolyTree) on the screen.
   - For visibility, the mouse's position is also displayed in the bottom-left corner.

2. **Real-Time Relationships**:
   - The relationships between the mouse rectangle and each geometric object are calculated in **real time**, at 60 ticks per second. This highlights the library's performance in dynamic scenarios.

3. **Real-Time Rendering**:
   - The geometries are redrawn every frame, running at 60 frames per second. This highlights the efficiency of the library for real-time rendering.

4. **Resizable Window**:
   - For high-resolution displays, you can resize the application window to make the pixels larger and improve visibility.

## Key Features Demonstrated

- Real-time calculation and visualization of geometric relationships.
- Dynamic rendering of various geometric objects, including:
  - Single points.
  - Line segments using Bresenham's line algorithm.
  - Rectangles and their edges.
  - Circles using Bresenham's circle algorithm.
  - Polygons and holes using `PolyTree`.

## System Requirements

- [Go](https://go.dev) version 1.23 or newer.
- This demo uses [Ebitengine](https://ebiten.org), a fantastic game library for Go, check it out!
- The window is set to an initial size of 700x500 but can be resized.

## Running the Demo

You should be able to run the demo directly from GitHub with the command:

```bash
go run github.com/mikenye/geom2d/demos/polytree_relationship@latest
```

If this doesn't work:

1. Clone the repository:
   ```bash
   git clone https://github.com/mikenye/geom2d
   cd geom2d/demos/polytree_relationship
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the demo:
   ```bash
   go run .
   ```
