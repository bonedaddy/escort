package pkg

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
)

// performs DEFLATE compression and base64 encoding
// of the input data, returning an array of the segmented
// data to store.
//
// segment specifies the segment identifier that is used to
// separate the index from the data.
func Compress(
	data []byte,
	segment string,
) ([]string, error) {
	buffer := new(bytes.Buffer)
	writer, err := flate.NewWriter(buffer, flate.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}
	// 250 is the maximum size of data in byteswe will return per record
	parts := Chunks(base64.StdEncoding.EncodeToString(buffer.Bytes()), 250)
	segmentedData := make([]string, 0, len(parts))
	for i, part := range parts {
		segmentedData = append(segmentedData, fmt.Sprintf("%v%s%s\n", i, segment, part))
	}
	return segmentedData, nil
}

// Chunks is used to split a string into segments of chunkSize
// https://stackoverflow.com/questions/25686109/split-string-by-length-in-golang
func Chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}
