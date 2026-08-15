package service

import (
	"context"

	"github.com/1260124186-cc/expense-review-cli-20260815/internal/domain"
)

type Repository interface {
	Load(context.Context, string) (domain.Batch, error)
}

type Writer interface {
	Write(string, string) error
}

type Reviewer struct {
	repository Repository
	writer     Writer
}

func New(repository Repository, writer Writer) Reviewer {
	return Reviewer{repository: repository, writer: writer}
}

func (r Reviewer) Review(ctx context.Context, inputPath string) (domain.Review, error) {
	batch, err := r.repository.Load(ctx, inputPath)
	if err != nil {
		return domain.Review{}, err
	}
	return domain.ReviewBatch(ctx, batch)
}

func (r Reviewer) ReviewAndRender(ctx context.Context, inputPath string) (string, error) {
	review, err := r.Review(ctx, inputPath)
	if err != nil {
		return "", err
	}
	return Render(review), nil
}

func (r Reviewer) Write(path, rendered string) error {
	return r.writer.Write(path, rendered)
}
