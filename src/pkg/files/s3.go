package files

import (
	"account-summary/src/interfaces/libs"
	"os"
)

type s3Loader struct{}

func NewS3Loader() libs.FileLoaderInterface {
	return &s3Loader{}
}

func (c *s3Loader) LoadFile(path string) (*os.File, error) {
	//TODO: Implement S3 loading
	return nil, nil
}
