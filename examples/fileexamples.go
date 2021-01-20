package examples

import (
	"fmt"
	"os"
	"path/filepath"
)

func TestWalkFilesExample() {

	var files []string
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Println("...", path)
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}

}
