package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) imageUploaderHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	origPath := "/tmp/" + file.Filename
	err = c.SaveUploadedFile(file, origPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = r.imageUploader.EnqueueImage(c.Request.Context(), origPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "enqueued"})
}
