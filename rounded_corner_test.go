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

func TestRoundedCornerProcessor(t *testing.T) {
	// 创建一个测试图像 (100x100 的红色方块)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // 红色
		}
	}

	// 测试用例
	tests := []struct {
		name   string
		radius int
	}{
		{"Zero Radius", 0},
		{"Small Radius", 10},
		{"Medium Radius", 25},
		{"Large Radius", 50},
		{"Oversized Radius", 100}, // 应该被限制为50
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// 创建处理器
			processor := NewRoundedCornerProcessor(test.radius)

			// 处理图像
			result, err := processor.Process(img)
			if err != nil {
				t.Fatalf("处理图像时出错: %v", err)
			}

			// 验证结果是否为RGBA图像
			_, ok := result.(*image.RGBA)
			if !ok {
				t.Errorf("结果应该是*image.RGBA类型，但得到了%T", result)
			}

			// 验证图像尺寸未变
			if result.Bounds().Dx() != 100 || result.Bounds().Dy() != 100 {
				t.Errorf("图像尺寸应该保持不变，期望100x100，得到%dx%d", 
					result.Bounds().Dx(), result.Bounds().Dy())
			}

			// 验证圆角效果
			if test.radius > 0 {
				// 检查四个角是否透明
				corners := []struct{ x, y int }{
					{0, 0},                   // 左上
					{99, 0},                  // 右上
					{0, 99},                  // 左下
					{99, 99},                 // 右下
				}

				for _, corner := range corners {
					r, g, b, a := result.At(corner.x, corner.y).RGBA()
					if a != 0 {
						t.Errorf("角点(%d,%d)应该是透明的，但得到了RGBA(%d,%d,%d,%d)", 
							corner.x, corner.y, r>>8, g>>8, b>>8, a>>8)
					}
				}

				// 检查中心点是否保持原色
				r, g, b, a := result.At(50, 50).RGBA()
				if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
					t.Errorf("中心点应该保持红色，但得到了RGBA(%d,%d,%d,%d)", 
						r>>8, g>>8, b>>8, a>>8)
				}
			}
		})
	}
}