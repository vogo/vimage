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

	"github.com/stretchr/testify/require"
)

func TestDrawRect_Filled(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)
	processor := NewDrawRectProcessor(rect, color.RGBA{0, 255, 0, 255}, true)
	result, err := processor.Process(img)
	require.NoError(t, err)

	// Check a point inside the filled rectangle
	c := result.At(50, 50).(color.RGBA)
	require.Equal(t, uint8(0), c.R)
	require.Equal(t, uint8(255), c.G)
	require.Equal(t, uint8(0), c.B)
}

func TestDrawRect_Stroked(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)
	processor := NewDrawRectProcessor(rect, color.RGBA{0, 255, 0, 255}, false)
	result, err := processor.Process(img)
	require.NoError(t, err)

	// Check a point on the edge (should have color)
	c := result.At(20, 20).(color.RGBA)
	require.True(t, c.G > 0, "Edge should have green color")

	// Check a point well inside (should be mostly background)
	// Note: Due to stroke width, center might still be transparent
	c = result.At(50, 50).(color.RGBA)
	// Center of stroked rectangle should be transparent (original background)
	require.Equal(t, uint8(0), c.R)
	require.Equal(t, uint8(0), c.G)
	require.Equal(t, uint8(0), c.B)
}

func TestDrawRect_WithFillColor(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)

	// Blue border, yellow fill
	borderColor := color.RGBA{0, 0, 255, 255}
	fillColor := color.RGBA{255, 255, 0, 255}
	processor := NewDrawRectProcessorWithFillColor(rect, borderColor, fillColor)
	result, err := processor.Process(img)
	require.NoError(t, err)

	// Check a point inside (should be yellow)
	c := result.At(50, 50).(color.RGBA)
	require.True(t, c.R > 200 && c.G > 200 && c.B < 50, "Inside should be yellow")

	// Check a point on the edge (should have blue from border)
	c = result.At(20, 20).(color.RGBA)
	require.True(t, c.B > 0, "Edge should have blue color from border")
}
