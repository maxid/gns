[English](README.md) | 简体中文

# gns

一个根命名空间，派生出所有写共享状态的子系统所需的前缀：缓存键、迁移历史表、
CosID 号段、Kafka topic/group 和集群 channel。

> 状态：**实验性脚手架**（v0.0.0-dev）。API 尚不稳定。

## 为什么

两个部署共用 Redis、Kafka 或数据库时不能互相覆盖；迁移时又必须保留既有名称。
`gns` 用一个根命名空间，并允许每个子系统覆盖：

```yaml
namespace:
  root: app
  cache: app
  flyway: app
  cosid: none       # 保留既有 CosID 表的逻辑段名
  cluster: app
  timeSeries: app
```

```go
set, _ := gns.New(cfg)
set.RedisKey("token", "abc")          // app:token:abc
set.FlywayTable("schema_migrations")  // app_schema_migrations
set.Segment("sys_dict")               // cosid: none 时仍是 sys_dict
set.Topic("iothub_raw")               // app.iothub_raw
```

子系统为空时继承 `root`。名称只允许小写字母、数字、`_`、`-`；
`cosid: none` 是保留历史段名的显式逃生口。

## 安装

```bash
go get github.com/maxid/gns
```

## 示例

```bash
go run ./examples/derive
```

## 许可证

MIT。见 [LICENSE](LICENSE)。
