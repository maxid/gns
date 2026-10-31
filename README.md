English | [简体中文](README.zh-CN.md)

# gns

Namespace as an infrastructure primitive: a validated namespace value, context
propagation, and a pluggable store for distributed definitions.

> Status: **experimental scaffold** (v0.0.0-dev). The API is not stable yet.

## Scope

`gns` deliberately knows only three things:

- **Namespace** — a name plus application-defined labels.
- **Context** — `With` / `From` / `MustFrom`.
- **Store** — the SPI used to persist namespace definitions, with an in-memory
  implementation. Redis (or any other backend) is supplied by the application or
  framework, for example through `gcache`; gns has no Redis dependency.

It does **not** derive cache keys, Flyway table names, CosID segments, Kafka
topics or cluster channels. Those naming conventions belong to the library that
owns the subsystem, or to the application configuration.

## Context

```go
ns, _ := gns.New("app-dev")
ctx := gns.With(context.Background(), ns)

current, ok := gns.From(ctx)   // ok == true
name := current.Name           // app-dev
```

A zero namespace is not a scope: `With` leaves the context unchanged and
`From` returns `ok == false`. There is no implicit default namespace.

## Store SPI

```go
type Store interface {
    Load(ctx context.Context, name string) (Namespace, error)
    Save(ctx context.Context, ns Namespace) error
    Delete(ctx context.Context, name string) error
}
```

`gns.NewMemoryStore()` is the in-process implementation. For a shared deployment,
implement the same interface over Redis (a gable/gcache adapter is the intended
path) and let the application choose the backend and its configuration.

```go
store := gns.NewMemoryStore()
_ = store.Save(ctx, gns.Namespace{
    Name:   "app-dev",
    Labels: map[string]string{"env": "dev"},
})

manager, _ := gns.NewManager(store)
ctx, err := manager.Bind(ctx, "app-dev") // missing namespace is an error
```

`Manager.Bind` fails closed: it never substitutes an empty or default namespace.

## Install

```bash
go get github.com/maxid/gns
```

## Example

```bash
go run ./examples/context
```

## License

MIT. See [LICENSE](LICENSE).
