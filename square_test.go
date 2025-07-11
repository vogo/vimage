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
)

func TestSquareProcessor(t *testing.T) {
	// 创建一个非正方形测试图片 (200x100)
	img := image.NewRGBA(image.Rect(0, 0, 200, 100))

	// 填充不同区域的颜色以便于测试
	// 左侧区域填充红色
	for y := 0; y < 100; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // 红色
		}
	}

	// 中间区域填充绿色
	for y := 0; y < 100; y++ {
		for x := 50; x < 150; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255}) // 绿色
		}
	}

	// 右侧区域填充蓝色
	for y := 0; y < 100; y++ {
		for x := 150; x < 200; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255}) // 蓝色
		}
	}

	// 测试居中裁剪
	processor := NewSquareProcessor("center")
	result, err := processor.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// 验证结果是正方形
	bounds := result.Bounds()
	if bounds.Dx() != bounds.Dy() {
		t.Errorf("Result should be square, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	// 验证尺寸是较小的一边
	if bounds.Dx() != 100 {
		t.Errorf("Square size should be 100, got %d", bounds.Dx())
	}

	// 验证中心点颜色是绿色（居中裁剪）
	if c := result.At(50, 50).(color.RGBA); c.G != 255 || c.R != 0 || c.B != 0 {
		t.Errorf("Center pixel should be green, got %v", c)
	}

	// 测试左侧裁剪
	processorLeft := NewSquareProcessor("left")
	resultLeft, err := processorLeft.Process(img)
	if err != nil {
		t.Fatalf("Process with left position failed: %v", err)
	}

	// 验证左侧裁剪的中心点是红色或绿色
	if c := resultLeft.At(25, 50).(color.RGBA); c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("Left area pixel should be red, got %v", c)
	}

	// 测试右侧裁剪
	processorRight := NewSquareProcessor("right")
	resultRight, err := processorRight.Process(img)
	if err != nil {
		t.Fatalf("Process with right position failed: %v", err)
	}

	// 验证右侧裁剪的中心点是蓝色或绿色
	if c := resultRight.At(75, 50).(color.RGBA); c.B != 255 || c.R != 0 || c.G != 0 {
		t.Errorf("Right area pixel should be blue, got %v", c)
	}

	// 测试已经是正方形的图片
	squareImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			squareImg.Set(x, y, color.RGBA{255, 255, 0, 255}) // 黄色
		}
	}

	resultSquare, err := processor.Process(squareImg)
	if err != nil {
		t.Fatalf("Process square image failed: %v", err)
	}

	// 验证尺寸没有变化
	squareBounds := resultSquare.Bounds()
	if squareBounds.Dx() != 100 || squareBounds.Dy() != 100 {
		t.Errorf("Square image size should remain 100x100, got %dx%d", squareBounds.Dx(), squareBounds.Dy())
	}
}
