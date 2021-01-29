package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// ConvertByteToMB ...
func ConvertByteTo(n int64) string {
	switch {
	case n == 0:
		return ""
	case n < 1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024) + "KB"
	case n < 1024*1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024) + "MB"
	default:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024/1024) + "GB"
	}
}

// https://stackoverflow.com/questions/49057032/recursively-zipping-a-directory-in-golang
func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))

		// Chinese garbled code
		path, err := utf8ToGBK(relPath)
		if err != nil {
			return err
		}

		zipFile, err := myZip.Create(path)
		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

func utf8ToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewEncoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}
