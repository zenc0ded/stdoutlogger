package stdoutlogger

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// simple value to string extraction for types, no style is currently applied to these
func stringify(v log.Value) string {
	strval := ""
	// just do something sane for now
	switch v.Kind() {
	case log.KindString:
		// mirror charm and quote string values
		strval = strconv.Quote(v.AsString())
	case log.KindInt64:
		strval = fmt.Sprint(v.AsInt64())
	case log.KindFloat64:
		strval = fmt.Sprint(v.AsFloat64())
	case log.KindBool:
		strval = fmt.Sprint(v.AsBool())
	case log.KindBytes:
		strval = fmt.Sprint(v.AsBytes())
	case log.KindMap:
		m := v.AsMap()
		strval = "{"
		for i, kv := range m {
			if i > 0 {
				strval += ","
			}
			strval += fmt.Sprintf("%s%s%s", kv.Key, ":", stringify(kv.Value))
		}
		strval += "}"
	case log.KindSlice:
		s := v.AsSlice()
		strval = "["
		for i, e := range s {
			if i > 0 {
				strval += ","
			}
			strval += stringify(e)
		}
		strval += "]"
	case log.KindEmpty:
		strval = "nil"
	default:
		strval = "INVALID_TYPE"
	}
	return strval
}

// levelText renders the severity using the configured style. Severities
// without an explicit style (e.g. INFO2, TRACE) fall back to the style of the
// base level for their range.
func (e *exporter) levelText(sev log.Severity) string {
	st := e.style
	if style, ok := st.Levels[sev]; ok {
		return style.Renderer(e.renderer).Render()
	}
	var base log.Severity
	switch {
	case sev >= log.SeverityFatal:
		base = log.SeverityFatal
	case sev >= log.SeverityError:
		base = log.SeverityError
	case sev >= log.SeverityWarn:
		base = log.SeverityWarn
	case sev >= log.SeverityInfo:
		base = log.SeverityInfo
	case sev >= log.SeverityTrace:
		base = log.SeverityDebug
	default:
		return sev.String()
	}
	if style, ok := st.Levels[base]; ok {
		return style.Renderer(e.renderer).Render()
	}
	return sev.String()
}

// moduleFromFunction extracts the package name from a fully qualified function
// name such as "github.com/foo/bar.(*T).Method" -> "bar".
func moduleFromFunction(fn string) string {
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}
	if i := strings.Index(fn, "."); i >= 0 {
		fn = fn[:i]
	}
	return fn
}

// renderRecord formats a record as a single newline terminated line.
func (e *exporter) renderRecord(record sdklog.Record) string {
	st := e.style
	re := e.renderer
	// try to track source from attributes to better output it
	var module, filename string
	var lineno int64

	attrSep := st.Separator.Renderer(re).Render("=")

	// build msg string, start with record body and append attrs
	body := record.Body()
	bodyText := stringify(body)
	if body.Kind() == log.KindString {
		bodyText = body.AsString()
	}
	var msg strings.Builder
	msg.WriteString(st.Message.Renderer(re).Render(bodyText))

	// pull out certain attributes for source location, render the rest
	record.WalkAttributes(func(kv log.KeyValue) bool {
		k := kv.Key
		v := kv.Value
		// NOTE: these attrs are used by the otelslog bridge, are there others used as standard?
		switch k {
		case string(semconv.CodeFilepathKey):
			filename = path.Base(v.AsString())

		case string(semconv.CodeLineNumberKey):
			lineno = v.AsInt64()

		case string(semconv.CodeFunctionKey):
			module = moduleFromFunction(v.AsString())

		default:
			msg.WriteString(" ")
			msg.WriteString(st.styleForKey(k).Renderer(re).Render(k))
			msg.WriteString(attrSep)
			msg.WriteString(st.styleForValue(k).Renderer(re).Render(stringify(v)))
		}
		return true
	})

	var line strings.Builder
	if e.timestamps {
		ts := record.Timestamp()
		if ts.IsZero() {
			ts = record.ObservedTimestamp()
		}
		if !ts.IsZero() {
			line.WriteString(st.Timestamp.Renderer(re).Render(ts.Format("2006/01/02 15:04:05")))
			line.WriteString(" ")
		}
	}
	line.WriteString(e.levelText(record.Severity()))
	line.WriteString(" ")
	if e.source && (module != "" || filename != "") {
		var source string
		switch {
		case filename == "":
			source = fmt.Sprintf("<%s>", module)
		case module == "":
			source = fmt.Sprintf("<%s:%d>", filename, lineno)
		default:
			source = fmt.Sprintf("<%s/%s:%d>", module, filename, lineno)
		}
		line.WriteString(st.Caller.Renderer(re).Render(source))
		line.WriteString(" ")
	}
	line.WriteString(msg.String())
	line.WriteString("\n")
	return line.String()
}
