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
)

// Direction 表示马赛克开始的方向
type Direction string

const (
	DirectionLeft   Direction = "left"   // 从左侧开始
	DirectionRight  Direction = "right"  // 从右侧开始
	DirectionTop    Direction = "top"    // 从顶部开始
	DirectionBottom Direction = "bottom" // 从底部开始
)

// MosaicRegion 表示一个需要添加马赛克的区域
type MosaicRegion struct {
	FromX int // 区域左上角X坐标
	FromY int // 区域左上角Y坐标
	ToX   int // 区域右下角X坐标
	ToY   int // 区域右下角Y坐标
}

// MosaicImageWithOptions 对图片指定区域添加马赛克效果，支持指定百分比和方向
// img: 原始图片字节数据
// regions: 需要添加马赛克的区域列表
// mosaicPercent: 马赛克区域百分比，范围0-1，表示要处理的区域比例
// startDirection: 开始马赛克的方向（left, right, top, bottom）
// 返回: 处理后的图片字节数据和错误信息
func MosaicImageWithOptions(img []byte, regions []*MosaicRegion, mosaicPercent float32, startDirection Direction) ([]byte, error) {
	// 使用新的处理器框架
	processor := NewMosaicProcessor(regions, mosaicPercent, startDirection)
	processors := []ImageProcessor{processor}

	// 处理图片
	return ProcessImage(img, processors, nil)
}

// calculateMosaicRegion 根据百分比和方向计算实际需要马赛克的区域
// fromX, fromY: 区域左上角坐标
// toX, toY: 区域右下角坐标
// percent: 马赛克区域百分比，范围0-1
// direction: 开始马赛克的方向
// 返回: 实际需要马赛克的区域坐标 (actualFromX, actualFromY, actualToX, actualToY)
func calculateMosaicRegion(fromX, fromY, toX, toY int, percent float32, direction Direction) (int, int, int, int) {
	// 计算区域宽度和高度
	width := toX - fromX
	height := toY - fromY

	// 根据方向和百分比计算实际区域
	switch direction {
	case DirectionLeft:
		// 从左侧开始，计算实际宽度
		actualWidth := int(float32(width) * percent)
		return fromX, fromY, fromX + actualWidth, toY

	case DirectionRight:
		// 从右侧开始，计算实际宽度
		actualWidth := int(float32(width) * percent)
		return toX - actualWidth, fromY, toX, toY

	case DirectionTop:
		// 从顶部开始，计算实际高度
		actualHeight := int(float32(height) * percent)
		return fromX, fromY, toX, fromY + actualHeight

	case DirectionBottom:
		// 从底部开始，计算实际高度
		actualHeight := int(float32(height) * percent)
		return fromX, toY - actualHeight, toX, toY

	default:
		// 默认处理整个区域
		return fromX, fromY, toX, toY
	}
}

// MosaicImage 对图片指定区域添加马赛克效果 (向后兼容版本)
// img: 原始图片字节数据
// regions: 需要添加马赛克的区域列表
// 返回: 处理后的图片字节数据和错误信息
func MosaicImage(img []byte, regions []*MosaicRegion) ([]byte, error) {
	return MosaicImageWithOptions(img, regions, 1.0, DirectionLeft)
}

// MosaicImageSingle 对图片单个区域添加马赛克效果 (向后兼容版本)
func MosaicImageSingle(img []byte, fromX, fromY, toX, toY int) ([]byte, error) {
	regions := []*MosaicRegion{
		{
			FromX: fromX,
			FromY: fromY,
			ToX:   toX,
			ToY:   toY,
		},
	}
	return MosaicImageWithOptions(img, regions, 1.0, DirectionLeft)
}

// MosaicImageSingleWithOptions 对图片单个区域添加马赛克效果，支持指定百分比和方向
func MosaicImageSingleWithOptions(img []byte, fromX, fromY, toX, toY int, mosaicPercent float32, startDirection Direction) ([]byte, error) {
	regions := []*MosaicRegion{
		{
			FromX: fromX,
			FromY: fromY,
			ToX:   toX,
			ToY:   toY,
		},
	}
	return MosaicImageWithOptions(img, regions, mosaicPercent, startDirection)
}

// clampUint8 确保值在0-255范围内
func clampUint8(value int16) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

// MosaicProcessor 马赛克处理器
type MosaicProcessor struct {
	Regions        []*MosaicRegion // 马赛克区域
	MosaicPercent  float32         // 马赛克区域百分比 (0-1)
	StartDirection Direction       // 开始方向
}

// Process 实现ImageProcessor接口
func (p *MosaicProcessor) Process(img image.Image) (image.Image, error) {
	// 获取图片边界
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 创建新的RGBA图像
	dstImg := image.NewRGBA(bounds)

	// 复制原图像到新图像
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dstImg.Set(x, y, img.At(x, y))
		}
	}

	// 处理每个马赛克区域
	for _, region := range p.Regions {
		// 验证坐标范围
		fromX := region.FromX
		fromY := region.FromY
		toX := region.ToX
		toY := region.ToY

		if fromX < 0 {
			fromX = 0
		}
		if fromY < 0 {
			fromY = 0
		}
		if toX > width {
			toX = width
		}
		if toY > height {
			toY = height
		}
		if fromX >= toX || fromY >= toY {
			// 跳过无效的区域
			continue
		}

		// 根据百分比和方向计算实际需要马赛克的区域
		actualFromX, actualFromY, actualToX, actualToY := calculateMosaicRegion(
			fromX, fromY, toX, toY, p.MosaicPercent, p.StartDirection)

		// 应用马赛克效果
		mosaicSize := 10 // 马赛克块大小
		if (actualToX-actualFromX)/10 > mosaicSize {
			mosaicSize = (actualToX - actualFromX) / 10
		}
		if (actualToY-actualFromY)/10 > mosaicSize {
			mosaicSize = (actualToY - actualFromY) / 10
		}

		for y := actualFromY; y < actualToY; y += mosaicSize {
			for x := actualFromX; x < actualToX; x += mosaicSize {
				// 计算当前块的边界
				blockEndX := x + mosaicSize
				blockEndY := y + mosaicSize
				if blockEndX > actualToX {
					blockEndX = actualToX
				}
				if blockEndY > actualToY {
					blockEndY = actualToY
				}

				// 计算块内像素的平均颜色
				var totalR, totalG, totalB, totalA uint32
				pixelCount := 0

				for blockY := y; blockY < blockEndY; blockY++ {
					for blockX := x; blockX < blockEndX; blockX++ {
						r, g, b, a := img.At(blockX, blockY).RGBA()
						totalR += r
						totalG += g
						totalB += b
						totalA += a
						pixelCount++
					}
				}

				// 计算平均颜色
				if pixelCount > 0 {
					// 计算原始平均颜色
					avgR := uint8(totalR / uint32(pixelCount) / 256)
					avgG := uint8(totalG / uint32(pixelCount) / 256)
					avgB := uint8(totalB / uint32(pixelCount) / 256)
					avgA := uint8(totalA / uint32(pixelCount) / 256)

					// 添加随机偏移量，防止去马赛克技术还原
					// 使用块的坐标作为随机种子，确保同一位置的偏移量一致
					randomSeed := (x*1103515245 + y*69069) & 0x7fffffff

					// 生成-10到10之间的随机偏移量
					offsetR := int8(randomSeed%21 - 10)
					offsetG := int8((randomSeed>>8)%21 - 10)
					offsetB := int8((randomSeed>>16)%21 - 10)

					// 应用偏移量并确保值在0-255范围内
					finalR := clampUint8(int16(avgR) + int16(offsetR))
					finalG := clampUint8(int16(avgG) + int16(offsetG))
					finalB := clampUint8(int16(avgB) + int16(offsetB))

					// 将带有随机偏移的颜色应用到整个块
					for blockY := y; blockY < blockEndY; blockY++ {
						for blockX := x; blockX < blockEndX; blockX++ {
							dstImg.SetRGBA(blockX, blockY, color.RGBA{R: finalR, G: finalG, B: finalB, A: avgA})
						}
					}
				}
			}
		}
	}

	return dstImg, nil
}

// NewMosaicProcessor 创建新的马赛克处理器
func NewMosaicProcessor(regions []*MosaicRegion, mosaicPercent float32, startDirection Direction) *MosaicProcessor {
	// 验证参数
	if mosaicPercent < 0 || mosaicPercent > 1 {
		mosaicPercent = 1.0 // 默认处理整个区域
	}

	return &MosaicProcessor{
		Regions:        regions,
		MosaicPercent:  mosaicPercent,
		StartDirection: startDirection,
	}
}
