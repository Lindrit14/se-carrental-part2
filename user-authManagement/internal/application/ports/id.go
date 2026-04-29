package ports

// IDGenerator produces opaque, globally-unique identifiers (e.g. UUIDv7).
type IDGenerator interface {
	New() string
}
