package main

import (
	"fmt"

	"github.com/maxid/gns"
)

func main() {
	set, err := gns.New(gns.Config{Root: "app", CosID: "none"})
	if err != nil {
		panic(err)
	}
	fmt.Println("cache:", set.CacheName())
	fmt.Println("flyway:", set.FlywayTable("schema_migrations"))
	fmt.Println("cosid:", set.Segment("sys_dict"))
	fmt.Println("cluster:", set.ClusterName())
	fmt.Println("topic:", set.Topic("iothub_raw"))
	fmt.Println("redis:", set.RedisKey("token", "abc"))
}
