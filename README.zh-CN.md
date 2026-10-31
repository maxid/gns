[English](README.md) | 简体中文

# gns

把 namespace 当作基础设施原语：校验后的 namespace 值、context 传递，以及可插拔的
分布式定义存储。

> 状态：**实验性脚手架**（v0.0.0-dev）。API 尚不稳定。

## 边界

`gns` 只做三件事：

- **Namespace**：名称 + 应用自定义 labels；
- **Context**：`With` / `From` / `MustFrom`；
- **Store**：持久化 namespace 定义的 SPI，自带内存实现；Redis 等后端由应用或框架
  （例如通过 `gcache`）实现，gns 不依赖 Redis。

它**不**派生 cache key、Flyway 表名、CosID 段名、Kafka topic 或 cluster channel。
这些命名约定属于拥有该子系统的库，或应用自己的配置。

## Context

```go
ns, _ := gns.New("app-dev")
ctx := gns.With(context.Background(), ns)

current, ok := gns.From(ctx)   // ok == true
name := current.Name           // app-dev
```

零值 namespace 不是有效作用域：`With` 不会修改 context，`From` 返回
`ok == false`。gns 不提供隐式默认 namespace。

## Store SPI

```go
type Store interface {
    Load(ctx context.Context, name string) (Namespace, error)
    Save(ctx context.Context, ns Namespace) error
    Delete(ctx context.Context, name string) error
    List(ctx context.Context) ([]Namespace, error)
}
```

`gns.NewMemoryStore()` 是进程内实现。分布式部署实现同一接口即可（推荐用
gable/gcache 适配 Redis），后端与配置项由应用决定。

```go
store := gns.NewMemoryStore()
_ = store.Save(ctx, gns.Namespace{
    Name:   "app-dev",
    Labels: map[string]string{"env": "dev"},
})

manager, _ := gns.NewManager(store)
ctx, err := manager.Bind(ctx, "app-dev") // 不存在就是错误，不回落
```

`Manager.Bind` fail closed：不会用空 namespace 或默认值替代。

## 安装

```bash
go get github.com/maxid/gns
```

## 示例

```bash
go run ./examples/context
```

## 许可证

MIT。见 [LICENSE](LICENSE)。
