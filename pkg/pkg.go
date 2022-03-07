package pkg

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

// Core implementation of escort
type Core struct {
	xorKey            *byte
	segmentIdentifier string
	segmentSize       int
}

func NewCore(
	xorKey *byte,
	segmentIdentifier string,
	segmentSize int,
) *Core {
	return &Core{
		xorKey:            xorKey,
		segmentIdentifier: segmentIdentifier,
		segmentSize:       segmentSize,
	}
}

// Trick optionally XOR's the binary data before performing
// DEFLATE compressiong, followed by base64 encoding, and segmenting
// the encoded data.
func (c *Core) Trick(
	data []byte,
) ([]string, error) {
	if c.xorKey != nil {
		c.xorData(data)
	}
	// if there is a key, perform xor encryption on the data first
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
	parts := Chunks(base64.StdEncoding.EncodeToString(buffer.Bytes()), c.segmentSize)
	segmentedData := make([]string, 0, len(parts))
	for i, part := range parts {
		segmentedData = append(segmentedData, fmt.Sprintf("%v%s%s\n", i, c.segmentIdentifier, part))
	}
	return segmentedData, nil
}

// Trick parses the compressed, encoded data segments (tricked data)
// into it's original form
func (c *Core) TurnOut(
	dataSegments []string,
) ([]byte, error) {
	type segment struct {
		// in a segment, this is the data portion minus the segment identifier
		data []byte
		// in a segment, this is the index portion
		index int64
	}
	var (
		// unsorted
		segments        = make([]segment, 0, len(dataSegments))
		totalSegmentLen int
	)
	for _, dataSegment := range dataSegments {
		totalSegmentLen += len(dataSegment)
		parts := strings.Split(dataSegment, c.segmentIdentifier)
		if len(parts) < 2 {
			return nil, fmt.Errorf("encountered invalid  %s", dataSegment)
		}
		segmentIndex, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse segment identifier %s", err)
		}
		segments = append(segments, segment{
			data:  []byte(parts[1]),
			index: segmentIndex,
		})
		totalSegmentLen += len(parts[1])
	}
	// sort segments
	sort.Slice(segments, func(i, j int) bool {
		return segments[i].index < segments[j].index
	})
	var (
		buf bytes.Buffer
	)
	buf.Grow(totalSegmentLen)
	for _, segment := range segments {
		// todo(bonedaddy): should we care about the amount of data written
		_, err := buf.Write(segment.data)
		if err != nil {
			return nil, fmt.Errorf("failed to write to buffer %s", err)
		}
	}
	compressedData, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode data %s", err)
	}
	// decompress the data
	// todo(bonedaddy): not very memory efficient
	decompressedData, err := ioutil.ReadAll(flate.NewReader(bytes.NewReader(compressedData)))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data %s", err)
	}
	buf.Reset()
	// check to see if we need to xor it first
	if c.xorKey != nil {
		c.xorData(decompressedData)
	}
	return decompressedData, nil
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

func (c *Core) xorData(data []byte) []byte {
	xorKey := *c.xorKey
	for idx, b := range data {
		data[idx] = b ^ xorKey
	}
	return data
}
