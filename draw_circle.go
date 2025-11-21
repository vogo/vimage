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
)

// DrawCircleProcessor draws a circle on the image using fogleman/gg.
type DrawCircleProcessor struct {
	X, Y, Radius int
	Color        color.Color
	Fill         bool // If true, fill the circle with the same color as the border
}

// NewDrawCircleProcessor creates a new DrawCircleProcessor.
// If fill is true, the circle will be filled with the same color as the border.
func NewDrawCircleProcessor(x, y, radius int, c color.Color, fill bool) *DrawCircleProcessor {
	return &DrawCircleProcessor{
		X:      x,
		Y:      y,
		Radius: radius,
		Color:  c,
		Fill:   fill,
	}
}

// ContextProcess draws a circle on the image using the gg.Context.
func (p *DrawCircleProcessor) ContextProcess(ctx *ImageProcessContext) error {
	dc := ctx.DC()

	// Set the color
	dc.SetColor(p.Color)

	// Draw the circle
	dc.DrawCircle(float64(p.X), float64(p.Y), float64(p.Radius))

	// Fill or stroke based on the Fill parameter
	if p.Fill {
		dc.Fill()
	} else {
		dc.Stroke()
	}

	return nil
}

// Process implements the Processor interface for backward compatibility.
func (p *DrawCircleProcessor) Process(img image.Image) (image.Image, error) {
	return ContextProcess(img, []ContextProcessor{p})
}
