# Compression Package

The **compression** package provides a simple way to compress and decompress data using different algorithms. It supports three compression modes:

- **None**: No compression is applied.
- **Snappy**: Uses the Snappy algorithm.
- **GZIP**: Uses the GZIP algorithm.

Each compression mode implements the `Compressor` interface that defines the following methods:

- `Compression() Compression`  
  Returns the corresponding compression constant.
  
- `Compress(data []byte) ([]byte, error)`  
  Compresses the given data and returns the compressed data or an error.
  
- `Decompress(compressedData []byte) ([]byte, error)`  
  Decompresses the provided data and returns the original data or an error.

## Functionalities

- **Compressor Interface**:  
  Defined in [`compressor.go`](./compressor.go), this interface specifies the methods required for any compression implementation.

- **Compression Constants**:  
  The type `Compression` (defined in [`compression.go`](./compression.go)) represents the type of compression. The three supported constants are:
  - `None`
  - `Snappy`
  - `GZIP`
  
- **Compressor Implementations**:
  - **NoneCompressor**:  
    A no-op compressor that returns the original data.  
    (See: [`none.go`](./none.go))
    
  - **SnappyCompressor**:  
    Uses Google's Snappy library for compression and base64 encodes the result.  
    (See: [`snappy.go`](./snappy.go))
    
  - **GzipCompressor**:  
    Uses Go's standard `gzip` library for compression and base64 encodes the compressed data.  
    (See: [`gzip.go`](./gzip.go))
    
- **Compressor Factory**:  
  The function `GetCompressionParser` (in [`compression.go`](./compression.go)) returns a `Compressor` implementation based on the provided compression constant. If an unsupported compression type is requested, it defaults to the `None` compressor.

## How to Use

1. **Import the package** in your code:
   ```go
   import "github.com/omniful/go_commons/compression"
   ```

2. **Select a Compression Mode**:  
   Choose one of the provided compression constants: `compression.None`, `compression.Snappy`, or `compression.GZIP`.

3. **Get the Compressor**:  
   Use `GetCompressionParser` to retrieve the appropriate compressor implementation.
   ```go
   compressor := compression.GetCompressionParser(compression.GZIP)
   ```

4. **Compress and Decompress Data**:  
   Call the `Compress` method to compress your data, and the `Decompress` method to revert to the original data.
   ```go
   func main() {
       // Sample data
       data := []byte("Hello, World!")
       
       // Pick a compressor; here we use GZIP compression
       compressor := compression.GetCompressionParser(compression.GZIP)
       
       // Compress the data
       compressedData, err := compressor.Compress(data)
       if err != nil {
           // Handle error
           panic(err)
       }
       fmt.Printf("Compressed: %s\n", compressedData)
       
       // Decompress the data
       decompressedData, err := compressor.Decompress(compressedData)
       if err != nil {
           // Handle error
           panic(err)
       }
       fmt.Printf("Decompressed: %s\n", decompressedData)
   }
   ```

## Summary

The **compression** package is designed to provide a consistent interface for data compression with three modes—`None`, `Snappy`, and `GZIP`. The available implementations adhere to the `Compressor` interface and can be retrieved via the `GetCompressionParser` factory method, making it easy to switch between different compression algorithms. 