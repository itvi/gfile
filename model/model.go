package model

import (
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type File struct {
	Name         string
	IsDir        bool
	Size         int64
	LastModified time.Time
	Path         string
}

func GetFiles(dir, absolutePath string) []File {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	var myFiles []File
	for _, file := range files {
		relativePath := strings.Replace(absolutePath, dir, "", -1)
		if dir == absolutePath {
			relativePath = "."
		}
		f := File{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			Size:         file.Size(),
			LastModified: file.ModTime(),
			Path:         relativePath,
		}
		myFiles = append(myFiles, f)
	}
	return myFiles
}
