# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

vimage is a Go image processing library (图像处理工具库) that provides rich image manipulation capabilities including resizing, cropping, mosaic effects, watermarks, noise generation, and CAPTCHA generation.

## Build Commands

```bash
# Run all checks and tests (complete build pipeline)
make build

# Individual commands
make format          # Format code with goimports, go fmt, and gofumpt
make license-check   # Verify Apache 2.0 license headers
make lint           # Run golangci-lint
make test           # Run all tests with verbose output

# Run a single test
go test -v -run TestName
go test -v -run TestName ./path/to/package
```

## Required Development Tools

Install these tools before contributing:

```bash
go install github.com/vogo/license-header-checker/cmd/license-header-checker@latest
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Architecture Overview

### Processor Pattern

The library is built around a **Processor** interface that enables composable image transformations:

```go
type Processor interface {
    Process(img image.Image) (image.Image, error)
}

type ContextProcessor interface {
    ContextProcess(ctx *ImageProcessContext) error
}
```

Key architectural concepts:

1. **Processor Chain**: Multiple processors can be chained together using `ProcessImage()` or `Process()` functions. Each processor transforms the image and passes it to the next processor.

2. **Two Processor Types**:
   - **Regular Processors** (`Processor` interface): Most processors (Cut, Zoom, Circle, etc.) implement this. They take an `image.Image`, transform it, and return a new `image.Image`.
   - **Context Processors** (`ContextProcessor` interface): Drawing processors (DrawCircle, DrawRect, Text, Watermark) implement this. They operate on `ImageProcessContext` which wraps `gg.Context` from fogleman/gg library for drawing operations. Use `ContextProcess()` function to execute a chain of ContextProcessors.

3. **Processing Functions**:
   - `Process(img, []Processor)` - Applies regular processor chain to an image.Image
   - `ProcessImage(imgData, []Processor, options)` - Higher-level API that decodes bytes, applies regular processors, and encodes back to bytes
   - `ContextProcess(img, []ContextProcessor)` - Applies context processor chain for drawing operations

### Core Processor Categories

- **Geometric Transformations** (Processor): Cut, Zoom (with 6 modes: Exact, Ratio, Width, Height, Max, Min)
- **Shape Processors** (Processor): CutSquare, CutCircle, RoundedCorner
- **Drawing Processors** (ContextProcessor): DrawCircle, DrawRect - These implement the ContextProcessor interface
- **Effects** (Processor): Mosaic, Noise, Rotate
- **Overlays** (ContextProcessor): Watermark, Overlay, Text - These implement the ContextProcessor interface
- **Generation**: Captcha, Table (standalone functions, not processors)

### Important Implementation Details

1. **CutCircleProcessor requires square images** - The input image must have equal width and height, otherwise it returns an error.

2. **Zoom Modes** (zoom.go:28-44):
   - `ZoomModeExact`: Scale to exact width/height
   - `ZoomModeRatio`: Scale by ratio (e.g., 0.5 for 50%)
   - `ZoomModeWidth`: Scale to width, adjust height proportionally
   - `ZoomModeHeight`: Scale to height, adjust width proportionally
   - `ZoomModeMax`: Scale by maximum dimension while preserving aspect ratio
   - `ZoomModeMin`: Scale by minimum dimension while preserving aspect ratio

3. **Cut Positions** (cut.go:26-39): "center", "top", "bottom", "left", "right"

4. **Image Format Handling** (processor.go:78-113): The `ProcessImage` function automatically detects input format (JPEG/PNG) and preserves it in output. JPEG quality is configurable via `ProcessorOptions`.

### Typical Processor Composition Pattern

Common workflow: Square crop → Additional processing (circle, zoom, etc.)

```go
processors := []vimage.Processor{
    vimage.NewCutSquareProcessor("center"),  // First, make it square
    vimage.NewCutCircleProcessor(),          // Then apply circle crop
}
result, err := vimage.ProcessImage(imgData, processors, nil)
```

## Code Standards

- All files must include Apache License 2.0 header (checked by `make license-check`)
- Code must pass `golangci-lint` checks
- All exported functions and types require documentation comments (Chinese or English)
- Use `fogleman/gg` for drawing operations (ContextProcessor implementations)
- Use `golang.org/x/image/draw` for scaling operations
- Use Go 1.22+ range-over-int syntax (`for i := range N`) instead of traditional C-style loops

## Writing Tests

When creating tests, create test images programmatically using `image.NewRGBA()` instead of loading from files:

```go
func TestExample(t *testing.T) {
    // Create a test image programmatically
    img := image.NewRGBA(image.Rect(0, 0, 200, 200))

    // Fill with test colors
    for y := range 200 {
        for x := range 200 {
            img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
        }
    }

    // Test your processor
    processor := NewMyProcessor()
    result, err := processor.Process(img)
    require.NoError(t, err)

    // Validate results using assertions on pixel colors or image properties
    assert.Equal(t, expectedValue, actualValue)
}
```

Use `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require` for test assertions.
