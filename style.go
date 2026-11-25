package stdoutlogger

// Copyright (c) 2022-2023 Charmbracelet, Inc
// https://github.com/charmbracelet/log/blob/main/styles.go
//
// Copied to avoid pulling in additional dependencies :\
// Would be nice to keep it common and build up some common styles in a module.
// Will try to revisit this to see if can be reused in future.
//

import (
	"maps"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/otel/log"
)

// Styles defines the styles for the text logger.
type Styles struct {
	// Timestamp is the style for timestamps.
	Timestamp lipgloss.Style

	// Caller is the style for source caller.
	Caller lipgloss.Style

	// Prefix is not used by this exporter. It is kept to match
	// charmbracelet/log.
	Prefix lipgloss.Style

	// Message is the style for messages.
	Message lipgloss.Style

	// Key is the style for keys.
	Key lipgloss.Style

	// Value is the style for values.
	Value lipgloss.Style

	// Separator is the style for separators.
	Separator lipgloss.Style

	// Levels are the styles for each level.
	Levels map[log.Severity]lipgloss.Style

	// Keys overrides styles for specific keys.
	Keys map[string]lipgloss.Style

	// Values overrides value styles for specific keys.
	Values map[string]lipgloss.Style
}

func (s *Styles) styleForKey(key string) lipgloss.Style {
	if kStyle, ok := s.Keys[key]; ok {
		return kStyle
	}
	return s.Key
}

func (s *Styles) styleForValue(key string) lipgloss.Style {
	if vStyle, ok := s.Values[key]; ok {
		return vStyle
	}
	return s.Value
}

// clone returns a copy of s that does not share maps with s.
func (s Styles) clone() Styles {
	s.Levels = maps.Clone(s.Levels)
	s.Keys = maps.Clone(s.Keys)
	s.Values = maps.Clone(s.Values)
	return s
}

// DefaultStyles returns the default styles.
func DefaultStyles() *Styles {
	return &Styles{
		Timestamp: lipgloss.NewStyle(),
		Caller:    lipgloss.NewStyle().Faint(true),
		Prefix:    lipgloss.NewStyle().Bold(true).Faint(true),
		Message:   lipgloss.NewStyle(),
		Key:       lipgloss.NewStyle().Faint(true),
		Value:     lipgloss.NewStyle(),
		Separator: lipgloss.NewStyle().Faint(true),
		Levels: map[log.Severity]lipgloss.Style{
			log.SeverityDebug: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.SeverityDebug.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color("63")),
			log.SeverityInfo: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.SeverityInfo.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color("86")),
			log.SeverityWarn: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.SeverityWarn.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color("192")),
			log.SeverityError: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.SeverityError.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color("204")),
			log.SeverityFatal: lipgloss.NewStyle().
				SetString(strings.ToUpper(log.SeverityFatal.String())).
				Bold(true).
				MaxWidth(4).
				Foreground(lipgloss.Color("134")),
		},
		Keys: map[string]lipgloss.Style{
			"err": lipgloss.NewStyle().Foreground(lipgloss.Color("204")),
		},
		Values: map[string]lipgloss.Style{
			"err": lipgloss.NewStyle().Bold(true),
		},
	}
}
