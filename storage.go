package radix

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	basePath   string = "/storage"
	currentDir string
)

func init() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		Error(err.Error())
		return
	}
	_, err = os.Stat(filepath.Join(currentDir, basePath))
	if os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(currentDir, basePath), 0755)
		if err != nil {
			Error(err.Error())
			return
		}
		Info("Created a folder: %s", basePath)
	}
}

func handleFilename(filename *string) {
	files := strings.Split(*filename, "/")
	*filename = files[len(files)-1]
}

func file(filename string) string {
	return filepath.Join(currentDir, basePath, filename)
}

func SaveJSON(filename string, data Data) error {
	handleFilename(&filename)
	bytes, err := json.MarshalIndent(&data, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file(filename), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadJSON(filename string, data Data) error {
	handleFilename(&filename)
	bytes, err := ioutil.ReadFile(file(filename))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}
	return nil
}

func RemoveJSON(filename string) error {
	handleFilename(&filename)
	err := os.Remove(file(filename))
	if err != nil {
		return err
	}
	return nil
}
