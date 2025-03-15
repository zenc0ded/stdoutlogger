package stdoutlogger

import (
	"io"
	"os"
)

// Composable options, in same style as otel stdout module

var (
	defaultWriter     io.Writer = os.Stdout
	defaultTimestamps           = true
	defaultSource               = true
	defaultStyle                = DefaultStyles()
)

type config struct {
	Writer     io.Writer // output to something other than os.Stdout
	Timestamps bool      // prefix log lines with the timestamp
	Source     bool      // include any source references in the msg
	Style      Styles    // style to apply to log output
}

func newConfig(options []Option) config {
	cfg := config{
		Writer:     defaultWriter,
		Timestamps: defaultTimestamps,
		Source:     defaultSource,
		Style:      *defaultStyle,
	}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}
	return cfg
}

type Option interface {
	apply(config) config
}

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

func WithoutTimestamps() Option {
	return timestampsOption(false)
}

type timestampsOption bool

func (o timestampsOption) apply(cfg config) config {
	cfg.Timestamps = bool(o)
	return cfg
}

func WithStyle(style Styles) Option {
	return styleOption(style)
}

type styleOption Styles

func (o styleOption) apply(cfg config) config {
	cfg.Style = Styles(o)
	return cfg
}

func WithoutSource() Option {
	return sourceOption(false)
}

type sourceOption bool

func (o sourceOption) apply(cfg config) config {
	cfg.Source = bool(o)
	return cfg
}
