package filters

import (
	"context"
	"os/exec"
	"time"
)

// Git filters files not tracked by Git
type Git struct {
	// The max time to wait for git
	Timeout time.Duration
}

// Watched returns false if file is ignored by Git
func (f *Git) Watched(file string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout())
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "check-ignore", file)
	err := cmd.Run()
	// If command fails, file is not in .gitignore
	return err != nil
}

func (f *Git) timeout() time.Duration {
	if f.Timeout != 0 {
		return f.Timeout
	}
	return time.Second * 5
}
