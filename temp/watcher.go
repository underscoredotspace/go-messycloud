package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func messyWatcher(path string) error {

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
				log.Println("Event:", uint(event.Op))
				switch event.Op {
				case fsnotify.Write:
					log.Println("Edit:", event.Name)
				case fsnotify.Create:
					log.Println("Create:", event.Name)
					if err := isFolder(event.Name); err == nil {
						if err = addWatch(*watcher, event.Name); err != nil {
							log.Println("Error adding watch for folder: ", event.Name)
						}
					}
				case fsnotify.Remove, fsnotify.Op(12):
					log.Println("Delete:", event.Name)
					if err := isFolder(event.Name); err == nil {
						if err = watcher.Remove(event.Name); err != nil {
							log.Println("Error removing watch for folder: ", event.Name)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Printf("error: %s\n\n", err)
			}
		}
	}()

	err = addWatch(*watcher, "/var/cloud")
	if err != nil {
		log.Fatal(err)
	}
	<-done

	if err != nil {
		return err
	}
	return nil
}

func addWatch(watcher fsnotify.Watcher, path string) error {
	if err := watcher.Add(path); err != nil {
		return err
	}
	log.Println("New watch added:", path)
	return nil
}

func isFolder(path string) error {
	// Does folderpath exist?
	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Is folderpath a folder?
	if !finfo.IsDir() {
		return fmt.Errorf("%s: not a folder", path)
	}
	return nil
}
