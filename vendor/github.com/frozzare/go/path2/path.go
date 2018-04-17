package path2

import (
	"os"

	"github.com/kardianos/osext"
)

// RealPath will return right path for file, looks at the
// given file first and then looks in the executable folder.
// If no file or directory is found it will return the input value.
func RealPath(file string) string {
	if _, err := os.Stat(file); err == nil {
		return file
	}

	path, _ := osext.ExecutableFolder()

	if _, err := os.Stat(path + "/" + file); os.IsNotExist(err) {
		return file
	}

	return path + "/" + file
}
