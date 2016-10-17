package main

import (
    "stathat.com/c/jconfig"
    "os"
    "log"
    "fmt"
)

func main() {
    config := jconfig.LoadConfig("./settings.json")
    datafolder := config.GetString("data_folder")

    if err := CheckDataFolder(datafolder); err == nil {
        log.Printf("Folder %s is valid", datafolder)
    } else {
        log.Fatalln("Folder is not valid:", err)
    }   
}

func CheckDataFolder(folderpath string) (err error) {
    // ##TODO## also check persmissions
    finfo, err := os.Stat(folderpath)
    if err == nil {
        if !finfo.IsDir() {
            err = fmt.Errorf("Specified path is not a folder")
        }
    }
    return
}