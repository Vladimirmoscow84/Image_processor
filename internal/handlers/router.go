package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type imageUploader interface {
	AddImage(ctx context.Context, img *model.Image) (int, error)
	EnqueueImage(ctx context.Context, imageID int) error
}

type imageGetter interface {
	GetImage(ctx context.Context, id int) (*model.Image, error)
}

type imageDeleter interface {
	DeleteImage(ctx context.Context, image *model.Image) error
}

type Router struct {
	Router        *ginext.Engine
	imageUploader imageUploader
	imageGetter   imageGetter
	imageDeleter  imageDeleter
}

func New(router *ginext.Engine, imageUploader imageUploader, imageGetter imageGetter, imageDeleter imageDeleter) *Router {
	return &Router{
		Router:        router,
		imageUploader: imageUploader,
		imageGetter:   imageGetter,
		imageDeleter:  imageDeleter,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/upload", r.imageUploaderHandler)
	r.Router.GET("/image/:id", r.imageGetterHandler)
	r.Router.DELETE("/image/:id", r.imageDeleterHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}
