package keydb

import (
	"context"
	"encoding/json"
	"fmt"
)

// see internal/domain/messenger/bus.go
func (bus *Stream) Send(ctx context.Context, msg any) {
	bytes, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	err = bus.Publish(ctx, "objects", string(bytes), &Header{"type", fmt.Sprintf("%T", msg)}) // todo: type

	if err != nil {
		panic(err)
	}
}
