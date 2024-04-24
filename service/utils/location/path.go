package location

import (
	"log"
	"fmt"
	"os"
	
	"path/filepath"
	"runtime"
	
)

//GetCurrentAbPath 
func GetCurrentAbPath() (dir string, err error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Failed to get directory")
	}
	dir = filepath.Dir(filename)
	rootDir := filepath.Join(dir, "..", "..")
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(err)
		return "", fmt.Errorf("Failed to get directory")
	}
	return absRootDir, nil
}

// IsDir 
func IsDir(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// CreateDir 
func CreateDir(dirName string) bool {
	err := os.Mkdir(dirName, 755)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
