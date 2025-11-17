package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

func (r *Router) imageUploaderHandler(c *ginext.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	tempPath := "/tmp" + image.Filename
	err = c.SaveUploadedFile(image, tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = r.imageUploader.EnqueueImage(c.Request.Context(), tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "enqueued"})
}
