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

package examples

import (
	"fmt"
	"image"
	"image/color"

	"github.com/vogo/vimage"
)

func ExampleDrawCircleProcessor() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	x0, y0, r := 50, 50, 20
	processor := vimage.NewDrawCircleProcessor(x0, y0, r, color.RGBA{255, 0, 0, 255}, false)
	result, _ := processor.Process(img)

	c := result.At(x0+r-1, y0).(color.RGBA)
	// Due to anti-aliasing, the edge color may not be exactly 255
	fmt.Println(c.R > 0, c.G, c.B)
	// Output: true 0 0
}

func ExampleDrawCircleProcessor_filled() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	x0, y0, r := 50, 50, 20
	processor := vimage.NewDrawCircleProcessor(x0, y0, r, color.RGBA{255, 0, 0, 255}, true)
	result, _ := processor.Process(img)

	// Check center of the filled circle
	c := result.At(x0, y0).(color.RGBA)
	fmt.Println(c.R, c.G, c.B)
	// Output: 255 0 0
}

func ExampleDrawRectProcessor() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)
	processor := vimage.NewDrawRectProcessor(rect, color.RGBA{0, 255, 0, 255}, true)
	result, _ := processor.Process(img)

	c := result.At(50, 50).(color.RGBA)
	fmt.Println(c.R, c.G, c.B)
	// Output: 0 255 0
}

func ExampleDrawRectProcessor_stroked() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)
	processor := vimage.NewDrawRectProcessor(rect, color.RGBA{0, 255, 0, 255}, false)
	result, _ := processor.Process(img)

	// Check edge point
	c := result.At(20, 20).(color.RGBA)
	fmt.Println(c.G > 0)
	// Output: true
}

func ExampleDrawRectProcessor_withFillColor() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	rect := image.Rect(20, 20, 80, 80)
	// Blue border, yellow fill
	processor := vimage.NewDrawRectProcessorWithFillColor(
		rect,
		color.RGBA{0, 0, 255, 255},   // blue border
		color.RGBA{255, 255, 0, 255}, // yellow fill
	)
	result, _ := processor.Process(img)

	c := result.At(50, 50).(color.RGBA)
	// Inside should be yellow
	fmt.Println(c.R > 200 && c.G > 200 && c.B < 50)
	// Output: true
}
