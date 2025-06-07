package keydb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"time"

	"github.com/redis/go-redis/v9"
)

type Stream struct {
	client *redis.Client
}

func NewStream(client *redis.Client) *Stream {
	return &Stream{
		client: client,
	}
}

func (bus *Stream) Publish(ctx context.Context, stream string, payload string, opts ...option) error {
	maxLen := int64(0)
	headers := make(map[string]string)

	for _, opt := range opts {
		switch opt.typ() {
		case "maxlen":
			maxLen = int64(opt.val().(int64))
		case "header":
			maps.Copy(headers, opt.val().(map[string]string))
		}
	}

	headersBytes, err := json.Marshal(headers)

	if err != nil {
		return err
	}

	err = bus.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]any{
			"payload": payload,
			"headers": string(headersBytes),
		},
		MaxLen: maxLen,
		Approx: true, // MAXLEN ~
	}).Err()

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf(`redis: %s`, err.Error()), "stream", stream, "payload", payload, "headers", headers)
		return err
	} else {
		slog.DebugContext(ctx, fmt.Sprintf(`redis: send %s`, stream), "stream", stream, "payload", payload, "headers", headers)
		return nil
	}
}

func (bus *Stream) Consume(ctx context.Context, group string, stream string, callback func(context.Context, *Message) error, opts ...option) error {
	startId := "$" // start: "$" - latest, "0" - earliest
	delAfterAck := false

	for _, opt := range opts {
		switch opt.typ() {
		case "start_id":
			startId = opt.val().(string)
		case "del_after_ack":
			delAfterAck = opt.val().(bool)
		}
	}

	err := bus.client.XGroupCreateMkStream(ctx, stream, group, startId).Err()

	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf(`redis: %s`, err))

		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			return err
		}
	}

	// https://docs.keydb.dev/docs/streams-intro/#creating-a-consumer-group (see Ruby code)
	nextId := "0" // неподтвержденные сообщения
	checkBacklog := true

	for {
		if !checkBacklog {
			nextId = ">" // новые сообщения
		}

		xStreams, err := bus.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    group,
			Streams:  []string{stream, nextId},
			Consumer: "0",
			Count:    1,
			Block:    0,
			NoAck:    false,
		}).Result()

		if err != nil {
			return err
		}

		if len(xStreams[0].Messages) == 0 {
			checkBacklog = false
			time.Sleep(time.Second)
		}

		for _, xMessage := range xStreams[0].Messages {
			payload, hasPayload := xMessage.Values["payload"].(string)

			if !hasPayload {
				slog.ErrorContext(ctx, fmt.Sprintf(`redis: unknown message "%s"`, xMessage.ID), "values", xMessage.Values)

				if err = bus.client.XAck(ctx, stream, group, xMessage.ID).Err(); err != nil {
					slog.ErrorContext(ctx, err.Error())
					return err
				}

				continue
			}

			headers, hasHeaders := xMessage.Values["headers"].(string)

			if !hasHeaders {
				slog.ErrorContext(ctx, fmt.Sprintf(`redis: unknown message "%s"`, xMessage.ID), "values", xMessage.Values)

				if err = bus.client.XAck(ctx, stream, group, xMessage.ID).Err(); err != nil {
					slog.ErrorContext(ctx, err.Error())
					return err
				}

				continue
			}

			headersMap := make(map[string]string)
			json.Unmarshal([]byte(headers), &headersMap)

			slog.DebugContext(ctx, fmt.Sprintf(`redis: read "%s"`, xMessage.ID), "group", group, "stream", stream, "payload", payload, "headers", headers)

			message := &Message{
				client:  bus.client,
				group:   &group,
				stream:  &stream,
				id:      &xMessage.ID,
				payload: &payload,
				headers: headersMap,
			}

			if err = callback(ctx, message); err != nil {
				return err
			}

			if !message.acked {
				logMsg := fmt.Sprintf(`redis: message not acked (%s)`, xMessage.ID)
				slog.ErrorContext(ctx, logMsg)
				return errors.New(logMsg)
			}

			if delAfterAck {
				if err = bus.client.XDel(ctx, stream, xMessage.ID).Err(); err != nil {
					slog.ErrorContext(ctx, err.Error())
					return err
				}
			}
		}
	}
}

type Message struct {
	client  *redis.Client
	group   *string
	stream  *string
	acked   bool
	id      *string
	payload *string
	headers map[string]string
}

func (msg *Message) Id() string {
	return *msg.id
}

func (msg *Message) Payload() string {
	return *msg.payload
}

func (msg *Message) Headers() map[string]string {
	return msg.headers
}

func (msg *Message) Ack(ctx context.Context) error {
	if msg.acked {
		return nil
	}

	if err := msg.client.XAck(ctx, *msg.stream, *msg.group, *msg.id).Err(); err != nil {
		return err
	}

	msg.acked = true
	return nil
}
