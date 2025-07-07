package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"oms/database"
	"oms/models"
	"oms/utils"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
)

type OrderController struct {
	S3Uploader   *utils.S3UploaderImpl
	SQSPublisher *utils.SQSPublisherImpl
	OrderRepo    *database.OrderRepository
}

// validatecsv
func (h *OrderController) validateCSVContent(fileContent []byte) error {
	reader := bytes.NewReader(fileContent)
	csvReader := csv.NewReader(reader)

	header, err := csvReader.Read()
	if err != nil {
		return err
	}

	fmt.Println("Found headers:", header)

	requiredColumns := []string{"sku", "location", "tenant_id", "seller_id"}

	headerMap := make(map[string]bool)
	for _, col := range header {
		normalizedCol := strings.ToLower(strings.TrimSpace(col))
		headerMap[normalizedCol] = true
	}

	var missingColumns []string
	for _, requiredCol := range requiredColumns {
		if !headerMap[requiredCol] {
			missingColumns = append(missingColumns, requiredCol)
		}
	}

	if len(missingColumns) > 0 {
		errorMsg := fmt.Sprintf("missing required columns: %s", strings.Join(missingColumns, ", "))
		fmt.Println("ERROR:", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return fmt.Errorf("CSV file must contain at least one data row")
	}

	fmt.Println("CSV validation passed: found", len(rows), "data rows with required columns")
	return nil
}

// uploadcsv
func (h *OrderController) UploadCSV(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.file_not_found"),
		})
		return
	}
	defer file.Close()

	fmt.Println("File received:", header.Filename, "Size:", header.Size, "bytes")

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.invalid_format"),
		})
		return
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("ERROR: Failed to read file content:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "general.internal_error"),
		})
		return
	}

	fmt.Println("File content read:", len(fileContent), "bytes")

	if err := h.validateCSVContent(fileContent); err != nil {
		fmt.Println("ERROR: CSV validation failed:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.validation_failed") + ": " + err.Error(),
		})
		return
	}

	fmt.Println("CSV validation completed successfully for file:", header.Filename)

	s3Path, err := h.S3Uploader.UploadFile(c.Request.Context(), fileContent, header.Filename)
	if err != nil {
		fmt.Println("ERROR: Failed to upload file to S3:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "s3.upload_failed"),
		})
		return
	}

	fmt.Println("Publishing S3 path to SQS:", s3Path)
	err = h.SQSPublisher.PublishS3Path(c.Request.Context(), s3Path)
	if err != nil {
		fmt.Println("ERROR: Failed to publish S3 path to SQS:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "sqs.publish_failed") + ": " + err.Error(),
		})
		return
	}

	fmt.Println("Successfully queued CSV for processing:", s3Path)

	c.JSON(http.StatusOK, gin.H{
		"message":    "CSV file uploaded and queued for processing successfully",
		"s3_path":    s3Path,
		"filename":   header.Filename,
		"size":       len(fileContent),
		"queued":     true,
		"queue_name": h.SQSPublisher.GetQueueName(),
	})
}

// Listorders
func (h *OrderController) ListOrders(c *gin.Context) {
	filters := make(map[string]string)

	if tenantID := c.Query("tenant_id"); tenantID != "" {
		filters["tenant_id"] = tenantID
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		filters["seller_id"] = sellerID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	fmt.Println("Retrieving orders with filters:", filters, "Page:", page, "Limit:", limit)

	orders, err := h.OrderRepo.GetOrdersByFilter(c.Request.Context(), filters, limit, offset)
	if err != nil {
		fmt.Println("ERROR: Failed to retrieve orders from database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	response := gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(orders),
		},
		"filters": filters,
	}

	fmt.Println("Successfully retrieved", len(orders), "orders")
	c.JSON(http.StatusOK, response)
}

// getorderbyid
func (h *OrderController) GetOrderByID(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order ID is required",
		})
		return
	}

	fmt.Println("Retrieving order with ID:", orderID)

	order, err := h.OrderRepo.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		fmt.Println("ERROR: Failed to retrieve order:", err)

		if strings.Contains(err.Error(), "order not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Order not found with ID: %s", orderID),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve order",
		})
		return
	}

	fmt.Println("Successfully retrieved order - OrderID:", orderID, "Status:", order.Status)
	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// updateorder
func (h *OrderController) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	newStatus := models.OrderStatus(request.Status)
	if !newStatus.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order status: " + request.Status})
		return
	}

	err := h.OrderRepo.UpdateOrderStatus(c.Request.Context(), orderID, newStatus)
	if err != nil {
		fmt.Println("ERROR: Failed to update order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status: " + err.Error()})
		return
	}

	order, err := h.OrderRepo.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		fmt.Println("ERROR: Order updated but failed to retrieve:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order updated but failed to retrieve: " + err.Error()})
		return
	}

	fmt.Println("Order status updated successfully - OrderID:", orderID, "NewStatus:", newStatus)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"order":   order,
	})
}
