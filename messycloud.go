package main

import (
	"io/ioutil"
	"log"

	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"stathat.com/c/jconfig"
)

var messyWatcher fsnotify.Watcher

func main() {
	config := jconfig.LoadConfig("./settings.json")
	datafolder := config.GetString("datafolder")
	if datafolder == "" {
		log.Fatalln("datafolder not specifed or empty in settings.json")
	}

	// Set up messyWatcher object and handle events
	err := setupWatch()
	if err != nil {
		log.Fatalln(err)
	}

	// Watch the datafolder, and recurse as required
	if err := startWatch(datafolder); err != nil {
		log.Fatalln(err)
	}
}

func startWatch(datafolder string) error {
	// Returns list of items in datafolder
	items, err := ioutil.ReadDir(datafolder)
	if err != nil {
		return err
	}

	// Watch for changes
	addWatch(datafolder)

	for _, item := range items {
		if item.IsDir() {
			// addToDB(datafolder)
			if err := startWatch(filepath.Join(datafolder, item.Name())); err != nil {
				log.Println(err)
			}
		} else if item.Name()[0:1] == "." {
			log.Printf("Hidden file \"%s\" ignored", item.Name())
		} else {
			// addToDB(datafolder & item.Name())
			log.Println(item.Name(), "is a file")
		}
	}
	return nil
}

func addWatch(path string) error {
	if err := messyWatcher.Add(path); err != nil {
		return err
	}
	log.Println("New watch added:", path)
	return nil
}

func setupWatch() error {
	messyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer messyWatcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-messyWatcher.Events:
				log.Println("Event:", uint(event.Op))
				// switch event.Op {
				// case fsnotify.Write:
				// 	log.Println("Edit:", event.Name)
				// case fsnotify.Create:
				// 	log.Println("Create:", event.Name)
				// 	if err := isFolder(event.Name); err == nil {
				// 		if err = addWatch(event.Name); err != nil {
				// 			log.Println("Error adding watch for folder: ", event.Name)
				// 		}
				// 	}
				// case fsnotify.Remove, fsnotify.Op(12):
				// 	log.Println("Delete:", event.Name)
				// 	if err := isFolder(event.Name); err == nil {
				// 		if err = messyWatcher.Remove(event.Name); err != nil {
				// 			log.Println("Error removing watch for folder: ", event.Name)
				// 		}
				// 	}
				// }
			case err := <-messyWatcher.Errors:
				log.Printf("error: %s\n\n", err)
			}
		}
	}()
	<-done
	return nil
}
