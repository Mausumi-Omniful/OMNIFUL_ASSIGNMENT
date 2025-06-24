package compression

type NoneCompressor struct {
}

func NewNoneParser() NoneCompressor {
	return NoneCompressor{}
}

func (s NoneCompressor) Compression() Compression {
	return None
}

func (s NoneCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (s NoneCompressor) Decompress(compressedData []byte) ([]byte, error) {
	return compressedData, nil
}
