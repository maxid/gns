package gns

import "testing"

func TestDerivedNamespaces(t *testing.T) {
	set, err := New(Config{Root: "app"})
	if err != nil {
		t.Fatal(err)
	}
	if set.CacheName() != "app" || set.FlywayTable("schema_migrations") != "app_schema_migrations" {
		t.Fatalf("set = %+v", set)
	}
	if set.CosIDName() != "app" || set.SegmentKey("sys_dict") != "app.sys_dict" {
		t.Fatalf("cosid = %q/%q", set.CosIDName(), set.SegmentKey("sys_dict"))
	}
	if set.RedisKey("user", "1") != "app:user:1" {
		t.Fatalf("key = %q", set.RedisKey("user", "1"))
	}
	if set.Topic("iothub_raw") != "app.iothub_raw" || set.Group("raw") != "app.raw" {
		t.Fatalf("topic/group = %q/%q", set.Topic("iothub_raw"), set.Group("raw"))
	}
	if set.Table("device") != "app_device" {
		t.Fatalf("table = %q", set.Table("device"))
	}
}

func TestOverridesAndNone(t *testing.T) {
	set, err := New(Config{Root: "app", Cache: "app_cache", Flyway: "app_db", CosID: "none", Cluster: "app_cluster", TimeSeries: "app_ts"})
	if err != nil {
		t.Fatal(err)
	}
	if set.CacheName() != "app_cache" || set.FlywayTable("schema_migrations") != "app_db_schema_migrations" {
		t.Fatalf("set = %+v", set)
	}
	if set.CosIDName() != "" || set.SegmentKey("sys_dict") != "sys_dict" {
		t.Fatalf("cosid none = %q/%q", set.CosIDName(), set.SegmentKey("sys_dict"))
	}
	if set.ClusterName() != "app_cluster" {
		t.Fatalf("cluster = %q", set.ClusterName())
	}
}

func TestValidationAndJoin(t *testing.T) {
	for _, cfg := range []Config{
		{Root: "Bad Root"},
		{Root: "app", Cache: "has space"},
		{Root: "app", CosID: "bad.dot"},
		{Root: "app", Flyway: "bad.dot"},
	} {
		if err := cfg.Validate(); err == nil {
			t.Fatalf("expected error for %+v", cfg)
		}
	}
	if _, err := New(Config{Root: ""}); err != nil {
		t.Fatal(err)
	}
	flywayNone, err := New(Config{Root: "app", Flyway: "none"})
	if err != nil {
		t.Fatal(err)
	}
	if got := flywayNone.FlywayTable("flyway_schema_history"); got != "flyway_schema_history" {
		t.Fatalf("flyway none = %q", got)
	}
	if got := Join(":", "a", "", "b", "c:"); got != "a:b:c" {
		t.Fatalf("join = %q", got)
	}
}
