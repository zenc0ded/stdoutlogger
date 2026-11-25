package stdoutlogger_test

import (
	"context"
	"os"
	"time"

	"github.com/zenc0ded/stdoutlogger"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func ExampleNew() {
	exp, err := stdoutlogger.New(
		stdoutlogger.WithWriter(os.Stdout),
		stdoutlogger.WithoutSource(),
	)
	if err != nil {
		panic(err)
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)),
	)
	defer provider.Shutdown(context.Background())

	var r otellog.Record
	r.SetTimestamp(time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC))
	r.SetSeverity(otellog.SeverityInfo)
	r.SetBody(attribute.StringValue("hello"))
	r.AddAttributes(attribute.String("user", "dan"), attribute.Int("count", 3))
	provider.Logger("example").Emit(context.Background(), r)

	// Output: 2025/01/02 03:04:05 INFO hello user="dan" count=3
}
