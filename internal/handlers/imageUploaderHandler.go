package handlers

import (
	"log"
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

	log.Println("UPLOAD: received file:", file.Filename)

	uploadDir := "data/uploads"
	os.MkdirAll(uploadDir, 0755)

	origPath := filepath.Join(uploadDir, file.Filename)

	log.Println("UPLOAD: saving to:", origPath)

	err = c.SaveUploadedFile(file, origPath)
	if err != nil {
		log.Println("UPLOAD: SaveUploadedFile error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("UPLOAD: saved OK")

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
