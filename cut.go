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
)

// CutPosition 定义切割位置
type CutPosition string

const (
	// CutPositionCenter 居中切割
	CutPositionCenter CutPosition = "center"
	// CutPositionTop 从顶部切割
	CutPositionTop CutPosition = "top"
	// CutPositionBottom 从底部切割
	CutPositionBottom CutPosition = "bottom"
	// CutPositionLeft 从左侧切割
	CutPositionLeft CutPosition = "left"
	// CutPositionRight 从右侧切割
	CutPositionRight CutPosition = "right"
)

// CutProcessor 图像切割处理器
// 从原始图像中切割出指定区域（支持矩形和正方形）
type CutProcessor struct {
	// 目标宽度和高度（0表示自动检测）
	Width  int
	Height int
	// 切割位置
	Position CutPosition
	// 自定义切割区域（如果指定了，则忽略Position）
	X int // 左上角X坐标
	Y int // 左上角Y坐标
	// 是否使用自定义区域
	UseCustomRegion bool
	// 是否为正方形模式（自动使用较小边）
	SquareMode bool
}

// Process 实现Processor接口
func (p *CutProcessor) Process(img image.Image) (image.Image, error) {
	// 获取原始图片尺寸
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// 确定目标尺寸
	width := p.Width
	height := p.Height

	if p.SquareMode {
		// 正方形模式：如果没有指定尺寸，使用较小边
		if width == 0 && height == 0 {
			size := min(origWidth, origHeight)
			width = size
			height = size
		} else if width > 0 && height == 0 {
			// 指定了宽度，高度使用相同值
			height = width
		} else if width == 0 && height > 0 {
			// 指定了高度，宽度使用相同值
			width = height
		}
		// 如果宽高都指定了，确保它们相等
		if width != height {
			return nil, fmt.Errorf("正方形模式下宽度和高度必须相等: %dx%d", width, height)
		}

		// 如果已经是正方形且尺寸匹配，直接返回
		if origWidth == origHeight && origWidth == width && !p.UseCustomRegion {
			return img, nil
		}
	}

	// 验证目标尺寸
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("无效的切割尺寸: %dx%d", width, height)
	}

	// 检查目标尺寸是否超过原始尺寸
	if width > origWidth || height > origHeight {
		return nil, fmt.Errorf("切割尺寸(%dx%d)超过原始尺寸(%dx%d)",
			width, height, origWidth, origHeight)
	}

	// 计算切割的起始位置
	var x, y int

	if p.UseCustomRegion {
		// 使用自定义区域
		x = p.X
		y = p.Y

		// 验证自定义区域是否有效
		if x < 0 || y < 0 || x+width > origWidth || y+height > origHeight {
			return nil, fmt.Errorf("无效的切割区域: 起点(%d,%d), 尺寸(%dx%d), 原始尺寸(%dx%d)",
				x, y, width, height, origWidth, origHeight)
		}
	} else {
		// 根据位置计算起始点
		switch p.Position {
		case CutPositionTop:
			// 从顶部开始，水平居中
			x = (origWidth - width) / 2
			y = 0
		case CutPositionBottom:
			// 从底部开始，水平居中
			x = (origWidth - width) / 2
			y = origHeight - height
		case CutPositionLeft:
			// 从左侧开始，垂直居中
			x = 0
			y = (origHeight - height) / 2
		case CutPositionRight:
			// 从右侧开始，垂直居中
			x = origWidth - width
			y = (origHeight - height) / 2
		default: // CutPositionCenter 或其他值
			// 居中切割
			x = (origWidth - width) / 2
			y = (origHeight - height) / 2
		}
	}

	// 尝试使用SubImage以提高性能
	subImg, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if ok {
		return subImg.SubImage(image.Rect(x, y, x+width, y+height)), nil
	}

	// 如果不支持SubImage，手动复制像素
	cutImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			cutImg.Set(dx, dy, img.At(x+dx, y+dy))
		}
	}

	return cutImg, nil
}

// NewCutProcessor 创建新的矩形切割处理器（使用预定义位置）
func NewCutProcessor(width, height int, position CutPosition) *CutProcessor {
	return &CutProcessor{
		Width:           width,
		Height:          height,
		Position:        position,
		UseCustomRegion: false,
		SquareMode:      false,
	}
}

// NewCutProcessorWithRegion 创建新的切割处理器（使用自定义区域）
func NewCutProcessorWithRegion(width, height, x, y int) *CutProcessor {
	return &CutProcessor{
		Width:           width,
		Height:          height,
		X:               x,
		Y:               y,
		UseCustomRegion: true,
		SquareMode:      false,
	}
}

// NewSquareCutProcessor 创建正方形切割处理器
func NewSquareCutProcessor(size int, position CutPosition) *CutProcessor {
	return NewCutProcessor(size, size, position)
}

// NewCutSquareProcessor 创建正方形切割处理器（自动使用较小边）
func NewCutSquareProcessor(position string) *CutProcessor {
	return &CutProcessor{
		Width:           0, // 自动检测
		Height:          0, // 自动检测
		Position:        CutPosition(position),
		UseCustomRegion: false,
		SquareMode:      true,
	}
}

// NewCutSquareProcessorWithSize 创建指定尺寸的正方形切割处理器
func NewCutSquareProcessorWithSize(size int, position string) *CutProcessor {
	return &CutProcessor{
		Width:           size,
		Height:          size,
		Position:        CutPosition(position),
		UseCustomRegion: false,
		SquareMode:      true,
	}
}

// NewCutSquareProcessorWithRegion 创建使用自定义坐标的正方形切割处理器
func NewCutSquareProcessorWithRegion(size, x, y int) *CutProcessor {
	return &CutProcessor{
		Width:           size,
		Height:          size,
		X:               x,
		Y:               y,
		UseCustomRegion: true,
		SquareMode:      true,
	}
}

// Type alias for backward compatibility
type CutSquareProcessor = CutProcessor
