package compression

type Compressor interface {
	Compression() Compression
	Compress(data []byte) ([]byte, error)
	Decompress(compressedData []byte) ([]byte, error)
}
