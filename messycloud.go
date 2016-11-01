package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"stathat.com/c/jconfig"
)

func main() {
	var newWatchPath = make(chan string)

	config := jconfig.LoadConfig("./settings.json")
	datafolder := config.GetString("datafolder")
	if datafolder == "" {
		log.Fatalln("datafolder not specifed or empty in settings.json")
	}

	// Set up messyWatcher object and handle events
	setupWatch(newWatchPath)

	// Watch the datafolder, and recurse as required
	if err := indexFolder(newWatchPath, datafolder); err != nil {
		log.Fatalln(err)
	}
}

func indexFolder(newWatchPath chan string, datafolder string) error {
	// Returns list of items in datafolder
	items, err := ioutil.ReadDir(datafolder)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.IsDir() {
			// addToDB(datafolder)
			newWatchPath <- filepath.Join(datafolder, item.Name())
		} else if item.Name()[0:1] == "." {
			log.Printf("Hidden file \"%s\" ignored", item.Name())
		} else {
			// addToDB(datafolder & item.Name())
			log.Println(item.Name(), "is a file")
		}
	}

	// Watch for changes
	newWatchPath <- datafolder
	return nil
}

func setupWatch(newWatchPath chan string) {
	done := make(chan bool)
	messyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer messyWatcher.Close()

	go func() {
		for {
			select {
			case event := <-messyWatcher.Events:
				log.Println("Event:", event)
			case err := <-messyWatcher.Errors:
				log.Printf("error: %s\n\n", err)
			case newPath := <-newWatchPath:
				if err := messyWatcher.Add(newPath); err != nil {
					log.Println("Error adding watch for", newPath)
				}
				log.Println("New watch added for", newPath)
			}
		}
	}()
	<-done
}
