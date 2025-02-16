package errors

import (
	"context"
	"errors"
	"log/slog"
)

type FactoryImpl struct {
	logger *slog.Logger
}

func NewFactory(logger *slog.Logger) *FactoryImpl {
	return &FactoryImpl{logger}
}

func (f *FactoryImpl) NewContext(ctx context.Context, text string) error {
	f.logger.ErrorContext(ctx, text)
	return errors.New(text)
}
