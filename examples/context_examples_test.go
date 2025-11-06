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
	"image/draw"
	"log"
	"testing"

	"github.com/vogo/vimage"
)

type circleProcessor struct {
	x, y, r int
	color   color.Color
}

func (p *circleProcessor) ContextProcess(ctx *vimage.ImageProcessContext) error {
	dc := ctx.DC()
	dc.SetColor(p.color)
	dc.DrawCircle(float64(p.x), float64(p.y), float64(p.r))
	dc.Fill()
	return nil
}

func ExampleContextProcess() {
	// Create a blank image
	img := image.NewRGBA(image.Rect(0, 0, 400, 200))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	// Define a series of processors to be applied in parallel
	processors := []vimage.ContextProcessor{
		&circleProcessor{x: 100, y: 100, r: 50, color: color.RGBA{R: 255, A: 255}},
		&circleProcessor{x: 200, y: 100, r: 40, color: color.RGBA{G: 255, A: 255}},
		&circleProcessor{x: 300, y: 100, r: 30, color: color.RGBA{B: 255, A: 255}},
	}

	// Process the image in parallel
	processedImg, err := vimage.ContextProcess(img, processors)
	if err != nil {
		log.Fatalf("Failed to process image in parallel: %v", err)
	}

	// Save the resulting image
	if err := saveImage(processedImg, "parallel_processed.png"); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	fmt.Println("Parallel processing example finished. Check parallel_processed.png")
}

func TestExampleContextProcessParallel(t *testing.T) {
	ExampleContextProcess()
}
