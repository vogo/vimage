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
	"math"
)

// RoundedCornerProcessor 实现圆角处理器
// 将图片的四个角切割成圆角，角的大小可以通过半径参数控制
type RoundedCornerProcessor struct {
	// 圆角半径，单位为像素
	Radius int
}

// NewRoundedCornerProcessor 创建新的圆角处理器
// radius: 圆角半径，单位为像素
func NewRoundedCornerProcessor(radius int) *RoundedCornerProcessor {
	// 确保半径为正数
	if radius < 0 {
		radius = 0
	}

	return &RoundedCornerProcessor{
		Radius: radius,
	}
}

// Process 实现ImageProcessor接口
// 将图片的四个角切割成圆角，角外部分变为透明
func (p *RoundedCornerProcessor) Process(img image.Image) (image.Image, error) {
	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建新的RGBA图像（支持透明度）
	dst := image.NewRGBA(bounds)

	// 如果半径为0，直接返回原图
	if p.Radius <= 0 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				dst.Set(x, y, img.At(x, y))
			}
		}
		return dst, nil
	}

	// 确保半径不超过图片宽高的一半
	radius := p.Radius
	if radius > width/2 {
		radius = width / 2
	}
	if radius > height/2 {
		radius = height / 2
	}

	// 处理每个像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// 检查是否在四个角的圆角区域外
			if isInRoundedCorner(x, y, bounds, radius) {
				// 在圆角内部，保留原始像素
				dst.Set(x, y, img.At(x, y))
			} else {
				// 在圆角外部，设置为透明
				dst.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	return dst, nil
}

// isInRoundedCorner 判断像素是否在圆角区域内
// 返回true表示在圆角内部（保留），false表示在圆角外部（透明）
func isInRoundedCorner(x, y int, bounds image.Rectangle, radius int) bool {

	// 左上角
	if x < bounds.Min.X+radius && y < bounds.Min.Y+radius {
		// 计算到圆心的距离
		distance := math.Sqrt(math.Pow(float64(x-(bounds.Min.X+radius)), 2) + 
			math.Pow(float64(y-(bounds.Min.Y+radius)), 2))
		return distance <= float64(radius)
	}

	// 右上角
	if x >= bounds.Max.X-radius && y < bounds.Min.Y+radius {
		distance := math.Sqrt(math.Pow(float64(x-(bounds.Max.X-radius-1)), 2) + 
			math.Pow(float64(y-(bounds.Min.Y+radius)), 2))
		return distance <= float64(radius)
	}

	// 左下角
	if x < bounds.Min.X+radius && y >= bounds.Max.Y-radius {
		distance := math.Sqrt(math.Pow(float64(x-(bounds.Min.X+radius)), 2) + 
			math.Pow(float64(y-(bounds.Max.Y-radius-1)), 2))
		return distance <= float64(radius)
	}

	// 右下角
	if x >= bounds.Max.X-radius && y >= bounds.Max.Y-radius {
		distance := math.Sqrt(math.Pow(float64(x-(bounds.Max.X-radius-1)), 2) + 
			math.Pow(float64(y-(bounds.Max.Y-radius-1)), 2))
		return distance <= float64(radius)
	}

	// 不在四个角，保留原像素
	return true
}