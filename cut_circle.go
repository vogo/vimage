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
	"errors"
	"image"
	"image/color"
	"math"
)

// CutCircleProcessor implements the Processor interface for circular image cropping
type CutCircleProcessor struct{}

func NewCutCircleProcessor() *CutCircleProcessor {
	return &CutCircleProcessor{}
}

// Process cuts the image into a circle.
func (p *CutCircleProcessor) Process(img image.Image) (image.Image, error) {
	return Circle(img)
}

// Circle crops the image into a circle, making pixels outside the circle transparent
// If the image is not square, returns an error
func Circle(img image.Image) (image.Image, error) {
	// Check if image is square
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width != height {
		return nil, errors.New("image must be square for circular cropping")
	}

	// Get circle radius (default to half of width/height)
	radius := float64(width) / 2

	// Create new RGBA image with transparency
	dst := image.NewRGBA(bounds)
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
			if distance <= radius {
				dst.Set(x, y, img.At(x, y))
			} else {
				// Set transparent
				dst.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	return dst, nil
}
