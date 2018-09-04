package executors

// Executor executes commands
type Executor interface {
	Restart() error
}
