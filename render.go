package stdoutlogger

import (
	"fmt"
	"os"
	"path"
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
		strval = fmt.Sprintf("\"%s\"", v.AsString())
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

func (e *exporter) renderRecord(record sdklog.Record) {
	st := e.style
	re := e.renderer
	// try to track source from attributes to better output it
	var module string = "main"
	var filename string
	var lineno int64

	attrSep := st.Separator.Renderer(re).Render("=")

	// build msg string, start with record body and append attrs
	var msg strings.Builder
	msg.WriteString(st.Message.Renderer(re).Render(record.Body().AsString()))

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
			// FIXME: handle empty split?
			module = path.Base(strings.Split(v.AsString(), ".")[0])

		default:
			msg.WriteString(" ")
			msg.WriteString(st.styleForKey(k).Renderer(re).Render(k))
			msg.WriteString(attrSep)
			msg.WriteString(st.styleForValue(k).Renderer(re).Render(stringify(v)))
		}
		return true
	})

	if e.timestamps {
		os.Stdout.WriteString(st.Timestamp.Renderer(re).Render(record.Timestamp().Format("2006/01/02 15:04:05")) + " ")
	}
	os.Stdout.WriteString(st.Levels[record.Severity()].Render(record.Severity().String()))
	os.Stdout.WriteString(" ")
	if e.source {
		source := fmt.Sprintf("<%s/%s:%d>", module, filename, lineno)
		os.Stdout.WriteString(st.Caller.Renderer(re).Render(source))
		os.Stdout.WriteString(" ")
	}
	os.Stdout.WriteString(msg.String())

	os.Stdout.WriteString("\n")
}
