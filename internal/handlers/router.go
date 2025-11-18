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
type listImageGetter interface {
	GetAllImages(ctx context.Context) ([]*model.Image, error)
}

type imageDeleter interface {
	DeleteImage(ctx context.Context, image *model.Image) error
}

type Router struct {
	Router          *ginext.Engine
	imageUploader   imageUploader
	imageGetter     imageGetter
	imageDeleter    imageDeleter
	listImageGetter listImageGetter
}

func New(router *ginext.Engine, imageUploader imageUploader, imageGetter imageGetter, imageDeleter imageDeleter, listImageGetter listImageGetter) *Router {
	return &Router{
		Router:          router,
		imageUploader:   imageUploader,
		imageGetter:     imageGetter,
		imageDeleter:    imageDeleter,
		listImageGetter: listImageGetter,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/upload", r.imageUploaderHandler)
	r.Router.GET("/image/:id", r.imageGetterHandler)
	r.Router.GET("/images", r.listImagesHandler)
	r.Router.DELETE("/image/:id", r.imageDeleterHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
	r.Router.Static("/data", "./data")

}
