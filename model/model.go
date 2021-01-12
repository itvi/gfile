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

// GetFiles get files from absolute path
/*
dir: 		  .|G:/test
absolutePath: .|G:/test/a/t.txt
*/
func GetFiles(dir, absolutePath string) []File {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		log.Println(err)
		return nil
	}

	relativePath := strings.Replace(absolutePath, dir, "", -1)
	if relativePath == "/" {
		relativePath = ""
	}

	var myFiles []File
	for _, file := range files {
		f := File{
			Name:         file.Name(),
			IsDir:        file.IsDir(),
			Size:         file.Size(),
			LastModified: file.ModTime(),
			Path:         relativePath + "/" + file.Name(),
		}
		myFiles = append(myFiles, f)
	}
	return myFiles
}
