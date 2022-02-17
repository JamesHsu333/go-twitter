package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/file"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type fileRepository struct {
	cfg *config.Config
}

func NewFileRepository(cfg *config.Config) file.FileRepository {
	return &fileRepository{cfg: cfg}
}

func (f *fileRepository) PutObject(ctx context.Context, input models.UploadInput) (*string, error) {
	_, span := tracer.NewSpan(ctx, "fileRepository.PutObject", nil)
	defer span.End()

	filepath := path.Join(f.cfg.File.FilePath, f.generateFileName(input.Name))

	dst, err := os.Create(filepath)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "fileRepository.PutObject.os.Create")
	}
	defer dst.Close()

	if _, err = io.Copy(dst, input.File); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "fileRepository.PutObject.io.Copy")
	}
	return &filepath, nil
}

func (f *fileRepository) RemoveObject(ctx context.Context, filepath string) error {
	_, span := tracer.NewSpan(ctx, "fileRepository.RemoveObject", nil)
	defer span.End()

	err := os.Remove(filepath)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "fileRepository.RemoveObject.os.Remove")
	}

	return nil
}

func (f *fileRepository) generateFileName(fileName string) string {
	uid := uuid.New().String()
	return fmt.Sprintf("%s-%s", uid, fileName)
}
