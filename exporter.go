package stdoutlogger

import (
	"context"

	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/otel/sdk/log"
)

type exporter struct {
	renderer   *lipgloss.Renderer
	style      *Styles
	timestamps bool
	source     bool
}

func New(options ...Option) (log.Exporter, error) {
	cfg := newConfig(options)

	e := exporter{
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
		e.renderRecord(record)
	}
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	return nil
}

func (e *exporter) ForceFlush(context.Context) error {
	return nil
}
