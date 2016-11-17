package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatalln(err)
	}
	defer session.Close()

	if err := updateMessyDatabase(fileList, session); err != nil {
		log.Fatalln(err)
	}

	dbFileList, err := getfromMessyDatabase(session)
	if err != nil {
		log.Fatalln(err)
	}

	var dbFileItem messyfolders
	for _, dbFileItem = range dbFileList {
		fmt.Println(dbFileItem.Path, dbFileItem.ModTime)
	}

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

func updateMessyDatabase(fileList []messyfolders, session *mgo.Session) error {
	collection := session.DB("colin").C("messyfiles")
	for _, fileItem := range fileList {
		if err := collection.Insert(fileItem); err != nil {
			return err
		}
	}

	return nil
}

func getfromMessyDatabase(session *mgo.Session) ([]messyfolders, error) {
	var fileList []messyfolders

	collection := session.DB("colin").C("messyfiles")
	// err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	if err := collection.Find(bson.M{}).All(&fileList); err != nil {
		log.Fatal(err)
	}

	return fileList, nil
}
