package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func addChildren(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			panic(err)
		}

		if !info.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}

func openFile(f string) *os.File {
	r, err := os.Open(f)

	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}

	return r
}
