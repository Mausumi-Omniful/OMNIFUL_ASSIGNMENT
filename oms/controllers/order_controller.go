package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"oms/database"
	"oms/models"
	"oms/utils"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

type OrderController struct {
	S3Uploader   *utils.S3UploaderImpl
	SQSPublisher *utils.SQSPublisherImpl
	OrderRepo    *database.OrderRepository
}

// validateCSVContent validates that the CSV file contains the required columns
func (h *OrderController) validateCSVContent(fileContent []byte) error {
	// Create a reader from the file content
	reader := bytes.NewReader(fileContent)
	csvReader := csv.NewReader(reader)

	// Read the header row
	header, err := csvReader.Read()
	if err != nil {
		return err
	}

	log.Infof("Found headers: %v", header)

	// Define required columns
	requiredColumns := []string{"sku", "location", "tenant_id", "seller_id"}

	// Create a map for faster lookup
	headerMap := make(map[string]bool)
	for _, col := range header {
		normalizedCol := strings.ToLower(strings.TrimSpace(col))
		headerMap[normalizedCol] = true
	}

	// Check if all required columns are present
	var missingColumns []string
	for _, requiredCol := range requiredColumns {
		if !headerMap[requiredCol] {
			missingColumns = append(missingColumns, requiredCol)
		}
	}

	if len(missingColumns) > 0 {
		errorMsg := fmt.Sprintf("missing required columns: %s", strings.Join(missingColumns, ", "))
		log.Errorf("Validation failed: %s", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Check if file has at least one data row
	rows, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return fmt.Errorf("CSV file must contain at least one data row")
	}

	log.Infof("CSV validation passed: found %d data rows with required columns", len(rows))
	return nil
}

func (h *OrderController) UploadCSV(c *gin.Context) {
	// Get file from multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.file_not_found"),
		})
		return
	}
	defer file.Close()

	log.Infof("File received: name=%s, size=%d bytes", header.Filename, header.Size)

	// Validate file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.invalid_format"),
		})
		return
	}

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		log.WithError(err).Error("Failed to read file content")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "general.internal_error"),
		})
		return
	}

	log.Infof("File content read: %d bytes", len(fileContent))

	// Validate CSV content
	if err := h.validateCSVContent(fileContent); err != nil {
		log.WithError(err).Error("CSV validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": i18n.Translate(c.Request.Context(), "csv.validation_failed") + ": " + err.Error(),
		})
		return
	}

	log.Infof("CSV validation completed successfully for file: %s", header.Filename)

	// Upload file to S3
	s3Path, err := h.S3Uploader.UploadFile(c.Request.Context(), fileContent, header.Filename)
	if err != nil {
		log.WithError(err).Error("Failed to upload file to S3")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "s3.upload_failed"),
		})
		return
	}

	// Publish S3 path to SQS for processing
	log.Infof("Publishing S3 path to SQS: %s", s3Path)
	err = h.SQSPublisher.PublishS3Path(c.Request.Context(), s3Path)
	if err != nil {
		log.WithError(err).Error("Failed to publish S3 path to SQS")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": i18n.Translate(c.Request.Context(), "sqs.publish_failed") + ": " + err.Error(),
		})
		return
	}

	log.Infof("Successfully queued CSV for processing: %s", s3Path)

	c.JSON(http.StatusOK, gin.H{
		"message":    "CSV file uploaded and queued for processing successfully",
		"s3_path":    s3Path,
		"filename":   header.Filename,
		"size":       len(fileContent),
		"queued":     true,
		"queue_name": h.SQSPublisher.GetQueueName(),
	})
}

// ListOrders handles GET requests to retrieve orders with filtering and pagination
func (h *OrderController) ListOrders(c *gin.Context) {
	// Parse query parameters
	filters := make(map[string]string)

	// Add filters if provided
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		filters["tenant_id"] = tenantID
	}
	if sellerID := c.Query("seller_id"); sellerID != "" {
		filters["seller_id"] = sellerID
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Parse date filters
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Calculate offset
	offset := (page - 1) * limit

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	log.Infof("Retrieving orders with filters: %v, page: %d, limit: %d", filters, page, limit)

	// Get orders from repository
	orders, err := h.OrderRepo.GetOrdersByFilter(c.Request.Context(), filters, limit, offset)
	if err != nil {
		log.WithError(err).Error("❌ Failed to retrieve orders from database")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	// Build response
	response := gin.H{
		"orders": orders,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(orders), // For now, just return count of current page
		},
		"filters": filters,
	}

	log.Infof("✅ Successfully retrieved %d orders", len(orders))
	c.JSON(http.StatusOK, response)
}

// GetOrderByID handles GET requests to retrieve a specific order by ID
func (h *OrderController) GetOrderByID(c *gin.Context) {
	// Get order ID from URL parameter
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order ID is required",
		})
		return
	}

	log.Infof("🔍 Retrieving order with ID: %s", orderID)

	// Get order from repository
	order, err := h.OrderRepo.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		log.WithError(err).Errorf("❌ Failed to retrieve order from database - OrderID: %s", orderID)

		// Check if it's a "not found" error
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

	log.Infof("✅ Successfully retrieved order - OrderID: %s, Status: %s", orderID, order.Status)
	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// UpdateOrderStatus updates the status of an order
func (h *OrderController) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderID")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order ID is required"})
		return
	}

	// Parse request body
	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	// Validate status
	newStatus := models.OrderStatus(request.Status)
	if !newStatus.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order status: " + request.Status})
		return
	}

	// Update order status atomically
	err := h.OrderRepo.UpdateOrderStatus(c.Request.Context(), orderID, newStatus)
	if err != nil {
		log.WithError(err).Errorf("❌ Failed to update order status - OrderID: %s", orderID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status: " + err.Error()})
		return
	}

	// Get updated order to return
	order, err := h.OrderRepo.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		log.WithError(err).Errorf("❌ Failed to get updated order - OrderID: %s", orderID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order updated but failed to retrieve: " + err.Error()})
		return
	}

	log.Infof("✅ Order status updated successfully - OrderID: %s, NewStatus: %s", orderID, newStatus)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"order":   order,
	})
}

// TestEndpoint is a simple test endpoint to verify routing
func (h *OrderController) TestEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "Test endpoint working!",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// HealthCheck provides a health check endpoint for monitoring
func (h *OrderController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "oms-service",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}
