package pkg

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompressNoXor(t *testing.T) {
	dataToCompress := "hello world 420 blaze it putin can ligma deez nuts"
	for i := 0; i < 10; i++ {
		dataToCompress = dataToCompress + dataToCompress
	}
	dataToCompressOrigChecksum := sha256.Sum256([]byte(dataToCompress))
	core := NewCore(nil, "|", 250)
	segmentedData, err := core.Trick([]byte(dataToCompress))
	require.NoError(t, err)
	require.Len(t, segmentedData, 2)
	require.Equal(t, "0|7MvRCcIwGAbAVb4RRFyo2qCB31Q0Rej0DuHrDXCPVrXlu71rzeV8yrWWo6XPvPbZR27LSPX7c8na2pGxzw9BEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEARBEA\n", segmentedData[0])
	require.Equal(t, "1|RBEARBEARBEARBEARBEARBEARBEARB/C9+AQAA//8=\n", segmentedData[1])
	trickedData, err := core.TurnOut(segmentedData)
	require.NoError(t, err)
	trickedDataChecksum := sha256.Sum256(trickedData)
	require.Equal(t, dataToCompressOrigChecksum, trickedDataChecksum)
}

func TestCompressXor(t *testing.T) {
	var letters = "abcdefghijklmnopqrstuvwxyz"
	dataToCompress := "hello world 420 blaze it putin can ligma deez nuts"
	for i := 0; i < 10; i++ {
		dataToCompress = dataToCompress + dataToCompress
	}
	dataToCompressOrigChecksum := sha256.Sum256([]byte(dataToCompress))
	for _, letter := range letters {
		fmt.Println("using xorkey ", string(letter))
		xorKey := byte(letter)
		core := NewCore(&xorKey, "|", 250)
		segmentedData, err := core.Trick([]byte(dataToCompress))
		require.NoError(t, err)
		require.Len(t, segmentedData, 2)
		trickedData, err := core.TurnOut(segmentedData)
		require.NoError(t, err)
		trickedDataChecksum := sha256.Sum256(trickedData)
		require.Equal(t, dataToCompressOrigChecksum, trickedDataChecksum)
	}

}
