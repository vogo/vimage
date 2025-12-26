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

func TestCompressProcessor_CompressImage(t *testing.T) {
	bigFilePath := "/tmp/big_image.jpg"
	outputPath := "/tmp/big_image_compress.jpg"

	bigFileBytes, err := os.ReadFile(bigFilePath)
	if err != nil {
		t.Skipf("failed to open big image file | path: %s | err: %v", bigFilePath, err)
	}

	// Create a compress processor
	processor := vimage.NewCompressProcessor(1024*100, vimage.CompressStrategyDefault)

	// Process the image file bytes
	resultBytes, err := processor.ProcessFile(bigFileBytes)
	if err != nil {
		t.Fatalf("Compress processing failed: %v", err)
	}

	// Verify the result is not nil
	if resultBytes == nil {
		t.Fatal("Compress result is nil")
	}

	// Save the compressed image to a file
	if err := os.WriteFile(outputPath, resultBytes, 0o644); err != nil {
		t.Fatalf("failed to save compress image | path: %s | err: %v", outputPath, err)
	}
}
