package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
)

type GzipCompressor struct {
}

func NewGzipCompressor() GzipCompressor {
	return GzipCompressor{}
}

func (g GzipCompressor) Compression() Compression {
	return GZIP
}

func (g GzipCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Create a buffer to store compressed data
	var compressed bytes.Buffer

	// Create gzip writer
	gzipWriter := gzip.NewWriter(&compressed)

	// Write data to the gzip writer
	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}

	// Close the writer to flush compressed data
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(compressed.Bytes())
	return []byte(encoded), nil
}

func (g GzipCompressor) Decompress(compressedData []byte) ([]byte, error) {
	if len(compressedData) == 0 {
		return nil, nil
	}

	// Decode the Base64 string
	compressed, err := base64.StdEncoding.DecodeString(string(compressedData))
	if err != nil {
		return nil, err
	}

	// Create a buffer from compressed data
	bytesReader := bytes.NewReader(compressed)

	// Create gzip reader
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Read all decompressed data
	decompressed, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
