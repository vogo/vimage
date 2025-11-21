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

// DrawRectProcessor draws a rectangle on the image using fogleman/gg.
type DrawRectProcessor struct {
	Rect      image.Rectangle
	Color     color.Color  // Border/stroke color
	FillColor *color.Color // Optional fill color (nil = no fill, use Color for both)
	Fill      bool         // If true, fill the rectangle
}

// NewDrawRectProcessor creates a new DrawRectProcessor.
// If fill is true, the rectangle will be filled with the same color as the border.
func NewDrawRectProcessor(rect image.Rectangle, c color.Color, fill bool) *DrawRectProcessor {
	return &DrawRectProcessor{
		Rect:      rect,
		Color:     c,
		FillColor: nil, // Use same color for fill
		Fill:      fill,
	}
}

// NewDrawRectProcessorWithFillColor creates a DrawRectProcessor with separate border and fill colors.
// The rectangle will be filled with fillColor and stroked with borderColor.
func NewDrawRectProcessorWithFillColor(rect image.Rectangle, borderColor, fillColor color.Color) *DrawRectProcessor {
	return &DrawRectProcessor{
		Rect:      rect,
		Color:     borderColor,
		FillColor: &fillColor,
		Fill:      true,
	}
}

// ContextProcess draws a rectangle on the image using the gg.Context.
func (p *DrawRectProcessor) ContextProcess(ctx *ImageProcessContext) error {
	dc := ctx.DC()

	// Draw the rectangle
	dc.DrawRectangle(
		float64(p.Rect.Min.X),
		float64(p.Rect.Min.Y),
		float64(p.Rect.Dx()),
		float64(p.Rect.Dy()),
	)

	if p.Fill {
		// Set fill color
		if p.FillColor != nil {
			dc.SetColor(*p.FillColor)
		} else {
			dc.SetColor(p.Color)
		}

		// Fill first
		if p.FillColor != nil {
			// Fill and preserve the path
			dc.FillPreserve()
			// Then stroke with border color
			dc.SetColor(p.Color)
			dc.Stroke()
		} else {
			// Just fill with the same color
			dc.Fill()
		}
	} else {
		// Just stroke
		dc.SetColor(p.Color)
		dc.Stroke()
	}

	return nil
}

// Process implements the Processor interface for backward compatibility.
func (p *DrawRectProcessor) Process(img image.Image) (image.Image, error) {
	return ContextProcess(img, []ContextProcessor{p})
}
