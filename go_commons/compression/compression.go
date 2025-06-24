package compression

type Compression int8

var (
	None   Compression = 0
	Snappy Compression = 1
	GZIP   Compression = 2
)

var compressionCompressor = map[Compression]Compressor{
	None:   NewNoneParser(),
	Snappy: NewSnappyCompressor(),
	GZIP:   NewGzipCompressor(),
}

func GetCompressionParser(c Compression) Compressor {
	parser, ok := compressionCompressor[c]
	if ok {
		return parser
	}

	return compressionCompressor[None]
}
