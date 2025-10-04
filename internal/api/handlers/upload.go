package handlers

import (
	"dfood/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService service.UploadService
}

func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// Image Management
func (h *UploadHandler) UploadProfileImage(c *gin.Context) {
	// TODO: Implement upload profile images
	c.JSON(200, gin.H{"message": "Upload profile image - TODO"})
}

func (h *UploadHandler) UploadFoodImage(c *gin.Context) {
	// TODO: Implement upload food images
	c.JSON(200, gin.H{"message": "Upload food image - TODO"})
}

func (h *UploadHandler) UploadRestaurantImage(c *gin.Context) {
	// TODO: Implement upload restaurant images
	c.JSON(200, gin.H{"message": "Upload restaurant image - TODO"})
}

func (h *UploadHandler) DeleteImage(c *gin.Context) {
	// TODO: Implement delete uploaded image
	c.JSON(200, gin.H{"message": "Delete image - TODO"})
}
