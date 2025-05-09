package errors

import "context"

type Factory interface {
	NewContext(ctx context.Context, text string) error
	WrapContext(ctx context.Context, format string, err error) error
}
