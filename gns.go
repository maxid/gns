package gns

import (
	"fmt"
	"regexp"
	"strings"
)

// Config is the namespace configuration. Empty subsystem values inherit Root.
type Config struct {
	Root       string `yaml:"root" json:"root"`
	Cache      string `yaml:"cache,omitempty" json:"cache,omitempty"`
	Flyway     string `yaml:"flyway,omitempty" json:"flyway,omitempty"`
	CosID      string `yaml:"cosid,omitempty" json:"cosid,omitempty"` // "none" keeps logical segment names
	Cluster    string `yaml:"cluster,omitempty" json:"cluster,omitempty"`
	TimeSeries string `yaml:"timeSeries,omitempty" json:"timeSeries,omitempty"`
}

// DefaultConfig returns a configuration rooted at root.
func DefaultConfig(root string) Config {
	if strings.TrimSpace(root) == "" {
		root = "gable"
	}
	return Config{Root: root}
}

var nameRE = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)

// Normalize fills empty subsystem values from Root.
func (c Config) Normalize() Config {
	if strings.TrimSpace(c.Root) == "" {
		c.Root = "gable"
	}
	c.Root = strings.TrimSpace(c.Root)
	for _, target := range []*string{&c.Cache, &c.Flyway, &c.CosID, &c.Cluster, &c.TimeSeries} {
		*target = strings.TrimSpace(*target)
		if *target == "" {
			*target = c.Root
		}
	}
	return c
}

// Validate checks every namespace that becomes part of a key or table name.
func (c Config) Validate() error {
	c = c.Normalize()
	if !nameRE.MatchString(c.Root) {
		return fmt.Errorf("gns: root %q must match ^[a-z][a-z0-9_-]*$", c.Root)
	}
	for label, value := range map[string]string{
		"cache": c.Cache, "flyway": c.Flyway, "cluster": c.Cluster, "timeSeries": c.TimeSeries,
	} {
		if !nameRE.MatchString(value) {
			return fmt.Errorf("gns: %s %q must match ^[a-z][a-z0-9_-]*$", label, value)
		}
	}
	if c.CosID != "none" && !nameRE.MatchString(c.CosID) {
		return fmt.Errorf("gns: cosid %q must match ^[a-z][a-z0-9_-]*$ or be \"none\"", c.CosID)
	}
	return nil
}

// Set is a validated, derived namespace set.
type Set struct {
	Root       string
	Cache      string
	Flyway     string
	CosID      string
	Cluster    string
	TimeSeries string
}

// New validates cfg and returns the derived set.
func New(cfg Config) (Set, error) {
	cfg = cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return Set{}, err
	}
	return Set{
		Root:       cfg.Root,
		Cache:      cfg.Cache,
		Flyway:     cfg.Flyway,
		CosID:      cfg.CosID,
		Cluster:    cfg.Cluster,
		TimeSeries: cfg.TimeSeries,
	}, nil
}

// Default returns the default set.
func Default() Set {
	set, _ := New(DefaultConfig("gable"))
	return set
}

// CacheName is the namespace passed to gcache.
func (s Set) CacheName() string { return s.Cache }

// FlywayName is the prefix of migration history tables.
func (s Set) FlywayName() string { return s.Flyway }

// FlywayTable renders the migration history table for base.
func (s Set) FlywayTable(base string) string { return s.Flyway + "_" + sanitize(base) }

// CosIDName is the segment prefix; "none" keeps logical names.
func (s Set) CosIDName() string { return s.CosID }

// ClusterName is the cluster/gossip namespace.
func (s Set) ClusterName() string { return s.Cluster }

// TimeSeriesName is the time-series table prefix.
func (s Set) TimeSeriesName() string { return s.TimeSeries }

// Segment renders a CosID segment name. "none" keeps the logical name so a
// deployment that predates namespacing keeps its existing ID ranges.
func (s Set) Segment(logical string) string {
	if s.CosID == "none" || logical == "" {
		return logical
	}
	return s.CosID + "." + logical
}

// Table renders a time-series (or other engine-specific) table name.
func (s Set) Table(base string) string { return s.TimeSeries + "_" + sanitize(base) }

// RedisKey prefixes a logical key with the cache namespace.
func (s Set) RedisKey(parts ...string) string {
	clean := make([]string, 0, len(parts)+1)
	clean = append(clean, s.Cache)
	for _, part := range parts {
		part = strings.Trim(part, ":")
		if part != "" {
			clean = append(clean, part)
		}
	}
	return strings.Join(clean, ":")
}

// Join joins non-empty parts with separator, trimming separators at the edges.
func Join(separator string, parts ...string) string {
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, separator)
		if part != "" {
			clean = append(clean, part)
		}
	}
	return strings.Join(clean, separator)
}

// Topic and Group render Kafka names using the cluster namespace.
func (s Set) Topic(logical string) string { return s.Cluster + "." + logical }
func (s Set) Group(logical string) string { return s.Cluster + "." + logical }

// Channel renders a Redis pub/sub or internal cluster channel.
func (s Set) Channel(logical string) string { return s.Cluster + ":" + logical }

func sanitize(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, char := range value {
		switch {
		case char >= 'a' && char <= 'z', char >= '0' && char <= '9', char == '_':
			builder.WriteRune(char)
		default:
			builder.WriteByte('_')
		}
	}
	out := strings.Trim(builder.String(), "_")
	if out == "" {
		return "default"
	}
	return out
}
