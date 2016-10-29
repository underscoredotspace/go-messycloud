package main

import (
  "fmt"
  "os"
)

func watchFiles(foldername string) (error) {
  if err := validDataFolder(foldername); err !=nil {
    return err
  }

  AddWatch(foldername)

/*
  for each item in foldername {
    addToDB(item)
    if item == folder {
      WatchFiles(item.path)
    }
  }    
*/  
}

func AddWatch(path) {
  
}

func ValidDataFolder(path string) (error) {    
  // Does path exist?
  if finfo, err := os.Stat(path); err == nil {
    // Is path a folder?
    if finfo.IsDir() {
      return nil
    } else {
      return fmt.Errorf("%s: not a folder", path)
    }
  } else {
    return err
  }
}