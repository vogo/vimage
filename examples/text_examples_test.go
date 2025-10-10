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
	"image"
	"image/color"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/vogo/vimage"
	"golang.org/x/image/font/basicfont"
)

// TestRotatedText 展示如何使用TextProcessor旋转文本
func TestRotatedText(t *testing.T) {
	// 创建一个测试图像
	img := createBackgroundImage(400, 300)

	// 创建一个文本处理器，添加普通文本
	textProcessor := vimage.NewTextProcessor(vimage.TextOptions{
		Text:     "普通文本",
		Position: image.Point{50, 50},
		Font:     basicfont.Face7x13,
		Color:    color.RGBA{255, 0, 0, 255}, // 红色
	})

	// 处理图像
	result, err := textProcessor.Process(img)
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	// 创建一个文本处理器，添加旋转45度的文本
	textProcessor = vimage.NewTextProcessor(vimage.TextOptions{
		Text:     "旋转45度",
		Position: image.Point{200, 100},
		Font:     basicfont.Face7x13,
		Color:    color.RGBA{0, 0, 255, 255}, // 蓝色
		Angle:    45,                         // 旋转45度
	})

	// 处理图像
	result, err = textProcessor.Process(result)
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	// 使用WithAngle方法添加旋转90度的文本
	textProcessor = vimage.NewTextProcessor(vimage.TextOptions{
		Text:     "旋转90度",
		Position: image.Point{300, 150},
		Font:     basicfont.Face7x13,
		Color:    color.RGBA{0, 255, 0, 255}, // 绿色
	}).WithAngle(90) // 旋转90度

	// 处理图像
	result, err = textProcessor.Process(result)
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	// 保存结果
	outputPath := "../build/output_rotated_text.png"
	if err := saveImage(result, outputPath); err != nil {
		t.Logf("保存图像失败: %v", err)
	}

	t.Logf("旋转文本示例已保存到: %s", outputPath)
}

// TestTextWithProcessorChain 展示如何将文本处理器与其他处理器组合使用
func TestTextWithProcessorChain(t *testing.T) {
	// 创建一个测试图像
	img := createBackgroundImage(400, 300)

	// 创建处理器链
	processors := []vimage.Processor{
		// 先添加一个旋转45度的文本
		vimage.NewTextProcessor(vimage.TextOptions{
			Text:     "先添加文本",
			Position: image.Point{100, 100},
			Font:     basicfont.Face7x13,
			Color:    color.RGBA{255, 0, 0, 255}, // 红色
			Angle:    45,                         // 旋转45度
		}),
		// 然后旋转整个图像
		vimage.NewRotateProcessor(30).WithBackground(color.RGBA{240, 240, 240, 255}),
		// 最后添加一个水平文本
		vimage.NewTextProcessor(vimage.TextOptions{
			Text:     "再添加文本",
			Position: image.Point{150, 150},
			Font:     basicfont.Face7x13,
			Color:    color.RGBA{0, 0, 255, 255}, // 蓝色
		}),
	}

	// 使用处理器链处理图像
	result, err := vimage.Process(img, processors)
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	// 保存结果
	outputPath := "../build/output_text_with_processor_chain.png"
	if err := saveImage(result, outputPath); err != nil {
		t.Logf("保存图像失败: %v", err)
	}

	t.Logf("文本处理器链示例已保存到: %s", outputPath)
}

// TestWrappedText 展示限制文本宽度并自动换行
func TestWrappedText(t *testing.T) {
	// 创建一个测试图像
	img := createBackgroundImage(500, 300)

	longText := "这是一段较长的文本，用于演示在给定的最大宽度内自动换行的效果。你可以调整最大宽度来观察不同的换行行为。"

	// 创建文本处理器，限制宽度并自动换行
	textProcessor := vimage.NewTextProcessor(vimage.TextOptions{
		Text:        longText,
		Position:    image.Point{50, 60},
		Font:        truetype.NewFace(vimage.LoadHarmonyOSSansSCBlack(), &truetype.Options{Size: 16}),
		Color:       color.Black,
		MaxWidth:    200,  // 限制最大宽度为200像素
		LineSpacing: 1.5,  // 行距倍数
		CharWrap:    true, // 启用按字符换行，适合中文
		// Align 默认为左对齐
	})

	// 处理图像
	result, err := textProcessor.Process(img)
	if err != nil {
		t.Fatalf("处理图像失败: %v", err)
	}

	// 保存结果
	outputPath := "../build/output_wrapped_text.png"
	if err := saveImage(result, outputPath); err != nil {
		t.Logf("保存图像失败: %v", err)
	}

	t.Logf("限制宽度自动换行示例已保存到: %s", outputPath)
}

// 创建一个背景图像
func createBackgroundImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充浅灰色背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{240, 240, 240, 255})
		}
	}

	return img
}
