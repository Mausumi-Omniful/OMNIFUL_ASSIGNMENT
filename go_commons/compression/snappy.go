package compression

import (
	"encoding/base64"
	"github.com/golang/snappy"
)

type SnappyCompressor struct {
}

func NewSnappyCompressor() SnappyCompressor {
	return SnappyCompressor{}
}

func (s SnappyCompressor) Compression() Compression {
	return Snappy
}

// Compress compresses the input data using Snappy compression
func (s SnappyCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Create a buffer to store compressed data
	compressed := snappy.Encode(nil, data)

	// Base64-encode the compressed data
	encoded := base64.StdEncoding.EncodeToString(compressed)

	return []byte(encoded), nil
}

// Decompress decompresses Snappy-compressed data
func (s SnappyCompressor) Decompress(compressedData []byte) ([]byte, error) {
	if len(compressedData) == 0 {
		return nil, nil
	}

	// Decode the Base64-encoded string
	base64Decode, err := base64.StdEncoding.DecodeString(string(compressedData))
	if err != nil {
		return nil, err
	}

	// Decode the compressed data
	decoded, err := snappy.Decode(nil, base64Decode)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}
