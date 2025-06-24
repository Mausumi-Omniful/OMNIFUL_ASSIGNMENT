package requests

type PdfMergeRequest struct {
	PdfURLs           PdfURLs `json:"pdf_urls"`
	MergedPDFFilename string  `json:"merged_pdf_filename"`
}
type PdfURLs []PdfURL

type PdfURL struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// PrintAWBLabelRequest represents the request payload for generating an AWB label.
type PrintAWBLabelRequest struct {
	Service               string         `json:"service"`                  // Name of the service requesting the label (e.g., oms).
	InputTemplateFileName string         `json:"input_template_file_name"` // S3 path or identifier for the input template to use.
	Data                  []TemplateData `json:"data"`                     // Data used to populate the template. Supports up to 100 entries for multiple shipment boxes, each requiring its own AWB label. Exceeding the limit triggers a validation error.
	BarcodeConfig         BarcodeConfig  `json:"barcode_config"`           // Configuration for rendering barcodes on the label.
	PdfConfig             PdfConfig      `json:"pdf_config"`               // PDF generation settings (format, margins, etc.).
	OutputFileName        string         `json:"output_file_name"`         // Name of the generated PDF file. Can be empty. Note: Provided keys should be present in the template.
}

// TemplateData represents the overall template structure used for generating AWB labels.
// Add new attributes here as needed.
type TemplateData struct {
	EntityDetails   EntityDetails     `json:"entity_details"`   // Metadata related to the entity (e.g., order, return_order)
	ShipmentDetails ShipmentDetails   `json:"shipment_details"` // Information about the shipment
	SenderDetails   AddressDetails    `json:"sender_details"`   // Sender's contact and address info
	ReceiverDetails AddressDetails    `json:"receiver_details"` // Receiver's contact and address info
	BarcodeFields   map[string]string `json:"barcode_fields"`   // Fields for which we need to generate barcode
}

// EntityDetails contains key information about the entity(order,return_order).
type EntityDetails struct {
	EntityNumber string `json:"entity_number"` // Unique identifier for the entity
	PaymentType  string `json:"payment_type"`  // Payment mode (e.g., cod)
	TotalValue   string `json:"total_value"`   // Total value of the entity
	TotalDue     string `json:"total_due"`     // TotalDue represents the total payable amount for the entity
	CreatedAt    string `json:"created_at"`    // Timestamp when the entity was created
}

// ShipmentDetails holds specific information about the shipment.
type ShipmentDetails struct {
	AWBNumber        string         `json:"awb_number"`         // Air Waybill number for tracking
	CourierPartner   CourierPartner `json:"courier_partner"`    // Courier service handling the shipment
	NumberOfPackages int64          `json:"number_of_packages"` // Total packages in the shipment
	Weight           float64        `json:"weight"`             // Weight of the shipment (e.g., 2.0)
	Description      string         `json:"description"`        // Optional description of the shipment contents
	Remarks          string         `json:"remarks"`            // Remarks contains any additional notes or instructions related to the shipment.
	CreatedAt        string         `json:"created_at"`         // Timestamp when the shipment was created
}

// AddressDetails represents a generic address block used for sender and receiver.
type AddressDetails struct {
	Name          string `json:"name"`    // Full name of the contact person
	Address       string `json:"address"` // Full address
	State         string `json:"state"`
	City          string `json:"city"`
	Country       string `json:"country"`
	ContactNumber string `json:"contact_number"` // Phone number of the contact person
}

// CourierPartner defines attributes of the logistics provider.
type CourierPartner struct {
	Name string `json:"name"` // Name of the courier partner
	Logo string `json:"logo"` // URL for the courier partner's logo
}

// BarcodeConfig defines the configuration for rendering barcodes.
// If not provided, default values will be applied by the Lambda function.
type BarcodeConfig struct {
	Scale       float64 `json:"scale"`        // Scaling factor for the barcode image.
	Height      float64 `json:"height"`       // Height of the barcode.
	TextXAlign  string  `json:"text_x_align"` // Horizontal alignment of barcode text (e.g., "center").
	IncludeText bool    `json:"include_text"` // Whether to include human-readable text below the barcode.
}

// PdfConfig defines the configuration for generating the output PDF.
// If not provided, default values will be applied by the Lambda function.
type PdfConfig struct {
	Format          string       `json:"format"`           // Page format for the PDF (e.g., "A6", "A4").
	PrintBackground bool         `json:"print_background"` // Whether to print the background graphics/colors.
	Margin          MarginConfig `json:"margin"`           // Margin settings for the PDF.
}

// MarginConfig defines page margins for the PDF.
type MarginConfig struct {
	Top    string `json:"top"`    // Top margin (e.g., "10px").
	Right  string `json:"right"`  // Right margin.
	Bottom string `json:"bottom"` // Bottom margin.
	Left   string `json:"left"`   // Left margin.
}
