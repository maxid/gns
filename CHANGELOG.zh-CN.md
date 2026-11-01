[English](CHANGELOG.md) | 简体中文

# 变更日志

本文件记录本项目的重要变更。

## [Unreleased]

- 破坏性变更：gns 现在是 namespace/context/Store 原语。业务命名辅助（cache key、Flyway 表、CosID 段、Kafka 名称、cluster channel）已移除。
- `Namespace` 值/labels、context 辅助、`Store` SPI、`MemoryStore` 与 `Manager.Bind`。
