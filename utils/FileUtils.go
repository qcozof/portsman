package utils

import (
	"os"
	"strings"
	"sync"
)

type FileUtils struct {
}

func (_ *FileUtils) Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic("Exists:" + err.Error())
	}
	return true
}

func (_ *FileUtils) VerifyExt(contentType string, extArr []string) bool {
	for _, val := range extArr {
		if contentType == strings.Trim(val, "") {
			return true
		}
	}
	return false
}

var mutex sync.Mutex

func (_ *FileUtils) ReadFile(path string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	bt, err := os.ReadFile(path)
	return string(bt), err
}

func (_ *FileUtils) Write(path string, res []byte, perm os.FileMode) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return os.WriteFile(path, res, perm)
}

func (_ *FileUtils) Append(path string, res string, perm os.FileMode) error {

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(res); err != nil {
		return err
	}
	return nil
}

func (_ *FileUtils) Delete(file string) error {
	return os.Remove(file)
}
