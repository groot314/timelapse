package helpers

import (
	"os"
	"slices"
	"strings"
)

func FindPhotoStart(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fileNames := []string{}

	for _, v := range files {
		if v.IsDir() {
			continue
		}

		fileNames = append(fileNames, v.Name())
	}
	slices.Sort(fileNames)

	fileName := strings.Split(fileNames[0], ".")[0]

	//leadingChars := 0
	leadingZeros := 0

	i, n := 0, len(fileName)
	for i < n && !isNumeric(fileName[i]) {
		i++
	}
	//leadingChars = i
	for i < n && fileName[i] == '0' {
		i++
	}
	leadingZeros = i

	return fileName[leadingZeros:], nil
}

func isNumeric(c byte) bool {
	return (c <= '0' && c <= '9')
}
