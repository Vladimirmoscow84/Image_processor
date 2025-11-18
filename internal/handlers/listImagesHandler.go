package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Router) listImagesHandler(c *gin.Context) {
	images, err := r.listImageGetter.GetAllImages(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := []gin.H{}
	for _, img := range images {
		result = append(result, gin.H{
			"id":            img.ID,
			"status":        img.Status,
			"thumbnailPath": img.ThumbnailPath,
		})
	}

	c.JSON(http.StatusOK, result)
}
