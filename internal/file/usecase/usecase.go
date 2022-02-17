package usecase

import (
	"context"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/file"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
)

type fileUC struct {
	cfg      *config.Config
	fileRepo file.FileRepository
	logger   logger.Logger
}

func NewFileUseCase(cfg *config.Config, fileRepo file.FileRepository, logger logger.Logger) file.UseCase {
	return &fileUC{cfg: cfg, fileRepo: fileRepo, logger: logger}
}

func (u *fileUC) PutObject(ctx context.Context, input models.UploadInput) (*string, error) {
	ctx, span := tracer.NewSpan(ctx, "fileUC.PutObject", nil)
	defer span.End()

	return u.fileRepo.PutObject(ctx, input)
}
func (u *fileUC) RemoveObject(ctx context.Context, fileName string) error {
	ctx, span := tracer.NewSpan(ctx, "fileUC.RemoveObject", nil)
	defer span.End()

	return u.fileRepo.RemoveObject(ctx, fileName)
}
