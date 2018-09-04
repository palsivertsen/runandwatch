package watchers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type fsNotify struct {
	changes chan string
}

// NewFsNotify returns a Watcher that recorivly monitors for changes in the given directory
func NewFsNotify(root string) Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	if err := recursiveAddDirs(watcher, root); err != nil {
		panic(err)
	}

	changes := make(chan string)

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					panic("runandwatch: events channel closed")
				}
				if isDir(event.Name) {
					switch {
					case event.Op&fsnotify.Create == fsnotify.Create:
						if err := watcher.Add(event.Name); err != nil {
							log.Print("runandwatch: failed to watch dir: ", err)
						}
					case event.Op&fsnotify.Remove == fsnotify.Remove:
						if err := watcher.Remove(event.Name); err != nil {
							log.Print("runandwatch: failed to watch dir: ", err)
						}
					}
				} else if !isFile(event.Name) {
					break
				}
				changes <- event.Name
			case err, ok := <-watcher.Errors:
				if !ok {
					panic("runandwatch: error channel closed")
				}
				log.Print("runandwatch: error watching file: ", err)
			}
		}
	}()
	return &fsNotify{
		changes: changes,
	}
}

// Changes returns any file of directory has changed
func (w *fsNotify) Changes() <-chan string {
	return w.changes
}

func recursiveAddDirs(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, walk(watcher))
}

func walk(watcher *fsnotify.Watcher) filepath.WalkFunc {
	return func(path string, info os.FileInfo, ierr error) error {
		if !info.IsDir() {
			return ierr
		}
		return watcher.Add(path)
	}
}

func isDir(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		if err.(*os.PathError).Err.Error() == "no such file or directory" {
			return false
		}
		log.Panic("runandwatch: could not stat file: ", err)
	}
	return info.IsDir()
}

func isFile(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		if err.(*os.PathError).Err.Error() == "no such file or directory" {
			return false
		}
		log.Panic("runandwatch: could not stat file: ", err)
	}
	return info.Mode().IsRegular()
}
