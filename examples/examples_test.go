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

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/vogo/vimage"
)

func TestResizeLocalFile(t *testing.T) {
	// Create a test file
	b, err := os.ReadFile("/tmp/avatar.jpg")
	if err != nil {
		t.Skipf("Create failed: %v", err)
	}

	// Decode the image
	img, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Skipf("Decode failed: %v", err)
	}
	fmt.Println(format)

	processor := vimage.NewResizeProcessor(100, 100)
	img, err = processor.Process(img)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	outputFile := "/tmp/avatar_resized.jpg"
	os.Remove(outputFile)
	f, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	defer f.Close()
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 90})
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
}

func TestSquareProcessorLocalFile(t *testing.T) {
	// 读取本地文件进行测试
	b, err := os.ReadFile("/tmp/avatar.jpg")
	if err != nil {
		t.Skipf("ReadFile failed: %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Skipf("Decode failed: %v", err)
	}

	// 测试不同裁剪位置
	positions := []string{"center", "top", "bottom", "left", "right"}

	for _, pos := range positions {
		processor := vimage.NewSquareProcessor(pos)
		result, err := processor.Process(img)
		if err != nil {
			t.Fatalf("Process failed with position %s: %v", pos, err)
		}

		// 验证结果是正方形
		bounds := result.Bounds()
		if bounds.Dx() != bounds.Dy() {
			t.Errorf("Result should be square with position %s, got %dx%d", pos, bounds.Dx(), bounds.Dy())
		}

		// 保存处理后的图片
		outputFile := "/tmp/avatar_square_" + pos + ".jpg"
		os.Remove(outputFile)
		f, err := os.Create(outputFile)
		if err != nil {
			t.Logf("Create failed for position %s: %v", pos, err)
			continue
		}

		err = jpeg.Encode(f, result, &jpeg.Options{Quality: 90})
		f.Close()
		if err != nil {
			t.Logf("Encode failed for position %s: %v", pos, err)
		}
	}
}

func TestMosaicLocalImage(t *testing.T) {
	testImg, err := os.ReadFile("/tmp/test_cert.jpeg")
	if err != nil {
		t.Skipf("读取测试图片失败: %v", err)
	}

	// 使用向后兼容函数
	result, err := vimage.MosaicImageSingle(testImg, 683, 355, 872, 380)
	if err != nil {
		t.Fatalf("马赛克处理失败: %v", err)
	}
	if err := os.WriteFile("/tmp/test_cert_mosaic.jpeg", result, 0o644); err != nil {
		t.Fatalf("保存马赛克图片失败: %v", err)
	}
}

// SquareImage 将图片裁剪为正方形
// imgData: 原始图片字节数据
// position: 裁剪位置 ("center", "top", "bottom", "left", "right")
// 返回: 处理后的图片字节数据和错误信息
func SquareImage(imgData []byte, position string) ([]byte, error) {
	// 创建处理器链
	processors := []vimage.ImageProcessor{
		// 添加正方形裁剪处理器
		vimage.NewSquareProcessor(position),
	}

	// 处理图片
	return vimage.ProcessImage(imgData, processors, nil)
}

// SquareAndResizeImage 将图片裁剪为正方形并调整大小
// imgData: 原始图片字节数据
// position: 裁剪位置 ("center", "top", "bottom", "left", "right")
// size: 目标尺寸（正方形的边长）
// 返回: 处理后的图片字节数据和错误信息
func SquareAndResizeImage(imgData []byte, position string, size int) ([]byte, error) {
	// 创建处理器链
	processors := []vimage.ImageProcessor{
		// 先裁剪为正方形
		vimage.NewSquareProcessor(position),
		// 再调整大小
		vimage.NewResizeProcessor(size, size),
	}

	// 处理图片
	return vimage.ProcessImage(imgData, processors, nil)
}

// SquareAndCircleImage 将图片裁剪为正方形并应用圆形裁剪
// imgData: 原始图片字节数据
// position: 裁剪位置 ("center", "top", "bottom", "left", "right")
// 返回: 处理后的图片字节数据和错误信息
func SquareAndCircleImage(imgData []byte, position string) ([]byte, error) {
	// 解码图片
	srcImg, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, err
	}

	// 先裁剪为正方形
	squareProcessor := vimage.NewSquareProcessor(position)
	squareImg, err := squareProcessor.Process(srcImg)
	if err != nil {
		return nil, err
	}

	// 再应用圆形裁剪
	circleProcessor := &vimage.CircleProcessor{}
	circleImg, err := circleProcessor.Process(squareImg)
	if err != nil {
		return nil, err
	}

	// 编码图片
	buf := new(bytes.Buffer)
	options := &vimage.DefaultProcessorOptions

	// 根据原始格式编码
	switch format {
	case "jpeg":
		err = jpeg.Encode(buf, circleImg, &jpeg.Options{Quality: options.Quality})
	case "png":
		err = png.Encode(buf, circleImg)
	default:
		// 默认使用PNG格式
		err = png.Encode(buf, circleImg)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
