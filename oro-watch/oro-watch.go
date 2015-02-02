package main

import (
	"gopkg.in/fsnotify.v1"
	"log"
	"path/filepath"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	curpath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(curpath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
