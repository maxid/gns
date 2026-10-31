// Package gns carries a namespace value in a context and persists namespace
// definitions through a pluggable Store.
//
// It deliberately does not derive cache keys, migration table names, CosID
// segments, Kafka topics or any other business naming convention. Those belong
// to the library or application that owns the subsystem.
//
// Status: experimental scaffold (v0.0.0-dev). The API is not stable yet.
package gns
