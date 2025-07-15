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
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestOverlayProcessor 测试图层叠加处理器
func TestOverlayProcessor(t *testing.T) {
	// 创建底图
	baseImg := createTestImageForProcessor(400, 300)

	// 创建叠加图（一个红色方块）
	overlayImgData := createOverlayTestImage(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 200})

	// 解码叠加图像
	overlayImg, _, err := image.Decode(bytes.NewReader(overlayImgData))
	if err != nil {
		t.Fatalf("解码叠加图像失败: %v", err)
	}

	// 创建处理器链
	processors := []Processor{
		// 添加图层叠加处理器
		NewOverlayProcessor(overlayImg, 50, 50, 0.8, 1.0),
	}

	// 处理图片
	result, err := ProcessImage(baseImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	// 验证返回的是有效的图片数据
	_, _, err = image.Decode(bytes.NewReader(result))
	if err != nil {
		t.Fatalf("处理后的图片无法解码: %v", err)
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_overlay_processor.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_overlay_processor.png")
	}
}

// TestOverlayProcessorWithPosition 测试使用预设位置的图层叠加处理器
func TestOverlayProcessorWithPosition(t *testing.T) {
	// 创建底图
	baseImg := createTestImageForProcessor(400, 300)

	// 创建叠加图（一个蓝色方块）
	overlayImgData := createOverlayTestImage(80, 80, color.RGBA{R: 0, G: 0, B: 255, A: 180})

	// 解码叠加图像
	overlayImg, _, err := image.Decode(bytes.NewReader(overlayImgData))
	if err != nil {
		t.Fatalf("解码叠加图像失败: %v", err)
	}

	// 测试不同位置
	positions := []string{"center", "top-left", "top-right", "bottom-left", "bottom-right", "top-center", "bottom-center"}

	for _, position := range positions {
		// 创建处理器链
		processors := []Processor{
			// 添加图层叠加处理器
			NewOverlayProcessorWithPosition(overlayImg, position, 0.9, 1.0),
		}

		// 处理图片
		result, err := ProcessImage(baseImg, processors, nil)
		if err != nil {
			t.Fatalf("图片处理失败（位置 %s）: %v", position, err)
		}

		// 保存测试图片
		filename := "/tmp/test_overlay_" + position + ".png"
		err = os.WriteFile(filename, result, 0o644)
		if err != nil {
			t.Logf("Warning: Could not save test image: %v", err)
		} else {
			t.Logf("处理后的图片已保存到: %s", filename)
		}
	}
}

// TestOverlayProcessorWithScale 测试缩放功能
func TestOverlayProcessorWithScale(t *testing.T) {
	// 创建底图
	baseImg := createTestImageForProcessor(400, 300)

	// 创建叠加图（一个绿色方块）
	overlayImgData := createOverlayTestImage(100, 100, color.RGBA{R: 0, G: 255, B: 0, A: 200})

	// 解码叠加图像
	overlayImg, _, err := image.Decode(bytes.NewReader(overlayImgData))
	if err != nil {
		t.Fatalf("解码叠加图像失败: %v", err)
	}

	// 测试不同缩放比例
	scales := []float64{0.5, 1.0, 1.5, 2.0}

	for _, scale := range scales {
		// 创建处理器链
		processors := []Processor{
			// 添加图层叠加处理器
			NewOverlayProcessor(overlayImg, 150, 100, 1.0, scale),
		}

		// 处理图片
		result, err := ProcessImage(baseImg, processors, nil)
		if err != nil {
			t.Fatalf("图片处理失败（缩放比例 %.1f）: %v", scale, err)
		}

		// 保存测试图片
		filename := "/tmp/test_overlay_scale_" + fmt.Sprintf("%.1f", scale) + ".png"
		err = os.WriteFile(filename, result, 0o644)
		if err != nil {
			t.Logf("Warning: Could not save test image: %v", err)
		} else {
			t.Logf("处理后的图片已保存到: %s", filename)
		}
	}
}

// TestMultipleOverlays 测试多个叠加图层
func TestMultipleOverlays(t *testing.T) {
	// 创建底图
	baseImg := createTestImageForProcessor(400, 300)

	// 创建多个叠加图
	overlayImgData1 := createOverlayTestImage(120, 120, color.RGBA{R: 255, G: 0, B: 0, A: 150})
	overlayImgData2 := createOverlayTestImage(80, 80, color.RGBA{R: 0, G: 255, B: 0, A: 150})
	overlayImgData3 := createOverlayTestImage(60, 60, color.RGBA{R: 0, G: 0, B: 255, A: 150})

	// 解码叠加图像
	overlayImg1, _, err := image.Decode(bytes.NewReader(overlayImgData1))
	if err != nil {
		t.Fatalf("解码叠加图像1失败: %v", err)
	}
	overlayImg2, _, err := image.Decode(bytes.NewReader(overlayImgData2))
	if err != nil {
		t.Fatalf("解码叠加图像2失败: %v", err)
	}
	overlayImg3, _, err := image.Decode(bytes.NewReader(overlayImgData3))
	if err != nil {
		t.Fatalf("解码叠加图像3失败: %v", err)
	}

	// 创建处理器链
	processors := []Processor{
		// 添加多个图层叠加处理器
		NewOverlayProcessorWithPosition(overlayImg1, "top-left", 0.8, 1.0),
		NewOverlayProcessorWithPosition(overlayImg2, "bottom-right", 0.8, 1.0),
		NewOverlayProcessorWithPosition(overlayImg3, "center", 0.8, 1.0),
	}

	// 处理图片
	result, err := ProcessImage(baseImg, processors, nil)
	if err != nil {
		t.Fatalf("图片处理失败: %v", err)
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_multiple_overlays.png", result, 0o644)
	if err != nil {
		t.Logf("Warning: Could not save test image: %v", err)
	} else {
		t.Logf("处理后的图片已保存到: /tmp/test_multiple_overlays.png")
	}
}

// 创建测试用的叠加图像
func createOverlayTestImage(width, height int, bgColor color.RGBA) []byte {
	// 创建一个彩色图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充颜色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 编码为PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}
