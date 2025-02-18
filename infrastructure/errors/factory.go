package errors

import (
	"context"
	"errors"
	"fmt"
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

func (f *FactoryImpl) WrapContext(ctx context.Context, format string, err error) error {
	wrapErr := fmt.Errorf(format, err)
	f.logger.ErrorContext(ctx, wrapErr.Error())
	return wrapErr
}
