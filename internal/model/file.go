package model

import (
	"database/sql"
	"fmt"
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

type FileModel struct {
	DB *sql.DB
}

// GetFiles get files from absolute path
/*
dir: 		  .|G:/test
absolutePath: .|G:/test/a/t.txt
*/
func GetFiles(dir, absolutePath string) []*File {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		log.Println(err)
		return nil
	}

	relativePath := strings.Replace(absolutePath, dir, "", -1)
	if relativePath == "/" {
		relativePath = ""
	}

	var myFiles []*File
	for _, file := range files {
		f := &File{
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

func (m *FileModel) DeleteFileIndex() error {
	q := `DELETE FROM files;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	return err
}

func (m *FileModel) Search(q string) ([]*File, error) {
	s := `SELECT name,isdir,size,last_modified,path 
	FROM files WHERE name LIKE '%` + q + `%'`
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*File
	layout := "2006-01-02 15:04:05"

	for rows.Next() {
		file := &File{}
		var lastModified string
		if err = rows.Scan(&file.Name, &file.IsDir, &file.Size,
			&lastModified, &file.Path); err != nil {
			fmt.Println("rows scan error:", err)
			return nil, err
		}
		// parse datatime -> 2021-01-04 08:25:32.629566+08:00
		last := strings.Split(lastModified, "+")[0]
		lastModifiedDate, err := time.Parse(layout, last)
		if err != nil {
			fmt.Println(err)
			return files, err
		}
		file.LastModified = lastModifiedDate
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
