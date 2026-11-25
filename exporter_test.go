package stdoutlogger

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// emit sends one record through a real SDK pipeline into the exporter and
// returns whatever was written to the configured writer.
func emit(t *testing.T, build func() otellog.Record, opts ...Option) string {
	t.Helper()
	var buf bytes.Buffer
	exp, err := New(append([]Option{WithWriter(&buf)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	p := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	p.Logger("test").Emit(context.Background(), build())
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func rec(sev otellog.Severity, body string, attrs ...otellog.KeyValue) func() otellog.Record {
	return func() otellog.Record {
		var r otellog.Record
		r.SetTimestamp(time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC))
		r.SetSeverity(sev)
		r.SetBody(otellog.StringValue(body))
		r.AddAttributes(attrs...)
		return r
	}
}

func TestWritesToConfiguredWriter(t *testing.T) {
	out := emit(t, rec(otellog.SeverityInfo, "hello"), WithoutSource())
	if !strings.Contains(out, "hello") {
		t.Fatalf("expected output in configured writer, got %q", out)
	}
}

func TestBasicFormat(t *testing.T) {
	out := emit(t, rec(otellog.SeverityInfo, "hello", otellog.String("k", "v"), otellog.Int("n", 3)), WithoutSource())
	want := "2025/01/02 03:04:05 INFO hello k=\"v\" n=3\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestWithoutTimestamps(t *testing.T) {
	out := emit(t, rec(otellog.SeverityWarn, "x"), WithoutSource(), WithoutTimestamps())
	if out != "WARN x\n" {
		t.Fatalf("got %q", out)
	}
}

func TestSourceFromAttributes(t *testing.T) {
	out := emit(t, rec(otellog.SeverityInfo, "x",
		otellog.String("code.filepath", "/a/b/file.go"),
		otellog.Int("code.lineno", 42),
		otellog.String("code.function", "github.com/foo/bar.Func"),
	), WithoutTimestamps())
	want := "INFO <bar/file.go:42> x\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestSourceFunctionMain(t *testing.T) {
	out := emit(t, rec(otellog.SeverityInfo, "x",
		otellog.String("code.function", "main.main"),
		otellog.String("code.filepath", "m.go"),
		otellog.Int("code.lineno", 1),
	), WithoutTimestamps())
	if out != "INFO <main/m.go:1> x\n" {
		t.Fatalf("got %q", out)
	}
}

func TestSourceMissingAttributesOmitsLocation(t *testing.T) {
	out := emit(t, rec(otellog.SeverityInfo, "x"), WithoutTimestamps())
	if out != "INFO x\n" {
		t.Fatalf("got %q", out)
	}
}

func TestEmptyFunctionDoesNotPanic(t *testing.T) {
	emit(t, rec(otellog.SeverityInfo, "x", otellog.String("code.function", "")), WithoutTimestamps())
}

func TestAllSeveritiesRenderALevel(t *testing.T) {
	for sev := otellog.SeverityTrace; sev <= otellog.SeverityFatal4; sev++ {
		out := emit(t, rec(sev, "m"), WithoutSource(), WithoutTimestamps())
		if strings.HasPrefix(out, " ") || !strings.HasSuffix(out, " m\n") || len(out) <= len(" m\n") {
			t.Errorf("severity %v: bad output %q", sev, out)
		}
	}
}

func TestUndefinedSeverity(t *testing.T) {
	out := emit(t, rec(otellog.SeverityUndefined, "m"), WithoutSource(), WithoutTimestamps())
	if strings.HasPrefix(out, " ") {
		t.Fatalf("empty level column: %q", out)
	}
}

func TestZeroTimestampNotPrinted(t *testing.T) {
	out := emit(t, func() otellog.Record {
		var r otellog.Record
		r.SetSeverity(otellog.SeverityInfo)
		r.SetBody(otellog.StringValue("m"))
		return r
	}, WithoutSource())
	if strings.Contains(out, "0001/01/01") {
		t.Fatalf("zero time leaked: %q", out)
	}
}

func TestNonStringBody(t *testing.T) {
	out := emit(t, func() otellog.Record {
		var r otellog.Record
		r.SetSeverity(otellog.SeverityInfo)
		r.SetBody(otellog.Int64Value(7))
		return r
	}, WithoutSource(), WithoutTimestamps())
	if !strings.Contains(out, "7") {
		t.Fatalf("non-string body dropped: %q", out)
	}
}

func TestSingleWritePerRecord(t *testing.T) {
	w := &countingWriter{}
	exp, _ := New(WithWriter(w))
	p := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	p.Logger("t").Emit(context.Background(), rec(otellog.SeverityInfo, "m")())
	if w.n != 1 {
		t.Fatalf("want 1 write per record (atomic lines), got %d", w.n)
	}
}

type countingWriter struct{ n int }

func (c *countingWriter) Write(b []byte) (int, error) { c.n++; return len(b), nil }

func TestExportCancelledContext(t *testing.T) {
	exp, _ := New(WithWriter(&bytes.Buffer{}))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := exp.Export(ctx, []sdklog.Record{{}}); err == nil {
		t.Fatal("expected context error")
	}
}

func TestWriteErrorReturned(t *testing.T) {
	exp, _ := New(WithWriter(errWriter{}))
	if err := exp.Export(context.Background(), []sdklog.Record{{}}); err == nil {
		t.Fatal("expected write error to be returned")
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, bytes.ErrTooLarge }

func TestStringifyEscapesStrings(t *testing.T) {
	got := stringify(otellog.StringValue("a\"b\nc"))
	if got != `"a\"b\nc"` {
		t.Fatalf("got %s", got)
	}
}

func TestStringifyKinds(t *testing.T) {
	cases := []struct {
		v    otellog.Value
		want string
	}{
		{otellog.Int64Value(-1), "-1"},
		{otellog.Float64Value(1.5), "1.5"},
		{otellog.BoolValue(true), "true"},
		{otellog.Value{}, "nil"},
		{otellog.SliceValue(otellog.IntValue(1), otellog.StringValue("a")), `[1,"a"]`},
		{otellog.MapValue(otellog.Int("a", 1), otellog.Bool("b", false)), `{a:1,b:false}`},
		{otellog.BytesValue([]byte{1, 2}), "[1 2]"},
	}
	for _, c := range cases {
		if got := stringify(c.v); got != c.want {
			t.Errorf("got %q want %q", got, c.want)
		}
	}
}

func TestNilWriterDoesNotPanic(t *testing.T) {
	exp, err := New(WithWriter(nil))
	if err != nil {
		return
	}
	_ = exp.Export(context.Background(), []sdklog.Record{{}})
}

func TestShutdownAndFlush(t *testing.T) {
	exp, _ := New()
	if err := exp.ForceFlush(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := exp.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}
