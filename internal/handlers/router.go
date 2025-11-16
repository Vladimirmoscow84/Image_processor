package handlers

import (
	"context"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
)

type imageLoader interface {
	ProcessAndSaveImage(ctx context.Context, origPath string) (*model.Image, error)
}

type imageReader interface {
}

type imageModifyer interface {
}

type imageQueue interface {
}
