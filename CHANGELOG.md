# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

- Breaking: gns is now a namespace/context/Store primitive. Business naming
  helpers (cache keys, Flyway tables, CosID segments, Kafka names, cluster
  channels) were removed.
- `Namespace` value with labels, context helpers, `Store` SPI, `MemoryStore` and
  `Manager.Bind`.
