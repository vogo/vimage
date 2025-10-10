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
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

func TestCircleProcessor(t *testing.T) {
	// Create a test square image (100x100)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
		}
	}

	// Test with default radius
	processor := &CircleProcessor{}
	result, err := processor.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify center pixel is still red
	if c := result.At(50, 50).(color.RGBA); c.R != 255 || c.A != 255 {
		t.Errorf("Center pixel should be opaque red, got %v", c)
	}

	// Verify corner pixel is transparent
	if c := result.At(0, 0).(color.RGBA); c.A != 0 {
		t.Errorf("Corner pixel should be transparent, got %v", c)
	}

	// Test with non-square image (should error)
	nonSquare := image.NewRGBA(image.Rect(0, 0, 100, 80))
	_, err = processor.Process(nonSquare)
	if err == nil {
		t.Error("Expected error for non-square image")
	}
}

func TestCircleProcessLocalFile(t *testing.T) {
	b, err := os.ReadFile("build/avatar.jpg")
	if err != nil {
		t.Skipf("ReadFile failed: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Skipf("Decode failed: %v", err)
	}
	processor := &CircleProcessor{}
	result, err := processor.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}
	outputFile := "build/avatar_circle.jpg"
	_ = os.Remove(outputFile)
	f, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer func() { _ = f.Close() }()
	err = jpeg.Encode(f, result, &jpeg.Options{Quality: 90})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
}
