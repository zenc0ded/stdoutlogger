package stdoutlogger

import (
	"io"
	"os"
)

// Composable options, in same style as otel stdout module

type config struct {
	Writer     io.Writer // output to something other than os.Stdout
	Timestamps bool      // prefix log lines with the timestamp
	Source     bool      // include any source references in the msg
	Style      Styles    // style to apply to log output
}

func newConfig(options []Option) config {
	cfg := config{
		Writer:     os.Stdout,
		Timestamps: true,
		Source:     true,
		Style:      *DefaultStyles(),
	}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}
	return cfg
}

// Option configures the exporter. Use the With functions to create one.
type Option interface {
	apply(config) config
}

// WithWriter sets where output is written. The default is os.Stdout.
// A nil writer is ignored.
func WithWriter(w io.Writer) Option {
	return writerOption{w}
}

type writerOption struct {
	W io.Writer
}

func (o writerOption) apply(cfg config) config {
	cfg.Writer = o.W
	return cfg
}

// WithoutTimestamps removes the timestamp from each line.
func WithoutTimestamps() Option {
	return timestampsOption(false)
}

type timestampsOption bool

func (o timestampsOption) apply(cfg config) config {
	cfg.Timestamps = bool(o)
	return cfg
}

// WithStyle sets the styles used for output. Start from DefaultStyles and
// change what you need. The maps in style are copied.
func WithStyle(style Styles) Option {
	return styleOption(style)
}

type styleOption Styles

func (o styleOption) apply(cfg config) config {
	cfg.Style = Styles(o).clone()
	return cfg
}

// WithoutSource removes the source location from each line.
func WithoutSource() Option {
	return sourceOption(false)
}

type sourceOption bool

func (o sourceOption) apply(cfg config) config {
	cfg.Source = bool(o)
	return cfg
}
