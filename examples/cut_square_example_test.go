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
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/vogo/vimage"
)

func TestCutSquareProcessor(t *testing.T) {
	// 构造一张非正方形的测试图片，避免依赖本地文件
	src := image.NewRGBA(image.Rect(0, 0, 300, 200))
	// 填充背景
	draw.Draw(src, src.Bounds(), &image.Uniform{C: color.RGBA{240, 240, 240, 255}}, image.Point{}, draw.Src)
	// 画出左右/上下不同色块，便于人工查看裁剪结果
	for y := 0; y < 200; y++ {
		for x := 0; x < 300; x++ {
			switch {
			case x < 50:
				src.Set(x, y, color.RGBA{255, 0, 0, 255})
			case x >= 250:
				src.Set(x, y, color.RGBA{0, 0, 255, 255})
			case y < 30:
				src.Set(x, y, color.RGBA{0, 255, 0, 255})
			case y >= 170:
				src.Set(x, y, color.RGBA{255, 255, 0, 255})
			}
		}
	}

	// 测试不同裁剪位置
	positions := []string{"center", "top", "bottom", "left", "right"}

	for _, pos := range positions {
		processor := vimage.NewCutSquareProcessor(pos)
		result, err := processor.Process(src)
		if err != nil {
			t.Fatalf("Process failed with position %s: %v", pos, err)
		}

		// 验证结果是正方形
		bounds := result.Bounds()
		if bounds.Dx() != bounds.Dy() {
			t.Errorf("Result should be square with position %s, got %dx%d", pos, bounds.Dx(), bounds.Dy())
		}

		tmp := t.TempDir()
		outputFile := filepath.Join(tmp, "avatar_square_"+pos+".jpg")
		f, err := os.Create(outputFile)
		if err != nil {
			t.Logf("Create failed for position %s: %v", pos, err)
			continue
		}
		defer func() { _ = f.Close() }()

		if err = jpeg.Encode(f, result, &jpeg.Options{Quality: 90}); err != nil {
			t.Logf("Encode failed for position %s: %v", pos, err)
		}
	}
}
