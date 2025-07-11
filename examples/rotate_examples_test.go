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
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/vogo/vimage"
)

// Example_rotateImage 展示如何使用RotateProcessor旋转图像
func TestRotateImage(t *testing.T) {
	// 创建一个测试图像
	img := createColorBlockImage(200, 100)

	// 创建一个旋转处理器，旋转45度
	rotateProcessor := vimage.NewRotateProcessor(45)

	// 处理图像
	result, err := rotateProcessor.Process(img)
	if err != nil {
		fmt.Printf("旋转图像失败: %v\n", err)
		return
	}

	// 保存结果
	outputPath := "/tmp/output_rotate_45.png"
	if err := saveImage(result, outputPath); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("图像已旋转45度并保存到: %s\n", outputPath)

	// 使用自定义背景色旋转90度
	rotateProcessor = vimage.NewRotateProcessor(90).WithBackground(color.RGBA{255, 255, 0, 255}) // 黄色背景
	result, err = rotateProcessor.Process(img)
	if err != nil {
		fmt.Printf("旋转图像失败: %v\n", err)
		return
	}

	// 保存结果
	outputPath = "/tmp/output_rotate_90_yellow_bg.png"
	if err := saveImage(result, outputPath); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("图像已旋转90度（黄色背景）并保存到: %s\n", outputPath)

	// 保持原始尺寸旋转180度
	rotateProcessor = vimage.NewRotateProcessor(180).WithKeepSize(true)
	result, err = rotateProcessor.Process(img)
	if err != nil {
		fmt.Printf("旋转图像失败: %v\n", err)
		return
	}

	// 保存结果
	outputPath = "/tmp/output_rotate_180_keep_size.png"
	if err := saveImage(result, outputPath); err != nil {
		fmt.Printf("保存图像失败: %v\n", err)
		return
	}

	fmt.Printf("图像已旋转180度（保持原始尺寸）并保存到: %s\n", outputPath)

	// Output:
	// 图像已旋转45度并保存到: output_rotate_45.png
	// 图像已旋转90度（黄色背景）并保存到: output_rotate_90_yellow_bg.png
	// 图像已旋转180度（保持原始尺寸）并保存到: output_rotate_180_keep_size.png
}

// 创建一个彩色块测试图像
func createColorBlockImage(width, height int) image.Image {
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

// 保存图像到文件
func saveImage(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 根据文件扩展名选择编码格式
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case ".png":
		return png.Encode(file, img)
	default:
		// 默认使用PNG格式
		return png.Encode(file, img)
	}
}
