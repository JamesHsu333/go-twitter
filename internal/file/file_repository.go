package file

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
)

type FileRepository interface {
	PutObject(ctx context.Context, input models.UploadInput) (*string, error)
	RemoveObject(ctx context.Context, fileName string) error
}
