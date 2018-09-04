package executors

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

// Command executes commands
type Command struct {
	WorkingDir string
	Cmd        []string
	mux        sync.Mutex
	cancel     context.CancelFunc
}

// Restart the command
func (e *Command) Restart() error {
	e.mux.Lock()
	defer e.mux.Unlock()

	if len(e.Cmd) == 0 {
		return errors.New("no command")
	}

	if e.cancel != nil {
		e.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	cmd := exec.CommandContext(ctx, e.Cmd[0], e.Cmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not run command: %s", err.Error())
	}

	go cmd.Wait()
	return nil
}
