package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"stathat.com/c/jconfig"
)

func main() {
	config := jconfig.LoadConfig("./settings.json")
	datafolder := config.GetString("datafolder")
	if datafolder == "" {
		log.Fatalln("datafolder not specifed or empty in settings.json")
	}

	watchFolderList, err := indexDataFolder(datafolder)
	if err != nil {
		log.Fatalln(err)
	}
	watchFolders(watchFolderList)
}

func indexDataFolder(folderpath string) ([]string, error) {
	watchFolderList := []string{folderpath}

	// Returns list of items in datafolder
	items, err := ioutil.ReadDir(folderpath)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Name()[0:1] == "." {
			//log.Printf("Hidden item \"%s\" ignored", item.Name())
		} else {
			if item.IsDir() {
				// addToDB(datafolder)
				recurseItem, err := indexDataFolder(filepath.Join(folderpath, item.Name()))
				if err != nil {
					log.Fatalln(err)
				}
				for _, subItem := range recurseItem {
					watchFolderList = append(watchFolderList, subItem)
				}
			} else {
				// addToDB(datafolder & item.Name())
				//log.Println(item.Name(), "is a file")
			}
		}
	}

	return watchFolderList, err
}

func watchFolders(watchFolderList []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, item := range watchFolderList {
		if err := watcher.Add(item); err != nil {
			log.Println("Failed to add new watch: ", item, err)
		}
		log.Println("New watch added:", item)
	}

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
	<-done
}
