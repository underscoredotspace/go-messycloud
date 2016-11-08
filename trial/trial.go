package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"stathat.com/c/jconfig"
)

type fsItem int

const (
	messyFile fsItem = iota
	messyFolder
)

type messyItem struct {
	mtime  time.Time
	fsItem fsItem
}

type messyFolders map[string]messyItem

func main() {
	config := jconfig.LoadConfig("../settings.json")
	datafolder := config.GetString("datafolder")
	if datafolder == "" {
		log.Fatalln("datafolder not specifed or empty in settings.json")
	}

	myfolders := make(messyFolders)
	myfolders[datafolder] = messyItem{time.Now(), messyFolder}
	indexFolder(datafolder, myfolders)
	for path, item := range myfolders {
		fmt.Println(path, "is a", item.fsItem, "last edited", item.mtime)
	}
}

func indexFolder(indexPath string, messyfolder messyFolders) {
	var fullPath string
	items, err := ioutil.ReadDir(indexPath)
	if err != nil {
		log.Println(err)
	}
	for _, item := range items {
		if item.Name()[0:1] != "." {
			fullPath = filepath.Join(indexPath, item.Name())
			if item.IsDir() {
				messyfolder[fullPath] = messyItem{time.Now(), messyFolder}
				indexFolder(fullPath, messyfolder)
			} else {
				messyfolder[fullPath] = messyItem{time.Now(), messyFile}
			}
		}
	}
}
