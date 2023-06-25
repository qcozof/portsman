package utils

import (
	"os"
)

type DirectoryUtils struct {
}

func (_ *DirectoryUtils) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (svc *DirectoryUtils) CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := svc.PathExists(v)
		if err != nil {
			return err
		}
		if !exist {

			if err = os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}
