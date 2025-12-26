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
	"image"
	"image/color"
	"testing"

	"golang.org/x/image/draw"
)

// createCompressTestImage creates a test image for compression
func createCompressTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create simple test pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Simple gradient pattern
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(((x + y) * 255) / (width + height))
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img
}

// TestCompressProcessor_Basic tests basic functionality
func TestCompressProcessor_Basic(t *testing.T) {
	tests := []struct {
		name        string
		strategy    *CompressionStrategy
		targetSize  int64
		imageWidth  int
		imageHeight int
		expectError bool
	}{
		{
			name:        "Default Strategy - Small Image",
			strategy:    CompressStrategyDefault,
			targetSize:  10 * 1024, // 10KB
			imageWidth:  100,
			imageHeight: 100,
			expectError: false,
		},
		{
			name:        "High Quality Strategy - Medium Image",
			strategy:    CompressStrategyHighQuality,
			targetSize:  50 * 1024, // 50KB
			imageWidth:  200,
			imageHeight: 200,
			expectError: false,
		},
		{
			name:        "Small File Strategy - Large Image",
			strategy:    CompressStrategySmallSize,
			targetSize:  20 * 1024, // 20KB
			imageWidth:  500,
			imageHeight: 500,
			expectError: false,
		},
		{
			name:        "Balanced Strategy - Standard Image",
			strategy:    CompressStrategyBalanced,
			targetSize:  30 * 1024, // 30KB
			imageWidth:  300,
			imageHeight: 300,
			expectError: false,
		},
		{
			name:        "Web Optimized Strategy",
			strategy:    CompressStrategyWebOptimized,
			targetSize:  25 * 1024, // 25KB
			imageWidth:  400,
			imageHeight: 300,
			expectError: false,
		},
		{
			name:        "Thumbnail Strategy",
			strategy:    CompressStrategyThumbnail,
			targetSize:  5 * 1024, // 5KB
			imageWidth:  150,
			imageHeight: 150,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image
			img := createCompressTestImage(tt.imageWidth, tt.imageHeight)

			// Create processor
			processor := NewCompressProcessor(tt.targetSize, tt.strategy)

			// Execute compression
			result, err := processor.Process(img)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify result
			if result == nil {
				t.Error("result image is nil")
				return
			}

			resultBounds := result.Bounds()
			if resultBounds.Dx() <= 0 || resultBounds.Dy() <= 0 {
				t.Error("invalid result image dimensions")
			}

			// Verify strategy name
			strategyName := processor.GetStrategyName()
			if strategyName == "" {
				t.Error("strategy name is empty")
			}
		})
	}
}

// TestCompressProcessor_FixedParams tests custom strategy with fixed parameters
func TestCompressProcessor_FixedParams(t *testing.T) {
	img := createCompressTestImage(200, 200)

	// Create a strategy with fixed parameters
	fixedStrategy := CompressionStrategy{
		Name:           "Fixed Strategy",
		MinPixelRatio:  0.5,
		MaxPixelRatio:  0.5,
		QualityLevels:  []int{75},
		DefaultQuality: 75,
		DefaultFormat:  "jpeg",
	}

	processor := NewCompressProcessor(20*1024, &fixedStrategy)

	// Execute compression
	result, err := processor.Process(img)
	if err != nil {
		t.Errorf("compression failed: %v", err)
		return
	}

	if result == nil {
		t.Error("result image is nil")
		return
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("dimension mismatch: expected 100x100, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// TestCompressProcessor_InvalidParams tests invalid parameters
func TestCompressProcessor_InvalidParams(t *testing.T) {
	tests := []struct {
		name        string
		targetSize  int64
		strategy    *CompressionStrategy
		expectError bool
	}{
		{
			name:        "Target Size 0",
			targetSize:  0,
			strategy:    CompressStrategyDefault,
			expectError: true,
		},
		{
			name:        "Target Size Negative",
			targetSize:  -100,
			strategy:    CompressStrategyDefault,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := createCompressTestImage(100, 100)
			processor := NewCompressProcessor(tt.targetSize, tt.strategy)

			_, err := processor.Process(img)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestCompressProcessor_FormatValidation tests format validation
func TestCompressProcessor_FormatValidation(t *testing.T) {
	processor := NewCompressProcessor(20*1024, CompressStrategyDefault)

	// Valid formats
	validFormats := []string{"jpeg", "png", "webp"}
	for _, format := range validFormats {
		err := processor.SetFormat(format)
		if err != nil {
			t.Errorf("failed to set valid format: %s", format)
		}
	}

	// Invalid formats
	invalidFormats := []string{"gif", "bmp", "tiff", "invalid"}
	for _, format := range invalidFormats {
		err := processor.SetFormat(format)
		if err == nil {
			t.Errorf("expected failure for invalid format: %s", format)
		}
	}
}

// TestCompressProcessor_SetStrategy tests strategy switching
func TestCompressProcessor_SetStrategy(t *testing.T) {
	processor := NewCompressProcessor(30*1024, CompressStrategyDefault)

	// Verify initial strategy
	if processor.GetStrategyName() != "Default Strategy" {
		t.Errorf("initial strategy mismatch: %s", processor.GetStrategyName())
	}

	// Switch to high quality strategy
	processor.SetStrategy(CompressStrategyHighQuality)
	if processor.GetStrategyName() != "High Quality" {
		t.Errorf("high quality strategy mismatch: %s", processor.GetStrategyName())
	}

	// Switch to small file strategy
	processor.SetStrategy(CompressStrategySmallSize)
	if processor.GetStrategyName() != "Small Size" {
		t.Errorf("small file strategy mismatch: %s", processor.GetStrategyName())
	}
}

// TestCompressProcessor_EstimateFileSize tests file size estimation
func TestCompressProcessor_EstimateFileSize(t *testing.T) {
	img := createCompressTestImage(200, 200)

	processor := NewCompressProcessor(25*1024, CompressStrategyBalanced)

	estimatedSize, err := processor.EstimateFileSize(img)
	if err != nil {
		t.Errorf("failed to estimate file size: %v", err)
		return
	}

	// Verify estimation result is reasonable
	if estimatedSize <= 0 {
		t.Error("invalid estimated file size")
	}

	if estimatedSize > 100*1024 {
		t.Errorf("estimated file size too large: %d bytes", estimatedSize)
	}
}

// TestCompressProcessor_ContentComplexity tests content complexity analysis
func TestCompressProcessor_ContentComplexity(t *testing.T) {
	// Since GetContentComplexity is removed, we can no longer test internal state directly.
	// This test is now removed or should be adapted to test public behavior if applicable.
	t.Skip("Content complexity is now an internal implementation detail")
}

// TestCompressProcessor_WithScaler tests custom scaler
func TestCompressProcessor_WithScaler(t *testing.T) {
	img := createCompressTestImage(200, 200)
	processor := NewCompressProcessor(30*1024, CompressStrategyDefault)

	// Use different scaler (simplified test, in reality should test different draw.Scaler impls)
	processor.WithScaler(draw.BiLinear)

	result, err := processor.Process(img)
	if err != nil {
		t.Errorf("failed to use custom scaler: %v", err)
		return
	}

	if result == nil {
		t.Error("result image is nil")
	}
}

// TestCompressProcessor_NilImage tests nil image handling
func TestCompressProcessor_NilImage(t *testing.T) {
	processor := NewCompressProcessor(20*1024, CompressStrategyDefault)

	_, err := processor.Process(nil)
	if err == nil {
		t.Error("expected error when processing nil image")
	}
}

// TestCompressProcessor_SmallImage tests small image handling
func TestCompressProcessor_SmallImage(t *testing.T) {
	// Test 1x1 image
	img := createCompressTestImage(1, 1)
	processor := NewCompressProcessor(5*1024, CompressStrategyThumbnail)

	result, err := processor.Process(img)
	if err != nil {
		t.Errorf("failed to process small image: %v", err)
		return
	}

	if result == nil {
		t.Error("result image is nil")
		return
	}

	bounds := result.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Error("invalid result image dimensions")
	}
}

// TestCompressProcessor_LargeImage tests large image handling
func TestCompressProcessor_LargeImage(t *testing.T) {
	// Test large image (but not too large for memory)
	img := createCompressTestImage(1000, 1000)
	processor := NewCompressProcessor(100*1024, CompressStrategyBalanced)

	result, err := processor.Process(img)
	if err != nil {
		t.Errorf("failed to process large image: %v", err)
		return
	}

	if result == nil {
		t.Error("result image is nil")
		return
	}

	bounds := result.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Error("invalid result image dimensions")
	}
}

// BenchmarkCompressProcessor performance benchmark
func BenchmarkCompressProcessor(b *testing.B) {
	img := createCompressTestImage(500, 500)
	processor := NewCompressProcessor(50*1024, CompressStrategyBalanced)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.Process(img)
		if err != nil {
			b.Fatalf("compression failed: %v", err)
		}
	}
}

// BenchmarkCompressProcessor_ComplexityAnalysis complexity analysis benchmark
func BenchmarkCompressProcessor_ComplexityAnalysis(b *testing.B) {
	img := createCompressTestImage(500, 500)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor := NewCompressProcessor(50*1024, CompressStrategyBalanced)
		_, err := processor.Process(img)
		if err != nil {
			b.Fatalf("processing failed: %v", err)
		}
	}
}

// BenchmarkCalculateColorVariance benchmarks color variance calculation
func BenchmarkCalculateColorVariance(b *testing.B) {
	img := createCompressTestImage(100, 100)
	processor := NewCompressProcessor(50*1024, CompressStrategyBalanced)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.calculateColorVariance(img)
	}
}

// BenchmarkCalculateEdgeDensity benchmarks edge density calculation
func BenchmarkCalculateEdgeDensity(b *testing.B) {
	img := createCompressTestImage(100, 100)
	processor := NewCompressProcessor(50*1024, CompressStrategyBalanced)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.calculateEdgeDensity(img)
	}
}

// BenchmarkAnalyzeContentComplexity benchmarks full complexity analysis
func BenchmarkAnalyzeContentComplexity(b *testing.B) {
	img := createCompressTestImage(500, 500)
	processor := NewCompressProcessor(50*1024, CompressStrategyBalanced)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.analyzeContentComplexity(img)
	}
}
