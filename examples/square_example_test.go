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
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/vogo/vimage"
)

func TestSquareProcessorLocalFile(t *testing.T) {
	// 读取本地文件进行测试
	b, err := os.ReadFile("build/avatar.jpg")
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
		outputFile := "../build/avatar_square_" + pos + ".jpg"
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
