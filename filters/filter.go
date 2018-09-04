package filters

// A Filter for files
type Filter interface {
	Watched(file string) bool
}
