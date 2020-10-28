package main

import (
	"path/filepath"
)

type DiscoveryClient interface {
	getFiles(path string) []string
}

type BasicDiscoveryClient struct {
}

func (instance *BasicDiscoveryClient) getFiles(path string) []string {

	if len(path) == 0 {
		return make([]string, 0)
	}

	file := openFile(path)

	fi, err := file.Stat()

	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		return []string{path}
	}

	var children []string

	err = filepath.Walk(path, addChildren(&children))

	if err != nil {
		panic(err)
	}

	return children
}
