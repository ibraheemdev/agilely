package test

import (
	"os"
	"path"
	"runtime"
)

// MoveToRoot : execute tests root directory
func MoveToRoot() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
