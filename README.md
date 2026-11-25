# stdoutlogger

An OpenTelemetry log exporter that prints records to stdout (or any writer) as styled, human readable lines. It combines [charmbracelet/log](https://github.com/charmbracelet/log) styling with the [otel stdoutlog exporter](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/stdout/stdoutlog).

## Install

```
go get github.com/zenc0ded/stdoutlogger
```

Requires Go 1.26 or newer and `go.opentelemetry.io/otel/log` and `sdk/log` v1.47 or newer.

## Usage

```go
exp, err := stdoutlogger.New()
if err != nil {
	panic(err)
}
provider := log.NewLoggerProvider(log.WithProcessor(log.NewSimpleProcessor(exp)))
defer provider.Shutdown(context.Background())
global.SetLoggerProvider(provider)

slog.SetDefault(otelslog.NewLogger("app", otelslog.WithSource(true)))
slog.Info("hello", "user", "dan")
```

Output:

```
2025/01/02 03:04:05 INFO <main/main.go:20> hello user="dan"
```

A complete program is in [example/](example/main.go).

## Options

| Option | Effect |
|---|---|
| `WithWriter(w)` | Write to `w` instead of `os.Stdout` |
| `WithoutTimestamps()` | Omit the timestamp |
| `WithoutSource()` | Omit the source location |
| `WithStyle(s)` | Use custom styles. Start from `DefaultStyles()` |

## Notes

- The source location comes from the `code.function.name`, `code.file.path` and `code.line.number` attributes (or the older `code.function`, `code.filepath` and `code.lineno`). If they are missing, no location is printed.
- Colors are only used when the writer is a terminal.
- Each record is written with a single `Write` call.

## License

MIT. See [LICENSE](LICENSE).
