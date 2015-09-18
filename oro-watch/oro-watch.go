package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Rebuild(configpath string) error {
	return nil
}

func main() {

	var pageDir string
	pageDir = os.Args[1]
	if len(pageDir) == 0 {
		panic("Usage: oro-watch *directory*")
	}

	var configFiles []string

	// Add yaml files to watchlist
	var files []os.FileInfo
	files, err := ioutil.ReadDir(pageDir)
	if err != nil {
		fmt.Println("Error", err)
		panic(err)
	}
	var fnm string
	for _, file := range files {
		fnm = file.Name()
		if filepath.Ext(fnm) == ".yaml" {
			log.Println("Watching", fnm)
			configFiles = append(configFiles, fnm)
		}
	}

	// Start a watcher

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

				if event.Op&fsnotify.Create == fsnotify.Create {

					if filepath.Ext(event.Name) == ".yaml" {

						configFiles = append(configFiles, event.Name)

						log.Println("Now watching", event.Name)
					}

				} else if event.Op&fsnotify.Rename == fsnotify.Rename {

					for i, name := range configFiles {
						if name == event.Name {

							configFiles = append(
								configFiles[:i],
								configFiles[i+1:]...)

							configFiles = append(configFiles, event.Name)

							log.Println("Renamed", event.Name)
							break
						}
					}

				} else if event.Op&fsnotify.Remove == fsnotify.Remove {

					for i, name := range configFiles {
						if name == event.Name {

							configFiles = append(
								configFiles[:i],
								configFiles[i+1:]...)

							log.Println("Removed", event.Name)
							break
						}
					}

				} else if event.Op&fsnotify.Write == fsnotify.Write {

					log.Println("Modified", event.Name)

					for _, name := range configFiles {
						if name == event.Name {

							err := Rebuild(event.Name)
							if err != nil {
								log.Println("ERROR:", err)
							}

							log.Println("  Building", event.Name)
							break
						}
					}

				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	path, err := filepath.Abs(pageDir)
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
