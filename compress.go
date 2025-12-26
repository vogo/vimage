/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vimage

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"

	"golang.org/x/image/draw"
)

// CompressionStrategy configures the compression strategy
type CompressionStrategy struct {
	Name             string   // Strategy name
	Description      string   // Strategy description
	MinPixelRatio    float64  // Minimum pixel ratio
	MaxPixelRatio    float64  // Maximum pixel ratio
	PreferredRatio   float64  // Preferred pixel ratio
	QualityLevels    []int    // Quality levels options
	DefaultQuality   int      // Default quality
	PreferredFormats []string // Preferred formats
	DefaultFormat    string   // Default format
	MaxIterations    int      // Maximum iterations
	Tolerance        float64  // Tolerance range
}

// Predefined strategies
var (
	CompressStrategyDefault = &CompressionStrategy{
		Name:             "Default Strategy",
		Description:      "General compression suitable for most scenarios",
		MinPixelRatio:    0.3,
		MaxPixelRatio:    1.0,
		PreferredRatio:   0.7,
		QualityLevels:    []int{95, 90, 85, 80, 75},
		DefaultQuality:   85,
		PreferredFormats: []string{"jpeg", "webp", "png"},
		DefaultFormat:    "jpeg",
		MaxIterations:    3,
		Tolerance:        0.1, // ±10%
	}

	CompressStrategyHighQuality = &CompressionStrategy{
		Name:             "High Quality",
		Description:      "Prioritize visual quality, file size is secondary",
		MinPixelRatio:    0.8,
		MaxPixelRatio:    1.0,
		PreferredRatio:   0.9,
		QualityLevels:    []int{98, 95, 92, 90},
		DefaultQuality:   95,
		PreferredFormats: []string{"png", "webp", "jpeg"},
		DefaultFormat:    "png",
		MaxIterations:    2,
		Tolerance:        0.2, // ±20%
	}

	CompressStrategySmallSize = &CompressionStrategy{
		Name:             "Small Size",
		Description:      "Minimize file size, quality loss is acceptable",
		MinPixelRatio:    0.2,
		MaxPixelRatio:    0.6,
		PreferredRatio:   0.4,
		QualityLevels:    []int{85, 75, 65, 60, 55},
		DefaultQuality:   70,
		PreferredFormats: []string{"webp", "jpeg"},
		DefaultFormat:    "webp",
		MaxIterations:    5,
		Tolerance:        0.05, // ±5%
	}

	CompressStrategyBalanced = &CompressionStrategy{
		Name:             "Balanced Strategy",
		Description:      "Balance between quality and file size",
		MinPixelRatio:    0.4,
		MaxPixelRatio:    0.8,
		PreferredRatio:   0.6,
		QualityLevels:    []int{90, 85, 80, 75, 70},
		DefaultQuality:   80,
		PreferredFormats: []string{"jpeg", "webp"},
		DefaultFormat:    "jpeg",
		MaxIterations:    3,
		Tolerance:        0.08, // ±8%
	}

	CompressStrategyWebOptimized = &CompressionStrategy{
		Name:             "Web Optimized",
		Description:      "Optimized for web loading speed",
		MinPixelRatio:    0.5,
		MaxPixelRatio:    0.9,
		PreferredRatio:   0.7,
		QualityLevels:    []int{85, 80, 75, 70},
		DefaultQuality:   80,
		PreferredFormats: []string{"webp", "jpeg"},
		DefaultFormat:    "webp",
		MaxIterations:    2,
		Tolerance:        0.1, // ±10%
	}

	CompressStrategyThumbnail = &CompressionStrategy{
		Name:             "Thumbnail Strategy",
		Description:      "Generate high quality thumbnails",
		MinPixelRatio:    0.1,
		MaxPixelRatio:    0.3,
		PreferredRatio:   0.2,
		QualityLevels:    []int{90, 85, 80},
		DefaultQuality:   85,
		PreferredFormats: []string{"jpeg", "webp"},
		DefaultFormat:    "jpeg",
		MaxIterations:    1,
		Tolerance:        0.15, // ±15%
	}
)

// CompressProcessor smart compression processor
type CompressProcessor struct {
	// Basic configuration
	TargetFileSize int64
	Strategy       *CompressionStrategy
	Format         string

	// Scaler
	scaler draw.Scaler
}

// NewCompressProcessor creates a new compression processor
func NewCompressProcessor(targetFileSize int64, strategy *CompressionStrategy) *CompressProcessor {
	if strategy == nil {
		strategy = CompressStrategyDefault
	}
	return &CompressProcessor{
		TargetFileSize: targetFileSize,
		Strategy:       strategy,
		Format:         strategy.DefaultFormat,
		scaler:         draw.BiLinear,
	}
}

// SetStrategy sets the compression strategy
func (p *CompressProcessor) SetStrategy(strategy *CompressionStrategy) {
	p.Strategy = strategy

	if p.Format == "" {
		p.Format = strategy.DefaultFormat
	}
}

// SetFormat sets the output format
func (p *CompressProcessor) SetFormat(format string) error {
	if format != "jpeg" && format != "png" && format != "webp" {
		return fmt.Errorf("unsupported format: %s", format)
	}
	p.Format = format
	return nil
}

// WithScaler sets the scaling algorithm
func (p *CompressProcessor) WithScaler(scaler draw.Scaler) *CompressProcessor {
	p.scaler = scaler
	return p
}

// Process implements the Processor interface
func (p *CompressProcessor) Process(img image.Image) (image.Image, error) {
	// Parameter validation
	if p.TargetFileSize <= 0 {
		return nil, fmt.Errorf("target file size must be greater than 0")
	}

	if img == nil {
		return nil, fmt.Errorf("input image cannot be nil")
	}

	// Execute compression
	result, _, err := p.doProcess(img)
	return result, err
}

// ProcessFile processes image file bytes
func (p *CompressProcessor) ProcessFile(input []byte) ([]byte, error) {
	// Check if file size is already within target
	if int64(len(input)) <= p.TargetFileSize {
		return input, nil
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Process image
	result, quality, err := p.doProcess(img)
	if err != nil {
		return nil, err
	}

	// Encode to specified format
	var buf bytes.Buffer
	switch p.Format {
	case "jpeg":
		err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, result)
	case "webp":
		// WebP encoding simplified, use jpeg in actual project if webp lib not available
		err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: quality})
	default:
		err = jpeg.Encode(&buf, result, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}

	return buf.Bytes(), nil
}

// doProcess calculates optimal parameters and compresses the image
func (p *CompressProcessor) doProcess(img image.Image) (image.Image, int, error) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Analyze content complexity
	contentComplexity := p.analyzeContentComplexity(img)

	// Calculate optimal parameters based on strategy
	optimalPixelRatio, optimalQuality, err := p.optimizeByStrategy(width, height, contentComplexity)
	if err != nil {
		return nil, 0, err
	}

	// Execute compression
	result, err := p.compress(img, optimalPixelRatio)
	if err != nil {
		return nil, 0, err
	}

	return result, optimalQuality, nil
}

// optimizeByStrategy optimizes parameters based on strategy
func (p *CompressProcessor) optimizeByStrategy(width, height int, contentComplexity float64) (float64, int, error) {
	currentPixels := float64(width * height)
	var optimalPixelRatio float64
	var optimalQuality int

	// Estimate base pixel compression ratio
	baseRatio := p.estimateBasePixelRatio(currentPixels, p.Strategy.DefaultQuality, contentComplexity)

	// Apply strategy constraints
	if p.Strategy.MinPixelRatio > 0 && baseRatio < p.Strategy.MinPixelRatio {
		optimalPixelRatio = p.Strategy.MinPixelRatio
	} else if p.Strategy.MaxPixelRatio > 0 && baseRatio > p.Strategy.MaxPixelRatio {
		optimalPixelRatio = p.Strategy.MaxPixelRatio
	} else {
		optimalPixelRatio = baseRatio
	}

	// Adjust quality based on content complexity
	optimalQuality = p.Strategy.DefaultQuality
	if contentComplexity > 1.5 {
		// Complex content, increase quality
		optimalQuality = p.findHigherQuality(optimalQuality)
	} else if contentComplexity < 0.8 {
		// Simple content, decrease quality
		optimalQuality = p.findLowerQuality(optimalQuality)
	}

	return optimalPixelRatio, optimalQuality, nil
}

// estimateBasePixelRatio estimates base pixel compression ratio
func (p *CompressProcessor) estimateBasePixelRatio(currentPixels float64, quality int, contentComplexity float64) float64 {
	// Estimate based on empirical formula
	// Base compression efficiency for different formats (bytes/pixel)
	formatEfficiency := map[string]float64{
		"jpeg": 0.12,
		"png":  0.25,
		"webp": 0.08,
	}

	efficiency := formatEfficiency[p.Format]
	if efficiency == 0 {
		efficiency = 0.12 // Default JPEG
	}

	// Quality factor
	qualityFactor := float64(quality) / 100.0
	if qualityFactor > 0.9 {
		qualityFactor = 1.8 // High quality requires more space
	} else if qualityFactor > 0.8 {
		qualityFactor = 1.2
	} else if qualityFactor > 0.7 {
		qualityFactor = 0.8
	} else {
		qualityFactor = 0.5
	}

	// Calculate target pixel count
	targetPixels := float64(p.TargetFileSize) / (efficiency * qualityFactor)

	// Consider content complexity
	targetPixels = targetPixels / contentComplexity

	// Calculate compression ratio
	return math.Sqrt(targetPixels / currentPixels)
}

// findHigherQuality finds a higher quality level
func (p *CompressProcessor) findHigherQuality(currentQuality int) int {
	for _, quality := range p.Strategy.QualityLevels {
		if quality > currentQuality {
			return quality
		}
	}
	if len(p.Strategy.QualityLevels) > 0 {
		return p.Strategy.QualityLevels[0] // Return highest quality
	}
	return currentQuality
}

// findLowerQuality finds a lower quality level
func (p *CompressProcessor) findLowerQuality(currentQuality int) int {
	for i := len(p.Strategy.QualityLevels) - 1; i >= 0; i-- {
		if p.Strategy.QualityLevels[i] < currentQuality {
			return p.Strategy.QualityLevels[i]
		}
	}
	if len(p.Strategy.QualityLevels) > 0 {
		return p.Strategy.QualityLevels[len(p.Strategy.QualityLevels)-1] // Return lowest quality
	}
	return currentQuality
}

// compress executes compression
func (p *CompressProcessor) compress(img image.Image, pixelRatio float64) (image.Image, error) {
	bounds := img.Bounds()
	origWidth, origHeight := bounds.Dx(), bounds.Dy()

	// Calculate target dimensions
	targetWidth := int(float64(origWidth) * pixelRatio)
	targetHeight := int(float64(origHeight) * pixelRatio)

	// Ensure minimum dimensions are 1x1
	if targetWidth < 1 {
		targetWidth = 1
	}
	if targetHeight < 1 {
		targetHeight = 1
	}

	// Validate target dimensions
	if targetWidth <= 0 || targetHeight <= 0 {
		return nil, fmt.Errorf("invalid target dimensions: %dx%d", targetWidth, targetHeight)
	}

	// Create target image
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Execute scaling
	p.scaler.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst, nil
}

// analyzeContentComplexity analyzes image content complexity
func (p *CompressProcessor) analyzeContentComplexity(img image.Image) float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width <= 0 || height <= 0 {
		return 1.0 // Default complexity
	}

	// Sampling analysis to avoid processing large images
	sampleSize := 100
	if width > sampleSize || height > sampleSize {
		// Scale down for sampling
		img = p.createSampleImage(img, sampleSize)
	}

	// Calculate color variance and edge density
	colorVariance := p.calculateColorVariance(img)
	edgeDensity := p.calculateEdgeDensity(img)

	// Combined complexity (0.5 - 2.0)
	complexity := 0.6 + 0.4*colorVariance + 0.6*edgeDensity

	// Limit range
	if complexity < 0.5 {
		complexity = 0.5
	} else if complexity > 2.0 {
		complexity = 2.0
	}

	return complexity
}

// createSampleImage creates a sample image
func (p *CompressProcessor) createSampleImage(img image.Image, maxSize int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Calculate scaling ratio
	ratio := float64(maxSize) / float64(math.Max(float64(width), float64(height)))
	if ratio >= 1.0 {
		return img // No scaling needed
	}

	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	sample := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	p.scaler.Scale(sample, sample.Bounds(), img, bounds, draw.Over, nil)

	return sample
}

// calculateColorVariance calculates color variance
// Optimized: single-pass using Welford's online algorithm
func (p *CompressProcessor) calculateColorVariance(img image.Image) float64 {
	bounds := img.Bounds()

	// Fast path: use type assertion for direct pixel access
	if rgba, ok := img.(*image.RGBA); ok {
		return p.calculateColorVarianceFast(rgba)
	}

	// Standard path: single-pass Welford's algorithm
	var meanR, meanG, meanB float64
	var m2R, m2G, m2B float64
	var count int64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			count++

			// Update R channel using Welford's algorithm
			deltaR := float64(r) - meanR
			meanR += deltaR / float64(count)
			m2R += deltaR * (float64(r) - meanR)

			// Update G channel
			deltaG := float64(g) - meanG
			meanG += deltaG / float64(count)
			m2G += deltaG * (float64(g) - meanG)

			// Update B channel
			deltaB := float64(b) - meanB
			meanB += deltaB / float64(count)
			m2B += deltaB * (float64(b) - meanB)
		}
	}

	if count == 0 {
		return 0
	}

	// Calculate variance from M2 values
	variance := (m2R + m2G + m2B) / (65535.0 * 65535.0 * 3.0 * float64(count))

	// Normalize to 0-1
	if variance > 0.1 {
		return 1.0
	}
	return variance * 10.0
}

// calculateColorVarianceFast calculates color variance with direct pixel access
// Fast path for *image.RGBA - approximately 3-5x faster than interface calls
func (p *CompressProcessor) calculateColorVarianceFast(rgba *image.RGBA) float64 {
	bounds := rgba.Bounds()
	pix := rgba.Pix
	stride := rgba.Stride

	var meanR, meanG, meanB float64
	var m2R, m2G, m2B float64
	var count int64

	// Direct pixel buffer access with Welford's algorithm
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			i := y*stride + x*4
			r := uint32(pix[i])
			g := uint32(pix[i+1])
			b := uint32(pix[i+2])
			count++

			// Update R channel using Welford's algorithm
			deltaR := float64(r) - meanR
			meanR += deltaR / float64(count)
			m2R += deltaR * (float64(r) - meanR)

			// Update G channel
			deltaG := float64(g) - meanG
			meanG += deltaG / float64(count)
			m2G += deltaG * (float64(g) - meanG)

			// Update B channel
			deltaB := float64(b) - meanB
			meanB += deltaB / float64(count)
			m2B += deltaB * (float64(b) - meanB)
		}
	}

	if count == 0 {
		return 0
	}

	// Calculate variance from M2 values (note: 8-bit values so use 255.0 instead of 65535.0)
	variance := (m2R + m2G + m2B) / (255.0 * 255.0 * 3.0 * float64(count))

	// Normalize to 0-1
	if variance > 0.1 {
		return 1.0
	}
	return variance * 10.0
}

// calculateEdgeDensity calculates edge density
// Optimized: pre-compute grayscale cache to avoid repeated conversions
func (p *CompressProcessor) calculateEdgeDensity(img image.Image) float64 {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width < 3 || height < 3 {
		return 0.0
	}

	// Fast path: use type assertion for direct pixel access
	if rgba, ok := img.(*image.RGBA); ok {
		return p.calculateEdgeDensityFast(rgba)
	}

	// Pre-compute grayscale values to avoid repeated RGBA conversions
	// Memory cost: ~10KB for 100x100 sample (negligible)
	grayValues := make([][]int, height)
	for y := range height {
		grayValues[y] = make([]int, width)
		for x := range width {
			r, g, b, _ := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			grayValues[y][x] = int((r*299 + g*587 + b*114) / 256000)
		}
	}

	// Use cached grayscale values for edge detection
	var edgeCount int
	var total int

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			center := grayValues[y][x]

			// Inline max calculation for better performance
			maxDiff := abs(grayValues[y][x-1] - center)
			if diff := abs(grayValues[y][x+1] - center); diff > maxDiff {
				maxDiff = diff
			}
			if diff := abs(grayValues[y-1][x] - center); diff > maxDiff {
				maxDiff = diff
			}
			if diff := abs(grayValues[y+1][x] - center); diff > maxDiff {
				maxDiff = diff
			}

			if maxDiff > 30 {
				edgeCount++
			}
			total++
		}
	}

	if total == 0 {
		return 0.0
	}

	return float64(edgeCount) / float64(total)
}

// calculateEdgeDensityFast calculates edge density with direct pixel access
// Fast path for *image.RGBA - approximately 3-5x faster than interface calls
func (p *CompressProcessor) calculateEdgeDensityFast(rgba *image.RGBA) float64 {
	bounds := rgba.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width < 3 || height < 3 {
		return 0.0
	}

	pix := rgba.Pix
	stride := rgba.Stride

	// Pre-compute grayscale values with direct pixel buffer access
	grayValues := make([][]int, height)
	for y := range height {
		grayValues[y] = make([]int, width)
		for x := range width {
			i := y*stride + x*4
			r := uint32(pix[i])
			g := uint32(pix[i+1])
			b := uint32(pix[i+2])
			grayValues[y][x] = int((r*299 + g*587 + b*114) / 1000)
		}
	}

	// Use cached grayscale values for edge detection
	var edgeCount int
	var total int

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			center := grayValues[y][x]

			// Inline max calculation for better performance
			maxDiff := abs(grayValues[y][x-1] - center)
			if diff := abs(grayValues[y][x+1] - center); diff > maxDiff {
				maxDiff = diff
			}
			if diff := abs(grayValues[y-1][x] - center); diff > maxDiff {
				maxDiff = diff
			}
			if diff := abs(grayValues[y+1][x] - center); diff > maxDiff {
				maxDiff = diff
			}

			if maxDiff > 30 {
				edgeCount++
			}
			total++
		}
	}

	if total == 0 {
		return 0.0
	}

	return float64(edgeCount) / float64(total)
}

// EstimateFileSize estimates compressed file size (for testing)
func (p *CompressProcessor) EstimateFileSize(img image.Image) (int64, error) {
	// Process image
	compressed, quality, err := p.doProcess(img)
	if err != nil {
		return 0, err
	}

	// Encode to specified format and calculate size
	var buf bytes.Buffer
	switch p.Format {
	case "jpeg":
		err = jpeg.Encode(&buf, compressed, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, compressed)
	case "webp":
		// WebP encoding simplified, use jpeg in actual project if webp lib not available
		err = jpeg.Encode(&buf, compressed, &jpeg.Options{Quality: quality})
	default:
		err = jpeg.Encode(&buf, compressed, &jpeg.Options{Quality: quality})
	}

	if err != nil {
		return 0, fmt.Errorf("encoding failed: %v", err)
	}

	return int64(buf.Len()), nil
}

// GetStrategyName gets current strategy name
func (p *CompressProcessor) GetStrategyName() string {
	return p.Strategy.Name
}
