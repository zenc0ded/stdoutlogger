# Changelog

## Unreleased

### Fixed
- Output goes to the writer given to `WithWriter`. It was always written to `os.Stdout`.
- `Export` returns write errors.
- Each record is written with one `Write` call.
- Package name is parsed correctly from the function name for import paths.
- Source location attributes using the current semantic convention names (`code.file.path`, `code.function.name`, `code.line.number`) are recognized. Before, they were printed as regular attributes.
- No source location is printed when the record has none.
- Severities without their own style (such as `INFO2`) use the style of their base level.
- Zero timestamps are not printed. The observed timestamp is used if set.
- Non-string bodies are printed.
- String values are quoted and escaped.
- `WithWriter(nil)` falls back to `os.Stdout`.
- Style maps are no longer shared between exporters.

### Changed
- Go 1.26 or newer is required.
- Moved to the stable `otel/log` and `otel/sdk/log` v1.47 API, which uses `attribute.Value`. This also fixes GO-2026-6615 (BatchProcessor busy-spin).
- All dependencies updated. `example/` builds against the local module through a `replace` directive.

### Added
- Typed slice values (bool, int64, float64, string) are printed.
- Package documentation, an example and tests.
- MIT license.
