package resourceenvironment

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// This file contains functions used to read/write pipeline environment data from/to disk.
// The content of the file is the value. For the custom parameters this could for example also be a JSON representation of a more complex value.
// The file representation looks as follows:

// <pipeline env path>/

// <pipeline env path>/artifactVersion

// <pipeline env path>/git/
// <pipeline env path>/git/branch
// <pipeline env path>/git/commitId
// <pipeline env path>/git/commitMessage
// <pipeline env path>/git/repositoryUrl -> TODO: storing function(s) with ssh and https getters

// <pipeline env path>/github/
// <pipeline env path>/github/owner
// <pipeline env path>/github/repository

// <pipeline env path>/custom/
// <pipeline env path>/custom/<parameter>

// SetParameter sets any parameter in the pipeline environment or another environment stored in the file system
func SetParameter(path, name, value string) error {
	paramPath := filepath.Join(path, name)
	return writeToDisk(paramPath, []byte(value))
}

// GetParameter reads any parameter from the pipeline environment or another environment stored in the file system
func GetParameter(path, name string) string {
	paramPath := filepath.Join(path, name)
	return readFromDisk(paramPath)
}

func writeToDisk(filename string, data []byte) error {

	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(filename), 0700)
	}

	//ToDo: make sure to not overwrite file but rather add another file? Create error if already existing?
	return ioutil.WriteFile(filename, data, 0700)
}

func readFromDisk(filename string) string {
	//ToDo: if multiple files exist, read from latest file
	v, err := ioutil.ReadFile(filename)
	val := string(v)
	if err != nil {
		val = ""
	}
	return val
}
