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

	"golang.org/x/image/font/basicfont"
)

func TestTextProcessor(t *testing.T) {
	// 创建一个测试图像
	img := createTextTestImage(400, 300)

	// 测试普通文本
	t.Run("NormalText", func(t *testing.T) {
		// 创建文本处理器
		processor := NewTextProcessor(TextOptions{
			Text:     "Hello, World!",
			Position: image.Point{50, 50},
			Font:     basicfont.Face7x13,
			Color:    color.RGBA{255, 0, 0, 255}, // 红色
		})

		// 处理图片
		_, err := processor.Process(img)
		if err != nil {
			t.Fatalf("处理图片失败: %v", err)
		}
	})

	// 测试旋转文本
	t.Run("RotatedText", func(t *testing.T) {
		// 创建文本处理器
		processor := NewTextProcessor(TextOptions{
			Text:     "Rotated Text!",
			Position: image.Point{200, 150},
			Font:     basicfont.Face7x13,
			Color:    color.RGBA{0, 0, 255, 255}, // 蓝色
			Angle:    45,                         // 旋转45度
		})

		// 处理图片
		_, err := processor.Process(img)
		if err != nil {
			t.Fatalf("处理图片失败: %v", err)
		}
	})

	// 测试使用WithAngle方法
	t.Run("WithAngle", func(t *testing.T) {
		// 创建文本处理器
		processor := NewTextProcessor(TextOptions{
			Text:     "With Angle Method",
			Position: image.Point{200, 200},
			Font:     basicfont.Face7x13,
			Color:    color.RGBA{0, 255, 0, 255}, // 绿色
		}).WithAngle(90) // 旋转90度

		// 处理图片
		_, err := processor.Process(img)
		if err != nil {
			t.Fatalf("处理图片失败: %v", err)
		}
	})
}

// 创建一个测试图像
func createTextTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充浅灰色背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}

	return img
}
