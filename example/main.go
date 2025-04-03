package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/zenc0ded/stdoutlogger"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

func main() {
	// style can be overriden or added to
	customStyle := stdoutlogger.DefaultStyles()
	customStyle.Keys["one"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#00EEFF")).Underline(true)
	customStyle.Values["one"] = lipgloss.NewStyle().Foreground(lipgloss.Color("#EE0088")).Bold(true)

	exp, err := stdoutlogger.New(
		stdoutlogger.WithWriter(os.Stderr),   // change from stdout
		stdoutlogger.WithoutSource(),         // don't print source references
		stdoutlogger.WithStyle(*customStyle), // use (slightly) modified style
	)
	if err != nil {
		panic(err)
	}

	processor := log.NewSimpleProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	global.SetLoggerProvider(provider)

	logger := otelslog.NewLogger(
		"example",
		otelslog.WithSource(true),
		otelslog.WithVersion("1.0"),
	)
	slog.SetDefault(logger)

	slog.Debug("Hello!")
	slog.Info("Hello!")
	slog.Warn("Hello!")
	slog.Error("Hello!")
	slog.Info("with attrs", "one", 1, "two", "2", "three", 3.0, "err", err)

	newLogger := logger.With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("version", runtime.Version()),
		),
	)
	newLogger.Info("group with attrs", "one", 1, "two", "2", "three", 3.0, "err", err)
}
