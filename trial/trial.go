package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"stathat.com/c/jconfig"
)

type messyfolders struct {
	Path    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func main() {
	config := jconfig.LoadConfig("../settings.json")
	datafolder := config.GetString("datafolder")
	if datafolder == "" {
		log.Fatalln("datafolder not specifed or empty in settings.json")
	}

	fileList, err := getMessyStructure(datafolder)
	if err != nil {
		log.Fatalln(err)
	}

	output, err := json.Marshal(fileList)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

func getMessyStructure(datafolder string) ([]messyfolders, error) {
	var fileList []messyfolders
	err := filepath.Walk(datafolder, func(path string, finfo os.FileInfo, err error) error {
		if finfo.Name()[0:1] != "." {
			fileList = append(fileList, messyfolders{path, finfo.Size(), finfo.Mode(), finfo.ModTime(), finfo.IsDir()})
		}
		return nil
	})
	return fileList, err
}
