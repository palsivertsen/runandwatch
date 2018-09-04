package watchers

// A Watcher watches filesystem for changes
type Watcher interface {
	Changes() <-chan string
}
