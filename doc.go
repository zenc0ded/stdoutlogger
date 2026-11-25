// Package stdoutlogger provides an OpenTelemetry log exporter that prints
// records to a writer (os.Stdout by default) as styled, human readable lines.
//
// Each record is written as one line:
//
//	2006/01/02 15:04:05 INFO <pkg/file.go:12> message key="value" n=3
//
// The source location is built from the code.function, code.filepath and
// code.lineno attributes, which the otelslog bridge sets when source
// logging is enabled. Those attributes are not repeated in the key/value list.
//
// Output is styled with lipgloss. Colors are only used when the writer is a
// terminal.
//
// Use New to create the exporter and pass it to a log SDK processor:
//
//	exp, err := stdoutlogger.New()
//	if err != nil {
//		return err
//	}
//	provider := log.NewLoggerProvider(log.WithProcessor(log.NewSimpleProcessor(exp)))
package stdoutlogger
