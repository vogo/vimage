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
	"fmt"
	"image"
	"image/color"
	"testing"
)

func TestRotateProcessor(t *testing.T) {
	img := createRotateTestImage(100, 50)

	// 测试不同角度的旋转
	angles := []float64{90, 180, 270, 45}

	for _, angle := range angles {
		t.Run("Rotate"+fmt.Sprintf("%v", angle), func(t *testing.T) {
			// 创建旋转处理器
			processor := NewRotateProcessor(angle)

			// 处理图片
			result, err := processor.Process(img)
			if err != nil {
				t.Fatalf("旋转处理失败: %v", err)
			}

			// 验证结果不为空
			if result == nil {
				t.Fatal("旋转结果为空")
			}

			// 对于90度和270度旋转，验证宽高是否交换（非KeepSize模式）
			if angle == 90 || angle == 270 {
				origBounds := img.Bounds()
				resultBounds := result.Bounds()

				// 允许1像素的误差（由于舍入）
				if absDiff(resultBounds.Dx()-origBounds.Dy()) > 1 || absDiff(resultBounds.Dy()-origBounds.Dx()) > 1 {
					t.Errorf("旋转%v度后尺寸不正确: 期望约 %dx%d, 实际 %dx%d",
						angle, origBounds.Dy(), origBounds.Dx(), resultBounds.Dx(), resultBounds.Dy())
				}
			}

			// 测试保持原始尺寸
			processor = NewRotateProcessor(angle).WithKeepSize(true)
			result, err = processor.Process(img)
			if err != nil {
				t.Fatalf("保持尺寸旋转处理失败: %v", err)
			}

			// 验证尺寸是否保持不变
			origBounds := img.Bounds()
			resultBounds := result.Bounds()
			if origBounds.Dx() != resultBounds.Dx() || origBounds.Dy() != resultBounds.Dy() {
				t.Errorf("保持尺寸旋转后尺寸不正确: 期望 %dx%d, 实际 %dx%d",
					origBounds.Dx(), origBounds.Dy(), resultBounds.Dx(), resultBounds.Dy())
			}

			// 测试自定义背景色
			redBg := color.RGBA{255, 0, 0, 255} // 红色背景
			processor = NewRotateProcessor(angle).WithBackground(redBg)
			result, err = processor.Process(img)
			if err != nil {
				t.Fatalf("自定义背景色旋转处理失败: %v", err)
			}

			// 验证结果不为空
			if result == nil {
				t.Fatal("自定义背景色旋转结果为空")
			}
		})
	}
}

// 创建一个简单的测试图像
func createRotateTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 绘制一个简单的图案，便于观察旋转效果
	// 左上角为红色
	for y := 0; y < height/2; y++ {
		for x := 0; x < width/2; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	// 右上角为绿色
	for y := 0; y < height/2; y++ {
		for x := width / 2; x < width; x++ {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}

	// 左下角为蓝色
	for y := height / 2; y < height; y++ {
		for x := 0; x < width/2; x++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	// 右下角为黄色
	for y := height / 2; y < height; y++ {
		for x := width / 2; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 255, 0, 255})
		}
	}

	return img
}

// 计算绝对值差值
func absDiff(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
