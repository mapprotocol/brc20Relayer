package startstore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const sep = ":"

func fileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func StoreLatestStart(file string, start, count uint64) error {
	// Create dir if it does not exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		dir, _ := filepath.Split(file)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Write bytes to file
	data := []byte(fmt.Sprintf("%d%s%d", start, sep, count))
	err := ioutil.WriteFile(file, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ReadLatestStart(file string) (start, count uint64, err error) {
	// If it exists, load and return
	exists, err := fileExists(file)
	if err != nil {
		return 0, 0, err
	}
	if exists {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return 0, 0, err
		}
		ns := strings.Split(strings.TrimSpace(string(data)), sep)
		start, err = strconv.ParseUint(strings.TrimSpace(ns[0]), 10, 0)
		if err != nil {
			return 0, 0, err
		}
		count, err = strconv.ParseUint(strings.TrimSpace(ns[0]), 10, 0)
		if err != nil {
			return 0, 0, err
		}
		return 0, 0, nil
	}
	// Otherwise just return 0
	return start, count, nil
}
