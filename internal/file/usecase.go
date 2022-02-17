package file

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
)

// Follow usecase interface
type UseCase interface {
	PutObject(ctx context.Context, input models.UploadInput) (*string, error)
	RemoveObject(ctx context.Context, fileName string) error
}
