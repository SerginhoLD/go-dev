package errors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

func New(text string) error {
	slog.Error(text)
	return errors.New(text)
}

func NewContext(ctx context.Context, text string) error {
	slog.ErrorContext(ctx, text)
	return errors.New(text)
}

func WrapContext(ctx context.Context, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	slog.ErrorContext(ctx, err.Error())
	return err
}
