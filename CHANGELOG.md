# Changelog

## Unreleased

### Fixed
- Output goes to the writer given to `WithWriter`. It was always written to `os.Stdout`.
- `Export` returns write errors.
- Each record is written with one `Write` call.
- Package name is parsed correctly from `code.function` for import paths.
- No source location is printed when the record has none.
- Severities without their own style (such as `INFO2`) use the style of their base level.
- Zero timestamps are not printed. The observed timestamp is used if set.
- Non-string bodies are printed.
- String values are quoted and escaped.
- `WithWriter(nil)` falls back to `os.Stdout`.
- Style maps are no longer shared between exporters.

### Changed
- Go 1.26 or newer is required.
- OpenTelemetry dependencies updated. `otel/log` and `otel/sdk/log` stay at v0.20.0.

### Added
- Package documentation, an example and tests.
- MIT license.
