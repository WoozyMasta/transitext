# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][],
and this project adheres to [Semantic Versioning][].

<!--
## Unreleased

### Added
### Changed
### Removed
-->

## [0.1.2][] - 2026-03-07

### Added

* Added schema-friendly config package
  `github.com/woozymasta/transitext/config` with aggregated wrapper
  and provider options for documentation and validation flows.

### Changed

* Added `jsonschema` constraints for core request/batch contracts, wrapper
  options, and built-in provider options.
* Updated public option field comments to be clearer for end-user configuration.
* Fixed package documentation example to match current
  `RetryOptions` and `RateLimitOptions` API.

[0.1.2]: https://github.com/WoozyMasta/transitext/compare/v0.1.1...v0.1.2

## [0.1.1][] - 2026-03-06

### Added

* Added `langmap` for language normalization and provider support checks, so
  user-facing language names/aliases can be resolved to provider-ready codes
  before translation.

[0.1.1]: https://github.com/WoozyMasta/transitext/compare/v0.1.0...v0.1.1

## [0.1.0][] - 2026-03-05

### Added

* First public release

[0.1.0]: https://github.com/WoozyMasta/transitext/tree/v0.1.0

<!--links-->
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
