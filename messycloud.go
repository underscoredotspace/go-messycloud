package main

import (
    "stathat.com/c/jconfig"
    "os"
    "log"
    "fmt"
)

func main() {
    config := jconfig.LoadConfig("./settings.json")
    datafolder := config.GetString("datafolder")
    if datafolder == "" {
      log.Fatalln("datafolder not specifed in settings.json")
    }

    if err := ValidDataFolder(datafolder); err == nil {
      log.Printf("Folder %s is valid\n", datafolder)
    } else {
      log.Fatalln(err)
    }
}

func ValidDataFolder(folderpath string) (error) {    
    // Does folderpath exist?
    if finfo, err := os.Stat(folderpath); err == nil {
      // Is folderpath a folder?
      if finfo.IsDir() {
        // ##TODO## - Do I have required permissions to folderpath?
        return nil
      } else {
        return fmt.Errorf("%s: not a folder", folderpath)
      }
    } else {
      return err
    }
}
