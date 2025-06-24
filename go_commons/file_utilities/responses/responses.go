package responses

type FailedPdfs []FailedPdf

type FailedPdf struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type PdfMergeResponse struct {
	FailedPdfs   FailedPdfs `json:"failed_pdfs"`
	MergedPdfURL string     `json:"merged_pdf_url"`
}

type PrintAWBLabelResponse struct {
	Message string `json:"message"`
	FileURL string `json:"file_url"`
}
