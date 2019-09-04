package test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Item is a dataset item to be use for testing/benchmarking of ImageProcessor implementations
type Item struct {
	InputFile    []byte
	InputWidth   uint
	InputHeight  uint
	OutputFile   []byte // if nil then don't test by "OutputFile"
	OutputWidth  uint
	OutputHeight uint
}

func panicIfError(err error) {
	if err == nil {
		return
	}

	panic(err)
}

// Dataset returns an in-memory dataset of pictures to be used for testing/benchmarking
// of ImageProcessor implementations
func Dataset() (result []Item) {
	datasetPath := `../common/test/dataset`

	datasetDir, err := os.Open(datasetPath)
	panicIfError(err)

	fileNames, err := datasetDir.Readdirnames(0)
	panicIfError(err)

	for _, fileName := range fileNames {
		words := strings.Split(strings.Split(fileName, `.`)[0], `_`)
		size := words[1]
		sizeParts := strings.Split(size, `x`)

		width64, err := strconv.ParseUint(sizeParts[0], 10, 64)
		panicIfError(err)

		height64, err := strconv.ParseUint(sizeParts[1], 10, 64)
		panicIfError(err)

		bytes, err := ioutil.ReadFile(filepath.Join(datasetPath, fileName))
		panicIfError(err)

		result = append(result, Item{
			InputFile:    bytes,
			InputWidth:   uint(width64),
			InputHeight:  uint(height64),
			OutputWidth:  uint(width64) * 2 / 3,
			OutputHeight: uint(height64) * 2 / 3,
		})

		result = append(result, Item{
			InputFile:    bytes,
			InputWidth:   uint(width64),
			InputHeight:  uint(height64),
			OutputWidth:  uint(width64) * 3 / 2,
			OutputHeight: uint(height64) * 3 / 2,
		})
	}

	return
}
