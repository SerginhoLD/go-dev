package messenger

import "context"

type Bus interface {
	Send(ctx context.Context, msg any)
}
