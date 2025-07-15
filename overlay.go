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
	"errors"
	"image"
	"image/color"

	"github.com/fogleman/gg"
)

// OverlayProcessor 图层叠加处理器
type OverlayProcessor struct {
	OverlayImage image.Image // 叠加图像
	X            int         // 叠加位置X坐标
	Y            int         // 叠加位置Y坐标
	Opacity      float64     // 不透明度 (0-1)
	Scale        float64     // 缩放比例 (0-n)
	Position     string      // 预设位置 ("center", "top-left", "bottom-right" 等)
}

// Process 实现Processor接口
func (p *OverlayProcessor) Process(img image.Image) (image.Image, error) {
	// 获取底图边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建新的上下文
	dc := gg.NewContext(width, height)

	// 绘制原图
	dc.DrawImage(img, 0, 0)

	// 使用叠加图像
	if p.OverlayImage == nil {
		return nil, errors.New("未提供叠加图像")
	}

	// 获取叠加图像
	overlayImg := p.OverlayImage

	// 计算叠加图像的尺寸
	overlayBounds := overlayImg.Bounds()
	overlayWidth := float64(overlayBounds.Dx())
	overlayHeight := float64(overlayBounds.Dy())

	// 应用缩放
	if p.Scale != 0 && p.Scale != 1 {
		newWidth := int(overlayWidth * p.Scale)
		newHeight := int(overlayHeight * p.Scale)

		// 创建临时上下文进行缩放
		tempDc := gg.NewContext(newWidth, newHeight)
		tempDc.DrawImage(overlayImg, 0, 0)
		overlayImg = tempDc.Image()

		// 更新尺寸
		overlayBounds = overlayImg.Bounds()
		overlayWidth = float64(overlayBounds.Dx())
		overlayHeight = float64(overlayBounds.Dy())
	}

	// 计算叠加位置
	var x, y float64

	switch p.Position {
	case "top-left":
		x, y = 0, 0
	case "top-right":
		x, y = float64(width)-overlayWidth, 0
	case "bottom-left":
		x, y = 0, float64(height)-overlayHeight
	case "bottom-right":
		x, y = float64(width)-overlayWidth, float64(height)-overlayHeight
	case "center":
		x, y = float64(width)/2-overlayWidth/2, float64(height)/2-overlayHeight/2
	case "top-center":
		x, y = float64(width)/2-overlayWidth/2, 0
	case "bottom-center":
		x, y = float64(width)/2-overlayWidth/2, float64(height)-overlayHeight
	case "left-center":
		x, y = 0, float64(height)/2-overlayHeight/2
	case "right-center":
		x, y = float64(width)-overlayWidth, float64(height)/2-overlayHeight/2
	default:
		// 使用指定的坐标
		x, y = float64(p.X), float64(p.Y)
	}

	// 处理透明度
	// 在 gg 库中，没有直接设置图像透明度的方法
	// 我们可以通过创建一个新的 RGBA 图像并调整每个像素的 alpha 值来实现
	if p.Opacity > 0 && p.Opacity < 1 {
		// 获取叠加图像的边界
		overlayBounds := overlayImg.Bounds()

		// 创建一个新的 RGBA 图像
		adjustedImg := image.NewRGBA(overlayBounds)

		// 遍历每个像素，调整 alpha 值
		for y := overlayBounds.Min.Y; y < overlayBounds.Max.Y; y++ {
			for x := overlayBounds.Min.X; x < overlayBounds.Max.X; x++ {
				// 获取原始颜色
				c := overlayImg.At(x, y)
				r, g, b, a := c.RGBA()

				// 调整 alpha 值
				a = uint32(float64(a) * p.Opacity)

				// 设置新颜色
				adjustedImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
			}
		}

		// 更新叠加图像
		overlayImg = adjustedImg
	}

	// 绘制叠加图像
	dc.DrawImage(overlayImg, int(x), int(y))

	return dc.Image(), nil
}

// NewOverlayProcessor 创建新的图层叠加处理器
func NewOverlayProcessor(overlayImage image.Image, x, y int, opacity, scale float64) *OverlayProcessor {
	// 验证参数
	if opacity < 0 || opacity > 1 {
		opacity = 1.0 // 默认完全不透明
	}
	if scale <= 0 {
		scale = 1.0 // 默认不缩放
	}

	return &OverlayProcessor{
		OverlayImage: overlayImage,
		X:            x,
		Y:            y,
		Opacity:      opacity,
		Scale:        scale,
	}
}

// NewOverlayProcessorWithPosition 创建新的图层叠加处理器（使用预设位置）
func NewOverlayProcessorWithPosition(overlayImage image.Image, position string, opacity, scale float64) *OverlayProcessor {
	// 验证参数
	if opacity < 0 || opacity > 1 {
		opacity = 1.0 // 默认完全不透明
	}
	if scale <= 0 {
		scale = 1.0 // 默认不缩放
	}

	return &OverlayProcessor{
		OverlayImage: overlayImage,
		Position:     position,
		Opacity:      opacity,
		Scale:        scale,
	}
}
