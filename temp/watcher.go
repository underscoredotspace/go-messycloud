package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
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
				log.Println("Event:", event)
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
