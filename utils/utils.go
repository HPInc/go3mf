package utils

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func GetRootProjectDir() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	return dir, err
}

func WithRootProjectDir(callback func(string)) {
	dir, err := GetRootProjectDir()
	if err != nil {
		panic(err)
	}
	callback(dir)
}

func WithProjectDirAndTestTempDirRemoveAtEnd(fileFolder string, callback func(testFileDir string, testTempDir string)) {
	dir, err := GetRootProjectDir()
	if err != nil {
		panic(err)
	}
	testTempDir, err := ioutil.TempDir("", "testcase")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(testTempDir)
	defer os.Remove(testTempDir)

	callback(filepath.Join(dir, fileFolder), testTempDir)
}
