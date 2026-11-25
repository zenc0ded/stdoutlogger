package stdoutlogger

import (
	"context"
	"io"

	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/otel/sdk/log"
)

type exporter struct {
	writer     io.Writer
	renderer   *lipgloss.Renderer
	style      *Styles
	timestamps bool
	source     bool
}

func New(options ...Option) (log.Exporter, error) {
	cfg := newConfig(options)

	e := exporter{
		writer:     cfg.Writer,
		renderer:   lipgloss.NewRenderer(cfg.Writer),
		style:      &cfg.Style,
		timestamps: cfg.Timestamps,
		source:     cfg.Source,
	}
	return &e, nil
}

func (e *exporter) Export(ctx context.Context, records []log.Record) error {
	for _, record := range records {
		if err := ctx.Err(); err != nil {
			return err
		}
		// one Write per record keeps lines intact when the writer is shared
		if _, err := io.WriteString(e.writer, e.renderRecord(record)); err != nil {
			return err
		}
	}
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	return nil
}

func (e *exporter) ForceFlush(context.Context) error {
	return nil
}
