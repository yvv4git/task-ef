package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrGetRurrentFile = errors.New("can't get current file path")
	ErrChangeDir      = errors.New("can't change directory")
)

func ChangeDirToProjectRoot(reversePath string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ErrGetRurrentFile
	}
	projectRoot := filepath.Join(filepath.Dir(filename), reversePath)

	if err := os.Chdir(projectRoot); err != nil {
		return ErrChangeDir
	}

	return nil
}
