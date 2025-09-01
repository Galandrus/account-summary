package libs

import "os"

type FileLoaderInterface interface {
	LoadFile(path string) (*os.File, error)
}
