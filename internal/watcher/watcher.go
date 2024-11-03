package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type ChanPayload struct {
	Error    error
	Filepath string
}

func WatchDir(dir string, channel chan<- ChanPayload) error {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer func() {
		err := watcher.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					channel <- ChanPayload{Error: err}
				} else {
					if event.Has(fsnotify.Create) {
						channel <- ChanPayload{Error: nil, Filepath: event.Name}
					}
				}
			case err, _ := <-watcher.Errors:
				channel <- ChanPayload{Error: err}
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	<-make(chan struct{})
	return nil
}
