# CSV Package

The CSV package offers a comprehensive solution for processing CSV files in Go. It is designed to handle different CSV sources and destinations including local files, AWS S3, raw byte streams, and URLs. The package also provides tools for sanitizing CSV content, converting CSV rows into maps, decoding compressed files, and batch processing data.

## Features

### CSV Reading
- **Multiple Sources**: Read CSV data from local files, AWS S3, raw bytes, or via a URL.
- **Batch Processing**: Process CSV records in batches to reduce memory usage.
- **Header Parsing**: Automatically parse and sanitize headers. Supports using custom sanitization functions.
- **Conversion**: Easily convert rows to maps or JSON objects.

### CSV Writing
- **Flexible Destination**: Write CSV files to local storage or upload them directly to AWS S3.
- **BOM Support**: Optionally write Byte Order Mark (BOM) at the beginning of the CSV output.
- **Header Management**: Ensure headers are written once and include options for dynamic file naming (with randomization).

### Data Sanitization
- **Sanitizers**: Built-in functions to trim spaces, remove BOM characters, convert case, and remove undesired asterisk characters.
- **Customizable**: Configure header and data row sanitizers through options when initializing CSV readers or writers.

### Downloading and Decoding
- **Downloader**: Download CSV files via HTTP with support for file chunking.
- **Decoder**: Support for decoding files compressed in gzip or zip format. The package automatically selects the appropriate decoder based on file encoding.

### Importer Utility
- **CSV Importer**: Combines reading and writing capabilities and allows you to process CSV data in batches.
- **Notification of Failed Entries**: Option to notify or process failed rows with a custom writer.
- **Flexible Options**: Configure the importer using functional options for integrating custom readers and writers.

### Configuration Options
- **Options Pattern**: Use functional options to configure CSV readers and writers with batch sizes, sanitizers, source details, and more.
- **Multiple Sources & Destinations**: The package supports different source types (e.g., `s3`, `local`, `raw-bytes`, `raw-url`) and destination types (`s3`, `local`).

## Usage Examples

### CSV Reader Example

```go
csvReader, err := csv.NewCommonCSV(
    csv.WithBatchSize(100),
    csv.WithSource(csv.Local),
    csv.WithLocalFileInfo("data.csv"),
    csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
    csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
)
if err != nil {
    log.Fatal(err)
}
err = csvReader.InitializeReader(context.TODO())
if err != nil {
    log.Fatal(err)
}

for !csvReader.IsEOF() {
    var records Records
    records, err = csvReader.ReadNextBatch()
    if err != nil {
        log.Fatal(err)
    }
    // Process the records
    fmt.Println(records)
}
```

### CSV Writer Example

```go
destination, err := csv.NewDestination(
    csv.WithType(csv.DestinationS3),
    csv.WithBucketName("my-bucket"),
    csv.WithUploadDirectory("uploads"),
)
if err != nil {
    log.Fatal(err)
}
csvWriter, err := csv.NewCommonCSVWriter(csv.WithWriterDestination(*destination))
if err != nil {
    log.Fatal(err)
}
err = csvWriter.Initialize()
if err != nil {
    log.Fatal(err)
}

csvWriter.SetHeaders([]string{"Name", "Age", "Email"})
err = csvWriter.WriteNextBatch(Records{
    {"Alice", "30", "alice@example.com"},
    {"Bob", "25", "bob@example.com"},
})
if err != nil {
    log.Fatal(err)
}
err = csvWriter.Close(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

### CSV Importer Example

```go
importer, err := csv.NewCSVImporter(
    csv.WithNotifyFailedRows(true),
    csv.WithReader(csvReader),
    csv.WithWriter(csvWriter),
)
if err != nil {
    log.Fatal(err)
}
err = importer.Initialize(context.TODO())
if err != nil {
    log.Fatal(err)
}

for !importer.IsEOF() {
    var data interface{}
    err = importer.ParseNextBatch(&data)
    if err != nil {
       log.Printf("Error processing batch: %v", err)
       break
    }
    // Process the imported data
    fmt.Println(data)
}
err = importer.Close(context.TODO())
if err != nil {
    log.Fatal(err)
}
```

### Reading CSV from a Raw URL Example

```go
csvReaderURL, err := csv.NewCommonCSV(
   csv.WithSource(csv.RawURL),
   csv.WithRawURLInfo("http://example.com/data.csv.gz", csv.EncodingGzip),
   csv.WithDelimiter(','),
)
if err != nil {
   log.Fatal(err)
}
err = csvReaderURL.InitializeReader(context.TODO())
if err != nil {
   log.Fatal(err)
}
records, err := csvReaderURL.ReadNextBatch()
if err != nil {
   log.Fatal(err)
}
fmt.Println(records)
```

### Custom Sanitizer Example

```go
customSanitizer := func(strSlice []string) []string {
    for i, val := range strSlice {
       strSlice[i] = strings.Trim(val, "#")
    }
    return strSlice
}

csvReaderCustom, err := csv.NewCommonCSV(
    csv.WithBatchSize(50),
    csv.WithSource(csv.Local),
    csv.WithLocalFileInfo("custom.csv"),
    csv.WithHeaderSanitizers(customSanitizer, csv.SanitizeToLower),
)
if err != nil {
    log.Fatal(err)
}
err = csvReaderCustom.InitializeReader(context.TODO())
if err != nil {
    log.Fatal(err)
}
header, err := csvReaderCustom.GetHeaders()
if err != nil {
   log.Fatal(err)
}
fmt.Println("Custom Headers:", header)
``` 