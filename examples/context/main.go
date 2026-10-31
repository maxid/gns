package main

import (
	"context"
	"fmt"

	"github.com/maxid/gns"
)

func main() {
	ns, _ := gns.New("app-dev")
	ctx := gns.With(context.Background(), ns)
	current := gns.MustFrom(ctx)
	fmt.Println("context:", current.Name)

	store := gns.NewMemoryStore()
	_ = store.Save(ctx, gns.Namespace{Name: "app-dev", Labels: map[string]string{"env": "dev"}})
	manager, _ := gns.NewManager(store)
	bound, _ := manager.Bind(ctx, "app-dev")
	fmt.Println("stored:", gns.MustFrom(bound).Name)
}
