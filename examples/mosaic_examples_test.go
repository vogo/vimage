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
	"os"
	"testing"

	"github.com/vogo/vimage"
)

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
		t.Logf("保存马赛克图片失败: %v", err)
	}
}
