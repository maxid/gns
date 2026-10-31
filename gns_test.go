package gns

import (
	"context"
	"errors"
	"testing"
)

func TestNamespaceValidation(t *testing.T) {
	for _, name := range []string{"app", "app-dev", "team.alpha", "A_1"} {
		if _, err := New(name); err != nil {
			t.Fatalf("New(%q): %v", name, err)
		}
	}
	for _, name := range []string{"", "-bad", "has space", "a/b", string(make([]byte, 200))} {
		if _, err := New(name); err == nil {
			t.Fatalf("New(%q) succeeded", name)
		}
	}
}

func TestContextRoundTrip(t *testing.T) {
	ns, _ := New("app-dev")
	ns.Labels = map[string]string{"env": "dev"}
	ctx := With(context.Background(), ns)
	got, ok := From(ctx)
	if !ok || got.Name != "app-dev" || got.Labels["env"] != "dev" {
		t.Fatalf("From = %+v, %v", got, ok)
	}
	got.Labels["env"] = "mutated"
	again, _ := From(ctx)
	if again.Labels["env"] != "dev" {
		t.Fatal("context namespace labels were not isolated")
	}

	if _, ok := From(context.Background()); ok {
		t.Fatal("empty context returned a namespace")
	}
	if ctx := With(context.Background(), Namespace{}); ctx == nil {
		t.Fatal("With returned nil")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("MustFrom did not panic")
		}
	}()
	_ = MustFrom(context.Background())
}

func TestMemoryStore(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	ns, _ := New("app")
	ns.Labels = map[string]string{"owner": "team"}
	if err := store.Save(ctx, ns); err != nil {
		t.Fatal(err)
	}
	got, err := store.Load(ctx, "app")
	if err != nil || got.Labels["owner"] != "team" {
		t.Fatalf("Load = %+v, %v", got, err)
	}
	got.Labels["owner"] = "changed"
	again, _ := store.Load(ctx, "app")
	if again.Labels["owner"] != "team" {
		t.Fatal("store returned its internal map")
	}
	if _, err := store.Load(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load missing = %v", err)
	}

	second, _ := New("other")
	if err := store.Save(ctx, second); err != nil {
		t.Fatal(err)
	}
	list, err := store.List(ctx)
	if err != nil || len(list) != 2 || list[0].Name != "app" || list[1].Name != "other" {
		t.Fatalf("List = %+v, %v", list, err)
	}
	if err := store.Delete(ctx, "app"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(ctx, "app"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Load after delete = %v", err)
	}
	if err := store.Delete(ctx, "app"); err != nil {
		t.Fatalf("delete missing = %v", err)
	}
}

func TestManagerBind(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryStore()
	ns, _ := New("app")
	if err := store.Save(ctx, ns); err != nil {
		t.Fatal(err)
	}
	mgr, err := NewManager(store)
	if err != nil {
		t.Fatal(err)
	}
	bound, err := mgr.Bind(ctx, "app")
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := From(bound); !ok || got.Name != "app" {
		t.Fatalf("bound = %+v, %v", got, ok)
	}
	if _, err := mgr.Bind(ctx, "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Bind missing = %v", err)
	}
	if _, err := NewManager(nil); err == nil {
		t.Fatal("nil store accepted")
	}
}
