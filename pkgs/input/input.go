package input

import (
	"os"
	"path/filepath"
	"strings"
)

type Tree struct {
	Files []string
}

func CreateTree(path string) *Tree {
	tree := &Tree{}
	tree.mapFiles(path)
	return tree
}

func (t *Tree) mapFiles(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic("error mapping root")
	}

	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())

		if f.IsDir() {
			deepPath := fullPath
			t.mapFiles(deepPath)
		} else {
			if strings.HasSuffix(f.Name(), ".go") {
				t.Files = append(t.Files, fullPath)
			}
		}
	}
}
