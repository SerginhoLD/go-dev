package logger

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"testing"
)

func TestLog(t *testing.T) {
	var tests = []struct {
		attrs []any
		json  string
	}{
		{
			attrs: []any{slog.Int("statusCode", 201), slog.String("orderCode", "Abc")},
			json:  `"level":"INFO","msg":"Test","context":"{\"orderCode\":\"Abc\",\"statusCode\":201}"}`,
		},
		{
			attrs: []any{slog.Int("statusCode", 404), slog.String("orderCode", "Cde"), slog.String("err", "Not Found"), slog.String("name", "Foo")},
			json:  `"level":"INFO","msg":"Test","err":"Not Found","context":"{\"name\":\"Foo\",\"orderCode\":\"Cde\",\"statusCode\":404}"}`,
		},
		{
			attrs: []any{slog.Int("statusCode", 500), slog.Int("err", 2)},
			json:  `"level":"INFO","msg":"Test","context":"{\"err\":2,\"statusCode\":500}"}`,
		},
		{
			attrs: []any{slog.Group("payload", slog.Int("statusCode", 200))},
			json:  `"level":"INFO","msg":"Test","context":"{\"payload\":{\"statusCode\":200}}"}`,
		},
		{
			attrs: []any{slog.Group("payload", slog.Group("two", slog.String("statusCode", "201")), slog.String("orderCode", "Abc"))},
			json:  `"level":"INFO","msg":"Test","context":"{\"payload\":{\"orderCode\":\"Abc\",\"two\":{\"statusCode\":\"201\"}}}"}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Example %d", i), func(t *testing.T) {
			var b bytes.Buffer
			logger := NewLogger(NewHandler(io.Writer(&b)))
			logger.Info("Test", tt.attrs...)

			r := regexp.MustCompile(`{"time":".+?",(.+)\n`)
			got := r.ReplaceAllString(b.String(), "$1")

			if tt.json != got {
				t.Errorf("Expected `%s`, got `%s`", tt.json, got)
			}
		})
	}
}

func TestRequestId(t *testing.T) {
	t.Run("X-Request-ID", func(t *testing.T) {
		var b bytes.Buffer
		logger := NewLogger(NewHandler(io.Writer(&b)))
		logger.DebugContext(context.WithValue(context.Background(), "X-Request-ID", "01951302-6642-7007-9539-4c0cc944e4eb"), "Test")

		r := regexp.MustCompile(`{"time":".+?",(.+)\n`)
		got := r.ReplaceAllString(b.String(), "$1")

		expect := `"level":"DEBUG","msg":"Test","X-Request-ID":"01951302-6642-7007-9539-4c0cc944e4eb"}`

		if expect != got {
			t.Errorf("Expected `%s`, got `%s`", expect, got)
		}
	})
}
