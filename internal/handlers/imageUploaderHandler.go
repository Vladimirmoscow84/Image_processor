package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
	"github.com/gin-gonic/gin"
)

func (r *Router) imageUploaderHandler(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	root := os.Getenv("FILE_STORAGE_ROOT")
	if root == "" {
		root = "./data"
	}

	err = os.MkdirAll(root, 0o755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	origPath := filepath.Join(root, file.Filename)

	err = c.SaveUploadedFile(file, origPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	imgModel := &model.Image{
		OriginalPath: origPath,
		Status:       "enqueued",
	}
	id, err := r.imageUploader.AddImage(c.Request.Context(), imgModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = r.imageUploader.EnqueueImage(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "enqueued", "id": id})
}
