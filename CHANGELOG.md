# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Initial scaffold: root namespace with per-subsystem overrides, cache keys,
  Flyway tables, CosID segments, Kafka names and cluster channels.
- `Segment` renamed to `SegmentKey` (fully-qualified key); `CosIDName()` is the
  allocator namespace and returns `""` for `cosid: none`.
- `flyway: none` keeps the base history table name without a prefix.
