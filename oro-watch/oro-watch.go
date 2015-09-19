package main

import (
	"github.com/njwilson23/orogenesis"
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Rebuild(configpath string) error {
	return nil
}

/* Add the source files for a page to a watcher */
func AddPageSources(w *fsnotify.Watcher, configDir string, config *orogenesis.Page) error {
	var err error
	err = nil
	if len(config.TemplatePath) != 0 {
		err := w.Add(filepath.Join(configDir, config.TemplatePath))
		log.Println(" ", filepath.Join(configDir, config.TemplatePath))
		if err != nil {
			log.Println(err)
			log.Println(config.TemplatePath)
		}
	}
	if len(config.HeaderPath) != 0 {
		err := w.Add(config.HeaderPath)
		log.Println("DBG: added", config.HeaderPath)
		if err != nil {
			log.Println(err)
			log.Println(" ", config.HeaderPath)
		}
	}
	if len(config.NavPath) != 0 {
		err := w.Add(config.NavPath)
		log.Println("DBG: added", config.NavPath)
		if err != nil {
			log.Println(err)
			log.Println(" ", config.NavPath)
		}
	}
	if len(config.BodyPath) != 0 {
		err := w.Add(config.BodyPath)
		log.Println("DBG: added", config.BodyPath)
		if err != nil {
			log.Println(err)
			log.Println(" ", config.BodyPath)
		}
	}
	if len(config.FooterPath) != 0 {
		err := w.Add(config.FooterPath)
		log.Println("DBG: added", config.FooterPath)
		if err != nil {
			log.Println(err)
			log.Println(" ", config.FooterPath)
		}
	}
	return err
}

// Check the source files for *configpath* and update a hash table with dependency relationships
func UpdateSourceDependencies(sourceHash map[string]string, configpath string, config *orogenesis.Page) error {
	sourceHash[configpath] = configpath
	if len(config.TemplatePath) != 0 {
		sourceHash[config.TemplatePath] = configpath
	}
	if len(config.HeaderPath) != 0 {
		sourceHash[config.HeaderPath] = configpath
	}
	if len(config.NavPath) != 0 {
		sourceHash[config.NavPath] = configpath
	}
	if len(config.BodyPath) != 0 {
		sourceHash[config.BodyPath] = configpath
	}
	if len(config.FooterPath) != 0 {
		sourceHash[config.FooterPath] = configpath
	}
	return nil
}

func main() {

	var configDir string
	configDir = os.Args[1]
	if len(configDir) == 0 {
		panic("Usage: oro-watch *directory*")
	}

	// Watcher to watch sources and trigger page builds
	source_watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer source_watcher.Close()

	// Check for existing yaml configuration files, build a source hash map,
	// and set a watcher on the sources
	var files []os.FileInfo
	files, err = ioutil.ReadDir(configDir)
	if err != nil {
		log.Println("ERROR:", err)
		panic(err)
	}

	var fnm string
	var config *orogenesis.Page
	sourceHash := make(map[string]string)

	for _, file := range files {
		fnm = file.Name()

		if filepath.Ext(fnm) == ".yaml" {

			fnm = filepath.Join(configDir, fnm)
			log.Println("watching", fnm)

			config, err = orogenesis.ReadConfig(fnm)
			if err != nil {
				log.Fatal(err)
			}

			err = AddPageSources(source_watcher, configDir, config)
			if err != nil {
				log.Fatal(err)
			}

			err = UpdateSourceDependencies(sourceHash, fnm, config)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Watcher to check for configuration files
	config_watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer config_watcher.Close()

	done := make(chan bool) // Channel will never be written to, so program will run as a daemon

	// Hook up the source file watcher
	go func() {
		var configpath string
		for {
			select {
			case event := <-source_watcher.Events:

				if event.Op&fsnotify.Remove == fsnotify.Remove {

					// Saving a file can generale delete events. Check whether
					// file actually deleted before acting
					if _, err = os.Stat(event.Name); err == nil {
						source_watcher.Add(event.Name)
					} else {
						log.Println("DBG: Remove", event.Name)
						delete(sourceHash, event.Name)
						err := source_watcher.Remove(event.Name)
						if err != nil {
							log.Println(err)
						}
					}
				} else {

					// Write, Rename, Chmod
					configpath = sourceHash[event.Name]
					log.Println("DBG:", event.Name)
					if len(configpath) == 0 {
						log.Println("Error: key for", event.Name, "not found")
					} else {

						config, err = orogenesis.ReadConfig(configpath)

						if err != nil {
							log.Fatal(err)

						} else {
							log.Println("building", configpath)
							_, err = orogenesis.BuildPage(configpath, config)
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}

			case err := <-source_watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// Hook up the config directory watcher
	go func() {
		for {
			select {
			case event := <-config_watcher.Events:

				if event.Op&fsnotify.Create == fsnotify.Create {

					if filepath.Ext(event.Name) == ".yaml" {

						config, err = orogenesis.ReadConfig(fnm)
						if err != nil {
							log.Fatal(err)
						}
						AddPageSources(source_watcher, configDir, config)
						log.Println("Now watching", event.Name)
					}

				} else if event.Op&fsnotify.Rename == fsnotify.Rename {

					config, err = orogenesis.ReadConfig(fnm)
					if err != nil {
						log.Fatal(err)
					}
					AddPageSources(source_watcher, configDir, config)
					log.Println("Renamed", event.Name)

				} else if event.Op&fsnotify.Remove == fsnotify.Remove {

					log.Println("Removed", event.Name)

				} else if event.Op&fsnotify.Write == fsnotify.Write {

					//pass

				}

			case err := <-config_watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	configDir, err = filepath.Abs(configDir)
	if err != nil {
		log.Fatal(err)
	}
	err = config_watcher.Add(configDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
