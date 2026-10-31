English | [简体中文](README.zh-CN.md)

# gns

One namespace, derived for every subsystem that writes shared state: cache
keys, migration history tables, CosID segments, Kafka topics/groups and cluster
channels.

> Status: **experimental scaffold** (v0.0.0-dev). The API is not stable yet.

## Why

Two deployments sharing Redis, Kafka or a database must not collide. They also
must keep their existing names when they migrate. `gns` keeps one root
namespace and lets each subsystem override it:

```yaml
namespace:
  root: app
  cache: app
  flyway: app
  cosid: none       # keep logical segment names for an existing CosID table
  cluster: app
  timeSeries: app
```

```go
set, _ := gns.New(cfg)
set.RedisKey("token", "abc")          // app:token:abc
set.FlywayTable("schema_migrations")  // app_schema_migrations
set.Segment("sys_dict")               // sys_dict when cosid: none
set.Topic("iothub_raw")               // app.iothub_raw
```

Empty subsystem values inherit `root`. Names are lowercase letters, digits,
`_` and `-`; `cosid: none` is the explicit escape for legacy segment names.

## Install

```bash
go get github.com/maxid/gns
```

## Example

```bash
go run ./examples/derive
```

## License

MIT. See [LICENSE](LICENSE).
