/*
 * Copyright 2025 Bronya0 <tangssst@163.com>.
 * Author Github: https://github.com/Bronya0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// GzipCompress 函数：压缩数据
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf) // 创建 gzip writer，写入到 buffer

	_, err := gzipWriter.Write(data) // 将原始数据写入 gzip writer 进行压缩
	if err != nil {
		return nil, fmt.Errorf("gzip write error: %w", err)
	}

	err = gzipWriter.Close() // **重要**: 关闭 writer，完成 gzip 流
	if err != nil {
		return nil, fmt.Errorf("gzip close error: %w", err)
	}

	return buf.Bytes(), nil // 返回 buffer 中的压缩数据
}

// GzipDecompress 函数：解压缩数据
func GzipDecompress(compressedData []byte) ([]byte, error) {
	buf := bytes.NewReader(compressedData) // 从压缩数据创建 reader
	gzipReader, err := gzip.NewReader(buf) // 创建 gzip reader，从 buffer 读取压缩数据
	if err != nil {
		return nil, fmt.Errorf("gzip reader creation error: %w", err)
	}
	defer gzipReader.Close() // 确保 reader 在使用后关闭

	decompressedData, err := io.ReadAll(gzipReader) // 读取所有解压缩后的数据
	if err != nil {
		return nil, fmt.Errorf("gzip read error: %w", err)
	}

	return decompressedData, nil // 返回解压缩后的数据
}
