package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type fileSetup []fileSpec

type fileSpec struct {
	fileExt      string
	fileContents []byte
}

func createTestFiles(files []fileSpec, t *testing.T) ([]string, []file) {
	t.Helper()
	var createdTempFilesLocations []string
	var tempFiles []file
	for i, f := range files {
		tempFile, err := ioutil.TempFile(".", fmt.Sprintf("test*.%v", f.fileExt))
		if err != nil {
			t.Errorf("Unable to create temp file")
		}
		_, err = tempFile.Write(f.fileContents)
		if err != nil {
			t.Errorf("Unable to write file contexts to %v", tempFile.Name())
		}
		createdTempFilesLocations = append(createdTempFilesLocations, tempFile.Name())
		tempFiles = append(tempFiles, file{
			Path:  tempFile.Name(),
			Order: i,
			Data:  f.fileContents,
		})
		tempFile.Close()
	}
	return createdTempFilesLocations, tempFiles
}

func deleteTestFiles(fileLocations []string) {
	for _, file := range fileLocations {
		os.Remove(file)
	}
}
