package input

import (
	"io/fs"
	"path/filepath"
	"strings"
)

// ScanGoFiles walks the directory tree rooted at root and returns a slice
// with the paths of all files that end with ".go".
//
// It skips common vendored and hidden directories (those starting with '.',
// also "vendor" and "node_modules") to avoid scanning unnecessary files.
func ScanGoFiles(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// propagate errors from WalkDir
			return err
		}

		// If it's a directory, decide whether to skip it.
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return fs.SkipDir
			}
			return nil
		}

		// If it's a file and has .go extension, keep it.
		if strings.HasSuffix(d.Name(), ".go") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
