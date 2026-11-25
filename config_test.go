package stdoutlogger

import (
	"bytes"
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewConfigDefaults(t *testing.T) {
	cfg := newConfig(nil)
	if cfg.Writer != os.Stdout {
		t.Error("default writer should be os.Stdout")
	}
	if !cfg.Timestamps || !cfg.Source {
		t.Error("timestamps and source should default to on")
	}
}

func TestOptions(t *testing.T) {
	var buf bytes.Buffer
	cfg := newConfig([]Option{WithWriter(&buf), WithoutTimestamps(), WithoutSource()})
	if cfg.Writer != &buf {
		t.Error("WithWriter not applied")
	}
	if cfg.Timestamps {
		t.Error("WithoutTimestamps not applied")
	}
	if cfg.Source {
		t.Error("WithoutSource not applied")
	}
}

func TestWithStyle(t *testing.T) {
	s := DefaultStyles()
	s.Keys["custom"] = lipgloss.NewStyle().Bold(true)
	cfg := newConfig([]Option{WithStyle(*s)})
	if _, ok := cfg.Style.Keys["custom"]; !ok {
		t.Error("WithStyle not applied")
	}
}

// Changing the style of one config must not leak into later configs.
func TestDefaultStyleIsNotShared(t *testing.T) {
	a := newConfig(nil)
	a.Style.Keys["leak"] = lipgloss.NewStyle().Bold(true)
	a.Style.Values["leak"] = lipgloss.NewStyle().Bold(true)
	delete(a.Style.Keys, "err")

	b := newConfig(nil)
	if _, ok := b.Style.Keys["leak"]; ok {
		t.Error("Keys map is shared between configs")
	}
	if _, ok := b.Style.Values["leak"]; ok {
		t.Error("Values map is shared between configs")
	}
	if _, ok := b.Style.Keys["err"]; !ok {
		t.Error("deleting from one config affected another")
	}
}

// The exporter must not see later changes made to the Styles passed in.
func TestWithStyleIsCopied(t *testing.T) {
	s := DefaultStyles()
	cfg := newConfig([]Option{WithStyle(*s)})
	s.Keys["late"] = lipgloss.NewStyle()
	if _, ok := cfg.Style.Keys["late"]; ok {
		t.Error("WithStyle kept a reference to the caller's maps")
	}
}

func TestStyleLookups(t *testing.T) {
	s := DefaultStyles()
	custom := lipgloss.NewStyle().Bold(true)
	s.Keys["k"] = custom
	s.Values["v"] = custom
	if s.styleForKey("k").GetBold() != true {
		t.Error("key override not used")
	}
	if s.styleForKey("other").GetFaint() != true {
		t.Error("default key style not used")
	}
	if s.styleForValue("v").GetBold() != true {
		t.Error("value override not used")
	}
	if s.styleForValue("other").GetBold() {
		t.Error("default value style should not be bold")
	}
}

func TestDefaultStylesCoverBaseLevels(t *testing.T) {
	s := DefaultStyles()
	if len(s.Levels) != 5 {
		t.Errorf("expected 5 level styles, got %d", len(s.Levels))
	}
}
