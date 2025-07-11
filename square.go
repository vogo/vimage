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
)

// SquareProcessor 正方形裁剪处理器
// 计算图片长宽，裁剪掉更长或更高的部分，使图片变为正方形
type SquareProcessor struct {
	// 可以添加裁剪位置的选项，如居中裁剪、从顶部裁剪等
	Position string // 裁剪位置："center"(默认), "top", "bottom", "left", "right"
}

// Process 实现ImageProcessor接口
func (p *SquareProcessor) Process(img image.Image) (image.Image, error) {
	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 如果已经是正方形，直接返回
	if width == height {
		return img, nil
	}

	// 确定目标尺寸（取较小的一边）
	size := width
	if height < width {
		size = height
	}

	// 计算裁剪的起始位置
	var x, y int
	switch p.Position {
	case "top":
		// 从顶部开始，水平居中
		x = (width - size) / 2
		y = 0
	case "bottom":
		// 从底部开始，水平居中
		x = (width - size) / 2
		y = height - size
	case "left":
		// 从左侧开始，垂直居中
		x = 0
		y = (height - size) / 2
	case "right":
		// 从右侧开始，垂直居中
		x = width - size
		y = (height - size) / 2
	default: // "center" 或其他值
		// 居中裁剪
		x = (width - size) / 2
		y = (height - size) / 2
	}

	// 创建子图像
	squareImg := image.NewRGBA(image.Rect(0, 0, size, size))

	// 复制像素
	for dy := 0; dy < size; dy++ {
		for dx := 0; dx < size; dx++ {
			squareImg.Set(dx, dy, img.At(x+dx, y+dy))
		}
	}

	return squareImg, nil
}

// NewSquareProcessor 创建新的正方形裁剪处理器
func NewSquareProcessor(position string) *SquareProcessor {
	// 如果位置参数无效，使用默认值
	validPositions := map[string]bool{
		"center": true,
		"top":    true,
		"bottom": true,
		"left":   true,
		"right":  true,
	}

	if !validPositions[position] {
		position = "center" // 默认居中裁剪
	}

	return &SquareProcessor{
		Position: position,
	}
}
