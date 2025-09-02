package files

import (
	"account-summary/src/interfaces/libs"
	"fmt"
	"log"
	"os"
)

type localLoader struct{}

func NewLocalLoader() libs.FileLoaderInterface {
	return &localLoader{}
}

func (c *localLoader) LoadFile(path string) (*os.File, error) {
	fileCsv, err := os.Open(path)
	if err != nil {
		log.Default().Printf("Error to open csv file: %v\n", err)
		return nil, fmt.Errorf("error to open csv file: %v", err)
	}

	return fileCsv, nil
}
