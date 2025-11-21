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
	"os"
	"testing"

	"github.com/vogo/vimage"
)

// TestCutProcessor 测试切割处理器
func TestCutProcessor(t *testing.T) {
	// 打开测试图片
	imgData, err := os.ReadFile("/tmp/avatar.jpg")
	if err != nil {
		t.Skipf("读取测试图片失败: %v", err)
	}

	// 创建切割处理器 - 居中切割
	processors := []vimage.Processor{
		vimage.NewCutProcessor(100, 100, vimage.CutPositionCenter),
	}

	// 处理图片
	resultData, err := vimage.ProcessImage(imgData, processors, nil)
	if err != nil {
		t.Fatalf("处理图片失败: %v", err)
	}

	// 保存结果
	err = os.WriteFile("/tmp/cut_center.jpg", resultData, 0o644)
	if err != nil {
		t.Logf("保存结果图片失败: %v", err)
	}

	fmt.Println("居中切割结果已保存到: /tmp/cut_center.jpg")

	// 创建切割处理器 - 自定义区域切割
	processors = []vimage.Processor{
		vimage.NewCutProcessorWithRegion(100, 100, 50, 50),
	}

	// 处理图片
	resultData, err = vimage.ProcessImage(imgData, processors, nil)
	if err != nil {
		t.Fatalf("处理图片失败: %v", err)
	}

	// 保存结果
	err = os.WriteFile("/tmp/cut_custom.jpg", resultData, 0o644)
	if err != nil {
		t.Logf("保存结果图片失败: %v", err)
	}

	fmt.Println("自定义区域切割结果已保存到: /tmp/cut_custom.jpg")
}

// TestZoomProcessor 测试缩放处理器
func TestZoomProcessor(t *testing.T) {
	// 打开测试图片
	imgData, err := os.ReadFile("/tmp/avatar.jpg")
	if err != nil {
		t.Skipf("读取测试图片失败: %v", err)
	}

	// 测试不同的缩放模式
	testCases := []struct {
		name      string
		processor vimage.Processor
	}{
		{
			name:      "exact",
			processor: vimage.NewZoomProcessor(400, 300),
		},
		{
			name:      "ratio",
			processor: vimage.NewZoomRatioProcessor(0.5),
		},
		{
			name:      "width",
			processor: vimage.NewZoomWidthProcessor(400),
		},
		{
			name:      "height",
			processor: vimage.NewZoomHeightProcessor(300),
		},
		{
			name:      "max",
			processor: vimage.NewZoomMaxProcessor(400),
		},
		{
			name:      "min",
			processor: vimage.NewZoomMinProcessor(300),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 处理图片
			processors := []vimage.Processor{tc.processor}
			resultData, err := vimage.ProcessImage(imgData, processors, nil)
			if err != nil {
				t.Fatalf("处理图片失败: %v", err)
			}

			// 保存结果
			outputPath := fmt.Sprintf("/tmp/zoom_%s.jpg", tc.name)
			err = os.WriteFile(outputPath, resultData, 0o644)
			if err != nil {
				t.Logf("保存结果图片失败: %v", err)
			}

			fmt.Printf("%s缩放结果已保存到: %s\n", tc.name, outputPath)
		})
	}
}

// TestCombinedCutAndZoom 测试组合使用切割和缩放
func TestCombinedCutAndZoom(t *testing.T) {
	// 打开测试图片
	imgData, err := os.ReadFile("/tmp/avatar.jpg")
	if err != nil {
		t.Skipf("读取测试图片失败: %v", err)
	}

	// 先切割再缩放
	processors := []vimage.Processor{
		vimage.NewCutProcessor(100, 100, vimage.CutPositionCenter),
		vimage.NewZoomRatioProcessor(0.5),
	}

	// 处理图片
	resultData, err := vimage.ProcessImage(imgData, processors, nil)
	if err != nil {
		t.Fatalf("处理图片失败: %v", err)
	}

	// 保存结果
	err = os.WriteFile("/tmp/cut_then_zoom.jpg", resultData, 0o644)
	if err != nil {
		t.Logf("保存结果图片失败: %v", err)
	}

	fmt.Println("先切割再缩放结果已保存到: /tmp/cut_then_zoom.jpg")
}
