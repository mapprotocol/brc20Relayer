package startstore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func StoreLatestStart(file string, number uint64) error {
	// Create dir if it does not exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		dir, _ := filepath.Split(file)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Write bytes to file
	data := []byte(strconv.FormatUint(number, 10))
	err := ioutil.WriteFile(file, data, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ReadLatestStart(file string) (uint64, error) {
	// If it exists, load and return
	exists, err := fileExists(file)
	if err != nil {
		return 0, err
	}
	if exists {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return 0, err
		}
		block, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 0)
		if err != nil {
			return 0, err
		}
		return block, nil
	}
	// Otherwise just return 0
	return 0, nil
}
