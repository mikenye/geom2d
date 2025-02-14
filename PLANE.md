# Geom2D Plane Specification

## Purpose

The `Plane` type serves as a container for geometric primitives (`Point`, `LineSegment`, etc.), enforcing a consistent precision level set by the user. It provides deduplication, boundary enforcement, and memory management while allowing efficient spatial computations.

## Core Features & Behavior

1. **Precision Control**
    - The `Plane` ensures all geometry follows a consistent numerical precision set by the user.
    - Precision is defined as a scale factor (default: `1e-10`), allowing floating-point coordinates to be stored as integers.
    - All floating-point inputs are scaled and rounded to maintain precision consistency.
    - Conversions:
      - User Input: `float64`
      - Internal Storage: `int64`
      - User Output: `float64`
2. **Object Deduplication & Identity Management**
    - The `Plane` deduplicates geometric objects using `sync.Map`, preventing duplicate Point or LineSegment instances.
    - Users interacting with the same logical point always receive the same pointer reference.
    - Object creation is idempotent—calling `NewPoint(1.23, 4.56)` twice will return the same instance.
[//]: # (TODO: We may end up using a regular map... Revisit this pnce created.)
3. **Boundary Constraints**
    - The `Plane` defines a bounding box (`minX`, `maxX`, `minY`, `maxY`) depending on the set precision.
    - Any object that falls outside these bounds results in an error.
    - Boundary enforcement applies to all geometric primitives in the plane.
4. **Memory Management**
    - The `Plane` manages object lifecycle to prevent memory leaks and unnecessary allocations.
    - Two approaches ensure efficient memory handling:
      - Manual deletion: Users can explicitly remove objects with `RemovePoint()`, `RemoveLineSegment()`, etc.
      - Automatic GC cleanup: `runtime.SetFinalizer()` is used to remove objects from `sync.Map` when they are dereferenced.
5. **Retrieval & Queries**
    - `Plane` provides methods to retrieve all stored objects:
      - `AllPoints() []*Point`
      - `AllLineSegments() []*LineSegment`
    - These functions scan `sync.Map` and return slices of active objects.

## 🔍 How Does Precision Affect Bounds?

The maximum and minimum bounds of the plane depend on **precision** because all geometry is scaled to **integer coordinates** to prevent floating-point precision errors.

### ✅ Formula for Bounds Calculation
\[
\text{maxFloat} = \frac{\text{math.MaxInt64}}{\text{scaleFactor}}
\]
\[
\text{minFloat} = -\text{maxFloat}
\]

Where:
- `scaleFactor = 1 / precision`
- `math.MaxInt64 ≈ 9.22 × 10^{18}` (largest signed `int64` value)
- This ensures that when we **convert back to integers**, we **don’t exceed `int64` limits**.

### How Precision Affects Bounds

| **Precision** (`p`) | **Scale Factor** (`1 / p`) | **Max Bound** (`math.MaxInt64 / scaleFactor`) |
|---------------------|--------------------------|----------------------------------|
| `1e-1` (`0.1`)      | `10`                      | `≈ 9.22 × 10^{17}` |
| `1e-3` (`0.001`)    | `1,000`                   | `≈ 9.22 × 10^{15}` |
| `1e-6` (`0.000001`) | `1,000,000`               | `≈ 9.22 × 10^{12}` |
| `1e-9` (`0.000000001`) | `1,000,000,000`       | `≈ 9.22 × 10^{9}` |
| `1e-12` (`0.000000000001`) | `1,000,000,000,000` | `≈ 9.22 × 10^{6}` |
| `1e-18` (`0.000000000000000001`) | `10^{18}` | `≈ 9.22` (very small bounds) |

- **Smaller precision means smaller bounds** (higher detail, but less area).  
- **Higher precision gives more range, but less granularity.**  
- **`math.MaxInt64 / scaleFactor` ensures no integer overflows.**  
