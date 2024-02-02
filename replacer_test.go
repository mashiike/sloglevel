package sloglevel_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/mashiike/sloglevel"
)

const (
	Notice slog.Level = slog.LevelInfo + 2
)

func Example() {
	buf := new(bytes.Buffer)
	h := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		ReplaceAttr: sloglevel.NewAttrReplacer(
			sloglevel.AddLevel(Notice, "NOTICE"),
			sloglevel.ToLower(),
			sloglevel.ChangeKey("severity"),
			sloglevel.NextAttrReplacer(func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				return a
			}),
		),
	})
	slog.SetDefault(slog.New(h))
	ctx := context.Background()
	slog.Log(ctx, slog.LevelDebug, "aaaaaa")
	slog.Log(ctx, slog.LevelInfo, "bbbbbb")
	slog.Log(ctx, Notice, "ccccc")
	slog.Log(ctx, slog.LevelWarn, "ddddd")
	slog.Log(ctx, slog.LevelError, "eeeee")
	fmt.Println(buf.String())
	// Output:
	// {"severity":"debug","msg":"aaaaaa"}
	// {"severity":"info","msg":"bbbbbb"}
	// {"severity":"notice","msg":"ccccc"}
	// {"severity":"warn","msg":"ddddd"}
	// {"severity":"error","msg":"eeeee"}
}
