package main

import (
	"github.com/njwilson23/orogenesis"
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Rebuild is not implemented. When it is, it will force a page rebuild.
func Rebuild(configpath string) error {
	return nil
}

/* AddPageSources add the source files for a page to a watcher. */
func AddPageSources(w *fsnotify.Watcher, configPath string, config *orogenesis.Page) error {

	//configDir := filepath.Dir(configPath)
	err := w.Add(configPath)
	if err != nil {
		log.Println("ERROR:", err)
		log.Println(configPath)
	} else {
		log.Println(" added source:", configPath)
	}

	if len(config.TemplatePath) != 0 {
		err := w.Add(config.TemplatePath)
		if err != nil {
			log.Println("ERROR:", err)
			log.Println(config.TemplatePath)
		} else {
			log.Println(" added source:", config.TemplatePath)
		}
	}
	if len(config.HeaderPath) != 0 {
		err := w.Add(config.HeaderPath)
		if err != nil {
			log.Println("ERROR:", err)
			log.Println(" ", config.HeaderPath)
		} else {
			log.Println(" added source:", config.HeaderPath)
		}
	}
	if len(config.NavPath) != 0 {
		err := w.Add(config.NavPath)
		if err != nil {
			log.Println("ERROR:", err)
			log.Println(" ", config.NavPath)
		} else {
			log.Println(" added source:", config.NavPath)
		}
	}
	if len(config.BodyPath) != 0 {
		err := w.Add(config.BodyPath)
		if err != nil {
			log.Println("ERROR:", err)
			log.Println(" ", config.BodyPath)
		} else {
			log.Println(" added source:", config.BodyPath)
		}

	}
	if len(config.FooterPath) != 0 {
		err := w.Add(config.FooterPath)
		if err != nil {
			log.Println("ERROR:", err)
			log.Println(" ", config.FooterPath)
		} else {
			log.Println(" added source:", config.FooterPath)
		}
	}
	return err
}

// UpdateDependencyHash checks the source files for *configpath* and updates a
// hash table with dependency relationships
func UpdateDependencyHash(sourceHash map[string]string, configpath string, config *orogenesis.Page) error {

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
	sourceWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer sourceWatcher.Close()

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

			err = AddPageSources(sourceWatcher, fnm, config)
			if err != nil {
				log.Fatal(err)
			}

			err = UpdateDependencyHash(sourceHash, fnm, config)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Watcher to check for configuration files
	configWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer configWatcher.Close()

	/* Channel will never be written to, so program will run as a daemon */
	done := make(chan bool)

	/* Hook up the source file watcher

	This watches a list of source files, each of which is also a key in
	*sourceHash*. If one is written to, the watcher triggers BuildPage.

	If a source other than a *.yaml disappears, the watcher will attempt to
	re-add it.

	If a *.yaml disappears, every source file that it depended on is removed. */
	go func() {
		var configPath string
		var err error
		for {
			select {
			case event := <-sourceWatcher.Events:

				if (event.Op&fsnotify.Remove == fsnotify.Remove) ||
					(event.Op&fsnotify.Rename == fsnotify.Rename) {

					// Saving a file can generate delete events. Check whether
					// file actually deleted before acting
					if _, err = os.Stat(event.Name); err == nil {
						sourceWatcher.Add(event.Name)
					} else {

						log.Println("DBG: Remove", event.Name)

						if filepath.Ext(event.Name) == ".yaml" {
							for k, v := range sourceHash {
								if v == event.Name {
									delete(sourceHash, k)
									err := sourceWatcher.Remove(k)
									if err != nil {
										log.Println("ERROR:", err)
									}
									log.Println("DBG: scrubbed", k)
								}
							}
						} else {

							delete(sourceHash, event.Name)
							err := sourceWatcher.Remove(event.Name)
							if err != nil {
								log.Println("ERROR:", err)
							}
							log.Println("DBG: scrubbed", event.Name)
						}

					}
				} else {

					// Write, Chmod
					configPath = sourceHash[event.Name]
					if len(configPath) == 0 {

						log.Println("ERROR: key for", event.Name, "not found")

					} else {

						config, err = orogenesis.ReadConfig(configPath)
						if err != nil {

							log.Fatal(err)

						} else {

							log.Println("  building", configPath)
							_, err = orogenesis.BuildPage(configPath, config)
							if err != nil {
								log.Fatal("ERROR:", err)
							}

						}
					}
				}

			case err := <-sourceWatcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	/* Hook up the config directory watcher

	This is responsible for watching the directory specified as a
	command-line argument for new YAML files. If one appears, it and its
	sources are added to the source watcher. If one disappears, it and its
	sources are removed from the source watcher. */
	go func() {
		var err error
		for {
			select {
			case event := <-configWatcher.Events:

				if event.Op&fsnotify.Create == fsnotify.Create {

					if filepath.Ext(event.Name) == ".yaml" {

						config, err = orogenesis.ReadConfig(event.Name)
						if err != nil {
							log.Fatal(err)
						}
						err = AddPageSources(sourceWatcher, event.Name, config)
						if err != nil {
							log.Println("ERROR:", err)
						}
						err = UpdateDependencyHash(sourceHash, event.Name, config)
						if err != nil {
							log.Fatal(err)
						}
						log.Println("Now watching", event.Name)
					}

				} else if event.Op&fsnotify.Rename == fsnotify.Rename {

					//if filepath.Ext(event.Name) == ".yaml" {
					//	config, err = orogenesis.ReadConfig(event.Name)
					//	if err != nil {
					//		log.Fatal(err)
					//	}
					//	AddPageSources(sourceWatcher, event.Name, config)
					//	log.Println("Renamed", event.Name)
					//}
					log.Println("Renamed", event.Name, "(no-op)")

				} else if event.Op&fsnotify.Remove == fsnotify.Remove {

					// TODO: remove from source watcher and source hash
					log.Println("Removed", event.Name, "(no-op)")

				} else if event.Op&fsnotify.Write == fsnotify.Write {

					// pass - source watcher responsible for rebuilds

				}

			case err := <-configWatcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	configDir, err = filepath.Abs(configDir)
	if err != nil {
		log.Fatal(err)
	}
	err = configWatcher.Add(configDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
