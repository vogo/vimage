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
	"testing"
)

func TestCutProcessor_Rectangle(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))

	// Test cutting a rectangle from center
	processor := NewCutProcessor(100, 80, CutPositionCenter)
	result, err := processor.Process(img)
	if err != nil {
		t.Fatal(err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 80 {
		t.Errorf("Expected 100x80, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestCutProcessor_WithRegion(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))

	// Test cutting with custom region
	processor := NewCutProcessorWithRegion(60, 40, 10, 20)
	result, err := processor.Process(img)
	if err != nil {
		t.Fatal(err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 60 || bounds.Dy() != 40 {
		t.Errorf("Expected 60x40, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestCutSquareProcessor_AutoSize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	positions := []string{"center", "top", "bottom", "left", "right"}

	for _, pos := range positions {
		t.Run(pos, func(t *testing.T) {
			processor := NewCutSquareProcessor(pos)
			result, err := processor.Process(img)
			if err != nil {
				t.Fatal(err)
			}

			// Check that result is square
			bounds := result.Bounds()
			if bounds.Dx() != bounds.Dy() {
				t.Errorf("Result should be square, got %dx%d", bounds.Dx(), bounds.Dy())
			}

			// Check that size is the smaller dimension (100)
			if bounds.Dx() != 100 {
				t.Errorf("Expected size 100, got %d", bounds.Dx())
			}
		})
	}
}

func TestCutSquareProcessor_WithSize(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))

	// Cut a 50x50 square from center
	processor := NewCutSquareProcessorWithSize(50, "center")
	result, err := processor.Process(img)
	if err != nil {
		t.Fatal(err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestCutSquareProcessor_WithRegion(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))

	// Cut a 60x60 square starting at (10, 20)
	processor := NewCutSquareProcessorWithRegion(60, 10, 20)
	result, err := processor.Process(img)
	if err != nil {
		t.Fatal(err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 60 || bounds.Dy() != 60 {
		t.Errorf("Expected 60x60, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestSquareCutProcessor_LegacyAPI(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 200, 150))

	// Test legacy NewSquareCutProcessor API
	processor := NewSquareCutProcessor(80, CutPositionCenter)
	result, err := processor.Process(img)
	if err != nil {
		t.Fatal(err)
	}

	bounds := result.Bounds()
	if bounds.Dx() != 80 || bounds.Dy() != 80 {
		t.Errorf("Expected 80x80, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
