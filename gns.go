package gns

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Errors returned by gns.
var (
	ErrInvalidNamespace = errors.New("gns: invalid namespace")
	ErrNotFound         = errors.New("gns: namespace not found")
)

var nameRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// Namespace is a deployment or application namespace. Labels are application
// defined metadata; gns never interprets them.
type Namespace struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// New validates name and returns a namespace without labels.
func New(name string) (Namespace, error) {
	ns := Namespace{Name: strings.TrimSpace(name)}
	if err := ns.Validate(); err != nil {
		return Namespace{}, err
	}
	return ns, nil
}

// Validate checks the namespace name.
func (n Namespace) Validate() error {
	if !nameRE.MatchString(n.Name) {
		return fmt.Errorf("%w: %q (want letters, digits, '.', '_' or '-')", ErrInvalidNamespace, n.Name)
	}
	return nil
}

// String returns the namespace name.
func (n Namespace) String() string { return n.Name }

// IsZero reports whether the namespace has no name.
func (n Namespace) IsZero() bool { return n.Name == "" }

// Clone returns a copy with an independent label map.
func (n Namespace) Clone() Namespace {
	out := Namespace{Name: n.Name}
	if len(n.Labels) > 0 {
		out.Labels = make(map[string]string, len(n.Labels))
		for k, v := range n.Labels {
			out.Labels[k] = v
		}
	}
	return out
}

type contextKey struct{}

// With returns a child context carrying ns. A zero namespace leaves ctx
// unchanged; callers must not use an empty namespace as a real scope.
func With(ctx context.Context, ns Namespace) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if ns.IsZero() {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, ns.Clone())
}

// From returns the namespace carried by ctx.
func From(ctx context.Context) (Namespace, bool) {
	if ctx == nil {
		return Namespace{}, false
	}
	ns, ok := ctx.Value(contextKey{}).(Namespace)
	if !ok || ns.IsZero() {
		return Namespace{}, false
	}
	return ns.Clone(), true
}

// MustFrom returns the namespace or panics when the context has none.
func MustFrom(ctx context.Context) Namespace {
	ns, ok := From(ctx)
	if !ok {
		panic("gns: context has no namespace")
	}
	return ns
}

// Store is the namespace persistence SPI. Implementations must be safe for
// concurrent use; a Redis-backed implementation can live in a deployment or
// framework package, so gns itself stays free of a Redis dependency.
type Store interface {
	// Load returns the namespace stored under name.
	Load(ctx context.Context, name string) (Namespace, error)
	// Save creates or replaces a namespace.
	Save(ctx context.Context, ns Namespace) error
	// Delete removes a namespace; deleting a missing name is not an error.
	Delete(ctx context.Context, name string) error
}

// MemoryStore is the in-process Store implementation. A zero value is ready to
// use.
type MemoryStore struct {
	mu     sync.RWMutex
	values map[string]Namespace
}

// NewMemoryStore returns an empty in-memory store.
func NewMemoryStore() *MemoryStore { return &MemoryStore{} }

// Load implements Store.
func (s *MemoryStore) Load(_ context.Context, name string) (Namespace, error) {
	if !nameRE.MatchString(name) {
		return Namespace{}, fmt.Errorf("%w: %q", ErrInvalidNamespace, name)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ns, ok := s.values[name]
	if !ok {
		return Namespace{}, fmt.Errorf("%w: %s", ErrNotFound, name)
	}
	return ns.Clone(), nil
}

// Save implements Store.
func (s *MemoryStore) Save(_ context.Context, ns Namespace) error {
	if err := ns.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.values == nil {
		s.values = make(map[string]Namespace)
	}
	s.values[ns.Name] = ns.Clone()
	return nil
}

// Delete implements Store.
func (s *MemoryStore) Delete(_ context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, name)
	return nil
}

// List implements Store.
func (s *MemoryStore) List(_ context.Context) ([]Namespace, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Namespace, 0, len(s.values))
	for _, ns := range s.values {
		out = append(out, ns.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// Manager resolves persisted namespaces and binds them to contexts.
type Manager struct {
	store Store
}

// NewManager creates a manager over store.
func NewManager(store Store) (*Manager, error) {
	if store == nil {
		return nil, errors.New("gns: nil store")
	}
	return &Manager{store: store}, nil
}

// Load resolves a namespace from the store.
func (m *Manager) Load(ctx context.Context, name string) (Namespace, error) {
	return m.store.Load(ctx, name)
}

// Bind resolves name and returns a context carrying it. A missing namespace is
// an error, never a silent fallback.
func (m *Manager) Bind(ctx context.Context, name string) (context.Context, error) {
	ns, err := m.store.Load(ctx, name)
	if err != nil {
		return nil, err
	}
	return With(ctx, ns), nil
}

// Save persists a namespace.
func (m *Manager) Save(ctx context.Context, ns Namespace) error {
	return m.store.Save(ctx, ns)
}
