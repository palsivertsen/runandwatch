package main

import (
	"flag"
	"log"
	"runandwatch/executors"
	"runandwatch/filters"
	"runandwatch/watchers"
)

/*
runadnwatch go run main.go
runandwatch bash -c "go run main.go && go test ./..."
*/

func main() {
	var (
		workingDir = flag.String("workingDir", ".", "file or directory to watch")
		logPrefix  = flag.String("logPrefix", "runandwatch: ", "Prefix for runandwatch's internal log statements")
	)
	flag.Parse()

	log.SetPrefix(*logPrefix)
	log.SetFlags(log.Lmicroseconds)

	var cmd []string

	cmd = flag.Args()

	filter := &filters.Git{}
	watcher := watchers.NewFsNotify(*workingDir)
	executor := &executors.Command{
		WorkingDir: *workingDir,
		Cmd:        cmd,
	}

	run(filter, watcher, executor)
}

func run(filter filters.Filter, watcher watchers.Watcher, executor executors.Executor) {
	log.Println("Starting")
	if err := executor.Restart(); err != nil {
		log.Print(err)
	}
	for file := range watcher.Changes() {
		log.Print("Restarting")
		if !filter.Watched(file) {
			continue
		}
		if err := executor.Restart(); err != nil {
			log.Print("runandwatch: error executing command: ", err)
		}
	}
}
