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
	"os"
	"testing"
)

func TestGenTableImage_Basic(t *testing.T) {
	headers := []string{"姓名", "年龄", "职业"}
	rows := [][]string{
		{"张三", "25", "工程师"},
		{"李四", "30", "设计师"},
		{"王五", "28", "产品经理"},
	}

	buf, err := GenMultipleRowsTableImage(nil, headers, rows, nil)
	if err != nil {
		t.Fatalf("GenTableImage failed: %v", err)
	}

	if buf == nil {
		t.Fatal("Expected non-nil buffer")
	}

	if buf.Len() == 0 {
		t.Fatal("Expected non-empty buffer")
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_basic_table.png", buf.Bytes(), 0o644)
	if err != nil {
		t.Logf("warning: failed to save basic table test image | output_path: /tmp/test_basic_table.png | err: %v", err)
	} else {
		t.Logf("basic table test image saved successfully | output_path: /tmp/test_basic_table.png")
	}

	t.Logf("basic table generated successfully | image_size_bytes: %d", buf.Len())
}

func TestGenMultipleColumnsTableImage_Basic(t *testing.T) {
	headers := []string{"姓名", "年龄", "职业"}
	data := [][]string{
		{"张三", "25", "工程师"},
		{"李四", "30", "设计师"},
		{"王五", "28", "产品经理"},
		{"赵六", "32", "架构师"},
	}

	buf, err := GenMultipleColumnsTableImage(nil, headers, data)
	if err != nil {
		t.Fatalf("GenMultipleColumnsTableImage failed: %v", err)
	}

	if buf == nil {
		t.Fatal("Expected non-nil buffer")
	}

	if buf.Len() == 0 {
		t.Fatal("Expected non-empty buffer")
	}

	// 保存测试图片
	err = os.WriteFile("/tmp/test_columns_table.png", buf.Bytes(), 0o644)
	if err != nil {
		t.Logf("warning: failed to save columns table test image | output_path: /tmp/test_columns_table.png | err: %v", err)
	} else {
		t.Logf("columns table test image saved successfully | output_path: /tmp/test_columns_table.png")
	}

	t.Logf("columns table generated successfully | image_size_bytes: %d", buf.Len())
}
